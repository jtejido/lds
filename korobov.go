package lds

import (
	"github.com/jtejido/grand"
)

// The parameters are the modulus n and the multiplier a,
// for an arbitrary integer 1 <= a < n. When a is outside the
// interval [1,n), then we replace a by (a mod n) in all calculations. The
// number of points is n, their dimension is s, and they are defined by:
//
// ui= (i/n)(1, a, a2, , as−1) mod 1 for i=0,…,n-1
// N. M. Korobov. The  Approximate  Computation  of  Multiple  Integrals. 1959.
// I. H. Sloan and S. Joe. Lattice Methods for Multiple Integration. 1994.
// P. L’Ecuyer, C. Lemieux. Variance Reduction via Lattice Rules. 2000.
// V. Sinescu, P. L’Ecuyer. Variance Bounds and Existence Results for Randomly Shifted LatticeRules. 2012.
type KorobovLattice struct {
	Rank1Lattice
	genA, genT uint32
}

// Instantiates a Korobov lattice point set with modulus n and multiplier a in dimension s.
// Generator panics when s is exhausted.
func NewKorobovLattice(n, a, s int, src grand.Source) (*KorobovLattice, error) {
	return NewShiftedKorobovLattice(n, a, s, 0, src)
}

// Instantiates a shifted Korobov lattice point set with modulus
// n and multiplier a in dimension s. The first t coordinates of a standard
// Korobov lattice are dropped as described above. The case t=0 corresponds
// to the standard Korobov lattice.
func NewShiftedKorobovLattice(n, a, s, t int, src grand.Source) (*KorobovLattice, error) {
	if n < 0 {
		panic("number of points cannot be less than 0")
	}

	if s < 0 {
		panic("# of dimension cannot be less than 0")
	}

	if t < 0 {
		panic("shift cannot be less than 0")
	}

	kl := new(KorobovLattice)
	kl.spi = kl
	if src == nil {
		src = internalSrc
	}

	kl.src = src

	kl.genAs = make([]uint32, 0)
	kl.n = n
	kl.genA = uint32(a)
	kl.dim = s
	kl.lv = make([]uint32, s)

	kl.init(n, t)
	kl.Randomize()
	return kl, nil
}

func (kl *KorobovLattice) init(n, t int) {
	a := kl.genA % uint32(n)
	if kl.genA < 0 {
		a += uint32(n)
	}

	kl.genT = uint32(t)
	B := make([]uint32, kl.dim)
	B[0] = 1

	for j := 0; j < t; j++ {
		B[0] *= a
	}

	kl.lv[0] = B[0]

	for j := 1; j < kl.dim; j++ {
		B[j] = (a * B[j-1]) % uint32(n)
		kl.lv[j] = B[j]
	}
}

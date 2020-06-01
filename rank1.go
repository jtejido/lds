package lds

import (
	"github.com/jtejido/grand"
	"math"
)

// Implements point sets specified by integration lattices of rank 1.
// I. H. Sloan and S. Joe. Lattice Methods for Multiple Integration. 1994
type Rank1Lattice struct {
	base
	genAs []uint32
	lv    []uint32
}

// number n points, lattice vector a, dimension s
// One selects an arbitrary positive integer n and an s-dimensional integer
// vector (a_0,…,a_{s-1}). [Usually,a_0=1 and 0 <= a_j < n for each j;
// when the a_j are outside the interval [0,n), then we replace a_j by (a_j mod n) in all calculations.]
//
// The points are defined by :
//
// ui= (i/n)(a0, a1, , as−1) mod 1
// for i=0,…,n-1. These n points are distinct provided that n and the a_j’s have no common factor.
func NewRank1Lattice(n int, a []uint32, s int, src grand.Source) (*Rank1Lattice, error) {
	if n < 0 {
		panic("number of points cannot be less than 0")
	}

	if s < 0 {
		panic("# of dimension cannot be less than 0")
	}

	if a == nil {
		panic("lattice vector cannot be nil")
	}

	if len(a) != s {
		panic("lattice vector length should be equal to dimension")
	}

	if src == nil {
		src = internalSrc
	}

	r1l := new(Rank1Lattice)
	r1l.spi = r1l
	r1l.src = src
	r1l.dim = s
	r1l.lv = make([]uint32, s)
	r1l.genAs = make([]uint32, s)
	for j := 0; j < s; j++ {
		r1l.genAs[j] = a[j]
	}

	r1l.n = n

	for j := 0; j < r1l.dim; j++ {
		amod := (r1l.genAs[j] % uint32(n))
		if r1l.genAs[j] < 0 {
			amod += uint32(n)
		}

		r1l.lv[j] = amod
	}

	r1l.Randomize()
	return r1l, nil
}

func (r1l *Rank1Lattice) addRandomShift(d1, d2 int) {
	if 0 == d2 {
		d2 = int(math.Max(1., float64(r1l.dim)))
	}
	if r1l.shift == nil {
		r1l.shift = make([]uint32, d2)
		r1l.capacityShift = d2
	} else if d2 > r1l.capacityShift {
		d3 := int(math.Max(4., float64(r1l.capacityShift)))
		for d2 > d3 {
			d3 *= 2
		}
		temp := make([]uint32, d3)
		r1l.capacityShift = d3
		for i := 0; i < d1; i++ {
			temp[i] = r1l.shift[i]
		}
		r1l.shift = temp
	}
	r1l.dimShift = d2

	for i := d1; i < d2; i++ {
		r1l.shift[i] = r1l.src.Uint32()
	}
}

func (r1l *Rank1Lattice) Randomize() {
	r1l.addRandomShift(0, r1l.dim)
}

func (r1l *Rank1Lattice) Uint32() uint32 {
	return r1l.uint32()
}

func (r1l *Rank1Lattice) uint32() uint32 {
	if r1l.pointIdx >= r1l.n || r1l.coordIdx >= r1l.dim {
		r1l.outOfBounds()
	}

	x := (uint32(r1l.pointIdx) * r1l.lv[r1l.coordIdx]) % 1

	if r1l.shift != nil {
		if r1l.coordIdx >= r1l.dimShift {
			r1l.addRandomShift(r1l.dimShift, r1l.coordIdx+1)
		}

		x += r1l.shift[r1l.coordIdx]
	}
	r1l.coordIdx++
	return x
}

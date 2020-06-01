package lds

import (
	"fmt"
	"github.com/jtejido/grand"
)

const (
	MAXDIM_NX = 32
)

// Implements digital sequences based on the Niederreiter-Xing
// sequence in base 2.
//
// H. Niederreiter and C. Xing. Nets, (t,s)-Sequences, and Algebraic Geometry. 1998
type NiedXing struct {
	digitalSequenceBaseTwo
}

func NewNiedXing(k, dim int, src grand.Source, t RandomType) (*NiedXing, error) {

	if (dim < 1) || (dim > MAXDIM_NX) {
		return nil, fmt.Errorf("Dimension for NiedXingSequenceBase2 must be > 1 and <= %v", MAXDIM_NX)
	}
	if k >= MAXBITS {
		return nil, fmt.Errorf("k should be less than 32")
	}

	if src == nil {
		src = internalSrc
	}

	nx := new(NiedXing)
	nx.spi = nx
	nx.src = src
	nx.nCols = uint(k)
	nx.n = int(1 << k)
	nx.dim = dim

	nx.vec = make([]uint32, dim*k)
	nx.cachedCurrentPoints = make([]uint32, nx.dim+1)

	var start int

	if dim <= 4 {
		start = 0
	} else {
		start = ((nx.dim * (nx.dim - 1) / 2) - 6) * NUMCOLS
	}

	for j := 0; j < nx.dim; j++ {
		for c := 0; c < k; c++ {
			x := niedXing[start+j*NUMCOLS+c]
			x <<= 1
			nx.vec[j*k+c] = x
		}
	}

	nx.RandomType = t
	nx.Randomize()
	nx.addShiftToCache()
	return nx, nil
}

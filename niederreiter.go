package lds

import (
	"errors"
	"fmt"
	"github.com/jtejido/grand"
)

const (
	NUMCOLS     = 30
	MAXDIM_NIED = 318
	MAXBITS     = 32
)

// Implements digital sequences based on the Niederreiter-Xing sequence in base 2.
// H. Niederreiter and C. Xing. Nets, (t,s)-Sequences, and Algebraic Geometry. 1998
type Niederreiter struct {
	digitalSequenceBaseTwo
}

func NewNiederreiter(k, dim int, src grand.Source, t RandomType) (*Niederreiter, error) {
	if (dim < 1) || (dim > MAXDIM_NIED) {
		return nil, errors.New(fmt.Sprintf("Dimension for NiederreiterSequence must be > 1 and <= %v", MAXDIM_NIED))
	}
	if k >= 32 {
		return nil, errors.New(fmt.Sprintf("k should be less than 32"))
	}
	if src == nil {
		src = internalSrc
	}

	nx := new(Niederreiter)
	nx.spi = nx
	nx.src = src
	nx.nCols = uint(k)
	nx.n = int(1 << k)
	nx.dim = dim

	nx.vec = make([]uint32, dim*k)
	nx.cachedCurrentPoints = make([]uint32, nx.dim+1)

	for j := 0; j < nx.dim; j++ {
		for c := 0; c < k; c++ {
			nx.vec[j*k+c] = niedMat[j*NUMCOLS+c] << 1
			//nx.vec[j*k+c] >>= MAXBITS
		}
	}
	nx.RandomType = t
	nx.Randomize()
	nx.addShiftToCache()
	return nx, nil
}

package lds

import (
	"fmt"
	"github.com/jtejido/grand"
)

const (
	MAXDEGREE_SOBOL = 18
)

// https://en.wikipedia.org/wiki/Sobol_sequence
// 32-bit Sobol sequence. After 2^32 points, the sequence terminates.
type Sobol struct {
	digitalSequenceBaseTwo
}

// number of points 2^k, dimension dim
func NewSobol(k, dim int, src grand.Source, t RandomType) (*Sobol, error) {
	if (dim < 1) || (dim > len(poly)) {
		return nil, fmt.Errorf("Dimension for SobolSequence must be > 0 and <= %v", len(poly))
	}

	if src == nil {
		src = internalSrc
	}

	s := new(Sobol)
	s.spi = s
	s.src = src
	s.nCols = uint(k)
	s.n = int(1 << k)
	s.dim = dim
	s.vec = make([]uint32, dim*k)
	s.cachedCurrentPoints = make([]uint32, s.dim+1)

	// the first dimension, j = 0.
	for c := 0; c < k; c++ {
		s.vec[c] = (1 << (k - c - 1))
	}

	// the other dimensions j > 0.
	for j := 1; j < s.dim; j++ {
		// if a direction number file was provided, use it
		polynomial := poly[j]
		// find the degree of primitive polynomial f_j

		degree := MAXDEGREE_SOBOL
		for ((polynomial >> (degree)) & 1) == 0 {
			degree--
		}
		// Get initial direction numbers m_{j,0},..., m_{j,degree-1}.
		start := j * k
		for c := 0; c < degree && c < k; c++ {
			m_i := minit[j-1][c]
			s.vec[start+c] = m_i << (MAXBITS - c - 1)
		}

		// Compute the following ones via the recursion.
		for c := degree; c < k; c++ {
			nextCol := s.vec[start+c-degree] >> degree
			for i := 0; i < degree; i++ {
				if ((polynomial >> i) & 1) == 1 {
					nextCol ^= s.vec[start+c-degree+i]
				}
			}
			s.vec[start+c] = nextCol
		}
	}
	s.RandomType = t
	s.Randomize()
	s.addShiftToCache()
	return s, nil
}

package lds

import (
	"github.com/jtejido/grand"
	"math"
)

type digitalSequenceBaseTwo struct {
	digitalNetBase
	origVec, vec, digitalShift, cachedCurrentPoints []uint32
}

func (d *digitalSequenceBaseTwo) stripedMatrixScramble() {
	if d.origVec == nil {
		d.origVec = d.vec
		d.vec = make([]uint32, d.dim*int(d.nCols))
	}

	scrambleMat := make([]uint32, MAXBITS)

	var i uint
	for i = 0; i < MAXBITS; i++ {
		scrambleMat[i] = ((1<<32 - 1) >> i) << i
	}
	for j := 0; j < d.dim; j++ {
		d.leftMultiplyMat(j, scrambleMat)
	}
}

func (d *digitalSequenceBaseTwo) leftMatrixScramble() {
	var dd uint
	if d.origVec == nil {
		d.origVec = d.vec
		d.vec = make([]uint32, d.dim*int(d.nCols))
	}

	scrambleMat := make([][]uint32, d.dim)
	for j := 0; j < d.dim; j++ {
		if scrambleMat[j] == nil {
			scrambleMat[j] = make([]uint32, MAXBITS)
		}
		scrambleMat[j][0] = 1<<32 - 1
		for dd = 1; dd < MAXBITS; dd++ {
			scrambleMat[j][dd] = d.src.Uint32() << dd
		}
	}

	for j := 0; int(j) < d.dim; j++ {
		d.leftMultiplyMat(j, scrambleMat[j])
	}
}

func (d *digitalSequenceBaseTwo) leftMultiplyMat(j int, Mj []uint32) {
	var dd, col uint32

	for c := 0; c < int(d.nCols); c++ {
		col = 0

		for dd = 0; dd < MAXBITS; dd++ {
			col ^= (Mj[dd] & d.origVec[j*int(d.nCols)+c]) >> dd
		}

		d.vec[j*int(d.nCols)+c] = col
	}
}

func (d *digitalSequenceBaseTwo) rightMatrixScramble() {
	var c uint
	gen := grand.New(d.src)
	if d.origVec == nil {
		d.origVec = d.vec
		d.vec = make([]uint32, d.dim*int(d.nCols))
	}

	scrambleMat := make([]uint32, MAXBITS)
	var boundInt int32
	for c = 0; c < d.nCols; c++ {
		boundInt += (1 << c)
		scrambleMat[c] = (1 | uint32(gen.Int31n(boundInt))) << (MAXBITS - c - 1)
	}

	for j := 0; j < d.dim; j++ {
		d.rightMultiplyMat(j, scrambleMat)
	}
}

func (d *digitalSequenceBaseTwo) rightMultiplyMat(j int, Mj []uint32) {
	for c := 0; c < int(d.nCols); c++ {
		mask := uint32(1<<32 - 1)
		col := d.origVec[j*int(d.nCols)+c]
		for r := 0; r < c; r++ {
			// If bit (outDigits - 1 - r) of Mj[c] is 1, add column r
			if (Mj[c] & mask) != 0 {
				col ^= d.origVec[j*int(d.nCols)+r]
			}
			mask >>= 1
		}
		d.vec[j*int(d.nCols)+c] = col
	}
}

func (d *digitalSequenceBaseTwo) Seed(seed int64) {
	d.src.Seed(seed)
	d.Restart()
}

func (d *digitalSequenceBaseTwo) addRandomShift(d1, d2 int) {
	if 0 == d2 {
		d2 = int(math.Max(1., float64(d.dim)))
	}

	if d.digitalShift == nil {
		d.digitalShift = make([]uint32, d2)
		d.capacityShift = d2
	} else if d2 > d.capacityShift {
		d3 := int(math.Max(4., float64(d.capacityShift)))
		for d2 > d3 {
			d3 *= 2
		}
		temp := make([]uint32, d3)
		d.capacityShift = d3
		for i := 0; i < d1; i++ {
			temp[i] = d.digitalShift[i]
		}
		d.digitalShift = temp
	}

	for i := d1; i < d2; i++ {
		d.digitalShift[i] = d.src.Uint32()
	}

	d.dimShift = d2
}

func (d *digitalSequenceBaseTwo) Randomize() {
	if d.RandomType == LeftMatrixScrambling {
		d.leftMatrixScramble()
	} else if d.RandomType == StripedMatrixScrambling {
		d.stripedMatrixScramble()
	} else if d.RandomType == RightMatrixScrambling {
		d.rightMatrixScramble()
	}

	d.addRandomShift(0, d.dim)
}

func (d *digitalSequenceBaseTwo) Restart() {
	d.resetCurPointIndex()
}

func (d *digitalSequenceBaseTwo) resetCurPointIndex() {
	d.addShiftToCache()
	d.pointIdx = 0
	d.coordIdx = 0
	d.resetState()
}

func (d *digitalSequenceBaseTwo) Jump() {
	d.resetToNextPoint()
}

func (d *digitalSequenceBaseTwo) resetToNextPoint() int {
	pos := 0 // Will be position of change in Gray code,
	// = pos. of first 0 in binary code of point index.
	for ((d.pointIdx >> pos) & 1) != 0 {
		pos++
	}
	if pos < int(d.nCols) {
		for j := 0; j < d.dim; j++ {
			d.cachedCurrentPoints[j] ^= d.vec[j*int(d.nCols)+pos]
		}
	}
	d.coordIdx = 0
	d.pointIdx++
	d.resetState()
	return d.pointIdx
}

func (d *digitalSequenceBaseTwo) setCurPointIndex(i int) {
	if i == 0 {
		d.resetCurPointIndex()
		return
	}
	// Out of order computation, must recompute the cached current
	// point from scratch.
	d.pointIdx = i
	d.coordIdx = 0
	d.resetState()
	d.addShiftToCache()

	grayCode := i ^ (i >> 1)
	pos := 0 // Position of the bit that is examined.
	for (grayCode >> pos) != 0 {
		if ((grayCode >> pos) & 1) != 0 {
			for j := 0; j < d.dim; j++ {
				d.cachedCurrentPoints[j] ^= d.vec[j*int(d.nCols)+pos]
			}
		}
		pos++
	}
}

func (d *digitalSequenceBaseTwo) addShiftToCache() {
	if d.digitalShift == nil {
		for j := 0; j < d.dim; j++ {
			d.cachedCurrentPoints[j] = 0
		}
	} else {
		if d.dimShift < d.dim {
			d.addRandomShift(d.dimShift, d.dim)
		}
		for j := 0; j < d.dim; j++ {
			d.cachedCurrentPoints[j] = d.digitalShift[j]
		}
	}
}

func (d *digitalSequenceBaseTwo) Uint32() uint32 {
	return d.uint32()
}

func (d *digitalSequenceBaseTwo) uint32() uint32 {
	if d.pointIdx >= d.n || d.coordIdx >= d.dim {
		d.outOfBounds()
	}

	t := d.coordIdx
	d.coordIdx++
	return d.cachedCurrentPoints[t]
}

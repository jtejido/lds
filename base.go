package lds

import (
	"github.com/jtejido/grand"
)

type RandomType int

const (
	// Default Random type
	RandomShift RandomType = iota

	// H. Faure and S. Tezuka. Another Random Scrambling of Digital $(t,s)$-sequences. 2002.
	LeftMatrixScrambling

	// H. Faure and S. Tezuka. Another Random Scrambling of Digital $(t,s)$-sequences. 2002.
	// H. S. Hong and F. H. Hickernell. Algorithm 823: Implementing Scrambled Digital Sequences. 2003
	RightMatrixScrambling

	// A. B. Owen. Variance with Alternative Scramblings of Digital Nets. 2003.
	StripedMatrixScrambling
)

type base struct {
	baseSource32
	dimShift, capacityShift, dim, n, pointIdx, coordIdx int
	shift                                               []uint32
	src                                                 grand.Source
}

func (b *base) Points() int {
	return b.n
}

func (b *base) Dimensions() int {
	return b.dim
}

func (b *base) Seed(seed int64) {
	b.src.Seed(seed)
	b.Restart()
}

func (b base) outOfBounds() {
	if b.pointIdx >= b.n {
		panic("Not enough points available")
	}

	panic("Not enough coordinates available")
}

func (b *base) Restart() {
	b.resetCurPointIndex()
}

func (b *base) RestartSubstream() {
	b.resetCurCoordIndex()
}

func (b *base) Jump() {
	b.resetToNextPoint()
}

func (b *base) resetCurPointIndex() {
	b.setCurPointIndex(0)
}

func (b *base) resetToNextPoint() int {
	b.setCurPointIndex(b.pointIdx + 1)
	return b.pointIdx
}

func (b *base) setCurPointIndex(i int) {
	b.pointIdx = i
	b.resetCurCoordIndex()
}

func (b *base) resetCurCoordIndex() {
	b.setCurCoordIndex(0)
	b.resetState()
}

func (b *base) setCurCoordIndex(j int) {
	b.coordIdx = j
}

type digitalNetBase struct {
	base
	nCols, nRows, b uint
	RandomType      RandomType
}

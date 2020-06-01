package lds

import (
	"github.com/jtejido/grand"
	"github.com/jtejido/grand/source64"
)

// This is used for seeding any implementations here
var (
	internalSrc grand.Source
)

// Internal
// All parent structs should follow this.
type source32 interface {
	Uint32() uint32
}

func init() {
	internalSrc = source64.NewSplitMix64(12345)
}

// This serves as the base struct for holding the starting point for the stream.
// all structs embedding this should handle their own starting point of the stream.
type baseSource32 struct {
	// Golang doesn't override or pickup methods from child (in OOP sense),
	// thus it is required to assign any source32 implementors here to be used by Bool().
	spi source32
	// booleanSource caches the most recent uint32
	// booleanBitMask is the bit mask of the boolean source to obtain the boolean bit.
	// This begins at the least significant bit and is gradually shifted upwards until overflow to zero.
	// When zero a new boolean source should be created and the mask set to the least significant bit (i.e. 1).
	booleanBitMask, booleanSource uint32
	// stream stores the starting point for the stream.
	stream []uint32
}

func (bs32 *baseSource32) Uint32() uint32 { return bs32.spi.Uint32() }

// restores starting values for Bool()
func (bs32 *baseSource32) resetState() {
	bs32.booleanSource = 0
	bs32.booleanBitMask = 0
}

// Generates a boolean value
func (bs32 *baseSource32) Bool() bool {
	bs32.booleanBitMask <<= 1
	if bs32.booleanBitMask == 0 {
		bs32.booleanBitMask = 1
		bs32.booleanSource = bs32.Uint32()
	}

	return (bs32.booleanSource & bs32.booleanBitMask) != 0
}

// Embeds baseSource32 and store substreams.
type baseJumpableSource32 struct {
	baseSource32
	//  substream stores the starting point of the current substream.
	substream []uint32
}

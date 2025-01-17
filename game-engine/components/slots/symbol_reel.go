package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// NewMultiSymbolReels instantiates a new set of multiple symbol reel sets.
// This can be used to randomly alternate between different symbol reel sets.
func NewMultiSymbolReels(sets ...*SymbolReels) *MultiSymbolReels {
	return &MultiSymbolReels{flag: -1, sets: sets}
}

// WithFlag initializes the flag to report the use of this set.
func (msr *MultiSymbolReels) WithFlag(flag int) *MultiSymbolReels {
	msr.flag = flag
	return msr
}

// Sets returns the configured symbol reel sets.
func (msr *MultiSymbolReels) Sets() []*SymbolReels {
	return msr.sets
}

// Spin can be used to fill all reels on the grid with random symbols from the configured symbol reels.
// This function implements the Spinner interface.
func (msr *MultiSymbolReels) Spin(spin *Spin, out utils.Indexes) {
	set := spin.prng.IntN(len(msr.sets))
	if msr.flag >= 0 {
		spin.roundFlags[msr.flag] = set
	}
	msr.sets[set].Spin(spin, out)
}

// NewSymbolReels instantiates a new set of symbol reels.
// This can be used to fill an entire grid using the configured symbol reels.
func NewSymbolReels(reels ...*SymbolReel) *SymbolReels {
	return &SymbolReels{flag: -1, reels: reels}
}

// WithFlag initializes the flag to report the use of this set.
func (sr *SymbolReels) WithFlag(flag int) *SymbolReels {
	sr.flag = flag
	return sr
}

// Reel returns one of the configured symbol reels.
// Note that reels are 1-based, so reel 0 does not exist and will throw a panic)
func (sr *SymbolReels) Reel(reel uint8) utils.Indexes {
	return sr.reels[reel-1].Reel()
}

// Spin can be used to fill all reels on the grid with random symbols from the configured symbol reels.
// This function implements the Spinner interface.
func (sr *SymbolReels) Spin(spin *Spin, out utils.Indexes) {
	rows, mask := spin.rowCount, spin.mask

	if sr.flag >= 0 {
		spin.roundFlags[sr.flag] = 1
	}

	var offs int
	for ix := range sr.reels {
		reel := sr.reels[ix]
		end := offs + int(mask[ix])
		reel.Fill(spin, out[offs:end])
		offs += rows
	}
}

// NewSymbolReel instantiates a new symbol reel.
// This can be used to fill a single reel on the grid with random symbols.
func NewSymbolReel(rows uint8, stops ...utils.Index) *SymbolReel {
	rs := &SymbolReel{rows: rows, max: len(stops), flag: -1}
	rs.stops = make(utils.Indexes, rs.max+int(rows)-1)

	copy(rs.stops, stops)
	copy(rs.stops[rs.max:], stops) // accommodate circular references.

	return rs
}

// WithFlag initializes the flag to report the use of this symbol reel.
func (sr *SymbolReel) WithFlag(flag int) *SymbolReel {
	sr.flag = flag
	return sr
}

// Reel returns the configured symbol reel.
func (sr *SymbolReel) Reel() utils.Indexes {
	return sr.stops[:sr.max]
}

// Fill can be used to fill a reel on the grid with random symbols from the symbol reel.
func (sr *SymbolReel) Fill(spin *Spin, out utils.Indexes) {
	stop := spin.prng.IntN(sr.max)
	if sr.flag >= 0 {
		spin.roundFlags[sr.flag] = stop
	}
	copy(out, sr.stops[stop:])
}

// MultiSymbolReels represents a set of multiple symbol reel sets.
type MultiSymbolReels struct {
	flag int
	sets []*SymbolReels
}

// SymbolReels represents a set of symbol reels.
type SymbolReels struct {
	flag  int
	reels []*SymbolReel
}

// SymbolReel represents a symbol reel which can be used to fill a grid reel with random symbols.
type SymbolReel struct {
	rows  uint8
	flag  int
	max   int
	stops utils.Indexes
}

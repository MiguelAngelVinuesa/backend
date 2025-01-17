package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/interfaces"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// AcquireReel instantiates a new slot machine reel from the memory pool.
// noRepeat indicates that randomly generated indexes do not repeat consecutively.
// e.g. after a reel "spin" consecutive rows are guaranteed not to contain the same symbol if noRepeat == true.
// The function panics if symbols is nil.
// Reels are not safe to be used across concurrent go-routines.
func AcquireReel(index, rows uint8, symbols *SymbolSet, noRepeat uint8) *Reel {
	r := reelProducer.Acquire().(*Reel)
	r.index = index
	r.rows = int(rows)
	r.symbols = symbols

	switch {
	case noRepeat == 0:
		r.weighting = utils.AcquireWeighting()
	case rows == 3 && noRepeat == 2:
		r.weighting = utils.AcquireWeightingDedup3()
	case rows == 4 && noRepeat == 3:
		r.weighting = utils.AcquireWeightingDedup4()
	case rows == 5 && noRepeat == 4:
		r.weighting = utils.AcquireWeightingDedup5()
	default:
		panic("AcquireReel: invalid noRepeat config")
	}

	for _, symbol := range r.symbols.symbols {
		if int(r.index) < len(symbol.weights) {
			r.weighting.AddWeight(symbol.id, symbol.weights[r.index])
		}
	}

	return r
}

// Spin performs a reel "spin". It fills the given slice with random symbol numbers generated through
// the weightings using the given PRNG. If len(in) > rowCount, only rowCount positions are filled from index 0 onwards.
func (r *Reel) Spin(prng interfaces.Generator, in utils.Indexes) {
	if l := len(in); r.rows > l {
		r.weighting.FillRandom(prng, l, in)
	} else {
		r.weighting.FillRandom(prng, r.rows, in)
	}
}

// RandomIndex retrieves a single random symbol index.
func (r *Reel) RandomIndex(prng interfaces.Generator) utils.Index {
	return r.weighting.RandomIndex(prng)
}

// Reel contains the details of a slot machine reel.
// Keep fields ordered by ascending SizeOf().
type Reel struct {
	index     uint8
	rows      int
	symbols   *SymbolSet
	weighting utils.WeightedGenerator
	pool.Object
}

// reelProducer is the memory pool for slot machine reels.
// Make sure to initialize all slices appropriately!
var reelProducer = pool.NewProducer(func() (pool.Objecter, func()) {
	r := &Reel{}
	return r, r.reset
})

// reset clears trhe reel.
func (r *Reel) reset() {
	if r != nil {
		r.index = 0
		r.rows = 0
		r.symbols = nil

		if r.weighting != nil {
			r.weighting.Release()
			r.weighting = nil
		}
	}
}

// Reels is a convenience type for a slice of reels.
type Reels []*Reel

// ReleaseReels releases all reels and returns an empty slice.
func ReleaseReels(list Reels) Reels {
	if list == nil {
		return nil
	}
	for ix := range list {
		if r := list[ix]; r != nil {
			r.Release()
			list[ix] = nil
		}
	}
	return list[:0]
}

// PurgeReels returns the input cleared to zero length or a new slice if it doesn't have the requested capacity.
func PurgeReels(list Reels, capacity int) Reels {
	list = ReleaseReels(list)
	if cap(list) < capacity {
		return make(Reels, 0, capacity)
	}
	return list
}

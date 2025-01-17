package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func (a *ReviseAction) doDeduplication(spin *Spin) bool {
	l := len(spin.Indexes())
	maxID := spin.symbols.maxID
	max := uint8(a.dupWeights.RandomIndex(spin.prng))

	// for a second spin we increase max with the count of existing sticky symbols.
	if spin.Kind() == SecondSpin || spin.Kind() == SecondFreeSpin {
		s := spin.StickyCount()
		max += s

		// prevent against single (or no) sticky symbol selected!
		switch s {
		case 0:
			max += 2
		case 1:
			max += 1
		}
	}

	if max == 0 || (max == 1 && len(spin.symbols.symbols) < l) {
		panic(consts.MsgInvalidDeduplication)
	}

	// more than 100 symbols is insane! this way it's on the stack, not the heap!
	counts := make([]uint8, 100)
	for _, id := range spin.indexes {
		counts[id] = counts[id] + 1
	}

	valid := func(id1, id2 utils.Index) bool {
		if id1 == id2 || counts[id1] >= max {
			return false
		}
		if s := spin.GetSymbol(id1); s == nil || s.kind != Standard {
			return false
		}
		return true
	}

	var dup bool
	for ix, c := range counts {
		for c > max {
			dup = true

			// pick a random symbol to replace the duplicate.
			id := utils.Index(spin.prng.IntN(int(maxID)) + 1)

			// make sure it's not the same as we're replacing, has a low count, exists and is a standard symbol!
			for !valid(id, utils.Index(ix)) {
				if id >= maxID {
					id = 1
				} else {
					id++
				}
			}

			// locate, at random, one of the duplicates.
			offset := spin.prng.IntN(l)
			for spin.indexes[offset] != utils.Index(ix) || spin.sticky[offset] {
				if offset++; offset >= l {
					offset = 0
				}
			}

			// replace it.
			spin.indexes[offset] = id
			counts[id] = counts[id] + 1
			c--
		}
	}

	return dup
}

func (a *ReviseAction) doDedupSymbol(spin *Spin) bool {
	var dup bool

	f := func(reel int) {
		var found bool
		min := reel * spin.rowCount
		max := min + spin.rowCount
		for offset := min; offset < max; offset++ {
			if !found {
				found = spin.indexes[offset] == a.symbol
			} else {
				// here we need to de-duplicate it.
				for spin.indexes[offset] == a.symbol {
					dup = true
					spin.indexes[offset] = spin.reels[reel].weighting.RandomIndex(spin.prng)
				}
			}
		}
	}

	if len(a.dedupReels) > 0 {
		for ix := range a.dedupReels {
			f(int(a.dedupReels[ix]) - 1)
		}
	} else {
		for reel := 0; reel < spin.reelCount; reel++ {
			f(reel)
		}
	}

	return dup
}

package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// OnGridContains returns a filter to test if a symbol exists in the grid.
// Note that reels are 1-based!
func OnGridContains(symbol utils.Index, reels ...int) SpinDataFilterer {
	return func(spin *Spin) bool {
		if len(reels) == 0 {
			for offset := range spin.indexes {
				if spin.indexes[offset] == symbol {
					return true
				}
			}
		} else {
			for ix := range reels {
				reel := reels[ix] - 1
				offset := reel * spin.rowCount
				high := offset + int(spin.mask[reel])
				for ; offset < high; offset++ {
					if spin.indexes[offset] == symbol {
						return true
					}
				}
			}
		}

		return false
	}
}

// OnGridNotContains returns a filter to test if a symbol does not exist in the grid.
// Note that reels are 1-based!
func OnGridNotContains(symbol utils.Index, reels ...int) SpinDataFilterer {
	return func(spin *Spin) bool {
		if len(reels) == 0 {
			for offset := range spin.indexes {
				if spin.indexes[offset] == symbol {
					return false
				}
			}
		} else {
			for ix := range reels {
				reel := reels[ix] - 1
				offset := reel * spin.rowCount
				high := offset + int(spin.mask[reel])
				for ; offset < high; offset++ {
					if spin.indexes[offset] == symbol {
						return false
					}
				}
			}
		}

		return true
	}
}

// OnGridCount returns a filter to test if a symbol appears count times in the grid.
// Note that reels are 1-based!
func OnGridCount(symbol utils.Index, count int, reels ...int) SpinDataFilterer {
	return func(spin *Spin) bool {
		var c int
		if len(reels) == 0 {
			for offset := range spin.indexes {
				if spin.indexes[offset] == symbol {
					c++
				}
			}
		} else {
			for ix := range reels {
				reel := reels[ix] - 1
				offset := reel * spin.rowCount
				high := offset + int(spin.mask[reel])
				for ; offset < high; offset++ {
					if spin.indexes[offset] == symbol {
						c++
					}
				}
			}
		}
		return c == count
	}
}

// OnGridCounts returns a filter to test if a symbol appears one of the counts times in the grid.
// Note that reels are 1-based!
func OnGridCounts(symbol utils.Index, counts []int, reels ...int) SpinDataFilterer {
	return func(spin *Spin) bool {
		var c int
		if len(reels) == 0 {
			for offset := range spin.indexes {
				if spin.indexes[offset] == symbol {
					c++
				}
			}
		} else {
			for ix := range reels {
				reel := reels[ix] - 1
				offset := reel * spin.rowCount
				high := offset + int(spin.mask[reel])
				for ; offset < high; offset++ {
					if spin.indexes[offset] == symbol {
						c++
					}
				}
			}
		}

		for ix := range counts {
			if c == counts[ix] {
				return true
			}
		}
		return false
	}
}

// OnGridShape returns a filter to test if a symbol appears on all the given shape coordinates in the grid.
func OnGridShape(symbol utils.Index, shape utils.UInt8s) SpinDataFilterer {
	return func(spin *Spin) bool {
		for ix := range shape {
			if spin.indexes[shape[ix]] != symbol {
				return false
			}
		}
		return true
	}
}

// OnGridNoWilds return true if the spin grid contains no wilds.
func OnGridNoWilds() SpinDataFilterer {
	return func(spin *Spin) bool {
		return spin.newWilds == 0
	}
}

// OnGridNoScatters return true if the spin grid contains no scatters.
func OnGridNoScatters() SpinDataFilterer {
	return func(spin *Spin) bool {
		return spin.newScatters == 0
	}
}

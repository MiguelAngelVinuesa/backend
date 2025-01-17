package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// PreventPaylines removes possible paylines by de-duplicating symbols and (optionally) removing wilds from the first adjoining reels.
// The direction must be equal to how the paylines are awarded.
// removeWilds indicate if wild symbols should also be removed from the second reel.
// allowDupes indicate if duplicate symbols are allowed on the same reel.
func (s *Spin) PreventPaylines(direction PayDirection, removeWilds, allowDupes bool) bool {
	switch direction {
	case PayLTR:
		return s.doPreventPaylines(removeWilds, allowDupes, 0, 1)
	case PayRTL:
		return s.doPreventPaylines(removeWilds, allowDupes, s.reelCount-1, s.reelCount-2)
	case PayBoth:
		return s.doPreventPaylines(removeWilds, allowDupes, 0, 1) || s.doPreventPaylines(removeWilds, allowDupes, s.reelCount-1, s.reelCount-2)
	default:
		return false
	}
}

// PreventPaylines2 removes possible paylines by de-duplicating symbols and (optionally) removing wilds from the 2nd & 3rd reels.
// The direction must be equal to how the paylines are awarded.
// removeWilds indicate if wild symbols should also be removed from the second reel.
// allowDupes indicate if duplicate symbols are allowed on the same reel.
func (s *Spin) PreventPaylines2(direction PayDirection, removeWilds, allowDupes bool) bool {
	switch direction {
	case PayLTR:
		return s.doPreventPaylines(removeWilds, allowDupes, 1, 2)
	case PayRTL:
		return s.doPreventPaylines(removeWilds, allowDupes, s.reelCount-2, s.reelCount-3)
	case PayBoth:
		return s.doPreventPaylines(removeWilds, allowDupes, 1, 2) || s.doPreventPaylines(removeWilds, allowDupes, s.reelCount-2, s.reelCount-3)
	default:
		return false
	}
}

func (s *Spin) doPreventPaylines(removeWilds, allowDupes bool, first, second int) bool {
	reels := s.reels
	if s.altActive {
		reels = s.altReels
	}

	rows := s.rowCount

	if !removeWilds || s.HasSticky() {
		if first < second {
			for reel := first; reel <= second; reel++ {
				offset := reel * rows
				for row := 0; row < rows; row++ {
					if s.sticky[offset] {
						return false // no removal if there are stickies in the reels!
					}
					if !removeWilds {
						if id := s.indexes[offset]; id != utils.NullIndex {
							if symbol := s.GetSymbol(id); symbol != nil && symbol.isWild {
								return false // no removal if there are wilds in the reels!
							}
						}
					}
					offset++
				}
			}
		} else {
			for reel := first; reel >= second; reel-- {
				offset := reel * rows
				for row := 0; row < rows; row++ {
					if s.sticky[offset] {
						return false // no removal if there are stickies in the reels!
					}
					if !removeWilds {
						if id := s.indexes[offset]; id != utils.NullIndex {
							if symbol := s.GetSymbol(id); symbol != nil && symbol.isWild {
								return false // no removal if there are wilds in the reels!
							}
						}
					}
					offset++
				}
			}

		}
	}

	replaced := make(utils.Indexes, 0, 32)

	var id utils.Index
	validSecond := func() bool {
		offset := second * rows
		if id == utils.NullIndex {
			return true
		}

		symbol := s.GetSymbol(id)
		if symbol == nil {
			return false
		}
		if symbol.isScatter && !symbol.isWild {
			return true
		}

		for _, r := range replaced {
			if r == id {
				return false
			}
		}

		for max := offset + rows; offset < max; offset++ {
			if s.indexes[offset] == id {
				return false
			}
		}

		return !symbol.isWild
	}

	validFirst := func() bool {
		offset := first * rows
		if id == utils.NullIndex {
			return true
		}

		symbol := s.GetSymbol(id)
		if symbol == nil {
			return false
		}
		if symbol.isScatter && !symbol.isWild {
			return true
		}

		for _, r := range replaced {
			if r == id {
				return false
			}
		}

		for max := offset + rows; offset < max; offset++ {
			if s.indexes[offset] == id {
				return false
			}
		}

		return !symbol.isWild
	}

	dupeFirst := func() bool {
		if allowDupes {
			return false
		}
		offset := first * rows
		for max := offset + rows; offset < max; offset++ {
			if s.indexes[offset] == id {
				return true
			}
		}
		return false
	}

	dupeSecond := func() bool {
		if allowDupes {
			return false
		}
		offset := second * rows
		for max := offset + rows; offset < max; offset++ {
			if s.indexes[offset] == id {
				return true
			}
		}
		return false
	}

	var removed bool

	offset := first * rows
	for max := offset + rows; offset < max; offset++ {
		replaced = replaced[:0]
		id = s.indexes[offset]
		if !validSecond() {
			for {
				replaced = append(replaced, id)
				id = reels[first].RandomIndex(s.prng)
				if validSecond() && !dupeFirst() {
					break
				}
			}
			s.indexes[offset] = id
			removed = true
		}
	}

	if removeWilds {
		offset = second * rows
		for max := offset + rows; offset < max; offset++ {
			replaced = replaced[:0]
			id = s.indexes[offset]
			if !validFirst() {
				for {
					replaced = append(replaced, id)
					id = reels[second].RandomIndex(s.prng)
					if validFirst() && !dupeSecond() {
						break
					}
				}
				s.indexes[offset] = id
				removed = true
			}
		}
	}

	return removed
}

// PreventBonus removes possible payouts caused by bonus symbols.
// allowDupes indicates if duplicate symbols are allowed on the same reel.
func (s *Spin) PreventBonus(allowDupes bool) bool {
	bonus := s.BonusSymbol()
	if bonus == utils.NullIndex || bonus == utils.MaxIndex {
		return false
	}

	reels := s.reels
	if s.altActive {
		reels = s.altReels
	}

	rows := s.rowCount

	var id utils.Index
	dupe := func(reel int) bool {
		offset := reel * rows
		for max := offset + rows; offset < max; offset++ {
			if s.indexes[offset] == id {
				return true
			}
		}
		return false
	}

	var removed bool
	var offset int

	for reel := 0; reel < s.reelCount; reel++ {
		for row := 0; row < rows; row++ {
			if id = s.indexes[offset]; id == bonus {
				for {
					id = reels[reel].RandomIndex(s.prng)
					if id != bonus && (allowDupes || !dupe(reel)) {
						break
					}
				}
				s.indexes[offset] = id
				removed = true
			}
			offset++
		}
	}

	return removed
}

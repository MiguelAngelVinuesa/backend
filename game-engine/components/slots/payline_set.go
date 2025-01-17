package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// PaylineSet contains all paylines for a slot machine game.
// It is optimized to quickly determine the payouts for a spin.
type PaylineSet struct {
	direction PayDirection
	highest   bool
	paylines  Paylines
	tree      payOffsets
}

// NewPaylineSet instantiates a new set of paylines.
func NewPaylineSet(directions PayDirection, highestPayout bool, paylines ...*Payline) *PaylineSet {
	s := &PaylineSet{direction: directions, highest: highestPayout, paylines: paylines}
	if directions == PayRTL {
		s.buildTreeRTL()
	} else {
		s.buildTreeLTR() // LTR & Both!
	}
	return s
}

// Directions returns the payline direction(s).
func (s *PaylineSet) Directions() PayDirection {
	return s.direction
}

// HighestPayout defines if the highest payout feature is active.
func (s *PaylineSet) HighestPayout() bool {
	return s.highest
}

// Paylines returns the paylines contained within this set.
func (s *PaylineSet) Paylines() Paylines {
	return s.paylines
}

// GetPayouts detects any payouts in the given spin data and adds them to the result.
func (s *PaylineSet) GetPayouts(spin *Spin, result *results.Result) bool {
	mult := len(spin.multipliers) == len(spin.indexes)
	if mult {
		return s.testTreeMult(spin, result)
	}
	return s.testTree(spin, result)
}

func (s *PaylineSet) buildTreeLTR() {
	var offs *payOffset
	for ix := range s.paylines {
		p := s.paylines[ix]
		s.tree, offs = findOffset(s.tree, s.highest, s.direction, p.offsets[0], 1)
		offs.addOffsetLTR(p, 1, 2)
	}
}

func (s *PaylineSet) buildTreeRTL() {
	var offs *payOffset
	for ix := range s.paylines {
		p := s.paylines[ix]
		reel := len(p.offsets) - 1
		s.tree, offs = findOffset(s.tree, s.highest, s.direction, p.offsets[reel], 1)
		offs.addOffsetRTL(p, reel-1, 2)
	}
}

func (s *PaylineSet) testTree(spin *Spin, result *results.Result) bool {
	var found, got bool
	for ix := range s.tree {
		offs := s.tree[ix]
		if symbol := spin.GetSymbol(spin.indexes[offs.offset]); symbol != nil {
			if symbol.isWild {
				got = offs.testBranches(spin, result, symbol, 1, 1, utils.NewMultiplier(symbol.multiplier))
			} else {
				got = offs.testBranches(spin, result, symbol, 0, 0, 1.0)
			}
		}
		found = found || got
	}
	return found
}

func (s *PaylineSet) testTreeMult(spin *Spin, result *results.Result) bool {
	var found, got bool
	for ix := range s.tree {
		offs := s.tree[ix]
		if symbol := spin.GetSymbol(spin.indexes[offs.offset]); symbol != nil {
			m := utils.NewMultiplier(float64(spin.multipliers[offs.offset]))
			if symbol.isWild {
				got = offs.testBranchesMult(spin, result, symbol, 1, 1, m, m)
			} else {
				got = offs.testBranchesMult(spin, result, symbol, 0, 0, m, m)
			}
		}
		found = found || got
	}
	return found
}

func (o *payOffset) testBranches(spin *Spin, result *results.Result, symbol *Symbol, wilds, allWilds int, multiplier float64) bool {
	if len(o.branches) == 0 {
		return o.testPayouts(o, spin, result, symbol, wilds, allWilds, multiplier, 1.0)
	}

	var found, got bool

	for ix := range o.branches {
		offs := o.branches[ix]
		s := spin.GetSymbol(spin.indexes[offs.offset])
		if s == nil {
			got = o.testPayouts(offs, spin, result, symbol, wilds, allWilds, multiplier, 1.0)
		} else {
			switch {
			case s == symbol:
				if symbol.isWild {
					m := utils.NewMultiplier(multiplier, s.multiplier)
					got = offs.testBranches(spin, result, symbol, wilds+1, allWilds+1, m)
				} else {
					got = offs.testBranches(spin, result, symbol, wilds, allWilds, multiplier)
				}

			case s.isWild:
				m := utils.NewMultiplier(multiplier, s.multiplier)
				if symbol.isWild {
					got = offs.testBranches(spin, result, symbol, wilds+1, allWilds+1, m)
				} else {
					got = offs.testBranches(spin, result, symbol, wilds, allWilds+1, m)
				}

			case s.WildFor(symbol.id):
				m := utils.NewMultiplier(multiplier, s.multiplier)
				got = offs.testBranches(spin, result, symbol, wilds, allWilds, m)

			case symbol.isWild, symbol.WildFor(s.ID()):
				got = offs.testBranches(spin, result, s, wilds, allWilds, multiplier)

			default:
				got = o.testPayouts(offs, spin, result, symbol, wilds, allWilds, multiplier, 1.0)

			}
		}

		if got {
			found = true
		}
	}

	return found
}

func (o *payOffset) testBranchesMult(spin *Spin, result *results.Result, symbol *Symbol, wilds, allWilds int, multiplier, startMultiplier float64) bool {
	if len(o.branches) == 0 {
		return o.testPayouts(o, spin, result, symbol, wilds, allWilds, multiplier, startMultiplier)
	}

	var found, got bool

	for ix := range o.branches {
		offs := o.branches[ix]
		s := spin.GetSymbol(spin.indexes[offs.offset])
		if s == nil {
			got = o.testPayouts(offs, spin, result, symbol, wilds, allWilds, multiplier, startMultiplier)
		} else {
			switch {
			case s == symbol:
				m := utils.NewMultiplier(multiplier, float64(spin.multipliers[offs.offset]))
				if symbol.isWild {
					m *= utils.NewMultiplier(s.multiplier)
					got = offs.testBranchesMult(spin, result, symbol, wilds+1, allWilds+1, m, m)
				} else {
					got = offs.testBranchesMult(spin, result, symbol, wilds, allWilds, m, startMultiplier)
				}

			case s.isWild:
				m := utils.NewMultiplier(multiplier, float64(spin.multipliers[offs.offset]), s.multiplier)
				if symbol.isWild {
					got = offs.testBranchesMult(spin, result, symbol, wilds+1, allWilds+1, m, m)
				} else {
					got = offs.testBranchesMult(spin, result, symbol, wilds, allWilds+1, m, startMultiplier)
				}

			case s.WildFor(symbol.id):
				m := utils.NewMultiplier(multiplier, float64(spin.multipliers[offs.offset]), s.multiplier)
				got = offs.testBranchesMult(spin, result, symbol, wilds, allWilds, m, startMultiplier)

			case symbol.isWild, symbol.WildFor(s.ID()):
				m := utils.NewMultiplier(multiplier, float64(spin.multipliers[offs.offset]))
				got = offs.testBranchesMult(spin, result, s, wilds, allWilds, m, startMultiplier)

			default:
				got = o.testPayouts(offs, spin, result, symbol, wilds, allWilds, multiplier, startMultiplier)

			}
		}

		if got {
			found = true
		}
	}

	return found
}

func (o *payOffset) testPayouts(pay *payOffset, spin *Spin, result *results.Result, symbol *Symbol, wilds, allWilds int, multiplier, wildMultiplier float64) bool {
	symbolID := symbol.id
	count := o.level
	factor := symbol.Payout(count)

	if o.highest && wilds > 1 {
		// highest payout only!
		symbols := spin.GetSymbols()
		if highest := symbols.bestWildPay[wilds]; highest > factor {
			symbolID = symbols.bestWildSym[wilds]
			count = uint8(wilds)
			factor = highest
			multiplier = wildMultiplier
		}
	}

	if factor <= 0 && o.direction != PayBoth {
		return false
	}

	var found bool

	for ix := range pay.paylines {
		p := pay.paylines[ix]
		if o.direction == PayBoth {
			if !paylineDone(p.id, result) {
				// test each payline in reverse to see if it gives a beter result.
				if symbolID2, count2, allWilds2, factor2, multiplier2 := o.testReverse(spin, p); factor2*multiplier2 > factor*multiplier {
					result.AddPayouts(WinlinePayoutFromData(factor2, utils.NewMultiplier(multiplier2, spin.getMultiplier(allWilds2)), symbolID2, count2, PayRTL, p.id, p.rows))
					spin.markWinline(count2, PayRTL, p.rows)
					found = true
				} else if factor > 0 {
					result.AddPayouts(WinlinePayoutFromData(factor, utils.NewMultiplier(multiplier, spin.getMultiplier(allWilds)), symbolID, count, PayLTR, p.id, p.rows))
					spin.markWinline(count, PayLTR, p.rows)
					found = true
				}
			}
		} else {
			result.AddPayouts(WinlinePayoutFromData(factor, utils.NewMultiplier(multiplier, spin.getMultiplier(allWilds)), symbolID, count, o.direction, p.id, p.rows))
			spin.markWinline(count, o.direction, p.rows)
			found = true
		}
	}

	return found
}

func paylineDone(id uint8, result *results.Result) bool {
	for ix := range result.Payouts {
		if p, ok := result.Payouts[ix].(*SpinPayout); ok {
			if p.paylineID == id {
				return true
			}
		}
	}
	return false
}

func (o *payOffset) testReverse(spin *Spin, p *Payline) (utils.Index, uint8, int, float64, float64) {
	reelIX := len(p.offsets) - 1
	symbol := spin.GetSymbol(spin.indexes[p.offsets[reelIX]])
	if symbol == nil {
		return 0, 0, 0, 0, 0
	}

	count := uint8(1)

	var startWilds, allWilds int
	if symbol.isWild {
		startWilds = 1
		allWilds = 1
	}

	for reelIX--; reelIX >= 0; reelIX-- {
		if s := spin.GetSymbol(spin.indexes[p.offsets[reelIX]]); s == nil {
			reelIX = -1
		} else {
			switch {
			case s == symbol:
				if symbol.isWild {
					startWilds++
					allWilds++
				}
				count++

			case s.isWild:
				if symbol.isWild {
					startWilds++
				}
				allWilds++
				count++

			case s.WildFor(symbol.id):
				count++

			case symbol.isWild, symbol.WildFor(s.id):
				symbol = s
				count++

			default:
				reelIX = -1
			}
		}
	}

	symbolID := symbol.id
	factor := symbol.Payout(count)

	if o.highest && startWilds > 1 {
		// highest payout only!
		symbols := spin.GetSymbols()
		if highest := symbols.bestWildPay[startWilds]; highest > factor {
			count = uint8(startWilds)
			allWilds = startWilds
			factor = highest
			symbolID = symbols.bestWildSym[startWilds]
		}
	}

	// we determine the real multiplier after a potential highest payout switch.
	multiplier := 1.0
	if len(spin.multipliers) == len(spin.indexes) {
		reelIX = len(p.offsets) - 1
		for ix := uint8(0); ix < count; ix++ {
			multiplier = utils.NewMultiplier(multiplier, float64(spin.multipliers[p.offsets[reelIX]]))
			reelIX--
		}
	}

	return symbolID, count, allWilds, factor, multiplier
}

func (o *payOffset) addOffsetLTR(p *Payline, reel int, level uint8) {
	o.paylines = append(o.paylines, p)

	var offs *payOffset
	o.branches, offs = findOffset(o.branches, o.highest, o.direction, p.offsets[reel], level)

	if reel < len(p.offsets)-1 {
		offs.addOffsetLTR(p, reel+1, level+1)
	} else {
		offs.paylines = append(offs.paylines, p)
	}
}

func (o *payOffset) addOffsetRTL(p *Payline, reel int, level uint8) {
	o.paylines = append(o.paylines, p)

	var offs *payOffset
	o.branches, offs = findOffset(o.branches, o.highest, o.direction, p.offsets[reel], level)

	if reel > 0 {
		offs.addOffsetRTL(p, reel-1, level+1)
	} else {
		offs.paylines = append(offs.paylines, p)
	}
}

func findOffset(list payOffsets, highest bool, direction PayDirection, offset int, level uint8) (payOffsets, *payOffset) {
	var offs *payOffset

	if list == nil {
		list = make(payOffsets, 0)
	} else {
		for ix := range list {
			offs = list[ix]
			if offs.offset == offset {
				break
			} else {
				offs = nil
			}
		}
	}

	if offs == nil {
		offs = newOffset(highest, direction, offset, level)
		list = append(list, offs)
	}

	return list, offs
}

func newOffset(highest bool, direction PayDirection, offset int, level uint8) *payOffset {
	return &payOffset{
		highest:   highest,
		direction: direction,
		level:     level,
		offset:    offset,
		branches:  make(payOffsets, 0),
		paylines:  make(Paylines, 0),
	}
}

type payOffset struct {
	highest   bool
	direction PayDirection
	level     uint8
	offset    int
	branches  payOffsets
	paylines  Paylines
}

type payOffsets []*payOffset

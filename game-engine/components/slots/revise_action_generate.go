package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func (a *ReviseAction) doGenerateSymbol(spin *Spin) bool {
	l := len(spin.indexes)
	if spin.StickyCount() > uint8(l/2) {
		return false // not enough space left!
	}

	if !a.genAllowOld {
		for ix := range spin.indexes {
			if spin.indexes[ix] == a.symbol && !spin.sticky[ix] {
				return false // already done!
			}
		}
	}

	return a.doGenerateSymbol2(spin)
}

type validator struct {
	testMask   bool
	testReels  bool
	testDupes1 bool
	testDupes2 bool
	gridSize   int
	rows       int
	spin       *Spin
	a          *ReviseAction
}

func (v validator) valid(offs int) bool {
	if v.spin.indexes[offs] == v.a.symbol || v.spin.sticky[offs] {
		return false
	}
	if v.testMask && !v.validMask(offs) {
		return false
	}
	if v.testReels && !v.validReel(offs) {
		return false
	}
	if v.testDupes1 && !v.validNoDupes1(offs) {
		return false
	}
	if v.testDupes2 && !v.validNoDupes2(offs) {
		return false
	}
	return true
}

func (v validator) validMask(offs int) bool {
	reel := offs / v.rows
	row := uint8(offs - (reel * v.rows))
	return row < v.spin.mask[reel]
}

func (v validator) validReel(offs int) bool {
	reel := uint8(offs/v.rows) + 1
	return v.a.generateReels.Contains(reel)
}

func (v validator) validNoDupes1(offs int) bool {
	offs = (offs / v.rows) * v.rows
	end := offs + v.rows
	for offs < end {
		if v.spin.indexes[offs] == v.a.symbol {
			return false
		}
		offs++
	}
	return true
}

func (v validator) validNoDupes2(offs int) bool {
	offs = (offs / v.rows) * v.rows
	end := offs + v.rows
	for offs < end {
		s := v.spin.indexes[offs]
		if s == v.a.symbol {
			return false
		}
		for ix := range v.a.morphFor {
			if s == v.a.morphFor[ix] {
				return false
			}
		}
		offs++
	}
	return true
}

func (a *ReviseAction) doGenerateSymbol2(spin *Spin) bool {
	v := validator{
		testMask:   spin.gridDef.haveMask,
		testReels:  len(a.generateReels) > 0,
		testDupes1: !a.genAllowDupes && len(a.morphFor) == 0,
		testDupes2: !a.genAllowDupes && len(a.morphFor) > 0,
		gridSize:   spin.reelCount * spin.rowCount,
		rows:       spin.rowCount,
		spin:       spin,
		a:          a,
	}

	hasMult := len(spin.multipliers) > 0
	jump := calculateJump(v.gridSize)

	var generated uint8
	for ix := range a.symbolChances {
		chance := a.symbolChances[ix]
		if chance < 99.99 && !spin.TestChance2(a.ModifyChance(chance, spin)) {
			break
		}

		offs := spin.prng.IntN(v.gridSize)

		if chance > 99.99 {
			for !v.valid(offs) {
				if offs += jump; offs >= v.gridSize {
					offs -= v.gridSize
				}
			}
		} else {
			deadlock := 35
			for !v.valid(offs) && deadlock > 0 {
				if offs += jump; offs >= v.gridSize {
					offs -= v.gridSize
				}
				deadlock--
			}
			if deadlock == 0 {
				break
			}
		}

		generated++
		spin.indexes[offs] = a.symbol

		if a.genMultiplier != nil {
			if m := a.genMultiplier.RandomIndex(spin.prng); m > 1 {
				if len(spin.multipliers) < v.gridSize {
					spin.multipliers = utils.PurgeUInt16s(spin.multipliers, cap(spin.indexes))[:v.gridSize]
					clear(spin.multipliers)
					hasMult = true
				}
				spin.multipliers[offs] = uint16(m)
			} else if hasMult {
				spin.multipliers[offs] = 0
			}
		}
	}

	if generated == 0 {
		return false
	}

	if s := spin.GetSymbol(a.symbol); s != nil {
		if s.isWild {
			spin.newWilds += generated
		}
		if s.isScatter {
			spin.newScatters += generated
		}
		if s.isHero {
			spin.newHeroes += generated
		}
	}

	return true
}

func (a *ReviseAction) doGenerateReelSymbols(spin *Spin) bool {
	prng := spin.prng

	count := int(a.symbolWeights.RandomIndex(prng))
	if count == 0 {
		return false
	}

	rows, hasMult := spin.rowCount, len(spin.multipliers) > 0
	done := make([]bool, spin.reelCount, 16)

	for ix := 0; ix < count; ix++ {
		reel := a.reelWeights.RandomIndex(prng) - 1
		for done[reel] {
			reel = a.reelWeights.RandomIndex(prng) - 1
		}
		done[reel] = true

		row := prng.IntN(rows)
		offs := int(reel)*rows + row
		spin.indexes[offs] = a.symbol

		if hasMult {
			spin.multipliers[offs] = 0
		}
	}

	if s := spin.GetSymbol(a.symbol); s != nil {
		if s.isWild {
			spin.newWilds += uint8(count)
		}
		if s.isScatter {
			spin.newScatters += uint8(count)
		}
		if s.isHero {
			spin.newHeroes += uint8(count)
		}
	}

	return true
}

func (a *ReviseAction) doGenerateSymbols(spin *Spin) bool {
	v := validator{
		testMask:   spin.gridDef.haveMask,
		testReels:  len(a.generateReels) > 0,
		testDupes1: !a.genAllowDupes && len(a.morphFor) == 0,
		testDupes2: !a.genAllowDupes && len(a.morphFor) > 0,
		gridSize:   spin.reelCount * spin.rowCount,
		rows:       spin.rowCount,
		spin:       spin,
		a:          a,
	}

	hasMult := len(spin.multipliers) > 0
	jump := calculateJump(v.gridSize)

	var generated bool
	for ix := range a.symbolChances {
		if !spin.TestChance2(a.ModifyChance(a.symbolChances[ix], spin)) {
			break
		}

		offs := spin.prng.IntN(v.gridSize)
		deadlock := 35

		for !v.valid(offs) {
			if deadlock--; deadlock == 0 {
				break
			}
			if offs += jump; offs >= v.gridSize {
				offs -= v.gridSize
			}
		}

		if deadlock == 0 {
			break
		}

		generated = true
		symbol := a.genWeights.RandomIndex(spin.prng)
		spin.indexes[offs] = symbol

		if s := spin.GetSymbol(symbol); s != nil {
			if s.isWild {
				spin.newWilds++
			}
			if s.isScatter {
				spin.newScatters++
			}
			if s.isHero {
				spin.newHeroes++
			}
		}

		if hasMult {
			spin.multipliers[offs] = 0
		}
	}

	return generated
}

func (a *ReviseAction) doGenerateBonus(spin *Spin) bool {
	symbol := spin.bonusSymbol
	if spin.kind != FreeSpin || symbol == 0 || symbol == utils.MaxIndex || int(symbol) >= len(a.bonusChances) {
		return false
	}

	if !spin.TestChance2(a.ModifyChance(a.bonusChances[symbol], spin)) {
		return false
	}

	rows, reels := spin.rowCount, spin.reelCount
	done := make([]bool, 0, 10)[:reels]
	var count uint8

	for reel := range done {
		if spin.slots.hotReelsAsBonusSymbol && spin.hot[reel] {
			done[reel] = true
			count++
		} else {
			offset := reel * rows
			end := offset + rows
			for offset < end {
				if spin.indexes[offset] == spin.bonusSymbol {
					done[reel] = true
					count++
					break
				}
				offset++
			}
		}
	}

	if count >= a.bonusCount {
		return false
	}

	n := reels * rows
	for count < a.bonusCount {
		offset := spin.prng.IntN(n)
		reel := offset / rows
		if !done[reel] {
			spin.indexes[offset] = spin.bonusSymbol
			done[reel] = true
			count++
		}
	}

	return true
}

func (a *ReviseAction) doGenerateShape(spin *Spin) bool {
	if !spin.TestChance2(a.ModifyChance(a.shapeChance, spin)) {
		return false
	}

	// select a random symbol.
	symbol := a.shapeWeights.RandomIndex(spin.prng)

	// select a random center location.
	var c int
	if len(a.shapeCenters) > 1 {
		c = spin.prng.IntN(len(a.shapeCenters))
	}
	center := a.shapeCenters[c]

	rows := spin.rowCount

	if !a.genAllowOld {
		// check if the shape hits a sticky symbol.
		var already uint8
		for ix := range a.shapeGrid {
			reel, row := center[0]+a.shapeGrid[ix][0], center[1]+a.shapeGrid[ix][1]
			offset := reel*rows + row
			if spin.sticky[offset] {
				symbol = spin.indexes[offset]
				already++
			}
		}
		if spin.CountSymbol(symbol) > already {
			return false
		}
	}

	// fill the shape.
	var generated uint8
	for ix := range a.shapeGrid {
		reel, row := center[0]+a.shapeGrid[ix][0], center[1]+a.shapeGrid[ix][1]
		spin.indexes[reel*rows+row] = symbol
		generated++
	}

	if s := spin.GetSymbol(symbol); s != nil {
		if s.isWild {
			spin.newWilds += generated
		}
		if s.isScatter {
			spin.newScatters += generated
		}
		if s.isHero {
			spin.newHeroes += generated
		}
	}

	return true
}

func calculateJump(gridSize int) int {
	jump := 7
	switch {
	case gridSize <= 3:
		return 1
	case gridSize <= 5:
		return 3
	case gridSize <= 7:
		return 5
	default:
		for gridSize%jump == 0 {
			jump += 2
		}
		return jump
	}
}

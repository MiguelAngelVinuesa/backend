package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// testPaylines tests the configured paylines on the given spin, adds any wins to the given slice and returns the wins.
func (a *PayoutAction) testPaylines(spin *Spin, res *results.Result) bool {
	if s := spin.paylines; s != nil {
		return s.GetPayouts(spin, res)
	}
	return false
}

// testAllPaylines tests all paylines on the given spin, adds any wins to the given slice and returns the wins.
func (a *PayoutAction) testAllPaylines(spin *Spin, res *results.Result, highestPayout bool) bool {
	var hits int
	reelSize := spin.ReelSize(0)

	test := &allPaylineTester{
		highestPayout: highestPayout,
		multipliers:   len(spin.multipliers) > 0,
		reelCount:     spin.reelCount,
		rowCount:      spin.rowCount,
		spin:          spin,
		symbols:       spin.GetSymbols(),
		res:           res,
		rows:          make(utils.UInt8s, 16)[:spin.reelCount],
		mults:         make([]float64, 16)[:spin.reelCount],
	}

	if !test.multipliers {
		for row := 0; row < reelSize; row++ {
			if symbol := test.symbols.GetSymbol(spin.indexes[row]); symbol != nil {
				test.rows[0] = uint8(row)
				test.mults[0] = symbol.multiplier
				if symbol.isWild {
					hits += test.checkRemainingReels(1, 1, symbol)
				} else {
					hits += test.checkRemainingReels(1, 0, symbol)
				}
			}
		}
	} else {
		for row := 0; row < reelSize; row++ {
			if symbol := test.symbols.GetSymbol(spin.indexes[row]); symbol != nil {
				test.rows[0] = uint8(row)
				test.mults[0] = utils.NewMultiplier(symbol.multiplier, float64(spin.multipliers[row]))
				if symbol.isWild {
					hits += test.checkRemainingReelsMulti(1, 1, symbol)
				} else {
					hits += test.checkRemainingReelsMulti(1, 0, symbol)
				}
			}
		}
	}

	return hits > 0
}

type allPaylineTester struct {
	highestPayout bool
	multipliers   bool
	reelCount     int
	rowCount      int
	spin          *Spin
	symbols       *SymbolSet
	res           *results.Result
	rows          utils.UInt8s
	mults         []float64
}

// checkRemainingReels is a recursive function for the "remaining" reels during testing of the AllPaylines feature.
// The function returns the number of payouts recorded.
// from is the reel from which to test onwards. If it matches the reelCount the recursive nature of the function
// has passed through all reels and therefor warrants the maximum payout for the symbol.
// If no payouts were recorded during recursive calls, it means there are no matching symbols on the following reels,
// so the current reelCount is the maximum matched for this symbol. If the number of matching reels warrants a payout,
// then a new winline is stored and a count of 1 will be returned blocking lower counts from also recording a winline.
func (a *allPaylineTester) checkRemainingReels(from, wilds int, symbol *Symbol) int {
	if from == a.reelCount {
		return a.checkHighestPayout(symbol, from, wilds)
	}

	var hits int
	reelSize := a.spin.ReelSize(uint8(from))
	offset := from * a.rowCount

	for row := 0; row < reelSize; row++ {
		if next := a.symbols.GetSymbol(a.spin.indexes[offset]); next != nil {
			a.rows[from] = uint8(row)
			a.mults[from] = next.multiplier

			switch {
			case symbol.IsWild(), symbol.WildFor(next.id):
				if next.isWild {
					hits += a.checkRemainingReels(from+1, wilds+1, next)
				} else {
					hits += a.checkRemainingReels(from+1, wilds, next)
				}

			case next.IsWild(), next == symbol, next.WildFor(symbol.id):
				if next.isWild {
					hits += a.checkRemainingReels(from+1, wilds+1, symbol)
				} else {
					hits += a.checkRemainingReels(from+1, wilds, symbol)
				}
			}

			a.rows[from] = 0
			a.mults[from] = 0
		}

		offset++
	}

	if hits == 0 {
		return a.checkHighestPayout(symbol, from, wilds)
	}
	return hits
}

// checkRemainingReelsMulti works like checkRemainingReels() but for spins with multipliers.
func (a *allPaylineTester) checkRemainingReelsMulti(from, wilds int, symbol *Symbol) int {
	if from == a.reelCount {
		return a.checkHighestPayout(symbol, from, wilds)
	}

	var hits int
	reelSize := a.spin.ReelSize(uint8(from))
	offset := from * a.rowCount

	for row := 0; row < reelSize; row++ {
		if next := a.symbols.GetSymbol(a.spin.indexes[offset]); next != nil {
			a.rows[from] = uint8(row)
			a.mults[from] = utils.NewMultiplier(next.multiplier, float64(a.spin.multipliers[offset]))

			switch {
			case symbol.IsWild(), symbol.WildFor(next.id):
				if next.isWild {
					hits += a.checkRemainingReelsMulti(from+1, wilds+1, next)
				} else {
					hits += a.checkRemainingReelsMulti(from+1, wilds, next)
				}

			case next.IsWild(), next == symbol, next.WildFor(symbol.id):
				if next.isWild {
					hits += a.checkRemainingReelsMulti(from+1, wilds+1, symbol)
				} else {
					hits += a.checkRemainingReelsMulti(from+1, wilds, symbol)
				}
			}

			a.rows[from] = 0
			a.mults[from] = 0
		}

		offset++
	}

	if hits == 0 {
		return a.checkHighestPayout(symbol, from, wilds)
	}
	return hits
}

func (a *allPaylineTester) checkHighestPayout(symbol *Symbol, count, wilds int) int {
	var starting int
	var hwp float64
	var hps utils.Index
	if a.highestPayout && wilds > 0 {
		if starting = a.startingWilds(count); starting >= 2 {
			hwp = a.symbols.bestWildPay[starting]
			if hwp > 0 {
				hps = a.symbols.bestWildSym[starting]
			}
		}
	}

	c := uint8(count)
	p := symbol.Payout(c)

	switch {
	case p > 0 && p >= hwp:
		// normal symbol payout.
		a.res.AddPayouts(AllPaylinePayout(p, utils.NewMultiplier(a.getMultiplier(count), a.spin.getMultiplier(wilds)), symbol.id, c, a.rows))
		a.spin.markAllPayline(c, a.rows)
		return 1

	case hwp > 0:
		// highest paying symbol from starting wilds.
		c = uint8(starting)
		a.res.AddPayouts(AllPaylinePayout(hwp, utils.NewMultiplier(a.getMultiplier(starting), a.spin.getMultiplier(starting)), hps, c, a.rows))
		a.spin.markAllPayline(c, a.rows)
		return 1
	}

	return 0
}

func (a *allPaylineTester) getMultiplier(count int) float64 {
	return utils.NewMultiplier(a.mults[:count]...)
}

func (a *allPaylineTester) startingWilds(count int) int {
	var out int

	for ix := 0; ix < count; ix++ {
		offset := ix*a.rowCount + int(a.rows[ix])
		if symbol := a.symbols.GetSymbol(a.spin.indexes[offset]); symbol == nil || !symbol.IsWild() {
			break
		}
		out++
	}

	return out
}

package slots

import (
	"math"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
)

func (a *PayoutAction) testRemovePayouts(spin *Spin, res *results.Result) bool {
	switch a.remPayoutMech {
	case 1, 3:
		if t := math.Round(res.Total*100) / 100; t != 0.0 && t >= a.remMinFactor && t <= a.remMaxFactor {
			return a.chanceRemovePayouts(spin, res, a.remPayoutChance)
		}

	case 2:
		if n := spin.CountBonusSymbol(); n > 2 {
			symbol := spin.GetSymbol(spin.bonusSymbol)
			if symbol == nil {
				panic(consts.MsgSymbolNotFound)
			}
			if t := symbol.Payout(n); t != 0.0 && t >= a.remMinFactor && t <= a.remMaxFactor {
				return a.chanceRemoveBonus(spin, res, a.remPayoutChance)
			}
		}
	}

	return false
}

func (a *PayoutAction) testRemovePayoutBands(spin *Spin, res *results.Result) bool {
	t := math.Round(res.Total*100) / 100
	if t == 0.0 {
		return false
	}

	var b *RemovePayoutBand
	for ix := range a.remBands {
		b = &a.remBands[ix]
		if t >= b.MinPayout && t < b.MaxPayout {
			break
		}
		b = nil
	}

	if b == nil {
		return false
	}

	return a.chanceRemovePayouts(spin, res, b.RemoveChance)
}

func (a *PayoutAction) chanceRemovePayouts(spin *Spin, res *results.Result, chance float64) bool {
	if n := float64(spin.prng.IntN(10000)) / 100.0; n >= chance {
		return false
	}

	switch a.remPayoutMech {
	case 1:
		if !spin.PreventPaylines(a.remPayoutDir, a.remPayoutWilds, a.remDupes) {
			return false
		}
	case 2:
		if !spin.PreventBonus(a.remDupes) {
			return false
		}
	case 3:
		if !spin.PreventPaylines2(a.remPayoutDir, a.remPayoutWilds, a.remDupes) {
			return false
		}
	}

	// fix the initial array in the result.
	if sr, ok := res.Data.(*SpinResult); ok {
		copy(sr.initial, spin.indexes)
	}

	// remove the payouts.
	spin.resetPayouts()
	res.ReleasePayouts()

	return true
}

func (a *PayoutAction) chanceRemoveBonus(spin *Spin, res *results.Result, chance float64) bool {
	if n := float64(spin.prng.IntN(10000)) / 100.0; n >= chance {
		return false
	}

	if !spin.PreventBonus(a.remDupes) {
		return false
	}

	// fix the initial array in the result.
	if sr, ok := res.Data.(*SpinResult); ok {
		copy(sr.initial, spin.indexes)
	}

	return true
}

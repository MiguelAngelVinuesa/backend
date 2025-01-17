package bot

import (
	magic "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/analysis/simulation/slots"
	comp "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	game "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/slots"
	rslt "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	util "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/utils/conv"
)

func Conditions() map[string]*magic.Condition {
	return conditions
}

func MakeMatcher(key string, params map[string]any, game *game.Regular) magic.Matcher {
	switch key {
	case keyBonusSymbol:
		return testBonusSymbol(
			magic.SymbolIndex(AllSymbols(), conv.IntFromAny(params[fieldSymbol])),
		)

	case keySpinHotReel:
		return testSpinHotReel(
			conv.IntFromAny(params[fieldSequence]),
			conv.IntsFromAny(params[fieldReels]),
		)

	case keyAnyHotReel:
		return testAnyHotReel(
			conv.IntsFromAny(params[fieldReels]),
		)

	case keyBonusPayout:
		return testBonusPayout(
			conv.IntsFromAny(params[fieldReels]),
		)

	case keyFullBonusPayout:
		return testFullBonusPayout()

	default:
		return magic.MakeMatcher(key, params, AllSymbols(), game)
	}
}

const (
	keyBonusSymbol     = "bonus-symbol"
	keySpinHotReel     = "spin-hot-reel"
	keyAnyHotReel      = "any-hot-reel"
	keyBonusPayout     = "bonus-payout"
	keyFullBonusPayout = "full-bonus-payout"

	fieldSymbol   = "symbol"
	fieldSequence = "seq"
	fieldReels    = "reels..."
)

func testBonusSymbol(symbol util.Index) magic.Matcher {
	return func(res rslt.Results) bool {
		if len(res) < 2 {
			return false
		}

		if spin, ok := res[1].Data.(*comp.SpinResult); ok {
			if spin.BonusSymbol() == symbol {
				return true
			}
		}
		return false
	}
}

func testSpinHotReel(spin int, reels []int) magic.Matcher {
	spin--

	return func(res rslt.Results) bool {
		if len(res) <= spin {
			return false
		}

		if r, ok := res[spin].Data.(*comp.SpinResult); ok {
			if len(reels) == 0 {
				return len(r.Hot()) > 0
			}

			var count int
			for ix, reel := range reels {
				for _, hot := range r.Hot() {
					if int(hot) == reel {
						count++
						break
					}
				}
				if count != ix+1 {
					return false
				}
			}
			return true
		}

		return false
	}
}

func testAnyHotReel(reels []int) magic.Matcher {
	return func(res rslt.Results) bool {
		for ix := range res {
			if ix > 0 {
				if r, ok := res[ix].Data.(*comp.SpinResult); ok {
					if len(reels) == 0 {
						return len(r.Hot()) > 0
					}

					var count int
					for iy, reel := range reels {
						for _, hot := range r.Hot() {
							if int(hot) == reel {
								count++
								break
							}
						}
						if count != iy+1 {
							return false
						}
					}
					return true
				}
			}
		}
		return false
	}
}

func testFullBonusPayout() magic.Matcher {
	return testBonusPayout([]int{1, 2, 3, 4, 5})
}

func testBonusPayout(reels []int) magic.Matcher {
	return func(res rslt.Results) bool {
		for ix := range res {
			if ix > 0 {
				if bonusPayout(res[ix], reels) {
					return true
				}
			}
		}
		return false
	}
}

func bonusPayout(r *rslt.Result, reels []int) bool {
	if s, ok := r.Data.(*comp.SpinResult); ok {
		for iy := range r.Payouts {
			if p := r.Payouts[iy]; p.Kind() == rslt.SlotBonusSymbol {
				var count int
				for iz := range reels {
					if !bonusReel(s, reels[iz]) {
						break
					}
					count++
				}
				if count == len(reels) {
					return true
				}
			}
		}
	}
	return false
}

func bonusReel(s *comp.SpinResult, reel int) bool {
	for _, hot := range s.Hot() {
		if int(hot) == reel {
			return true
		}
	}

	grid := s.Initial()
	symbol := s.BonusSymbol()

	high := rows * reel
	for offs := high - rows; offs < high; offs++ {
		if grid[offs] == symbol {
			return true
		}
	}

	return false
}

var conditions = make(map[string]*magic.Condition)

func init() {
	conditions[keyBonusSymbol] = magic.NewCondition(keyBonusSymbol, magic.NewSymbolParam())
	conditions[keySpinHotReel] = magic.NewCondition(keySpinHotReel, magic.NewIntParam(fieldSequence, 2, 999), magic.NewReelsParam())
	conditions[keyAnyHotReel] = magic.NewCondition(keyAnyHotReel, magic.NewReelsParam())
	conditions[keyBonusPayout] = magic.NewCondition(keyBonusPayout, magic.NewReelsParam())
	conditions[keyFullBonusPayout] = magic.NewCondition(keyFullBonusPayout)

	for ix := range magic.Builtins {
		c := magic.Builtins[ix]
		conditions[c.Name] = c
	}
}

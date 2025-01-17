package owl

import (
	"strings"

	magic "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/analysis/simulation/slots"
	game "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/wheel"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/utils/conv"
)

func Conditions() map[string]*magic.Condition {
	return conditions
}

func MakeMatcher(key string, params map[string]any, game *game.Regular) magic.Matcher {
	switch key {
	case keyInstantBonus:
		return testInstantBonus(conv.StringFromAny(params[fieldResult]))
	case keyBonusWheel:
		return testBonusWheel(conv.StringFromAny(params[fieldResult]))
	default:
		return magic.MakeMatcher(key, params, AllSymbols(), game)
	}
}

const (
	keyInstantBonus = "instant-bonus"
	keyBonusWheel   = "bonus-wheel"

	fieldResult = "result"
)

func testInstantBonus(s string) magic.Matcher {
	var kind utils.Index
	switch strings.ToLower(s) {
	case "1", "random-wilds":
		kind = 1
	case "2", "wild-reels":
		kind = 2
	case "3", "bonus-wheel":
		kind = 3
	}

	return func(res results.Results) bool {
		if len(res) == 0 {
			return false
		}

		for ix := range res {
			if r := res[ix]; r.DataKind == results.InstantBonusData {
				if b, ok := res[ix].Data.(*results.BonusSelector); ok {
					if b.Chosen() == kind {
						return true
					}
				}
			}
			if ix > 3 {
				break
			}
		}

		return false
	}
}

func testBonusWheel(s string) magic.Matcher {
	var kind utils.Index
	switch strings.ToLower(s) {
	case "1", "4", "random-wilds":
		kind = 4
	case "2", "5", "wild-reels":
		kind = 5
	case "3", "6", "scatter-progressive":
		kind = 6
	}

	return func(res results.Results) bool {
		if len(res) == 0 {
			return false
		}

		for ix := range res {
			if r := res[ix]; r.DataKind == results.BonusWheelData {
				if b, ok := res[ix].Data.(*wheel.BonusWheelResult); ok {
					if b.Result() == kind {
						return true
					}
				}
			}
			if ix > 3 {
				break
			}
		}

		return false
	}
}

var conditions = make(map[string]*magic.Condition)

func init() {
	conditions[keyInstantBonus] = magic.NewCondition(keyInstantBonus, magic.NewStringParam(fieldResult))
	conditions[keyBonusWheel] = magic.NewCondition(keyBonusWheel, magic.NewStringParam(fieldResult))

	for ix := range magic.Builtins {
		c := magic.Builtins[ix]
		conditions[c.Name] = c
	}
}

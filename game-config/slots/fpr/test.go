package fpr

import (
	magic "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/analysis/simulation/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	game "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
)

func Conditions() map[string]*magic.Condition {
	return conditions
}

func MakeMatcher(key string, params map[string]any, game *game.Regular) magic.Matcher {
	switch key {
	case keyWildRespin:
		return testWildRespin()
	case keyBonus:
		return testBonus()
	case keyFrostyBonus:
		return testFrostyBonus()
	default:
		return magic.MakeMatcher(key, params, AllSymbols(), game)
	}
}

const (
	keyWildRespin  = "wild-respin"
	keyBonus       = "princess-spin-bonus"
	keyFrostyBonus = "frost-princess-spin-bonus"
)

func testWildRespin() magic.Matcher {
	return func(res results.Results) bool {
		if len(res) == 0 {
			return false
		}

		if spin, ok := res[0].Data.(*slots.SpinResult); ok {
			return scatterCount(spin) < 3 && isWildRespin(spin)
		}
		return false
	}
}

func testBonus() magic.Matcher {
	return func(res results.Results) bool {
		if len(res) == 0 {
			return false
		}

		if spin, ok := res[0].Data.(*slots.SpinResult); ok {
			return scatterCount(spin) >= 3 && !isWildRespin(spin)
		}
		return false
	}
}

func testFrostyBonus() magic.Matcher {
	return func(res results.Results) bool {
		if len(res) == 0 {
			return false
		}

		if spin, ok := res[0].Data.(*slots.SpinResult); ok {
			return scatterCount(spin) >= 3 && isWildRespin(spin)
		}
		return false
	}
}

func scatterCount(spin *slots.SpinResult) int {
	var c int
	for ix := range spin.Initial() {
		if spin.Initial()[ix] == scatter {
			c++
		}
	}
	return c
}

func isWildRespin(spin *slots.SpinResult) bool {
	return spin.Initial()[0] == wild && spin.Initial()[1] == wild && spin.Initial()[20] == wild && spin.Initial()[21] == wild
}

var conditions = make(map[string]*magic.Condition)

func init() {
	conditions[keyWildRespin] = magic.NewCondition(keyWildRespin)
	conditions[keyBonus] = magic.NewCondition(keyBonus)
	conditions[keyFrostyBonus] = magic.NewCondition(keyFrostyBonus)

	for ix := range magic.Builtins {
		c := magic.Builtins[ix]
		conditions[c.Name] = c
	}
}

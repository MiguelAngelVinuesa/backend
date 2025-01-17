package frm

import (
	magic "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/analysis/simulation/slots"
	game "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/slots"
)

func Conditions() map[string]*magic.Condition {
	return conditions
}

func MakeMatcher(key string, params map[string]any, game *game.Regular) magic.Matcher {
	switch key {
	// case keyWildRespin:
	// 	return testWildRespin()
	default:
		return magic.MakeMatcher(key, params, AllSymbols(), game)
	}
}

// const (
// 	keyWildRespin = "wild-respin"
// )

// func testWildRespin() magic.Matcher {
// 	return func(res results.Results) bool {
// 		if len(res) == 0 {
// 			return false
// 		}
//
// 		if spin, ok := res[0].Data.(*slots.SpinResult); ok {
// 			return scatterCount(spin) < 3 && isWildRespin(spin)
// 		}
// 		return false
// 	}
// }

var conditions = make(map[string]*magic.Condition)

func init() {
	//	conditions[keyWildRespin] = magic.NewCondition(keyWildRespin)

	for ix := range magic.Builtins {
		c := magic.Builtins[ix]
		conditions[c.Name] = c
	}
}

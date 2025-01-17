package ccb

import (
	magic "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/analysis/simulation/slots"
	comp "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	game "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/slots"
	rslt "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	util "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func Conditions() map[string]*magic.Condition {
	return conditions
}

func MakeMatcher(key string, params map[string]any, game *game.Regular) magic.Matcher {
	switch key {
	case keyFirstSuperX:
		return testFirstSuperX(
			magic.SymbolIndex(AllSymbols(), params[fieldSymbol]),
		)

	case keyFreeSuperX:
		return testFreeSuperX(
			magic.SymbolIndex(AllSymbols(), params[fieldSymbol]),
		)

	default:
		return magic.MakeMatcher(key, params, AllSymbols(), game)
	}
}

const (
	keyFirstSuperX = "first-super-x"
	keyFreeSuperX  = "free-super-x"

	fieldSymbol = "symbol"
)

func testFirstSuperX(symbol util.Index) magic.Matcher {
	return func(res rslt.Results) bool {
		if len(res) == 0 {
			return false
		}

		if r, ok := res[0].Data.(*comp.SpinResult); ok {
			return superX(r.Initial(), symbol)
		}
		return false
	}
}

func testFreeSuperX(symbol util.Index) magic.Matcher {
	return func(res rslt.Results) bool {
		if len(res) == 0 {
			return false
		}

		for ix := 2; ix < len(res); ix++ {
			if r, ok := res[ix].Data.(*comp.SpinResult); ok {
				if superX(r.Initial(), symbol) {
					return true
				}
			}
		}
		return false
	}
}

func superX(grid util.Indexes, symbol util.Index) bool {
	return hasSuperX(grid, symbol, 0) || hasSuperX(grid, symbol, 3) || hasSuperX(grid, symbol, 6)
}

func hasSuperX(grid util.Indexes, symbol util.Index, start int) bool {
	if symbol == util.NullIndex {
		symbol = grid[start]
	}

	if symbol != grid[start] ||
		symbol != grid[start+2] ||
		symbol != grid[start+4] ||
		symbol != grid[start+6] ||
		symbol != grid[start+8] {
		return false
	}

	var count int
	for ix := range grid {
		if symbol == grid[ix] {
			count++
		}
	}

	return count == 5
}

var conditions = make(map[string]*magic.Condition)

func init() {
	conditions[keyFirstSuperX] = magic.NewCondition(keyFirstSuperX, magic.NewSymbolParam())
	conditions[keyFreeSuperX] = magic.NewCondition(keyFreeSuperX, magic.NewSymbolParam())

	for ix := range magic.Builtins {
		c := magic.Builtins[ix]
		conditions[c.Name] = c
	}
}

package crw

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

var conditions = make(map[string]*magic.Condition)

func init() {
	for ix := range magic.Builtins {
		c := magic.Builtins[ix]
		conditions[c.Name] = c
	}
}

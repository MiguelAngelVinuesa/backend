package magic

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/utils/conv"

	comp "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	rslt "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	util "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// SpinGrid matches if the given spin matches the given grid.
// A zero symbol index in the grid matches any symbol in the result.
// If the given grid is too short, missing tiles always match.
func SpinGrid(spin int, grid util.Indexes) Matcher {
	spin--

	return func(res rslt.Results) bool {
		res = SpinResults(res)
		if len(res) < spin {
			return false
		}

		if r, ok := res[spin].Data.(*comp.SpinResult); ok {
			return matchGrid(r, grid)
		}
		return false
	}
}

// AnyGrid matches if the given spin matches the given grid.
// A zero symbol index in the grid matches any symbol in the result.
// If the given grid is too short, missing tiles always match.
func AnyGrid(grid util.Indexes) Matcher {
	return func(res rslt.Results) bool {
		for ix := range res {
			if r, ok := res[ix].Data.(*comp.SpinResult); ok {
				if matchGrid(r, grid) {
					return true
				}
			}
		}
		return false
	}
}

func matchGrid(r *comp.SpinResult, want util.Indexes) bool {
	got := r.Initial()
	for ix, symbol := range want {
		if !symbolMatch(symbol, got[ix]) {
			return false
		}
	}
	return true
}

func indexesFromAny(param any) util.Indexes {
	list := conv.IntsFromAny(param)
	out := make(util.Indexes, len(list))
	for ix := range list {
		out[ix] = util.Index(list[ix])
	}
	return out
}

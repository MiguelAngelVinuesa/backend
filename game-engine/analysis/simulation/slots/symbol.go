package magic

import (
	"strconv"

	comp "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	game "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/slots"
	rslt "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	util "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// SymbolIndex returns the ID of the symbol matching the input parameter.
// The function attempts to match on the symbol ID, symbol name and symbol resource.
// If no match is found, the invalid symbol ID is returned.
func SymbolIndex(symbols *comp.SymbolSet, param any) util.Index {
	name := util.StringFromAny(param)
	id, _ := strconv.Atoi(name)

	last := symbols.GetMaxSymbolID()
	for ix := util.Index(1); ix <= last; ix++ {
		if symbol := symbols.GetSymbol(ix); symbol != nil {
			if symbol.ID() == util.Index(id) || symbol.Name() == name || symbol.Resource() == name {
				return symbol.ID()
			}
		}
	}

	return util.MaxIndex
}

// SpinSymbolCount determines if the given spin contains exactly count symbols.
func SpinSymbolCount(spin, count int, symbol util.Index, game *game.Regular, reels []int) Matcher {
	spin--

	return func(res rslt.Results) bool {
		res = SpinResults(res)
		if len(res) < spin {
			return false
		}

		if r, ok := res[spin].Data.(*comp.SpinResult); ok {
			return countSymbols(symbol, r.Initial(), game, reels) == count
		}
		return false
	}
}

// SpinSymbolRange determines if the given spin contains count symbols between min and max.
func SpinSymbolRange(spin, min, max int, symbol util.Index, game *game.Regular, reels []int) Matcher {
	spin--

	return func(res rslt.Results) bool {
		res = SpinResults(res)
		if len(res) < spin {
			return false
		}

		if r, ok := res[spin].Data.(*comp.SpinResult); ok {
			c := countSymbols(symbol, r.Initial(), game, reels)
			return c >= min || (max < min || c <= max)
		}

		return false
	}
}

// AnySymbolCount determines if any spin contains exactly count symbols.
func AnySymbolCount(count int, symbol util.Index, game *game.Regular, reels []int) Matcher {
	return func(res rslt.Results) bool {
		for ix := range res {
			if r, ok := res[ix].Data.(*comp.SpinResult); ok {
				if countSymbols(symbol, r.Initial(), game, reels) == count {
					return true
				}
			}
		}
		return false
	}
}

// AnySymbolRange determines if any spin contains count symbols between min and max.
func AnySymbolRange(min, max int, symbol util.Index, game *game.Regular, reels []int) Matcher {
	return func(res rslt.Results) bool {
		for ix := range res {
			if r, ok := res[ix].Data.(*comp.SpinResult); ok {
				c := countSymbols(symbol, r.Initial(), game, reels)
				if c >= min || (max < min || c <= max) {
					return true
				}
			}
		}
		return false
	}
}

// AllSymbolCount determines if all spin together contain exactly count symbols.
func AllSymbolCount(count int, symbol util.Index, game *game.Regular, reels []int) Matcher {
	return func(res rslt.Results) bool {
		var c int
		for ix := range res {
			if r, ok := res[ix].Data.(*comp.SpinResult); ok {
				c += countSymbols(symbol, r.Initial(), game, reels)
			}
		}
		return count == c
	}
}

// AllSymbolRange determines if all spin together contain count symbols between min and max.
func AllSymbolRange(min, max int, symbol util.Index, game *game.Regular, reels []int) Matcher {
	return func(res rslt.Results) bool {
		var c int
		for ix := range res {
			if r, ok := res[ix].Data.(*comp.SpinResult); ok {
				c += countSymbols(symbol, r.Initial(), game, reels)
			}
		}
		return c >= min || (max < min || c <= max)
	}
}

func countSymbols(symbol util.Index, grid util.Indexes, game *game.Regular, reels []int) int {
	var count int

	if len(reels) == 0 {
		for _, s := range grid {
			if s == symbol {
				count++
			}
		}
		return count
	}

	reelCount, rowCount := game.Slots().ReelCount(), game.Slots().RowCount()

	for _, reel := range reels {
		if reel > 0 && reel <= reelCount {
			high := rowCount * reel
			for offs := high - rowCount; offs < high; offs++ {
				if grid[offs] == symbol {
					count++
				}
			}
		}
	}
	return count
}

func symbolMatch(s1, s2 util.Index) bool {
	return s1 == 0 || s1 == s2
}

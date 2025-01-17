package magic

import (
	"strings"

	comp "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	game "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/slots"
	rslt "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	util "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func SpinHasSymbolTypeCount(spin, count int, symbolType string, symbols *comp.SymbolSet, game *game.Regular, reels []int) Matcher {
	spin--
	compare := getSymbolCompare(symbolType)

	return func(res rslt.Results) bool {
		res = SpinResults(res)
		if len(res) < spin {
			return false
		}

		var c int
		if r, ok := res[spin].Data.(*comp.SpinResult); ok {
			c = countSymbolType(compare, r.Initial(), symbols, game, reels)
		}

		return count == c
	}
}

func SpinHasSymbolTypeRange(spin, min, max int, symbolType string, symbols *comp.SymbolSet, game *game.Regular, reels []int) Matcher {
	spin--
	compare := getSymbolCompare(symbolType)

	return func(res rslt.Results) bool {
		res = SpinResults(res)
		if len(res) < spin {
			return false
		}

		var c int
		if r, ok := res[spin].Data.(*comp.SpinResult); ok {
			c = countSymbolType(compare, r.Initial(), symbols, game, reels)
		}

		return c >= min && (max < min || c <= max)
	}
}

func AnyHasSymbolTypeCount(count int, symbolType string, symbols *comp.SymbolSet, game *game.Regular, reels []int) Matcher {
	compare := getSymbolCompare(symbolType)

	return func(res rslt.Results) bool {
		for ix := range res {
			var c int
			if r, ok := res[ix].Data.(*comp.SpinResult); ok {
				c = countSymbolType(compare, r.Initial(), symbols, game, reels)
			}
			if count == c {
				return true
			}
		}
		return false
	}
}

func AnyHasSymbolTypeRange(min, max int, symbolType string, symbols *comp.SymbolSet, game *game.Regular, reels []int) Matcher {
	compare := getSymbolCompare(symbolType)

	return func(res rslt.Results) bool {
		for ix := range res {
			var c int
			if r, ok := res[ix].Data.(*comp.SpinResult); ok {
				c = countSymbolType(compare, r.Initial(), symbols, game, reels)
			}
			if c >= min && (max < min || c <= max) {
				return true
			}
		}
		return false
	}
}

func AllHasSymbolTypeCount(count int, symbolType string, symbols *comp.SymbolSet, game *game.Regular, reels []int) Matcher {
	compare := getSymbolCompare(symbolType)

	return func(res rslt.Results) bool {
		var c int
		for ix := range res {
			if r, ok := res[ix].Data.(*comp.SpinResult); ok {
				c += countSymbolType(compare, r.Initial(), symbols, game, reels)
			}
		}
		return count == c
	}
}

func AllHasSymbolTypeRange(min, max int, symbolType string, symbols *comp.SymbolSet, game *game.Regular, reels []int) Matcher {
	compare := getSymbolCompare(symbolType)

	return func(res rslt.Results) bool {
		var c int
		for ix := range res {
			if r, ok := res[ix].Data.(*comp.SpinResult); ok {
				c += countSymbolType(compare, r.Initial(), symbols, game, reels)
			}
		}
		return c >= min && (max < min || c <= max)
	}
}

func countSymbolType(compare typeCompare, grid util.Indexes, symbols *comp.SymbolSet, game *game.Regular, reels []int) int {
	var count int

	if len(reels) == 0 {
		for ix := range grid {
			if symbol := symbols.GetSymbol(grid[ix]); symbol != nil && compare(symbol) {
				count++
			}
		}
		return count
	}

	reelCount, rowCount := game.Slots().ReelCount(), game.Slots().RowCount()

	for _, reel := range reels {
		if reel > 0 && reel <= reelCount {
			high := rowCount + reel
			for offs := high - rowCount; offs < high; offs++ {
				if symbol := symbols.GetSymbol(grid[offs]); symbol != nil && compare(symbol) {
					count++
				}
			}
		}
	}
	return count
}

type typeCompare func(symbol *comp.Symbol) bool

func getSymbolCompare(symbolType string) typeCompare {
	switch strings.ToLower(symbolType) {
	case "wild":
		return func(symbol *comp.Symbol) bool {
			return symbol.IsWild()
		}

	case "scatter":
		return func(symbol *comp.Symbol) bool {
			return symbol.IsScatter()
		}

	case "bomb":
		return func(symbol *comp.Symbol) bool {
			return symbol.IsBomb()
		}

	default:
		return func(symbol *comp.Symbol) bool {
			return false
		}
	}
}

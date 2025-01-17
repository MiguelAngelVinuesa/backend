package magic

import (
	comp "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	game "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/slots"
	rslt "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	util "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// SpinStickyCount determines if the given spin contains exactly count sticky tiles.
func SpinStickyCount(spin, count int, game *game.Regular, reels []int) Matcher {
	spin--

	return func(res rslt.Results) bool {
		res = SpinResults(res)
		if len(res) < spin {
			return false
		}

		if r, ok := res[spin].Data.(*comp.SpinResult); ok {
			if sticky := r.StickySymbol(); sticky > 0 {
				return countSymbols(sticky, r.Initial(), game, reels) == count
			}
		}
		return false
	}
}

// AnyStickyCount determines if any spin contains exactly count sticky tiles.
func AnyStickyCount(count int, game *game.Regular, reels []int) Matcher {
	return func(res rslt.Results) bool {
		for ix := range res {
			if r, ok := res[ix].Data.(*comp.SpinResult); ok {
				if sticky := r.StickySymbol(); sticky > 0 {
					return countSymbols(sticky, r.Initial(), game, reels) == count
				}
			}
		}
		return false
	}
}

// SpinStickySymbol determines if the given spin has the given sticky symbol.
func SpinStickySymbol(spin int, sticky util.Index) Matcher {
	spin--

	return func(res rslt.Results) bool {
		res = SpinResults(res)
		if len(res) < spin {
			return false
		}

		if r, ok := res[spin].Data.(*comp.SpinResult); ok {
			return r.StickySymbol() == sticky
		}
		return false
	}
}

// AnyStickySymbol determines if any spin has the given sticky symbol.
func AnyStickySymbol(sticky util.Index) Matcher {
	return func(res rslt.Results) bool {
		for ix := range res {
			if r, ok := res[ix].Data.(*comp.SpinResult); ok {
				return r.StickySymbol() == sticky
			}
		}
		return false
	}
}

// SpinStickyWild determines if the given spin has a count of sticky wilds between min and max.
func SpinStickyWild(spin int, min, max int, symbols *comp.SymbolSet) Matcher {
	spin--

	return func(res rslt.Results) bool {
		res = SpinResults(res)
		if len(res) < spin {
			return false
		}

		if r, ok := res[spin].Data.(*comp.SpinResult); ok {
			var count int
			for _, id := range r.Initial() {
				if s := symbols.GetSymbol(id); s != nil && s.IsWild() {
					count++
				}
			}
			return count >= min && (max < min || count <= max)
		}

		return false
	}
}

// AnyStickyWild determines if the given spin has a count of sticky wilds between min and max.
func AnyStickyWild(min, max int, symbols *comp.SymbolSet) Matcher {
	return func(res rslt.Results) bool {
		for ix := range res {
			if r, ok := res[ix].Data.(*comp.SpinResult); ok {
				var count int
				for _, id := range r.Initial() {
					if s := symbols.GetSymbol(id); s != nil && s.IsWild() {
						count++
					}
				}
				return count >= min && (max < min || count <= max)
			}
		}
		return false
	}
}

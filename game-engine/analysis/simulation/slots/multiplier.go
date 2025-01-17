package magic

import (
	comp "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	rslt "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	util "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// SpinMultiplierCount matches if the count of symbols with multiplier in the given spin is between min and max.
func SpinMultiplierCount(spin int, min, max int, symbols *comp.SymbolSet) Matcher {
	spin--

	return func(res rslt.Results) bool {
		res = SpinResults(res)
		if len(res) < spin {
			return false
		}

		if r, ok := res[spin].Data.(*comp.SpinResult); ok {
			var count int
			for ix, id := range r.Initial() {
				if s := symbols.GetSymbol(id); s != nil && util.ValidMultiplier(s.Multiplier()) {
					count++
				} else if ix < len(r.Multipliers()) && util.ValidMultiplier(float64(r.Multipliers()[ix])) {
					count++
				}
			}
			return count >= min && (max < min || count <= max)
		}
		return false
	}
}

// AnyMultiplierCount matches if the count of symbols with multiplier in any spin is between min and max.
func AnyMultiplierCount(min, max int, symbols *comp.SymbolSet) Matcher {
	return func(res rslt.Results) bool {
		for ix := range res {
			if r, ok := res[ix].Data.(*comp.SpinResult); ok {
				var count int
				for iy, id := range r.Initial() {
					if s := symbols.GetSymbol(id); s != nil && util.ValidMultiplier(s.Multiplier()) {
						count++
					} else if iy < len(r.Multipliers()) && util.ValidMultiplier(float64(r.Multipliers()[iy])) {
						count++
					}
				}
				return count >= min && (max < min || count <= max)
			}
		}
		return false
	}
}

// SpinMultiplierRange matches if the multiplier for the given spin is between min and max.
func SpinMultiplierRange(spin int, min, max uint16) Matcher {
	spin--

	return func(res rslt.Results) bool {
		res = SpinResults(res)
		if len(res) < spin {
			return false
		}

		if r, ok := res[spin].Data.(*comp.SpinResult); ok {
			for _, m := range r.Multipliers() {
				if m >= min && (max < min || m <= max) {
					return true
				}
			}
		}
		return false
	}
}

// AnyMultiplierRange matches if the multiplier for the last rslt is between min and max.
func AnyMultiplierRange(min, max uint16) Matcher {
	return func(res rslt.Results) bool {
		for ix := range res {
			if r, ok := res[ix].Data.(*comp.SpinResult); ok {
				for _, m := range r.Multipliers() {
					if m >= min && (max < min || m <= max) {
						return true
					}
				}
			}
		}
		return false
	}
}

// LastMultiplierRange matches if the multiplier for the last rslt is between min and max.
func LastMultiplierRange(min, max uint16) Matcher {
	return func(res rslt.Results) bool {
		res = SpinResults(res)
		if r, ok := res[len(res)-1].Data.(*comp.SpinResult); ok {
			for _, m := range r.Multipliers() {
				if m >= min && (max < min || m <= max) {
					return true
				}
			}
		}
		return false
	}
}

// SpinMultiplierSymbol matches if the given symbol appears in the given spin with a multiplier between min and max.
func SpinMultiplierSymbol(spin int, symbol util.Index, min, max float64, symbols *comp.SymbolSet) Matcher {
	spin--

	return func(res rslt.Results) bool {
		res = SpinResults(res)
		if len(res) < spin {
			return false
		}

		if r, ok := res[spin].Data.(*comp.SpinResult); ok {
			for iy, id := range r.Initial() {
				if s := symbols.GetSymbol(id); s != nil && symbolMatch(symbol, s.ID()) {
					if symbolMuliplierMinMax(s, min, max) || gridMuliplierMinMax(r.Multipliers(), iy, min, max) {
						return true
					}
				}

			}
		}
		return false
	}
}

// AnyMultiplierSymbol matches if the given symbol appears in the given spin with a multiplier between min and max.
func AnyMultiplierSymbol(symbol util.Index, min, max float64, symbols *comp.SymbolSet) Matcher {
	return func(res rslt.Results) bool {
		for ix := range res {
			if r, ok := res[ix].Data.(*comp.SpinResult); ok {
				for iy, id := range r.Initial() {
					if s := symbols.GetSymbol(id); s != nil && symbolMatch(symbol, s.ID()) {
						if symbolMuliplierMinMax(s, min, max) || gridMuliplierMinMax(r.Multipliers(), iy, min, max) {
							return true
						}
					}

				}
			}
		}
		return false
	}
}

// SpinMultiplierWild matches if a wild symbol appears in the given spin with a multiplier between min and max.
func SpinMultiplierWild(spin int, min, max float64, symbols *comp.SymbolSet) Matcher {
	spin--

	return func(res rslt.Results) bool {
		res = SpinResults(res)
		if len(res) < spin {
			return false
		}

		if r, ok := res[spin].Data.(*comp.SpinResult); ok {
			for iy, id := range r.Initial() {
				if s := symbols.GetSymbol(id); s != nil && s.IsWild() {
					if symbolMuliplierMinMax(s, min, max) || gridMuliplierMinMax(r.Multipliers(), iy, min, max) {
						return true
					}
				}

			}
		}
		return false
	}
}

// AnyMultiplierWild matches if a wild symbol appears in any spin with a multiplier between min and max.
func AnyMultiplierWild(min, max float64, symbols *comp.SymbolSet) Matcher {
	return func(res rslt.Results) bool {
		for ix := range res {
			if r, ok := res[ix].Data.(*comp.SpinResult); ok {
				for iy, id := range r.Initial() {
					if s := symbols.GetSymbol(id); s != nil && s.IsWild() {
						if symbolMuliplierMinMax(s, min, max) || gridMuliplierMinMax(r.Multipliers(), iy, min, max) {
							return true
						}
					}
				}
			}
		}
		return false
	}
}

func symbolMuliplierMinMax(symbol *comp.Symbol, min, max float64) bool {
	if m := symbol.Multiplier(); m >= min && (max < min || m <= max) {
		return true
	}
	return false
}

func gridMuliplierMinMax(multipliers []uint16, ix int, min, max float64) bool {
	if ix < len(multipliers) {
		if m := float64(multipliers[ix]); m >= min && (max < min || m <= max) {
			return true
		}
	}
	return false
}

// RoundMultiplier matches if the overal game multiplier on the last spin is between min and max.
func RoundMultiplier(min, max float64) Matcher {
	return func(res rslt.Results) bool {
		res = SpinResults(res)
		if len(res) < 0 {
			return false
		}

		if r, ok := res[len(res)-1].Data.(*comp.SpinResult); ok {
			m := r.Multiplier()
			return m >= min && (max < min || m <= max)
		}
		return false
	}
}

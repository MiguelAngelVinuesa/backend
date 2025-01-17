package rng_magic

import (
	"fmt"
	"strconv"

	magic "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/analysis/simulation/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func checkParams(cond *magic.Condition, symbols *slots.SymbolSet, params []string) ([]any, error) {
	out := make([]any, 0, 16)

	for ix, param := range cond.Parameters {
		var value string
		if ix < len(params) {
			value = params[ix]
		}

		switch param.Kind {
		case "INT":
			var i int

			if value != "" {
				i2, err := strconv.ParseInt(value, 10, 64)
				if err != nil {
					return nil, fmt.Errorf(msgInvalidInt, value)
				}

				i = int(i2)

				if float64(i) < param.Min || float64(i) > param.Max {
					return nil, fmt.Errorf(msgInvalidIntMinMax, i, param.Min, param.Max)
				}
			}

			out = append(out, i)

		case "FLOAT":
			var i float64
			if value != "" {
				i2, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return nil, fmt.Errorf(msgInvalidFloat, value)
				}
				i = i2
			}

			if i < param.Min || i > param.Max {
				return nil, fmt.Errorf(msgInvalidFloatMinMax, i, param.Min, param.Max)
			}

			out = append(out, i)

		case "STRING":
			out = append(out, value)

		case "REELS":
			reels := make([]int, 0, 8)
			for iy := ix; iy < len(params); iy++ {
				i, err := strconv.ParseInt(params[iy], 10, 64)
				if err != nil || i <= 0 || i > 10 {
					return nil, fmt.Errorf(msgInvalidReel, params[iy])
				}
				reels = append(reels, int(i))
			}

			out = append(out, reels)
			ix = len(cond.Parameters) // make sure REELS is always processed last!
		}
	}

	return out, nil
}

func checkSymbol(value string, symbols *slots.SymbolSet) (int, error) {
	switch value {
	case "":
		return 1, nil
	case "wild", "scatter", "bomb", "hero":
		return findSymbolName(value, symbols)
	}

	if i, err := strconv.ParseInt(value, 10, 64); err == nil {
		if symbols.GetSymbol(utils.Index(i)) == nil {
			return 0, fmt.Errorf(msgInvalidSymbolCode, value)
		}
		return int(i), nil
	}

	if i, err := findSymbolName(value, symbols); err == nil {
		return i, nil
	}
	return 0, fmt.Errorf(msgInvalidSymbolCode, value)
}

func findSymbolName(value string, symbols *slots.SymbolSet) (int, error) {
	high := symbols.GetMaxSymbolID()
	for ix := utils.Index(1); ix <= high; ix++ {
		if symbol := symbols.GetSymbol(ix); symbol != nil {
			if symbol.Name() == value || symbol.Resource() == value {
				return int(symbol.ID()), nil
			}
		}
	}

	var i int
	switch value {
	case "wild":
		i = findSymbolKind(slots.Wild, symbols)
	case "scatter":
		i = findSymbolKind(slots.Scatter, symbols)
	case "bomb":
		i = findSymbolKind(slots.Bomb, symbols)
	case "hero":
		i = findSymbolKind(slots.Hero, symbols)
	}

	if i == 0 {
		return 0, fmt.Errorf(msgInvalidSymbolCode, value)
	}
	return i, nil
}

func findSymbolKind(kind slots.SymbolKind, symbols *slots.SymbolSet) int {
	high := symbols.GetMaxSymbolID()
	for ix := utils.Index(1); ix <= high; ix++ {
		if symbol := symbols.GetSymbol(ix); symbol != nil {
			if symbol.Kind() == kind {
				return int(symbol.ID())
			}
		}
	}
	return 0
}

const (
	msgInvalidInt         = "invalid integer parameter [%s]"
	msgInvalidIntMinMax   = "invalid integer parameter [%d]; must be between %.0f and %.0f"
	msgInvalidFloat       = "invalid float parameter [%s]"
	msgInvalidFloatMinMax = "invalid float parameter [%.2f]; must be between %.2f and %.2f"
	msgInvalidSymbolCode  = "invalid symbol code [%s]"
	msgInvalidReel        = "invalid reel number [%s]"
)

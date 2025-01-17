package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func (a *ReviseAction) doReplaceSymbols(spin *Spin) bool {
	count := spin.CountSymbolInReels(a.symbol, a.detectReels)
	if count == 0 {
		return false
	}

	var repl bool
	var offset, max int

	maxID := spin.symbols.maxID

	replacements := make(utils.Indexes, 0, 100)
	for id := utils.Index(1); id <= maxID; id++ {
		if id == a.symbol {
			continue
		}
		ok := true
		for _, s := range a.replSymbols {
			if s == id {
				ok = false
				break
			}
		}
		if ok {
			replacements = append(replacements, id)
		}
	}

	valid := func(id utils.Index) bool {
		if a.genAllowDupes {
			return true
		}
		for ix := offset; ix < max; ix++ {
			if spin.indexes[ix] == id {
				return false
			}
		}
		return true
	}

	repL := len(replacements)

	for ix, symbol := range a.replSymbols {
		if spin.TestChance2(a.ModifyChance(a.replChances[ix], spin)) {
			for _, reel := range a.replReels {
				max = int(reel) * spin.rowCount
				for offset = max - spin.rowCount; offset < max; offset++ {
					if spin.indexes[offset] == symbol {
						id := replacements[spin.prng.IntN(repL)]
						for !valid(id) {
							id = replacements[spin.prng.IntN(repL)]
						}
						spin.indexes[offset] = id
						repl = true
					}
				}
			}
		}
	}

	return repl
}

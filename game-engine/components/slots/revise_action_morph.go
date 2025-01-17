package slots

func (a *ReviseAction) doMorphSymbols(spin *Spin) bool {
	var morphed bool

	// TODO: pick reels at random as now the left-most reels will take precedence when there are multiple symbols to morph!

	for ix := range a.morphChances {
		if spin.TestChance2(a.ModifyChance(a.morphChances[ix], spin)) {
			if len(a.morphReels) > 0 {
				for _, reel := range a.morphReels {
					if a.morphReelSymbols(spin, int(reel-1)) {
						morphed = true
						break
					}
				}
			} else {
				for reel := 0; reel < spin.reelCount; reel++ {
					if a.morphReelSymbols(spin, reel) {
						morphed = true
						break
					}
				}
			}
		} else {
			break
		}
	}

	return morphed
}

func (a *ReviseAction) morphReelSymbols(spin *Spin, reel int) bool {
	rows := spin.rowCount
	allSymbols := len(a.morphFor) == 0
	offset := reel * rows

	for max := offset + rows; offset < max; offset++ {
		id := spin.indexes[offset]
		if allSymbols || a.morphFor.Contains(id) {
			if s := spin.GetSymbol(id); s != nil && s.morphInto != 0 {
				spin.indexes[offset] = s.morphInto
				return true
			}
		}
	}
	return false
}

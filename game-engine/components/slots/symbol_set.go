package slots

import (
	"math"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/interfaces"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// SymbolSet represents the set of symbols for a slot machine.
type SymbolSet struct {
	symbols      Symbols                 // the list of symbols.
	sorted       Symbols                 // contains all symbols indexed by their ID.
	maxID        utils.Index             // highest ID of the symbols.
	haveBombs    bool                    // indicates if there are bomb symbols.
	minPayout    uint8                   // lowest number of consecutive symbols with payout.
	maxPayout    uint8                   // highest number of consecutive symbols with payout.
	bestWildSym  []utils.Index           // symbols with the best payout substituting for n wilds.
	bestWildPay  []float64               // corresponding payout substituting for n wilds.
	bonusWeights utils.WeightedGenerator // weighting for selecting a bonus symbol.
}

// NewSymbolSet instantiates a new symbol set and adds the given symbols to it.
// The function panics if one or more indexes of the given symbols are not unique.
// SymbolSets are immutable once created, so they are safe to use across concurrent go-routines.
func NewSymbolSet(symbols ...*Symbol) *SymbolSet {
	s := &SymbolSet{symbols: symbols}
	return s.init()
}

// SetBonusWeights sets the weighting for selecting a bonus symbol from the list of symbols.
func (s *SymbolSet) SetBonusWeights(weighting utils.WeightedGenerator) *SymbolSet {
	s.bonusWeights = weighting
	return s
}

// GetBonusSymbol returns a random bonus symbol using the bonus symbol weighting with the given PRNG.
func (s *SymbolSet) GetBonusSymbol(prng interfaces.Generator) utils.Index {
	return s.bonusWeights.RandomIndex(prng)
}

// GetSymbol returns the symbol matching the given index or nil if it doesn't exist.
func (s *SymbolSet) GetSymbol(index utils.Index) *Symbol {
	if index > s.maxID {
		return nil
	}
	return s.sorted[index]
}

// GetMaxSymbolID returns the highest symbol id.
func (s *SymbolSet) GetMaxSymbolID() utils.Index {
	return s.maxID
}

// GetBombSymbol returns the first symbol with the bomb kind.
func (s *SymbolSet) GetBombSymbol() *Symbol {
	for ix := range s.symbols {
		symbol := s.symbols[ix]
		if symbol.IsBomb() {
			return symbol
		}
	}
	return nil
}

func (s *SymbolSet) init() *SymbolSet {
	s.maxID = 0
	s.minPayout = math.MaxUint8
	s.maxPayout = 0
	s.bestWildSym = utils.PurgeIndexes(s.bestWildSym, 8)
	s.bestWildPay = utils.PurgeFloats(s.bestWildPay, 8)
	s.bonusWeights = nil

	for _, symbol := range s.symbols {
		if symbol.id > s.maxID {
			s.maxID = symbol.id
		}
		if symbol.IsBomb() {
			s.haveBombs = true
		}
		if !symbol.isScatter {
			for ix, payout := range symbol.payouts {
				if payout > 0 {
					reels := uint8(ix) + 1
					if reels < s.minPayout {
						s.minPayout = reels
					}
					if reels > s.maxPayout {
						s.maxPayout = reels
					}
				}
			}
		}
	}

	max := s.maxID + 1
	s.sorted = make(Symbols, max)
	for _, symbol := range s.symbols {
		id := symbol.id
		if s.sorted[id] != nil {
			panic(consts.MsgDuplicateSymbolIndex)
		}
		s.sorted[id] = symbol
	}

	if s.minPayout < s.maxPayout {
		maxReel := s.maxPayout + 1
		s.bestWildSym = utils.PurgeIndexes(s.bestWildSym, int(maxReel))[:maxReel]
		s.bestWildPay = utils.PurgeFloats(s.bestWildPay, int(maxReel))[:maxReel]
		for reels := uint8(0); reels < maxReel; reels++ {
			s.bestWildSym[reels] = utils.MaxIndex
			s.bestWildPay[reels] = 0
			if reels >= 2 {
				for ix := len(s.sorted) - 1; ix >= 0; ix-- {
					if symbol := s.sorted[ix]; symbol != nil {
						if p := symbol.Payout(reels); p > s.bestWildPay[reels] {
							s.bestWildSym[reels] = symbol.id
							s.bestWildPay[reels] = p
						}
					}
				}
			}
		}
	}

	return s
}

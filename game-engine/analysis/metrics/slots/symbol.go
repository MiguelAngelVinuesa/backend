package slots

import (
	"reflect"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// NewSymbol instantiates a new symbol metrics from the memory pool.
func NewSymbol(id utils.Index, name, resource string, reelCount int) *Symbol {
	s := symbolPool.Acquire().(*Symbol)

	s.ID = id
	s.Name = name
	s.Resource = resource
	s.TotalReels = utils.PurgeUInt64s(s.TotalReels, reelCount)[:reelCount]
	s.FirstReels = utils.PurgeUInt64s(s.FirstReels, reelCount)[:reelCount]
	s.SecondReels = utils.PurgeUInt64s(s.SecondReels, reelCount)[:reelCount]
	s.FreeReels = utils.PurgeUInt64s(s.FreeReels, reelCount)[:reelCount]
	s.FreeSecondReels = utils.PurgeUInt64s(s.FreeSecondReels, reelCount)[:reelCount]
	s.Payouts = utils.PurgeUInt64s(s.Payouts, 16)

	s.ResetData()
	return s
}

// Merge merges with the given input.
func (s *Symbol) Merge(other *Symbol) {
	if len(s.TotalReels) != len(other.TotalReels) {
		panic(consts.MsgAnalysisNonMatchingRounds)
	}

	s.TotalCount += other.TotalCount
	s.FirstCount += other.FirstCount
	s.SecondCount += other.SecondCount
	s.FreeCount += other.FreeCount
	s.FreeSecondCount += other.FreeSecondCount
	s.BonusCount += other.BonusCount
	s.StickyCount += other.StickyCount
	s.SuperCount += other.SuperCount

	for ix := 0; ix < len(s.TotalReels); ix++ {
		s.TotalReels[ix] += other.TotalReels[ix]
		s.FirstReels[ix] += other.FirstReels[ix]
		s.SecondReels[ix] += other.SecondReels[ix]
		s.FreeReels[ix] += other.FreeReels[ix]
		s.FreeSecondReels[ix] += other.FreeSecondReels[ix]
	}

	for ix := range other.Payouts {
		s.AddPayout(uint8(ix+2), other.Payouts[ix])
	}
}

// IncreaseFirst increases the symbol metrics for the first spin.
func (s *Symbol) IncreaseFirst(reel int) {
	s.TotalCount++
	s.TotalReels[reel]++
	s.FirstCount++
	s.FirstReels[reel]++
}

// IncreaseSecond increases the symbol metrics for the second spin.
func (s *Symbol) IncreaseSecond(reel int) {
	s.TotalCount++
	s.TotalReels[reel]++
	s.SecondCount++
	s.SecondReels[reel]++
}

// IncreaseFree increases the symbol metrics for a free spin.
func (s *Symbol) IncreaseFree(reel int) {
	s.TotalCount++
	s.TotalReels[reel]++
	s.FreeCount++
	s.FreeReels[reel]++
}

// IncreaseSecondFree increases the symbol metrics for a free second spin.
func (s *Symbol) IncreaseSecondFree(reel int) {
	s.TotalCount++
	s.TotalReels[reel]++
	s.FreeSecondCount++
	s.FreeSecondReels[reel]++
}

// IncreaseBonus increases the count for bonus symbol.
func (s *Symbol) IncreaseBonus() {
	s.BonusCount++
}

// IncreaseSticky increases the count for sticky symbol.
func (s *Symbol) IncreaseSticky() {
	s.StickyCount++
}

// IncreaseSuper increases the count for super symbol.
func (s *Symbol) IncreaseSuper() {
	s.SuperCount++
}

// AddPayout adds a payout length for the symbol.
func (s *Symbol) AddPayout(l uint8, count uint64) {
	if l < 2 {
		return
	}

	l -= 2
	need := int(l) - len(s.Payouts) + 1
	for ix := 0; ix < need; ix++ {
		s.Payouts = append(s.Payouts, 0)
	}

	s.Payouts[l] += count
}

// Symbol represents the metrics for a symbol.
type Symbol struct {
	ID              utils.Index `json:"id"`
	TotalCount      uint64      `json:"totalCount"`
	FirstCount      uint64      `json:"firstCount"`
	SecondCount     uint64      `json:"secondCount,omitempty"`
	FreeCount       uint64      `json:"freeCount"`
	FreeSecondCount uint64      `json:"freeSecondCount,omitempty"`
	BonusCount      uint64      `json:"bonusCount,omitempty"`
	StickyCount     uint64      `json:"stickyCount,omitempty"`
	SuperCount      uint64      `json:"superCount,omitempty"`
	TotalReels      []uint64    `json:"totalReels,omitempty"`
	FirstReels      []uint64    `json:"firstReels,omitempty"`
	SecondReels     []uint64    `json:"secondReels,omitempty"`
	FreeReels       []uint64    `json:"freeReels,omitempty"`
	FreeSecondReels []uint64    `json:"freeSecondReels,omitempty"`
	Payouts         []uint64    `json:"payouts,omitempty"`
	Name            string      `json:"name"`
	Resource        string      `json:"resource,omitempty"`
	pool.Object
}

var symbolPool = pool.NewProducer(func() (pool.Objecter, func()) {
	s := &Symbol{
		TotalReels:      make([]uint64, 0, 8),
		FirstReels:      make([]uint64, 0, 8),
		SecondReels:     make([]uint64, 0, 8),
		FreeReels:       make([]uint64, 0, 8),
		FreeSecondReels: make([]uint64, 0, 8),
		Payouts:         make([]uint64, 0, 16),
	}
	return s, s.reset
})

// reset clears the symbol metrics.
func (s *Symbol) reset() {
	if s != nil {
		s.ResetData()

		s.ID = 0
		s.Name = ""
		s.Resource = ""

		s.TotalReels = s.TotalReels[:0]
		s.FirstReels = s.FirstReels[:0]
		s.SecondReels = s.SecondReels[:0]
		s.FreeReels = s.FreeReels[:0]
		s.FreeSecondReels = s.FreeSecondReels[:0]
		s.Payouts = s.Payouts[:0]
	}
}

// ResetData resets the symbol metrics to initial state.
func (s *Symbol) ResetData() {
	s.TotalCount = 0
	s.FirstCount = 0
	s.SecondCount = 0
	s.FreeCount = 0
	s.FreeSecondCount = 0
	s.BonusCount = 0
	s.StickyCount = 0
	s.SuperCount = 0

	copy(s.TotalReels, consts.ClearUint64s)
	copy(s.FirstReels, consts.ClearUint64s)
	copy(s.SecondReels, consts.ClearUint64s)
	copy(s.FreeReels, consts.ClearUint64s)
	copy(s.FreeSecondReels, consts.ClearUint64s)
	copy(s.Payouts, consts.ClearUint64s)
}

// Equals is used internally for unit tests!
func (s *Symbol) Equals(other *Symbol) bool {
	return s.ID == other.ID &&
		s.TotalCount == other.TotalCount &&
		s.FirstCount == other.FirstCount &&
		s.SecondCount == other.SecondCount &&
		s.FreeCount == other.FreeCount &&
		s.FreeSecondCount == other.FreeSecondCount &&
		s.BonusCount == other.BonusCount &&
		s.StickyCount == other.StickyCount &&
		s.SuperCount == other.SuperCount &&
		s.Name == other.Name &&
		s.Resource == other.Resource &&
		reflect.DeepEqual(s.TotalReels, other.TotalReels) &&
		reflect.DeepEqual(s.FirstReels, other.FirstReels) &&
		reflect.DeepEqual(s.SecondReels, other.SecondReels) &&
		reflect.DeepEqual(s.FreeReels, other.FreeReels) &&
		reflect.DeepEqual(s.FreeSecondReels, other.FreeSecondReels) &&
		reflect.DeepEqual(s.Payouts, other.Payouts)

}

// Symbols is a convenience type for a slice of symbols.
type Symbols []*Symbol

// ReleaseSymbols release the symbols and returns an empty slice.
func ReleaseSymbols(list Symbols) Symbols {
	if list == nil {
		return nil
	}
	for ix := range list {
		if s := list[ix]; s != nil {
			s.Release()
			list[ix] = nil
		}
	}
	return list[:0]
}

// PurgeSymbols returns the input cleared to zero length or a new slice if it doesn't have the requested capacity.
func PurgeSymbols(list Symbols, capacity int) Symbols {
	list = ReleaseSymbols(list)
	if cap(list) < capacity {
		return make(Symbols, capacity)
	}
	return list
}

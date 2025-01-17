package slots

import (
	"reflect"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

const (
	maxCount = 49 // maximum for payout count; increase if a game grid with scatter payouts exceeds this size.
)

// NewScatterPayout instantiates a new scatter payout metrics from the memory pool.
func NewScatterPayout(maxSymbol utils.Index) *ScatterPayout {
	s := scatterPayoutPool.Acquire().(*ScatterPayout)
	s.SetMaxSymbol(maxSymbol)
	return s
}

func (s *ScatterPayout) SetMaxSymbol(maxSymbol utils.Index) {
	m := int(maxSymbol + 1)
	s.Symbols = utils.PurgeUInt64s(s.Symbols, m)[:m]

	if cap(s.SymbolLengths) < m {
		s.SymbolLengths = makeSymbolLengths(m)
	} else {
		s.SymbolLengths = s.SymbolLengths[:m]
	}
}

// Merge merges with the given input.
func (s *ScatterPayout) Merge(other *ScatterPayout) {
	if len(s.Symbols) != len(other.Symbols) || len(other.Lengths) > cap(s.Lengths) {
		panic(consts.MsgAnalysisNonMatchingRounds)
	}

	s.Count += other.Count

	s.Payouts.Merge(other.Payouts)

	for ix := range s.Symbols {
		s.Symbols[ix] += other.Symbols[ix]
	}

	l := len(other.Lengths)
	if l > len(s.Lengths) {
		s.Lengths = s.Lengths[:l]
	}
	for ix := range s.Lengths {
		if ix < l {
			s.Lengths[ix] += other.Lengths[ix]
		}
	}

	for ix := range s.SymbolLengths {
		l1, l2 := s.SymbolLengths[ix], other.SymbolLengths[ix]
		if len(l2) > len(l1) {
			l1 = l1[:len(l2)]
			s.SymbolLengths[ix] = l1
		}
		for iy := range l2 {
			l1[iy] += l2[iy]
		}
	}
}

// Increase updates the scatter payout metrics with the given input.
// The function panics if symbol > maxSymbol or count > maxCount.
func (s *ScatterPayout) Increase(symbol utils.Index, count uint8, payout float64) *ScatterPayout {
	s.Count++
	s.Payouts.Increase(payout)
	s.Symbols[symbol]++

	if int(count) >= len(s.Lengths) {
		s.Lengths = s.Lengths[:count+1]
	}
	s.Lengths[count]++

	if sl := s.SymbolLengths[symbol]; sl != nil {
		if int(count) >= len(sl) {
			sl = sl[:count+1]
			s.SymbolLengths[symbol] = sl
		}
		sl[count]++
	}

	return s
}

// ScatterPayout represents the metrics for a scatter payout.
type ScatterPayout struct {
	Count         uint64         `json:"count,omitempty"`
	Payouts       *MinMaxFloat64 `json:"payouts,omitempty"`
	Symbols       []uint64       `json:"symbols,omitempty"`
	Lengths       []uint64       `json:"lengths,omitempty"`
	SymbolLengths [][]uint64     `json:"symbolLengths,omitempty"`
	pool.Object
}

var scatterPayoutPool = pool.NewProducer(func() (pool.Objecter, func()) {
	s := &ScatterPayout{
		Symbols:       make([]uint64, 0, 16),
		Lengths:       make([]uint64, 0, maxCount+1),
		SymbolLengths: makeSymbolLengths(16),
		Payouts:       AcquireMinMaxFloat64(2),
	}
	return s, s.reset
})

func makeSymbolLengths(m int) [][]uint64 {
	out := make([][]uint64, m)
	for ix := range out {
		out[ix] = make([]uint64, 0, maxCount+1)
	}
	return out
}

func clearSymbolLengths(s [][]uint64) [][]uint64 {
	for ix := range s {
		clear(s[ix])
		s[ix] = s[ix][:0]
	}
	return s[:0]
}

// reset clears the scatter payout metrics.
func (s *ScatterPayout) reset() {
	if s != nil {
		s.ResetData()

		s.Symbols = s.Symbols[:0]
		s.Lengths = s.Lengths[:0]
		s.SymbolLengths = s.SymbolLengths[:0]
	}
}

// ResetData resets the scatter payout metrics to initial state.
func (s *ScatterPayout) ResetData() {
	s.Count = 0

	s.Payouts.ResetData()

	clear(s.Symbols)
	clear(s.Lengths)

	s.Lengths = s.Lengths[:0]
	s.SymbolLengths = clearSymbolLengths(s.SymbolLengths)
}

// Equals is used internally for unit tests!
func (s *ScatterPayout) Equals(other *ScatterPayout) bool {
	return s.Count == other.Count &&
		s.Payouts.Equals(other.Payouts) &&
		reflect.DeepEqual(s.Symbols, other.Symbols) &&
		reflect.DeepEqual(s.Lengths, other.Lengths) &&
		reflect.DeepEqual(s.SymbolLengths, other.SymbolLengths)
}

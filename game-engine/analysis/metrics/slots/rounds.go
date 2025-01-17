package slots

import (
	"reflect"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// NewRounds instantiates a new rounds metrics from the memory pool.
func NewRounds(symbols *slots.SymbolSet) *Rounds {
	r := roundsPool.Acquire().(*Rounds)

	r.symbols = symbols
	r.FreeSpinRounds = utils.PurgeUInt64s(r.FreeSpinRounds, 32)
	r.RefillRounds = utils.PurgeUInt64s(r.RefillRounds, 32)
	r.SuperRounds = utils.PurgeUInt64s(r.SuperRounds, 32)

	max, l := 16, 0
	if symbols != nil {
		max, l = 16, int(symbols.GetMaxSymbolID())+1
		if l > max {
			max = l
		}
	}
	r.SymbolsUsed = utils.PurgeUInt64s(r.SymbolsUsed, max)[:l]
	r.SymbolsNoFree = utils.PurgeUInt64s(r.SymbolsNoFree, max)[:l]
	r.SymbolsFree = utils.PurgeUInt64s(r.SymbolsFree, max)[:l]

	return r
}

// Merge merges with the given input.
func (r *Rounds) Merge(other *Rounds) {
	if r.symbols != other.symbols {
		panic("merge rounds with non-matching symbols")
	}

	r.Count += other.Count

	r.Bets.Merge(other.Bets)
	r.BetsNoFree.Merge(other.BetsNoFree)
	r.BetsFree.Merge(other.BetsFree)
	r.Wins.Merge(other.Wins)
	r.WinsNoFree.Merge(other.WinsNoFree)
	r.WinsFree.Merge(other.WinsFree)
	r.FreeSpins.Merge(other.FreeSpins)
	r.RefillSpins.Merge(other.RefillSpins)
	r.SuperSpins.Merge(other.SuperSpins)

	l1, l2 := len(r.FreeSpinRounds), len(other.FreeSpinRounds)
	if l2 > l1 {
		l1 = l2
	}
	r.FreeSpinRounds = ExpandUInt64s(r.FreeSpinRounds, l1)[:l1]

	l1, l2 = len(r.RefillRounds), len(other.RefillRounds)
	if l2 > l1 {
		l1 = l2
	}
	r.RefillRounds = ExpandUInt64s(r.RefillRounds, l1)[:l1]

	l1, l2 = len(r.SuperRounds), len(other.SuperRounds)
	if l2 > l1 {
		l1 = l2
	}
	r.SuperRounds = ExpandUInt64s(r.SuperRounds, l1)[:l1]

	for ix := range other.FreeSpinRounds {
		r.FreeSpinRounds[ix] += other.FreeSpinRounds[ix]
	}
	for ix := range other.RefillRounds {
		r.RefillRounds[ix] += other.RefillRounds[ix]
	}
	for ix := range other.SuperRounds {
		r.SuperRounds[ix] += other.SuperRounds[ix]
	}

	for ix := range r.SymbolsUsed {
		r.SymbolsUsed[ix] += other.SymbolsUsed[ix]
		r.SymbolsNoFree[ix] += other.SymbolsNoFree[ix]
		r.SymbolsFree[ix] += other.SymbolsFree[ix]
	}
}

// NewRound records a new round in the metrics.
func (r *Rounds) NewRound(bet, win int64, freeSpins, refillSpins, superSpins uint64, result results.Results) {
	r.Count++

	r.Bets.Increase(bet)
	r.Wins.Increase(win)

	if freeSpins > 0 {
		r.BetsFree.Increase(bet)
		r.WinsFree.Increase(win)
	} else {
		r.BetsNoFree.Increase(bet)
		r.WinsNoFree.Increase(win)
	}

	r.FreeSpins.Increase(freeSpins)
	r.RefillSpins.Increase(refillSpins)
	r.SuperSpins.Increase(superSpins)

	l := len(r.FreeSpinRounds)
	if int(freeSpins) >= l {
		l = int(freeSpins) + 1
	}
	r.FreeSpinRounds = ExpandUInt64s(r.FreeSpinRounds, l)[:l]
	r.FreeSpinRounds[freeSpins]++

	l = len(r.RefillRounds)
	if int(refillSpins) >= l {
		l = int(refillSpins) + 1
	}
	r.RefillRounds = ExpandUInt64s(r.RefillRounds, l)[:l]
	r.RefillRounds[refillSpins]++

	l = len(r.SuperRounds)
	if int(superSpins) >= l {
		l = int(superSpins) + 1
	}
	r.SuperRounds = ExpandUInt64s(r.SuperRounds, l)[:l]
	r.SuperRounds[superSpins]++

	used := make([]bool, 100)
	free := false

	for ix := range result {
		res := result[ix]
		if spin, ok := res.Data.(*slots.SpinResult); ok {
			if spin.Kind() == slots.FreeSpin || spin.Kind() == slots.FirstFreeSpin || spin.Kind() == slots.SecondFreeSpin {
				free = true
			}
			for _, id := range spin.Initial() {
				used[id] = true
			}
		}
	}

	l = len(r.SymbolsUsed)
	for ix, u := range used {
		if u && ix < l {
			r.SymbolsUsed[ix]++
		}
	}

	if free {
		for ix := range result {
			res := result[ix]
			if spin, ok := res.Data.(*slots.SpinResult); ok {
				if spin.Kind() == slots.FreeSpin || spin.Kind() == slots.FirstFreeSpin || spin.Kind() == slots.SecondFreeSpin {
					for _, id := range spin.Initial() {
						used[id] = true
					}
				}
			}
		}

		for ix, u := range used {
			if u && ix < l {
				r.SymbolsFree[ix]++
			}
		}
	} else {
		for ix := range result {
			res := result[ix]
			if spin, ok := res.Data.(*slots.SpinResult); ok {
				if spin.Kind() != slots.FreeSpin && spin.Kind() == slots.FirstFreeSpin && spin.Kind() == slots.SecondFreeSpin {
					for _, id := range spin.Initial() {
						used[id] = true
					}
				}
			}
		}

		for ix, u := range used {
			if u && ix < l {
				r.SymbolsNoFree[ix]++
			}
		}
	}
}

// Rounds contains the metrics for bets, wins and free spins.
type Rounds struct {
	Count          uint64 `json:"count"`
	symbols        *slots.SymbolSet
	Bets           *MinMaxInt64  `json:"bets,omitempty"`
	BetsNoFree     *MinMaxInt64  `json:"betsNoFree,omitempty"`
	BetsFree       *MinMaxInt64  `json:"betsFree,omitempty"`
	Wins           *MinMaxInt64  `json:"wins,omitempty"`
	WinsNoFree     *MinMaxInt64  `json:"winsNoFree,omitempty"`
	WinsFree       *MinMaxInt64  `json:"winsFree,omitempty"`
	FreeSpins      *MinMaxUInt64 `json:"freeSpins,omitempty"`
	RefillSpins    *MinMaxUInt64 `json:"refillSpins,omitempty"`
	SuperSpins     *MinMaxUInt64 `json:"superSpins,omitempty"`
	FreeSpinRounds []uint64      `json:"freeSpinRounds,omitempty"`
	RefillRounds   []uint64      `json:"refillRounds,omitempty"`
	SuperRounds    []uint64      `json:"superRounds,omitempty"`
	SymbolsUsed    []uint64      `json:"symbolsUsed,omitempty"`
	SymbolsNoFree  []uint64      `json:"symbolsUsedWithoutFree,omitempty"`
	SymbolsFree    []uint64      `json:"symbolsUsedWithFree,omitempty"`
	pool.Object
}

var roundsPool = pool.NewProducer(func() (pool.Objecter, func()) {
	r := &Rounds{
		Bets:           AcquireMinMaxInt64(),
		BetsNoFree:     AcquireMinMaxInt64(),
		BetsFree:       AcquireMinMaxInt64(),
		Wins:           AcquireMinMaxInt64(),
		WinsNoFree:     AcquireMinMaxInt64(),
		WinsFree:       AcquireMinMaxInt64(),
		FreeSpins:      AcquireMinMaxUInt64(),
		RefillSpins:    AcquireMinMaxUInt64(),
		SuperSpins:     AcquireMinMaxUInt64(),
		FreeSpinRounds: make([]uint64, 0, 32),
		RefillRounds:   make([]uint64, 0, 32),
		SuperRounds:    make([]uint64, 0, 32),
		SymbolsUsed:    make([]uint64, 0, 16),
		SymbolsNoFree:  make([]uint64, 0, 16),
		SymbolsFree:    make([]uint64, 0, 16),
	}
	return r, r.reset
})

// reset clears the round metrics.
func (r *Rounds) reset() {
	if r != nil {
		r.ResetData()

		r.FreeSpinRounds = r.FreeSpinRounds[:0]
		r.RefillRounds = r.RefillRounds[:0]
		r.SuperRounds = r.SuperRounds[:0]
		r.SymbolsUsed = r.SymbolsUsed[:0]
		r.SymbolsNoFree = r.SymbolsNoFree[:0]
		r.SymbolsFree = r.SymbolsFree[:0]
	}
}

// ResetData resets the metrics to initial state.
func (r *Rounds) ResetData() {
	r.Count = 0

	r.Bets.ResetData()
	r.BetsNoFree.ResetData()
	r.BetsFree.ResetData()
	r.Wins.ResetData()
	r.WinsNoFree.ResetData()
	r.WinsFree.ResetData()
	r.FreeSpins.ResetData()
	r.RefillSpins.ResetData()
	r.SuperSpins.ResetData()

	copy(r.FreeSpinRounds, consts.ClearUint64s)
	copy(r.RefillRounds, consts.ClearUint64s)
	copy(r.SuperRounds, consts.ClearUint64s)
	copy(r.SymbolsUsed, consts.ClearUint64s)
	copy(r.SymbolsNoFree, consts.ClearUint64s)
	copy(r.SymbolsFree, consts.ClearUint64s)
}

// Equals is used internally for unit-tests!
func (r *Rounds) Equals(other *Rounds) bool {
	return r.Count == other.Count &&
		r.Bets.Equals(other.Bets) &&
		r.BetsNoFree.Equals(other.BetsNoFree) &&
		r.BetsFree.Equals(other.BetsFree) &&
		r.Wins.Equals(other.Wins) &&
		r.WinsNoFree.Equals(other.WinsNoFree) &&
		r.WinsFree.Equals(other.WinsFree) &&
		r.FreeSpins.Equals(other.FreeSpins) &&
		r.RefillSpins.Equals(other.RefillSpins) &&
		r.SuperSpins.Equals(other.SuperSpins) &&
		reflect.DeepEqual(r.FreeSpinRounds, other.FreeSpinRounds) &&
		reflect.DeepEqual(r.RefillRounds, other.RefillRounds) &&
		reflect.DeepEqual(r.SuperRounds, other.SuperRounds) &&
		reflect.DeepEqual(r.SymbolsUsed, other.SymbolsUsed) &&
		reflect.DeepEqual(r.SymbolsNoFree, other.SymbolsNoFree) &&
		reflect.DeepEqual(r.SymbolsFree, other.SymbolsFree)
}

func ExpandUInt64s(input []uint64, capacity int) []uint64 {
	if cap(input) < capacity {
		output := make([]uint64, capacity)
		copy(output, input)
		return output
	}
	return input[:capacity]
}

package slots

import (
	"reflect"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// NewPayline instantiates a new payline metrics from the memory pool.
func NewPayline(id int, maxSymbol utils.Index, rowMap utils.UInt8s) *Payline {
	p := paylinePool.Acquire().(*Payline)
	p.ID = id
	p.RowMap = utils.CopyPurgeUInt8s(rowMap, p.RowMap, 0)

	m := maxSymbol + 1
	p.Symbols = utils.PurgeUInt64s(p.Symbols, int(m))[:m]
	p.Lengths = utils.PurgeUInt64s(p.Lengths, maxReels+1)

	if cap(p.SymbolLengths) < int(m) {
		p.SymbolLengths = makeSymbolLengths(int(m))
	} else {
		p.SymbolLengths = p.SymbolLengths[:m]
	}

	return p
}

// Merge merges with the given metrics.
func (p *Payline) Merge(other *Payline) {
	if len(p.Symbols) != len(other.Symbols) || len(other.Lengths) > cap(p.Lengths) {
		panic(consts.MsgAnalysisNonMatchingRounds)
	}

	p.Count += other.Count

	p.Payouts.Merge(other.Payouts)

	for ix := range p.Symbols {
		p.Symbols[ix] += other.Symbols[ix]
	}

	l := len(other.Lengths)
	if l > len(p.Lengths) {
		p.Lengths = p.Lengths[:l]
	}
	for ix := range p.Lengths {
		if ix < l {
			p.Lengths[ix] += other.Lengths[ix]
		}
	}

	for ix := range p.SymbolLengths {
		l1, l2 := p.SymbolLengths[ix], other.SymbolLengths[ix]
		if len(l2) > len(l1) {
			l1 = l1[:len(l2)]
			p.SymbolLengths[ix] = l1
		}
		for iy := range l2 {
			l1[iy] += l2[iy]
		}
	}
}

// Increase updates the payline metrics with the given input.
// The function panics if symbol > maxSymbol or count > maxReels.
func (p *Payline) Increase(symbol utils.Index, count uint8, payout float64) *Payline {
	p.Count++
	p.Payouts.Increase(payout)
	p.Symbols[symbol]++

	if int(count) >= len(p.Lengths) {
		p.Lengths = p.Lengths[:count+1]
	}
	p.Lengths[count]++

	if sl := p.SymbolLengths[symbol]; sl != nil {
		if int(count) >= len(sl) {
			sl = sl[:count+1]
			p.SymbolLengths[symbol] = sl
		}
		sl[count]++
	}

	return p
}

// Payline represents the metrics for a payline.
type Payline struct {
	ID            int            `json:"id,omitempty"`
	Count         uint64         `json:"count,omitempty"`
	Payouts       *MinMaxFloat64 `json:"payouts,omitempty"`
	RowMap        utils.UInt8s   `json:"rowMap,omitempty"`
	Symbols       []uint64       `json:"symbols,omitempty"`
	Lengths       []uint64       `json:"lengths,omitempty"`
	SymbolLengths [][]uint64     `json:"symbolsX,omitempty"`
	pool.Object
}

var paylinePool = pool.NewProducer(func() (pool.Objecter, func()) {
	p := &Payline{
		Payouts:       AcquireMinMaxFloat64(1),
		Symbols:       make([]uint64, 0, 16),
		Lengths:       make([]uint64, 0, maxReels+1),
		SymbolLengths: makeSymbolLengths(16),
	}
	return p, p.reset
})

// reset clear the payout metrics.
func (p *Payline) reset() {
	if p != nil {
		p.ResetData()

		p.ID = 0
		p.RowMap = nil
	}
}

// ResetData resets the payline metrics to initial state.
func (p *Payline) ResetData() {
	p.Count = 0

	p.Payouts.ResetData()

	clear(p.Symbols)
	clear(p.Lengths)

	p.Symbols = p.Symbols[:0]
	p.Lengths = p.Lengths[:0]
	p.SymbolLengths = clearSymbolLengths(p.SymbolLengths)
}

// Equals is used internally for unit tests!
func (p *Payline) Equals(other *Payline) bool {
	return p.ID == other.ID &&
		p.Count == other.Count &&
		p.Payouts.Equals(other.Payouts) &&
		reflect.DeepEqual(p.RowMap, other.RowMap) &&
		reflect.DeepEqual(p.Symbols, other.Symbols) &&
		reflect.DeepEqual(p.Lengths, other.Lengths) &&
		reflect.DeepEqual(p.SymbolLengths, other.SymbolLengths)
}

// Paylines is a convenience type for a slice of paylines.
type Paylines []*Payline

// ReleasePaylines releases all paylines and returns an empty slice.
func ReleasePaylines(list Paylines) Paylines {
	if list == nil {
		return nil
	}
	for ix := range list {
		if p := list[ix]; p != nil {
			p.Release()
			list[ix] = nil
		}
	}
	return list[:0]
}

// PurgePaylines returns the input cleared to zero length or a new slice if it doesn't have the requested capacity.
func PurgePaylines(list Paylines, capacity int) Paylines {
	list = ReleasePaylines(list)
	if cap(list) < capacity {
		return make(Paylines, capacity)
	}
	return list
}

const (
	maxReels = 15 // maximum possible count; increase when a game has more reels than this.
)

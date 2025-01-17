package slots

import (
	"fmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	analyse "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/analysis/metrics/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// NewPayouts instantiates a new payout metrics from the memory pool.
func NewPayouts(rowCount int, maxSymbol utils.Index, paylines []*slots.Payline, noPaylines bool) *Payouts {
	p := payoutsProvider.Acquire().(*Payouts)
	p.maxSymbol = maxSymbol
	p.rowCount = rowCount
	p.noPaylines = noPaylines

	p.WildPayouts.SetMaxSymbol(maxSymbol)
	p.ScatterPayouts.SetMaxSymbol(maxSymbol)
	p.BonusPayouts.SetMaxSymbol(maxSymbol)
	p.SuperPayouts.SetMaxSymbol(maxSymbol)
	p.OtherPayouts.SetMaxSymbol(maxSymbol)

	if len(paylines) > 0 {
		var m uint8
		for ix := range paylines {
			if l := paylines[ix]; l.ID() > m {
				m = l.ID()
			}
		}
		m++

		p.Paylines = analyse.PurgePaylines(p.Paylines, int(m))[:m]
		for ix := range paylines {
			l := paylines[ix]
			p.Paylines[l.ID()] = analyse.NewPayline(int(l.ID()), maxSymbol, l.RowMap())
		}
	}

	return p
}

// Merge merges with the given input.
func (p *Payouts) Merge(other *Payouts) {
	if p.rowCount != other.rowCount || p.maxSymbol != other.maxSymbol || len(p.Paylines) != len(other.Paylines) {
		panic(consts.MsgAnalysisNonMatchingRounds)
	}

	p.Count += other.Count
	p.Total += other.Total

	for ix := range p.Paylines {
		l := p.Paylines[ix]
		if l != nil {
			l.Merge(other.Paylines[ix])
		}
	}

	for k, v := range other.AllPaylines {
		l, ok := p.AllPaylines[k]
		if !ok {
			l = analyse.NewPayline(v.ID, p.maxSymbol, v.RowMap)
			p.AllPaylines[k] = l
		}
		l.Merge(v)
	}

	p.WildPayouts.Merge(other.WildPayouts)
	p.ScatterPayouts.Merge(other.ScatterPayouts)
	p.BonusPayouts.Merge(other.BonusPayouts)
	p.SuperPayouts.Merge(other.SuperPayouts)
	p.OtherPayouts.Merge(other.OtherPayouts)
}

func (p *Payouts) AvgPayout() float64 {
	return float64(p.Total) / float64(p.Count)
}

func (p *Payouts) analyseRound(win int64) {
	if win == 0 {
		return
	}

	p.Count++
	p.Total += win
}

func (p *Payouts) analyse(result *results.Result) {
	for ix := range result.Payouts {
		pay := result.Payouts[ix].(*slots.SpinPayout)
		switch pay.Kind() {
		case results.SlotWinline:
			if !p.noPaylines {
				if id := pay.PaylineID(); id > 0 {
					p.analysePayline(id, pay, 1)
				} else {
					p.analyseAllPayline(pay, 1)
				}
			}

		case results.SlotWilds:
			p.WildPayouts.Increase(pay.Symbol(), pay.Count(), pay.Total())

		case results.SlotScatters, results.SlotBombScatters:
			p.ScatterPayouts.Increase(pay.Symbol(), pay.Count(), pay.Total())

		case results.SlotBonusSymbol:
			p.BonusPayouts.Increase(pay.Symbol(), pay.Count(), pay.Total())

		case results.SlotSuperShape:
			p.SuperPayouts.Increase(pay.Symbol(), pay.Count(), pay.Total())

		default:
			p.OtherPayouts.Increase(pay.Symbol(), pay.Count(), pay.Total())
		}
	}
}

func (p *Payouts) analysePayline(id uint8, pay *slots.SpinPayout, multiplier float64) {
	l := p.Paylines[id]
	if l == nil {
		panic(fmt.Sprintf("bad input; payline [%d] not configured for analysis: %+v", id, pay))
	}
	l.Increase(pay.Symbol(), pay.Count(), pay.Total()*multiplier)
}

func (p *Payouts) analyseAllPayline(pay *slots.SpinPayout, multiplier float64) {
	id := pay.AllPaylineID()
	if id < 0 {
		panic(fmt.Sprintf("bad input; all paylines id cannot be < zero: %+v", pay))
	}

	l := p.AllPaylines[id]
	if l == nil {
		l = analyse.NewPayline(id, p.maxSymbol, pay.PayRows()[:pay.Count()])
	}
	p.AllPaylines[id] = l.Increase(pay.Symbol(), pay.Count(), pay.Total()*multiplier)
}

// Payouts contains payout metrics.
type Payouts struct {
	noPaylines     bool
	maxSymbol      utils.Index
	rowCount       int
	Count          uint64                   `json:"count,omitempty"`
	Total          int64                    `json:"total,omitempty"`
	WildPayouts    *analyse.ScatterPayout   `json:"wildPayouts,omitempty"`
	ScatterPayouts *analyse.ScatterPayout   `json:"scatterPayouts,omitempty"`
	BonusPayouts   *analyse.ScatterPayout   `json:"bonusSymbolPayouts,omitempty"`
	SuperPayouts   *analyse.ScatterPayout   `json:"superSymbolPayouts,omitempty"`
	OtherPayouts   *analyse.ScatterPayout   `json:"otherPayouts,omitempty"`
	Paylines       analyse.Paylines         `json:"paylines,omitempty"`
	AllPaylines    map[int]*analyse.Payline `json:"allPaylines,omitempty"`
	pool.Object
}

var payoutsProvider = pool.NewProducer(func() (pool.Objecter, func()) {
	p := &Payouts{
		WildPayouts:    analyse.NewScatterPayout(16),
		ScatterPayouts: analyse.NewScatterPayout(16),
		BonusPayouts:   analyse.NewScatterPayout(16),
		SuperPayouts:   analyse.NewScatterPayout(16),
		OtherPayouts:   analyse.NewScatterPayout(16),
		Paylines:       make(analyse.Paylines, 0, 16),
		AllPaylines:    make(map[int]*analyse.Payline, 256),
	}
	return p, p.reset
})

// reset clears the payout metrics.
func (p *Payouts) reset() {
	if p != nil {
		p.ResetData()
		p.Paylines = analyse.ReleasePaylines(p.Paylines)

		for id := range p.AllPaylines {
			if l, ok := p.AllPaylines[id]; ok && l != nil {
				l.Release()
				p.AllPaylines[id] = nil
			}
		}
		clear(p.AllPaylines)
	}
}

// ResetData resets the payout metrics to initial state.
func (p *Payouts) ResetData() {
	p.maxSymbol = 0
	p.rowCount = 0
	p.Count = 0
	p.Total = 0

	p.WildPayouts.ResetData()
	p.ScatterPayouts.ResetData()
	p.BonusPayouts.ResetData()
	p.SuperPayouts.ResetData()
	p.OtherPayouts.ResetData()

	for ix := range p.Paylines {
		if l := p.Paylines[ix]; l != nil {
			l.ResetData()
		}
	}

	for id := range p.AllPaylines {
		if l, ok := p.AllPaylines[id]; ok && l != nil {
			l.ResetData()
		}
	}
}

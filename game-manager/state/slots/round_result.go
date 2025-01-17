package slots

import (
	"fmt"
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/wheel"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/object"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"
)

// AcquireRoundResult instantiates a single spin result from the memory pool.
func AcquireRoundResult(spinSeq int, result *results.Result) *RoundResult {
	r := roundResultPool.Acquire().(*RoundResult)

	r.SpinSeq = spinSeq
	r.AwardedFreeGames = result.AwardedFreeGames
	r.FreeGames = result.FreeGames
	r.TotalPayout = result.Total
	r.Created = time.Now().UTC()

	if s := result.Data; s != nil {
		switch d := s.(type) {
		case *slots.SpinResult:
			r.SpinData = d.Clone().(*slots.SpinResult)
		case *results.InstantBonus:
			r.InstantBonus = d.Clone().(*results.InstantBonus)
		case *results.BonusSelector:
			r.BonusSelector = d.Clone().(*results.BonusSelector)
		case *wheel.BonusWheelResult:
			r.BonusWheel = d.Clone().(*wheel.BonusWheelResult)
		}
	}

	if s := result.State; s != nil {
		ss, ok := s.(*slots.SymbolsState)
		if !ok {
			panic("invalid result: cannot unmarshal symbolsState")
		}
		r.SymbolsState = ss.Clone().(*slots.SymbolsState)
	}

	l := len(result.Payouts)
	top := object.NormalizeSize(l, 8)
	r.Payouts = results.PurgePayouts(r.Payouts, top)[:l]
	for ix := 0; ix < l; ix++ {
		p := result.Payouts[ix]
		r.Payouts[ix] = p.Clone().(results.Payout)
	}

	l = len(result.Penalties)
	top = object.NormalizeSize(l, 8)
	r.Penalties = results.PurgePenalties(r.Penalties, top)[:l]
	for ix := 0; ix < l; ix++ {
		p := result.Penalties[ix]
		r.Penalties[ix] = p.Clone().(results.Penalty)
	}

	return r
}

// AcquireRoundResultFromJSON instantiates a single spin result from the memory pool using the given json data.
func AcquireRoundResultFromJSON(data []byte) *RoundResult {
	r := roundResultPool.Acquire().(*RoundResult)

	dec := zjson.AcquireDecoder(data)
	defer dec.Release()
	if ok := dec.Object(r); ok {
		return r
	}

	r.Release()
	return nil
}

// Result returns the round result as a game engine result.
func (r *RoundResult) Result() *results.Result {
	switch {
	case r.SpinData != nil:
		result := results.AcquireResult(r.SpinData, results.SpinData, r.Payouts...)
		if len(r.Penalties) > 0 {
			result.AddPenalties(r.Penalties...)
		}
		if r.SymbolsState != nil {
			result.AddState(r.SymbolsState)
		}
		result.AwardFreeGames(uint8(r.AwardedFreeGames))
		result.SetFreeGames(r.FreeGames)
		return result

	case r.InstantBonus != nil:
		return results.AcquireResult(r.InstantBonus, results.InstantBonusData, r.Payouts...)

	case r.BonusSelector != nil:
		return results.AcquireResult(r.BonusSelector, results.BonusSelectorData, r.Payouts...)

	case r.BonusWheel != nil:
		return results.AcquireResult(r.BonusWheel, results.BonusWheelData, r.Payouts...)

	default:
		return results.AcquireResult(nil, 0)

	}
}

// EncodeFields implements the zjson encoder interface.
func (r *RoundResult) EncodeFields(enc *zjson.Encoder) {
	enc.IntFieldOpt("spinSeq", r.SpinSeq)
	enc.Int64FieldOpt("balanceBefore", r.BalanceBefore)
	enc.Int64FieldOpt("balanceAfter", r.BalanceAfter)
	enc.Int64FieldOpt("bet", r.Bet)
	enc.Int64FieldOpt("win", r.Win)
	enc.Int64FieldOpt("totalWin", r.TotalWin)
	enc.Int64FieldOpt("progressiveWin", r.ProgressiveWin)
	enc.Int64FieldOpt("bonusWin", r.BonusWin)
	enc.Int64FieldOpt("spinWin", r.SpinWin)
	enc.Uint64FieldOpt("awardedFreeGames", r.AwardedFreeGames)
	enc.Uint64FieldOpt("freeGames", r.FreeGames)
	enc.FloatFieldOpt("totalPayout", r.TotalPayout, 'g', 2)
	enc.FloatFieldOpt("maxPayout", r.MaxPayout, 'g', 2)

	if r.SpinData != nil {
		enc.ObjectField("spinData", r.SpinData)
	}
	if r.InstantBonus != nil {
		enc.ObjectField("instantBonus", r.InstantBonus)
	}
	if r.BonusSelector != nil {
		enc.ObjectField("bonusSelector", r.BonusSelector)
	}
	if r.BonusWheel != nil {
		enc.ObjectField("bonusWheel", r.BonusWheel)
	}
	if r.SymbolsState != nil {
		enc.ObjectField("symbolsState", r.SymbolsState)
	}

	if len(r.Payouts) > 0 {
		enc.StartArrayField("payouts")
		for ix := range r.Payouts {
			enc.Object(r.Payouts[ix])
		}
		enc.EndArray()
	}

	if len(r.Penalties) > 0 {
		enc.StartArrayField("penalties")
		for ix := range r.Penalties {
			enc.Object(r.Penalties[ix])
		}
		enc.EndArray()
	}

	enc.TimestampField("created", r.Created)
}

// DecodeField implements the zjson decoder interface.De
func (r *RoundResult) DecodeField(dec *zjson.Decoder, key []byte) error {
	ok := true

	if string(key) == "spinSeq" {
		r.SpinSeq, ok = dec.Int()
	} else if string(key) == "balanceBefore" {
		r.BalanceBefore, ok = dec.Int64()
	} else if string(key) == "balanceAfter" {
		r.BalanceAfter, ok = dec.Int64()
	} else if string(key) == "bet" {
		r.Bet, ok = dec.Int64()
	} else if string(key) == "win" {
		r.Win, ok = dec.Int64()
	} else if string(key) == "totalWin" {
		r.TotalWin, ok = dec.Int64()
	} else if string(key) == "progressiveWin" {
		r.ProgressiveWin, ok = dec.Int64()
	} else if string(key) == "bonusWin" {
		r.BonusWin, ok = dec.Int64()
	} else if string(key) == "spinWin" {
		r.SpinWin, ok = dec.Int64()
	} else if string(key) == "awardedFreeGames" {
		r.AwardedFreeGames, ok = dec.Uint64()
	} else if string(key) == "freeGames" {
		r.FreeGames, ok = dec.Uint64()
	} else if string(key) == "totalPayout" {
		r.TotalPayout, ok = dec.Float()
	} else if string(key) == "maxPayout" {
		r.MaxPayout, ok = dec.Float()
	} else if string(key) == "spinData" {
		r.SpinData = slots.AcquireSpinResult(nil)
		ok = dec.Object(r.SpinData)
	} else if string(key) == "instantBonus" {
		r.InstantBonus = results.AcquireInstantBonusChoice("")
		ok = dec.Object(r.InstantBonus)
	} else if string(key) == "bonusSelector" {
		r.BonusSelector = results.AcquireBonusSelectorChoice(0, 0)
		ok = dec.Object(r.BonusSelector)
	} else if string(key) == "bonusWheel" {
		r.BonusWheel = wheel.AcquireBonusWheelResult(0, nil)
		ok = dec.Object(r.BonusWheel)
	} else if string(key) == "symbolsState" {
		r.SymbolsState = slots.AcquireSymbolsState(nil)
		ok = dec.Object(r.SymbolsState)
	} else if string(key) == "payouts" {
		ok = dec.Array(r.decodePayout)
	} else if string(key) == "penalties" {
		ok = dec.Array(r.decodePenalty)
	} else if string(key) == "created" {
		r.Created, ok = dec.Timestamp()
	} else {
		return fmt.Errorf("RoundResult.Decode: invalid field '%s'", string(key))
	}

	if ok {
		return nil
	}
	return dec.Error()
}

func (r *RoundResult) decodePayout(dec *zjson.Decoder) error {
	p := slots.WinlinePayoutFromData(0, 1, 0, 0, 0, 0, nil)
	if ok := dec.Object(p); ok {
		r.Payouts = append(r.Payouts, p)
		return nil
	}
	return dec.Error()

}

func (r *RoundResult) decodePenalty(dec *zjson.Decoder) error {
	p := &slots.SpinPenalty{}
	if ok := dec.Object(p); ok {
		r.Penalties = append(r.Penalties, p)
		return nil
	}
	return dec.Error()

}

// RoundResult contains the data for a single spin result with (re)calculated balances, progressive win amount.
// It also contains pointers to the actual spin result, instant bonus, bonus wheel, etc.
type RoundResult struct {
	SpinSeq          int                     `json:"spinSeq,omitempty"`
	BalanceBefore    int64                   `json:"balanceBefore,omitempty"`
	BalanceAfter     int64                   `json:"balanceAfter,omitempty"`
	Bet              int64                   `json:"bet,omitempty"`
	Win              int64                   `json:"win,omitempty"`
	TotalWin         int64                   `json:"totalWin,omitempty"`
	ProgressiveWin   int64                   `json:"progressiveWin,omitempty"`
	BonusWin         int64                   `json:"bonusWin,omitempty"`
	SpinWin          int64                   `json:"spinWin,omitempty"`
	AwardedFreeGames uint64                  `json:"awardedFreeGames,omitempty"`
	FreeGames        uint64                  `json:"freeGames,omitempty"`
	TotalPayout      float64                 `json:"totalPayout,omitempty"`
	MaxPayout        float64                 `json:"maxPayout,omitempty"`
	SpinData         *slots.SpinResult       `json:"spinData,omitempty"`
	InstantBonus     *results.InstantBonus   `json:"instantBonus,omitempty"`
	BonusSelector    *results.BonusSelector  `json:"bonusSelector,omitempty"`
	BonusWheel       *wheel.BonusWheelResult `json:"bonusWheel,omitempty"`
	SymbolsState     *slots.SymbolsState     `json:"symbolsState,omitempty"`
	Payouts          results.Payouts         `json:"payouts,omitempty"`
	Penalties        results.Penalties       `json:"penalties,omitempty"`
	Created          time.Time               `json:"created,omitempty"`
	pool.Object
}

// roundResultsPool is the memory pool for single spin results.
var roundResultPool = pool.NewProducer(func() (pool.Objecter, func()) {
	r := &RoundResult{
		Payouts: make(results.Payouts, 0, 8),
	}
	return r, r.reset
})

// Reset implements the Objecter.Reset interface.
func (r *RoundResult) reset() {
	if r != nil {
		if r.SpinData != nil {
			r.SpinData.Release()
			r.SpinData = nil
		}
		if r.InstantBonus != nil {
			r.InstantBonus.Release()
			r.InstantBonus = nil
		}
		if r.BonusSelector != nil {
			r.BonusSelector.Release()
			r.BonusSelector = nil
		}
		if r.BonusWheel != nil {
			r.BonusWheel.Release()
			r.BonusWheel = nil
		}
		if r.SymbolsState != nil {
			r.SymbolsState.Release()
			r.SymbolsState = nil
		}

		for ix := range r.Payouts {
			r.Payouts[ix].Release()
			r.Payouts[ix] = nil
		}

		for ix := range r.Penalties {
			r.Penalties[ix].Release()
			r.Penalties[ix] = nil
		}

		r.SpinSeq = 0
		r.BalanceBefore = 0
		r.BalanceAfter = 0
		r.Bet = 0
		r.Win = 0
		r.TotalWin = 0
		r.ProgressiveWin = 0
		r.BonusWin = 0
		r.SpinWin = 0
		r.AwardedFreeGames = 0
		r.FreeGames = 0
		r.TotalPayout = 0.0
		r.MaxPayout = 0.0
		r.Payouts = r.Payouts[:0]
		r.Penalties = r.Penalties[:0]
		r.Created = time.Time{}
	}
}

// RoundResults is a convenience type for a slice of single spin results.
type RoundResults []*RoundResult

// PurgeRoundResults returns input reset to zero length or a new slice if input is shorter than the requested capacity.
func PurgeRoundResults(input RoundResults, capacity int) RoundResults {
	if cap(input) < capacity {
		return make(RoundResults, 0, capacity)
	}
	return input[:0]
}

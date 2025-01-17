package slots

import (
	"fmt"
	"math"
	"reflect"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// NewSlotReduction instantiates a penalty for a "reverse win" spin with a fixed reduction.
func NewSlotReduction(symbol utils.Index, count uint8, reduction float64, spin *Spin) results.Penalty {
	p := initPenalty(results.SlotReduction, symbol, count, spin)
	p.reduction = reduction
	p.factor = -reduction
	return p
}

// NewSlotDivision instantiates a penalty for a "reverse win" spin which divides the current total payout.
func NewSlotDivision(symbol utils.Index, count uint8, division float64, spin *Spin) results.Penalty {
	p := initPenalty(results.SlotDivision, symbol, count, spin)
	p.division = division
	return p
}

func (p *SpinPenalty) AsPayout() results.Payout {
	payout := &SpinPayout{
		direction: PayScatter,
		count:     1,
		symbol:    p.symbol,
		factor:    p.factor,
		payMap:    nil,
		message:   p.message,
	}

	switch p.kind {
	case results.SlotReduction:
		payout.kind = results.SlotReducePenalty
	case results.SlotDivision:
		payout.kind = results.SlotDividePenalty
	}

	return payout
}

// Kind returns the penalty kind.
func (p *SpinPenalty) Kind() results.PenaltyKind {
	return p.kind
}

// Reduce returns the reduction for the penalty.
func (p *SpinPenalty) Reduce() float64 {
	return math.Round(p.reduction*10.0) / 10.0
}

// Divide returns the divisor for the penalty.
func (p *SpinPenalty) Divide() float64 {
	return math.Round(p.division*10.0) / 10.0
}

// SetFactor implements the Localizer.SetFactor interface.
func (p *SpinPenalty) SetFactor(factor float64) {
	p.factor = factor
}

// SetMessage implements the Localizer.SetMessage interface.
func (p *SpinPenalty) SetMessage(msg string) {
	p.message = msg
}

// EncodeFields implements the zjson.Encoder.EncodeFields interface.
func (p *SpinPenalty) EncodeFields(enc *zjson.Encoder) {
	enc.Uint8Field("kind", uint8(p.kind))
	enc.Uint8FieldOpt("count", p.count)
	enc.Uint16FieldOpt("symbol", uint16(p.symbol))
	enc.FloatField("reduction", math.Round(p.reduction*10)/10, 'g', -1)
	enc.FloatField("division", math.Round(p.division*10)/10, 'g', -1)
	enc.FloatField("factor", math.Round(p.factor*10)/10, 'g', -1)
	enc.StringFieldOpt("message", p.message)

	if len(p.payMap) > 0 {
		enc.StartArrayField("payMap")
		for ix := range p.payMap {
			enc.Uint64(uint64(p.payMap[ix]))
		}
		enc.EndArray()
	}
}

// DecodeField implements the zjson.Decoder.DecodeField interface.
func (p *SpinPenalty) DecodeField(dec *zjson.Decoder, key []byte) error {
	var ok, esc bool
	var i8 uint8
	var i16 uint16
	var b []byte

	if string(key) == "kind" {
		if i8, ok = dec.Uint8(); ok {
			p.kind = results.PenaltyKind(i8)
		}
	} else if string(key) == "symbol" {
		if i16, ok = dec.Uint16(); ok {
			p.symbol = utils.Index(i16)
		}
	} else if string(key) == "count" {
		p.count, ok = dec.Uint8()
	} else if string(key) == "reduction" {
		p.reduction, ok = dec.Float()
	} else if string(key) == "division" {
		p.division, ok = dec.Float()
	} else if string(key) == "factor" {
		p.factor, ok = dec.Float()
	} else if string(key) == "message" {
		if b, esc, ok = dec.String(); ok {
			if esc {
				p.message = string(dec.Unescaped(b))
			} else {
				p.message = string(b)
			}
		}
	} else if string(key) == "payMap" {
		dec.Array(func(dec *zjson.Decoder) error {
			if i8, ok = dec.Uint8(); ok {
				p.payMap = append(p.payMap, i8)
			}
			return dec.Error()
		})
	} else {
		return fmt.Errorf("SpinPenalty.DecodeField invalid field: %s", string(key))
	}

	if ok {
		return nil
	}
	return dec.Error()
}

// SpinPenalty contains the details of a spin reward.
// Keep fields ordered by ascending SizeOf().
type SpinPenalty struct {
	kind      results.PenaltyKind
	count     uint8
	symbol    utils.Index
	reduction float64
	division  float64
	factor    float64
	payMap    utils.UInt8s
	message   string
	pool.Object
}

// initPenalty initializes a new penalty from the memory pool.
func initPenalty(kind results.PenaltyKind, symbol utils.Index, count uint8, spin *Spin) *SpinPenalty {
	p := spinPenaltyProducer.Acquire().(*SpinPenalty)
	p.kind = kind
	p.count = count
	p.symbol = symbol
	p.payMap = spin.ScatterMap(symbol, count, p.payMap)
	return p
}

// spinPenaltyProducer is the memory pool for spin penalties.
// Make sure to initialize all slices appropriately!
var spinPenaltyProducer = pool.NewProducer(func() (pool.Objecter, func()) {
	p := &SpinPenalty{
		payMap: make(utils.UInt8s, 0, 8),
	}
	return p, p.reset
})

// reset clears the penalty.
func (p *SpinPenalty) reset() {
	if p != nil {
		p.kind = 0
		p.count = 0
		p.symbol = 0
		p.reduction = 0.0
		p.division = 0.0
		p.factor = 0.0
		p.message = ""

		clear(p.payMap)
		p.payMap = p.payMap[:0]
	}
}

func (p *SpinPenalty) Equals(other *SpinPenalty) bool {
	return p.kind == other.kind &&
		p.count == other.count &&
		p.symbol == other.symbol &&
		p.reduction == other.reduction &&
		p.division == other.division &&
		p.factor == other.factor &&
		p.message == other.message &&
		reflect.DeepEqual(p.payMap, other.payMap)
}

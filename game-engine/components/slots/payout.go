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

// WinlinePayoutFromData instantiates a new payout for a winline from an AllPaylines payline.
func WinlinePayoutFromData(factor, multiplier float64, symbol utils.Index, count uint8, direction PayDirection, paylineID uint8, rows utils.UInt8s) results.Payout {
	p := initPayout(results.SlotWinline, direction, count, symbol, factor, multiplier)
	p.paylineID = paylineID
	p.payRows = utils.CopyPurgeUInt8s(rows, p.payRows, 8)
	return p
}

// AllPaylinePayout instantiates a new payout for a winline from an AllPaylines payline.
func AllPaylinePayout(factor, multiplier float64, symbol utils.Index, count uint8, rows utils.UInt8s) results.Payout {
	p := initPayout(results.SlotWinline, PayLTR, count, symbol, factor, multiplier)
	p.payRows = utils.CopyPurgeUInt8s(rows, p.payRows, 8)
	return p
}

// ClusterPayout instantiates a payout for a cluster of symbols using the given payMap.
func ClusterPayout(factor, multiplier float64, symbol utils.Index, count uint8, payMap utils.UInt8s) results.Payout {
	p := initPayout(results.SlotCluster, PayCluster, count, symbol, factor, multiplier)
	p.payMap = utils.CopyPurgeUInt8s(payMap, p.payMap, 16)
	return p
}

// WildSymbolPayout instantiates a payout for a wild symbol with the given result.
func WildSymbolPayout(factor, multiplier float64, symbol utils.Index, count uint8, spin *Spin) results.Payout {
	p := initPayout(results.SlotWilds, PayScatter, count, symbol, factor, multiplier)
	p.payMap = spin.ScatterMap(symbol, count, p.payMap)
	return p
}

// WildSymbolPayoutWithMap instantiates a payout for a wild symbol with the given result and payMap.
func WildSymbolPayoutWithMap(factor, multiplier float64, symbol utils.Index, count uint8, payMap utils.UInt8s) results.Payout {
	p := initPayout(results.SlotWilds, PayScatter, count, symbol, factor, multiplier)
	p.payMap = utils.CopyPurgeUInt8s(payMap, p.payMap, 16)
	return p
}

// ScatterSymbolPayout instantiates a payout for a scatter symbol with the given result.
func ScatterSymbolPayout(factor, multiplier float64, symbol utils.Index, count uint8, spin *Spin) results.Payout {
	p := initPayout(results.SlotScatters, PayScatter, count, symbol, factor, multiplier)
	p.payMap = spin.ScatterMap(symbol, count, p.payMap)
	return p
}

// ScatterSymbolPayoutWithMap instantiates a payout for a scatter symbol with the given result and payMap.
func ScatterSymbolPayoutWithMap(factor, multiplier float64, symbol utils.Index, count uint8, payMap utils.UInt8s) results.Payout {
	p := initPayout(results.SlotScatters, PayScatter, count, symbol, factor, multiplier)
	p.payMap = utils.CopyPurgeUInt8s(payMap, p.payMap, 16)
	return p
}

// BombScatterPayout instantiates a payout for a scatter symbol after a bomb explosion.
func BombScatterPayout(factor, multiplier float64, symbol utils.Index, count uint8, spin *Spin) results.Payout {
	p := initPayout(results.SlotBombScatters, PayScatter, count, symbol, factor, multiplier)
	p.payMap = spin.ScatterMap(symbol, count, p.payMap)
	return p
}

// BonusSymbolPayout instantiates a payout for a matching bonus symbol during free spins.
func BonusSymbolPayout(factor, multiplier float64, symbol utils.Index, count uint8) results.Payout {
	p := initPayout(results.SlotBonusSymbol, PayScatter, count, symbol, factor, multiplier)
	return p
}

// SuperSymbolPayout instantiates a payout for a super shape win.
func SuperSymbolPayout(factor, multiplier float64, symbol utils.Index, count uint8, spin *Spin) results.Payout {
	p := initPayout(results.SlotSuperShape, PayScatter, count, symbol, factor, multiplier)
	p.payMap = spin.ScatterMap(symbol, count, p.payMap)
	return p
}

// Kind returns the reward kind.
func (p *SpinPayout) Kind() results.PayoutKind {
	return p.kind
}

// Count returns the number of symbols counted for the reward.
func (p *SpinPayout) Count() uint8 {
	return p.count
}

// Direction returns the direction for the reward.
func (p *SpinPayout) Direction() PayDirection {
	return p.direction
}

// PaylineID returns the payline id for the reward.
func (p *SpinPayout) PaylineID() uint8 {
	return p.paylineID
}

// AllPaylineID returns the payline id for the reward.
func (p *SpinPayout) AllPaylineID() int {
	var result int
	for _, row := range p.payRows[:p.count] {
		result *= 12
		result += int(row) + 1
	}
	return result
}

// Symbol returns the symbol for the reward.
func (p *SpinPayout) Symbol() utils.Index {
	return p.symbol
}

// Factor returns the win factor for the reward.
func (p *SpinPayout) Factor() float64 {
	return math.Round(p.factor*100.0) / 100.0
}

// Multiplier returns the multiplier for the reward.
func (p *SpinPayout) Multiplier() float64 {
	return math.Round(p.multiplier*100.0) / 100.0
}

// Total returns the total win factor (including the multiplier if valid).
func (p *SpinPayout) Total() float64 {
	payout := p.Factor()
	if utils.ValidMultiplier(p.multiplier) {
		return math.Round(payout*p.multiplier*100.0) / 100.0
	}
	return payout
}

// PayRows returns the row numbers of the symbol.
func (p *SpinPayout) PayRows() utils.UInt8s {
	return p.payRows
}

// PayMap returns the tile positions of the symbol.
func (p *SpinPayout) PayMap() utils.UInt8s {
	return p.payMap
}

// SetMessage implements the Localizer.SetMessage interface.
func (p *SpinPayout) SetMessage(msg string) {
	p.message = msg
}

// EncodeFields implements the zjson.Encoder.EncodeFields interface.
func (p *SpinPayout) EncodeFields(enc *zjson.Encoder) {
	enc.Uint8Field("kind", uint8(p.kind))
	enc.FloatField("payout", math.Round(p.factor*100)/100, 'g', -1)

	if utils.ValidMultiplier(p.multiplier) {
		enc.FloatField("multiplier", math.Round(p.multiplier*100)/100, 'g', -1)
	}

	enc.Uint16FieldOpt("symbol", uint16(p.symbol))
	enc.Uint8FieldOpt("count", p.count)
	enc.Uint8FieldOpt("direction", uint8(p.direction))
	enc.Uint8FieldOpt("paylineID", p.paylineID)
	enc.StringFieldOpt("message", p.message)

	if len(p.payRows) > 0 {
		enc.StartArrayField("paylineRows")
		for ix := range p.payRows {
			enc.Uint64(uint64(p.payRows[ix]))
		}
		enc.EndArray()
	}

	if len(p.payMap) > 0 {
		enc.StartArrayField("payMap")
		for ix := range p.payMap {
			enc.Uint64(uint64(p.payMap[ix]))
		}
		enc.EndArray()
	}
}

// DecodeField implements the zjson.Decoder.DecodeField interface.
func (p *SpinPayout) DecodeField(dec *zjson.Decoder, key []byte) error {
	var ok, esc bool
	var i8 uint8
	var i16 uint16
	var b []byte

	if string(key) == "kind" {
		if i8, ok = dec.Uint8(); ok {
			p.kind = results.PayoutKind(i8)
		}
	} else if string(key) == "payout" {
		p.factor, ok = dec.Float()
	} else if string(key) == "multiplier" {
		p.multiplier, ok = dec.Float()
	} else if string(key) == "symbol" {
		if i16, ok = dec.Uint16(); ok {
			p.symbol = utils.Index(i16)
		}
	} else if string(key) == "count" {
		p.count, ok = dec.Uint8()
	} else if string(key) == "direction" {
		if i8, ok = dec.Uint8(); ok {
			p.direction = PayDirection(i8)
		}
	} else if string(key) == "paylineID" {
		p.paylineID, ok = dec.Uint8()
	} else if string(key) == "message" {
		if b, esc, ok = dec.String(); ok {
			if esc {
				p.message = string(dec.Unescaped(b))
			} else {
				p.message = string(b)
			}
		}
	} else if string(key) == "paylineRows" {
		dec.Array(func(dec *zjson.Decoder) error {
			if i8, ok = dec.Uint8(); ok {
				p.payRows = append(p.payRows, i8)
			}
			return dec.Error()
		})
	} else if string(key) == "payMap" {
		dec.Array(func(dec *zjson.Decoder) error {
			if i8, ok = dec.Uint8(); ok {
				p.payMap = append(p.payMap, i8)
			}
			return dec.Error()
		})
	} else {
		return fmt.Errorf("SpinPayout.DecodeField invalid field: %s", string(key))
	}

	if ok {
		return nil
	}
	return dec.Error()
}

// SpinPayout contains the details of a spin reward.
// Keep fields ordered by ascending SizeOf().
type SpinPayout struct {
	kind       results.PayoutKind
	direction  PayDirection
	count      uint8
	paylineID  uint8
	symbol     utils.Index
	factor     float64
	multiplier float64
	payRows    utils.UInt8s
	payMap     utils.UInt8s
	message    string
	pool.Object
}

// initPayout initializes a new payout from the memory pool.
func initPayout(kind results.PayoutKind, direction PayDirection, count uint8, symbol utils.Index, factor, multiplier float64) *SpinPayout {
	p := spinPayoutProducer.Acquire().(*SpinPayout)
	p.kind = kind
	p.direction = direction
	p.count = count
	p.symbol = symbol
	p.factor = math.Round(factor*100.0) / 100.0
	p.multiplier = math.Round(multiplier*100.0) / 100.0
	return p
}

// spinPayoutProducer is the memory pool for spin payouts.
// Make sure to initialize all slices appropriately!
var spinPayoutProducer = pool.NewProducer(func() (pool.Objecter, func()) {
	p := &SpinPayout{
		payRows: make(utils.UInt8s, 0, 8),
		payMap:  make(utils.UInt8s, 0, 16),
	}
	return p, p.reset
})

// reset clears the payout.
func (p *SpinPayout) reset() {
	if p != nil {
		p.kind = 0
		p.count = 0
		p.direction = PayLTR
		p.paylineID = 0
		p.symbol = utils.NullIndex
		p.factor = 0.0
		p.multiplier = 0.0
		p.message = ""

		clear(p.payRows)
		p.payRows = p.payRows[:0]
		clear(p.payMap)
		p.payMap = p.payMap[:0]
	}
}

func (p *SpinPayout) Equals(other *SpinPayout) bool {
	return p.kind == other.kind &&
		p.count == other.count &&
		p.direction == other.direction &&
		p.paylineID == other.paylineID &&
		p.symbol == other.symbol &&
		p.factor == other.factor &&
		p.multiplier == other.multiplier &&
		p.message == other.message &&
		reflect.DeepEqual(p.payRows, other.payRows) &&
		reflect.DeepEqual(p.payMap, other.payMap)
}

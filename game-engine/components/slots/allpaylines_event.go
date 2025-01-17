package slots

import (
	"math"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/object"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

const allPaylinesEvent = results.AllPaylinesEvent

// AcquireAllPaylinesEvents instantiates one or more new all paylines animation events.
func AcquireAllPaylinesEvents(symbol utils.Index, reels, rows int, payouts results.Payouts) []results.Animator {
	events := make([]results.Animator, 0, 8)

	for count := 2; count <= reels; count++ {
		event := allPaylinesEventProducer.Acquire().(*AllPaylinesEvent)
		event.symbol = symbol
		event.count = uint8(count)

		const maxGrid uint8 = 200
		var maxOffset uint8

		r := uint8(rows)
		m := make(map[uint8]bool, maxGrid)

		for ix := range payouts {
			if p, ok := payouts[ix].(*SpinPayout); ok && p.kind == results.SlotWinline && p.paylineID == 0 && p.symbol == symbol && p.count == uint8(count) && len(p.payRows) > 0 {
				event.factor += p.Total()
				if utils.ValidMultiplier(p.multiplier) {
					event.paylines += uint16(math.Round(p.multiplier))
				} else {
					event.paylines++
				}
				for iy := uint8(0); iy < p.count; iy++ {
					offset := iy*r + p.payRows[iy]
					m[offset] = true
					if offset > maxOffset {
						maxOffset = offset
					}
				}
			}
		}

		if event.paylines == 0 {
			event.Release()
			continue
		}

		max := object.NormalizeSize(int(maxOffset), 8)
		event.payMap = utils.PurgeUInt8s(event.payMap, max)

		for ix := uint8(0); ix <= maxOffset; ix++ {
			if m[ix] {
				event.payMap = append(event.payMap, ix)
			}
		}

		events = append(events, event)
	}

	return events
}

// Kind implements the Animator.Kind interface.
func (e *AllPaylinesEvent) Kind() results.EventKind {
	return allPaylinesEvent
}

// SetMessage implements the Localizer.SetMessage interface.
func (e *AllPaylinesEvent) SetMessage(msg string) {
	e.message = msg
}

// EncodeFields implements the zjson.ObjectEncoder.EncodeFields interface.
func (e *AllPaylinesEvent) EncodeFields(enc *zjson.Encoder) {
	enc.Uint8Field("kind", uint8(allPaylinesEvent))
	enc.Uint16Field("symbol", uint16(e.symbol))
	enc.Uint8Field("count", e.count)
	enc.Uint16Field("paylines", e.paylines)
	enc.FloatField("factor", e.factor, 'g', -1)
	enc.StringFieldOpt("message", e.message)

	if len(e.payMap) > 0 {
		enc.StartArrayField("payMap")
		for ix := range e.payMap {
			enc.Uint64(uint64(e.payMap[ix]))
		}
		enc.EndArray()
	}
}

// AllPaylinesEvent is the animation event for a single symbol all paylines animation event.
type AllPaylinesEvent struct {
	count    uint8
	symbol   utils.Index
	paylines uint16
	factor   float64
	payMap   utils.UInt8s
	message  string
	pool.Object
}

// Factor returns the factor for the all paylines animation event.
func (e *AllPaylinesEvent) Factor() float64 {
	return e.factor
}

// Count returns the symbol count for the all paylines animation event.
func (e *AllPaylinesEvent) Count() uint8 {
	return e.count
}

// Paylines returns the # paylines for the all paylines animation event.
func (e *AllPaylinesEvent) Paylines() uint16 {
	return e.paylines
}

var allPaylinesEventProducer = pool.NewProducer(func() (pool.Objecter, func()) {
	e := &AllPaylinesEvent{
		payMap: make(utils.UInt8s, 0, 8),
	}
	return e, e.reset
})

// reset clear the all paylines animation event.
func (e *AllPaylinesEvent) reset() {
	if e != nil {
		e.count = 0
		e.symbol = 0
		e.paylines = 0
		e.factor = 0.0
		e.payMap = e.payMap[:0]
		e.message = ""
	}
}

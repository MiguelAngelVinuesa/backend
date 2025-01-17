package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

const payoutEvent = results.PayoutEvent

// AcquirePayoutEvent instantiates a new spin payout animation event.
func AcquirePayoutEvent(p *SpinPayout) results.Animator {
	e := payoutEventProducer.Acquire().(*PayoutEvent)

	e.payoutKind = p.kind
	e.count = p.count
	e.symbol = p.symbol
	e.factor = p.Total()

	e.payRows = utils.CopyPurgeUInt8s(p.payRows, e.payRows, 8)
	if l := int(p.count); l > 0 && p.direction == PayLTR && len(e.payRows) > l {
		e.payRows = e.payRows[:l]
	}

	e.payMap = utils.CopyPurgeUInt8s(p.payMap, e.payMap, 16)

	return e
}

// Kind implements the Animator.Kind interface.
func (e *PayoutEvent) Kind() results.EventKind {
	return payoutEvent
}

// SetMessage implements the Localizer.SetMessage interface.
func (e *PayoutEvent) SetMessage(msg string) {
	e.message = msg
}

// SetPaylines adds a set of paylines to the event.
func (e *PayoutEvent) SetPaylines(paylines *PaylineSet) {
	if e.paylines == nil {
		e.paylines = make([]utils.UInt8s, 0, len(paylines.paylines))
	}

	for ix := range paylines.paylines {
		e.paylines = append(e.paylines, paylines.paylines[ix].rows)
	}
}

// EncodeFields implements the zjson.ObjectEncoder.EncodeFields interface.
func (e *PayoutEvent) EncodeFields(enc *zjson.Encoder) {
	enc.Uint8Field("kind", uint8(payoutEvent))
	enc.Uint8Field("payoutKind", uint8(e.payoutKind))
	enc.Uint8Field("count", e.count)
	enc.Uint16Field("symbol", uint16(e.symbol))
	enc.FloatField("factor", e.factor, 'g', -1)
	enc.StringFieldOpt("message", e.message)

	if len(e.payRows) > 0 {
		enc.StartArrayField("paylineRows")
		for ix := range e.payRows {
			enc.Uint64(uint64(e.payRows[ix]))
		}
		enc.EndArray()
	}

	if len(e.payMap) > 0 {
		enc.StartArrayField("payMap")
		for ix := range e.payMap {
			enc.Uint64(uint64(e.payMap[ix]))
		}
		enc.EndArray()
	}

	if len(e.paylines) > 0 {
		enc.StartArrayField("paylines")
		for ix := range e.paylines {
			line := e.paylines[ix]

			enc.StartArray()
			for iy := range line {
				enc.Uint64(uint64(line[iy]))
			}
			enc.EndArray()
		}
		enc.EndArray()
	}
}

// PayoutEvent is the animation event for a spin payout.
type PayoutEvent struct {
	payoutKind results.PayoutKind
	count      uint8
	symbol     utils.Index
	factor     float64
	payRows    utils.UInt8s
	payMap     utils.UInt8s
	paylines   []utils.UInt8s
	message    string
	pool.Object
}

var payoutEventProducer = pool.NewProducer(func() (pool.Objecter, func()) {
	e := &PayoutEvent{
		payRows: make(utils.UInt8s, 0, 8),
		payMap:  make(utils.UInt8s, 0, 16),
	}
	return e, e.reset
})

func (e *PayoutEvent) Factor() float64 {
	return e.factor
}

func (e *PayoutEvent) Count() uint8 {
	return e.count
}

// reset clears the playout event.
func (e *PayoutEvent) reset() {
	if e != nil {
		e.payoutKind = 0
		e.count = 0
		e.symbol = 0
		e.factor = 0.0
		e.payRows = e.payRows[:0]
		e.payMap = e.payMap[:0]
		e.message = ""

		for ix := range e.paylines {
			e.paylines[ix] = nil
		}
		e.paylines = e.paylines[:0]
	}
}

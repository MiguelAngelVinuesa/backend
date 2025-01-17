package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/object"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

const refillEvent = results.RefillEvent

// AcquireRefillEvent instantiates a refill animation event.
func AcquireRefillEvent(spin *SpinResult) results.Animator {
	e := refillEventProducer.Acquire().(*RefillEvent)

	max := object.NormalizeSize(len(spin.initial), 16)
	e.refill = utils.PurgeUInt8s(e.refill, max)
	if l := len(spin.afterClear); l == len(spin.initial) {
		for ix := range spin.afterClear {
			if spin.afterClear[ix] == 0 {
				e.refill = append(e.refill, uint8(ix))
			}
		}
	}

	return e
}

// Kind implements the Animator.Kind interface.
func (e *RefillEvent) Kind() results.EventKind {
	return refillEvent
}

// EncodeFields implements the zjson.ObjectEncoder.EncodeFields interface.
func (e *RefillEvent) EncodeFields(enc *zjson.Encoder) {
	enc.Uint8Field("kind", uint8(refillEvent))

	if len(e.refill) > 0 {
		enc.StartArrayField("refill")
		for ix := range e.refill {
			enc.Uint64(uint64(e.refill[ix]))
		}
		enc.EndArray()
	}
}

// RefillEvent is the animation event for a refill spin.
type RefillEvent struct {
	refill utils.UInt8s
	pool.Object
}

var refillEventProducer = pool.NewProducer(func() (pool.Objecter, func()) {
	e := &RefillEvent{
		refill: make(utils.UInt8s, 0, 16),
	}
	return e, e.reset
})

// reset clears the refill event.
func (e *RefillEvent) reset() {
	if e != nil {
		e.refill = e.refill[:0]
	}
}

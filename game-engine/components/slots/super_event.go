package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/object"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

const superEvent = results.SuperEvent

// AcquireSuperEvent instantiates a super shape animation event.
func AcquireSuperEvent(spin *SpinResult) results.Animator {
	e := superEventProducer.Acquire().(*SuperEvent)

	e.first = true

	max := object.NormalizeSize(len(spin.initial), 16)
	e.shape = utils.PurgeUInt8s(e.shape, max)
	if len(spin.sticky) == len(spin.initial) {
		for ix := range spin.sticky {
			switch spin.sticky[ix] {
			case 1:
				e.first = false
			case 2:
				e.shape = append(e.shape, uint8(ix))
			}
		}
	}

	return e
}

// Kind implements the Animator.Kind interface.
func (e *SuperEvent) Kind() results.EventKind {
	return superEvent
}

// EncodeFields implements the zjson.ObjectEncoder.EncodeFields interface.
func (e *SuperEvent) EncodeFields(enc *zjson.Encoder) {
	enc.Uint8Field("kind", uint8(superEvent))

	if len(e.shape) > 0 {
		enc.StartArrayField("shape")
		for ix := range e.shape {
			enc.Uint64(uint64(e.shape[ix]))
		}
		enc.EndArray()
	}

	enc.IntBoolFieldOpt("first", e.first)
}

// SuperEvent is the animation event for a super shape spin.
type SuperEvent struct {
	first bool
	shape utils.UInt8s
	pool.Object
}

var superEventProducer = pool.NewProducer(func() (pool.Objecter, func()) {
	e := &SuperEvent{
		shape: make(utils.UInt8s, 0, 16),
	}
	return e, e.reset
})

// reset clears the super event.
func (e *SuperEvent) reset() {
	if e != nil {
		e.first = false
		e.shape = e.shape[:0]
	}
}

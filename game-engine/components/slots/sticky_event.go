package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/object"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

const stickyEvent = results.StickyEvent

// AcquireStickyEvent instantiates a refill animation event.
func AcquireStickyEvent(spin *SpinResult) results.Animator {
	e := stickyEventProducer.Acquire().(*StickyEvent)

	max := object.NormalizeSize(len(spin.initial), 16)
	e.stickies = utils.PurgeUInt8s(e.stickies, max)
	if len(spin.sticky) == len(spin.initial) {
		for ix := range spin.sticky {
			if spin.sticky[ix] != 0 {
				e.stickies = append(e.stickies, uint8(ix))
			}
		}
	}

	return e
}

// Kind implements the Animator.Kind interface.
func (e *StickyEvent) Kind() results.EventKind {
	return stickyEvent
}

// EncodeFields implements the zjson.ObjectEncoder.EncodeFields interface.
func (e *StickyEvent) EncodeFields(enc *zjson.Encoder) {
	enc.Uint8Field("kind", uint8(stickyEvent))

	if len(e.stickies) > 0 {
		enc.StartArrayField("stickies")
		for ix := range e.stickies {
			enc.Uint64(uint64(e.stickies[ix]))
		}
		enc.EndArray()
	}
}

// StickyEvent is the animation event for a refill spin.
type StickyEvent struct {
	stickies utils.UInt8s
	pool.Object
}

var stickyEventProducer = pool.NewProducer(func() (pool.Objecter, func()) {
	e := &StickyEvent{
		stickies: make(utils.UInt8s, 0, 16),
	}
	return e, e.reset
})

// reset clears the sticky event.
func (e *StickyEvent) reset() {
	if e != nil {
		e.stickies = e.stickies[:0]
	}
}

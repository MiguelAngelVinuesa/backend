package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/object"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

const clearEvent = results.ClearEvent

// AcquireClearEvent instantiates a clear tiles animation event.
func AcquireClearEvent(spin *SpinResult) results.Animator {
	e := clearEventProducer.Acquire().(*ClearEvent)

	max := object.NormalizeSize(len(spin.initial), 16)
	e.clear = utils.PurgeUInt8s(e.clear, max)

	if l := len(spin.afterClear); l == len(spin.initial) {
		for ix := range spin.afterClear {
			if spin.afterClear[ix] == 0 {
				e.clear = append(e.clear, uint8(ix))
			}
		}
	}

	return e
}

// Kind implements the Animator.Kind interface.
func (e *ClearEvent) Kind() results.EventKind {
	return clearEvent
}

// EncodeFields implements the zjson.ObjectEncoder.EncodeFields interface.
func (e *ClearEvent) EncodeFields(enc *zjson.Encoder) {
	enc.Uint8Field("kind", uint8(clearEvent))

	if len(e.clear) > 0 {
		enc.StartArrayField("clear")
		for ix := range e.clear {
			enc.Uint64(uint64(e.clear[ix]))
		}
		enc.EndArray()
	}
}

// ClearEvent is the animation event for clearing tiles.
type ClearEvent struct {
	clear utils.UInt8s
	pool.Object
}

var clearEventProducer = pool.NewProducer(func() (pool.Objecter, func()) {
	e := &ClearEvent{
		clear: make(utils.UInt8s, 0, 16),
	}
	return e, e.reset
})

// reset re-implements the objecter.ResetData interface.
func (e *ClearEvent) reset() {
	if e != nil {
		e.clear = e.clear[:0]
	}
}

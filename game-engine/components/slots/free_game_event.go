package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
)

const freeGameEvent = results.FreeGameEvent

// AcquireFreeGameEvent instantiates a free game(s) award animation event.
func AcquireFreeGameEvent(count uint64) results.Animator {
	e := freeGameEventProducer.Acquire().(*FreeGameEvent)
	e.count = count
	return e
}

// Kind implements the Animator.Kind interface.
func (e *FreeGameEvent) Kind() results.EventKind {
	return freeGameEvent
}

// EncodeFields implements the zjson.ObjectEncoder.EncodeFields interface.
func (e *FreeGameEvent) EncodeFields(enc *zjson.Encoder) {
	enc.Uint8Field("kind", uint8(freeGameEvent))
	enc.Uint64Field("count", e.count)
}

// FreeGameEvent is the animation event for awarded free game(s).
type FreeGameEvent struct {
	count uint64
	pool.Object
}

var freeGameEventProducer = pool.NewProducer(func() (pool.Objecter, func()) {
	e := &FreeGameEvent{}
	return e, e.reset
})

// reset clears the free game event.
func (e *FreeGameEvent) reset() {
	if e != nil {
		e.count = 0
	}
}

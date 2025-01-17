package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
)

// AcquireChoiceRequest instantiates a new player choice request animation event.
func AcquireChoiceRequest() results.Animator {
	e := choiceRequestProducer.Acquire().(*ChoiceRequest)
	return e
}

// Kind implements the Animator.Kind interface.
func (e *ChoiceRequest) Kind() results.EventKind {
	return results.ChoiceRequestEvent
}

// EncodeFields implements the zjson.ObjectEncoder.EncodeFields interface.
func (e *ChoiceRequest) EncodeFields(enc *zjson.Encoder) {
	enc.Uint8Field("kind", uint8(results.ChoiceRequestEvent))
}

// ChoiceRequest is the animation event to indicate the player has to make a choice.
type ChoiceRequest struct {
	pool.Object
}

var choiceRequestProducer = pool.NewProducer(func() (pool.Objecter, func()) {
	e := &ChoiceRequest{}
	return e, nil
})

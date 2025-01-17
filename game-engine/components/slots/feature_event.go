package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
)

// FeatureTransitionKind indicates a transition on the current "spin" or the next "spin".
type FeatureTransitionKind uint8

const (
	InstantBonusTeaser FeatureTransitionKind = iota + 1
	InstantBonusRequested
	InstantBonusResulted
	BonusWheelTransition
	ChestFeatureTransition
)

// AcquireFeatureTransition instantiates a new bonus feature transition animation event.
func AcquireFeatureTransition(transition FeatureTransitionKind) results.Animator {
	e := FeatureTransitionProducer.Acquire().(*FeatureTransition)
	e.transition = transition
	return e
}

// Kind implements the Animator.Kind interface.
func (e *FeatureTransition) Kind() results.EventKind {
	return results.FeatureTransitionEvent
}

// EncodeFields implements the zjson EncodeFields interface.
func (e *FeatureTransition) EncodeFields(enc *zjson.Encoder) {
	enc.Uint8Field("kind", uint8(results.FeatureTransitionEvent))
	enc.Uint8Field("transition", uint8(e.transition))
}

// FeatureTransition is the animation event to indicate a bonus feature transition happens.
type FeatureTransition struct {
	transition FeatureTransitionKind
	pool.Object
}

var FeatureTransitionProducer = pool.NewProducer(func() (pool.Objecter, func()) {
	e := &FeatureTransition{}
	return e, e.reset
})

func (e *FeatureTransition) reset() {
	if e != nil {
		e.transition = 0
	}
}

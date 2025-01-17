package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
)

const awardEvent = results.AwardEvent

// AcquireBonusGameEvent instantiates an award animation event.
func AcquireBonusGameEvent() results.Animator {
	e := awardEventProducer.Acquire().(*AwardEvent)
	e.bonusGame = true
	return e
}

// Kind implements the Animator.Kind interface.
func (e *AwardEvent) Kind() results.EventKind {
	return awardEvent
}

// EncodeFields implements the zjson.ObjectEncoder.EncodeFields interface.
func (e *AwardEvent) EncodeFields(enc *zjson.Encoder) {
	enc.Uint8Field("kind", uint8(awardEvent))
	enc.IntBoolFieldOpt("bonusGame", e.bonusGame)
}

// AwardEvent is the animation event for spin awards.
type AwardEvent struct {
	bonusGame bool
	pool.Object
}

var awardEventProducer = pool.NewProducer(func() (pool.Objecter, func()) {
	e := &AwardEvent{}
	return e, e.reset
})

// reset clears the award event.
func (e *AwardEvent) reset() {
	if e != nil {
		e.bonusGame = false
	}
}

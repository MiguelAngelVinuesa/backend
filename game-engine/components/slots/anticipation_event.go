package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
)

// AcquireReelAnticipation instantiates a new reel anticipation animation event.
func AcquireReelAnticipation(startReel, stopReel uint8) results.Animator {
	e := reelAnticipationProducer.Acquire().(*ReelAnticipation)
	e.startReel = startReel
	e.stopReel = stopReel
	return e
}

// Kind implements the Animator.Kind interface.
func (e *ReelAnticipation) Kind() results.EventKind {
	return results.ReelAnticipationEvent
}

// EncodeFields implements the zjson.ObjectEncoder.EncodeFields interface.
func (e *ReelAnticipation) EncodeFields(enc *zjson.Encoder) {
	enc.Uint8Field("kind", uint8(results.ReelAnticipationEvent))
	enc.Uint8FieldOpt("startReel", e.startReel)
	enc.Uint8FieldOpt("stopReel", e.stopReel)
}

// ReelAnticipation is the animation event for extending reel spins in anticipation of a special bonus.
// E.g. in BOT to anticipate dropping 3 or more scatters; in MGD to anticapte wild respin and/or 3 scatters.
type ReelAnticipation struct {
	startReel uint8 // the reel where the event should start (1-based).
	stopReel  uint8 // the reel where the event should stop (optional; 1-based).
	pool.Object
}

var reelAnticipationProducer = pool.NewProducer(func() (pool.Objecter, func()) {
	e := &ReelAnticipation{}
	return e, e.reset
})

// reset clears the anticipation event.
func (e *ReelAnticipation) reset() {
	if e != nil {
		e.startReel = 0
		e.stopReel = 0
	}
}

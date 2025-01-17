package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

type ReelEventKind uint8

const (
	HotReelEvent ReelEventKind = iota + 1
	ExpandReelEvent
	BlockedReelEvent
)

const reelEvent = results.ReelEvent

// AcquireHotReelEvent instantiates a hot reel animation.
func AcquireHotReelEvent(reel uint8, first bool) results.Animator {
	e := reelEventProducer.Acquire().(*ReelEvent)
	e.kind = HotReelEvent
	e.reel = reel
	e.first = first
	return e
}

// AcquireExpandReelEvent instantiates a hot reel animation.
func AcquireExpandReelEvent(reel uint8, rows utils.UInt8s, symbol utils.Index) results.Animator {
	e := reelEventProducer.Acquire().(*ReelEvent)
	e.kind = ExpandReelEvent
	e.reel = reel
	e.symbol = symbol
	e.rows = utils.CopyPurgeUInt8s(rows, e.rows, 2)
	return e
}

// AcquireBlockedReelEvent instantiates a hot reel animation.
func AcquireBlockedReelEvent(reel uint8, first bool) results.Animator {
	e := reelEventProducer.Acquire().(*ReelEvent)
	e.kind = BlockedReelEvent
	e.reel = reel
	e.first = first
	return e
}

// Kind implements the Animator.Kind interface.
func (e *ReelEvent) Kind() results.EventKind {
	return reelEvent
}

// Encode implements the zjson.ObjectEncoder.Encode interface.
func (e *ReelEvent) Encode(enc *zjson.Encoder) {
	enc.Object(e)
}

// EncodeFields implements the zjson.ObjectEncoder.EncodeFields interface.
func (e *ReelEvent) EncodeFields(enc *zjson.Encoder) {
	enc.Uint8Field("kind", uint8(reelEvent))
	enc.Uint8Field("reelKind", uint8(e.kind))
	enc.Uint8Field("reel", e.reel)
	enc.Uint16FieldOpt("symbol", uint16(e.symbol))

	if len(e.rows) > 0 {
		enc.StartArrayField("rows")
		for _, row := range e.rows {
			enc.Uint64(uint64(row))
		}
		enc.EndArray()
	}

	enc.IntBoolFieldOpt("firstTime", e.first)
}

// ReelEvent is the animation event for reel effects.
type ReelEvent struct {
	first  bool
	kind   ReelEventKind
	reel   uint8
	symbol utils.Index
	rows   utils.UInt8s
	pool.Object
}

var reelEventProducer = pool.NewProducer(func() (pool.Objecter, func()) {
	e := &ReelEvent{
		rows: make(utils.UInt8s, 0, 8),
	}
	return e, e.reset
})

// reset clears the reel event.
func (e *ReelEvent) reset() {
	if e != nil {
		e.first = false
		e.kind = 0
		e.reel = 0
		e.symbol = 0
		e.rows = e.rows[:0]
	}
}

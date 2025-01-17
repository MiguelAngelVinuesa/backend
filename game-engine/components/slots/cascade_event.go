package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/object"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

const cascadeEvent = results.CascadeEvent

// AcquireCascadeEvent instantiates a refill animation event.
func AcquireCascadeEvent(spin *SpinResult, reels, rows int) results.Animator {
	e := cascadeEventProducer.Acquire().(*CascadeEvent)

	max := object.NormalizeSize(len(spin.initial), 16)
	e.from = utils.PurgeUInt8s(e.from, max)
	e.to = utils.PurgeUInt8s(e.to, max)

	temp := make(utils.Indexes, 0, 100)[:reels*rows]
	copy(temp, spin.afterClear)

	for reel := 0; reel < reels; reel++ {
		top := reel * rows
		to := top + rows - 1
		for to >= top && temp[to] != utils.NullIndex {
			to--
		}

		from := to - 1
		for to > top && from >= top {
			for from >= top && temp[from] == utils.NullIndex {
				from--
			}

			if from >= top {
				e.from = append(e.from, uint8(from))
				e.to = append(e.to, uint8(to))
				temp[to] = temp[from]
				temp[from] = utils.NullIndex
			}

			to--
		}
	}

	return e
}

// Kind implements the Animator.Kind interface.
func (e *CascadeEvent) Kind() results.EventKind {
	return cascadeEvent
}

// EncodeFields implements the zjson.ObjectEncoder.EncodeFields interface.
func (e *CascadeEvent) EncodeFields(enc *zjson.Encoder) {
	enc.Uint8Field("kind", uint8(cascadeEvent))

	if len(e.from) > 0 {
		enc.StartArrayField("from")
		for ix := range e.from {
			enc.Uint64(uint64(e.from[ix]))
		}
		enc.EndArray()
	}
	if len(e.to) > 0 {
		enc.StartArrayField("to")
		for ix := range e.to {
			enc.Uint64(uint64(e.to[ix]))
		}
		enc.EndArray()
	}
}

// CascadeEvent is the animation event for a refill spin.
type CascadeEvent struct {
	from utils.UInt8s
	to   utils.UInt8s
	pool.Object
}

var cascadeEventProducer = pool.NewProducer(func() (pool.Objecter, func()) {
	e := &CascadeEvent{
		from: make(utils.UInt8s, 0, 16),
		to:   make(utils.UInt8s, 0, 16),
	}
	return e, e.reset
})

// reset clears the cascade event.
func (e *CascadeEvent) reset() {
	if e != nil {
		e.from = e.from[:0]
		e.to = e.to[:0]
	}
}

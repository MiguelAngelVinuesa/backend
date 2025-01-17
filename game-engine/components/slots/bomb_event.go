package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/object"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

const bombEvent = results.BombEvent

// AcquireBombEvent instantiates a bomb animation event.
func AcquireBombEvent(spin *SpinResult, reels, rows int, center Offsets, grid GridOffsets) results.Animator {
	e := bombEventProducer.Acquire().(*BombEvent)

	max := object.NormalizeSize(len(spin.initial), 16)
	haveExpand := len(spin.afterExpand) == len(spin.initial)

	e.center = uint8(center[0]*rows + center[1])
	e.grid = utils.PurgeUInt8s(e.grid, max)
	e.changed = utils.PurgeUInt8s(e.changed, max)

	for ix := range grid {
		p := grid[ix]
		x, y := center[0]+p[0], center[1]+p[1]
		if x >= 0 && x < reels && y >= 0 && y < rows {
			offset := uint8(x*rows + y)
			e.grid = append(e.grid, offset)

			if haveExpand && spin.initial[offset] != spin.afterExpand[offset] {
				e.changed = append(e.changed, offset)
			}
		}
	}

	return e
}

// Kind implements the Animator.Kind interface.
func (e *BombEvent) Kind() results.EventKind {
	return bombEvent
}

// EncodeFields implements the zjson.ObjectEncoder.EncodeFields interface.
func (e *BombEvent) EncodeFields(enc *zjson.Encoder) {
	enc.Uint8Field("kind", uint8(bombEvent))
	enc.Uint8Field("center", e.center)

	if len(e.grid) > 0 {
		enc.StartArrayField("grid")
		for ix := range e.grid {
			enc.Uint64(uint64(e.grid[ix]))
		}
		enc.EndArray()
	}

	if len(e.changed) > 0 {
		enc.StartArrayField("changed")
		for ix := range e.changed {
			enc.Uint64(uint64(e.changed[ix]))
		}
		enc.EndArray()
	}
}

// BombEvent is the animation event for a refill spin.
type BombEvent struct {
	center  uint8
	grid    utils.UInt8s
	changed utils.UInt8s
	pool.Object
}

var bombEventProducer = pool.NewProducer(func() (pool.Objecter, func()) {
	e := &BombEvent{
		grid:    make(utils.UInt8s, 0, 16),
		changed: make(utils.UInt8s, 0, 16),
	}
	return e, e.reset
})

// reset clears the bomb event.
func (e *BombEvent) reset() {
	if e != nil {
		e.center = 0
		e.grid = e.grid[:0]
		e.changed = e.changed[:0]
	}
}

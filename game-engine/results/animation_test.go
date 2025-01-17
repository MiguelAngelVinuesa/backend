package results

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventKind_String(t *testing.T) {
	assert.Equal(t, "reelanticipation", ReelAnticipationEvent.String())
	assert.Equal(t, "bonusanticipation", BonusAnticipationEvent.String())
	assert.Equal(t, "refill", RefillEvent.String())
	assert.Equal(t, "freegame", FreeGameEvent.String())
	assert.Equal(t, "payout", PayoutEvent.String())
	assert.Equal(t, "reel", ReelEvent.String())
	assert.Equal(t, "award", AwardEvent.String())
	assert.Equal(t, "bomb", BombEvent.String())
	assert.Equal(t, "super", SuperEvent.String())
	assert.Equal(t, "shooter", ShooterEvent.String())
	assert.Equal(t, "sticky", StickyEvent.String())
	assert.Equal(t, "clear", ClearEvent.String())
	assert.Equal(t, "cascade", CascadeEvent.String())
	assert.Equal(t, "megaways", AllPaylinesEvent.String())
}

func TestReleaseAnimations(t *testing.T) {
	t.Run("release animations", func(t *testing.T) {
		d1 := dummyAnimationProducer.Acquire().(*dummyAnimation)
		d2 := dummyAnimationProducer.Acquire().(*dummyAnimation)
		d3 := dummyAnimationProducer.Acquire().(*dummyAnimation)
		d4 := dummyAnimationProducer.Acquire().(*dummyAnimation)
		d5 := dummyAnimationProducer.Acquire().(*dummyAnimation)
		list := Animations{d1, d2, d3, d4, d5}

		got := ReleaseAnimations(list)
		require.NotNil(t, got)
		assert.Zero(t, len(got))

		got = got[:cap(got)]
		for ix := range got {
			assert.Nil(t, got[ix])
		}
	})
}

type dummyAnimation struct {
	pool.Object
}

func (e *dummyAnimation) Kind() EventKind                 { return EventKind(99) }
func (e *dummyAnimation) reset()                          {}
func (e *dummyAnimation) EncodeFields(enc *zjson.Encoder) { enc.Uint8Field("kind", 99) }

var dummyAnimationProducer = pool.NewProducer(func() (pool.Objecter, func()) { d := &dummyAnimation{}; return d, d.reset })

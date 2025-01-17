package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestAcquireHotReelEvent(t *testing.T) {
	testCases := []struct {
		name  string
		reel  uint8
		first bool
		want  string
	}{
		{"first 1", 1, true, `{"kind":5,"reelKind":1,"reel":1,"firstTime":1}`},
		{"first 3", 3, true, `{"kind":5,"reelKind":1,"reel":3,"firstTime":1}`},
		{"next 2", 2, false, `{"kind":5,"reelKind":1,"reel":2}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := AcquireHotReelEvent(tc.reel, tc.first).(*ReelEvent)
			require.NotNil(t, a)
			assert.Equal(t, HotReelEvent, a.kind)
			assert.Equal(t, tc.reel, a.reel)
			assert.Equal(t, tc.first, a.first)
			assert.Zero(t, a.symbol)
			assert.Empty(t, a.rows)

			enc := zjson.AcquireEncoder(1024)
			defer enc.Release()

			a.Encode(enc)
			got := enc.Bytes()
			assert.Equal(t, tc.want, string(got))
		})
	}
}

func TestAcquireExpandReelEvent(t *testing.T) {
	testCases := []struct {
		name   string
		reel   uint8
		symbol utils.Index
		rows   utils.UInt8s
		want   string
	}{
		{"first 1 3 0,2", 1, 3, utils.UInt8s{0, 2}, `{"kind":5,"reelKind":2,"reel":1,"symbol":3,"rows":[0,2]}`},
		{"first 3 6 1,2", 3, 6, utils.UInt8s{1, 2}, `{"kind":5,"reelKind":2,"reel":3,"symbol":6,"rows":[1,2]}`},
		{"next 2 10 0,1,2", 2, 10, utils.UInt8s{0, 1, 2}, `{"kind":5,"reelKind":2,"reel":2,"symbol":10,"rows":[0,1,2]}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := AcquireExpandReelEvent(tc.reel, tc.rows, tc.symbol).(*ReelEvent)
			require.NotNil(t, e)
			defer e.Release()

			assert.Equal(t, ExpandReelEvent, e.kind)
			assert.Equal(t, tc.reel, e.reel)
			assert.Equal(t, false, e.first)
			assert.Equal(t, tc.symbol, e.symbol)
			assert.EqualValues(t, tc.rows, e.rows)

			enc := zjson.AcquireEncoder(1024)
			defer enc.Release()

			e.Encode(enc)
			got := enc.Bytes()
			assert.Equal(t, tc.want, string(got))
		})
	}
}

func TestAcquireBlockedReelEvent(t *testing.T) {
	testCases := []struct {
		name  string
		reel  uint8
		first bool
		want  string
	}{
		{"first 1", 1, true, `{"kind":5,"reelKind":3,"reel":1,"firstTime":1}`},
		{"first 3", 3, true, `{"kind":5,"reelKind":3,"reel":3,"firstTime":1}`},
		{"next 2", 2, false, `{"kind":5,"reelKind":3,"reel":2}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := AcquireBlockedReelEvent(tc.reel, tc.first).(*ReelEvent)
			require.NotNil(t, a)
			assert.Equal(t, BlockedReelEvent, a.kind)
			assert.Equal(t, tc.reel, a.reel)
			assert.Equal(t, tc.first, a.first)
			assert.Zero(t, a.symbol)
			assert.Empty(t, a.rows)

			enc := zjson.AcquireEncoder(1024)
			defer enc.Release()

			a.Encode(enc)
			got := enc.Bytes()
			assert.Equal(t, tc.want, string(got))
		})
	}
}

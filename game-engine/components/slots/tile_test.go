package slots

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestNewTile(t *testing.T) {
	testCases := []struct {
		name       string
		offset     uint8
		symbol     utils.Index
		sticky     uint8
		multiplier uint16
	}{
		{name: "no sticky, no multiplier", offset: 1, symbol: 1},
		{name: "sticky, no multiplier", offset: 2, symbol: 2, sticky: 1},
		{name: "no sticky, multiplier", offset: 3, symbol: 3, multiplier: 5},
		{name: "sticky, multiplier", offset: 4, symbol: 4, sticky: 2, multiplier: 15},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tile := AcquireTile(tc.offset, tc.symbol, tc.sticky, tc.multiplier)
			require.NotNil(t, tile)
			defer tile.Release()

			assert.False(t, tile.IsJump())
			assert.Equal(t, tc.offset, tile.Offset())
			assert.Zero(t, tile.offsetFrom)
			assert.Equal(t, tc.symbol, tile.Symbol())
			assert.Equal(t, tc.sticky, tile.Sticky())
			assert.Equal(t, tc.multiplier, tile.Multiplier())

			enc := zjson.AcquireEncoder(1024)
			defer enc.Release()

			tile.Encode(enc)
			got := enc.Bytes()
			require.NotEmpty(t, got)

			dec := zjson.AcquireDecoder(got)
			defer dec.Release()

			tile2 := tileProducer.Acquire().(*Tile)
			defer tile2.Release()

			ok := dec.Object(tile2)
			require.True(t, ok)
			assert.NoError(t, dec.Error())
			assert.True(t, tile.DeepEqual(tile2))
		})
	}
}

func TestNewJumpedTile(t *testing.T) {
	testCases := []struct {
		name       string
		from       uint8
		to         uint8
		symbol     utils.Index
		sticky     uint8
		multiplier uint16
	}{
		{name: "no sticky, no multiplier", from: 1, to: 2, symbol: 1},
		{name: "sticky, no multiplier", from: 2, to: 4, symbol: 2, sticky: 1},
		{name: "no sticky, multiplier", from: 3, to: 6, symbol: 3, multiplier: 5},
		{name: "sticky, multiplier", from: 4, to: 8, symbol: 4, sticky: 2, multiplier: 15},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tile := AcquireJumpedTile(tc.from, tc.to, tc.symbol, tc.sticky, tc.multiplier)
			require.NotNil(t, tile)
			defer tile.Release()

			assert.True(t, tile.IsJump())
			assert.Equal(t, tc.to, tile.Offset())
			assert.Equal(t, tc.from, tile.offsetFrom)
			assert.Equal(t, tc.symbol, tile.Symbol())
			assert.Equal(t, tc.sticky, tile.Sticky())
			assert.Equal(t, tc.multiplier, tile.Multiplier())

			enc := zjson.AcquireEncoder(1024)
			defer enc.Release()

			tile.Encode(enc)
			got := enc.Bytes()
			require.NotEmpty(t, got)

			dec := zjson.AcquireDecoder(got)
			defer dec.Release()

			tile2 := tileProducer.Acquire().(*Tile)
			defer tile2.Release()

			ok := dec.Object(tile2)
			require.True(t, ok)
			assert.NoError(t, dec.Error())
			assert.True(t, tile.DeepEqual(tile2))
		})
	}
}

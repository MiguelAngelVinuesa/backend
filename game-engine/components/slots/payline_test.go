package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestNewPayline(t *testing.T) {
	testCases := []struct {
		name     string
		id       uint8
		rowCount uint8
		rows     utils.UInt8s
		offsets  []int
	}{
		{"3 empty", 0, 3, nil, nil},
		{"3 middle", 1, 3, utils.UInt8s{1, 1, 1, 1, 1}, []int{1, 4, 7, 10, 13}},
		{"3 top", 2, 3, utils.UInt8s{0, 0, 0, 0, 0}, []int{0, 3, 6, 9, 12}},
		{"3 bottom", 3, 3, utils.UInt8s{2, 2, 2, 2, 2}, []int{2, 5, 8, 11, 14}},
		{"3 cross", 4, 3, utils.UInt8s{0, 0, 1, 2, 2}, []int{0, 3, 7, 11, 14}},
		{"3 zigzag", 5, 3, utils.UInt8s{0, 1, 0, 1, 0}, []int{0, 4, 6, 10, 12}},
		{"5 middle", 1, 5, utils.UInt8s{2, 2, 2, 2, 2}, []int{2, 7, 12, 17, 22}},
		{"5 top", 2, 5, utils.UInt8s{0, 0, 0, 0, 0}, []int{0, 5, 10, 15, 20}},
		{"5 bottom", 3, 5, utils.UInt8s{4, 4, 4, 4, 4}, []int{4, 9, 14, 19, 24}},
		{"5 cross", 4, 5, utils.UInt8s{0, 1, 2, 3, 4}, []int{0, 6, 12, 18, 24}},
		{"5 zigzag", 5, 5, utils.UInt8s{1, 2, 3, 2, 1}, []int{1, 7, 13, 17, 21}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := NewPayline(tc.id, tc.rowCount, tc.rows...)
			require.NotNil(t, p)

			assert.Equal(t, tc.id, p.ID())
			if len(tc.rows) > 0 {
				assert.EqualValues(t, tc.rows, p.RowMap())
			}
			if len(tc.offsets) > 0 {
				assert.EqualValues(t, tc.offsets, p.offsets)
			}
		})
	}
}

func TestPurgePaylines(t *testing.T) {
	testCases := []struct {
		name     string
		paylines Paylines
		capacity int
		want     int
	}{
		{"nil, 2", nil, 2, 2},
		{"nil, 5", nil, 5, 5},
		{"2, 1", Paylines{&Payline{}, &Payline{}}, 1, 2},
		{"2, 2", Paylines{&Payline{}, &Payline{}}, 2, 2},
		{"2, 5", Paylines{&Payline{}, &Payline{}}, 5, 5},
		{"4, 1", Paylines{&Payline{}, &Payline{}, &Payline{}, &Payline{}}, 1, 4},
		{"4, 2", Paylines{&Payline{}, &Payline{}, &Payline{}, &Payline{}}, 2, 4},
		{"4, 5", Paylines{&Payline{}, &Payline{}, &Payline{}, &Payline{}}, 5, 5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			p := PurgePaylines(tc.paylines, tc.capacity)
			require.NotNil(t, p)
			assert.Zero(t, len(p))
			assert.Equal(t, tc.want, cap(p))
		})
	}
}

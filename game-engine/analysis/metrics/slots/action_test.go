package slots

import (
	"testing"

	"github.com/goccy/go-json"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAction(t *testing.T) {
	testCases := []struct {
		name       string
		actionID   int
		actionName string
		first      []bool
		second     []bool
		free       []bool
		free2      []bool
		want       *Action
		j          string
	}{
		{
			name:       "empty",
			actionID:   1,
			actionName: "something",
			want: &Action{
				ID:   1,
				Name: "something",
			},
		},
		{
			name:       "few first",
			actionID:   2,
			actionName: "something else",
			first:      []bool{false, false, true, true, false},
			want: &Action{
				ID:             2,
				Name:           "something else",
				TotalCount:     5,
				TotalTriggered: 2,
				FirstCount:     5,
				FirstTriggered: 2,
			},
		},
		{
			name:       "plenty first",
			actionID:   3,
			actionName: "hihi",
			first:      []bool{false, false, true, true, false, true, true, false, false, false, true},
			want: &Action{
				ID:             3,
				Name:           "hihi",
				TotalCount:     11,
				TotalTriggered: 5,
				FirstCount:     11,
				FirstTriggered: 5,
			},
		},
		{
			name:       "few second",
			actionID:   4,
			actionName: "hoho",
			second:     []bool{false, false, true, true, false},
			want: &Action{
				ID:              4,
				Name:            "hoho",
				TotalCount:      5,
				TotalTriggered:  2,
				SecondCount:     5,
				SecondTriggered: 2,
			},
		},
		{
			name:       "few free",
			actionID:   5,
			actionName: "hehe",
			free:       []bool{false, false, true, true, false},
			want: &Action{
				ID:             5,
				Name:           "hehe",
				TotalCount:     5,
				TotalTriggered: 2,
				FreeCount:      5,
				FreeTriggered:  2,
			},
		},
		{
			name:       "few second free",
			actionID:   6,
			actionName: "haha",
			free2:      []bool{false, false, true, true, false},
			want: &Action{
				ID:                  6,
				Name:                "haha",
				TotalCount:          5,
				TotalTriggered:      2,
				FreeSecondCount:     5,
				FreeSecondTriggered: 2,
			},
		},
		{
			name:       "few mixed",
			actionID:   7,
			actionName: "nono",
			first:      []bool{true, false, false, false},
			second:     []bool{false, true, true, false},
			free:       []bool{false, false, true, false},
			free2:      []bool{false, false, true, true},
			want: &Action{
				ID:                  7,
				Name:                "nono",
				TotalCount:          16,
				TotalTriggered:      6,
				FirstCount:          4,
				FirstTriggered:      1,
				SecondCount:         4,
				SecondTriggered:     2,
				FreeCount:           4,
				FreeTriggered:       1,
				FreeSecondCount:     4,
				FreeSecondTriggered: 2,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewAction(tc.actionID, tc.actionName, "", "")
			require.NotNil(t, a)
			defer a.Release()

			for _, r := range tc.first {
				a.IncreaseFirst(r)
			}
			for _, r := range tc.second {
				a.IncreaseSecond(r)
			}
			for _, r := range tc.free {
				a.IncreaseFree(r)
			}
			for _, r := range tc.free2 {
				a.IncreaseSecondFree(r)
			}

			if !tc.want.Equals(a) {
				assert.EqualValues(t, tc.want, a)
			}

			if tc.j != "" {
				j, err := json.Marshal(a)
				require.NoError(t, err)
				assert.Equal(t, tc.j, string(j))
			}

			n := a.Clone().(*Action)
			require.NotNil(t, n)
			defer n.Release()

			if !tc.want.Equals(n) {
				assert.EqualValues(t, tc.want, n)
			}

			n.ResetData()
			assert.Zero(t, n.TotalCount)
			assert.Zero(t, n.FirstCount)
			assert.Zero(t, n.SecondCount)
			assert.Zero(t, n.FreeCount)
			assert.Zero(t, n.FreeSecondCount)
			assert.Zero(t, n.TotalTriggered)
			assert.Zero(t, n.FirstTriggered)
			assert.Zero(t, n.SecondTriggered)
			assert.Zero(t, n.FreeTriggered)
			assert.Zero(t, n.FreeSecondTriggered)
		})
	}
}

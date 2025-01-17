package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"
)

func TestNewPlayerChoiceAction(t *testing.T) {
	testCases := []struct {
		name    string
		flag    int
		key     string
		values  []string
		results []int
		choices map[string]string
		trigger bool
		want    int
	}{
		{
			name:    "no choice",
			flag:    1,
			key:     "wing",
			values:  []string{"north", "south"},
			results: []int{1, 2},
			choices: map[string]string{},
		},
		{
			name:    "wrong key",
			flag:    2,
			key:     "wing",
			values:  []string{"north", "south"},
			results: []int{1, 2},
			choices: map[string]string{"wong": "south"},
		},
		{
			name:    "wrong value",
			flag:    3,
			key:     "wing",
			values:  []string{"north", "south"},
			results: []int{1, 2},
			choices: map[string]string{"wing": "east"},
		},
		{
			name:    "triggered 1",
			flag:    4,
			key:     "wing",
			values:  []string{"north", "south"},
			results: []int{1, 2},
			choices: map[string]string{"wing": "north"},
			trigger: true,
			want:    1,
		},
		{
			name:    "triggered 2",
			flag:    5,
			key:     "wing",
			values:  []string{"north", "south"},
			results: []int{1, 2},
			choices: map[string]string{"wing": "south"},
			trigger: true,
			want:    2,
		},
		{
			name:    "multiple, not triggered",
			flag:    6,
			key:     "wing",
			values:  []string{"north", "south"},
			results: []int{3, 4},
			choices: map[string]string{"heelo": "world", "sequence": "5", "wang": "north"},
		},
		{
			name:    "multiple, triggered",
			flag:    7,
			key:     "wing",
			values:  []string{"north", "south"},
			results: []int{5, 6},
			choices: map[string]string{"heelo": "world", "sequence": "5", "wing": "south"},
			trigger: true,
			want:    6,
		},
		{
			name:    "integer, not triggered",
			flag:    8,
			key:     "wing",
			values:  []string{"1", "2"},
			results: []int{7, 8},
			choices: map[string]string{"wing": "4"},
		},
		{
			name:    "integer, triggered",
			flag:    8,
			key:     "wing",
			values:  []string{"1", "2"},
			results: []int{7, 8},
			choices: map[string]string{"wing": "1"},
			trigger: true,
			want:    7,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewPlayerChoiceAction(tc.flag, tc.key, tc.values, tc.results)
			require.NotNil(t, a)

			assert.Equal(t, tc.flag, a.flag)

			prng := rng.NewRNG()
			defer prng.ReturnToPool()

			game := NewSlots(Grid(5, 3), WithSymbols(setF1))

			spin := AcquireSpin(game, prng)
			defer spin.Release()

			got := a.TestChoices(spin, tc.choices)
			if tc.trigger {
				assert.NotNil(t, got)
				assert.Equal(t, tc.want, spin.roundFlags[tc.flag])
			} else {
				assert.Nil(t, got)
			}
		})
	}
}

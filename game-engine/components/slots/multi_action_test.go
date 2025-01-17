package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/rng"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestNewMultiActionFlagValue(t *testing.T) {
	prng := rng.AcquireRNG()
	defer prng.ReturnToPool()

	sx := NewSymbol(10, WithKind(Wild))
	set := NewSymbolSet(sf1, sf2, sf3, sf4, sf5, sf6, sx)
	slots := NewSlots(Grid(5, 3), WithSymbols(set))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	a1 := NewGenerateSymbolAction(10, []float64{100})
	a2 := NewGenerateSymbolAction(10, []float64{100})
	a3 := NewGenerateSymbolAction(10, []float64{100})
	a4 := NewGenerateSymbolAction(10, []float64{100})
	a5 := NewGenerateSymbolAction(10, []float64{100})

	testCases := []struct {
		name   string
		flag   int
		params []any
		value  int
		want   bool
	}{
		{
			name:   "zero",
			flag:   1,
			params: []any{1, a1, 2, a2, 3, a3, 4, a4, 5, a5},
		},
		{
			name:   "a1",
			flag:   1,
			params: []any{1, a1, 2, a2, 3, a3, 4, a4, 5, a5},
			value:  1,
			want:   true,
		},
		{
			name:   "a2",
			flag:   1,
			params: []any{1, a1, 2, a2, 3, a3, 4, a4, 5, a5},
			value:  2,
			want:   true,
		},
		{
			name:   "a3",
			flag:   1,
			params: []any{1, a1, 2, a2, 3, a3, 4, a4, 5, a5},
			value:  3,
			want:   true,
		},
		{
			name:   "a4",
			flag:   1,
			params: []any{1, a1, 2, a2, 3, a3, 4, a4, 5, a5},
			value:  4,
			want:   true,
		},
		{
			name:   "a5",
			flag:   1,
			params: []any{1, a1, 2, a2, 3, a3, 4, a4, 5, a5},
			value:  5,
			want:   true,
		},
		{
			name:   "too big",
			flag:   1,
			params: []any{1, a1, 2, a2, 3, a3, 4, a4, 5, a5},
			value:  6,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewMultiActionFlagValue(tc.flag, tc.params...)
			require.NotNil(t, a)
			assert.Equal(t, len(tc.params)/2, len(a.actions))

			spin.indexes = utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3}
			spin.roundFlags[tc.flag] = tc.value

			got := a.Triggered(spin)
			if tc.want {
				assert.NotNil(t, got)
			} else {
				assert.Nil(t, got)
			}
		})
	}
}

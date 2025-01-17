package slots

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/wheel"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestNewInstantBonusAction(t *testing.T) {
	testCases := []struct {
		name    string
		chance  float64
		tease   bool
		choice  string
		options []string
	}{
		{
			name:   "basic",
			chance: 25.0,
		},
		{
			name:   "tease",
			chance: 12.5,
			tease:  true,
		},
		{
			name:    "choice",
			chance:  10,
			choice:  "bonus",
			options: []string{"left", "middle", "right"},
		},
	}

	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewInstantBonusAction(tc.chance)
			require.NotNil(t, a)

			if tc.choice != "" {
				a.WithPlayerChoice(tc.choice, tc.options...)
			}
			if tc.tease {
				a.WithTease()
			}

			assert.False(t, a.selector)
			assert.False(t, a.wheel)
			assert.Equal(t, tc.chance, a.chance)
			assert.Equal(t, tc.tease, a.tease)

			if tc.choice != "" {
				assert.True(t, a.playerChoice)
				assert.Equal(t, tc.choice, a.choice)
				assert.EqualValues(t, tc.options, a.options)
			} else {
				assert.False(t, a.playerChoice)
				assert.Empty(t, a.choice)
				assert.Empty(t, a.options)
			}

			counts := make(map[bool]int)
			for ix := 0; ix < 10000; ix++ {
				if a.Triggered(spin) != nil {
					counts[true] = counts[true] + 1
				} else {
					counts[false] = counts[false] + 1
				}
			}

			assert.NotZero(t, counts[true])
			assert.NotZero(t, counts[false])

			got := float64(counts[true]*100) / float64(counts[false]+counts[true])
			assert.GreaterOrEqual(t, got, tc.chance*0.9)
			assert.LessOrEqual(t, got, tc.chance*1.1)
		})
	}
}

func TestNewBonusSelectorAction(t *testing.T) {
	weightsDedup3 := utils.AcquireWeightingDedup3().AddWeights(utils.Indexes{1, 2, 3}, []float64{65, 25, 10})
	defer weightsDedup3.Release()

	testCases := []struct {
		name    string
		weights utils.WeightedGenerator
		count   int
		choice  int
		flag    int
		flags   []int
		want    bool
	}{
		{
			name:    "no choice",
			weights: weightsDedup3,
			count:   3,
			choice:  0,
			flag:    1,
			flags:   []int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			want:    true,
		},
		{
			name:    "first choice",
			weights: weightsDedup3,
			count:   3,
			choice:  0,
			flag:    1,
			flags:   []int{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			want:    true,
		},
		{
			name:    "second choice",
			weights: weightsDedup3,
			count:   3,
			choice:  0,
			flag:    1,
			flags:   []int{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			want:    true,
		},
		{
			name:    "third choice",
			weights: weightsDedup3,
			count:   3,
			choice:  0,
			flag:    1,
			flags:   []int{3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			want:    true,
		},
	}

	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			a := NewBonusSelectorAction(tc.weights, tc.count, tc.choice, tc.flag)
			require.NotNil(t, a)

			assert.True(t, a.selector)
			assert.False(t, a.wheel)
			assert.Equal(t, tc.count, a.selCount)
			assert.Equal(t, tc.choice, a.selChoiceFlag)
			assert.Equal(t, tc.flag, a.selFlag)

			copy(spin.roundFlags, tc.flags)

			a2 := a.Triggered(spin)
			require.Equal(t, a2, a)

			for ix := 0; ix < 1000; ix++ {
				got := a2.BonusSelect(spin)
				if !tc.want {
					require.Nil(t, got)
				} else {
					require.NotNil(t, got)

					data, ok := got.(*results.BonusSelector)
					require.True(t, ok)
					require.NotNil(t, data)

					player := tc.flags[tc.choice]
					assert.Equal(t, player, int(data.PlayerChoice()))

					res := data.Results()
					assert.Equal(t, 3, len(res))

					chosen := res[player-1]
					assert.Equal(t, chosen, data.Chosen())

					got.Release()
				}
			}
		})
	}
}

func TestNewInstantBonusWheelAction(t *testing.T) {
	indexes1 := utils.Indexes{1}
	weights1 := []float64{100}
	indexes3 := utils.Indexes{1, 2, 3}
	weights3 := []float64{65, 25, 10}
	indexes6 := utils.Indexes{1, 2, 3, 4, 5, 6}
	weights6 := []float64{57, 25, 10, 5, 2, 1}

	testCases := []struct {
		name    string
		flag    int
		indexes utils.Indexes
		weights []float64
	}{
		{
			name:    "one",
			flag:    1,
			indexes: indexes1,
			weights: weights1,
		},
		{
			name:    "three",
			flag:    3,
			indexes: indexes3,
			weights: weights3,
		},
		{
			name:    "six",
			flag:    6,
			indexes: indexes6,
			weights: weights6,
		},
	}

	prng := rng.NewRNG()
	defer prng.ReturnToPool()

	slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

	spin := AcquireSpin(slots, prng)
	defer spin.Release()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := utils.AcquireWeighting().AddWeights(tc.indexes, tc.weights)
			defer w.Release()

			a := NewInstantBonusWheelAction(tc.flag, w)
			require.NotNil(t, a)
			assert.False(t, a.selector)
			assert.True(t, a.wheel)
			assert.Equal(t, tc.flag, a.wheelFlag)
			assert.Equal(t, w, a.wheelWeights)

			spin.ResetSpin()

			for ix := 0; ix < 1000; ix++ {
				data := a.BonusGame(spin)
				require.NotNil(t, data)

				result, ok := data.(*wheel.BonusWheelResult)
				require.True(t, ok)
				require.NotNil(t, result)

				v := result.Result()
				assert.Equal(t, spin.roundFlags[tc.flag], int(v))
				assert.GreaterOrEqual(t, v, tc.indexes[0])
				assert.LessOrEqual(t, v, tc.indexes[len(tc.indexes)-1])

				for iy := range spin.roundFlags {
					if iy != tc.flag {
						assert.Zero(t, spin.roundFlags[iy])
					}
				}

				data.Release()
			}
		})
	}
}

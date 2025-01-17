package slots

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestNewMultiplierScaleAction(t *testing.T) {
	t.Run("new multiplier scale", func(t *testing.T) {
		a := NewMultiplierScaleAction(10, 1, 2, 2, 2, 5, 5, 5, 5, 5, 10, 10, 10, 10, 10, 10, 10, 50)
		require.NotNil(t, a)
		assert.Equal(t, utils.Index(10), a.triggerSymbol)
		assert.Equal(t, 1, a.firstLevel)

		slots := NewSlots(Grid(5, 3), WithSymbols(setF1))

		prng := rng.NewRNG()
		defer prng.ReturnToPool()

		spin := AcquireSpin(slots, prng)
		defer spin.Release()

		spin.indexes = utils.Indexes{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3}
		a2 := a.Triggered(spin)
		require.Nil(t, a2)
		assert.Zero(t, spin.progressLevel)
		assert.Zero(t, spin.multiplier)

		spin.indexes = utils.Indexes{1, 2, 3, 4, 5, 6, 10, 2, 3, 4, 5, 6, 1, 2, 3}
		a2 = a.Triggered(spin)
		require.NotNil(t, a2)
		assert.Equal(t, 1, spin.progressLevel)
		assert.Equal(t, 2.0, spin.multiplier)

		spin.indexes = utils.Indexes{1, 2, 3, 4, 5, 6, 10, 2, 3, 4, 5, 6, 1, 2, 3}
		a2 = a.Triggered(spin)
		require.NotNil(t, a2)
		assert.Equal(t, 2, spin.progressLevel)
		assert.Equal(t, 2.0, spin.multiplier)

		spin.indexes = utils.Indexes{1, 2, 3, 4, 5, 6, 10, 2, 3, 4, 5, 6, 1, 2, 3}
		a2 = a.Triggered(spin)
		require.NotNil(t, a2)
		assert.Equal(t, 3, spin.progressLevel)
		assert.Equal(t, 2.0, spin.multiplier)

		spin.indexes = utils.Indexes{1, 2, 3, 4, 5, 6, 10, 2, 3, 4, 5, 6, 1, 2, 3}
		a2 = a.Triggered(spin)
		require.NotNil(t, a2)
		assert.Equal(t, 4, spin.progressLevel)
		assert.Equal(t, 5.0, spin.multiplier)

		spin.indexes = utils.Indexes{10, 2, 3, 4, 5, 6, 10, 2, 3, 4, 5, 6, 10, 2, 3}
		a2 = a.Triggered(spin)
		require.NotNil(t, a2)
		assert.Equal(t, 7, spin.progressLevel)
		assert.Equal(t, 5.0, spin.multiplier)

		spin.indexes = utils.Indexes{10, 2, 3, 4, 5, 6, 10, 2, 3, 4, 5, 6, 10, 2, 3}
		a2 = a.Triggered(spin)
		require.NotNil(t, a2)
		assert.Equal(t, 10, spin.progressLevel)
		assert.Equal(t, 10.0, spin.multiplier)

		spin.indexes = utils.Indexes{10, 2, 3, 4, 5, 6, 10, 2, 3, 4, 5, 6, 10, 2, 3}
		a2 = a.Triggered(spin)
		require.NotNil(t, a2)
		assert.Equal(t, 13, spin.progressLevel)
		assert.Equal(t, 10.0, spin.multiplier)

		spin.indexes = utils.Indexes{10, 2, 3, 4, 5, 6, 10, 2, 3, 4, 5, 6, 10, 2, 3}
		a2 = a.Triggered(spin)
		require.NotNil(t, a2)
		assert.Equal(t, 16, spin.progressLevel)
		assert.Equal(t, 50.0, spin.multiplier)

		spin.indexes = utils.Indexes{10, 2, 3, 4, 5, 6, 10, 2, 3, 4, 5, 6, 10, 2, 3}
		a2 = a.Triggered(spin)
		require.NotNil(t, a2)
		assert.Equal(t, 16, spin.progressLevel)
		assert.Equal(t, 50.0, spin.multiplier)
	})
}

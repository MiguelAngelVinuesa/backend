package cards

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTable(t *testing.T) {
	t.Run("new table", func(t *testing.T) {
		d := NewDeck(SevenUp(), WithJokers(2))
		require.NotNil(t, d)
		defer d.Release()

		assert.Equal(t, 34, len(d.cards))

		b := NewTable(d, "John", "Billy", "Bob", "Jack", "Harry")
		require.NotNil(t, b)
		defer b.Release()

		assert.Equal(t, d, b.deck)
		assert.Equal(t, 5, len(b.players))
		assert.Equal(t, 5, len(b.hands))
		assert.Equal(t, 5, b.PlayerCount())
		assert.False(t, b.IsPlaying())
	})
}

func TestTable_AddPlayers(t *testing.T) {
	testCases := []struct {
		name    string
		initial []string
		add     []string
		want    []string
	}{
		{"none+none", nil, nil, []string{}},
		{"1+none", []string{"X"}, nil, []string{"X"}},
		{"2+none", []string{"Y", "X"}, nil, []string{"Y", "X"}},
		{"none+1", nil, []string{"Y"}, []string{"Y"}},
		{"none+2", nil, []string{"Y", "X"}, []string{"Y", "X"}},
		{"1+1", []string{"A"}, []string{"Y"}, []string{"A", "Y"}},
		{"2+2", []string{"B", "A"}, []string{"Y", "X"}, []string{"B", "A", "Y", "X"}},
	}

	d := NewDeck(SevenUp(), WithJokers(2))
	require.NotNil(t, d)
	defer d.Release()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b := NewTable(d, tc.initial...)
			require.NotNil(t, b)
			defer b.Release()

			assert.Equal(t, len(tc.initial), len(b.players))

			b2 := b.AddPlayers(tc.add...)
			assert.Equal(t, b, b2)
			assert.Equal(t, tc.want, b.players)
			assert.Equal(t, len(tc.want), len(b.hands))
		})
	}
}

func TestTable_AddPlayers_Error(t *testing.T) {
	t.Run("add players error", func(t *testing.T) {
		d := NewDeck(SevenUp(), WithJokers(2))
		require.NotNil(t, d)
		defer d.Release()

		assert.Equal(t, 34, len(d.cards))

		b := NewTable(d, "John", "Billy", "Bob", "Jack")
		require.NotNil(t, b)
		defer b.Release()

		require.Equal(t, 4, len(b.players))

		b.AddPlayers("Harry")
		require.Equal(t, 5, len(b.players))

		b.Burn()

		b.AddPlayers("Jill")
		require.Equal(t, 5, len(b.players))
	})
}

func TestTable_Player(t *testing.T) {
	t.Run("get player by seat", func(t *testing.T) {
		d := NewDeck(SevenUp(), WithJokers(2))
		require.NotNil(t, d)
		defer d.Release()

		assert.Equal(t, 34, len(d.cards))

		b := NewTable(d, "John", "Billy", "Bob", "Jack", "Harry")
		require.NotNil(t, b)
		defer b.Release()

		h := b.Player(-1)
		assert.Nil(t, h)

		h = b.Player(0)
		assert.NotNil(t, h)

		h = b.Player(4)
		assert.NotNil(t, h)

		h = b.Player(5)
		assert.Nil(t, h)
	})
}

func TestTable_PlayerByID(t *testing.T) {
	t.Run("get player by ID", func(t *testing.T) {
		d := NewDeck(SevenUp(), WithJokers(2))
		require.NotNil(t, d)
		defer d.Release()

		assert.Equal(t, 34, len(d.cards))

		b := NewTable(d, "John", "Billy", "Bob", "Jack", "Harry")
		require.NotNil(t, b)
		defer b.Release()

		i, h := b.PlayerByID("?")
		assert.Equal(t, -1, i)
		assert.Nil(t, h)

		i, h = b.PlayerByID("Jack")
		assert.Equal(t, 3, i)
		assert.NotNil(t, h)

		i, h = b.PlayerByID("Billy")
		assert.Equal(t, 1, i)
		assert.NotNil(t, h)
	})
}

func TestTable_Draw(t *testing.T) {
	t.Run("draw", func(t *testing.T) {
		d := NewDeck(SevenUp(), WithJokers(2))
		require.NotNil(t, d)
		defer d.Release()

		assert.Equal(t, 34, len(d.cards))

		b := NewTable(d, "John", "Billy", "Bob", "Jack", "Harry")
		require.NotNil(t, b)
		defer b.Release()

		assert.False(t, b.IsPlaying())

		b2 := b.Draw(2)
		assert.Equal(t, b, b2)
		assert.Equal(t, 24, b.Remaining())
		assert.True(t, b.IsPlaying())

		for ix := 0; ix < 5; ix++ {
			h := b.Player(ix)
			assert.NotNil(t, h)
			assert.Equal(t, 2, h.Size())
		}

		b2 = b.Draw(0)
		assert.Equal(t, b, b2)
		assert.Equal(t, 19, b.Remaining())
		assert.True(t, b.IsPlaying())

		for ix := 0; ix < 5; ix++ {
			h := b.Player(ix)
			assert.NotNil(t, h)
			assert.Equal(t, 3, h.Size())
		}
	})
}

func TestTable_DrawExhausted(t *testing.T) {
	t.Run("draw exhausted", func(t *testing.T) {
		d := NewDeck(SevenUp(), WithJokers(2))
		require.NotNil(t, d)
		defer d.Release()

		assert.Equal(t, 34, len(d.cards))

		b := NewTable(d, "John", "Billy", "Bob", "Jack", "Harry")
		require.NotNil(t, b)
		defer b.Release()

		b2 := b.Draw(7)
		assert.Equal(t, b, b2)
		assert.Zero(t, b.Remaining())
		assert.True(t, b.IsPlaying())

		for ix := 0; ix < 5; ix++ {
			h := b.Player(ix)
			assert.NotNil(t, h)
			if ix < 4 {
				assert.Equal(t, 7, h.Size())
			} else {
				assert.Equal(t, 6, h.Size())
			}
		}
	})
}

func TestTable_PlayerDraw(t *testing.T) {
	t.Run("player draw", func(t *testing.T) {
		d := NewDeck(SevenUp(), WithJokers(2))
		require.NotNil(t, d)
		defer d.Release()

		assert.Equal(t, 34, len(d.cards))

		b := NewTable(d, "John", "Billy", "Bob", "Jack", "Harry")
		require.NotNil(t, b)
		defer b.Release()

		b2 := b.Draw(2)
		assert.Equal(t, b, b2)
		assert.Equal(t, 24, b.Remaining())
		assert.True(t, b.IsPlaying())

		for ix := 0; ix < 5; ix++ {
			h := b.Player(ix)
			assert.NotNil(t, h)
			assert.Equal(t, 2, h.Size())
		}

		b2 = b.PlayerDraw(-1, 99)
		assert.Equal(t, b, b2)
		assert.True(t, b.IsPlaying())

		b2 = b.PlayerDraw(0, 1)
		assert.Equal(t, b, b2)
		h := b.Player(0)
		require.NotNil(t, h)
		assert.Equal(t, 3, h.Size())
		assert.True(t, b.IsPlaying())

		b2 = b.PlayerDraw(0, 0)
		assert.Equal(t, b, b2)
		assert.Equal(t, 4, h.Size())
		assert.True(t, b.IsPlaying())
	})
}

func TestTable_PlayerDrawExhausted(t *testing.T) {
	t.Run("player draw exhausted", func(t *testing.T) {
		d := NewDeck(SevenUp(), WithJokers(2))
		require.NotNil(t, d)
		defer d.Release()

		assert.Equal(t, 34, len(d.cards))

		b := NewTable(d, "John", "Billy", "Bob", "Jack", "Harry")
		require.NotNil(t, b)
		defer b.Release()

		b2 := b.Draw(6)
		assert.Equal(t, b, b2)
		assert.Equal(t, 4, b.Remaining())
		assert.True(t, b.IsPlaying())

		for ix := 0; ix < 5; ix++ {
			h := b.Player(ix)
			assert.NotNil(t, h)
			assert.Equal(t, 6, h.Size())
		}

		b2 = b.PlayerDraw(0, 5)
		assert.Equal(t, b, b2)
		h := b.Player(0)
		require.NotNil(t, h)
		assert.Equal(t, 10, h.Size())
		assert.True(t, b.IsPlaying())
	})
}

func TestTable_Burn(t *testing.T) {
	t.Run("draw", func(t *testing.T) {
		d := NewDeck(SevenUp(), WithJokers(2))
		require.NotNil(t, d)
		defer d.Release()

		assert.Equal(t, 34, len(d.cards))

		b := NewTable(d, "John", "Billy", "Bob", "Jack", "Harry")
		require.NotNil(t, b)
		defer b.Release()

		for ix := 0; ix < 12; ix++ {
			b.Burn()
			assert.True(t, b.IsPlaying())
		}

		b2 := b.Draw(5)
		assert.Equal(t, b, b2)
		assert.Zero(t, b.Remaining())
		assert.True(t, b.IsPlaying())

		for ix := 0; ix < 5; ix++ {
			h := b.Player(ix)
			assert.NotNil(t, h)
			if ix < 2 {
				assert.Equal(t, 5, h.Size())
			} else {
				assert.Equal(t, 4, h.Size())
			}
		}
	})
}

func TestTable_Shuffle(t *testing.T) {
	t.Run("draw", func(t *testing.T) {
		d := NewDeck(SevenUp(), WithJokers(2))
		require.NotNil(t, d)
		assert.Equal(t, 34, len(d.cards))
		defer d.Release()

		b := NewTable(d, "John", "Billy", "Bob", "Jack", "Harry")
		require.NotNil(t, b)
		defer b.Release()

		b2 := b.Draw(5)
		assert.Equal(t, b, b2)
		assert.Equal(t, 9, b.Remaining())
		assert.True(t, b.IsPlaying())

		for ix := 0; ix < 5; ix++ {
			h := b.Player(ix)
			assert.NotNil(t, h)
			assert.Equal(t, 5, h.Size())
		}

		b2 = b.Shuffle()
		assert.Equal(t, b, b2)
		assert.Equal(t, 34, b.Remaining())
		assert.False(t, b.IsPlaying())

		for ix := 0; ix < 5; ix++ {
			h := b.Player(ix)
			assert.NotNil(t, h)
			assert.Zero(t, h.Size())
		}
	})
}

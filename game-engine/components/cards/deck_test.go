package cards

import (
	"log"
	"math/rand"
	"sync"
	"testing"
	"time"

	rng2 "git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/interfaces"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/rng"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	rng.AcquireRNG = func() interfaces.Generator { return rng2.NewRNG() }
}

func TestNewDeck(t *testing.T) {
	custom := func() []*Card {
		return []*Card{
			NewCard(GameCard0+0, WithColor(Red), WithValue(100)),
			NewCard(GameCard0+1, WithColor(Black), WithValue(200)),
			NewCard(GameCard0+2, WithColor(Red), WithValue(300)),
			NewCard(GameCard0+3, WithColor(Black), WithValue(400)),
		}
	}

	testCases := []struct {
		name      string
		opts      []DeckOption
		size      int
		remaining int
		shuffled  bool
	}{
		{"empty", nil, 0, 0, false},
		{"standard", []DeckOption{StandardDeck()}, 52, 52, false},
		{"seven-up", []DeckOption{SevenUp()}, 32, 32, false},
		{"just jokers", []DeckOption{WithJokers(99)}, 12, 12, false},
		{"standard with 3 jokers", []DeckOption{StandardDeck(), WithJokers(3)}, 55, 55, false},
		{"only hearts", []DeckOption{WithCards(HeartsAll...)}, 13, 13, false},
		{"only diamonds", []DeckOption{WithCards(DiamondsAll...)}, 13, 13, false},
		{"only clubs", []DeckOption{WithCards(ClubsAll...)}, 13, 13, false},
		{"only spades", []DeckOption{WithCards(SpadesAll...)}, 13, 13, false},
		{"standard shuffled", []DeckOption{StandardDeck(), Shuffled()}, 52, 52, true},
		{"custom cards", []DeckOption{WithCustom(custom()...)}, 4, 4, false},
		{"all", []DeckOption{StandardDeck(), WithJokers(3), WithCustom(custom()...), Shuffled()}, 59, 59, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := NewDeck(tc.opts...)
			require.NotNil(t, d)
			defer d.Release()

			assert.Equal(t, tc.size, d.Size())
			assert.Equal(t, tc.remaining, d.Remaining())
			assert.Equal(t, tc.shuffled, d.Shuffled())

			if tc.shuffled {
				assert.NotEqualValues(t, d.cards, d.Remain())
			} else {
				assert.EqualValues(t, d.cards, d.Remain())
			}
		})
	}
}

func TestDeckFromPool(t *testing.T) {
	t.Run("deck from pool", func(t *testing.T) {
		max := 8192

		wg := sync.WaitGroup{}
		wg.Add(max)

		for ix := 0; ix < max; ix++ {
			go func() {
				d := NewDeck(StandardDeck(), WithJokers(3), WithCustom(NewCard(GameCard0, WithValue(100))))
				require.NotNil(t, d)
				require.True(t, d.Unique(), d.String())
				time.Sleep(time.Microsecond * time.Duration(1+rand.Int31n(400)))
				require.True(t, d.Unique(), d.String())
				d.Release()
				wg.Done()
			}()
		}

		wg.Wait()
	})
}

func TestDeck_Add(t *testing.T) {
	t.Run("add card", func(t *testing.T) {
		d := NewDeck()
		require.NotNil(t, d)
		defer d.Release()

		d.Add(HeartsAll...)
		assert.Equal(t, len(HeartsAll), d.Size())

		d.Add(JokerX, JokerY)
		assert.Equal(t, len(HeartsAll)+2, d.Size())
	})
}

func TestDeck_AddCustom(t *testing.T) {
	t.Run("add custom", func(t *testing.T) {
		d := NewDeck()
		require.NotNil(t, d)
		defer d.Release()

		d.AddCustom(NewCard(Heart2))
		d.AddCustom(NewCard(Diamond2))
		d.AddCustom(NewCard(Spade2))
		d.AddCustom(NewCard(Club2))
		assert.Equal(t, 4, d.Size())

		d.AddCustom(NewCard(GameCard0+0, WithColor(Black), WithValue(100)))
		d.AddCustom(NewCard(GameCard0+1, WithColor(Red), WithValue(200)))
		assert.Equal(t, 6, d.Size())
	})
}

func TestDeck_GetRandomCard(t *testing.T) {
	t.Run("random card 1/2 deck", func(t *testing.T) {
		d := NewDeck(StandardDeck(), WithJokers(3), WithCustom(NewCard(GameCard0, WithValue(100))))
		require.NotNil(t, d)
		require.True(t, d.Unique(), d.String())
		defer d.Release()

		max := 1000000
		size := d.Size()
		cnt := make(map[CardID]int)

		for ix := 0; ix < max; ix++ {
			if d.Remaining() <= d.Size()/2 {
				d.Reset()
				require.True(t, d.Unique(), d.String())
			}
			c := d.GetRandomCard()
			cnt[c.id] = cnt[c.id] + 1
		}

		assert.Equal(t, size, len(cnt))

		avg := max / size
		low := avg * 9 / 10
		high := avg * 11 / 10
		for _, c := range cnt {
			assert.GreaterOrEqual(t, c, low)
			assert.LessOrEqual(t, c, high)
		}
	})

	t.Run("random card 1/4 deck", func(t *testing.T) {
		d := NewDeck(StandardDeck(), WithJokers(3), WithCustom(NewCard(GameCard0, WithValue(100))))
		require.NotNil(t, d)
		require.True(t, d.Unique(), d.String())
		defer d.Release()

		max := 1000000
		size := d.Size()
		cnt := make(map[CardID]int)

		for ix := 0; ix < max; ix++ {
			if d.Remaining() <= d.Size()*3/4 {
				d.Reset()
				require.True(t, d.Unique(), d.String())
			}
			c := d.GetRandomCard()
			cnt[c.id] = cnt[c.id] + 1
		}

		assert.Equal(t, size, len(cnt))

		avg := max / size
		low := avg * 9 / 10
		high := avg * 11 / 10
		for _, c := range cnt {
			assert.GreaterOrEqual(t, c, low)
			assert.LessOrEqual(t, c, high)
		}
	})
}

func TestDeck_GetRandomCardEmpty(t *testing.T) {
	t.Run("random card empty", func(t *testing.T) {
		d := NewDeck(StandardDeck())
		require.NotNil(t, d)
		defer d.Release()

		for ix := 0; ix < d.Size(); ix++ {
			d.GetRandomCard()
		}

		c := d.GetRandomCard()
		assert.Nil(t, c)
	})
}

func TestDeck_Shuffle(t *testing.T) {
	t.Run("shuffle", func(t *testing.T) {
		d := NewDeck(SevenUp())
		require.NotNil(t, d)
		assert.False(t, d.shuffled)
		defer d.Release()

		d.Shuffle()
		assert.True(t, d.shuffled)
		assert.Equal(t, d.Remaining(), d.Size())
		assert.NotEqualValues(t, d.cards, d.remain)

		d.Reset()
		assert.False(t, d.shuffled)
		assert.Equal(t, d.cards, d.Remain())

		d.Shuffle()
		assert.True(t, d.shuffled)
		assert.Equal(t, d.Remaining(), d.Size())
		assert.NotEqualValues(t, d.cards, d.Remain())
	})
}

func TestDeck_ShuffleRemaining(t *testing.T) {
	t.Run("shuffle remaining", func(t *testing.T) {
		d := NewDeck(SevenUp())
		require.NotNil(t, d)
		assert.False(t, d.shuffled)
		defer d.Release()

		d.Shuffle()
		assert.True(t, d.shuffled)

		d.DrawMulti(20)
		assert.Equal(t, d.Remaining(), d.Size()-20)

		c1 := d.Remain()
		c2 := c1

		for ix := 0; ix < 10000; ix++ {
			d.ShuffleRemaining()
			assert.True(t, d.Shuffled())

			c3 := d.Remain()

			assert.NotEqualValues(t, c1, c3)
			assert.NotEqualValues(t, c2, c3)

			c2 = c3
		}
	})
}

func TestDeck_Cut(t *testing.T) {
	t.Run("cut", func(t *testing.T) {
		d := NewDeck(SevenUp())
		require.NotNil(t, d)
		defer d.Release()

		assert.Equal(t, d.Remaining(), d.Size())
		assert.Equal(t, d.cards, d.Remain())

		max := d.Size()

		d.Cut(17)
		assert.Equal(t, d.Remaining(), d.Size())
		assert.NotEqual(t, d.cards, d.Remain())

		d.Cut(max - 17)
		assert.Equal(t, d.Remaining(), d.Size())
		assert.Equal(t, d.cards, d.Remain())
	})
}

func TestDeck_CutPanic(t *testing.T) {
	t.Run("cut panic 1", func(t *testing.T) {
		defer func() {
			e := recover()
			require.NotNil(t, e)
		}()

		d := NewDeck(SevenUp())
		defer d.Release()
		d.Cut(0)
	})

	t.Run("cut panic 2", func(t *testing.T) {
		defer func() {
			e := recover()
			require.NotNil(t, e)
		}()

		d := NewDeck(SevenUp())
		defer d.Release()
		d.Cut(d.Size())
	})
}

func TestDeck_CutRandom(t *testing.T) {
	t.Run("cut random", func(t *testing.T) {
		d := NewDeck(SevenUp())
		require.NotNil(t, d)
		defer d.Release()

		assert.Equal(t, d.Remaining(), d.Size())
		assert.EqualValues(t, d.cards, d.Remain())

		d.CutRandom(0)
		assert.Equal(t, d.Remaining(), d.Size())
		assert.NotEqualValues(t, d.cards, d.Remain())

		d.CutRandom(5)
		assert.Equal(t, d.Remaining(), d.Size())
		assert.NotEqualValues(t, d.cards, d.Remain())

		d.CutRandom(2000)
		assert.Equal(t, d.Remaining(), d.Size())
		assert.NotEqualValues(t, d.cards, d.Remain())
	})
}

func TestDeck_CutRandomPanic(t *testing.T) {
	t.Run("cut random panic", func(t *testing.T) {
		defer func() {
			e := recover()
			require.NotNil(t, e)
		}()

		d := NewDeck(SevenUp(), Shuffled())
		defer d.Release()
		for ix := 0; ix < d.Size(); ix++ {
			d.Draw()
		}
		d.CutRandom(2)
	})
}

func TestDeck_Draw(t *testing.T) {
	t.Run("draw 1/2 deck", func(t *testing.T) {
		d := NewDeck(StandardDeck(), WithJokers(3), WithCustom(NewCard(GameCard0, WithValue(100))), Shuffled())
		require.NotNil(t, d)
		require.True(t, d.Unique(), d.String())
		defer d.Release()

		max := 1000000
		size := d.Size()
		cnt := make(map[CardID]int)

		for ix := 0; ix < max; ix++ {
			if d.Remaining() <= d.Size()/2 {
				d.Shuffle()
				require.True(t, d.Unique(), d.String())
			}
			c := d.Draw()
			cnt[c.id] = cnt[c.id] + 1
		}

		if len(cnt) != size {
			log.Printf("size: %d; cnt: %v\n", size, cnt)
		}

		assert.Equal(t, size, len(cnt))

		avg := max / size
		low := avg * 9 / 10
		high := avg * 11 / 10

		for _, c := range cnt {
			assert.GreaterOrEqual(t, c, low)
			assert.LessOrEqual(t, c, high)
		}
	})

	t.Run("draw 1/4 deck", func(t *testing.T) {
		d := NewDeck(StandardDeck(), WithJokers(3), WithCustom(NewCard(GameCard0, WithValue(100))), Shuffled())
		require.NotNil(t, d)
		defer d.Release()

		max := 1000000
		size := d.Size()
		cnt := make(map[CardID]int)

		for ix := 0; ix < max; ix++ {
			if d.Remaining() <= d.Size()*3/4 {
				d.Shuffle()
			}
			c := d.Draw()
			cnt[c.id] = cnt[c.id] + 1
		}

		assert.Equal(t, size, len(cnt))

		avg := max / size
		low := avg * 9 / 10
		high := avg * 11 / 10
		for _, c := range cnt {
			assert.GreaterOrEqual(t, c, low)
			assert.LessOrEqual(t, c, high)
		}
	})
}

func TestDeck_DrawEmpty(t *testing.T) {
	t.Run("draw empty", func(t *testing.T) {
		d := NewDeck(StandardDeck(), Shuffled())
		require.NotNil(t, d)
		defer d.Release()

		for ix := 0; ix < d.Size(); ix++ {
			d.Draw()
		}

		c := d.Draw()
		assert.Nil(t, c)
	})
}

func TestDeck_DrawMulti(t *testing.T) {
	t.Run("draw multi", func(t *testing.T) {
		d := NewDeck(SevenUp(), Shuffled())
		require.NotNil(t, d)
		defer d.Release()

		list := d.DrawMulti(3)
		require.NotNil(t, list)
		assert.Equal(t, 3, len(list))
		assert.Equal(t, 29, d.Remaining())

		list = d.DrawMulti(10)
		require.NotNil(t, list)
		assert.Equal(t, 10, len(list))
		assert.Equal(t, 19, d.Remaining())

		list = d.DrawMulti(14)
		require.NotNil(t, list)
		assert.Equal(t, 14, len(list))
		assert.Equal(t, 5, d.Remaining())

		list = d.DrawMulti(8)
		require.NotNil(t, list)
		assert.Equal(t, 5, len(list))
		assert.Equal(t, 0, d.Remaining())

		list = d.DrawMulti(4)
		require.NotNil(t, list)
		assert.Equal(t, 0, len(list))
		assert.Equal(t, 0, d.Remaining())
	})
}

func TestDeck_DrawInto(t *testing.T) {
	t.Run("draw multi", func(t *testing.T) {
		d := NewDeck(SevenUp(), Shuffled())
		require.NotNil(t, d)
		defer d.Release()

		list := make(Cards, 3, 32)
		list = d.DrawInto(list)
		require.NotNil(t, list)
		assert.Equal(t, 3, len(list))
		assert.Equal(t, 29, d.Remaining())

		list = list[:10]
		list = d.DrawInto(list)
		require.NotNil(t, list)
		assert.Equal(t, 10, len(list))
		assert.Equal(t, 19, d.Remaining())

		list = list[:14]
		list = d.DrawInto(list)
		require.NotNil(t, list)
		assert.Equal(t, 14, len(list))
		assert.Equal(t, 5, d.Remaining())

		list = list[:8]
		list = d.DrawInto(list)
		require.NotNil(t, list)
		assert.Equal(t, 5, len(list))
		assert.Equal(t, 0, d.Remaining())

		list = list[:4]
		list = d.DrawInto(list)
		require.NotNil(t, list)
		assert.Equal(t, 0, len(list))
		assert.Equal(t, 0, d.Remaining())
	})
}

func TestDeck_Burn(t *testing.T) {
	t.Run("burn", func(t *testing.T) {
		d := NewDeck(SevenUp(), Shuffled())
		require.NotNil(t, d)
		defer d.Release()

		i := d.Burn()
		assert.Equal(t, 31, i)
		assert.Equal(t, 31, d.Remaining())

		d.DrawMulti(5)

		i = d.Burn()
		assert.Equal(t, 25, i)
		assert.Equal(t, 25, d.Remaining())

		d.DrawMulti(10)

		i = d.Burn()
		assert.Equal(t, 14, i)
		assert.Equal(t, 14, d.Remaining())

		d.DrawMulti(7)

		i = d.Burn()
		assert.Equal(t, 6, i)
		assert.Equal(t, 6, d.Remaining())

		d.DrawMulti(7)

		i = d.Burn()
		assert.Equal(t, 0, i)
		assert.Equal(t, 0, d.Remaining())
	})
}

func TestDeck_Unique(t *testing.T) {
	testCases := []struct {
		name string
		add  []CardID
		want bool
	}{
		{"none", []CardID{}, true},
		{"joker", []CardID{Joker0}, true},
		{"jokers", []CardID{Joker0, Joker1, Joker2}, true},
		{"Ace of hearts", []CardID{HeartA}, false},
		{"4 of clubs", []CardID{Club4}, false},
		{"dup jokers", []CardID{Joker0, Joker5, Joker0}, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := NewDeck(StandardDeck(), Shuffled())
			require.NotNil(t, d)
			assert.True(t, d.Unique())

			d.Add(tc.add...)
			assert.Equal(t, tc.want, d.Unique())

			for ix := 0; ix < 10; ix++ {
				d.Shuffle()
				assert.Equal(t, tc.want, d.Unique())
			}
		})
	}
}

func BenchmarkRandomCardStandard(b *testing.B) {
	d := NewDeck(StandardDeck())
	defer d.Release()
	for i := 0; i < b.N; i++ {
		if d.Remaining() == 0 {
			d.Reset()
		}
		d.GetRandomCard()
	}
}

func BenchmarkDrawStandard(b *testing.B) {
	d := NewDeck(StandardDeck(), Shuffled())
	defer d.Release()
	for i := 0; i < b.N; i++ {
		if d.Remaining() == 0 {
			d.Shuffle()
		}
		d.Draw()
	}
}

func BenchmarkRandomCardSevenUp(b *testing.B) {
	d := NewDeck(SevenUp())
	defer d.Release()
	for i := 0; i < b.N; i++ {
		if d.Remaining() == 0 {
			d.Reset()
		}
		d.GetRandomCard()
	}
}

func BenchmarkDrawSevenUp(b *testing.B) {
	d := NewDeck(SevenUp(), Shuffled())
	defer d.Release()
	for i := 0; i < b.N; i++ {
		if d.Remaining() == 0 {
			d.Shuffle()
		}
		d.Draw()
	}
}

func BenchmarkShuffleFY(b *testing.B) {
	d := NewDeck(StandardDeck())
	defer d.Release()
	for i := 0; i < b.N; i++ {
		d.Shuffle()
	}
}

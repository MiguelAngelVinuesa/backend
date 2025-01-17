package poker

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/cards"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	H2 = cards.NewCard(cards.Heart2)
	H3 = cards.NewCard(cards.Heart3)
	H4 = cards.NewCard(cards.Heart4)
	H5 = cards.NewCard(cards.Heart5)
	H6 = cards.NewCard(cards.Heart6)
	H7 = cards.NewCard(cards.Heart7)
	H8 = cards.NewCard(cards.Heart8)
	H9 = cards.NewCard(cards.Heart9)
	HX = cards.NewCard(cards.HeartX)
	HJ = cards.NewCard(cards.HeartJ)
	HQ = cards.NewCard(cards.HeartQ)
	HK = cards.NewCard(cards.HeartK)
	HA = cards.NewCard(cards.HeartA)
	C2 = cards.NewCard(cards.Club2)
	C3 = cards.NewCard(cards.Club3)
	C4 = cards.NewCard(cards.Club4)
	C5 = cards.NewCard(cards.Club5)
	C6 = cards.NewCard(cards.Club6)
	CX = cards.NewCard(cards.ClubX)
	CJ = cards.NewCard(cards.ClubJ)
	CQ = cards.NewCard(cards.ClubQ)
	CK = cards.NewCard(cards.ClubK)
	CA = cards.NewCard(cards.ClubA)
	S2 = cards.NewCard(cards.Spade2)
	S3 = cards.NewCard(cards.Spade3)
	S4 = cards.NewCard(cards.Spade4)
	S5 = cards.NewCard(cards.Spade5)
	S6 = cards.NewCard(cards.Spade6)
	SX = cards.NewCard(cards.SpadeX)
	SJ = cards.NewCard(cards.SpadeJ)
	SQ = cards.NewCard(cards.SpadeQ)
	SK = cards.NewCard(cards.SpadeK)
	SA = cards.NewCard(cards.SpadeA)
	D4 = cards.NewCard(cards.Diamond4)
)

func TestNewRanked(t *testing.T) {
	t.Run("new ranked", func(t *testing.T) {
		r := NewRanked(nil)
		require.NotNil(t, r)
		assert.Equal(t, Rank(0), r.rank)
		assert.Equal(t, "No rank", r.rank.String())
	})
}

func TestNewRankHand(t *testing.T) {
	testCases := []struct {
		name   string
		cards  cards.Cards
		rank   Rank
		rankS  string
		ranked cards.Cards
		remain cards.Cards
	}{
		{"3, high card", cards.Cards{C4, SJ, CX}, HighCard, "High card", cards.Cards{SJ}, cards.Cards{CX, C4}},
		{"3, one pair", cards.Cards{C4, SJ, H4}, OnePair, "One pair", cards.Cards{H4, C4}, cards.Cards{SJ}},
		{"3, three of a kind", cards.Cards{C4, S4, H4}, ThreeOfAKind, "Three of a kind", cards.Cards{S4, H4, C4}, nil},
		{"3, flush", cards.Cards{H9, H3, H6}, Flush, "Flush", cards.Cards{H9, H6, H3}, nil},
		{"3, straight (1)", cards.Cards{S2, H3, SA}, Straight, "Straight", cards.Cards{H3, S2, SA}, nil},
		{"3, straight (2)", cards.Cards{C4, C6, H5}, Straight, "Straight", cards.Cards{C6, H5, C4}, nil},
		{"3, straight flush (1)", cards.Cards{H5, H6, H4}, StraightFlush, "Straight flush", cards.Cards{H6, H5, H4}, nil},
		{"3, straight flush (2)", cards.Cards{H3, HA, H2}, StraightFlush, "Straight flush", cards.Cards{H3, H2, HA}, nil},
		{"3, royal flush", cards.Cards{HK, HQ, HA}, RoyalFlush, "Royal flush", cards.Cards{HA, HK, HQ}, nil},
		{"5, high card", cards.Cards{C4, SJ, CX, H6, SA}, HighCard, "High card", cards.Cards{SA}, cards.Cards{SJ, CX, H6, C4}},
		{"5, one pair", cards.Cards{C4, SJ, H4, HQ, S6}, OnePair, "One pair", cards.Cards{H4, C4}, cards.Cards{HQ, SJ, S6}},
		{"5, two pair", cards.Cards{C4, SJ, H4, HJ, H9}, TwoPair, "Two pair", cards.Cards{SJ, HJ, H4, C4}, cards.Cards{H9}},
		{"5, three of a kind", cards.Cards{C4, H6, S4, SJ, H4}, ThreeOfAKind, "Three of a kind", cards.Cards{S4, H4, C4}, cards.Cards{SJ, H6}},
		{"5, four of a kind", cards.Cards{C4, C6, D4, S4, H4}, FourOfAKind, "Four of a kind", cards.Cards{S4, H4, C4, D4}, cards.Cards{C6}},
		{"5, full house", cards.Cards{C4, S4, H4, HX, SX}, FullHouse, "Full house", cards.Cards{S4, H4, C4, SX, HX}, nil},
		{"5, flush", cards.Cards{H9, HA, H3, HJ, H6}, Flush, "Flush", cards.Cards{HA, HJ, H9, H6, H3}, nil},
		{"5, straight (1)", cards.Cards{C4, C6, H5, C3, S2}, Straight, "Straight", cards.Cards{C6, H5, C4, C3, S2}, nil},
		{"5, straight (2)", cards.Cards{C4, SA, H5, C3, S2}, Straight, "Straight", cards.Cards{H5, C4, C3, S2, SA}, nil},
		{"5, straight flush (1)", cards.Cards{H4, H6, H5, H8, H7}, StraightFlush, "Straight flush", cards.Cards{H8, H7, H6, H5, H4}, nil},
		{"5, straight flush (1)", cards.Cards{H4, H2, H5, HA, H3}, StraightFlush, "Straight flush", cards.Cards{H5, H4, H3, H2, HA}, nil},
		{"5, royal flush", cards.Cards{HK, HX, HJ, HQ, HA}, RoyalFlush, "Royal flush", cards.Cards{HA, HK, HQ, HJ, HX}, nil},
		{"7, high card", cards.Cards{C4, SJ, C5, H6, SQ, CK, SA}, HighCard, "High card", cards.Cards{SA}, cards.Cards{CK, SQ, SJ, H6, C5, C4}},
		{"7, one pair", cards.Cards{CQ, C4, SJ, H4, SK, H9, S6}, OnePair, "One pair", cards.Cards{H4, C4}, cards.Cards{SK, CQ, SJ, H9, S6}},
		{"7, two pair", cards.Cards{C4, SJ, H4, HJ, H9, SA, H7}, TwoPair, "Two pair", cards.Cards{SJ, HJ, H4, C4}, cards.Cards{SA, H9, H7}},
		{"7, three of a kind", cards.Cards{C4, S5, H6, C2, S4, SJ, H4}, ThreeOfAKind, "Three of a kind", cards.Cards{S4, H4, C4}, cards.Cards{SJ, H6, S5, C2}},
		{"7, four of a kind (1)", cards.Cards{C4, C6, SQ, D4, S4, CA, H4}, FourOfAKind, "Four of a kind", cards.Cards{S4, H4, C4, D4}, cards.Cards{CA, SQ, C6}},
		{"7, four of a kind (2)", cards.Cards{C4, HA, SQ, D4, S4, CA, H4}, FourOfAKind, "Four of a kind", cards.Cards{S4, H4, C4, D4}, cards.Cards{HA, CA, SQ}},
		{"7, four of a kind (3)", cards.Cards{C4, HA, SA, D4, S4, CA, H4}, FourOfAKind, "Four of a kind", cards.Cards{S4, H4, C4, D4}, cards.Cards{SA, HA, CA}},
		{"7, full house (1)", cards.Cards{C4, S4, SA, H4, H9, HX, SX}, FullHouse, "Full house", cards.Cards{S4, H4, C4, SX, HX}, cards.Cards{SA, H9}},
		{"7, full house (2)", cards.Cards{C2, S2, S4, H2, H4, HX, SX}, FullHouse, "Full house", cards.Cards{S2, H2, C2, SX, HX}, cards.Cards{S4, H4}},
		{"7, full house (3)", cards.Cards{C4, S4, S2, H4, H2, HX, SX}, FullHouse, "Full house", cards.Cards{S4, H4, C4, SX, HX}, cards.Cards{S2, H2}},
		{"7, full house (4)", cards.Cards{CX, S4, S2, H4, H2, HX, SX}, FullHouse, "Full house", cards.Cards{SX, HX, CX, S4, H4}, cards.Cards{S2, H2}},
		{"7, full house (5)", cards.Cards{C4, S5, H5, C5, S4, SJ, H4}, FullHouse, "Full house", cards.Cards{S5, H5, C5, S4, H4}, cards.Cards{SJ, C4}},
		{"7, flush", cards.Cards{H9, HA, SA, SK, H3, HJ, H6}, Flush, "Flush", cards.Cards{HA, HJ, H9, H6, H3}, cards.Cards{SA, SK}},
		{"7, straight (1)", cards.Cards{C4, CJ, C6, S3, H5, C3, S2}, Straight, "Straight", cards.Cards{C6, H5, C4, S3, S2}, cards.Cards{CJ, C3}},
		{"7, straight (2)", cards.Cards{C4, SA, H5, HA, D4, C3, S2}, Straight, "Straight", cards.Cards{H5, C4, C3, S2, SA}, cards.Cards{HA, D4}},
		{"7, straight flush (1)", cards.Cards{H4, H6, H3, H2, H5, H8, H7}, StraightFlush, "Straight flush", cards.Cards{H8, H7, H6, H5, H4}, cards.Cards{H3, H2}},
		{"7, straight flush (1)", cards.Cards{H4, SA, D4, H2, H5, HA, H3}, StraightFlush, "Straight flush", cards.Cards{H5, H4, H3, H2, HA}, cards.Cards{SA, D4}},
		{"7, royal flush", cards.Cards{HK, HX, HJ, H9, SA, HQ, HA}, RoyalFlush, "Royal flush", cards.Cards{HA, HK, HQ, HJ, HX}, cards.Cards{SA, H9}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := RankHand(tc.cards)
			require.NotNil(t, r)
			assert.Equal(t, tc.rank, r.Rank())
			assert.Equal(t, tc.rankS, r.Rank().String())

			if tc.ranked != nil {
				require.Equal(t, len(tc.ranked), len(r.Ranked()))
				for ix, card := range tc.ranked {
					assert.Equal(t, card, r.Ranked()[ix])
				}
			}

			if tc.remain != nil {
				require.Equal(t, len(tc.remain), len(r.Remaining()))
				for ix, card := range tc.remain {
					assert.Equal(t, card, r.Remaining()[ix])
				}
			}
		})
	}
}

func TestRankHandEmpty(t *testing.T) {
	t.Run("rank empty hand", func(t *testing.T) {
		r := RankHand(nil)
		require.NotNil(t, r)
		assert.Equal(t, Rank(0), r.rank)
		assert.Equal(t, "No rank", r.rank.String())
	})
}

func TestRanked_Reload(t *testing.T) {
	t.Run("rank reload", func(t *testing.T) {
		r := RankHand(nil)
		require.NotNil(t, r)
		assert.Equal(t, Rank(0), r.rank)

		r2 := r.Reload(cards.Cards{H2, S2, SA, C3, C2, D4, HA})
		require.Equal(t, r2, r)
		assert.Equal(t, FullHouse, r.Rank())
		assert.Equal(t, cards.Cards{S2, H2, C2, SA, HA}, r.Ranked())
		assert.Equal(t, cards.Cards{D4, C3}, r.Remaining())

		r3 := r.Reload(cards.Cards{C2, SA})
		require.Equal(t, r3, r)
		assert.Equal(t, Rank(0), r.rank)
	})
}

func BenchmarkRankHighcard(b *testing.B) {
	hand := cards.Cards{C4, SJ, C5, H6, SQ, CK, SA}
	for i := 0; i < b.N; i++ {
		RankHand(hand).Release()
	}
}

func BenchmarkRankOnePair(b *testing.B) {
	hand := cards.Cards{C4, SJ, C5, H6, SQ, CJ, SA}
	for i := 0; i < b.N; i++ {
		RankHand(hand).Release()
	}
}

func BenchmarkRankTwoPair(b *testing.B) {
	hand := cards.Cards{C4, SJ, H4, H6, SQ, CJ, SA}
	for i := 0; i < b.N; i++ {
		RankHand(hand).Release()
	}
}

func BenchmarkRankThreeOfAKind(b *testing.B) {
	hand := cards.Cards{C4, SJ, H4, H6, SQ, D4, SA}
	for i := 0; i < b.N; i++ {
		RankHand(hand).Release()
	}
}

func BenchmarkRankFourOfAKind(b *testing.B) {
	hand := cards.Cards{C4, SJ, H4, H6, S4, D4, SA}
	for i := 0; i < b.N; i++ {
		RankHand(hand).Release()
	}
}

func BenchmarkRankFullHouse(b *testing.B) {
	hand := cards.Cards{C4, SJ, H4, S4, SQ, CJ, SA}
	for i := 0; i < b.N; i++ {
		RankHand(hand).Release()
	}
}

func BenchmarkRankFlush(b *testing.B) {
	hand := cards.Cards{HQ, H5, HJ, S4, H8, H7, SA}
	for i := 0; i < b.N; i++ {
		RankHand(hand).Release()
	}
}

func BenchmarkRankStraight(b *testing.B) {
	hand := cards.Cards{C4, H5, S6, S4, H8, H7, SA}
	for i := 0; i < b.N; i++ {
		RankHand(hand).Release()
	}
}

func BenchmarkRankStraightLowAce(b *testing.B) {
	hand := cards.Cards{C2, H5, SJ, S4, H3, D4, SA}
	for i := 0; i < b.N; i++ {
		RankHand(hand).Release()
	}
}

func BenchmarkRankStraightFlush(b *testing.B) {
	hand := cards.Cards{C4, H5, H6, H4, H8, H7, SA}
	for i := 0; i < b.N; i++ {
		RankHand(hand).Release()
	}
}

func BenchmarkRankStraightFlushLowAce(b *testing.B) {
	hand := cards.Cards{H2, H5, SA, S4, H3, H4, HA}
	for i := 0; i < b.N; i++ {
		RankHand(hand).Release()
	}
}

func BenchmarkRankRoyalFlush(b *testing.B) {
	hand := cards.Cards{SJ, HK, HA, HX, HJ, HQ, SA}
	for i := 0; i < b.N; i++ {
		RankHand(hand).Release()
	}
}

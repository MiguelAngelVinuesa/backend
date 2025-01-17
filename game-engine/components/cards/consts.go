package cards

// CardID represents the "id" of a card.
// "id div 16" represents the suit, and "id mod 16" represents the ordinal card value.
// The ordinal values 0, 1 and 15 (regardless of suit) are reserved for Joker cards.
// All values from 64 upwards are reserved for special cards that can be used in a game dependent way.
type CardID uint8

// Standard diamonds suit.
const (
	DiamondA CardID = iota + 1
	Diamond2
	Diamond3
	Diamond4
	Diamond5
	Diamond6
	Diamond7
	Diamond8
	Diamond9
	DiamondX
	DiamondJ
	DiamondQ
	DiamondK
)

// Standard clubs suit.
const (
	ClubA CardID = iota + 17
	Club2
	Club3
	Club4
	Club5
	Club6
	Club7
	Club8
	Club9
	ClubX
	ClubJ
	ClubQ
	ClubK
)

// Standard hearts suit.
const (
	HeartA CardID = iota + 33
	Heart2
	Heart3
	Heart4
	Heart5
	Heart6
	Heart7
	Heart8
	Heart9
	HeartX
	HeartJ
	HeartQ
	HeartK
)

// Standard spades suit.
const (
	SpadeA CardID = iota + 49
	Spade2
	Spade3
	Spade4
	Spade5
	Spade6
	Spade7
	Spade8
	Spade9
	SpadeX
	SpadeJ
	SpadeQ
	SpadeK
)

// Standard joker cards.
const (
	Joker0 CardID = 0
	Joker1 CardID = 16
	Joker2 CardID = 32
	Joker3 CardID = 48
	Joker4 CardID = 14
	Joker5 CardID = 30
	Joker6 CardID = 46
	Joker7 CardID = 62
	Joker8 CardID = 15
	Joker9 CardID = 31
	JokerX CardID = 47
	JokerY CardID = 63
)

// GameCard0 is the basis for game specific card id's.
const GameCard0 CardID = 64

// Sets of cards used to initialize a deck.
var (
	DiamondsAll = []CardID{DiamondA, Diamond2, Diamond3, Diamond4, Diamond5, Diamond6, Diamond7, Diamond8, Diamond9, DiamondX, DiamondJ, DiamondQ, DiamondK}
	ClubsAll    = []CardID{ClubA, Club2, Club3, Club4, Club5, Club6, Club7, Club8, Club9, ClubX, ClubJ, ClubQ, ClubK}
	HeartsAll   = []CardID{HeartA, Heart2, Heart3, Heart4, Heart5, Heart6, Heart7, Heart8, Heart9, HeartX, HeartJ, HeartQ, HeartK}
	SpadesAll   = []CardID{SpadeA, Spade2, Spade3, Spade4, Spade5, Spade6, Spade7, Spade8, Spade9, SpadeX, SpadeJ, SpadeQ, SpadeK}
	Diamonds7up = []CardID{Diamond7, Diamond8, Diamond9, DiamondX, DiamondJ, DiamondQ, DiamondK, DiamondA}
	Clubs7up    = []CardID{Club7, Club8, Club9, ClubX, ClubJ, ClubQ, ClubK, ClubA}
	Hearts7up   = []CardID{Heart7, Heart8, Heart9, HeartX, HeartJ, HeartQ, HeartK, HeartA}
	Spades7up   = []CardID{Spade7, Spade8, Spade9, SpadeX, SpadeJ, SpadeQ, SpadeK, SpadeA}
	JokersAll   = []CardID{Joker0, Joker1, Joker2, Joker3, Joker4, Joker5, Joker6, Joker7, Joker8, Joker9, JokerX, JokerY}
)

// Ordinal represents the ordinal value of a card.
type Ordinal uint8

// The default ordinals for cards.
const (
	Ace Ordinal = iota + 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

// Suit represents the suit of a card.
// Valid values are 1..4 for respectively diamonds, clubs, hearts and spades.
// Joker cards are suit independent and will have the invalid suit value 0.
type Suit uint8

// The default suites for cards.
// The ordinal value of a suit represents the importance of suits commonly used in some games.
const (
	Diamonds Suit = iota + 1
	Clubs
	Hearts
	Spades
)

// Color represents the color of a card.
// Valid values are 1..2 for respectively red and black.
// Joker cards have a valid color.
type Color uint8

// The default colors for cards.
const (
	Red Color = iota + 1
	Black
)

const (
	invalidCutPosition   = "invalid position for cut function"
	notEnoughCardsForCut = "not enough remaining cards for cutting"
)

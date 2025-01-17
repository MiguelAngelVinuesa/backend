package cards

import (
	"bytes"
	"fmt"
	"strconv"
	"sync"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
)

// Card represents a single card.
type Card struct {
	id      CardID
	ordinal Ordinal
	joker   bool
	suit    Suit
	color   Color
	value   int
}

var cardPool = sync.Pool{
	New: func() interface{} {
		return &Card{}
	},
}

// NewCard instantiates a new card from the memory pool.
// The opts can be used to change the card characteristics as required.
// Cards are immutable once created, and can be used from multiple go-routines simultaneously.
func NewCard(id CardID, opts ...CardOption) *Card {
	c := cardPool.Get().(*Card)
	c.id = id

	if id < GameCard0 {
		c.ordinal = Ordinal(id & 0xf)
		c.value = int(c.ordinal)
		c.joker = c.ordinal == 0 || c.ordinal == 14 || c.ordinal == 15

		s := int8(id >> 4)
		if s == 0 || s == 2 {
			c.color = Red
		} else {
			c.color = Black
		}

		if c.IsJoker() {
			c.value = 15
		} else {
			c.suit = Suit(s) + Diamonds
		}
	} else {
		c.ordinal = Ordinal(id)
		c.value = int(id)
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Release puts the card back into the memory pool.
func (c *Card) Release() {
	if c != nil {
		c.reset()
		cardPool.Put(c)
	}
}

func (c *Card) reset() {
	c.id = 0
	c.ordinal = 0
	c.value = 0
	c.joker = false
	c.color = 0
	c.suit = 0
}

// CardOption is the function signature for Card option functions.
type CardOption = func(*Card)

// WithSuit is the card option to set the suit.
func WithSuit(suit Suit) CardOption {
	return func(c *Card) {
		c.suit = suit
	}
}

// WithColor is the card option to set the color.
func WithColor(color Color) CardOption {
	return func(c *Card) {
		c.color = color
	}
}

// WithJoker is the card option to toggle the joker flag.
func WithJoker(joker bool) CardOption {
	return func(c *Card) {
		c.joker = joker
	}
}

// WithOrdinal is the card option to set the ordinal value.
func WithOrdinal(ordinal Ordinal) CardOption {
	return func(c *Card) {
		c.ordinal = ordinal
	}
}

// WithValue is the card option to set the card value.
func WithValue(value int) CardOption {
	return func(c *Card) {
		c.value = value
	}
}

// ID returns the id of the card.
func (c *Card) ID() CardID {
	return c.id
}

// Ordinal returns the ordinal value of the card.
// For standard cards, it is by default "id mod 16".
// For game dependent cards it is by default the same as the id.
func (c *Card) Ordinal() Ordinal {
	return c.ordinal
}

// IsJoker returns whether the card is a joker card.
// This is by default true for ordinal values 0, 1 and 15.
func (c *Card) IsJoker() bool {
	return c.joker
}

// IsAce returns whether the card is an ace.
func (c *Card) IsAce() bool {
	return c.ordinal == Ace
}

// Suit returns the suit of the card.
// For standard cards it is determined from "id div 16".
// For joker cards and game dependent cards it returns the invalid zero value by default.
func (c *Card) Suit() Suit {
	return c.suit
}

// Color returns the color of the card.
// For standard cards it is determined by default by the suit of the card.
// For game dependent cards it returns the invalid zero value by default.
func (c *Card) Color() Color {
	return c.color
}

// Value returns the given value for the card.
// It is the same as the ordinal value by default.
func (c *Card) Value() int {
	return c.value
}

// String implements the Stringer interface.
func (c *Card) String() string {
	var b bytes.Buffer
	b.WriteString("{id:")
	b.WriteString(strconv.Itoa(int(c.id)))
	b.WriteString(",ordinal:")
	b.WriteString(strconv.Itoa(int(c.ordinal)))
	b.WriteString(",value:")
	b.WriteString(strconv.Itoa(c.value))
	b.WriteString(",suit:")
	b.WriteString(strconv.Itoa(int(c.suit)))
	b.WriteString(",color:")
	b.WriteString(strconv.Itoa(int(c.color)))
	b.WriteString(fmt.Sprintf(",joker:%t", c.joker))
	b.WriteByte('}')
	return b.String()
}

// IsEmpty implements the zjson encoder interface.
func (c *Card) IsEmpty() bool {
	return false
}

// EncodeFields implements the zjson encoder interface.
func (c *Card) EncodeFields(enc *zjson.Encoder) {
	enc.Uint8Field("id", uint8(c.id))
	enc.Uint8Field("ordinal", uint8(c.ordinal))
	if c.joker {
		enc.IntBoolField("joker", c.joker)
	}
	enc.Uint8Field("suit", uint8(c.suit))
	enc.Uint8Field("color", uint8(c.color))
	enc.IntFieldOpt("value", c.value)
}

// DecodeField implements the zjson decoder interface.
func (c *Card) DecodeField(dec *zjson.Decoder, key []byte) error {
	var ok bool
	var i uint8

	if string(key) == "id" {
		if i, ok = dec.Uint8(); ok {
			c.id = CardID(i)
		}
	} else if string(key) == "ordinal" {
		if i, ok = dec.Uint8(); ok {
			c.ordinal = Ordinal(i)
		}
	} else if string(key) == "joker" {
		c.joker, ok = dec.Bool()
	} else if string(key) == "suit" {
		if i, ok = dec.Uint8(); ok {
			c.suit = Suit(i)
		}
	} else if string(key) == "color" {
		if i, ok = dec.Uint8(); ok {
			c.color = Color(i)
		}
	} else if string(key) == "value" {
		c.value, ok = dec.Int()
	}

	if ok {
		return nil
	}
	return dec.Error()
}

// Cards represents a slice of cards.
type Cards []*Card

// ResliceCards creates or resizes an existing slice of cards to the requested capacity with zero elements.
func ResliceCards(cards Cards, need int) Cards {
	if cards == nil || cap(cards) < need {
		return make(Cards, 0, need)
	}
	return cards[:0]
}

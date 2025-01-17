package cards

import (
	"bytes"
	"fmt"
	"strconv"
	"sync"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/interfaces"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/rng"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/sharedlib"
)

// Deck represents a deck of cards.
//
// There are two methods to use a deck:
//  1. Instantiate a deck and then use the GetRandomCard() function.
//     Use ResetData() to start again.
//     This method is cheap for resets but moderately expensive for drawing cards.
//  2. Instantiate a deck, use Shuffle() to randomize it, and then use the Draw() function(s).
//     This method is expensive for shuffles but very cheap for drawing cards.
//
// Mixing the methods is not recommended!
// Which method to use depends on your use-case.
// For games where only a few draws are needed, method 1 works best.
// For longer lasting games method 2 is superior.
// It is recommended to use method 2 with the Fisher-Yates shuffle mechanism.
//
// A Deck is safe to use across multiple go-routines simultaneously.
type Deck struct {
	mutex    sync.RWMutex
	prng     interfaces.Generator
	shuffler interfaces.Shuffler
	cards    Cards
	buffer   []int
	remain   []int
	shuffled bool
}

var deckPool = sync.Pool{
	New: func() interface{} {
		return &Deck{}
	},
}

// NewDeck instantiates a new deck of cards from the memory pool.
func NewDeck(opts ...DeckOption) *Deck {
	d := deckPool.Get().(*Deck)
	d.prng = rng.AcquireRNG()
	d.shuffler = sharedlib.AcquireShuffler()
	if d.cards == nil {
		d.cards = make(Cards, 0, 32)
	}
	d.reset()
	for _, opt := range opts {
		opt(d)
	}
	return d
}

// Release puts a deck back into the memory pool.
func (d *Deck) Release() {
	if d != nil {
		d.shuffler.Release()
		d.shuffler = nil
		d.prng.ReturnToPool()
		d.prng = nil

		for _, c := range d.cards {
			c.Release()
		}

		d.cards = d.cards[:0]
		d.buffer = d.buffer[:0]
		d.remain = d.buffer
		d.shuffled = false

		deckPool.Put(d)
	}
}

// Add adds one or more cards to the deck.
func (d *Deck) Add(cards ...CardID) *Deck {
	d.mutex.Lock()
	d.add(cards)
	d.reset()
	d.mutex.Unlock()
	return d
}

func (d *Deck) add(cards []CardID) {
	// mutex should already be locked
	for _, id := range cards {
		d.cards = append(d.cards, NewCard(id))
	}
}

// AddCustom adds one or more custom-built cards to the deck.
// Deck takes ownership of the cards. Do not use the cards after this.
func (d *Deck) AddCustom(cards ...*Card) *Deck {
	d.mutex.Lock()
	d.addCustom(cards)
	d.reset()
	d.mutex.Unlock()
	return d
}

func (d *Deck) addCustom(cards []*Card) {
	// mutex should already be locked
	d.cards = append(d.cards, cards...)
}

// Size returns the size of the deck. E.g. how many cards there are in total.
func (d *Deck) Size() int {
	d.mutex.RLock()
	defer d.mutex.RUnlock()
	return len(d.cards)
}

// Unique returns whether all cards in the deck have unique identifiers.
func (d *Deck) Unique() bool {
	d.mutex.RLock()
	u := d.unique()
	d.mutex.RUnlock()
	return u
}

func (d *Deck) unique() bool {
	max := len(d.cards)
	for ix := 0; ix < max-1; ix++ {
		for iy := ix + 1; iy < max; iy++ {
			if d.cards[ix].id == d.cards[iy].id {
				return false
			}
		}
	}
	return true
}

// Remaining returns the count of remaining cards that can be drawn.
func (d *Deck) Remaining() int {
	d.mutex.RLock()
	l := len(d.remain)
	d.mutex.RUnlock()
	return l
}

// Shuffled returns whether the deck has been shuffled.
func (d *Deck) Shuffled() bool {
	d.mutex.RLock()
	s := d.shuffled
	d.mutex.RUnlock()
	return s
}

// Reset resets a deck for use with GetRandomCard().
func (d *Deck) Reset() *Deck {
	d.mutex.Lock()
	d.reset()
	d.mutex.Unlock()
	return d
}

func (d *Deck) reset() {
	// mutex should already be locked
	l := len(d.cards)
	if cap(d.buffer) < l {
		d.buffer = make([]int, l)
	} else {
		d.buffer = d.buffer[:l]
	}
	for ix := range d.cards {
		d.buffer[ix] = ix
	}
	d.remain = d.buffer
	d.shuffled = false
}

// GetRandomCard returns a random card from the deck and reduces the number of remaining cards.
// The function returns nil if there are no remaining cards.
func (d *Deck) GetRandomCard() *Card {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	max := len(d.remain)
	if max == 0 {
		return nil
	}

	ix := d.prng.IntN(max)
	c := d.remain[ix]

	switch ix {
	case 0:
		d.remain = d.remain[1:]
	case max - 1:
		d.remain = d.remain[:max-1]
	default:
		d.remain = append(d.remain[:ix], d.remain[ix+1:]...)
	}

	return d.cards[c]
}

// Shuffle resets the deck and shuffles it into a random order.
func (d *Deck) Shuffle() *Deck {
	d.mutex.Lock()
	d.reset()
	d.shuffle()
	d.mutex.Unlock()
	return d
}

// ShuffleRemaining shuffles the remaining cards into a random order.
func (d *Deck) ShuffleRemaining() *Deck {
	d.mutex.Lock()
	d.shuffle()
	d.mutex.Unlock()
	return d
}

// Cut splits the remaining cards at the given position and recombines the two parts in swapped order.
// The function panics if position <= 0 or >= Remaining()-1
func (d *Deck) Cut(position int) *Deck {
	if position <= 0 || position >= d.Remaining()-1 {
		panic(invalidCutPosition)
	}

	d.mutex.Lock()
	d.remain = append(d.remain[position:], d.remain[:position]...)
	d.mutex.Unlock()
	return d
}

// CutRandom splits the remaining cards at a random position and recombines the two parts in swapped order.
// It repeats the process count times.
// Count must be between 1 and 1000. It is adjusted accordingly if the value is outside the limits.
// The function panics if there are less than 6 remaining cards.
func (d *Deck) CutRandom(count int) *Deck {
	if d.Remaining() < 6 {
		panic(notEnoughCardsForCut)
	}

	if count < 1 {
		count = 1
	}
	if count > 1000 {
		count = 1000
	}

	d.mutex.Lock()
	max := len(d.remain) - 2

	for ix := 0; ix < count; ix++ {
		position := d.prng.IntN(max) + 1
		d.remain = append(d.remain[position:], d.remain[:position]...)
	}

	d.mutex.Unlock()
	return d
}

// Draw returns a single card from the deck and reduces the remaining cards.
// It returns nil if there are no more remaining cards.
// Make sure to call Shuffle() before drawing the first card, or you will get the cards in the order they were added to the deck.
func (d *Deck) Draw() *Card {
	d.mutex.Lock()

	var c *Card
	if len(d.remain) > 0 {
		c = d.cards[d.remain[0]]
		d.remain = d.remain[1:]
	}

	d.mutex.Unlock()
	return c
}

// DrawMulti returns multiple cards from the deck and reduces the remaining cards.
// It will return an empty slice if Remaining() is already 0.
// It will only return the remaining cards if Remaining() is smaller than the requested count.
// Make sure to call Shuffle() first, or you will get the cards in the order they were added to the deck.
func (d *Deck) DrawMulti(count int) Cards {
	d.mutex.Lock()

	if count > len(d.remain) {
		count = len(d.remain)
	}
	c := make(Cards, count)
	for ix := range c {
		c[ix] = d.cards[d.remain[ix]]
	}
	d.remain = d.remain[count:]

	d.mutex.Unlock()
	return c
}

// DrawInto returns multiple cards from the deck and reduces the remaining cards.
// It will fill the given slice with at most Remaining() cards.
// Make sure to call Shuffle() first, or you will get the cards in the order they were added to the deck.
func (d *Deck) DrawInto(out Cards) Cards {
	d.mutex.Lock()

	if len(out) > len(d.remain) {
		out = out[:len(d.remain)]
	}
	for ix := range out {
		out[ix] = d.cards[d.remain[ix]]
	}
	d.remain = d.remain[len(out):]

	d.mutex.Unlock()
	return out
}

// Burn removes a single card from the deck and reduces the remaining cards.
// It returns the number of remaining cards.
// Make sure to call Shuffle() before burning the first card, or you will burn cards in the order they were added to the deck.
func (d *Deck) Burn() int {
	d.mutex.Lock()
	l := len(d.remain)
	if l > 0 {
		d.remain = d.remain[1:]
		l--
	}
	d.mutex.Unlock()
	return l
}

// Remain returns the remaining cards in the deck.
func (d *Deck) Remain() Cards {
	l := len(d.remain)
	out := make(Cards, l)
	for ix := range out {
		out[ix] = d.cards[d.remain[ix]]
	}
	return out
}

func (d *Deck) String() string {
	var b bytes.Buffer

	d.mutex.RLock()

	b.WriteString("size:")
	b.WriteString(strconv.Itoa(len(d.cards)))
	b.WriteString(fmt.Sprintf(" shuffled:%t", d.shuffled))
	b.WriteString(" remaining:")
	b.WriteString(strconv.Itoa(len(d.remain)))
	b.WriteString(" cards:[")

	for ix, c := range d.cards {
		if ix > 0 {
			b.WriteByte(',')
		}
		b.WriteString(c.String())
	}

	b.WriteByte(']')

	d.mutex.RUnlock()
	return b.String()
}

// shuffle will shuffle the remaining cards.
func (d *Deck) shuffle() {
	d.shuffler.Shuffle(d.remain)
	d.shuffled = true
}

// DeckOption is the function signature for Deck option functions.
type DeckOption func(*Deck)

// StandardDeck initializes the deck with all standard cards from the 4 suits.
func StandardDeck() DeckOption {
	return func(d *Deck) {
		d.add(DiamondsAll)
		d.add(ClubsAll)
		d.add(HeartsAll)
		d.add(SpadesAll)
		d.reset()
	}
}

// SevenUp initializes the deck with the standard cards from the 4 suits from 7 and up.
func SevenUp() DeckOption {
	return func(d *Deck) {
		d.add(Diamonds7up)
		d.add(Clubs7up)
		d.add(Hearts7up)
		d.add(Spades7up)
		d.reset()
	}
}

// WithJokers adds one or more jokers to the deck.
// The maximum is 12 jokers.
func WithJokers(count int) DeckOption {
	if count > 12 {
		count = 12
	}
	return func(d *Deck) {
		d.add(JokersAll[:count])
		d.reset()
	}
}

// WithCards adds one or more default cards to the deck.
func WithCards(cards ...CardID) DeckOption {
	return func(d *Deck) {
		d.add(cards)
		d.reset()
	}
}

// WithCustom adds one or more custom cards to the deck.
func WithCustom(cards ...*Card) DeckOption {
	return func(d *Deck) {
		d.addCustom(cards)
		d.reset()
	}
}

// Shuffled shuffles the deck.
// Use it as the last option!
func Shuffled() DeckOption {
	return func(d *Deck) {
		d.shuffle()
	}
}

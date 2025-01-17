package cards

import (
	"math"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/cards"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
)

// AcquireBonusCallSuit instantiates a bonus "call suit" game from the memory pool.
func AcquireBonusCallSuit() *BonusCallSuit {
	return bonusCallSuitPool.Acquire().(*BonusCallSuit)
}

// RequireParams implements the BonusRunner interface.
func (b *BonusCallSuit) RequireParams() bool { return true }

// Run implements the BonusRunner interface.
// The first parameter must be the players choice of card suit. If not supplied the default is spades.
// The second parameter must be the position where to cut the deck. If not supplied or invalid the deck is not cut.
// The function returns a new payout result with payout 4 or 0, depending on whether the players chosen suit was drawn or not.
func (b *BonusCallSuit) Run(_ *results.Result, params ...interface{}) (int, *results.Result) {
	choice := cards.Spades
	if len(params) > 0 {
		if n, ok := params[0].(int); ok {
			switch n {
			case int(cards.Diamonds):
				choice = cards.Diamonds
			case int(cards.Clubs):
				choice = cards.Clubs
			case int(cards.Hearts):
				choice = cards.Hearts
			}
		}
	}

	var cut int
	if len(params) > 1 {
		if n, ok := params[1].(int); ok && n < b.deck.Size()-1 {
			cut = n
		}
	}

	if cut > 0 {
		b.deck.Cut(cut)
	}
	card := b.deck.Draw()

	data := bonusCallSuitDataPool.Acquire().(*BonusCallSuitData)
	data.Choice = choice
	data.Cut = cut
	data.Card = card

	var payout float64
	if card.Suit() == choice {
		payout = 4
	}

	p := results.AcquirePlayerChoice(payout)
	defer p.Release()
	return int(math.Round(payout * 100)), results.AcquireResult(data, results.BonusCallSuitData, p)
}

// BonusCallSuit is a 4x or nothing bonus game where the player must guess the suit of a card to be drawn.
// The card is drawn from a standard 52-card "French" deck.
// The player has the option to "cut" the deck before the card is drawn.
type BonusCallSuit struct {
	deck *cards.Deck
	pool.Object
}

var bonusCallSuitPool = pool.NewProducer(func() (pool.Objecter, func()) {
	b := &BonusCallSuit{
		deck: cards.NewDeck(cards.StandardDeck(), cards.Shuffled()),
	}
	return b, b.reset
})

// reset clears the bonus.
func (b *BonusCallSuit) reset() {
	if b != nil {
		b.deck.Release()
		b.deck = nil
	}
}

// BonusCallSuitData represents the details of a bonus "call suit" game.
type BonusCallSuitData struct {
	cardOwned bool
	Choice    cards.Suit
	Cut       int
	Card      *cards.Card
	pool.Object
}

var bonusCallSuitDataPool = pool.NewProducer(func() (pool.Objecter, func()) {
	d := &BonusCallSuitData{}
	return d, d.reset
})

// reset clears the bonus data.
func (d *BonusCallSuitData) reset() {
	if d != nil {
		if d.cardOwned {
			d.Card.Release()
			d.cardOwned = false
		}
		d.Choice = cards.Suit(0)
		d.Cut = 0
		d.Card = nil
	}
}

// IsEmpty implements the zjson encoder interface.
func (d *BonusCallSuitData) IsEmpty() bool {
	return false
}

// EncodeFields implements the zjson encoder interface.
func (d *BonusCallSuitData) EncodeFields(enc *zjson.Encoder) {
	enc.Uint8Field("choice", uint8(d.Choice))
	enc.IntField("cut", d.Cut)
	enc.ObjectField("card", d.Card)
}

// Encode2 implements the PoolRCZ.Encode2 interface.
func (d *BonusCallSuitData) Encode2(enc *zjson.Encoder) {
	d.EncodeFields(enc)
}

// DecodeField implements the zjson decoder interface.
func (d *BonusCallSuitData) DecodeField(dec *zjson.Decoder, key []byte) error {
	var ok bool
	var i uint8

	if string(key) == "choice" {
		if i, ok = dec.Uint8(); ok {
			d.Choice = cards.Suit(i)
		}
	} else if string(key) == "cut" {
		d.Cut, ok = dec.Int()
	} else if string(key) == "card" {
		d.Card = cards.NewCard(0)
		d.cardOwned = true
		ok = dec.Object(d.Card)
	}

	if ok {
		return nil
	}
	return dec.Error()
}

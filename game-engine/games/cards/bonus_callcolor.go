package cards

import (
	"math"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/cards"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
)

// AcquireBonusCallColor instantiates a bonus "call color" game from the memory pool.
func AcquireBonusCallColor() *BonusCallColor {
	return bonusCallColorPool.Acquire().(*BonusCallColor)
}

// RequireParams implements the BonusRunner interface.
func (b *BonusCallColor) RequireParams() bool { return true }

// Run implements the BonusRunner interface.
// The first parameter must be the players choice of card color. If not supplied the default is red.
// The second parameter must be the position where to cut the deck. If not supplied or invalid the deck is not cut.
// The function returns a new payout result with payout 2 or 0, depending on whether the players chosen color was drawn or not.
func (b *BonusCallColor) Run(_ *results.Result, params ...interface{}) (int, *results.Result) {
	choice := cards.Red
	if len(params) > 0 {
		if n, ok := params[0].(int); ok && n == int(cards.Black) {
			choice = cards.Black
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

	data := bonusCallColorDataPool.Acquire().(*BonusCallColorData)
	data.Choice = choice
	data.Cut = cut
	data.Card = card

	var payout float64
	if card.Color() == choice {
		payout = 2
	}

	p := results.AcquirePlayerChoice(payout)
	defer p.Release()
	return int(math.Round(payout * 100)), results.AcquireResult(data, results.BonusCallColorData, p)
}

// BonusCallColor is a double or nothing bonus game where the player must guess the color of a card to be drawn.
// The card is drawn from a standard 52-card "French" deck.
// The player has the option to "cut" the deck before the card is drawn.
type BonusCallColor struct {
	deck *cards.Deck
	pool.Object
}

var bonusCallColorPool = pool.NewProducer(func() (pool.Objecter, func()) {
	b := &BonusCallColor{
		deck: cards.NewDeck(cards.StandardDeck(), cards.Shuffled()),
	}
	return b, b.reset
})

// reset clears the bonus game.
func (b *BonusCallColor) reset() {
	if b != nil {
		b.deck.Release()
		b.deck = nil
	}
}

// BonusCallColorData represents the details of a bonus "call color" game.
type BonusCallColorData struct {
	cardOwned bool
	Choice    cards.Color
	Cut       int
	Card      *cards.Card
	pool.Object
}

var bonusCallColorDataPool = pool.NewProducer(func() (pool.Objecter, func()) {
	b := &BonusCallColorData{}
	return b, b.reset
})

// reset clears the bonus data.
func (d *BonusCallColorData) reset() {
	if d != nil {
		if d.cardOwned {
			d.Card.Release()
			d.cardOwned = false
		}
		d.Choice = cards.Color(0)
		d.Cut = 0
		d.Card = nil
	}
}

// EncodeFields implements the zjson encoder interface.
func (d *BonusCallColorData) EncodeFields(enc *zjson.Encoder) {
	enc.Uint8Field("choice", uint8(d.Choice))
	enc.IntField("cut", d.Cut)
	enc.ObjectField("card", d.Card)
}

// Encode2 implements the PoolRCZ.Encode2 interface.
func (d *BonusCallColorData) Encode2(enc *zjson.Encoder) {
	d.EncodeFields(enc)
}

// DecodeField implements the zjson decoder interface.
func (d *BonusCallColorData) DecodeField(dec *zjson.Decoder, key []byte) error {
	var ok bool
	var i uint8

	if string(key) == "choice" {
		if i, ok = dec.Uint8(); ok {
			d.Choice = cards.Color(i)
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

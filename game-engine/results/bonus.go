package results

import (
	"fmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// InstantBonusKind represents the kind of instant bonus.
type InstantBonusKind uint8

const (
	// InstantBonusTeaser represent an instant bonus teaser. E.g. a dummy bonus.
	InstantBonusTeaser InstantBonusKind = iota + 1
	// InstantBonusChoice represents an instant bonus which requires a player choice.
	InstantBonusChoice
)

// AcquireInstantBonusTeaser instantiates an instant bonus teaser.
func AcquireInstantBonusTeaser() *InstantBonus {
	b := instantBonusPool.Acquire().(*InstantBonus)
	b.kind = InstantBonusTeaser
	return b
}

// AcquireInstantBonusChoice instantiates an instant bonus player choice.
func AcquireInstantBonusChoice(choice string, options ...string) *InstantBonus {
	b := instantBonusPool.Acquire().(*InstantBonus)
	b.kind = InstantBonusChoice
	b.choice = choice
	b.options = options
	return b
}

// Kind returns the bonus kind.
func (b *InstantBonus) Kind() InstantBonusKind {
	return b.kind
}

// Choice returns the bonus player choice.
func (b *InstantBonus) Choice() string {
	return b.choice
}

// Options returns the bonus player choice options.
func (b *InstantBonus) Options() []string {
	return b.options
}

// String implements the Stringer interface.
func (b *InstantBonus) String() string {
	switch b.kind {
	case InstantBonusTeaser:
		return fmt.Sprintf("Teaser")
	case InstantBonusChoice:
		return fmt.Sprintf("Choice - %s: %v", b.choice, b.options)
	default:
		return fmt.Sprintf("<Unknown>")
	}
}

// EncodeFields implements the zjson.EncodeFields interface.
func (b *InstantBonus) EncodeFields(enc *zjson.Encoder) {
	b.encode(enc, true)
}

// Encode2 implements the Objecter2 interface so the object can be used as result data.
func (b *InstantBonus) Encode2(enc *zjson.Encoder) {
	b.encode(enc, false)
}

func (b *InstantBonus) encode(enc *zjson.Encoder, withLog bool) {
	enc.Uint8Field("kind", uint8(b.kind))
	enc.StringFieldOpt("choice", b.choice)

	if len(b.options) > 0 {
		enc.StartArrayField("options")
		for ix := range b.options {
			enc.String(b.options[ix])
		}
		enc.EndArray()
	}

	b.PlayerChoices.EncodeChoices(enc)
	if withLog {
		b.PrngLog.EncodeEventLog(enc)
	}
}

// DecodeField implements the zjson.DecodeField interface.
func (b *InstantBonus) DecodeField(dec *zjson.Decoder, key []byte) error {
	var ok bool
	if string(key) == "kind" {
		var i uint8
		if i, ok = dec.Uint8(); ok {
			b.kind = InstantBonusKind(i)
		}
	} else if string(key) == "choice" {
		var data []byte
		var escaped bool
		if data, escaped, ok = dec.String(); ok {
			if escaped {
				b.choice = string(dec.Unescaped(data))
			} else {
				b.choice = string(data)
			}
		}
	} else if string(key) == "options" {
		ok = dec.Array(b.decodeOption)
	} else if string(key) == "playerChoices" {
		ok = b.PlayerChoices.DecodeChoices(dec)
	} else if string(key) == "events" {
		ok = dec.Array(b.PrngLog.DecodeEventLog)
	} else if string(key) == "rngIn" {
		ok = dec.Array(b.PrngLog.DecodeRngIn)
	} else if string(key) == "rngOut" {
		ok = dec.Array(b.PrngLog.DecodeRngOut)
	} else {
		return fmt.Errorf("BonusWheelResult: invalid field [%s]", key)
	}

	if ok {
		return nil
	}
	return dec.Error()
}

func (b *InstantBonus) decodeOption(dec *zjson.Decoder) error {
	if data, escaped, ok := dec.String(); ok {
		if escaped {
			b.options = append(b.options, string(dec.Unescaped(data)))
		} else {
			b.options = append(b.options, string(data))
		}
		return nil
	}
	return dec.Error()
}

// InstantBonus holds the details of an instant bonus. E.g. a bonus that is presented before the actual game.
type InstantBonus struct {
	kind    InstantBonusKind
	choice  string
	options []string
	PlayerChoices
	PrngLog
	pool.Object
}

// instantBonusPool is the memory pool for instant bonus objects.
var instantBonusPool = pool.NewProducer(func() (pool.Objecter, func()) {
	b := &InstantBonus{}
	b.PrngLog.Initialize()
	return b, b.reset
})

// reset clears the instant bonus.
func (b *InstantBonus) reset() {
	if b != nil {
		b.kind = 0
		b.choice = ""
		b.options = nil
		b.PlayerChoices.Reset()
		b.PrngLog.Reset()
	}
}

// AcquireBonusSelectorChoice instantiates a bonus selector result.
func AcquireBonusSelectorChoice(player uint8, chosen utils.Index, results ...utils.Index) *BonusSelector {
	b := bonusSelectorPool.Acquire().(*BonusSelector)
	b.playerChoice = player
	b.chosen = chosen
	for ix := range results {
		b.results = append(b.results, results[ix])
	}
	return b
}

// PlayerChoice returns the player choice.
func (b *BonusSelector) PlayerChoice() uint8 {
	return b.playerChoice
}

// Results returns the selected results.
func (b *BonusSelector) Results() utils.Indexes {
	return b.results
}

// Chosen returns the chosen result.
func (b *BonusSelector) Chosen() utils.Index {
	return b.chosen
}

// EncodeFields implements the zjson.EncodeFields interface.
func (b *BonusSelector) EncodeFields(enc *zjson.Encoder) {
	b.encode(enc, true)
}

// Encode2 implements the PoolRCZ.Encode2 interface.
func (b *BonusSelector) Encode2(enc *zjson.Encoder) {
	b.encode(enc, false)
}

func (b *BonusSelector) encode(enc *zjson.Encoder, withLog bool) {
	enc.Uint8Field("playerChoice", b.playerChoice)

	if len(b.results) > 0 {
		enc.StartArrayField("results")
		for ix := range b.results {
			enc.Uint64(uint64(b.results[ix]))
		}
		enc.EndArray()
	}

	enc.Uint16Field("chosen", uint16(b.chosen))

	b.PlayerChoices.EncodeChoices(enc)
	if withLog {
		b.PrngLog.EncodeEventLog(enc)
	}
}

// DecodeField implements the zjson.DecodeField interface.
func (w *BonusSelector) DecodeField(dec *zjson.Decoder, key []byte) error {
	var ok bool
	if string(key) == "result" {
		var i uint8
		if i, ok = dec.Uint8(); ok {
			w.playerChoice = i
		}
	} else if string(key) == "results" {
		ok = dec.Array(w.decodeResult)
	} else if string(key) == "chosen" {
		var i uint16
		if i, ok = dec.Uint16(); ok {
			w.chosen = utils.Index(i)
		}

	} else if string(key) == "playerChoices" {
		ok = w.PlayerChoices.DecodeChoices(dec)
	} else if string(key) == "events" {
		ok = dec.Array(w.PrngLog.DecodeEventLog)
	} else if string(key) == "rngIn" {
		ok = dec.Array(w.PrngLog.DecodeRngIn)
	} else if string(key) == "rngOut" {
		ok = dec.Array(w.PrngLog.DecodeRngOut)
	} else {
		return fmt.Errorf("BonusWheelResult: invalid field [%s]", key)
	}

	if ok {
		return nil
	}
	return dec.Error()
}

func (w *BonusSelector) decodeResult(dec *zjson.Decoder) error {
	if i, ok := dec.Uint16(); ok {
		w.results = append(w.results, utils.Index(i))
		return nil
	}
	return dec.Error()
}

type BonusSelector struct {
	playerChoice uint8
	results      utils.Indexes
	chosen       utils.Index
	PlayerChoices
	PrngLog
	pool.Object
}

// bonusSelectorPool is the memory pool for bonus selector objects.
var bonusSelectorPool = pool.NewProducer(func() (pool.Objecter, func()) {
	b := &BonusSelector{
		results: make(utils.Indexes, 0, 4),
	}
	b.PrngLog.Initialize()
	return b, b.reset
})

// reset clears the bonus selector.
func (b *BonusSelector) reset() {
	if b != nil {
		b.playerChoice = 0
		b.results = b.results[:0]
		b.chosen = 0
		b.PlayerChoices.Reset()
		b.PrngLog.Reset()
	}
}

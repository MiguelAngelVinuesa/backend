package results

import (
	"reflect"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
)

// PlayerChoices can hold the player choices.
type PlayerChoices struct {
	choices map[string]string // player choice(s) at the start of this spin.
}

// Reset resets the player choices.
func (p *PlayerChoices) Reset() {
	clear(p.choices)
}

// Equal is used internally for unit-tests.
func (p *PlayerChoices) Equal(other *PlayerChoices) bool {
	return reflect.DeepEqual(p.choices, other.choices)
}

// SetChoices sets the player choices.
func (p *PlayerChoices) SetChoices(choices map[string]string) {
	p.choices = choices
}

// Choices returns the player choices.
func (p *PlayerChoices) Choices() map[string]string {
	return p.choices
}

// EncodeChoices can be used to encode the player choices.
func (p *PlayerChoices) EncodeChoices(enc *zjson.Encoder) {
	enc.StringMapFieldOpt("playerChoices", p.choices)
}

// DecodeChoices can be used to decode the player choices.
func (p *PlayerChoices) DecodeChoices(dec *zjson.Decoder) bool {
	var ok bool
	p.choices, ok = dec.StringMap(p.choices)
	return ok
}

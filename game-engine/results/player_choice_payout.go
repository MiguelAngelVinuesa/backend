package results

import (
	"fmt"
	"math"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"
)

// AcquirePlayerChoice instantiates a new player choice payout.
func AcquirePlayerChoice(factor float64) Payout {
	p := choiceRewardProducer.Acquire().(*ChoiceReward)
	p.factor = factor
	return p
}

// Kind returns the reward kind.
func (p *ChoiceReward) Kind() PayoutKind {
	return PlayerChoice
}

// Factor returns the win factor for the reward.
func (p *ChoiceReward) Factor() float64 {
	return math.Round(p.factor*100.0) / 100.0
}

// Total returns the total win factor for the reward.
func (p *ChoiceReward) Total() float64 {
	return p.Factor()
}

// EncodeFields implements the zjson.Encoder.EncodeFields interface.
func (p *ChoiceReward) EncodeFields(enc *zjson.Encoder) {
	enc.Uint8Field("kind", uint8(PlayerChoice))
	enc.FloatField("payout", math.Round(p.Factor()*100)/100, 'g', -1)
}

// DecodeField implements the zjson.Decoder.DecodeField interface.
func (p *ChoiceReward) DecodeField(dec *zjson.Decoder, key []byte) error {
	var ok bool

	if string(key) == "kind" {
		_, ok = dec.Uint8()
	} else if string(key) == "payout" {
		p.factor, ok = dec.Float()
	} else {
		return fmt.Errorf("ChoiceReward.DecodeField invalid field: %s", string(key))
	}

	if ok {
		return nil
	}
	return dec.Error()
}

// ChoiceReward contains the details for a player choice payout.
type ChoiceReward struct {
	factor float64
	pool.Object
}

var choiceRewardProducer = pool.NewProducer(func() (pool.Objecter, func()) {
	c := &ChoiceReward{}
	return c, c.reset
})

// reset resets a reward to its initial state.
func (p *ChoiceReward) reset() {
	if p != nil {
		p.factor = 0.0
	}
}

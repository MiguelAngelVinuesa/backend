package wheel

import (
	"fmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/interfaces"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// AcquireBonusWheel instantiates a new bonus wheel game.
func AcquireBonusWheel(prng interfaces.Generator, weights utils.WeightedGenerator) *BonusWheel {
	b := bonusWheelProducer.Acquire().(*BonusWheel)
	b.prng = prng
	b.weights = weights
	return b
}

// RequireParams implements the BonusRunner interface.
func (w *BonusWheel) RequireParams() bool {
	return false
}

// Run implements the BonusRunner interface.
func (w *BonusWheel) Run(_ *results.Result, _ ...interface{}) (int, interfaces.Objecter2) {
	result := w.weights.RandomIndex(w.prng)
	return int(result), AcquireBonusWheelResult(result, w.weights)
}

// BonusWheel is a bonus game which spins a wheel.
type BonusWheel struct {
	weights utils.WeightedGenerator
	prng    interfaces.Generator
	pool.Object
}

// bonusWheelProducer is the memory pool for bonus wheel games.
var bonusWheelProducer = pool.NewProducer(func() (pool.Objecter, func()) {
	w := &BonusWheel{}
	return w, w.reset
})

// reset clears the bonus wheel game.
func (w *BonusWheel) reset() {
	if w != nil {
		w.prng = nil
		w.weights = nil
	}
}

// AcquireBonusWheelResult instantiates a new bonus wheel result.
func AcquireBonusWheelResult(result utils.Index, w utils.WeightedGenerator) *BonusWheelResult {
	b := bonusWheelResultProducer.Acquire().(*BonusWheelResult)
	b.result = result
	if w != nil {
		b.options = w.Options()
	}
	return b
}

// Result returns the generated result.
func (w *BonusWheelResult) Result() utils.Index {
	return w.result
}

// Options returns the possible results.
func (w *BonusWheelResult) Options() utils.Indexes {
	return w.options
}

// EncodeFields implements the zjson.Encoder interface.
func (w *BonusWheelResult) EncodeFields(enc *zjson.Encoder) {
	w.encode(enc, true)
}

// Encode2 implements the Objecter2 interface.
func (w *BonusWheelResult) Encode2(enc *zjson.Encoder) {
	w.encode(enc, false)
}

func (w *BonusWheelResult) encode(enc *zjson.Encoder, withLog bool) {
	enc.Uint16FieldOpt("result", uint16(w.result))

	if len(w.options) > 0 {
		enc.StartArrayField("options")
		for ix := range w.options {
			enc.Uint64(uint64(w.options[ix]))
		}
		enc.EndArray()
	}

	w.PlayerChoices.EncodeChoices(enc)
	if withLog {
		w.PrngLog.EncodeEventLog(enc)
	}
}

// DecodeField implements the zjson.DecodeField interface.
func (w *BonusWheelResult) DecodeField(dec *zjson.Decoder, key []byte) error {
	var ok bool
	if string(key) == "result" {
		var i uint16
		if i, ok = dec.Uint16(); ok {
			w.result = utils.Index(i)
		}
	} else if string(key) == "options" {
		ok = dec.Array(w.decodeOption)
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

func (w *BonusWheelResult) decodeOption(dec *zjson.Decoder) error {
	if i, ok := dec.Uint16(); ok {
		w.options = append(w.options, utils.Index(i))
		return nil
	}
	return dec.Error()
}

// BonusWheelResult represent the result of a bonus wheel game.
type BonusWheelResult struct {
	result  utils.Index
	options utils.Indexes
	results.PlayerChoices
	results.PrngLog
	pool.Object
}

// bonusWheelResultProducer is the memory pool for bonus wheel results.
var bonusWheelResultProducer = pool.NewProducer(func() (pool.Objecter, func()) {
	w := &BonusWheelResult{}
	w.PrngLog.Initialize()
	return w, w.reset
})

// reset clears the bonus wheel result.
func (w *BonusWheelResult) reset() {
	if w != nil {
		w.result = 0
		w.options = nil
		w.PlayerChoices.Reset()
		w.PrngLog.Reset()
	}
}

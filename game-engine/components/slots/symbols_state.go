package slots

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// AcquireSymbolsState instantiates a new state for the symbol set from the memory pool.
func AcquireSymbolsState(symbols *SymbolSet, exclude ...utils.Index) *SymbolsState {
	s := symbolsStateProducer.Acquire().(*SymbolsState)
	if symbols == nil {
		return s
	}

	max := int(symbols.maxID) + 1

	s.flagged = utils.PurgeBools(s.flagged, max)[:max]
	s.valid = utils.PurgeBools(s.valid, max)[:max]

	ex := func(id utils.Index) bool {
		for ix := range exclude {
			if exclude[ix] == id {
				return true
			}
		}
		return false
	}

	sorted := symbols.sorted
	for ix := range sorted {
		symbol := sorted[ix]
		s.valid[ix] = symbol != nil && !ex(symbol.ID())
	}

	return s
}

// DeepCopy creates a deep copy of the state.
func (s *SymbolsState) DeepCopy() *SymbolsState {
	n := symbolsStateProducer.Acquire().(*SymbolsState)
	n.flagged = utils.PurgeBools(n.flagged, cap(s.flagged))[:len(s.flagged)]
	n.valid = utils.PurgeBools(n.valid, cap(s.valid))[:len(s.valid)]
	copy(n.flagged, s.flagged)
	copy(n.valid, s.valid)
	return n
}

// ResetState resets the state.
func (s *SymbolsState) ResetState() {
	s.triggered = false
	clear(s.flagged)
}

// Flagged returns the symbol flags.
func (s *SymbolsState) Flagged() []bool {
	return s.flagged
}

// Valid returns which symbol flags are valid.
func (s *SymbolsState) Valid() []bool {
	return s.valid
}

// IsTriggered returns if the state has triggered an action.
func (s *SymbolsState) IsTriggered() bool {
	return s.triggered
}

// IsFlagged returns whether a symbol is flagged or not.
// This function will panic if the index is out of range.
func (s *SymbolsState) IsFlagged(symbol utils.Index) bool {
	return s.flagged[symbol]
}

// AllFlagged returns whether all symbols are flagged.
func (s *SymbolsState) AllFlagged() bool {
	for ix := 1; ix < len(s.flagged); ix++ {
		if s.valid[ix] && !s.flagged[ix] {
			return false
		}
	}
	return true
}

// SetFlagged sets the flagged state of a symbol.
// This function will panic if the index is out of range.
func (s *SymbolsState) SetFlagged(symbol utils.Index, flagged bool) {
	if s.valid[symbol] {
		s.flagged[symbol] = flagged
	}
}

// IsEmpty implements the zjson encoder interface.
func (s *SymbolsState) IsEmpty() bool {
	return false
}

// EncodeFields implements the zjson encoder interface.
func (s *SymbolsState) EncodeFields(enc *zjson.Encoder) {
	if len(s.flagged) > 0 {
		enc.StartArrayField("flagged")
		for ix := range s.flagged {
			enc.IntBool(s.flagged[ix])
		}
		enc.EndArray()
	}

	if len(s.valid) > 0 {
		enc.StartArrayField("valid")
		for ix := range s.valid {
			enc.IntBool(s.valid[ix])
		}
		enc.EndArray()
	}
}

// Encode2 implements the PoolRCZ.Encode2 interface.
func (s *SymbolsState) Encode2(enc *zjson.Encoder) {
	s.EncodeFields(enc)
}

// DecodeField implements the zjson decoder interface.
func (s *SymbolsState) DecodeField(dec *zjson.Decoder, key []byte) error {
	var ok bool

	if string(key) == "flagged" {
		ok = dec.Array(s.decodeFlagged)
	} else if string(key) == "valid" {
		ok = dec.Array(s.decodeValid)
	}

	if ok {
		return nil
	}
	return dec.Error()
}

func (s *SymbolsState) decodeFlagged(dec *zjson.Decoder) error {
	if b, ok := dec.IntBool(); ok {
		s.flagged = append(s.flagged, b)
		return nil
	}
	return dec.Error()
}

func (s *SymbolsState) decodeValid(dec *zjson.Decoder) error {
	if b, ok := dec.IntBool(); ok {
		s.valid = append(s.valid, b)
		return nil
	}
	return dec.Error()
}

// SymbolsState contains the state for a slot machine symbol set.
// SymbolsState is not safe for concurrent use across multiple go-routines.
// Note that index 0 belongs to the "empty" symbol, so it must always have neutral value.
type SymbolsState struct {
	triggered bool   // indicates if the state triggered an action.
	flagged   []bool // array of flags indexed by symbol id.
	valid     []bool // array of indicators if corresponding flag is valid symbol.
	pool.Object
}

// symbolsStateProducer is the memory pool for symbols states.
// Make sure to initialize all slices appropriately!
var symbolsStateProducer = pool.NewProducer(func() (pool.Objecter, func()) {
	s := &SymbolsState{
		flagged: make([]bool, 0, 16),
		valid:   make([]bool, 0, 16),
	}
	return s, s.reset
})

// reset clears the symbols state.
func (s *SymbolsState) reset() {
	if s != nil {
		clear(s.flagged)
		clear(s.valid)

		s.triggered = false
		s.flagged = s.flagged[:0]
		s.valid = s.valid[:0]
	}
}

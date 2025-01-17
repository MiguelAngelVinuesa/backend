package slots

import (
	"fmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"
)

// AcquireGameState instantiates a new game state from the memory pool with the given parameters.
func AcquireGameState(spin *slots.SpinState, symbols *slots.SymbolsState, bet int64) *GameState {
	gs := gameStatePool.Acquire().(*GameState)
	if spin != nil {
		gs.spin = spin.Clone().(*slots.SpinState)
	}
	if symbols != nil {
		gs.symbols = symbols.Clone().(*slots.SymbolsState)
	}
	gs.roundSeq = 1
	gs.bet = bet
	return gs
}

// AcquireGameStateFromJSON instantiates a new game state from the memory pool using the given json data.
func AcquireGameStateFromJSON(data []byte) (*GameState, error) {
	o, err := gameStatePool.AcquireFromJSON(data)
	if err == nil {
		return o.(*GameState), nil
	}
	o.Release()
	return nil, err
}

// RoundID returns the round id.
func (s *GameState) RoundID() string {
	return s.roundID
}

// Bet returns the bet amount.
func (s *GameState) Bet() int64 {
	return s.bet
}

// RoundSeq returns the round sequence number.
func (s *GameState) RoundSeq() int64 {
	return s.roundSeq
}

// NextOffset returns the offset for round/next calls.
func (s *GameState) NextOffset() int64 {
	return s.nextOffset
}

// SpinState returns the spin state.
func (s *GameState) SpinState() *slots.SpinState {
	return s.spin
}

// SymbolsState returns the symbols state.
func (s *GameState) SymbolsState() *slots.SymbolsState {
	return s.symbols
}

// SetSymbolsState returns the symbols state.
func (s *GameState) SetSymbolsState(state *slots.SymbolsState) {
	if s.symbols != nil {
		s.symbols.Release()
	}
	s.symbols = state.Clone().(*slots.SymbolsState)
}

// SetRoundID sets the current round identifier.
func (s *GameState) SetRoundID(roundID string) {
	s.roundID = roundID
}

// SetRoundSeq sets the current round sequence number.
func (s *GameState) SetRoundSeq(seq int64) {
	s.roundSeq = seq
}

// SetNextOffset sets the offset for retrieving the next spin result from the backing store.
func (s *GameState) SetNextOffset(offset int64) {
	s.nextOffset = offset
}

// EncodeFields implements the zjson encoder interface.
func (s *GameState) EncodeFields(enc *zjson.Encoder) {
	enc.StringFieldOpt("roundID", s.roundID)
	enc.Int64FieldOpt("roundSeq", s.roundSeq)
	enc.Int64FieldOpt("nextOffset", s.nextOffset)
	enc.Int64FieldOpt("bet", s.bet)
	if s.spin != nil {
		enc.ObjectField("spin", s.spin)
	}
	if s.symbols != nil {
		enc.ObjectField("symbols", s.symbols)
	}
}

// Decode unmarshalls the game state from JSON.
func (s *GameState) Decode(dec *zjson.Decoder) error {
	if dec.Object(s) {
		return nil
	}
	return dec.Error()
}

// DecodeField implements the zjson decoder interface.
func (s *GameState) DecodeField(dec *zjson.Decoder, key []byte) error {
	var ok, escaped bool
	var b []byte
	var i int64

	if string(key) == "roundID" {
		if b, escaped, ok = dec.String(); ok {
			if escaped {
				s.roundID = string(dec.Unescaped(b))
			} else {
				s.roundID = string(b)
			}
		}
	} else if string(key) == "roundSeq" {
		if i, ok = dec.Int64(); ok {
			s.roundSeq = i
		}
	} else if string(key) == "nextOffset" {
		if i, ok = dec.Int64(); ok {
			s.nextOffset = i
		}
	} else if string(key) == "bet" {
		if i, ok = dec.Int64(); ok {
			s.bet = i
		}
	} else if string(key) == "spin" {
		s.spin = slots.AcquireSpinState(nil)
		ok = dec.Object(s.spin)
	} else if string(key) == "symbols" {
		s.symbols = slots.AcquireSymbolsState(nil)
		ok = dec.Object(s.symbols)
	} else {
		return fmt.Errorf("GameState.Decode: invalid field '%s'", string(key))
	}

	if ok {
		return nil
	}
	return dec.Error()
}

// GameState represents the game state for a single player slots game session.
// It does not take ownership of the state objects.
// GameState is not safe for use across multiple go-routines.
type GameState struct {
	roundSeq   int64
	nextOffset int64
	bet        int64
	spin       *slots.SpinState
	symbols    *slots.SymbolsState
	roundID    string
	pool.Object
}

// gameStatePool is the memory pool for game states.
var gameStatePool = pool.NewProducer(func() (pool.Objecter, func()) {
	s := &GameState{}
	return s, s.reset
})

// reset clears the game state.
func (s *GameState) reset() {
	if s.spin != nil {
		s.spin.Release()
	}
	if s.symbols != nil {
		s.symbols.Release()
	}

	s.roundSeq = 0
	s.nextOffset = 0
	s.bet = 0
	s.spin = nil
	s.symbols = nil
	s.roundID = ""
}

package slots

import (
	"fmt"
	"sort"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/object"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"
)

// AcquireGamePrefs instantiates a new player state from the memory pool with the given parameters.
func AcquireGamePrefs(bet int64, flags *slots.SymbolsState) *GamePrefs {
	s := gamePrefsPool.Acquire().(*GamePrefs)
	s.symbolsCCB[bet] = flags.Clone().(*slots.SymbolsState)
	return s
}

// AcquireGamePrefsFromJSON instantiates a new player state from the memory pool using the given json data.
func AcquireGamePrefsFromJSON(data []byte) (*GamePrefs, error) {
	o, err := gamePrefsPool.AcquireFromJSON(data)
	if err == nil {
		return o.(*GamePrefs), nil
	}
	o.Release()
	return nil, err
}

// GamePrefs returns all game preferences.
func (s *GamePrefs) GamePrefs() map[string]string {
	return s.game
}

// GamePref returns a game preference.
func (s *GamePrefs) GamePref(key string) string {
	return s.game[key]
}

// SetGamePref updates a game preference.
func (s *GamePrefs) SetGamePref(key, value string) {
	s.game[key] = value
}

// AddStateCCB adds or modifies the state for a specific bet for ChaCha Bomb.
func (s *GamePrefs) AddStateCCB(bet int64, flags *slots.SymbolsState) {
	s.symbolsCCB[bet] = flags.Clone().(*slots.SymbolsState)
}

// GetStateCCB adds or modifies the state for a specific bet for ChaCha Bomb.
func (s *GamePrefs) GetStateCCB(bet int64) *slots.SymbolsState {
	return s.symbolsCCB[bet]
}

// EncodeFields implements the zjson.Encoder.EncodeFields interface.
func (s *GamePrefs) EncodeFields(enc *zjson.Encoder) {
	enc.StringMapFieldOpt("prefs", s.game)

	if s.symbolsCCB != nil {
		bets := make([]int64, 0, 100)
		for bet := range s.symbolsCCB {
			bets = append(bets, bet)
		}
		sort.Slice(bets, func(i, j int) bool { return bets[i] < bets[j] })

		enc.StartArrayField("ccbBets")
		for _, bet := range bets {
			enc.Int64(bet)
		}
		enc.EndArray()

		enc.StartArrayField("ccbStates")
		for _, bet := range bets {
			enc.Object(s.symbolsCCB[bet])
		}
		enc.EndArray()
	}
}

// DecodeField implements the zjson.Decoder.DecodeField interface.
func (s *GamePrefs) DecodeField(dec *zjson.Decoder, key []byte) error {
	var ok bool

	if string(key) == "prefs" {
		s.game, ok = dec.StringMap(s.game)
	} else if string(key) == "ccbBets" {
		dec.Array(s.decodeCCBbet)
	} else if string(key) == "ccbStates" {
		dec.Array(s.decodeCCBstate)
	}

	if ok {
		return nil
	}
	return dec.Error()
}

func (s *GamePrefs) decodeCCBbet(dec *zjson.Decoder) error {
	if i, ok := dec.Int64(); ok {
		if s.bets == nil {
			s.bets = stateBetsPool.Acquire().(*object.Int64sManager)
		}
		s.bets.Append(i)
		return nil
	}

	return dec.Error()
}

func (s *GamePrefs) decodeCCBstate(dec *zjson.Decoder) error {
	if s.symbolsCCB == nil {
		s.symbolsCCB = make(map[int64]*slots.SymbolsState)
	}

	ix := len(s.symbolsCCB)
	if ix >= len(s.bets.Items) {
		return fmt.Errorf("unmatched CCB states")
	}

	flags := slots.AcquireSymbolsState(nil)
	if ok := dec.Object(flags); ok {
		bet := s.bets.Items[ix]
		s.symbolsCCB[bet] = flags
		return nil
	}

	return dec.Error()
}

// GamePrefs represents the player state for a single player slots game session.
// It does not take ownership of given state objects; it simply clones any input.
// GamePrefs is not safe for use across multiple go-routines.
type GamePrefs struct {
	bets       *object.Int64sManager         // temporarily used during decoding
	game       map[string]string             // prefs on game level (e.g. bet, audio)
	symbolsCCB map[int64]*slots.SymbolsState // CCB game states per bet
	pool.Object
}

var stateBetsPool = object.NewInt64sProducer(8, 32, true)
var gamePrefsPool = pool.NewProducer(func() (pool.Objecter, func()) {
	p := &GamePrefs{
		bets:       stateBetsPool.Acquire().(*object.Int64sManager),
		game:       make(map[string]string, 8),
		symbolsCCB: make(map[int64]*slots.SymbolsState, 8),
	}
	return p, p.reset
})

// reset clears the game preferences.
func (s *GamePrefs) reset() {
	for k := range s.game {
		delete(s.game, k)
	}

	if s.symbolsCCB != nil {
		for bet := range s.symbolsCCB {
			s.symbolsCCB[bet].Release()
			delete(s.symbolsCCB, bet)
		}
	}

	if s.bets != nil {
		s.bets.Release()
		s.bets = nil
	}
}

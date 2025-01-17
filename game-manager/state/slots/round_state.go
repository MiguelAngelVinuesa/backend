package slots

import (
	"fmt"
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/object"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"
)

// AcquireRoundState instantiates a new round state from the memory pool.
func AcquireRoundState(sessionID, roundID string, spins int) *RoundState {
	s := roundStatePool.Acquire().(*RoundState)
	s.sessionID = sessionID
	s.roundID = roundID
	s.spins = spins

	if s.spins == 1 {
		// we don't get the "finish" call , so mark the single spin as played (1 second into the future)
		s.SpinPlayed(1, time.Now().UTC().Add(time.Second))
	}

	return s
}

// AcquireRoundStateFromJSON instantiates a new round state from the memory pool using the JSON input.
func AcquireRoundStateFromJSON(sessionID, roundID string, data []byte) (*RoundState, error) {
	o, err := roundStatePool.AcquireFromJSON(data)
	if err != nil {
		return nil, err
	}

	s := o.(*RoundState)
	s.sessionID = sessionID
	s.roundID = roundID
	return s, nil
}

// Spins returns the number of spins in the round.
func (s *RoundState) Spins() int {
	return s.spins
}

// PlayedFull return whether all spins of the round have been marked as played.
func (s *RoundState) PlayedFull() bool {
	return s.playedFull
}

// ResumePlay restores the PlayedFull flag for DSF.
func (s *RoundState) ResumePlay(spins int) {
	s.spins = spins
	if spins == 2 {
		s.playedFull = true
		s.replayedFull = false
		s.SpinPlayed(2, time.Now().UTC().Add(time.Second))
	} else {
		s.playedFull = false
		s.replayedFull = false
	}
}

// ReplayedFull return whether all spins of the round have been marked as re-played.
func (s *RoundState) ReplayedFull() bool {
	return s.replayedFull
}

// SpinPlayed logs the display time of a spin.
func (s *RoundState) SpinPlayed(spinSeq int, ts time.Time) {
	if s.played == nil {
		s.played = roundTimesPool.Acquire().(*object.TimesManager)
		for len(s.played.Items) < s.spins {
			s.played.Append(time.Time{})
		}
	}

	for len(s.played.Items) < spinSeq {
		s.played.Append(time.Time{})
	}

	s.played.Items[spinSeq-1] = ts.UnixMilli()
	s.playSeq = spinSeq
	s.playedFull = spinSeq >= s.spins
}

// SpinReplayed logs the replay time of a spin.
func (s *RoundState) SpinReplayed(spinSeq int, ts time.Time) {
	if s.replayed == nil {
		s.replayed = roundTimesPool.Acquire().(*object.TimesManager)
		for len(s.replayed.Items) < s.spins {
			s.replayed.Append(time.Time{})
		}
	}

	for len(s.replayed.Items) < spinSeq {
		s.replayed.Append(time.Time{})
	}

	s.replayed.Items[spinSeq-1] = ts.UnixMilli()
	s.replaySeq = spinSeq
	s.replayedFull = spinSeq >= s.spins
}

// EncodeFields implements the zjson.Encoder interface.
func (s *RoundState) EncodeFields(enc *zjson.Encoder) {
	enc.IntBoolFieldOpt("playedFull", s.playedFull)
	enc.IntBoolFieldOpt("replayedFull", s.replayedFull)
	enc.IntFieldOpt("spins", s.spins)
	enc.IntFieldOpt("playSeq", s.playSeq)
	enc.IntFieldOpt("replaySeq", s.replaySeq)
	enc.ArrayFieldOpt("flags", s.flags)
	enc.ArrayFieldOpt("played", s.played)
	enc.ArrayFieldOpt("replayed", s.replayed)
}

// DecodeField decodes a round state field from the decoder.
func (s *RoundState) DecodeField(dec *zjson.Decoder, key []byte) error {
	var ok bool

	if string(key) == "playedFull" {
		s.playedFull, ok = dec.IntBool()
	} else if string(key) == "replayedFull" {
		s.replayedFull, ok = dec.IntBool()
	} else if string(key) == "spins" {
		s.spins, ok = dec.Int()
	} else if string(key) == "playSeq" {
		s.playSeq, ok = dec.Int()
	} else if string(key) == "replaySeq" {
		s.replaySeq, ok = dec.Int()
	} else if string(key) == "flags" {
		if s.flags == nil {
			s.flags = roundFlagsPool.Acquire().(*object.IntsManager)
		}
		return s.flags.Decode(dec)
	} else if string(key) == "played" {
		if s.played == nil {
			s.played = roundTimesPool.Acquire().(*object.TimesManager)
		}
		return s.played.Decode(dec)
	} else if string(key) == "replayed" {
		if s.replayed == nil {
			s.replayed = roundTimesPool.Acquire().(*object.TimesManager)
		}
		return s.replayed.Decode(dec)
	} else {
		return fmt.Errorf("RoundState.DecodeField: invalid field '%s'", string(key))
	}

	if ok {
		return nil
	}
	return dec.Error()
}

var (
	roundFlagsPool = object.NewIntsProducer(16, 32, true)
	roundTimesPool = object.NewTimesProducer(16, 64, true)
)

// RoundState contains the details of a round state.
type RoundState struct {
	playedFull   bool                 // indicates inital play has completed.
	replayedFull bool                 // indicates last replay has completed.
	spins        int                  // total number of spins in the round.
	playSeq      int                  // current sequence number for initial play.
	replaySeq    int                  // current sequence number for replay.
	flags        *object.IntsManager  // flags for this round.
	played       *object.TimesManager // last display time for each spin (initial play).
	replayed     *object.TimesManager // last display time for each spin (replays).
	sessionID    string               // unique session id the round is part of.
	roundID      string               // unique id of the round.
	pool.Object                       // base Object to satisfy Objecter interface.
}

// roundStatePool is the memory pool for a round state.
var roundStatePool = pool.NewProducer(func() (pool.Objecter, func()) {
	s := &RoundState{}
	return s, s.reset
})

// reset clears the round state.
func (s *RoundState) reset() {
	if s != nil {
		s.playedFull = false
		s.replayedFull = false
		s.spins = 0
		s.playSeq = 0
		s.replaySeq = 0
		s.sessionID = ""
		s.roundID = ""

		if s.flags != nil {
			s.flags.Release()
			s.flags = nil
		}
		if s.played != nil {
			s.played.Release()
			s.played = nil
		}
		if s.replayed != nil {
			s.replayed.Release()
			s.replayed = nil
		}
	}
}

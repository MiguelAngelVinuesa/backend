package slots

import (
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git/consts"
)

// AcquireSessionState instantiates a new session state from the memory pool.
// It makes a deep copy of the input, so it's safe to call Release() on the supplied round.
func AcquireSessionState(round *Round) *SessionState {
	s := sessionStatePool.Acquire().(*SessionState)
	s.round = round.Clone().(*Round)
	s.balance = 1000000 + s.round.TotalWin() - s.round.Bet()
	s.created = time.Now()
	s.expires = s.created.Add(consts.DefaultSessionExpire)
	return s
}

// Balance returns the current balance for the session.
func (s *SessionState) Balance() int64 {
	s.Touch()
	return s.balance
}

// RoundResults returns the round results from the last round in the session.
func (s *SessionState) RoundResults() RoundResults {
	s.Touch()
	return s.round.roundResults
}

// Expired returns true if the session state has expired.
func (s *SessionState) Expired() bool {
	return s.expires.Before(time.Now())
}

// Touch sets a new expiration timestamp.
func (s *SessionState) Touch() {
	s.expires = time.Now().Add(consts.DefaultSessionExpire)
}

// SetRound stores a deep copy of the given round for the session and updates the balance.
func (s *SessionState) SetRound(round *Round) {
	s.Touch()
	if s.round != nil {
		s.round.Release()
	}
	s.round = round.Clone().(*Round)
	if len(s.round.roundResults) == 0 {
		panic("invalid round: no results")
	}
	if s.round.roundResults[0].SpinData != nil && s.round.roundResults[0].SpinData.Kind() != slots.SecondSpin {
		s.balance -= s.round.TotalBet()
	}
	s.balance += s.round.TotalWin()
}

// GetState returns a deep copy of the last game state for the session.
func (s *SessionState) GetState() *GameState {
	s.Touch()
	return s.round.GameState()
}

// SetState stores a deep copy of the last game state for the session.
func (s *SessionState) SetState(state *GameState) {
	s.Touch()
	s.round.SetGameState(state)
}

// SessionState keeps track of the session balance and last round played.
// It is not safe to be used across multiple go-routines.
type SessionState struct {
	balance int64
	round   *Round
	created time.Time
	expires time.Time
	pool.Object
}

// sessionStatePool is the memory pool for session states.
var sessionStatePool = pool.NewProducer(func() (pool.Objecter, func()) {
	s := &SessionState{}
	return s, s.reset
})

// reset clears the session state.
func (s *SessionState) reset() {
	if s != nil {
		if s.round != nil {
			s.round.Release()
			s.round = nil
		}

		s.balance = 0
		s.created = time.Time{}
		s.expires = time.Time{}
	}
}

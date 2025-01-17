package slots

import (
	"sync"
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git/consts"
	state "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git/state/slots"
)

// memory represents a game bet&state manager utilizing local memory as the backing store.
// It should be created once during app initialization, and only for development sessions.
type memory struct {
	mu       sync.RWMutex
	sessions map[string]*state.SessionState
}

// NewMemory instantiates a new game round manager using local memory.
func NewMemory() state.RoundManager {
	m := &memory{sessions: make(map[string]*state.SessionState, 256)}

	go func(m *memory) {
		t := time.NewTicker(time.Minute)
		for {
			select {
			case <-t.C:
				m.checkExpired()
			}
		}
	}(m)

	return m
}

// PostRound implements the RoundManager interface.
func (m *memory) PostRound(round *state.Round, _ bool) (string, int64, error) {
	return m.newRound(round)
}

// PostInitRound implements the RoundManager interface.
func (m *memory) PostInitRound(round *state.Round, _ bool) (string, int64, error) {
	return m.newRound(round)
}

// PostCompleteRound implements the RoundManager interface.
func (m *memory) PostCompleteRound(round *state.Round, _ *state.RoundState, _ bool) (string, int64, error) {
	return m.newRound(round)
}

func (m *memory) newRound(round *state.Round) (string, int64, error) {
	sessionID := round.SessionID()

	m.mu.Lock()

	s := m.sessions[sessionID]
	if s == nil {
		s = state.AcquireSessionState(round)
	} else {
		s.SetRound(round)
	}
	m.sessions[sessionID] = s

	m.mu.Unlock()
	return "1", s.Balance(), nil
}

// PostRoundNext implements the RoundManager interface.
// RoundID is ignored here, as we just assume it's correct.
func (m *memory) PostRoundNext(sessionID, _ string, _ *state.RoundState, spinSeq int) (*state.RoundResult, int64, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	s := m.sessions[sessionID]
	if s == nil {
		return nil, 0, consts.ErrSessionNotFound
	}

	spinSeq--
	results := s.RoundResults()
	if spinSeq < 0 || spinSeq >= len(results) {
		return nil, 0, consts.ErrSpinSeqNotFound
	}

	return results[spinSeq].Clone().(*state.RoundResult), s.Balance(), nil
}

// GetRoundState implements the RoundManager interface.
func (m *memory) GetRoundState(_, _ string) (*state.RoundState, error) {
	return nil, nil
}

// PutGameState implements the RoundManager interface.
func (m *memory) PutGameState(sessionID string, state *state.GameState) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	s := m.sessions[sessionID]
	if s == nil {
		return consts.ErrSessionNotFound
	}

	s.SetState(state)
	return nil
}

// GetGameState implements the RoundManager interface.
func (m *memory) GetGameState(sessionID string) (*state.GameState, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	s := m.sessions[sessionID]
	if s == nil {
		return nil, consts.ErrSessionNotFound
	}

	gs := s.GetState()
	if gs == nil {
		return nil, consts.ErrGameStateNotFound
	}

	return gs.Clone().(*state.GameState), nil
}

// GetGamePrefs implements the RoundManager interface.
func (m *memory) GetGamePrefs(_ string) (string, string, *state.GamePrefs, error) {
	return "", "", nil, nil
}

// PutGamePrefs implements the RoundManager interface.
func (m *memory) PutGamePrefs(_ string, _ *state.GamePrefs) error {
	return nil
}

// GetPlayerPrefs implements the RoundManager interface.
func (m *memory) GetPlayerPrefs(_ string) (map[string]string, error) {
	return nil, nil
}

// PutPlayerPrefs implements the RoundManager interface.
func (m *memory) PutPlayerPrefs(_ string, _ map[string]string) error {
	return nil
}

// checkExpired removes expired sessions from memory.
func (m *memory) checkExpired() {
	var keys []string
	if l := len(m.sessions); l > 4096 {
		keys = make([]string, len(m.sessions))
	} else {
		keys = make([]string, 4096)[:l]
	}

	m.mu.RLock()
	var ix int
	for k := range m.sessions {
		keys[ix] = k
		ix++
	}
	m.mu.RUnlock()

	for _, key := range keys {
		m.mu.Lock()
		if s := m.sessions[key]; s != nil {
			if s.Expired() {
				delete(m.sessions, key)
				s.Release()
			}
		}
		m.mu.Unlock()
	}
}

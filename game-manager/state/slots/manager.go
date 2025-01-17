package slots

import (
	"fmt"
)

// RoundManager is the interface for bet validation and state management for slots games.
type RoundManager interface {
	PostRound(r *Round, debug bool) (string, int64, error)
	PostInitRound(r *Round, debug bool) (string, int64, error)
	PostCompleteRound(r *Round, state *RoundState, debug bool) (string, int64, error)
	PostRoundNext(sessionID, roundID string, roundState *RoundState, spinSeq int) (*RoundResult, int64, error)

	GetRoundState(sessionID, roundID string) (*RoundState, error)

	PutGameState(sessionID string, state *GameState) error
	GetGameState(sessionID string) (*GameState, error)

	GetGamePrefs(sessionID string) (string, string, *GamePrefs, error)
	PutGamePrefs(sessionID string, state *GamePrefs) error

	GetPlayerPrefs(sessionID string) (map[string]string, error)
	PutPlayerPrefs(sessionID string, state map[string]string) error
}

type APIerror struct {
	Err     error
	Status  int
	Code    int
	Level   string
	Message string
}

func (e *APIerror) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return fmt.Sprintf("http %d - error %d [%s]: %s", e.Status, e.Code, e.Level, e.Message)
}

package metrics

import (
	"os"
	"sync"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts"
)

var (
	sessionMetrics bool = true
	mu             sync.Mutex
	sessions       map[string]struct{}
)

func MarkSession(sessionID string) {
	if !sessionMetrics {
		return
	}
	mu.Lock()
	sessions[sessionID] = struct{}{}
	mu.Unlock()
}

func ResetSessions() int {
	if !sessionMetrics {
		return 0
	}
	mu.Lock()
	i := len(sessions)
	sessions = make(map[string]struct{}, 256)
	mu.Unlock()
	return i
}

func init() {
	if s := os.Getenv(consts.EnvSessionMetrics); s != "" {
		sessionMetrics = s == "1"
	}
	if sessionMetrics {
		sessions = make(map[string]struct{}, 2048)
	}
}

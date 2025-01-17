package consts

import (
	"fmt"
	"time"
)

const (
	DefaultSessionExpire = 2 * time.Hour
)

var (
	ErrNotValidated      = fmt.Errorf("not validated")
	ErrSessionNotFound   = fmt.Errorf("session not found")
	ErrSpinSeqNotFound   = fmt.Errorf("spin sequence not found")
	ErrGameStateNotFound = fmt.Errorf("game state not found")
)

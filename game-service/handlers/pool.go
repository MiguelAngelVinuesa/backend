package handlers

import (
	"sync"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/api/models"
)

var (
	roundStartRequestPool  = sync.Pool{New: func() any { return &models.RoundStartRequest{} }}
	roundPaidRequestPool   = sync.Pool{New: func() any { return &models.RoundPaidRequest{} }}
	roundSecondRequestPool = sync.Pool{New: func() any { return &models.RoundSecondRequest{} }}
	roundResumeRequestPool = sync.Pool{New: func() any { return &models.RoundResumeRequest{} }}
	// SUPERVISED-BUILD-REMOVE-START
	roundDebugRequestPool = sync.Pool{New: func() any { return &models.RoundDebugRequest{} }}
	// SUPERVISED-BUILD-REMOVE-END
	roundNextRequestPool   = sync.Pool{New: func() any { return &models.RoundNextRequest{} }}
	roundFinishRequestPool = sync.Pool{New: func() any { return &models.RoundFinishRequest{} }}
)

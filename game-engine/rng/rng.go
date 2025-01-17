package rng

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/interfaces"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/sharedlib"
)

var AcquireRNG func() interfaces.Generator

func init() {
	AcquireRNG = sharedlib.AcquireRNG
}

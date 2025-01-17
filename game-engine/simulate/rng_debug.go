package simulate

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/prng.git/rng"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/interfaces"
)

// NewRNG instantiate a new RNG using the embedded prng package.
func NewRNG() interfaces.Generator {
	return rng.NewRNG()
}

package owl

import (
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/interfaces"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/simulate"
)

type RandomPlayer struct {
	prng interfaces.Generator
}

func NewRandomPlayer() *RandomPlayer {
	return &RandomPlayer{prng: simulate.NewRNG()}
}

func (p *RandomPlayer) Choices() map[string]string {
	switch p.prng.IntN(3) {
	case 0:
		return map[string]string{"selection": "left"}
	case 1:
		return map[string]string{"selection": "middle"}
	default:
		return map[string]string{"selection": "right"}
	}
}

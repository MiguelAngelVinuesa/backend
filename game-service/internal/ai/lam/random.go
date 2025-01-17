package lam

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
	if p.prng.IntN(10000) < 5000 {
		return map[string]string{"wing": "south"}
	}
	return map[string]string{"wing": "north"}
}

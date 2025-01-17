package ai

import "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"

type Player interface {
	RoundResult(bet, win int64, freeGames int)
}

type BetChanger interface {
	SelectBet() int64
}

type StickyChooser interface {
	SelectSticky(grid utils.Indexes, auto utils.Index, flags []bool) utils.Index
}

type ChoiceMaker interface {
	Choices() map[string]string
}

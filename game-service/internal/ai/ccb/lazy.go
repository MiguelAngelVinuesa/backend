package ccb

import "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"

type LazyPlayer struct{}

func (p *LazyPlayer) SelectSticky(_ utils.Indexes, auto utils.Index, _ []bool) utils.Index {
	return auto
}

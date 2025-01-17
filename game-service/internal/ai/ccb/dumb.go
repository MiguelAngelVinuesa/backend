package ccb

import "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"

type DumbPlayer struct{}

func (p *DumbPlayer) SelectSticky(grid utils.Indexes, _ utils.Index, flags []bool) utils.Index {
	max := utils.Index(len(flags) - 1)
	for id := utils.Index(1); id <= max; id++ {
		var count int
		for ix := range grid {
			if grid[ix] == id {
				count++
			}
		}
		if count == 1 && flags[id] {
			return id
		}
	}

	for id := utils.Index(1); id <= max; id++ {
		var count int
		for ix := range grid {
			if grid[ix] == id {
				count++
			}
		}
		if count == 1 {
			return id
		}
	}

	return 1
}

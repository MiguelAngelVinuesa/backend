package ccb

import "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"

type SmartPlayerHigh struct{}

func (p *SmartPlayerHigh) SelectSticky(grid utils.Indexes, auto utils.Index, flags []bool) utils.Index {
	max := utils.Index(len(flags) - 1)

	for min := 4; min > 0; min-- {
		for id := max; id > 0; id-- {
			var count int
			for ix := range grid {
				if grid[ix] == id {
					count++
				}
			}
			if count > min && flags[id] {
				return id
			}
		}
	}

	return auto
}

type SmartPlayerLow struct{}

func (p *SmartPlayerLow) SelectSticky(grid utils.Indexes, auto utils.Index, flags []bool) utils.Index {
	max := utils.Index(len(flags) - 1)

	for min := 4; min > 0; min-- {
		for id := utils.Index(1); id <= max; id++ {
			var count int
			for ix := range grid {
				if grid[ix] == id {
					count++
				}
			}
			if count > min && flags[id] {
				return id
			}
		}
	}

	return auto
}

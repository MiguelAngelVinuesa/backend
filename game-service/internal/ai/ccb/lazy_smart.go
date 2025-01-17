package ccb

import (
	"math/rand"
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/interfaces"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/simulate"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/ai"
)

type LazySmartPlayerHigh struct {
	prng interfaces.Generator
}

func NewLazySmartPlayerHigh() ai.StickyChooser {
	return &LazySmartPlayerHigh{prng: simulate.NewRNG()}
}

func (p *LazySmartPlayerHigh) SelectSticky(grid utils.Indexes, auto utils.Index, flags []bool) utils.Index {
	if p.prng.IntN(2) == 1 {
		return auto
	}

	max := utils.Index(len(flags) - 1)
	for min := 4; min > 0; min-- {
		for id := max; id > 0; id-- {
			var count int
			for ix := range grid {
				if grid[ix] == id {
					count++
				}
			}
			if count > min && !flags[id] {
				return id
			}
		}
	}

	return auto
}

type LazySmartPlayerLow struct {
	prng interfaces.Generator
}

func NewLazySmartPlayerLow() ai.StickyChooser {
	return &LazySmartPlayerLow{prng: simulate.NewRNG()}
}

func (p *LazySmartPlayerLow) SelectSticky(grid utils.Indexes, auto utils.Index, flags []bool) utils.Index {
	if p.prng.IntN(2) == 1 {
		return auto
	}

	max := utils.Index(len(flags) - 1)
	for min := 4; min > 0; min-- {
		for id := utils.Index(1); id <= max; id++ {
			var count int
			for ix := range grid {
				if grid[ix] == id {
					count++
				}
			}
			if count > min && !flags[id] {
				return id
			}
		}
	}

	return auto
}

func init() {
	rand.Seed(time.Now().UnixMicro())
}

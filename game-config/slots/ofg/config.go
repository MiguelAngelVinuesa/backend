package ofg

import (
	comp "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	game "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/slots"
	util "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

const (
	symbolCount   = 12
	reels         = 5
	rows          = 3
	direction     = comp.PayLTR
	highestPayout = true
	bonusBuyCost  = 150
	maxPayout     = 10000.0

	id01 = 1
	id02 = 2
	id03 = 3
	id04 = 4
	id05 = 5
	id06 = 6
	id07 = 7
	id08 = 8
	id09 = 9
	id10 = 10
	id11 = 11
	id12 = 12
	id13 = 13
	id14 = 14
	id15 = 15

	scatter = id11
	wild    = id12

	bonusBuyID        = 1
	initMultiplierID  = 2
	wilds92aID        = 10
	wilds92bID        = 11
	wilds92cID        = 12
	wilds92dID        = 13
	wilds92eID        = 14
	scatters92aID     = 15
	scatters92bID     = 16
	scatters92cID     = 17
	scatters92dID     = 18
	scatters92eID     = 19
	wilds94aID        = 20
	wilds94bID        = 21
	wilds94cID        = 22
	wilds94dID        = 23
	wilds94eID        = 24
	scatters94aID     = 25
	scatters94bID     = 26
	scatters94cID     = 27
	scatters94dID     = 28
	scatters94eID     = 29
	wilds96aID        = 30
	wilds96bID        = 31
	wilds96cID        = 32
	wilds96dID        = 33
	wilds96eID        = 34
	scatters96aID     = 35
	scatters96bID     = 36
	scatters96cID     = 37
	scatters96dID     = 38
	scatters96eID     = 39
	multiplierID      = 50
	winlinesID        = 60
	scatterPayouts2ID = 75
	scatterPayouts3ID = 76
	scatterPayouts4ID = 77
	scatterPayouts5ID = 78
	freeSpinsFirstID  = 80
	freeSpinsFreeID   = 81
	freeSpinsFlagID   = 91

	flagBonusBuy      = 0
	flagFreeSpins     = 1
	flagFirstReelSet1 = 2
	flagFirstReelSet2 = 3
	flagFirstReelSets = 4
	flagFreeReelSet   = 5
	flagScattersLevel = 6

	freeSpinsAwardFirst = 5
	freeSpinsAwardFree  = 5
)

var (
	weights92 = [symbolCount]comp.SymbolOption{
		comp.WithWeights(40, 40, 40, 40, 40),
		comp.WithWeights(40, 40, 40, 40, 40),
		comp.WithWeights(40, 40, 40, 40, 40),
		comp.WithWeights(40, 40, 40, 40, 40),
		comp.WithWeights(40, 40, 40, 40, 40),
		comp.WithWeights(40, 40, 40, 40, 40),
		comp.WithWeights(40, 40, 40, 40, 40),
		comp.WithWeights(40, 40, 40, 40, 40),
		comp.WithWeights(40, 40, 40, 40, 40),
		comp.WithWeights(40, 40, 40, 40, 40),
		comp.WithWeights(0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0),
	}

	wildWeights92a = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5, 0}, []float64{820, 200, 32, 6, 1, 2100})
	wildWeights92b = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5, 0}, []float64{3800, 490, 50, 9, 1, 2400})
	wildWeights92c = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5, 0}, []float64{3600, 430, 45, 8, 1, 2700})
	wildWeights92d = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5, 0}, []float64{3400, 380, 40, 7, 1, 3100})
	wildWeights92e = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5, 0}, []float64{3200, 340, 35, 6, 1, 3900})

	scatterWeights92a = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5, 0}, []float64{4800, 900, 180, 11, 1, 24000})
	scatterWeights92b = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5, 0}, []float64{16000, 3000, 180, 10, 1, 11000})
	scatterWeights92c = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5, 0}, []float64{20000, 2800, 560, 10, 1, 70000})
	scatterWeights92d = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5, 0}, []float64{23000, 2600, 620, 10, 1, 80000})
	scatterWeights92e = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5, 0}, []float64{25000, 2400, 680, 10, 1, 90000})

	weights94 = weights92

	wildWeights94a = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5, 0}, []float64{820, 200, 32, 6, 1, 2100})
	wildWeights94b = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5, 0}, []float64{3800, 490, 50, 9, 1, 2400})
	wildWeights94c = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5, 0}, []float64{3600, 430, 45, 8, 1, 2700})
	wildWeights94d = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5, 0}, []float64{3400, 380, 40, 7, 1, 3100})
	wildWeights94e = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5, 0}, []float64{3200, 340, 35, 6, 1, 3900})

	scatterWeights94a = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5, 0}, []float64{4800, 900, 180, 11, 1, 24000})
	scatterWeights94b = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5, 0}, []float64{16000, 3000, 180, 10, 1, 11000})
	scatterWeights94c = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5, 0}, []float64{20000, 2800, 560, 10, 1, 70000})
	scatterWeights94d = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5, 0}, []float64{23000, 2600, 620, 10, 1, 80000})
	scatterWeights94e = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5, 0}, []float64{25000, 2400, 680, 10, 1, 90000})

	weights96 = weights92

	wildWeights96a = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5, 0}, []float64{820, 200, 32, 6, 1, 2100})
	wildWeights96b = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5, 0}, []float64{3800, 490, 50, 9, 1, 2400})
	wildWeights96c = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5, 0}, []float64{3600, 430, 45, 8, 1, 2700})
	wildWeights96d = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5, 0}, []float64{3400, 380, 40, 7, 1, 3100})
	wildWeights96e = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5, 0}, []float64{3200, 340, 35, 6, 1, 3900})

	scatterWeights96a = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5, 0}, []float64{4800, 900, 180, 11, 1, 24000})
	scatterWeights96b = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5, 0}, []float64{16000, 3000, 180, 10, 1, 11000})
	scatterWeights96c = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5, 0}, []float64{20000, 2800, 560, 10, 1, 70000})
	scatterWeights96d = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5, 0}, []float64{23000, 2600, 620, 10, 1, 80000})
	scatterWeights96e = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5, 0}, []float64{25000, 2400, 680, 10, 1, 90000})
)

var (
	wildReelWeights1    = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5}, []float64{1, 1, 4, 5, 6})
	scatterReelWeights1 = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5}, []float64{6, 5, 4, 3, 2})
	wildReelWeights2    = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5}, []float64{7, 7, 8, 9, 9})
	scatterReelWeights2 = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5}, []float64{1, 1, 1, 1, 1})

	multiplierScale = []float64{2, 2, 2, 10, 10, 10, 10, 50, 50, 50, 50, 50, 150}

	multiplierFreeGames = map[int]uint8{
		10:  5,
		50:  5,
		150: 5,
	}

	firstReelSet1 = comp.NewSymbolReels(
		comp.NewSymbolReel(rows, 6, 5, 9, 3, 1, 2, 4, 9, 3, 2, 7, 4, 6, 10, 1, 5, 9, 4, 2, 7, 1, 5, 9, 4, 6, 8, 2, 3, 7, 6, 5, 8, 4, 3, 7, 2, 1, 10, 3, 6, 7, 1, 5, 9, 4, 2, 8, 6, 5, 9, 4, 3, 7, 5, 1, 9, 2, 1, 8, 4, 6, 10, 3, 2, 8, 5, 6, 9, 1, 3, 7, 4, 1, 8, 6, 5, 7, 4, 2, 1, 3, 8, 2, 4, 9, 5, 6, 7, 2, 3, 9, 1, 2, 8, 3, 1, 10, 4, 3, 7, 6, 5, 9, 1, 2, 8, 6, 3, 9, 5, 2, 8, 1, 4, 7, 9, 5, 1, 7, 3, 4, 9, 6, 2, 9, 8, 1, 6, 9, 4, 2, 8, 3, 5, 9, 7, 1, 2, 7, 3, 6, 9, 5, 4, 9, 8, 6, 1, 7, 2, 4, 8, 5, 3, 10, 9, 2, 1, 9, 3, 6, 9, 4, 5, 8, 7, 6, 3, 9, 5, 2, 8, 1, 4, 9, 7, 5, 1, 7, 3, 4, 9, 6, 2, 8, 9, 1, 6, 9, 4, 2, 8, 3, 5, 7, 1, 2, 7, 3, 6, 9, 5, 4, 8, 6, 1, 7, 2, 4, 8, 5, 3, 9, 10, 2, 1, 9, 3, 6, 9, 4, 5, 7, 8),
		comp.NewSymbolReel(rows, 6, 5, 10, 3, 1, 2, 4, 7, 3, 2, 10, 4, 6, 8, 1, 5, 9, 4, 2, 7, 1, 5, 10, 4, 6, 7, 2, 3, 8, 6, 5, 4, 3, 7, 2, 1, 10, 3, 6, 7, 1, 5, 9, 4, 2, 8, 6, 5, 10, 4, 3, 8, 5, 1, 10, 2, 1, 7, 4, 6, 3, 2, 8, 5, 6, 10, 1, 3, 7, 4, 1, 8, 6, 5, 7, 4, 2, 3, 1, 8, 2, 4, 9, 5, 6, 7, 2, 3, 1, 4, 8, 3, 1, 10, 2, 3, 7, 6, 5, 10, 1, 2, 8, 6, 3, 9, 5, 2, 8, 1, 4, 7, 10, 5, 1, 7, 3, 4, 10, 6, 2, 10, 8, 1, 6, 10, 4, 2, 8, 3, 5, 10, 7, 1, 2, 7, 3, 6, 10, 5, 4, 10, 8, 6, 1, 7, 2, 4, 8, 5, 3, 10, 9, 2, 1, 10, 3, 6, 10, 4, 5, 8, 7, 6, 3, 10, 5, 2, 8, 1, 4, 10, 7, 5, 1, 7, 3, 4, 10, 6, 2, 8, 1, 6, 10, 4, 2, 8, 3, 5, 7, 10, 1, 2, 7, 3, 6, 10, 5, 4, 8, 6, 1, 7, 2, 4, 8, 5, 3, 9, 10, 2, 1, 10, 3, 6, 10, 4, 5, 8, 7),
		comp.NewSymbolReel(rows, 6, 5, 7, 3, 1, 9, 10, 2, 4, 10, 3, 2, 9, 4, 6, 10, 9, 1, 5, 9, 4, 2, 7, 1, 5, 9, 4, 6, 8, 10, 2, 3, 7, 6, 5, 8, 4, 3, 7, 2, 1, 10, 9, 3, 6, 7, 1, 5, 9, 4, 2, 8, 6, 5, 10, 8, 4, 3, 7, 5, 1, 9, 2, 1, 8, 4, 6, 10, 9, 3, 2, 8, 5, 6, 9, 1, 3, 7, 4, 1, 8, 10, 6, 5, 7, 4, 2, 10, 1, 3, 8, 2, 4, 10, 5, 6, 7, 10, 2, 3, 9, 1, 2, 8, 3, 1, 10, 9, 4, 3, 7, 6, 5, 9, 1, 2, 8, 7, 6, 3, 10, 5, 2, 8, 1, 4, 7, 8, 5, 1, 7, 3, 4, 9, 6, 2, 10, 8, 1, 6, 9, 4, 2, 8, 3, 5, 10, 7, 1, 2, 7, 3, 6, 10, 5, 4, 9, 8, 6, 1, 7, 2, 4, 8, 5, 3, 10, 9, 2, 1, 10, 3, 6, 9, 4, 5, 8, 7, 6, 3, 10, 5, 2, 8, 1, 4, 9, 7, 5, 1, 7, 3, 4, 9, 6, 2, 8, 10, 1, 6, 9, 4, 2, 8, 3, 5, 7, 10, 1, 2, 7, 3, 6, 10, 5, 4, 8, 9, 6, 1, 7, 2, 4, 8, 5, 3, 9, 10, 2, 1, 10, 3, 6, 9, 4, 5, 7, 8),
		comp.NewSymbolReel(rows, 6, 5, 7, 3, 1, 9, 10, 2, 4, 10, 3, 2, 9, 4, 6, 10, 9, 1, 5, 9, 4, 2, 7, 1, 5, 9, 4, 6, 8, 10, 2, 3, 7, 6, 5, 8, 4, 3, 7, 2, 1, 10, 9, 3, 6, 7, 1, 5, 9, 4, 2, 8, 6, 5, 10, 8, 4, 3, 7, 5, 1, 9, 2, 1, 8, 4, 6, 10, 9, 3, 2, 8, 5, 6, 9, 1, 3, 7, 4, 1, 8, 10, 6, 5, 7, 4, 2, 10, 1, 3, 8, 2, 4, 10, 5, 6, 7, 10, 2, 3, 9, 1, 2, 8, 3, 1, 10, 9, 4, 3, 7, 6, 5, 9, 1, 2, 8, 7, 6, 3, 10, 5, 2, 8, 1, 4, 7, 8, 5, 1, 7, 3, 4, 9, 6, 2, 10, 8, 1, 6, 9, 4, 2, 8, 3, 5, 10, 7, 1, 2, 7, 3, 6, 10, 5, 4, 9, 8, 6, 1, 7, 2, 4, 8, 5, 3, 10, 9, 2, 1, 10, 3, 6, 9, 4, 5, 8, 7, 6, 3, 10, 5, 2, 8, 1, 4, 9, 7, 5, 1, 7, 3, 4, 9, 6, 2, 8, 10, 1, 6, 9, 4, 2, 8, 3, 5, 7, 10, 1, 2, 7, 3, 6, 10, 5, 4, 8, 9, 6, 1, 7, 2, 4, 8, 5, 3, 9, 10, 2, 1, 10, 3, 6, 9, 4, 5, 7, 8),
		comp.NewSymbolReel(rows, 6, 5, 7, 3, 1, 9, 10, 2, 4, 10, 3, 2, 9, 4, 6, 10, 9, 1, 5, 9, 4, 2, 7, 1, 5, 9, 4, 6, 8, 10, 2, 3, 7, 6, 5, 8, 4, 3, 7, 2, 1, 10, 9, 3, 6, 7, 1, 5, 9, 4, 2, 8, 6, 5, 10, 8, 4, 3, 7, 5, 1, 9, 2, 1, 8, 4, 6, 10, 9, 3, 2, 8, 5, 6, 9, 1, 3, 7, 4, 1, 8, 10, 6, 5, 7, 4, 2, 10, 1, 3, 8, 2, 4, 10, 5, 6, 7, 10, 2, 3, 9, 1, 2, 8, 3, 1, 10, 9, 4, 3, 7, 6, 5, 9, 1, 2, 8, 7, 6, 3, 10, 5, 2, 8, 1, 4, 7, 8, 5, 1, 7, 3, 4, 9, 6, 2, 10, 8, 1, 6, 9, 4, 2, 8, 3, 5, 10, 7, 1, 2, 7, 3, 6, 10, 5, 4, 9, 8, 6, 1, 7, 2, 4, 8, 5, 3, 10, 9, 2, 1, 10, 3, 6, 9, 4, 5, 8, 7, 6, 3, 10, 5, 2, 8, 1, 4, 9, 7, 5, 1, 7, 3, 4, 9, 6, 2, 8, 10, 1, 6, 9, 4, 2, 8, 3, 5, 7, 10, 1, 2, 7, 3, 6, 10, 5, 4, 8, 9, 6, 1, 7, 2, 4, 8, 5, 3, 9, 10, 2, 1, 10, 3, 6, 9, 4, 5, 7, 8),
	).WithFlag(flagFirstReelSet1)

	firstReelSet2 = comp.NewSymbolReels(
		comp.NewSymbolReel(rows, 6, 5, 10, 3, 1, 2, 4, 7, 3, 2, 10, 4, 6, 8, 1, 5, 9, 4, 2, 7, 1, 5, 10, 4, 6, 7, 2, 3, 8, 6, 5, 4, 3, 7, 2, 1, 10, 3, 6, 7, 1, 5, 9, 4, 2, 8, 6, 5, 10, 4, 3, 8, 5, 1, 10, 2, 1, 7, 4, 6, 3, 2, 8, 5, 6, 10, 1, 3, 7, 4, 1, 8, 6, 5, 7, 4, 2, 3, 1, 8, 2, 4, 9, 5, 6, 7, 2, 3, 1, 4, 8, 3, 1, 10, 2, 3, 7, 6, 5, 10, 1, 2, 8, 6, 3, 9, 5, 2, 8, 1, 4, 7, 10, 5, 1, 7, 3, 4, 10, 6, 2, 10, 8, 1, 6, 10, 4, 2, 8, 3, 5, 10, 7, 1, 2, 7, 3, 6, 10, 5, 4, 10, 8, 6, 1, 7, 2, 4, 8, 5, 3, 10, 9, 2, 1, 10, 3, 6, 10, 4, 5, 8, 7, 6, 3, 10, 5, 2, 8, 1, 4, 10, 7, 5, 1, 7, 3, 4, 10, 6, 2, 8, 1, 6, 10, 4, 2, 8, 3, 5, 7, 10, 1, 2, 7, 3, 6, 10, 5, 4, 8, 6, 1, 7, 2, 4, 8, 5, 3, 9, 10, 2, 1, 10, 3, 6, 10, 4, 5, 8, 7),
		comp.NewSymbolReel(rows, 6, 5, 9, 3, 1, 2, 4, 9, 3, 2, 7, 4, 6, 10, 1, 5, 9, 4, 2, 7, 1, 5, 9, 4, 6, 8, 2, 3, 7, 6, 5, 8, 4, 3, 7, 2, 1, 10, 3, 6, 7, 1, 5, 9, 4, 2, 8, 6, 5, 9, 4, 3, 7, 5, 1, 9, 2, 1, 8, 4, 6, 10, 3, 2, 8, 5, 6, 9, 1, 3, 7, 4, 1, 8, 6, 5, 7, 4, 2, 1, 3, 8, 2, 4, 9, 5, 6, 7, 2, 3, 9, 1, 2, 8, 3, 1, 10, 4, 3, 7, 6, 5, 9, 1, 2, 8, 6, 3, 9, 5, 2, 8, 1, 4, 7, 9, 5, 1, 7, 3, 4, 9, 6, 2, 9, 8, 1, 6, 9, 4, 2, 8, 3, 5, 9, 7, 1, 2, 7, 3, 6, 9, 5, 4, 9, 8, 6, 1, 7, 2, 4, 8, 5, 3, 10, 9, 2, 1, 9, 3, 6, 9, 4, 5, 8, 7, 6, 3, 9, 5, 2, 8, 1, 4, 9, 7, 5, 1, 7, 3, 4, 9, 6, 2, 8, 9, 1, 6, 9, 4, 2, 8, 3, 5, 7, 1, 2, 7, 3, 6, 9, 5, 4, 8, 6, 1, 7, 2, 4, 8, 5, 3, 9, 10, 2, 1, 9, 3, 6, 9, 4, 5, 7, 8),
		comp.NewSymbolReel(rows, 6, 5, 7, 3, 1, 9, 10, 2, 4, 10, 3, 2, 9, 4, 6, 10, 9, 1, 5, 9, 4, 2, 7, 1, 5, 9, 4, 6, 8, 10, 2, 3, 7, 6, 5, 8, 4, 3, 7, 2, 1, 10, 9, 3, 6, 7, 1, 5, 9, 4, 2, 8, 6, 5, 10, 8, 4, 3, 7, 5, 1, 9, 2, 1, 8, 4, 6, 10, 9, 3, 2, 8, 5, 6, 9, 1, 3, 7, 4, 1, 8, 10, 6, 5, 7, 4, 2, 10, 1, 3, 8, 2, 4, 10, 5, 6, 7, 10, 2, 3, 9, 1, 2, 8, 3, 1, 10, 9, 4, 3, 7, 6, 5, 9, 1, 2, 8, 7, 6, 3, 10, 5, 2, 8, 1, 4, 7, 8, 5, 1, 7, 3, 4, 9, 6, 2, 10, 8, 1, 6, 9, 4, 2, 8, 3, 5, 10, 7, 1, 2, 7, 3, 6, 10, 5, 4, 9, 8, 6, 1, 7, 2, 4, 8, 5, 3, 10, 9, 2, 1, 10, 3, 6, 9, 4, 5, 8, 7, 6, 3, 10, 5, 2, 8, 1, 4, 9, 7, 5, 1, 7, 3, 4, 9, 6, 2, 8, 10, 1, 6, 9, 4, 2, 8, 3, 5, 7, 10, 1, 2, 7, 3, 6, 10, 5, 4, 8, 9, 6, 1, 7, 2, 4, 8, 5, 3, 9, 10, 2, 1, 10, 3, 6, 9, 4, 5, 7, 8),
		comp.NewSymbolReel(rows, 6, 5, 7, 3, 1, 9, 10, 2, 4, 10, 3, 2, 9, 4, 6, 10, 9, 1, 5, 9, 4, 2, 7, 1, 5, 9, 4, 6, 8, 10, 2, 3, 7, 6, 5, 8, 4, 3, 7, 2, 1, 10, 9, 3, 6, 7, 1, 5, 9, 4, 2, 8, 6, 5, 10, 8, 4, 3, 7, 5, 1, 9, 2, 1, 8, 4, 6, 10, 9, 3, 2, 8, 5, 6, 9, 1, 3, 7, 4, 1, 8, 10, 6, 5, 7, 4, 2, 10, 1, 3, 8, 2, 4, 10, 5, 6, 7, 10, 2, 3, 9, 1, 2, 8, 3, 1, 10, 9, 4, 3, 7, 6, 5, 9, 1, 2, 8, 7, 6, 3, 10, 5, 2, 8, 1, 4, 7, 8, 5, 1, 7, 3, 4, 9, 6, 2, 10, 8, 1, 6, 9, 4, 2, 8, 3, 5, 10, 7, 1, 2, 7, 3, 6, 10, 5, 4, 9, 8, 6, 1, 7, 2, 4, 8, 5, 3, 10, 9, 2, 1, 10, 3, 6, 9, 4, 5, 8, 7, 6, 3, 10, 5, 2, 8, 1, 4, 9, 7, 5, 1, 7, 3, 4, 9, 6, 2, 8, 10, 1, 6, 9, 4, 2, 8, 3, 5, 7, 10, 1, 2, 7, 3, 6, 10, 5, 4, 8, 9, 6, 1, 7, 2, 4, 8, 5, 3, 9, 10, 2, 1, 10, 3, 6, 9, 4, 5, 7, 8),
		comp.NewSymbolReel(rows, 6, 5, 7, 3, 1, 9, 10, 2, 4, 10, 3, 2, 9, 4, 6, 10, 9, 1, 5, 9, 4, 2, 7, 1, 5, 9, 4, 6, 8, 10, 2, 3, 7, 6, 5, 8, 4, 3, 7, 2, 1, 10, 9, 3, 6, 7, 1, 5, 9, 4, 2, 8, 6, 5, 10, 8, 4, 3, 7, 5, 1, 9, 2, 1, 8, 4, 6, 10, 9, 3, 2, 8, 5, 6, 9, 1, 3, 7, 4, 1, 8, 10, 6, 5, 7, 4, 2, 10, 1, 3, 8, 2, 4, 10, 5, 6, 7, 10, 2, 3, 9, 1, 2, 8, 3, 1, 10, 9, 4, 3, 7, 6, 5, 9, 1, 2, 8, 7, 6, 3, 10, 5, 2, 8, 1, 4, 7, 8, 5, 1, 7, 3, 4, 9, 6, 2, 10, 8, 1, 6, 9, 4, 2, 8, 3, 5, 10, 7, 1, 2, 7, 3, 6, 10, 5, 4, 9, 8, 6, 1, 7, 2, 4, 8, 5, 3, 10, 9, 2, 1, 10, 3, 6, 9, 4, 5, 8, 7, 6, 3, 10, 5, 2, 8, 1, 4, 9, 7, 5, 1, 7, 3, 4, 9, 6, 2, 8, 10, 1, 6, 9, 4, 2, 8, 3, 5, 7, 10, 1, 2, 7, 3, 6, 10, 5, 4, 8, 9, 6, 1, 7, 2, 4, 8, 5, 3, 9, 10, 2, 1, 10, 3, 6, 9, 4, 5, 7, 8),
	).WithFlag(flagFirstReelSet2)

	firstReelSets = comp.NewMultiSymbolReels(firstReelSet1, firstReelSet2).WithFlag(flagFirstReelSets)

	freeReelSet = comp.NewSymbolReels(
		comp.NewSymbolReel(rows, 6, 5, 7, 3, 1, 9, 2, 4, 10, 3, 2, 9, 4, 6, 10, 1, 5, 9, 4, 2, 7, 1, 5, 9, 4, 6, 8, 2, 3, 7, 6, 5, 8, 4, 3, 7, 2, 1, 10, 3, 6, 7, 1, 5, 9, 4, 2, 8, 6, 5, 10, 4, 3, 7, 5, 1, 9, 2, 1, 8, 4, 6, 10, 3, 2, 8, 5, 6, 9, 1, 3, 7, 4, 1, 8, 6, 5, 7, 4, 2, 10, 1, 3, 8, 2, 4, 10, 5, 6, 7, 2, 3, 9, 1, 2, 8, 3, 1, 10, 4, 3, 7, 6, 5, 9, 1, 2, 8, 9, 6, 3, 10, 5, 2, 8, 1, 4, 7, 5, 1, 7, 3, 4, 9, 6, 2, 10, 1, 6, 9, 4, 2, 8, 3, 5, 10, 1, 2, 7, 3, 6, 10, 5, 4, 9, 8, 6, 1, 7, 2, 4, 8, 5, 3, 10, 2, 1, 10, 3, 6, 9, 4, 5, 8, 6, 3, 10, 5, 2, 8, 1, 4, 9, 7, 5, 1, 7, 3, 4, 9, 6, 2, 8, 1, 6, 9, 4, 2, 8, 3, 5, 7, 10, 1, 2, 7, 3, 6, 10, 5, 4, 8, 7, 6, 1, 7, 2, 4, 8, 5, 3, 9, 2, 1, 10, 3, 6, 9, 4, 5, 7, 8),
		comp.NewSymbolReel(rows, 6, 5, 7, 3, 1, 9, 2, 4, 10, 3, 2, 9, 4, 6, 10, 1, 5, 9, 4, 2, 7, 1, 5, 9, 4, 6, 8, 2, 3, 7, 6, 5, 8, 4, 3, 7, 2, 1, 10, 3, 6, 7, 1, 5, 9, 4, 2, 8, 6, 5, 10, 4, 3, 7, 5, 1, 9, 2, 1, 8, 4, 6, 10, 3, 2, 8, 5, 6, 9, 1, 3, 7, 4, 1, 8, 6, 5, 7, 4, 2, 10, 1, 3, 8, 2, 4, 10, 5, 6, 7, 2, 3, 9, 1, 2, 8, 3, 1, 10, 4, 3, 7, 6, 5, 9, 1, 2, 8, 9, 6, 3, 10, 5, 2, 8, 1, 4, 7, 5, 1, 7, 3, 4, 9, 6, 2, 10, 1, 6, 9, 4, 2, 8, 3, 5, 10, 1, 2, 7, 3, 6, 10, 5, 4, 9, 8, 6, 1, 7, 2, 4, 8, 5, 3, 10, 2, 1, 10, 3, 6, 9, 4, 5, 8, 6, 3, 10, 5, 2, 8, 1, 4, 9, 7, 5, 1, 7, 3, 4, 9, 6, 2, 8, 1, 6, 9, 4, 2, 8, 3, 5, 7, 10, 1, 2, 7, 3, 6, 10, 5, 4, 8, 7, 6, 1, 7, 2, 4, 8, 5, 3, 9, 2, 1, 10, 3, 6, 9, 4, 5, 7, 8),
		comp.NewSymbolReel(rows, 6, 5, 7, 3, 1, 9, 2, 4, 10, 3, 2, 9, 4, 6, 10, 1, 5, 9, 4, 2, 7, 1, 5, 9, 4, 6, 8, 2, 3, 7, 6, 5, 8, 4, 3, 7, 2, 1, 10, 3, 6, 7, 1, 5, 9, 4, 2, 8, 6, 5, 10, 4, 3, 7, 5, 1, 9, 2, 1, 8, 4, 6, 10, 3, 2, 8, 5, 6, 9, 1, 3, 7, 4, 1, 8, 6, 5, 7, 4, 2, 10, 1, 3, 8, 2, 4, 10, 5, 6, 7, 2, 3, 9, 1, 2, 8, 3, 1, 10, 4, 3, 7, 6, 5, 9, 1, 2, 8, 9, 6, 3, 10, 5, 2, 8, 1, 4, 7, 5, 1, 7, 3, 4, 9, 6, 2, 10, 1, 6, 9, 4, 2, 8, 3, 5, 10, 1, 2, 7, 3, 6, 10, 5, 4, 9, 8, 6, 1, 7, 2, 4, 8, 5, 3, 10, 2, 1, 10, 3, 6, 9, 4, 5, 8, 6, 3, 10, 5, 2, 8, 1, 4, 9, 7, 5, 1, 7, 3, 4, 9, 6, 2, 8, 1, 6, 9, 4, 2, 8, 3, 5, 7, 10, 1, 2, 7, 3, 6, 10, 5, 4, 8, 7, 6, 1, 7, 2, 4, 8, 5, 3, 9, 2, 1, 10, 3, 6, 9, 4, 5, 7, 8),
		comp.NewSymbolReel(rows, 6, 5, 7, 3, 1, 9, 2, 4, 10, 3, 2, 9, 4, 6, 10, 1, 5, 9, 4, 2, 7, 1, 5, 9, 4, 6, 8, 2, 3, 7, 6, 5, 8, 4, 3, 7, 2, 1, 10, 3, 6, 7, 1, 5, 9, 4, 2, 8, 6, 5, 10, 4, 3, 7, 5, 1, 9, 2, 1, 8, 4, 6, 10, 3, 2, 8, 5, 6, 9, 1, 3, 7, 4, 1, 8, 6, 5, 7, 4, 2, 10, 1, 3, 8, 2, 4, 10, 5, 6, 7, 2, 3, 9, 1, 2, 8, 3, 1, 10, 4, 3, 7, 6, 5, 9, 1, 2, 8, 9, 6, 3, 10, 5, 2, 8, 1, 4, 7, 5, 1, 7, 3, 4, 9, 6, 2, 10, 1, 6, 9, 4, 2, 8, 3, 5, 10, 1, 2, 7, 3, 6, 10, 5, 4, 9, 8, 6, 1, 7, 2, 4, 8, 5, 3, 10, 2, 1, 10, 3, 6, 9, 4, 5, 8, 6, 3, 10, 5, 2, 8, 1, 4, 9, 7, 5, 1, 7, 3, 4, 9, 6, 2, 8, 1, 6, 9, 4, 2, 8, 3, 5, 7, 10, 1, 2, 7, 3, 6, 10, 5, 4, 8, 7, 6, 1, 7, 2, 4, 8, 5, 3, 9, 2, 1, 10, 3, 6, 9, 4, 5, 7, 8),
		comp.NewSymbolReel(rows, 6, 5, 7, 3, 1, 9, 2, 4, 10, 3, 2, 9, 4, 6, 10, 1, 5, 9, 4, 2, 7, 1, 5, 9, 4, 6, 8, 2, 3, 7, 6, 5, 8, 4, 3, 7, 2, 1, 10, 3, 6, 7, 1, 5, 9, 4, 2, 8, 6, 5, 10, 4, 3, 7, 5, 1, 9, 2, 1, 8, 4, 6, 10, 3, 2, 8, 5, 6, 9, 1, 3, 7, 4, 1, 8, 6, 5, 7, 4, 2, 10, 1, 3, 8, 2, 4, 10, 5, 6, 7, 2, 3, 9, 1, 2, 8, 3, 1, 10, 4, 3, 7, 6, 5, 9, 1, 2, 8, 9, 6, 3, 10, 5, 2, 8, 1, 4, 7, 5, 1, 7, 3, 4, 9, 6, 2, 10, 1, 6, 9, 4, 2, 8, 3, 5, 10, 1, 2, 7, 3, 6, 10, 5, 4, 9, 8, 6, 1, 7, 2, 4, 8, 5, 3, 10, 2, 1, 10, 3, 6, 9, 4, 5, 8, 6, 3, 10, 5, 2, 8, 1, 4, 9, 7, 5, 1, 7, 3, 4, 9, 6, 2, 8, 1, 6, 9, 4, 2, 8, 3, 5, 7, 10, 1, 2, 7, 3, 6, 10, 5, 4, 8, 7, 6, 1, 7, 2, 4, 8, 5, 3, 9, 2, 1, 10, 3, 6, 9, 4, 5, 7, 8),
	).WithFlag(flagFreeReelSet)

	spinner = comp.NewFilteredSpinner(
		[]comp.SpinDataFilterer{comp.OnFreeSpin, comp.OnFirstSpin},
		[]comp.Spinner{freeReelSet, firstReelSets},
	)

	level1 = comp.OnRoundFlagValues(flagScattersLevel, 0, 1, 2, 3)
	level2 = comp.OnRoundFlagValues(flagScattersLevel, 4, 5, 6, 7)
	level3 = comp.OnRoundFlagValues(flagScattersLevel, 8, 9, 10, 11, 12)
	level4 = comp.OnRoundFlagAbove(flagScattersLevel, 12) // 13+

	n01 = comp.WithName("L6")
	n02 = comp.WithName("L5")
	n03 = comp.WithName("L4")
	n04 = comp.WithName("L3")
	n05 = comp.WithName("L2")
	n06 = comp.WithName("L1")
	n07 = comp.WithName("H4")
	n08 = comp.WithName("H3")
	n09 = comp.WithName("H2")
	n10 = comp.WithName("H1")
	n11 = comp.WithName("Scatter")
	n12 = comp.WithName("Wild")

	r01 = comp.WithResource("l6")
	r02 = comp.WithResource("l5")
	r03 = comp.WithResource("l4")
	r04 = comp.WithResource("l3")
	r05 = comp.WithResource("l2")
	r06 = comp.WithResource("l1")
	r07 = comp.WithResource("h4")
	r08 = comp.WithResource("h3")
	r09 = comp.WithResource("h2")
	r10 = comp.WithResource("h1")
	r11 = comp.WithResource("scatter")
	r12 = comp.WithResource("wild")

	scatterPayouts = []float64{0, 1, 4, 20, 100}
	wildPayouts    = []float64{0, 2, 15, 40, 100}

	p01 = comp.WithPayouts(0, 0, 0.5, 2, 6)
	p02 = comp.WithPayouts(0, 0, 0.5, 2, 6)
	p03 = comp.WithPayouts(0, 0, 0.6, 3, 8)
	p04 = comp.WithPayouts(0, 0, 0.6, 3, 8)
	p05 = comp.WithPayouts(0, 0, 0.8, 4, 10)
	p06 = comp.WithPayouts(0, 0, 0.8, 4, 10)
	p07 = comp.WithPayouts(0, 0, 1, 5, 15)
	p08 = comp.WithPayouts(0, 0, 1.5, 6, 20)
	p09 = comp.WithPayouts(0, 0.5, 2, 10, 30)
	p10 = comp.WithPayouts(0, 1, 3, 20, 50)
	p11 = comp.WithScatterPayouts(scatterPayouts...)
	p12 = comp.WithPayouts(wildPayouts...)

	p11111 = comp.NewPayline(id01, rows, 1, 1, 1, 1, 1)
	p22222 = comp.NewPayline(id03, rows, 2, 2, 2, 2, 2)
	p00000 = comp.NewPayline(id02, rows, 0, 0, 0, 0, 0)
	p01210 = comp.NewPayline(id04, rows, 0, 1, 2, 1, 0)
	p21012 = comp.NewPayline(id05, rows, 2, 1, 0, 1, 2)
	p10001 = comp.NewPayline(id07, rows, 1, 0, 0, 0, 1)
	p12221 = comp.NewPayline(id06, rows, 1, 2, 2, 2, 1)
	p21112 = comp.NewPayline(id09, rows, 2, 1, 1, 1, 2)
	p01110 = comp.NewPayline(id08, rows, 0, 1, 1, 1, 0)
	p00122 = comp.NewPayline(id11, rows, 0, 0, 1, 2, 2)
	p22100 = comp.NewPayline(id10, rows, 2, 2, 1, 0, 0)
	p12121 = comp.NewPayline(id13, rows, 1, 2, 1, 2, 1)
	p10101 = comp.NewPayline(id12, rows, 1, 0, 1, 0, 1)
	p21212 = comp.NewPayline(id15, rows, 2, 1, 2, 1, 2)
	p01010 = comp.NewPayline(id14, rows, 0, 1, 0, 1, 0)

	ids       = [symbolCount]util.Index{id01, id02, id03, id04, id05, id06, id07, id08, id09, id10, id11, id12}
	names     = [symbolCount]comp.SymbolOption{n01, n02, n03, n04, n05, n06, n07, n08, n09, n10, n11, n12}
	resources = [symbolCount]comp.SymbolOption{r01, r02, r03, r04, r05, r06, r07, r08, r09, r10, r11, r12}
	payouts   = [symbolCount]comp.SymbolOption{p01, p02, p03, p04, p05, p06, p07, p08, p09, p10, p11, p12}
	paylines  = comp.Paylines{p11111, p22222, p00000, p01210, p21012, p12221, p10001, p21112, p01110, p00122, p22100, p12121, p10101, p21212, p01010}
	flag0     = comp.NewRoundFlag(flagBonusBuy, "bonus buy")
	flag1     = comp.NewRoundFlag(flagFreeSpins, "free spins count")
	flag2     = comp.NewRoundFlag(flagFirstReelSet1, "first reel set (1)")
	flag3     = comp.NewRoundFlag(flagFirstReelSet2, "first reel set (2)")
	flag4     = comp.NewRoundFlag(flagFirstReelSets, "first reel sets")
	flag5     = comp.NewRoundFlag(flagFreeReelSet, "free reel set")
	flag6     = comp.NewRoundFlag(flagScattersLevel, "scatters level")
	flags     = comp.RoundFlags{flag0, flag1, flag2, flag3, flag4, flag5, flag6}
)

var (
	symbols *comp.SymbolSet

	actions92all     comp.SpinActions
	actions92first   comp.SpinActions
	actions92free    comp.SpinActions
	actions92firstBB comp.SpinActions
	actions92freeBB  comp.SpinActions
	actions94all     comp.SpinActions
	actions94first   comp.SpinActions
	actions94free    comp.SpinActions
	actions94firstBB comp.SpinActions
	actions94freeBB  comp.SpinActions
	actions96all     comp.SpinActions
	actions96first   comp.SpinActions
	actions96free    comp.SpinActions
	actions96firstBB comp.SpinActions
	actions96freeBB  comp.SpinActions

	slots92 *comp.Slots
	slots94 *comp.Slots
	slots96 *comp.Slots

	slots92params game.RegularParams
	slots94params game.RegularParams
	slots96params game.RegularParams
)

func initActions() {
	// bonus buy feature.
	bonusBuy := comp.NewPaidAction(comp.FreeSpins, freeSpinsAwardFirst, bonusBuyCost, scatter, 3).WithFlag(flagBonusBuy, bonusBuyID)
	bonusBuy.Describe(bonusBuyID, "bonus buy feature")

	// initialize bonus game multiplier.
	initMultiplier := comp.NewFirstMultiplierAction(1, multiplierScale[0]).WithFlag(flagScattersLevel)
	initMultiplier.Describe(initMultiplierID, "init bonus game multiplier")

	// generate wild symbols RTP 92.
	wilds92a := comp.NewGenerateReelSymbolAction(wild, wildWeights92a, wildReelWeights1)
	wilds92a.Describe(wilds92aID, "generate wilds - first spin - RTP 92")
	wilds92b := comp.NewGenerateReelSymbolAction(wild, wildWeights92b, wildReelWeights2)
	wilds92b.WithTriggerFilters(level1)
	wilds92b.Describe(wilds92bID, "generate wilds - free spins - lvl 1 - RTP 92")
	wilds92c := comp.NewGenerateReelSymbolAction(wild, wildWeights92c, wildReelWeights2)
	wilds92c.WithTriggerFilters(level2)
	wilds92c.Describe(wilds92cID, "generate wilds - free spins - lvl 2 - RTP 92")
	wilds92d := comp.NewGenerateReelSymbolAction(wild, wildWeights92d, wildReelWeights2)
	wilds92d.WithTriggerFilters(level3)
	wilds92d.Describe(wilds92dID, "generate wilds - free spins - lvl 3 - RTP 92")
	wilds92e := comp.NewGenerateReelSymbolAction(wild, wildWeights92e, wildReelWeights2)
	wilds92e.WithTriggerFilters(level4)
	wilds92e.Describe(wilds92eID, "generate wilds - free spins - lvl 4 - RTP 92")

	// generate wild symbols RTP 94.
	wilds94a := comp.NewGenerateReelSymbolAction(wild, wildWeights94a, wildReelWeights1)
	wilds94a.Describe(wilds94aID, "generate wilds - first spin - RTP 94")
	wilds94b := comp.NewGenerateReelSymbolAction(wild, wildWeights94b, wildReelWeights2)
	wilds94b.WithTriggerFilters(level1)
	wilds94b.Describe(wilds94bID, "generate wilds - free spins - lvl 1 - RTP 94")
	wilds94c := comp.NewGenerateReelSymbolAction(wild, wildWeights94c, wildReelWeights2)
	wilds94c.WithTriggerFilters(level2)
	wilds94c.Describe(wilds94cID, "generate wilds - free spins - lvl 2 - RTP 94")
	wilds94d := comp.NewGenerateReelSymbolAction(wild, wildWeights94d, wildReelWeights2)
	wilds94d.WithTriggerFilters(level3)
	wilds94d.Describe(wilds94dID, "generate wilds - free spins - lvl 3 - RTP 94")
	wilds94e := comp.NewGenerateReelSymbolAction(wild, wildWeights94e, wildReelWeights2)
	wilds94e.WithTriggerFilters(level4)
	wilds94e.Describe(wilds94eID, "generate wilds - free spins - lvl 4 - RTP 94")

	// generate wild symbols RTP 96.
	wilds96a := comp.NewGenerateReelSymbolAction(wild, wildWeights96a, wildReelWeights1)
	wilds96a.Describe(wilds96aID, "generate wilds - first spin - RTP 96")
	wilds96b := comp.NewGenerateReelSymbolAction(wild, wildWeights96b, wildReelWeights2)
	wilds96b.WithTriggerFilters(level1)
	wilds96b.Describe(wilds96bID, "generate wilds - free spins - lvl 1 - RTP 96")
	wilds96c := comp.NewGenerateReelSymbolAction(wild, wildWeights96c, wildReelWeights2)
	wilds96c.WithTriggerFilters(level2)
	wilds96c.Describe(wilds96cID, "generate wilds - free spins - lvl 2 - RTP 96")
	wilds96d := comp.NewGenerateReelSymbolAction(wild, wildWeights96d, wildReelWeights2)
	wilds96d.WithTriggerFilters(level3)
	wilds96d.Describe(wilds96dID, "generate wilds - free spins - lvl 3 - RTP 96")
	wilds96e := comp.NewGenerateReelSymbolAction(wild, wildWeights96e, wildReelWeights2)
	wilds96e.WithTriggerFilters(level4)
	wilds96e.Describe(wilds96eID, "generate wilds - free spins - lvl 4 - RTP 96")

	// generate scatter symbols RTP 92.
	scatters92a := comp.NewGenerateReelSymbolAction(scatter, scatterWeights92a, scatterReelWeights1)
	scatters92a.Describe(scatters92aID, "generate scatters - first spin - RTP 92")
	scatters92b := comp.NewGenerateReelSymbolAction(scatter, scatterWeights92b, scatterReelWeights2)
	scatters92b.WithTriggerFilters(level1)
	scatters92b.Describe(scatters92bID, "generate scatters - free spins - lvl 1 - RTP 92")
	scatters92c := comp.NewGenerateReelSymbolAction(scatter, scatterWeights92c, scatterReelWeights2)
	scatters92c.WithTriggerFilters(level2)
	scatters92c.Describe(scatters92cID, "generate scatters - free spins - lvl 2 - RTP 92")
	scatters92d := comp.NewGenerateReelSymbolAction(scatter, scatterWeights92d, scatterReelWeights2)
	scatters92d.WithTriggerFilters(level3)
	scatters92d.Describe(scatters92dID, "generate scatters - free spins - lvl 3 - RTP 92")
	scatters92e := comp.NewGenerateReelSymbolAction(scatter, scatterWeights92e, scatterReelWeights2)
	scatters92e.WithTriggerFilters(level4)
	scatters92e.Describe(scatters92eID, "generate scatters - free spins - lvl 4 - RTP 92")

	// generate scatter symbols RTP 94.
	scatters94a := comp.NewGenerateReelSymbolAction(scatter, scatterWeights94a, scatterReelWeights1)
	scatters94a.Describe(scatters94aID, "generate scatters - first spin - RTP 94")
	scatters94b := comp.NewGenerateReelSymbolAction(scatter, scatterWeights94b, scatterReelWeights2)
	scatters94b.WithTriggerFilters(level1)
	scatters94b.Describe(scatters94bID, "generate scatters - free spins - lvl 1 - RTP 94")
	scatters94c := comp.NewGenerateReelSymbolAction(scatter, scatterWeights94c, scatterReelWeights2)
	scatters94c.WithTriggerFilters(level2)
	scatters94c.Describe(scatters94cID, "generate scatters - free spins - lvl 2 - RTP 94")
	scatters94d := comp.NewGenerateReelSymbolAction(scatter, scatterWeights94d, scatterReelWeights2)
	scatters94d.WithTriggerFilters(level3)
	scatters94d.Describe(scatters94dID, "generate scatters - free spins - lvl 3 - RTP 94")
	scatters94e := comp.NewGenerateReelSymbolAction(scatter, scatterWeights94e, scatterReelWeights2)
	scatters94e.WithTriggerFilters(level4)
	scatters94e.Describe(scatters94eID, "generate scatters - free spins - lvl 4 - RTP 94")

	// generate scatter symbols RTP 96.
	scatters96a := comp.NewGenerateReelSymbolAction(scatter, scatterWeights96a, scatterReelWeights1)
	scatters96a.Describe(scatters96aID, "generate scatters - first spin - RTP 96")
	scatters96b := comp.NewGenerateReelSymbolAction(scatter, scatterWeights96b, scatterReelWeights2)
	scatters96b.WithTriggerFilters(level1)
	scatters96b.Describe(scatters96bID, "generate scatters - free spins - lvl 1 - RTP 96")
	scatters96c := comp.NewGenerateReelSymbolAction(scatter, scatterWeights96c, scatterReelWeights2)
	scatters96c.WithTriggerFilters(level2)
	scatters96c.Describe(scatters96cID, "generate scatters - free spins - lvl 2 - RTP 96")
	scatters96d := comp.NewGenerateReelSymbolAction(scatter, scatterWeights96d, scatterReelWeights2)
	scatters96d.WithTriggerFilters(level3)
	scatters96d.Describe(scatters96dID, "generate scatters - free spins - lvl 3 - RTP 96")
	scatters96e := comp.NewGenerateReelSymbolAction(scatter, scatterWeights96e, scatterReelWeights2)
	scatters96e.WithTriggerFilters(level4)
	scatters96e.Describe(scatters96eID, "generate scatters - free spins - lvl 4 - RTP 96")

	// calculate winlines.
	winlines := comp.NewPaylinesAction()
	winlines.Describe(winlinesID, "calculate winlines")

	// calculate wild payouts.
	scatterPayouts2 := comp.NewScatterPayoutAction(scatter, 2, scatterPayouts[1])
	scatterPayouts2.Describe(scatterPayouts2ID, "calculate scatter payouts x2")
	scatterPayouts3 := comp.NewScatterPayoutAction(scatter, 3, scatterPayouts[2]).WithAlternate(scatterPayouts2)
	scatterPayouts3.Describe(scatterPayouts3ID, "calculate scatter payouts x3")
	scatterPayouts4 := comp.NewScatterPayoutAction(scatter, 4, scatterPayouts[3]).WithAlternate(scatterPayouts3)
	scatterPayouts4.Describe(scatterPayouts4ID, "calculate scatter payouts x4")
	scatterPayouts5 := comp.NewScatterPayoutAction(scatter, 5, scatterPayouts[4]).WithAlternate(scatterPayouts4)
	scatterPayouts5.Describe(scatterPayouts5ID, "calculate scatter payouts x5")

	// award free spins.
	freeSpinsFirst := comp.NewScatterFreeSpinsAction(freeSpinsAwardFirst, false, scatter, 3, false)
	freeSpinsFirst.Describe(freeSpinsFirstID, "award free spins - first")
	freeSpinsFree := comp.NewScatterFreeSpinsAction(freeSpinsAwardFree, false, scatter, 3, false)
	freeSpinsFree.Describe(freeSpinsFreeID, "award free spins - free")

	// update round flag 3 marking sequence of free spin.
	freeSpinsFlag := comp.NewRoundFlagIncreaseAction(flagFreeSpins)
	freeSpinsFlag.Describe(freeSpinsFlagID, "count number of free spins (flag 1)")

	// round multiplier (goes into effect on the next spin!).
	multiplier := comp.NewMultiplierScaleAction(scatter, 1, multiplierScale...).WithFreeSpins(multiplierFreeGames).WithFlag(flagScattersLevel)
	multiplier.WithStage(comp.AwardBonuses)
	multiplier.Describe(multiplierID, "multiplier scale")

	actionsAall := comp.SpinActions{bonusBuy, initMultiplier}
	actionsAfirst := comp.SpinActions{bonusBuy, initMultiplier}
	actionsAfree := comp.SpinActions{}

	actionsBall := comp.SpinActions{winlines, scatterPayouts5, freeSpinsFirst, freeSpinsFree, freeSpinsFlag, multiplier}
	actionsBfirst := comp.SpinActions{winlines, scatterPayouts5, freeSpinsFirst}
	actionsBfree := comp.SpinActions{winlines, scatterPayouts5, freeSpinsFree, freeSpinsFlag, multiplier}

	actions92all = append(append(actionsAall, wilds92a, wilds92b, wilds92c, wilds92d, wilds92e, scatters92a, scatters92b, scatters92c, scatters92d, scatters92e), actionsBall...)
	actions92first = append(append(actionsAfirst, wilds92a, scatters92a), actionsBfirst...)
	actions92free = append(append(actionsAfree, wilds92b, wilds92c, wilds92d, wilds92e, scatters92b, scatters92c, scatters92d, scatters92e), actionsBfree...)
	actions92firstBB = append(append(actionsAfirst, wilds92a, scatters92a), actionsBfirst...)
	actions92freeBB = append(append(actionsAfree, wilds92b, wilds92c, wilds92d, wilds92e, scatters92b, scatters92c, scatters92d, scatters92e), actionsBfree...)

	actions94all = append(append(actionsAall, wilds94a, wilds94b, wilds94c, wilds94d, wilds94e, scatters94a, scatters94b, scatters94c, scatters94d, scatters94e), actionsBall...)
	actions94first = append(append(actionsAfirst, wilds94a, scatters94a), actionsBfirst...)
	actions94free = append(append(actionsAfree, wilds94b, wilds94c, wilds94d, wilds94e, scatters94b, scatters94c, scatters94d, scatters94e), actionsBfree...)
	actions94firstBB = append(append(actionsAfirst, wilds94a, scatters94a), actionsBfirst...)
	actions94freeBB = append(append(actionsAfree, wilds94b, wilds94c, wilds94d, wilds94e, scatters94b, scatters94c, scatters94d, scatters94e), actionsBfree...)

	actions96all = append(append(actionsAall, wilds96a, wilds96b, wilds96c, wilds96d, wilds96e, scatters96a, scatters96b, scatters96c, scatters96d, scatters96e), actionsBall...)
	actions96first = append(append(actionsAfirst, wilds96a, scatters96a), actionsBfirst...)
	actions96free = append(append(actionsAfree, wilds96b, wilds96c, wilds96d, wilds96e, scatters96b, scatters96c, scatters96d, scatters96e), actionsBfree...)
	actions96firstBB = append(append(actionsAfirst, wilds96a, scatters96a), actionsBfirst...)
	actions96freeBB = append(append(actionsAfree, wilds96b, wilds96c, wilds96d, wilds96e, scatters96b, scatters96c, scatters96d, scatters96e), actionsBfree...)
}

func initSlots(target float64, weights [symbolCount]comp.SymbolOption, actions1, actions2, actions3, actions4 []comp.SpinActioner) *comp.Slots {
	ss := make([]*comp.Symbol, symbolCount)
	for ix := range ss {
		id := ids[ix]
		switch id {
		case scatter:
			ss[ix] = comp.NewSymbol(id, names[ix], resources[ix], payouts[ix], weights[ix], comp.WithKind(comp.Scatter))
		case wild:
			ss[ix] = comp.NewSymbol(id, names[ix], resources[ix], payouts[ix], weights[ix], comp.WithKind(comp.Wild))
		default:
			ss[ix] = comp.NewSymbol(id, names[ix], resources[ix], payouts[ix], weights[ix])
		}
	}
	s := comp.NewSymbolSet(ss...)

	if symbols == nil {
		symbols = s
	}

	return comp.NewSlots(
		comp.Grid(reels, rows),
		comp.WithSpinner(spinner),
		comp.WithSymbols(s),
		comp.WithPaylines(direction, highestPayout, paylines...),
		comp.MaxPayout(maxPayout),
		comp.WithRoundMultiplier(),
		comp.WithMultiplierOnWildsOnly(),
		comp.WithBonusBuy(),
		comp.WithRTP(target),
		comp.WithRoundFlags(flags...),
		comp.WithActions(actions1, actions2, actions3, actions4),
	)
}

func init() {
	initActions()

	slots92 = initSlots(92, weights92, actions92first, actions92free, actions92firstBB, actions92freeBB)
	slots94 = initSlots(94, weights94, actions94first, actions94free, actions94firstBB, actions94freeBB)
	slots96 = initSlots(96, weights96, actions96first, actions96free, actions96firstBB, actions96freeBB)

	slots92params = game.RegularParams{Slots: slots92}
	slots94params = game.RegularParams{Slots: slots94}
	slots96params = game.RegularParams{Slots: slots96}
}

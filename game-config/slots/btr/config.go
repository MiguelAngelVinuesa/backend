package btr

import (
	"fmt"

	comp "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	game "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/slots"
	util "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

const (
	symbolCount  = 15
	reels        = 6
	rows         = 4
	direction    = comp.PayLTR
	maxPayout    = 5000.0
	bonusBuyCost = 150.0

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
	wild2   = id13
	wild3   = id14
	wild4   = id15

	bonusKind = 1

	// always add new ID's with a new unique number!
	bonusBuyID           = 1
	bonusExtraID         = 2
	firstWild92ID        = 101
	firstScatter92ID     = 111
	freeScatter92aID     = 112
	freeScatter92bID     = 122
	freeScatter92cID     = 132
	freeScatter92dID     = 142
	freeScatter92eID     = 152
	firstWild94ID        = 201
	firstScatter94ID     = 211
	freeScatter94aID     = 212
	freeScatter94bID     = 222
	freeScatter94cID     = 232
	freeScatter94dID     = 242
	freeScatter94eID     = 252
	firstWild96ID        = 301
	firstScatter96ID     = 311
	freeScatter96aID     = 312
	freeScatter96bID     = 322
	freeScatter96cID     = 332
	freeScatter96dID     = 342
	freeScatter96eID     = 352
	firstWild2ID         = 401
	freeWild234aID       = 411
	freeWild234bID       = 421
	freeWild234cID       = 431
	freeWild234dID       = 441
	freeWild234eID       = 451
	paylinesID           = 501
	removePayouts1ID     = 511
	removePayouts2ID     = 512
	removePayouts3ID     = 513
	award10ID            = 601
	award15ID            = 602
	award20ID            = 603
	award50ID            = 604
	retrigger3ID         = 611
	retrigger4ID         = 612
	retrigger5ID         = 613
	retrigger6ID         = 614
	countScattersID      = 621
	retriggerScatters3ID = 631
	markRetriggers3ID    = 641
	awardRetrigger3ID    = 651
	forceRetrigger3ID    = 652
	reduceRetriggersID   = 691
	stickies2ID          = 701
	payoutBands3ID       = 711
	payoutBands4ID       = 712
	payoutBands5ID       = 713
	payoutBands6ID       = 714
	payoutBands33ID      = 715
	payoutBands3BBID     = 721
	payoutBands4BBID     = 722
	payoutBands5BBID     = 723
	payoutBands6BBID     = 724
	payoutBands33BBID    = 725
	freeCounterID        = 791
	teaserID             = 901
	win100xID            = 911
	dummyScriptID        = 918
	scriptedSpinID       = 919
	maxWinID             = 921
	scriptedMaxWinID     = 929
	forceNoScattersID    = 930
	forceScrolls4ID      = 931
	forceBand5ID         = 932
	forceWildsX4ID       = 933
	script1ScattersR3ID  = 935
	script1Page1ID       = 936
	script1Page2ID       = 937
	script1Page3ID       = 938
	script2ForceFPID     = 941
	script3WildsX4ID     = 942
	script4WildsXID      = 943
	script5ScattersID    = 951
	script5WildsXID      = 952
	script5Retrigger1ID  = 953
	script5Retrigger2ID  = 954
	script6ForceFPID     = 961

	flagFreeSpinBands  = 0
	flagFreeSpinCount  = 1
	flagBonusBuy       = 2
	flagTriggerCount   = 3
	flagScriptedRound  = 4
	flagRetriggerCount = 5
)

var (
	// WIP WIP WIP
	weights92 = [symbolCount]comp.SymbolOption{
		comp.WithWeights(210, 230, 225, 225, 230, 210),
		comp.WithWeights(210, 230, 225, 225, 230, 210),
		comp.WithWeights(230, 225, 175, 175, 225, 230),
		comp.WithWeights(210, 220, 180, 180, 220, 210),
		comp.WithWeights(220, 185, 175, 175, 165, 220),
		comp.WithWeights(155, 180, 145, 145, 180, 155),
		comp.WithWeights(155, 180, 145, 145, 180, 155),
		comp.WithWeights(105, 125, 165, 165, 125, 105),
		comp.WithWeights(120, 95, 115, 115, 95, 120),
		comp.WithWeights(105, 100, 95, 95, 100, 105),
		comp.WithWeights(0, 0, 0, 0, 0, 0), // generated after initial spin.
		comp.WithWeights(0, 0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0, 0),
	}
	firstWildChances92        = []float64{1}
	scatterRegularChances92   = []float64{34, 28, 5.2, 33.6, 4.2, 8.2}
	scatterFreeSpinChances92a = [5][]float64{
		{20, 30, 0.3},
		{21, 31, 0.5, 75},
		{22, 37, 1, 75, 0.1},
		{30, 35, 1.5, 75, 1, 0.5},
		{35, 40, 2, 75, 5, 2}}
	scatterFreeSpinChances92b = [5][]float64{
		{20, 30, 1},
		{21, 31, 6, 1},
		{22, 37, 14.4, 4, 0.1},
		{30, 35, 19.4, 12, 1, 0.5},
		{35, 40, 28, 15, 5, 2}}
	scatterFreeSpinChances92c = [5][]float64{
		{20, 30, 1},
		{21, 31, 6, 1},
		{22, 37, 21.1, 4, 0.1},
		{30, 35, 23.1, 12, 1, 0.5},
		{35, 40, 31, 15, 5, 2}}
	scatterFreeSpinChances92d = [5][]float64{
		{20, 30, 1},
		{21, 31, 6, 1},
		{22, 37, 23.1, 4, 0.1},
		{30, 35, 25.1, 12, 1, 0.5},
		{35, 40, 33, 15, 5, 2}}
	scatterFreeSpinChances92e = [5][]float64{
		{20, 30, 1},
		{21, 31, 5, 1},
		{22, 37, 8, 4, 0.1},
		{30, 35, 17.5, 12, 1, 0.5},
		{35, 40, 25, 15, 5, 2}}

	// WIP WIP WIP
	weights94 = [symbolCount]comp.SymbolOption{
		comp.WithWeights(210, 230, 225, 225, 230, 210),
		comp.WithWeights(210, 230, 225, 225, 230, 210),
		comp.WithWeights(230, 225, 175, 175, 225, 230),
		comp.WithWeights(210, 220, 180, 180, 220, 210),
		comp.WithWeights(220, 185, 175, 175, 165, 220),
		comp.WithWeights(155, 180, 145, 145, 180, 155),
		comp.WithWeights(155, 180, 145, 145, 180, 155),
		comp.WithWeights(105, 125, 165, 165, 125, 105),
		comp.WithWeights(120, 95, 115, 115, 95, 120),
		comp.WithWeights(105, 100, 95, 95, 100, 105),
		comp.WithWeights(0, 0, 0, 0, 0, 0), // generated after initial spin.
		comp.WithWeights(0, 0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0, 0),
	}
	firstWildChances94        = []float64{1}
	scatterRegularChances94   = []float64{34, 28, 5.5, 33.6, 4.2, 8.2}
	scatterFreeSpinChances94a = [5][]float64{
		{20, 30, 0.3},
		{21, 31, 0.5, 75},
		{22, 37, 1, 75, 0.1},
		{30, 35, 1.5, 75, 1, 0.5},
		{35, 40, 2, 75, 5, 2}}
	scatterFreeSpinChances94b = [5][]float64{
		{20, 30, 1},
		{21, 31, 6, 1},
		{22, 37, 16, 4, 0.1},
		{30, 35, 21, 12, 1, 0.5},
		{35, 40, 28, 15, 5, 2}}
	scatterFreeSpinChances94c = [5][]float64{
		{20, 30, 1},
		{21, 31, 6, 1},
		{22, 37, 22, 4, 0.1},
		{30, 35, 24, 12, 1, 0.5},
		{35, 40, 31, 15, 5, 2}}
	scatterFreeSpinChances94d = [5][]float64{
		{20, 30, 1},
		{21, 31, 6, 1},
		{22, 37, 24, 4, 0.1},
		{30, 35, 26, 12, 1, 0.5},
		{35, 40, 33, 15, 5, 2}}
	scatterFreeSpinChances94e = [5][]float64{
		{20, 30, 1},
		{21, 31, 5, 1},
		{22, 37, 8, 4, 0.1},
		{30, 35, 18, 12, 1, 0.5},
		{35, 40, 25, 15, 5, 2}}

	// WIP WIP WIP
	weights96 = [symbolCount]comp.SymbolOption{
		comp.WithWeights(210, 230, 225, 225, 230, 210),
		comp.WithWeights(210, 230, 225, 225, 230, 210),
		comp.WithWeights(230, 225, 175, 175, 225, 230),
		comp.WithWeights(210, 220, 180, 180, 220, 210),
		comp.WithWeights(220, 185, 175, 175, 165, 220),
		comp.WithWeights(155, 180, 145, 145, 180, 155),
		comp.WithWeights(155, 180, 145, 145, 180, 155),
		comp.WithWeights(105, 125, 165, 165, 125, 105),
		comp.WithWeights(120, 95, 115, 115, 95, 120),
		comp.WithWeights(105, 100, 95, 95, 100, 105),
		comp.WithWeights(0, 0, 0, 0, 0, 0), // generated after initial spin.
		comp.WithWeights(0, 0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0, 0),
	}
	firstWildChances96        = []float64{1.4}
	scatterRegularChances96   = []float64{34, 28, 5.74, 33.6, 4.2, 8.2}
	scatterFreeSpinChances96a = [5][]float64{
		{20, 30, 0.3},
		{21, 31, 0.5, 75},
		{22, 37, 1, 75, 0.1},
		{30, 35, 1.8, 75, 1, 0.5},
		{35, 40, 2.2, 75, 5, 2}}
	scatterFreeSpinChances96b = [5][]float64{
		{20, 30, 1},
		{21, 31, 6, 1},
		{22, 37, 17.3, 4, 0.1},
		{30, 35, 22.3, 12, 1, 0.5},
		{35, 40, 28, 15, 5, 2}}
	scatterFreeSpinChances96c = [5][]float64{
		{20, 30, 1},
		{21, 31, 6, 1},
		{22, 37, 23.4, 4, 0.1},
		{30, 35, 25.4, 12, 1, 0.5},
		{35, 40, 31, 15, 5, 2}}
	scatterFreeSpinChances96d = [5][]float64{
		{20, 30, 1},
		{21, 31, 6, 1},
		{22, 37, 25, 4, 0.1},
		{30, 35, 26, 12, 1, 0.5},
		{35, 40, 33, 15, 5, 2}}
	scatterFreeSpinChances96e = [5][]float64{
		{20, 30, 1},
		{21, 31, 5, 1},
		{22, 37, 10, 4, 0.1},
		{30, 35, 18.5, 12, 1, 0.5},
		{35, 40, 25, 15, 5, 2}}
)

var (
	// free-spin wild reels.
	freeWildReels = []uint8{2, 3, 4, 5}

	// free spin payout band weights.
	freeSpinBandWeights3   = []float64{150, 75, 15, 1}
	freeSpinBandWeights4   = []float64{25, 125, 35, 2}
	freeSpinBandWeights5   = []float64{0, 95, 25, 4}
	freeSpinBandWeights6   = []float64{0, 55, 65, 34, 7}
	freeSpinBandWeights3BB = []float64{0, 45, 95}
	freeSpinBandWeights4BB = []float64{0, 55, 135, 4}
	freeSpinBandWeights5BB = []float64{0, 45, 125, 8}
	freeSpinBandWeights6BB = []float64{0, 35, 165, 34, 12}

	// weights to remove payouts based on total payout factor during first spins.
	reduceBandsFirst = []comp.RemovePayoutBand{
		{MinPayout: 0, MaxPayout: 2.5, RemoveChance: 70},
		{MinPayout: 2.5, MaxPayout: 5, RemoveChance: 32.2},
		{MinPayout: 5, MaxPayout: 25, RemoveChance: 42.5},
		{MinPayout: 25, MaxPayout: 30, RemoveChance: 15},
		{MinPayout: 30, MaxPayout: 45, RemoveChance: 58},
		{MinPayout: 45, MaxPayout: 60, RemoveChance: 15},
		{MinPayout: 60, MaxPayout: 85, RemoveChance: 1},
		{MinPayout: 85, MaxPayout: 200, RemoveChance: 75},
		{MinPayout: 200, MaxPayout: maxPayout, RemoveChance: 85},
	}

	// weights to remove payouts based on total payout factor during first few free spins.
	reduceBandsFree1 = []comp.RemovePayoutBand{
		{MinPayout: 0.0, MaxPayout: 30, RemoveChance: 13},
		{MinPayout: 30, MaxPayout: maxPayout, RemoveChance: 42},
	}

	// weights to remove payouts based on total payout factor during free spins.
	reduceBandsFree2 = []comp.RemovePayoutBand{
		{MinPayout: 0, MaxPayout: 5, RemoveChance: 10},
		{MinPayout: 5, MaxPayout: 10, RemoveChance: 5},
		{MinPayout: 10, MaxPayout: maxPayout, RemoveChance: 12},
	}

	// chance to get 2x/3x/4x wild in first spins.
	firstWildChances = []float64{33, 36.2, 12.8, 1.6}

	// chance to get 2x/3x/4x wild in free spins per band.
	freeWildChancesA = [5][]float64{
		{4.1},
		{8.7, 0.7},
		{13.5, 2.4, 0.5},
		{16, 4.3, 0.8},
		{21, 8, 2, 0.5},
	}
	freeWildChancesB = [5][]float64{
		{4.1},
		{8.7, 0.7},
		{13.7, 2.6, 0.5},
		{16.4, 4.3, 0.8},
		{21.4, 8, 2, 0.5},
	}
	freeWildChancesC = [5][]float64{
		{4.1},
		{8.7, 0.7},
		{15, 3, 0.5},
		{17.5, 4.3, 0.8},
		{24, 8, 2, 0.5},
	}
	freeWildChancesD = [5][]float64{
		{4.1},
		{8.7, 0.7},
		{16.5, 2, 0.5},
		{20, 4.3, 0.8},
		{26, 8, 2, 0.5},
	}
	freeWildChancesE = [5][]float64{
		{4.1},
		{8.7, 0.7},
		{13.5, 2.4, 0.5},
		{16, 4.3, 0.8},
		{21, 8, 2, 0.5},
	}

	modBonus3S = comp.NewMultiFunc(
		comp.NewDivideFunc(50, 0, func(spin *comp.Spin) float64 {
			if spin.SpinSeq() < 4 {
				return 100
			}
			return 50 + (0.9 * spin.TotalPayout())
		}),
	)

	modBonus4S = comp.NewMultiFunc(
		comp.NewDivideFunc(33, 0, func(spin *comp.Spin) float64 {
			if spin.SpinSeq() < 8 {
				return 100
			}
			return 33 + (1.6 * spin.TotalPayout())
		}),
	)

	modBonus5S = comp.NewMultiFunc(
		comp.NewDivideFunc(36, 0, func(spin *comp.Spin) float64 {
			if spin.SpinSeq() < 11 {
				return 100
			}
			return 36 + (2 * spin.TotalPayout())
		}),
	)

	modBonus6S = comp.NewMultiFunc(
		comp.NewDivideFunc(15, 0, func(spin *comp.Spin) float64 {
			if spin.SpinSeq() < 22 {
				return 100
			}
			return 15 + (1.3 * spin.TotalPayout())
		}),
	)

	modBonus33S = comp.NewMultiFunc(
		comp.NewDivideFunc(36, 0, func(spin *comp.Spin) float64 {
			if spin.SpinSeq() < 14 {
				return 100
			}
			return 36 + (2 * spin.TotalPayout())
		}),
	)

	modBonus3W = comp.NewMultiFunc(
		comp.NewPowerFunc(10.5, 0, func(spin *comp.Spin) float64 {
			seq := float64(spin.SpinSeq()) - 5
			if seq > 0 {
				return seq / 4.85
			}
			return seq / 24
		}),
		comp.NewDivideFunc(20, 0, func(spin *comp.Spin) float64 {
			return 20 + (1.8 * spin.TotalPayout())
		}),
	)

	modBonus4W = comp.NewMultiFunc(
		comp.NewPowerFunc(15.5, 0, func(spin *comp.Spin) float64 {
			seq := float64(spin.SpinSeq()) - 9
			if seq > 0 {
				return seq / 7.25
			}
			return seq / 35
		}),
		comp.NewDivideFunc(1, 0, func(spin *comp.Spin) float64 {
			if c := spin.CountSymbols(util.Indexes{wild, wild2, wild3, wild4}); c > 3 {
				return float64(c) * 1.5 / 4.5
			}
			return 1
		}),
		comp.NewDivideFunc(40, 0, func(spin *comp.Spin) float64 {
			return 40 + (1.5 * spin.TotalPayout())
		}),
	)

	modBonus5W = comp.NewMultiFunc(
		comp.NewPowerFunc(22.7, 0, func(spin *comp.Spin) float64 {
			seq := float64(spin.SpinSeq()) - 12
			if seq > 0 {
				return seq / 7.3
			}
			return seq / 50
		}),
		comp.NewDivideFunc(1, 0, func(spin *comp.Spin) float64 {
			if c := spin.CountSymbols(util.Indexes{wild, wild2, wild3, wild4}); c > 2 {
				return float64(c) * 3 / 5
			}
			return 1
		}),
		comp.NewDivideFunc(60, 0, func(spin *comp.Spin) float64 {
			return 60 + (1.3 * spin.TotalPayout())
		}),
	)

	modBonus6W = comp.NewMultiFunc(
		comp.NewPowerFunc(40, 0, func(spin *comp.Spin) float64 {
			seq := float64(spin.SpinSeq()) - 35
			if seq > 0 {
				return seq / 11
			}
			return seq / 100
		}),
		comp.NewDivideFunc(1, 0, func(spin *comp.Spin) float64 {
			if c := spin.CountSymbols(util.Indexes{wild, wild2, wild3, wild4}); c > 3 {
				return float64(c) * 1.5 / 4
			}
			return 1
		}),
		comp.NewDivideFunc(80, 0, func(spin *comp.Spin) float64 {
			return 80 + (1.1 * spin.TotalPayout())
		}),
	)

	modBonus33W = comp.NewMultiFunc(
		comp.NewPowerFunc(22.7, 0, func(spin *comp.Spin) float64 {
			seq := float64(spin.SpinSeq()) - 12
			if seq > 0 {
				return seq / 6.9
			}
			return seq / 50
		}),
		comp.NewDivideFunc(1, 0, func(spin *comp.Spin) float64 {
			if c := spin.CountSymbols(util.Indexes{wild, wild2, wild3, wild4}); c > 2 {
				return float64(c) * 3 / 5
			}
			return 1
		}),
		comp.NewDivideFunc(60, 0, func(spin *comp.Spin) float64 {
			return 60 + (1.3 * spin.TotalPayout())
		}),
	)

	wildWeights = util.AcquireWeighting().AddWeights(util.Indexes{wild, wild2, wild3, wild4}, []float64{200, 50, 5, 1})

	fakeFullPageChance  = 0.5
	fakeFullPageWeights = util.AcquireWeighting().AddWeights(util.Indexes{8, 9, 10}, []float64{40, 33, 27})

	win120SpinChance = 6.7
	wins120x         = []comp.FakeSpin{
		{Indexes: util.Indexes{10, 0, 0, 10, 0, 10, 0, 10, 0, 0, 10, 0, 10, 0, 0, 10, 0, 0, 10, 0, 10, 0, 0, 0}, ReplaceSymbols: util.Indexes{10, wild, wild2, wild3, wild4}},
		{Indexes: util.Indexes{10, 0, 10, 0, 0, 10, 10, 0, 0, 0, 0, 10, 0, 10, 0, 10, 0, 10, 0, 0, 0, 10, 0, 0}, ReplaceSymbols: util.Indexes{10, wild, wild2, wild3, wild4}},
		{Indexes: util.Indexes{10, 10, 0, 0, 10, 0, 0, 10, 0, 0, 10, 0, 10, 0, 10, 0, 0, 0, 10, 0, 0, 0, 0, 10}, ReplaceSymbols: util.Indexes{10, wild, wild2, wild3, wild4}},
		{Indexes: util.Indexes{0, 10, 0, 10, 0, 0, 10, 10, 0, 10, 0, 0, 10, 0, 10, 0, 10, 0, 0, 0, 10, 0, 0, 0}, ReplaceSymbols: util.Indexes{10, wild, wild2, wild3, wild4}},
		{Indexes: util.Indexes{0, 10, 10, 0, 0, 10, 0, 10, 10, 0, 0, 0, 0, 10, 0, 10, 0, 10, 0, 0, 0, 0, 10, 0}, ReplaceSymbols: util.Indexes{10, wild, wild2, wild3, wild4}},
		{Indexes: util.Indexes{0, 0, 10, 10, 0, 10, 10, 0, 0, 10, 0, 0, 0, 0, 10, 10, 0, 0, 0, 10, 0, 10, 0, 0}, ReplaceSymbols: util.Indexes{10, wild, wild2, wild3, wild4}},
		{Indexes: util.Indexes{10, 0, 0, 10, 0, 0, 10, 10, 0, 0, 0, 10, 0, 10, 0, 10, 0, 0, 0, 10, 10, 0, 0, 0}, ReplaceSymbols: util.Indexes{10, wild, wild2, wild3, wild4}},
		{Indexes: util.Indexes{10, 0, 10, 0, 0, 10, 10, 0, 0, 10, 0, 0, 10, 0, 0, 10, 10, 0, 0, 0, 0, 0, 10, 0}, ReplaceSymbols: util.Indexes{10, wild, wild2, wild3, wild4}},
		{Indexes: util.Indexes{10, 10, 0, 0, 10, 10, 0, 0, 0, 0, 10, 0, 10, 0, 10, 0, 0, 0, 10, 0, 0, 10, 0, 0}, ReplaceSymbols: util.Indexes{10, wild, wild2, wild3, wild4}},
		{Indexes: util.Indexes{0, 10, 0, 10, 0, 0, 10, 10, 0, 10, 0, 0, 0, 10, 10, 0, 0, 0, 0, 10, 0, 0, 0, 10}, ReplaceSymbols: util.Indexes{10, wild, wild2, wild3, wild4}},
		{Indexes: util.Indexes{0, 10, 10, 0, 10, 0, 10, 0, 0, 0, 0, 10, 0, 10, 0, 10, 10, 0, 0, 0, 0, 0, 10, 0}, ReplaceSymbols: util.Indexes{10, wild, wild2, wild3, wild4}},
		{Indexes: util.Indexes{0, 0, 10, 10, 10, 0, 0, 10, 10, 0, 0, 0, 0, 0, 10, 10, 0, 10, 0, 0, 10, 0, 0, 0}, ReplaceSymbols: util.Indexes{10, wild, wild2, wild3, wild4}},
	}

	scriptedSpinWeights = []float64{1, 99}

	maxWinChance = 0.08
	maxWins      = []comp.FakeSpin{
		{Indexes: util.Indexes{10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10}},
		{Indexes: util.Indexes{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9}},
		{Indexes: util.Indexes{8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8}},
	}

	scriptedMaxWinWeights = []float64{1, 999}

	script1Page1Grid = util.Indexes{10, 3, 3, 3, 10, 3, 3, 3, 10, 10, 10, 10, 10, 10, 10, 10, 10, 3, 3, 3, 10, 3, 3, 3}
	script1Page2Grid = util.Indexes{4, 10, 10, 4, 10, 10, 10, 10, 10, 4, 4, 10, 10, 4, 4, 10, 10, 10, 10, 10, 4, 10, 10, 4}
	script1Page3Grid = util.Indexes{10, 10, 10, 10, 10, 10, 10, 10, 10, 5, 10, 5, 10, 5, 10, 5, 10, 10, 10, 5, 10, 10, 10, 5}
	script6Grid      = util.Indexes{10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10, 10}

	script2SymbolWeights = util.AcquireWeighting().AddWeights(util.Indexes{8, 9, 10}, []float64{40, 33, 27})

	script1Weight = 1.2
	script2Weight = 1.2
	script3Weight = 1.2
	script4Weight = 1.2
	script5Weight = 1.2
	script6Weight = 1.2
	script9Weight = 10000.0
)

var (
	n01 = comp.WithName("Amethyst")
	n02 = comp.WithName("Topaz")
	n03 = comp.WithName("Emerald")
	n04 = comp.WithName("Citrine")
	n05 = comp.WithName("Ruby")
	n06 = comp.WithName("Grapes")
	n07 = comp.WithName("Cherries")
	n08 = comp.WithName("Gold")
	n09 = comp.WithName("Diamond")
	n10 = comp.WithName("Betic")
	n11 = comp.WithName("Free Spins")
	n12 = comp.WithName("Wild Star")
	n13 = comp.WithName("Wild Star x2")
	n14 = comp.WithName("Wild Star x3")
	n15 = comp.WithName("Wild Star x4")

	r01 = comp.WithResource("l5")
	r02 = comp.WithResource("l4")
	r03 = comp.WithResource("l3")
	r04 = comp.WithResource("l2")
	r05 = comp.WithResource("l1")
	r06 = comp.WithResource("h5")
	r07 = comp.WithResource("h4")
	r08 = comp.WithResource("h3")
	r09 = comp.WithResource("h2")
	r10 = comp.WithResource("h1")
	r11 = comp.WithResource("bonus")
	r12 = comp.WithResource("wild")
	r13 = comp.WithResource("wild2")
	r14 = comp.WithResource("wild3")
	r15 = comp.WithResource("wild4")

	p01 = comp.WithPayouts(0, 0, 0.2, 0.5, 0.8, 1.5)
	p02 = comp.WithPayouts(0, 0, 0.2, 0.5, 0.8, 1.5)
	p03 = comp.WithPayouts(0, 0, 0.2, 0.5, 0.8, 1.5)
	p04 = comp.WithPayouts(0, 0, 0.2, 0.8, 1, 2)
	p05 = comp.WithPayouts(0, 0, 0.2, 0.8, 1, 2)
	p06 = comp.WithPayouts(0, 0, 0.5, 1, 2, 3)
	p07 = comp.WithPayouts(0, 0, 0.5, 1, 2, 3)
	p08 = comp.WithPayouts(0, 0, 0.5, 1, 2, 5)
	p09 = comp.WithPayouts(0, 0, 1, 2, 3, 10)
	p10 = comp.WithPayouts(0, 0, 1, 3, 5, 15)
	p11 = comp.WithPayouts()
	p12 = p11
	p13 = p11
	p14 = p11
	p15 = p11

	ids       = [symbolCount]util.Index{id01, id02, id03, id04, id05, id06, id07, id08, id09, id10, id11, id12, id13, id14, id15}
	names     = [symbolCount]comp.SymbolOption{n01, n02, n03, n04, n05, n06, n07, n08, n09, n10, n11, n12, n13, n14, n15}
	resources = [symbolCount]comp.SymbolOption{r01, r02, r03, r04, r05, r06, r07, r08, r09, r10, r11, r12, r13, r14, r15}
	payouts   = [symbolCount]comp.SymbolOption{p01, p02, p03, p04, p05, p06, p07, p08, p09, p10, p11, p12, p13, p14, p15}

	freeSpinBands3   = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4}, freeSpinBandWeights3)
	freeSpinBands4   = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4}, freeSpinBandWeights4)
	freeSpinBands5   = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4}, freeSpinBandWeights5)
	freeSpinBands6   = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5}, freeSpinBandWeights6)
	freeSpinBands3BB = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3}, freeSpinBandWeights3BB)
	freeSpinBands4BB = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4}, freeSpinBandWeights4BB)
	freeSpinBands5BB = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4}, freeSpinBandWeights5BB)
	freeSpinBands6BB = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5}, freeSpinBandWeights6BB)

	freeSpinBand1 = util.AcquireWeighting().AddWeights(util.Indexes{1}, []float64{1})

	retriggerWeights3 = util.AcquireWeighting().AddWeights(util.Indexes{1, 0}, []float64{7, 93})

	flag0 = comp.NewRoundFlag(flagFreeSpinBands, "free spin band")
	flag1 = comp.NewRoundFlag(flagFreeSpinCount, "free spin count")
	flag2 = comp.NewRoundFlag(flagBonusBuy, "bonus buy type")
	flag3 = comp.NewRoundFlag(flagTriggerCount, "free spins trigger count")
	flag4 = comp.NewRoundFlag(flagScriptedRound, "scripted round ID")
	flag5 = comp.NewRoundFlag(flagRetriggerCount, "free spins retrigger count")
	flags = comp.RoundFlags{flag0, flag1, flag2, flag3, flag4, flag5}

	actions92All     comp.SpinActions
	actions92First   comp.SpinActions
	actions92Free    comp.SpinActions
	actions92FirstBB comp.SpinActions
	actions92FreeBB  comp.SpinActions
	actions94All     comp.SpinActions
	actions94First   comp.SpinActions
	actions94Free    comp.SpinActions
	actions94FirstBB comp.SpinActions
	actions94FreeBB  comp.SpinActions
	actions96All     comp.SpinActions
	actions96First   comp.SpinActions
	actions96Free    comp.SpinActions
	actions96FirstBB comp.SpinActions
	actions96FreeBB  comp.SpinActions
	actionsAllB      comp.SpinActions

	scriptedRounds *comp.ScriptedRoundSelector
)

var (
	symbols *comp.SymbolSet

	freeCounter *comp.RoundFlagAction

	slots92 *comp.Slots
	slots94 *comp.Slots
	slots96 *comp.Slots

	slots92params game.RegularParams
	slots94params game.RegularParams
	slots96params game.RegularParams
)

func initActions() {
	actions92All = make(comp.SpinActions, 0, 32)
	actions94All = make(comp.SpinActions, 0, 32)
	actions96All = make(comp.SpinActions, 0, 32)

	// bonus buy feature.
	bonusBuy1 := comp.NewPaidAction(comp.FreeSpins, 0, bonusBuyCost, scatter, 3).WithFlag(flagBonusBuy, bonusKind).WithBonusKind(bonusKind)
	bonusBuy1.Describe(bonusBuyID, "bonus buy")

	bonusExtraScatter := comp.NewGenerateSymbolAction(scatter, []float64{37.35, 14, 1}).AllowPrevious().GenerateNoDupes()
	bonusExtraScatter.WithTriggerFilters(comp.OnRoundFlagAbove(flagBonusBuy, 0), comp.OnRoundFlagValue(flagScriptedRound, 0))
	bonusExtraScatter.Describe(bonusExtraID, "generate 4th/5th/6th scatter - bonus buy")

	// generate wilds RTP 92.
	firstWild92 := comp.NewGenerateSymbolAction(wild, firstWildChances92)
	firstWild92.Describe(firstWild92ID, "generate wilds - first spin - RTP 92")

	actions92All = append(actions92All, firstWild92)

	// generate scatters RTP 92.
	firstScatter92 := comp.NewGenerateSymbolAction(scatter, scatterRegularChances92).GenerateNoDupes()
	firstScatter92.Describe(firstScatter92ID, "generate scatters - first spin - RTP 92")

	freeScatter92a := buildFreeScatterAction(92, freeScatter92aID, scatterFreeSpinChances92a, modBonus3S, "3 scatters")
	freeScatter92a.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 3))
	freeScatter92a.Describe(freeScatter92aID+5, "multi-select generate scatters - RTP 92 - 3 scatters")

	freeScatter92b := buildFreeScatterAction(92, freeScatter92bID, scatterFreeSpinChances92b, modBonus4S, "4 scatters")
	freeScatter92b.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 4))
	freeScatter92b.Describe(freeScatter92bID+5, "multi-select generate scatters - RTP 92 - 4 scatters")

	freeScatter92c := buildFreeScatterAction(92, freeScatter92cID, scatterFreeSpinChances92c, modBonus5S, "5 scatters")
	freeScatter92c.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 5))
	freeScatter92c.Describe(freeScatter92cID+5, "multi-select generate scatters - RTP 92 - 5 scatters")

	freeScatter92d := buildFreeScatterAction(92, freeScatter92dID, scatterFreeSpinChances92d, modBonus6S, "6 scatters")
	freeScatter92d.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 6))
	freeScatter92d.Describe(freeScatter92dID+5, "multi-select generate scatters - RTP 92 - 6 scatters")

	freeScatter92e := buildFreeScatterAction(92, freeScatter92eID, scatterFreeSpinChances92e, modBonus33S, "3+3 scatters")
	freeScatter92e.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 13))
	freeScatter92e.Describe(freeScatter92eID+5, "multi-select generate scatters - RTP 92 - 3+3 scatters")

	actions92All = append(actions92All, firstScatter92, freeScatter92a, freeScatter92b, freeScatter92c, freeScatter92d, freeScatter92e)

	// generate wilds RTP 94.
	firstWild94 := comp.NewGenerateSymbolAction(wild, firstWildChances94)
	firstWild94.Describe(firstWild94ID, "generate wilds - first spin - RTP 94")

	actions94All = append(actions94All, firstWild94)

	// generate scatters RTP 94.
	firstScatter94 := comp.NewGenerateSymbolAction(scatter, scatterRegularChances94).GenerateNoDupes()
	firstScatter94.Describe(firstScatter94ID, "generate scatters - first spin - RTP 94")

	freeScatter94a := buildFreeScatterAction(94, freeScatter94aID, scatterFreeSpinChances94a, modBonus3S, "3 scatters")
	freeScatter94a.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 3))
	freeScatter94a.Describe(freeScatter94aID+5, "multi-select generate scatters - RTP 94 - 3 scatters")

	freeScatter94b := buildFreeScatterAction(94, freeScatter94bID, scatterFreeSpinChances94b, modBonus4S, "4 scatters")
	freeScatter94b.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 4))
	freeScatter94b.Describe(freeScatter94bID+5, "multi-select generate scatters - RTP 94 - 4 scatters")

	freeScatter94c := buildFreeScatterAction(94, freeScatter94cID, scatterFreeSpinChances94c, modBonus5S, "5 scatters")
	freeScatter94c.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 5))
	freeScatter94c.Describe(freeScatter94cID+5, "multi-select generate scatters - RTP 94 - 5 scatters")

	freeScatter94d := buildFreeScatterAction(94, freeScatter94dID, scatterFreeSpinChances94d, modBonus6S, "6 scatters")
	freeScatter94d.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 6))
	freeScatter94d.Describe(freeScatter94dID+5, "multi-select generate scatters - RTP 94 - 6 scatters")

	freeScatter94e := buildFreeScatterAction(94, freeScatter94eID, scatterFreeSpinChances94e, modBonus33S, "3+3 scatters")
	freeScatter94e.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 13))
	freeScatter94e.Describe(freeScatter94eID+5, "multi-select generate scatters - RTP 94 - 3+3 scatters")

	actions94All = append(actions94All, firstScatter94, freeScatter94a, freeScatter94b, freeScatter94c, freeScatter94d, freeScatter94e)

	// generate wilds RTP 96.
	firstWild96 := comp.NewGenerateSymbolAction(wild, firstWildChances96)
	firstWild96.Describe(firstWild96ID, "generate wilds - first spin - RTP 96")

	actions96All = append(actions96All, firstWild96)

	// generate scatters RTP 96.
	firstScatter96 := comp.NewGenerateSymbolAction(scatter, scatterRegularChances96).GenerateNoDupes()
	firstScatter96.Describe(firstScatter96ID, "generate scatters - first spin - RTP 96")

	freeScatter96a := buildFreeScatterAction(96, freeScatter96aID, scatterFreeSpinChances96a, modBonus3S, "3 scatters")
	freeScatter96a.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 3))
	freeScatter96a.Describe(freeScatter96aID+5, "multi-select generate scatters - RTP 96 - 3 scatters")

	freeScatter96b := buildFreeScatterAction(96, freeScatter96bID, scatterFreeSpinChances96b, modBonus4S, "4 scatters")
	freeScatter96b.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 4))
	freeScatter96b.Describe(freeScatter96bID+5, "multi-select generate scatters - RTP 96 - 4 scatters")

	freeScatter96c := buildFreeScatterAction(96, freeScatter96cID, scatterFreeSpinChances96c, modBonus5S, "5 scatters")
	freeScatter96c.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 5))
	freeScatter96c.Describe(freeScatter96cID+5, "multi-select generate scatters - RTP 96 - 5 scatters")

	freeScatter96d := buildFreeScatterAction(96, freeScatter96dID, scatterFreeSpinChances96d, modBonus6S, "6 scatters")
	freeScatter96d.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 6))
	freeScatter96d.Describe(freeScatter96dID+5, "multi-select generate scatters - RTP 96 - 6 scatters")

	freeScatter96e := buildFreeScatterAction(96, freeScatter96eID, scatterFreeSpinChances96e, modBonus33S, "3+3 scatters")
	freeScatter96e.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 13))
	freeScatter96e.Describe(freeScatter96eID+5, "multi-select generate scatters - RTP 96 - 3+3 scatters")

	actions96All = append(actions96All, firstScatter96, freeScatter96a, freeScatter96b, freeScatter96c, freeScatter96d, freeScatter96e)

	// generate 2x/3x/4x wilds on first spin.
	firstWild234 := comp.NewGenerateSymbolsAction(wildWeights, firstWildChances)
	firstWild234.Describe(firstWild2ID, "generate any wild - first spin")

	// generate 2x/3x/4x wilds on free spins.
	freeWild234a := buildWildsAction(freeWild234aID, freeWildChancesA, modBonus3W, "free spins (3 scatters)", freeWildReels...)
	freeWild234a.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 3))
	freeWild234a.Describe(freeWild234aID+5, "multi-select generate any wild - free spins (3 scatters)")

	freeWild234b := buildWildsAction(freeWild234bID, freeWildChancesB, modBonus4W, "free spins (4 scatters)", freeWildReels...)
	freeWild234b.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 4))
	freeWild234b.Describe(freeWild234bID+5, "multi-select generate any wild - free spins (4 scatters)")

	freeWild234c := buildWildsAction(freeWild234cID, freeWildChancesC, modBonus5W, "free spins (5 scatters)", freeWildReels...)
	freeWild234c.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 5))
	freeWild234c.Describe(freeWild234cID+5, "multi-select generate any wild - free spins (5 scatters)")

	freeWild234d := buildWildsAction(freeWild234dID, freeWildChancesD, modBonus6W, "free spins (6 scatters)", freeWildReels...)
	freeWild234d.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 6))
	freeWild234d.Describe(freeWild234dID+5, "multi-select generate any wild - free spins (6 scatters)")

	freeWild234e := buildWildsAction(freeWild234eID, freeWildChancesE, modBonus33W, "free spins (3+3 scatters)", freeWildReels...)
	freeWild234e.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 13))
	freeWild234e.Describe(freeWild234eID+5, "multi-select generate any wild - free spins (3+3 scatters)")

	// all paylines.
	paylines := comp.NewAllPaylinesAction(true)
	paylines.Describe(paylinesID, "all paylines")

	// reduce payouts by band.
	reducePayouts1 := comp.NewRemovePayoutBandsAction(3, direction, true, true, reduceBandsFirst)
	reducePayouts1.Describe(removePayouts1ID, "remove payouts by band (first spins)")

	reducePayouts2 := comp.NewRemovePayoutBandsAction(3, direction, false, true, reduceBandsFree1)
	reducePayouts2.WithTriggerFilters(comp.OnRoundFlagBelow(flagFreeSpinCount, 7))
	reducePayouts2.Describe(removePayouts2ID, "remove payouts by band (first few free spins)")

	reducePayouts3 := comp.NewRemovePayoutBandsAction(3, direction, false, true, reduceBandsFree2)
	reducePayouts3.WithTriggerFilters(comp.OnRoundFlagAbove(flagFreeSpinCount, 6))
	reducePayouts3.Describe(removePayouts3ID, "remove payouts by band (free spins)")

	// award free spins.
	award10 := comp.NewScatterFreeSpinsAction(10, false, scatter, 3, false)
	award10.Describe(award10ID, "award 10 free spins - first spin - 3 scatters")
	award15 := comp.NewScatterFreeSpinsAction(15, false, scatter, 4, false).WithAlternate(award10)
	award15.Describe(award15ID, "award 15 free spins - first spin - 4 scatters")
	award20 := comp.NewScatterFreeSpinsAction(20, false, scatter, 5, false).WithAlternate(award15)
	award20.Describe(award20ID, "award 20 free spins - first spin - 5 scatters")
	award50 := comp.NewScatterFreeSpinsAction(50, false, scatter, 6, false).WithAlternate(award20)
	award50.Describe(award50ID, "award 50 free spins - first spin - 6 scatters")

	// count scatters in first spin.
	countScatters := comp.NewRoundFlagSymbolCountAction(flagTriggerCount, scatter) // important values are 3/4/5/6.
	countScatters.Describe(countScattersID, "count scatters into flag 3")

	// randomly select zero or more free spin retriggers.
	retriggerScatters3 := comp.NewRoundFlagWeightedAction(flagRetriggerCount, retriggerWeights3)
	retriggerScatters3.WithStage(comp.AwardBonuses)
	retriggerScatters3.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 3), comp.OnSpinSequence(1))
	retriggerScatters3.Describe(retriggerScatters3ID, "award retriggers into flag 5 - 3 scatters")

	// mark awarded free spin retriggers in flag 3.
	markRetriggers3 := comp.NewRoundFlagSet(flagTriggerCount, 13)
	markRetriggers3.WithStage(comp.AwardBonuses)
	markRetriggers3.WithTriggerFilters(comp.OnRoundFlagValue(flagRetriggerCount, 1), comp.OnRoundFlagValue(flagTriggerCount, 3))
	markRetriggers3.Describe(markRetriggers3ID, "mark retriggers into flag 3 - 3 scatters")

	// award retriggers.
	awardRetrigger3 := comp.NewGenerateSymbolAction(scatter, []float64{90, 85, 95}).GenerateNoDupes()
	awardRetrigger3.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 13), comp.OnRoundFlagAbove(flagRetriggerCount, 0), comp.OnSpinSequenceAbove(7))
	awardRetrigger3.Describe(awardRetrigger3ID, "award retriggers - spin 7+ - 3+3 scatters")

	// force retriggers.
	forceRetrigger3 := comp.NewGenerateSymbolAction(scatter, []float64{100, 100, 100}).GenerateNoDupes()
	forceRetrigger3.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 13), comp.OnRoundFlagAbove(flagRetriggerCount, 0), comp.OnSpinSequence(11))
	forceRetrigger3.Describe(forceRetrigger3ID, "force retriggers - spin 11 - 3+3 scatters")

	// reduce retriggers.
	reduceRetriggers := comp.NewRoundFlagDecreaseAction(flagRetriggerCount)
	reduceRetriggers.WithStage(comp.AwardBonuses)
	reduceRetriggers.WithTriggerFilters(comp.OnRoundFlagAbove(flagRetriggerCount, 0), comp.OnGridCounts(scatter, []int{3, 4, 5, 6}))
	reduceRetriggers.Describe(reduceRetriggersID, "reduce retrigger count")

	// award free spins from retriggers.
	retrigger3 := comp.NewScatterFreeSpinsAction(10, false, scatter, 3, false)
	retrigger3.Describe(retrigger3ID, "award 10 free spins - free spin - 3 scatters")
	retrigger4 := comp.NewScatterFreeSpinsAction(15, false, scatter, 4, false).WithAlternate(retrigger3)
	retrigger4.Describe(retrigger4ID, "award 15 free spins - free spin - 4 scatters")
	retrigger5 := comp.NewScatterFreeSpinsAction(20, false, scatter, 5, false).WithAlternate(retrigger4)
	retrigger5.Describe(retrigger5ID, "award 20 free spins - free spin - 5 scatters")
	retrigger6 := comp.NewScatterFreeSpinsAction(50, false, scatter, 6, false).WithAlternate(retrigger5)
	retrigger6.Describe(retrigger6ID, "award 50 free spins - free spin - 6 scatters")

	// mark sticky wilds.
	stickies := comp.NewStickySymbolsAction(wild, wild2, wild3, wild4)
	stickies.Describe(stickies2ID, "make wild symbols sticky (free spins)")

	// round flag 0 for payout bands (base game).
	payoutBands3 := comp.NewRoundFlagWeightedAction(flagFreeSpinBands, freeSpinBands3)
	payoutBands3.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 3))
	payoutBands3.Describe(payoutBands3ID, "weighted payout band - 3 scatters (flag 0)")
	payoutBands4 := comp.NewRoundFlagWeightedAction(flagFreeSpinBands, freeSpinBands4)
	payoutBands4.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 4))
	payoutBands4.Describe(payoutBands4ID, "weighted payout band - 4 scatters (flag 0)")
	payoutBands5 := comp.NewRoundFlagWeightedAction(flagFreeSpinBands, freeSpinBands5)
	payoutBands5.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 5))
	payoutBands5.Describe(payoutBands5ID, "weighted payout band - 5 scatters (flag 0)")
	payoutBands6 := comp.NewRoundFlagWeightedAction(flagFreeSpinBands, freeSpinBands6)
	payoutBands6.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 6))
	payoutBands6.Describe(payoutBands6ID, "weighted payout band - 6 scatters (flag 0)")

	payoutBands33 := comp.NewRoundFlagWeightedAction(flagFreeSpinBands, freeSpinBands3)
	payoutBands33.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 13))
	payoutBands33.Describe(payoutBands33ID, "weighted payout band - 3+3 scatters (flag 0)")

	// round flag 0 for payout bands (base game).
	payoutBands3BB := comp.NewRoundFlagWeightedAction(flagFreeSpinBands, freeSpinBands3BB)
	payoutBands3BB.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 3))
	payoutBands3BB.Describe(payoutBands3BBID, "weighted payout band - BB - 3 scatters (flag 0)")
	payoutBands4BB := comp.NewRoundFlagWeightedAction(flagFreeSpinBands, freeSpinBands4BB)
	payoutBands4BB.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 4))
	payoutBands4BB.Describe(payoutBands4BBID, "weighted payout band - BB - 4 scatters (flag 0)")
	payoutBands5BB := comp.NewRoundFlagWeightedAction(flagFreeSpinBands, freeSpinBands5BB)
	payoutBands5BB.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 5))
	payoutBands5BB.Describe(payoutBands5BBID, "weighted payout band - BB - 5 scatters (flag 0)")
	payoutBands6BB := comp.NewRoundFlagWeightedAction(flagFreeSpinBands, freeSpinBands6BB)
	payoutBands6BB.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 6))
	payoutBands6BB.Describe(payoutBands6BBID, "weighted payout band - BB - 6 scatters (flag 0)")

	payoutBands33BB := comp.NewRoundFlagWeightedAction(flagFreeSpinBands, freeSpinBands3BB)
	payoutBands33BB.WithTriggerFilters(comp.OnRoundFlagValue(flagTriggerCount, 13))
	payoutBands33BB.Describe(payoutBands33BBID, "weighted payout band - BB - 3+3 scatters (flag 0)")

	// update round flag 3 marking sequence of free spin.
	freeCounter = comp.NewRoundFlagIncreaseAction(flagFreeSpinCount)
	freeCounter.Describe(freeCounterID, "count number of free spins (flag 1)")

	// grid teasers.
	teaserFP := comp.NewCustomAction(comp.AwardBonuses, comp.GridModified, fakeFullPage, nil)
	teaserFP.WithTriggerFilters(comp.OnFirstSpin, comp.OnZeroPayouts(), comp.OnGridCount(scatter, 0))
	teaserFP.Describe(teaserID, "teaser almost full page")

	// forced 120x wins.
	forceWin120x := comp.NewFakeSpinAction(win120SpinChance, wins120x...)
	forceWin120x.WithStage(comp.ReviseGrid)
	forceWin120x.Describe(win100xID, "forced 100x win - Betic")

	// forced max wins.
	forcedMaxWin := comp.NewFakeSpinAction(maxWinChance, maxWins...)
	forcedMaxWin.WithStage(comp.ReviseGrid)
	forcedMaxWin.Describe(maxWinID, "forced max win - Betic")

	// dummy script; does absolutely nothing.
	dummyScript := comp.NewFakeSpinAction(0)
	dummyScript.WithStage(comp.ReviseGrid)
	dummyScript.Describe(dummyScriptID, "dummy script")

	// scripted spins.
	scriptedSpin := comp.NewMultiActionWeighted(util.AcquireWeighting().AddWeights(util.Indexes{1, 2}, scriptedSpinWeights), forceWin120x, dummyScript)
	scriptedSpin.WithTriggerFilters(comp.OnFirstSpin, comp.OnGridNoWilds(), comp.OnGridCount(scatter, 0), comp.OnGridCount(10, 0))
	scriptedSpin.Describe(scriptedSpinID, "force scripted spin")

	scriptedMaxWin := comp.NewMultiActionWeighted(util.AcquireWeighting().AddWeights(util.Indexes{1, 2}, scriptedMaxWinWeights), forcedMaxWin, dummyScript)
	scriptedMaxWin.WithTriggerFilters(comp.OnFirstSpin, comp.OnGridNoWilds(), comp.OnGridCount(scatter, 0), comp.OnGridCount(10, 0))
	scriptedMaxWin.Describe(scriptedMaxWinID, "force scripted max win")

	// remove scatters for combined bonus buy + scripted round.
	forceNoScatters := comp.NewCustomAction(comp.ReviseGrid, comp.GridModified, removeScatters, nil)
	forceNoScatters.Describe(forceNoScattersID, "remove scatters - bonus buy (scripted round)")

	// force 3-6 scatters for scripted rounds.
	forceScatters3to6 := comp.NewGenerateSymbolAction(scatter, []float64{100, 100, 100, 45, 20, 10}).GenerateNoDupes()
	forceScatters3to6.Describe(forceScrolls4ID, "force 3-6 scatters (scripted round)")

	// force payout band for scripted rounds.
	forceBand1 := comp.NewRoundFlagWeightedAction(flagFreeSpinBands, freeSpinBand1)
	forceBand1.Describe(forceBand5ID, "force payout band 5 (flag 0; scripted round)")

	// force random 4x wilds for scripted rounds.
	forceWildsX4 := comp.NewGenerateSymbolAction(wild4, []float64{100, 100, 99.98, 99.98, 99.98, 40, 10}, freeWildReels...).AllowPrevious()
	forceWildsX4.WithTriggerFilters(comp.OnSpinSequence(11))
	forceWildsX4.Describe(forceWildsX4ID, "force 4x wilds - spin 11 (scripted round)")

	// force 3-6 scatters in row 3 - script 1.
	script1ScatterR3 := comp.NewCustomAction(comp.ReviseGrid, comp.GridModified, forceRow3Scatters, nil)
	script1ScatterR3.Describe(script1ScattersR3ID, "force 3-6 scatters row 3 (script 1)")

	// forced pages - script 1.
	script1Page1 := comp.NewCustomAction(comp.ReviseGrid, comp.GridModified, replaceGridFunction(script1Page1Grid), nil)
	script1Page1.WithTriggerFilters(comp.OnSpinSequence(2))
	script1Page1.Describe(script1Page1ID, "force page (script 1, free spin 1)")

	script1Page2 := comp.NewCustomAction(comp.ReviseGrid, comp.GridModified, replaceGridFunction(script1Page2Grid), nil)
	script1Page2.WithTriggerFilters(comp.OnSpinSequence(3))
	script1Page2.Describe(script1Page2ID, "force page (script 1, free spin 2)")

	script1Page3 := comp.NewCustomAction(comp.ReviseGrid, comp.GridModified, replaceGridFunction(script1Page3Grid), nil)
	script1Page3.WithTriggerFilters(comp.OnSpinSequence(4))
	script1Page3.Describe(script1Page3ID, "force page (script 1, free spin 3)")

	// forced full page - script 2.
	script2ForceFP := comp.NewCustomAction(comp.ReviseGrid, comp.GridModified, script2FullPage, nil)
	script2ForceFP.WithTriggerFilters(comp.OnSpinSequenceAbove(5))
	script2ForceFP.Describe(script2ForceFPID, "force full page (script 2)")

	// random 4x wilds - script 3.
	script3WildsX4 := comp.NewGenerateSymbolAction(wild4, []float64{72, 45, 30, 10}, freeWildReels...).AllowPrevious()
	script3WildsX4.Describe(script3WildsX4ID, "random 4x wilds (script 3)")

	// random multiplier wilds - script 4.
	weightsScrip4WildsX := util.AcquireWeighting().AddWeights(util.Indexes{wild2, wild3, wild4}, []float64{45, 30, 25})
	script4WildsX := comp.NewGenerateSymbolsAction(weightsScrip4WildsX, []float64{88, 45, 35, 10}, freeWildReels...).AllowPrevious()
	script4WildsX.Describe(script4WildsXID, "random multiplier wilds (script 4)")

	// random wilds - script 5.
	weightsScrip5WildsX := util.AcquireWeighting().AddWeights(util.Indexes{wild, wild2, wild3, wild4}, []float64{70, 15, 10, 5})
	script5WildsX := comp.NewGenerateSymbolsAction(weightsScrip5WildsX, []float64{17, 10, 1}, freeWildReels...).AllowPrevious()
	script5WildsX.Describe(script5WildsXID, "random multiplier wilds (script 5)")

	// random scatters - script 5.
	script5Scatters := comp.NewGenerateSymbolAction(scatter, []float64{60, 45, 25, 40, 8, 0.5}).GenerateNoDupes()
	script5Scatters.WithTriggerFilters(comp.OnBetweenFreeSpins(6, 99))
	script5Scatters.Describe(script5ScattersID, "random scatters (script 5)")

	// random retrigger - script 5.
	script5Retrigger1 := comp.NewGenerateSymbolAction(scatter, []float64{90, 90, 90, 50, 10, 0.5}).GenerateNoDupes().AllowPrevious()
	script5Retrigger1.WithTriggerFilters(comp.OnBetweenFreeSpins(1, 5))
	script5Retrigger1.Describe(script5Retrigger1ID, "attempt retrigger (script 5)")

	// force retrigger - script 5.
	script5Retrigger2 := comp.NewGenerateSymbolAction(scatter, []float64{100, 100, 100, 70, 20, 0.5}).GenerateNoDupes().AllowPrevious()
	script5Retrigger2.WithTriggerFilters(comp.OnRemainingFreeSpins(0))
	script5Retrigger2.Describe(script5Retrigger2ID, "force retrigger (script 5)")

	// forced full page - script 6.
	script6ForceFP := comp.NewCustomAction(comp.ReviseGrid, comp.GridModified, replaceGridFunction(script6Grid), nil)
	script6ForceFP.Describe(script6ForceFPID, "force full page (script 6)")

	// scripted rounds - scenario 1 configuration (max payout).
	script1 := comp.NewScriptedRound(1, script1Weight,
		comp.SpinActions{paylines, award50, countScatters, forceBand1},
		comp.SpinActions{freeCounter, script1Page1, script1Page2, script1Page3, paylines},
	).WithBonusBuys(bonusKind)

	// scripted rounds - scenario 2 configuration (max payout).
	script2 := comp.NewScriptedRound(2, script2Weight,
		comp.SpinActions{forceScatters3to6, paylines, award50, countScatters, forceBand1},
		comp.SpinActions{freeCounter, script2ForceFP, paylines},
	).WithBonusBuys(bonusKind)

	// scripted rounds - scenario 3 configuration (max payout).
	script3 := comp.NewScriptedRound(3, script3Weight,
		comp.SpinActions{forceScatters3to6, paylines, award50, countScatters, forceBand1},
		comp.SpinActions{freeCounter, script3WildsX4, forceWildsX4, paylines, stickies},
	).WithBonusBuys(bonusKind)

	// scripted rounds - scenario 4 configuration (max payout).
	script4 := comp.NewScriptedRound(4, script4Weight,
		comp.SpinActions{forceScatters3to6, paylines, award50, countScatters, forceBand1},
		comp.SpinActions{freeCounter, script4WildsX, forceWildsX4, paylines, stickies},
	).WithBonusBuys(bonusKind)

	// scripted rounds - scenario 5 configuration (max payout).
	script5 := comp.NewScriptedRound(5, script5Weight,
		comp.SpinActions{forceScatters3to6, paylines, award50, countScatters, forceBand1},
		comp.SpinActions{freeCounter, script5WildsX, script5Scatters, script5Retrigger1, script5Retrigger2, paylines, retrigger6, stickies},
	).WithBonusBuys(bonusKind)

	// scripted rounds - scenario 6 configuration (max payout - single spin).
	script6 := comp.NewScriptedRound(6, script6Weight,
		comp.SpinActions{script6ForceFP, paylines, stickies},
		nil,
	)

	// scripted rounds - scenario 9 configuration (dummy to fill up the weighting).
	script9 := comp.NewScriptedRound(9, script9Weight, nil, nil)

	// scripted rounds accumulated.
	scriptedRounds = comp.NewScriptedRoundSelector(0.1, script1, script2, script3, script4, script5, script6, script9)
	scriptedRounds.WithSpinFlag(flagScriptedRound).WithBonusBuys(bonusKind).WithBonusChances(flagBonusBuy, 21.1)

	forceAll := comp.SpinActions{
		script1ScatterR3, forceScatters3to6, forceBand1, forceWildsX4,
		script1Page1, script1Page2, script1Page3,
		script2ForceFP, script3WildsX4, script4WildsX,
		script5WildsX, script5Scatters, script5Retrigger1, script5Retrigger2,
		script6ForceFP,
	}

	// gather it all together.
	actions92All = append(actions92All, bonusBuy1, bonusExtraScatter, freeCounter, awardRetrigger3, forceRetrigger3)
	actions92FirstA := comp.SpinActions{firstWild92, firstScatter92}
	actions92FreeA := comp.SpinActions{freeCounter, awardRetrigger3, forceRetrigger3, freeScatter92a, freeScatter92b, freeScatter92c, freeScatter92d, freeScatter92e}
	actions92FirstBonusA := comp.SpinActions{bonusBuy1, bonusExtraScatter}
	actions92FreeBonusA := comp.SpinActions{freeCounter, awardRetrigger3, forceRetrigger3, freeScatter92a, freeScatter92b, freeScatter92c, freeScatter92d, freeScatter92e}

	actions94All = append(actions94All, bonusBuy1, bonusExtraScatter, freeCounter, awardRetrigger3, forceRetrigger3)
	actions94FirstA := comp.SpinActions{firstWild94, firstScatter94}
	actions94FreeA := comp.SpinActions{freeCounter, awardRetrigger3, forceRetrigger3, freeScatter94a, freeScatter94b, freeScatter94c, freeScatter94d, freeScatter94e}
	actions94FirstBonusA := comp.SpinActions{bonusBuy1, bonusExtraScatter}
	actions94FreeBonusA := comp.SpinActions{freeCounter, awardRetrigger3, forceRetrigger3, freeScatter94a, freeScatter94b, freeScatter94c, freeScatter94d, freeScatter94e}

	actions96All = append(actions96All, bonusBuy1, bonusExtraScatter, freeCounter, awardRetrigger3, forceRetrigger3)
	actions96FirstA := comp.SpinActions{firstWild96, firstScatter96}
	actions96FreeA := comp.SpinActions{freeCounter, awardRetrigger3, forceRetrigger3, freeScatter96a, freeScatter96b, freeScatter96c, freeScatter96d, freeScatter96e}
	actions96FirstBonusA := comp.SpinActions{bonusBuy1, bonusExtraScatter}
	actions96FreeBonusA := comp.SpinActions{freeCounter, awardRetrigger3, forceRetrigger3, freeScatter96a, freeScatter96b, freeScatter96c, freeScatter96d, freeScatter96e}

	actionsAllB = append(actionsAllB,
		payoutBands3, payoutBands4, payoutBands5, payoutBands6, payoutBands33,
		payoutBands3BB, payoutBands4BB, payoutBands5BB, payoutBands6BB, payoutBands33BB,
		firstWild234, freeWild234a, freeWild234b, freeWild234c, freeWild234d, freeWild234e,
		forceWin120x, scriptedSpin, forcedMaxWin, scriptedMaxWin,
		paylines, reducePayouts1, reducePayouts2, reducePayouts3, teaserFP,
		award50, countScatters, retriggerScatters3, markRetriggers3, reduceRetriggers, retrigger6, stickies,
	)
	actionsFirstB := comp.SpinActions{
		payoutBands3, payoutBands4, payoutBands5, payoutBands6, payoutBands33,
		firstWild234, scriptedSpin, scriptedMaxWin, paylines, reducePayouts1, award50,
		countScatters, retriggerScatters3, markRetriggers3, teaserFP,
	}
	actionsFreeB := comp.SpinActions{
		freeWild234a, freeWild234b, freeWild234c, freeWild234d, freeWild234e, paylines,
		reduceRetriggers, reducePayouts2, reducePayouts3, retrigger6, stickies,
	}
	actionsFirstBonusB := comp.SpinActions{
		payoutBands3BB, payoutBands4BB, payoutBands5BB, payoutBands6BB, payoutBands33BB,
		paylines, reducePayouts1, award50,
		countScatters, retriggerScatters3, markRetriggers3,
	}
	actionsFreeBonusB := comp.SpinActions{
		freeWild234a, freeWild234b, freeWild234c, freeWild234d, freeWild234e, paylines,
		reduceRetriggers, reducePayouts2, reducePayouts3, retrigger6, stickies,
	}

	actions92All = append(append(actions92All, actionsAllB...), forceAll...)
	actions92First = append(actions92FirstA, actionsFirstB...)
	actions92Free = append(actions92FreeA, actionsFreeB...)
	actions92FirstBB = append(actions92FirstBonusA, actionsFirstBonusB...)
	actions92FreeBB = append(actions92FreeBonusA, actionsFreeBonusB...)

	actions94All = append(append(actions94All, actionsAllB...), forceAll...)
	actions94First = append(actions94FirstA, actionsFirstB...)
	actions94Free = append(actions94FreeA, actionsFreeB...)
	actions94FirstBB = append(actions94FirstBonusA, actionsFirstBonusB...)
	actions94FreeBB = append(actions94FreeBonusA, actionsFreeBonusB...)

	actions96All = append(append(actions96All, actionsAllB...), forceAll...)
	actions96First = append(actions96FirstA, actionsFirstB...)
	actions96Free = append(actions96FreeA, actionsFreeB...)
	actions96FirstBB = append(actions96FirstBonusA, actionsFirstBonusB...)
	actions96FreeBB = append(actions96FreeBonusA, actionsFreeBonusB...)
}

func buildFreeScatterAction(rtp, id int, chances [5][]float64, mod comp.ChanceModifier, suffix string) *comp.MultiActionFlagValue {
	if suffix != "" {
		suffix = " - " + suffix
	}

	var list [5]*comp.ReviseAction
	for i := range list {
		a := comp.NewGenerateSymbolAction(scatter, chances[i]).GenerateNoDupes()
		a.WithChanceModifier(mod)
		a.Describe(id+i, fmt.Sprintf("generate scatters - free spin band %d - RTP %d%s", i+1, rtp, suffix))
		list[i] = a
	}

	a := comp.NewMultiActionFlagValue(flagFreeSpinBands, 1, list[0], 2, list[1], 3, list[2], 4, list[3], 5, list[4])

	switch rtp {
	case 92:
		actions92All = append(actions92All, list[0], list[1], list[2], list[3], list[4], a)
	case 94:
		actions94All = append(actions94All, list[0], list[1], list[2], list[3], list[4], a)
	case 96:
		actions96All = append(actions96All, list[0], list[1], list[2], list[3], list[4], a)
	}

	return a
}

func buildWildsAction(id int, chances [5][]float64, mod comp.ChanceModifier, suffix string, reels ...uint8) *comp.MultiActionFlagValue {
	if suffix != "" {
		suffix = " - " + suffix
	}

	var list [5]*comp.ReviseAction
	for i := range list {
		a := comp.NewGenerateSymbolsAction(wildWeights, chances[i], reels...)
		a.WithChanceModifier(mod)
		a.Describe(id+i, fmt.Sprintf("generate any wild - free spin band %d%s", i+1, suffix))
		list[i] = a
	}

	a := comp.NewMultiActionFlagValue(flagFreeSpinBands, 1, list[0], 2, list[1], 3, list[2], 4, list[3], 5, list[4])
	actionsAllB = append(actionsAllB, list[0], list[1], list[2], list[3], list[4], a)
	return a
}

func fakeFullPage(spin *comp.Spin) bool {
	if !spin.TestChance2(fakeFullPageChance) {
		return false
	}

	prng, indexes := spin.PRNG(), spin.Indexes()

	symbol := fakeFullPageWeights.RandomIndex(prng)
	exclude2 := spin.TestChance2(55)
	count := 12 + prng.IntN(9)

	var reel *comp.Reel
	var from, to int

	if exclude2 {
		reel = spin.Reels()[1]
		from, to = 4, 8
	} else {
		reel = spin.Reels()[2]
		from, to = 8, 12
	}

	for ix := from; ix < to; ix++ {
		id := indexes[ix]
		deadlock := 11

		for deadlock > 0 && (id == symbol || id == wild || id == wild2 || id == wild3 || id == wild4) {
			id = reel.RandomIndex(prng)
			indexes[ix] = id
			deadlock--
		}

		if deadlock == 0 {
			return false
		}
	}

	count -= int(spin.CountSymbol(symbol))

	l := len(indexes)
	deadlock := 11

	for deadlock > 0 && count > 0 {
		if ix := prng.IntN(l); ix < from || ix >= to {
			if id := indexes[ix]; id != symbol {
				indexes[ix] = symbol
				count--
			} else {
				deadlock--
			}
		}
	}

	return true
}

func removeScatters(spin *comp.Spin) bool {
	reelCount, rowCount := spin.GridSize()
	indexes, gen, prng := spin.Indexes(), spin.Reels(), spin.PRNG()

	var found bool

	reel1 := make(util.Indexes, 4)
	copy(reel1, indexes)

	newID := func(reel int) util.Index {
		for {
			id := gen[reel-1].RandomIndex(prng)
			if reel != 2 || !reel1.Contains(id) {
				return id
			}
		}
	}

	for reel := 1; reel <= reelCount; reel++ {
		m := reel * rowCount
		for offs := m - rowCount; offs < m; offs++ {
			if indexes[offs] == scatter {
				indexes[offs] = newID(reel)
			}
		}
	}

	return found
}

func forceRow3Scatters(spin *comp.Spin) bool {
	prng, indexes := spin.PRNG(), spin.Indexes()
	remove := 3 - prng.IntN(4)

	offsets := []int{2, 6, 10, 14, 18, 22}
	for remove > 0 {
		i := prng.IntN(len(offsets))
		offsets[i] = 0
		remove--
	}

	for _, offs := range offsets {
		if offs > 0 {
			indexes[offs] = scatter
		}
	}

	spin.CountSpecials()

	return true
}

func replaceGridFunction(replace util.Indexes) func(*comp.Spin) bool {
	return func(spin *comp.Spin) bool {
		grid := spin.Indexes()
		for ix := range replace {
			grid[ix] = replace[ix]
		}
		return true
	}
}

func script2FullPage(spin *comp.Spin) bool {
	seq := spin.SpinSeq()
	if seq < 11 && !spin.TestChance2(21.0*float64(seq/4)) {
		return false
	}

	symbol := script2SymbolWeights.RandomIndex(spin.PRNG())
	grid := spin.Indexes()
	for ix := range grid {
		grid[ix] = symbol
	}

	return true
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
		case wild2:
			ss[ix] = comp.NewSymbol(id, names[ix], resources[ix], payouts[ix], weights[ix], comp.WithKind(comp.Wild), comp.WithMultiplier(2))
		case wild3:
			ss[ix] = comp.NewSymbol(id, names[ix], resources[ix], payouts[ix], weights[ix], comp.WithKind(comp.Wild), comp.WithMultiplier(3))
		case wild4:
			ss[ix] = comp.NewSymbol(id, names[ix], resources[ix], payouts[ix], weights[ix], comp.WithKind(comp.Wild), comp.WithMultiplier(4))
		default:
			ss[ix] = comp.NewSymbol(id, names[ix], resources[ix], payouts[ix], weights[ix])
		}
	}
	s1 := comp.NewSymbolSet(ss...)

	if symbols == nil {
		symbols = s1
	}

	return comp.NewSlots(
		comp.Grid(reels, rows),
		comp.WithSymbols(s1),
		comp.PayDirections(direction),
		comp.MaxPayout(maxPayout),
		comp.WithRTP(target),
		comp.WithBonusBuy(flagBonusBuy),
		comp.WithRoundFlags(flags...),
		comp.WithActions(actions1, actions2, actions3, actions4),
		comp.WithScriptedRoundSelector(scriptedRounds),
	)
}

func init() {
	initActions()

	slots92 = initSlots(92.0, weights92, actions92First, actions92Free, actions92FirstBB, actions92FreeBB)
	slots94 = initSlots(94.0, weights94, actions94First, actions94Free, actions94FirstBB, actions94FreeBB)
	slots96 = initSlots(96.0, weights96, actions96First, actions96Free, actions96FirstBB, actions96FreeBB)

	slots92params = game.RegularParams{Slots: slots92}
	slots94params = game.RegularParams{Slots: slots94}
	slots96params = game.RegularParams{Slots: slots96}
}

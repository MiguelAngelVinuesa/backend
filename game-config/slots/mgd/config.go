package mgd

import (
	"math"

	comp "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	game "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/slots"
	util "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

const (
	symbolCount = 13
	reels       = 6
	rows        = 4
	direction   = comp.PayLTR
	maxPayout   = 30000.0

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

	scatter = id09
	wild    = id10
	wild2   = id11
	wild3   = id12
	wild4   = id13

	// always add new ID's with a new unique number!
	wildRespin92aID       = 1
	wildRespin92bID       = 2
	wildRespin92cID       = 3
	wildRespin92dID       = 4
	wildRespin92eID       = 5
	wildRespin92fID       = 6
	wildRespin92gID       = 7
	firstWild92ID         = 8
	freeWild92ID          = 9
	firstScatter92ID      = 11
	freeScatter92aID      = 12
	freeScatter92bID      = 13
	freeScatter92cID      = 14
	freeScatter92dID      = 15
	freeScatter92eID      = 16
	freeScatter92ID       = 19
	wildRespin94aID       = 21
	wildRespin94bID       = 22
	wildRespin94cID       = 23
	wildRespin94dID       = 24
	wildRespin94eID       = 25
	wildRespin94fID       = 26
	wildRespin94gID       = 27
	firstWild94ID         = 28
	freeWild94ID          = 29
	firstScatter94ID      = 31
	freeScatter94aID      = 32
	freeScatter94bID      = 33
	freeScatter94cID      = 34
	freeScatter94dID      = 35
	freeScatter94eID      = 36
	freeScatter94ID       = 39
	wildRespin96aID       = 41
	wildRespin96bID       = 42
	wildRespin96cID       = 43
	wildRespin96dID       = 44
	wildRespin96eID       = 45
	wildRespin96fID       = 46
	wildRespin96gID       = 47
	firstWild96ID         = 48
	freeWild96ID          = 49
	firstScatter96ID      = 51
	freeScatter96aID      = 52
	freeScatter96bID      = 53
	freeScatter96cID      = 54
	freeScatter96dID      = 55
	freeScatter96eID      = 56
	freeScatter96ID       = 59
	freeWild2aID          = 81
	freeWild2bID          = 82
	freeWild2cID          = 83
	freeWild2dID          = 84
	freeWild2eID          = 85
	freeWild2ID           = 86
	freeWild3aID          = 87
	freeWild3bID          = 88
	freeWild3cID          = 89
	freeWild3dID          = 90
	freeWild3eID          = 91
	freeWild3ID           = 92
	freeWild4aID          = 93
	freeWild4bID          = 94
	freeWild4cID          = 95
	freeWild4dID          = 96
	freeWild4eID          = 97
	freeWild4ID           = 98
	paylinesID            = 101
	removePayouts1ID      = 102
	removePayouts2ID      = 103
	removePayouts3ID      = 104
	teaser1ID             = 111
	teaser2ID             = 112
	teaser3ID             = 113
	teaser4ID             = 114
	teaser5ID             = 115
	teaserFPID            = 118
	teaserNoScattersID    = 119
	teaserScattersWildsID = 120
	teaserZeroPayWildsID  = 121
	award5ID              = 131
	award8ID              = 132
	retrigger1ID          = 133
	retrigger2ID          = 134
	retrigger3ID          = 135
	retrigger4ID          = 136
	refillID              = 141
	stickies1ID           = 142
	stickies2ID           = 143
	payoutBandsID         = 151
	freeCounterID         = 152
	wildRespinFlagID      = 153

	forceWildRespinID    = 900
	forceBand5ID         = 905
	forceNoPayoutsID     = 906
	forceScrolls4ID      = 911
	forceScrolls4S6ID    = 912
	forceScrollsS10UpID  = 913
	forceScrollsS9UpID   = 914
	forceScrollsS2to8ID  = 915
	forceChaliceFPID     = 920
	forcePotionFPID      = 921
	forceDaggerFPID      = 922
	forceBookFPID        = 923
	forceMedusaFPID      = 924
	forceDeathFPID       = 925
	forceSuccubusFPID    = 926
	forceDemonGoatFPID   = 927
	forceW1S1ID          = 950
	force4RWID           = 951
	force4RWS8UpID       = 952
	forceW4S2R5ID        = 961
	forceW4S3R4ID        = 962
	forceW4S4R5ID        = 963
	forceW4S5R4ID        = 964
	forceW4S7R2ID        = 965
	forceW4S8R3ID        = 966
	forceW4S9UpRandomID  = 991
	forceW4S12UpRandomID = 992

	flagFreeSpinBands = 0
	flagFreeSpinCount = 1
	flagWildRespin    = 2
)

var (
	// WIP WIP WIP
	weights92 = [symbolCount]comp.SymbolOption{
		comp.WithWeights(210, 230, 225, 225, 230, 210),
		comp.WithWeights(230, 225, 175, 175, 225, 230),
		comp.WithWeights(210, 220, 180, 180, 220, 210),
		comp.WithWeights(220, 185, 175, 175, 165, 220),
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
	wildRespinChance92       = 2.32
	fakeRespin1Chance92      = 0.51
	fakeRespin2Chance92      = 0.52
	fakeRespin3Chance92      = 0.53
	fakeRespin4Chance92      = 0.39
	fakeRespin5Chance92      = 0.4
	fakeRespin6Chance92      = 0.41
	firstWildChances92       = []float64{17.7, 15.3, 3}
	freeWildChances92        = []float64{42.0, 27.1, 19, 3}
	scatterRegularChances92  = []float64{23.9, 23.3, 14.2, 7}
	scatterFreeSpinChances92 = [5][]float64{
		{6},
		{8, 1},
		{12, 5, 3},
		{23.3, 11, 4, 1},
		{31, 13.5, 8, 4}}

	// WIP WIP WIP
	weights94 = [symbolCount]comp.SymbolOption{
		comp.WithWeights(210, 230, 225, 225, 230, 210),
		comp.WithWeights(230, 225, 175, 175, 225, 230),
		comp.WithWeights(210, 220, 180, 180, 220, 210),
		comp.WithWeights(220, 185, 175, 175, 165, 220),
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
	wildRespinChance94       = 2.33
	fakeRespin1Chance94      = 0.51
	fakeRespin2Chance94      = 0.52
	fakeRespin3Chance94      = 0.53
	fakeRespin4Chance94      = 0.39
	fakeRespin5Chance94      = 0.4
	fakeRespin6Chance94      = 0.41
	firstWildChances94       = []float64{18, 16, 3}
	freeWildChances94        = []float64{42.0, 27.1, 19, 3}
	scatterRegularChances94  = []float64{24.2, 23.8, 14.2, 7}
	scatterFreeSpinChances94 = [5][]float64{
		{6},
		{8, 1},
		{12, 5, 3},
		{23.3, 11, 4, 1},
		{31, 13.5, 8, 4}}

	// WIP WIP WIP
	weights96 = [symbolCount]comp.SymbolOption{
		comp.WithWeights(210, 230, 225, 225, 230, 210),
		comp.WithWeights(230, 225, 175, 175, 225, 230),
		comp.WithWeights(210, 220, 180, 180, 220, 210),
		comp.WithWeights(220, 185, 175, 175, 165, 220),
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
	wildRespinChance96       = 2.34
	fakeRespin1Chance96      = 0.51
	fakeRespin2Chance96      = 0.52
	fakeRespin3Chance96      = 0.53
	fakeRespin4Chance96      = 0.39
	fakeRespin5Chance96      = 0.4
	fakeRespin6Chance96      = 0.41
	firstWildChances96       = []float64{18.8, 16.3, 3}
	freeWildChances96        = []float64{42.0, 27.1, 19, 3}
	scatterRegularChances96  = []float64{24.6, 24.3, 14.2, 7}
	scatterFreeSpinChances96 = [5][]float64{
		{6},
		{8, 1},
		{12, 5, 3},
		{23.3, 11, 4, 1},
		{31, 13.5, 8, 4}}
)

var (
	// free spin payout band weights.
	freeSpinBandWeights = []float64{165, 75, 25, 9, 2.3}

	// weights to remove payouts based on total payout factor during first spins.
	reduceBandsFirst = []comp.RemovePayoutBand{
		{MinPayout: 0, MaxPayout: 2.5, RemoveChance: 72},
		{MinPayout: 2.5, MaxPayout: 5, RemoveChance: 36.5},
		{MinPayout: 5, MaxPayout: 10, RemoveChance: 47},
		{MinPayout: 10, MaxPayout: 40, RemoveChance: 25},
		{MinPayout: 40, MaxPayout: 125, RemoveChance: 42},
		{MinPayout: 125, MaxPayout: 200, RemoveChance: 15},
		{MinPayout: 200, MaxPayout: maxPayout, RemoveChance: 30},
	}

	// weights to remove payouts based on total payout factor during first few free spins.
	reduceBandsFree1 = []comp.RemovePayoutBand{
		{MinPayout: 0.0, MaxPayout: 30, RemoveChance: 15},
		{MinPayout: 30, MaxPayout: maxPayout, RemoveChance: 27},
	}

	// weights to remove payouts based on total payout factor during free spins.
	reduceBandsFree2 = []comp.RemovePayoutBand{
		{MinPayout: 0, MaxPayout: 5, RemoveChance: 5},
		{MinPayout: 5, MaxPayout: 10, RemoveChance: 2},
		{MinPayout: 10, MaxPayout: maxPayout, RemoveChance: 8},
	}

	// chance to get 2x wild in free spins per band.
	freeWild2Chances1 = []float64{2.2}
	freeWild2Chances2 = []float64{7}
	freeWild2Chances3 = []float64{10}
	freeWild2Chances4 = []float64{19, 3}
	freeWild2Chances5 = []float64{45, 6}

	// chance to get 3x wild in free spins per band.
	freeWild3Chances1 = []float64{0.6}
	freeWild3Chances2 = []float64{1.8}
	freeWild3Chances3 = []float64{3.2}
	freeWild3Chances4 = []float64{7, 1}
	freeWild3Chances5 = []float64{18, 3}

	// chance to get 4x wild in free spins per band.
	freeWild4Chances1 = []float64{0.1}
	freeWild4Chances2 = []float64{0.4}
	freeWild4Chances3 = []float64{2}
	freeWild4Chances4 = []float64{9, 1}
	freeWild4Chances5 = []float64{14, 2}

	script1weight = 0.14
	script2weight = 0.22
	script3weight = 0.22
	script4weight = 0.22
	script5weight = 0.1
	script6weight = 0.1
	script9weight = 999.0
)

var (
	gridMask = []uint8{2, 3, 4, 4, 3, 2}

	n01 = comp.WithName("Chalice")
	n02 = comp.WithName("Potion")
	n03 = comp.WithName("Daggers")
	n04 = comp.WithName("Book")
	n05 = comp.WithName("Medusa")
	n06 = comp.WithName("Death")
	n07 = comp.WithName("Succubus")
	n08 = comp.WithName("Demon Goat")
	n09 = comp.WithName("Scatter Scroll")
	n10 = comp.WithName("Wild Devil")
	n11 = comp.WithName("Wild Devil x2")
	n12 = comp.WithName("Wild Devil x3")
	n13 = comp.WithName("Wild Devil x4")

	r01 = comp.WithResource("l4")
	r02 = comp.WithResource("l3")
	r03 = comp.WithResource("l2")
	r04 = comp.WithResource("l1")
	r05 = comp.WithResource("h4")
	r06 = comp.WithResource("h3")
	r07 = comp.WithResource("h2")
	r08 = comp.WithResource("h1")
	r09 = comp.WithResource("bonus")
	r10 = comp.WithResource("devil")
	r11 = comp.WithResource("devil2")
	r12 = comp.WithResource("devil3")
	r13 = comp.WithResource("devil4")

	p01 = comp.WithPayouts(0, 0, 0.2, 0.5, 0.8, 1.5)
	p02 = comp.WithPayouts(0, 0, 0.2, 0.5, 0.8, 1.5)
	p03 = comp.WithPayouts(0, 0, 0.2, 0.8, 1, 2)
	p04 = comp.WithPayouts(0, 0, 0.2, 0.8, 1, 2)
	p05 = comp.WithPayouts(0, 0, 0.5, 1, 2, 3)
	p06 = comp.WithPayouts(0, 0, 0.5, 1, 2, 5)
	p07 = comp.WithPayouts(0, 0, 1, 2, 3, 10)
	p08 = comp.WithPayouts(0, 0, 1, 3, 5, 15)
	p09 = comp.WithPayouts()
	p10 = p09
	p11 = p09
	p12 = p09
	p13 = p09

	ids       = [symbolCount]util.Index{id01, id02, id03, id04, id05, id06, id07, id08, id09, id10, id11, id12, id13}
	names     = [symbolCount]comp.SymbolOption{n01, n02, n03, n04, n05, n06, n07, n08, n09, n10, n11, n12, n13}
	resources = [symbolCount]comp.SymbolOption{r01, r02, r03, r04, r05, r06, r07, r08, r09, r10, r11, r12, r13}
	payouts   = [symbolCount]comp.SymbolOption{p01, p02, p03, p04, p05, p06, p07, p08, p09, p10, p11, p12, p13}

	wildRespinGrid   = comp.GridOffsets{{0, 0}, {0, 1}, {5, 0}, {5, 1}}
	fakeRespinGrid1  = comp.GridOffsets{{0, 0}, {0, 1}}
	fakeRespinGrid2  = comp.GridOffsets{{0, 0}, {0, 1}, {5, 0}}
	fakeRespinGrid3  = comp.GridOffsets{{0, 0}, {0, 1}, {5, 1}}
	fakeRespinGrid4  = comp.GridOffsets{{5, 0}, {5, 1}}
	fakeRespinGrid5  = comp.GridOffsets{{0, 0}, {5, 0}, {5, 1}}
	fakeRespinGrid6  = comp.GridOffsets{{0, 1}, {5, 0}, {5, 1}}
	wildRespinCenter = comp.GridOffsets{{0, 0}}
	wildRespin       = util.AcquireWeighting().AddWeights(util.Indexes{wild}, []float64{1})
	scatterReels     = util.UInt8s{2, 3, 4, 5}
	scatterReels2    = []int{2, 3, 4, 5}
	freeSpinBands    = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5}, freeSpinBandWeights)
	freeSpinBand5    = util.AcquireWeighting().AddWeights(util.Indexes{1}, []float64{1})

	flag0 = comp.NewRoundFlag(flagFreeSpinBands, "free spin band")
	flag1 = comp.NewRoundFlag(flagFreeSpinCount, "free spin count")
	flag2 = comp.NewRoundFlag(flagWildRespin, "wild respin landed")
	flags = comp.RoundFlags{flag0, flag1, flag2}

	fakeSpinChance = 0.5

	teasers1 = []comp.FakeSpin{
		{Indexes: util.Indexes{8, 8, 0, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 0, 8, 8, 0, 0}, MatchReels: util.UInt8s{1, 6}, MatchSymbol: 10, MatchInverse: true, ReplaceSymbols: util.Indexes{8, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{7, 7, 0, 0, 0, 0, 0, 0, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 0, 7, 7, 0, 0}, MatchReels: util.UInt8s{1, 6}, MatchSymbol: 10, MatchInverse: true, ReplaceSymbols: util.Indexes{7, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{6, 6, 0, 0, 0, 0, 0, 0, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 0, 6, 6, 0, 0}, MatchReels: util.UInt8s{1, 6}, MatchSymbol: 10, MatchInverse: true, ReplaceSymbols: util.Indexes{6, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{5, 5, 0, 0, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 0, 5, 5, 0, 0}, MatchReels: util.UInt8s{1, 6}, MatchSymbol: 10, MatchInverse: true, ReplaceSymbols: util.Indexes{5, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{8, 8, 0, 0, 8, 8, 8, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 0, 8, 8, 0, 0}, MatchReels: util.UInt8s{1, 6}, MatchSymbol: 10, MatchInverse: true, ReplaceSymbols: util.Indexes{8, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{7, 7, 0, 0, 7, 7, 7, 0, 0, 0, 0, 0, 7, 7, 7, 7, 7, 7, 7, 0, 7, 7, 0, 0}, MatchReels: util.UInt8s{1, 6}, MatchSymbol: 10, MatchInverse: true, ReplaceSymbols: util.Indexes{7, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{6, 6, 0, 0, 6, 6, 6, 0, 0, 0, 0, 0, 6, 6, 6, 6, 6, 6, 6, 0, 6, 6, 0, 0}, MatchReels: util.UInt8s{1, 6}, MatchSymbol: 10, MatchInverse: true, ReplaceSymbols: util.Indexes{6, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{5, 5, 0, 0, 5, 5, 5, 0, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 0, 5, 5, 0, 0}, MatchReels: util.UInt8s{1, 6}, MatchSymbol: 10, MatchInverse: true, ReplaceSymbols: util.Indexes{5, 10}, ReplaceReels: util.UInt8s{3}},
	}

	teasers2 = []comp.FakeSpin{
		{Indexes: util.Indexes{8, 8, 0, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 0, 10, 0, 0, 0}, MatchReels: util.UInt8s{6}, ReplaceSymbols: util.Indexes{8, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{7, 7, 0, 0, 0, 0, 0, 0, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 0, 10, 0, 0, 0}, MatchReels: util.UInt8s{6}, ReplaceSymbols: util.Indexes{7, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{6, 6, 0, 0, 0, 0, 0, 0, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 0, 10, 0, 0, 0}, MatchReels: util.UInt8s{6}, ReplaceSymbols: util.Indexes{6, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{5, 5, 0, 0, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 0, 10, 0, 0, 0}, MatchReels: util.UInt8s{6}, ReplaceSymbols: util.Indexes{5, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{8, 8, 0, 0, 8, 8, 8, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 0, 10, 0, 0, 0}, MatchReels: util.UInt8s{6}, ReplaceSymbols: util.Indexes{8, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{7, 7, 0, 0, 7, 7, 7, 0, 0, 0, 0, 0, 7, 7, 7, 7, 7, 7, 7, 0, 10, 0, 0, 0}, MatchReels: util.UInt8s{6}, ReplaceSymbols: util.Indexes{7, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{6, 6, 0, 0, 6, 6, 6, 0, 0, 0, 0, 0, 6, 6, 6, 6, 6, 6, 6, 0, 10, 0, 0, 0}, MatchReels: util.UInt8s{6}, ReplaceSymbols: util.Indexes{6, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{5, 5, 0, 0, 5, 5, 5, 0, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 0, 10, 0, 0, 0}, MatchReels: util.UInt8s{6}, ReplaceSymbols: util.Indexes{5, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{8, 8, 0, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 0, 0, 10, 0, 0}, MatchReels: util.UInt8s{6}, ReplaceSymbols: util.Indexes{8, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{7, 7, 0, 0, 0, 0, 0, 0, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 0, 0, 10, 0, 0}, MatchReels: util.UInt8s{6}, ReplaceSymbols: util.Indexes{7, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{6, 6, 0, 0, 0, 0, 0, 0, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 0, 0, 10, 0, 0}, MatchReels: util.UInt8s{6}, ReplaceSymbols: util.Indexes{6, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{5, 5, 0, 0, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 0, 0, 10, 0, 0}, MatchReels: util.UInt8s{6}, ReplaceSymbols: util.Indexes{5, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{8, 8, 0, 0, 8, 8, 8, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 0, 0, 10, 0, 0}, MatchReels: util.UInt8s{6}, ReplaceSymbols: util.Indexes{8, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{7, 7, 0, 0, 7, 7, 7, 0, 0, 0, 0, 0, 7, 7, 7, 7, 7, 7, 7, 0, 0, 10, 0, 0}, MatchReels: util.UInt8s{6}, ReplaceSymbols: util.Indexes{7, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{6, 6, 0, 0, 6, 6, 6, 0, 0, 0, 0, 0, 6, 6, 6, 6, 6, 6, 6, 0, 0, 10, 0, 0}, MatchReels: util.UInt8s{6}, ReplaceSymbols: util.Indexes{6, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{5, 5, 0, 0, 5, 5, 5, 0, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 0, 0, 10, 0, 0}, MatchReels: util.UInt8s{6}, ReplaceSymbols: util.Indexes{5, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{10, 0, 0, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 0, 8, 8, 0, 0}, MatchReels: util.UInt8s{1}, ReplaceSymbols: util.Indexes{8, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{10, 0, 0, 0, 0, 0, 0, 0, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 0, 7, 7, 0, 0}, MatchReels: util.UInt8s{1}, ReplaceSymbols: util.Indexes{7, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{10, 0, 0, 0, 0, 0, 0, 0, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 0, 6, 6, 0, 0}, MatchReels: util.UInt8s{1}, ReplaceSymbols: util.Indexes{6, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{10, 0, 0, 0, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 0, 5, 5, 0, 0}, MatchReels: util.UInt8s{1}, ReplaceSymbols: util.Indexes{5, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{10, 0, 0, 0, 8, 8, 8, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 0, 8, 8, 0, 0}, MatchReels: util.UInt8s{1}, ReplaceSymbols: util.Indexes{8, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{10, 0, 0, 0, 7, 7, 7, 0, 0, 0, 0, 0, 7, 7, 7, 7, 7, 7, 7, 0, 7, 7, 0, 0}, MatchReels: util.UInt8s{1}, ReplaceSymbols: util.Indexes{7, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{10, 0, 0, 0, 6, 6, 6, 0, 0, 0, 0, 0, 6, 6, 6, 6, 6, 6, 6, 0, 6, 6, 0, 0}, MatchReels: util.UInt8s{1}, ReplaceSymbols: util.Indexes{6, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{10, 0, 0, 0, 5, 5, 5, 0, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 0, 5, 5, 0, 0}, MatchReels: util.UInt8s{1}, ReplaceSymbols: util.Indexes{5, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{0, 10, 0, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 0, 8, 8, 0, 0}, MatchReels: util.UInt8s{1}, ReplaceSymbols: util.Indexes{8, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{0, 10, 0, 0, 0, 0, 0, 0, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 0, 7, 7, 0, 0}, MatchReels: util.UInt8s{1}, ReplaceSymbols: util.Indexes{7, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{0, 10, 0, 0, 0, 0, 0, 0, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 0, 6, 6, 0, 0}, MatchReels: util.UInt8s{1}, ReplaceSymbols: util.Indexes{6, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{0, 10, 0, 0, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 0, 5, 5, 0, 0}, MatchReels: util.UInt8s{1}, ReplaceSymbols: util.Indexes{5, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{0, 10, 0, 0, 8, 8, 8, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 0, 8, 8, 0, 0}, MatchReels: util.UInt8s{1}, ReplaceSymbols: util.Indexes{8, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{0, 10, 0, 0, 7, 7, 7, 0, 0, 0, 0, 0, 7, 7, 7, 7, 7, 7, 7, 0, 7, 7, 0, 0}, MatchReels: util.UInt8s{1}, ReplaceSymbols: util.Indexes{7, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{0, 10, 0, 0, 6, 6, 6, 0, 0, 0, 0, 0, 6, 6, 6, 6, 6, 6, 6, 0, 6, 6, 0, 0}, MatchReels: util.UInt8s{1}, ReplaceSymbols: util.Indexes{6, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{0, 10, 0, 0, 5, 5, 5, 0, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 0, 5, 5, 0, 0}, MatchReels: util.UInt8s{1}, ReplaceSymbols: util.Indexes{5, 10}, ReplaceReels: util.UInt8s{3}},
	}

	teasers3 = []comp.FakeSpin{
		{Indexes: util.Indexes{8, 8, 0, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{6}, ReplaceSymbols: util.Indexes{8, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{7, 7, 0, 0, 0, 0, 0, 0, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{6}, ReplaceSymbols: util.Indexes{7, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{6, 6, 0, 0, 0, 0, 0, 0, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{6}, ReplaceSymbols: util.Indexes{6, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{5, 5, 0, 0, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{6}, ReplaceSymbols: util.Indexes{5, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{8, 8, 0, 0, 8, 8, 8, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{6}, ReplaceSymbols: util.Indexes{8, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{7, 7, 0, 0, 7, 7, 7, 0, 0, 0, 0, 0, 7, 7, 7, 7, 7, 7, 7, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{6}, ReplaceSymbols: util.Indexes{7, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{6, 6, 0, 0, 6, 6, 6, 0, 0, 0, 0, 0, 6, 6, 6, 6, 6, 6, 6, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{6}, ReplaceSymbols: util.Indexes{6, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{5, 5, 0, 0, 5, 5, 5, 0, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{6}, ReplaceSymbols: util.Indexes{5, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{10, 10, 0, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 0, 8, 8, 0, 0}, MatchReels: util.UInt8s{1}, ReplaceSymbols: util.Indexes{8, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{10, 10, 0, 0, 0, 0, 0, 0, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 0, 7, 7, 0, 0}, MatchReels: util.UInt8s{1}, ReplaceSymbols: util.Indexes{7, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{10, 10, 0, 0, 0, 0, 0, 0, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 0, 6, 6, 0, 0}, MatchReels: util.UInt8s{1}, ReplaceSymbols: util.Indexes{6, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{10, 10, 0, 0, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 0, 5, 5, 0, 0}, MatchReels: util.UInt8s{1}, ReplaceSymbols: util.Indexes{5, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{10, 10, 0, 0, 8, 8, 8, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 0, 8, 8, 0, 0}, MatchReels: util.UInt8s{1}, ReplaceSymbols: util.Indexes{8, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{10, 10, 0, 0, 7, 7, 7, 0, 0, 0, 0, 0, 7, 7, 7, 7, 7, 7, 7, 0, 7, 7, 0, 0}, MatchReels: util.UInt8s{1}, ReplaceSymbols: util.Indexes{7, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{10, 10, 0, 0, 6, 6, 6, 0, 0, 0, 0, 0, 6, 6, 6, 6, 6, 6, 6, 0, 6, 6, 0, 0}, MatchReels: util.UInt8s{1}, ReplaceSymbols: util.Indexes{6, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{10, 10, 0, 0, 5, 5, 5, 0, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 0, 5, 5, 0, 0}, MatchReels: util.UInt8s{1}, ReplaceSymbols: util.Indexes{5, 10}, ReplaceReels: util.UInt8s{3}},
	}

	teasers4 = []comp.FakeSpin{
		{Indexes: util.Indexes{10, 0, 0, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{8, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{10, 0, 0, 0, 0, 0, 0, 0, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{7, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{10, 0, 0, 0, 0, 0, 0, 0, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{6, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{10, 0, 0, 0, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{5, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{10, 0, 0, 0, 8, 8, 8, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{8, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{10, 0, 0, 0, 7, 7, 7, 0, 0, 0, 0, 0, 7, 7, 7, 7, 7, 7, 7, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{7, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{10, 0, 0, 0, 6, 6, 6, 0, 0, 0, 0, 0, 6, 6, 6, 6, 6, 6, 6, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{6, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{10, 0, 0, 0, 5, 5, 5, 0, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{5, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{0, 10, 0, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{8, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{0, 10, 0, 0, 0, 0, 0, 0, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{7, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{0, 10, 0, 0, 0, 0, 0, 0, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{6, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{0, 10, 0, 0, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{5, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{0, 10, 0, 0, 8, 8, 8, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{8, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{0, 10, 0, 0, 7, 7, 7, 0, 0, 0, 0, 0, 7, 7, 7, 7, 7, 7, 7, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{7, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{0, 10, 0, 0, 6, 6, 6, 0, 0, 0, 0, 0, 6, 6, 6, 6, 6, 6, 6, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{6, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{0, 10, 0, 0, 5, 5, 5, 0, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{5, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{10, 10, 0, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 0, 10, 0, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{8, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{10, 10, 0, 0, 0, 0, 0, 0, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 0, 10, 0, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{7, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{10, 10, 0, 0, 0, 0, 0, 0, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 0, 10, 0, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{6, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{10, 10, 0, 0, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 0, 10, 0, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{5, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{10, 10, 0, 0, 8, 8, 8, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 0, 10, 0, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{8, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{10, 10, 0, 0, 7, 7, 7, 0, 0, 0, 0, 0, 7, 7, 7, 7, 7, 7, 7, 0, 10, 0, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{7, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{10, 10, 0, 0, 6, 6, 6, 0, 0, 0, 0, 0, 6, 6, 6, 6, 6, 6, 6, 0, 10, 0, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{6, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{10, 10, 0, 0, 5, 5, 5, 0, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 0, 10, 0, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{5, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{10, 10, 0, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 0, 0, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{8, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{10, 10, 0, 0, 0, 0, 0, 0, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 0, 0, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{7, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{10, 10, 0, 0, 0, 0, 0, 0, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 0, 0, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{6, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{10, 10, 0, 0, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 0, 0, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{5, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{10, 10, 0, 0, 8, 8, 8, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 0, 0, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{8, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{10, 10, 0, 0, 7, 7, 7, 0, 0, 0, 0, 0, 7, 7, 7, 7, 7, 7, 7, 0, 0, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{7, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{10, 10, 0, 0, 6, 6, 6, 0, 0, 0, 0, 0, 6, 6, 6, 6, 6, 6, 6, 0, 0, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{6, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{10, 10, 0, 0, 5, 5, 5, 0, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 0, 0, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{5, 10}, ReplaceReels: util.UInt8s{3}},
	}

	teasers5 = []comp.FakeSpin{
		{Indexes: util.Indexes{10, 10, 0, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 8, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{8, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{10, 10, 0, 0, 0, 0, 0, 0, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 7, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{7, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{10, 10, 0, 0, 0, 0, 0, 0, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{6, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{10, 10, 0, 0, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 5, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{5, 10}, ReplaceReels: util.UInt8s{2}},
		{Indexes: util.Indexes{10, 10, 0, 0, 8, 8, 8, 0, 0, 0, 0, 0, 8, 8, 8, 8, 8, 8, 8, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{8, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{10, 10, 0, 0, 7, 7, 7, 0, 0, 0, 0, 0, 7, 7, 7, 7, 7, 7, 7, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{7, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{10, 10, 0, 0, 6, 6, 6, 0, 0, 0, 0, 0, 6, 6, 6, 6, 6, 6, 6, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{6, 10}, ReplaceReels: util.UInt8s{3}},
		{Indexes: util.Indexes{10, 10, 0, 0, 5, 5, 5, 0, 0, 0, 0, 0, 5, 5, 5, 5, 5, 5, 5, 0, 10, 10, 0, 0}, MatchReels: util.UInt8s{1, 6}, ReplaceSymbols: util.Indexes{5, 10}, ReplaceReels: util.UInt8s{3}},
	}

	fullPageR1to4shape       = comp.GridOffsets{{0, 0}, {0, 1}, {0, 2}, {1, 0}, {1, 1}, {1, 2}, {1, 3}, {2, 0}, {2, 1}, {2, 2}, {2, 3}, {3, 0}, {3, 1}, {3, 2}}
	fullPageR1offset         = comp.GridOffsets{{1, 0}}
	fullPage4Xweights        = util.AcquireWeighting().AddWeights(util.Indexes{wild4}, []float64{100})
	fullPageChaliceWeights   = util.AcquireWeighting().AddWeights(util.Indexes{1}, []float64{100})
	fullPagePotionWeights    = util.AcquireWeighting().AddWeights(util.Indexes{2}, []float64{100})
	fullPageDaggerWeights    = util.AcquireWeighting().AddWeights(util.Indexes{3}, []float64{100})
	fullPageBookWeights      = util.AcquireWeighting().AddWeights(util.Indexes{4}, []float64{100})
	fullPageMedusaWeights    = util.AcquireWeighting().AddWeights(util.Indexes{5}, []float64{100})
	fullPageDeathWeights     = util.AcquireWeighting().AddWeights(util.Indexes{6}, []float64{100})
	fullPageSuccubusWeights  = util.AcquireWeighting().AddWeights(util.Indexes{7}, []float64{100})
	fullPageDemonGoatWeights = util.AcquireWeighting().AddWeights(util.Indexes{8}, []float64{100})

	actions92all   comp.SpinActions
	actions92first comp.SpinActions
	actions92free  comp.SpinActions
	actions94all   comp.SpinActions
	actions94first comp.SpinActions
	actions94free  comp.SpinActions
	actions96all   comp.SpinActions
	actions96first comp.SpinActions
	actions96free  comp.SpinActions

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
	// generate (fake) wild respin & scatters RTP 92.
	wildRespin92a := comp.NewGenerateShapeAction(wildRespinChance92, wildRespinGrid, wildRespinCenter, wildRespin)
	wildRespin92a.Describe(wildRespin92aID, "generate wild respin - RTP 92")
	wildRespin92b := comp.NewGenerateShapeAction(fakeRespin6Chance92, fakeRespinGrid6, wildRespinCenter, wildRespin).WithAlternate(wildRespin92a)
	wildRespin92b.Describe(wildRespin92bID, "generate fake respin right (3) - RTP 92")
	wildRespin92c := comp.NewGenerateShapeAction(fakeRespin5Chance92, fakeRespinGrid5, wildRespinCenter, wildRespin).WithAlternate(wildRespin92b)
	wildRespin92c.Describe(wildRespin92cID, "generate fake respin right (2) - RTP 92")
	wildRespin92d := comp.NewGenerateShapeAction(fakeRespin4Chance92, fakeRespinGrid4, wildRespinCenter, wildRespin).WithAlternate(wildRespin92c)
	wildRespin92d.Describe(wildRespin92dID, "generate fake respin right (1) - RTP 92")
	wildRespin92e := comp.NewGenerateShapeAction(fakeRespin3Chance92, fakeRespinGrid3, wildRespinCenter, wildRespin).WithAlternate(wildRespin92d)
	wildRespin92e.Describe(wildRespin92eID, "generate fake respin left (3) - RTP 92")
	wildRespin92f := comp.NewGenerateShapeAction(fakeRespin2Chance92, fakeRespinGrid2, wildRespinCenter, wildRespin).WithAlternate(wildRespin92e)
	wildRespin92f.Describe(wildRespin92fID, "generate fake respin left (2) - RTP 92")
	wildRespin92g := comp.NewGenerateShapeAction(fakeRespin1Chance92, fakeRespinGrid1, wildRespinCenter, wildRespin).WithAlternate(wildRespin92f)
	wildRespin92g.WithTriggerFilters(comp.OnFirstSpin) // leave this in place as it must not be called during refills!!!
	wildRespin92g.Describe(wildRespin92gID, "generate fake respin left (1) - RTP 92")

	firstWild92 := comp.NewGenerateSymbolAction(wild, firstWildChances92).GenerateNoDupes()
	firstWild92.WithTriggerFilters(comp.OnFirstSpin) // leave this in place as it must not be called during refills!!!
	firstWild92.Describe(firstWild92ID, "generate wilds - first spin - RTP 92")

	freeWild92 := comp.NewGenerateSymbolAction(wild, freeWildChances92, scatterReels...)
	freeWild92.Describe(freeWild92ID, "generate wilds - free spins - RTP 92")

	firstScatter92 := comp.NewGenerateSymbolAction(scatter, scatterRegularChances92, scatterReels...).GenerateNoDupes()
	firstScatter92.WithTriggerFilters(comp.OnFirstSpin) // leave this in place as it must not be called during refills!!!
	firstScatter92.Describe(firstScatter92ID, "generate scatters - first spin - RTP 92")

	freeScatter92a := comp.NewGenerateSymbolAction(scatter, scatterFreeSpinChances92[0], scatterReels...).GenerateNoDupes()
	freeScatter92a.Describe(freeScatter92aID, "generate scatters - free spin band 1 - RTP 92")
	freeScatter92b := comp.NewGenerateSymbolAction(scatter, scatterFreeSpinChances92[1], scatterReels...).GenerateNoDupes()
	freeScatter92b.Describe(freeScatter92bID, "generate scatters - free spin band 2 - RTP 92")
	freeScatter92c := comp.NewGenerateSymbolAction(scatter, scatterFreeSpinChances92[2], scatterReels...).GenerateNoDupes()
	freeScatter92c.Describe(freeScatter92cID, "generate scatters - free spin band 3 - RTP 92")
	freeScatter92d := comp.NewGenerateSymbolAction(scatter, scatterFreeSpinChances92[3], scatterReels...).GenerateNoDupes()
	freeScatter92d.Describe(freeScatter92dID, "generate scatters - free spin band 4 - RTP 92")
	freeScatter92e := comp.NewGenerateSymbolAction(scatter, scatterFreeSpinChances92[4], scatterReels...).GenerateNoDupes()
	freeScatter92e.Describe(freeScatter92eID, "generate scatters - free spin band 5 - RTP 92")

	freeScatter92 := comp.NewMultiActionFlagValue(flagFreeSpinBands, 1, freeScatter92a, 2, freeScatter92b, 3, freeScatter92c, 4, freeScatter92d, 5, freeScatter92e)
	freeScatter92.Describe(freeScatter92ID, "multi-select generate scatters - RTP 92")

	// generate (fake) wild respin & scatters RTP 94.
	wildRespin94a := comp.NewGenerateShapeAction(wildRespinChance94, wildRespinGrid, wildRespinCenter, wildRespin)
	wildRespin94a.Describe(wildRespin94aID, "generate wild respin - RTP 94")
	wildRespin94b := comp.NewGenerateShapeAction(fakeRespin6Chance94, fakeRespinGrid6, wildRespinCenter, wildRespin).WithAlternate(wildRespin94a)
	wildRespin94b.Describe(wildRespin94bID, "generate fake respin right (3) - RTP 94")
	wildRespin94c := comp.NewGenerateShapeAction(fakeRespin5Chance94, fakeRespinGrid5, wildRespinCenter, wildRespin).WithAlternate(wildRespin94b)
	wildRespin94c.Describe(wildRespin94cID, "generate fake respin right (2) - RTP 94")
	wildRespin94d := comp.NewGenerateShapeAction(fakeRespin4Chance94, fakeRespinGrid4, wildRespinCenter, wildRespin).WithAlternate(wildRespin94c)
	wildRespin94d.Describe(wildRespin94dID, "generate fake respin right (1) - RTP 94")
	wildRespin94e := comp.NewGenerateShapeAction(fakeRespin3Chance94, fakeRespinGrid3, wildRespinCenter, wildRespin).WithAlternate(wildRespin94d)
	wildRespin94e.Describe(wildRespin94eID, "generate fake respin left (3) - RTP 94")
	wildRespin94f := comp.NewGenerateShapeAction(fakeRespin2Chance94, fakeRespinGrid2, wildRespinCenter, wildRespin).WithAlternate(wildRespin94e)
	wildRespin94f.Describe(wildRespin94fID, "generate fake respin left (2) - RTP 94")
	wildRespin94g := comp.NewGenerateShapeAction(fakeRespin1Chance94, fakeRespinGrid1, wildRespinCenter, wildRespin).WithAlternate(wildRespin94f)
	wildRespin94g.WithTriggerFilters(comp.OnFirstSpin) // leave this in place as it must not be called during refills!!!
	wildRespin94g.Describe(wildRespin94gID, "generate fake respin left (1) - RTP 94")

	firstWild94 := comp.NewGenerateSymbolAction(wild, firstWildChances94).GenerateNoDupes()
	firstWild94.WithTriggerFilters(comp.OnFirstSpin) // leave this in place as it must not be called during refills!!!
	firstWild94.Describe(firstWild94ID, "generate wilds - first spin - RTP 94")

	freeWild94 := comp.NewGenerateSymbolAction(wild, freeWildChances94, scatterReels...)
	freeWild94.Describe(freeWild94ID, "generate wilds - free spins - RTP 94")

	firstScatter94 := comp.NewGenerateSymbolAction(scatter, scatterRegularChances94, scatterReels...).GenerateNoDupes()
	firstScatter94.WithTriggerFilters(comp.OnFirstSpin) // leave this in place as it must not be called during refills!!!
	firstScatter94.Describe(firstScatter94ID, "generate scatters - first spin - RTP 94")

	freeScatter94a := comp.NewGenerateSymbolAction(scatter, scatterFreeSpinChances94[0], scatterReels...).GenerateNoDupes()
	freeScatter94a.Describe(freeScatter94aID, "generate scatters - free spin band 1 - RTP 94")
	freeScatter94b := comp.NewGenerateSymbolAction(scatter, scatterFreeSpinChances94[1], scatterReels...).GenerateNoDupes()
	freeScatter94b.Describe(freeScatter94bID, "generate scatters - free spin band 2 - RTP 94")
	freeScatter94c := comp.NewGenerateSymbolAction(scatter, scatterFreeSpinChances94[2], scatterReels...).GenerateNoDupes()
	freeScatter94c.Describe(freeScatter94cID, "generate scatters - free spin band 3 - RTP 94")
	freeScatter94d := comp.NewGenerateSymbolAction(scatter, scatterFreeSpinChances94[3], scatterReels...).GenerateNoDupes()
	freeScatter94d.Describe(freeScatter94dID, "generate scatters - free spin band 4 - RTP 94")
	freeScatter94e := comp.NewGenerateSymbolAction(scatter, scatterFreeSpinChances94[4], scatterReels...).GenerateNoDupes()
	freeScatter94e.Describe(freeScatter94eID, "generate scatters - free spin band 5 - RTP 94")

	freeScatter94 := comp.NewMultiActionFlagValue(flagFreeSpinBands, 1, freeScatter94a, 2, freeScatter94b, 3, freeScatter94c, 4, freeScatter94d, 5, freeScatter94e)
	freeScatter94.Describe(freeScatter94ID, "multi-select generate scatters - RTP 94")

	// generate (fake) wild respin & scatters RTP 96.
	wildRespin96a := comp.NewGenerateShapeAction(wildRespinChance96, wildRespinGrid, wildRespinCenter, wildRespin)
	wildRespin96a.Describe(wildRespin96aID, "generate wild respin - RTP 96")
	wildRespin96b := comp.NewGenerateShapeAction(fakeRespin6Chance96, fakeRespinGrid6, wildRespinCenter, wildRespin).WithAlternate(wildRespin96a)
	wildRespin96b.Describe(wildRespin96bID, "generate fake respin right (3) - RTP 96")
	wildRespin96c := comp.NewGenerateShapeAction(fakeRespin5Chance96, fakeRespinGrid5, wildRespinCenter, wildRespin).WithAlternate(wildRespin96b)
	wildRespin96c.Describe(wildRespin96cID, "generate fake respin right (2) - RTP 96")
	wildRespin96d := comp.NewGenerateShapeAction(fakeRespin4Chance96, fakeRespinGrid4, wildRespinCenter, wildRespin).WithAlternate(wildRespin96c)
	wildRespin96d.Describe(wildRespin96dID, "generate fake respin right (1) - RTP 96")
	wildRespin96e := comp.NewGenerateShapeAction(fakeRespin3Chance96, fakeRespinGrid3, wildRespinCenter, wildRespin).WithAlternate(wildRespin96d)
	wildRespin96e.Describe(wildRespin96eID, "generate fake respin left (3) - RTP 96")
	wildRespin96f := comp.NewGenerateShapeAction(fakeRespin2Chance96, fakeRespinGrid2, wildRespinCenter, wildRespin).WithAlternate(wildRespin96e)
	wildRespin96f.Describe(wildRespin96fID, "generate fake respin left (2) - RTP 96")
	wildRespin96g := comp.NewGenerateShapeAction(fakeRespin1Chance96, fakeRespinGrid1, wildRespinCenter, wildRespin).WithAlternate(wildRespin96f)
	wildRespin96g.WithTriggerFilters(comp.OnFirstSpin) // leave this in place as it must not be called during refills!!!
	wildRespin96g.Describe(wildRespin96gID, "generate fake respin left (1) - RTP 96")

	firstWild96 := comp.NewGenerateSymbolAction(wild, firstWildChances96).GenerateNoDupes()
	firstWild96.WithTriggerFilters(comp.OnFirstSpin) // leave this in place as it must not be called during refills!!!
	firstWild96.Describe(firstWild96ID, "generate wilds - first spin - RTP 96")

	freeWild96 := comp.NewGenerateSymbolAction(wild, freeWildChances96, scatterReels...)
	freeWild96.Describe(freeWild96ID, "generate wilds - free spins - RTP 96")

	firstScatter96 := comp.NewGenerateSymbolAction(scatter, scatterRegularChances96, scatterReels...).GenerateNoDupes()
	firstScatter96.WithTriggerFilters(comp.OnFirstSpin) // leave this in place as it must not be called during refills!!!
	firstScatter96.Describe(firstScatter96ID, "generate scatters - first spin - RTP 96")

	freeScatter96a := comp.NewGenerateSymbolAction(scatter, scatterFreeSpinChances96[0], scatterReels...).GenerateNoDupes()
	freeScatter96a.Describe(freeScatter96aID, "generate scatters - free spin band 1 - RTP 96")
	freeScatter96b := comp.NewGenerateSymbolAction(scatter, scatterFreeSpinChances96[1], scatterReels...).GenerateNoDupes()
	freeScatter96b.Describe(freeScatter96bID, "generate scatters - free spin band 2 - RTP 96")
	freeScatter96c := comp.NewGenerateSymbolAction(scatter, scatterFreeSpinChances96[2], scatterReels...).GenerateNoDupes()
	freeScatter96c.Describe(freeScatter96cID, "generate scatters - free spin band 3 - RTP 96")
	freeScatter96d := comp.NewGenerateSymbolAction(scatter, scatterFreeSpinChances96[3], scatterReels...).GenerateNoDupes()
	freeScatter96d.Describe(freeScatter96dID, "generate scatters - free spin band 4 - RTP 96")
	freeScatter96e := comp.NewGenerateSymbolAction(scatter, scatterFreeSpinChances96[4], scatterReels...).GenerateNoDupes()
	freeScatter96e.Describe(freeScatter96eID, "generate scatters - free spin band 5 - RTP 96")

	freeScatter96 := comp.NewMultiActionFlagValue(flagFreeSpinBands, 1, freeScatter96a, 2, freeScatter96b, 3, freeScatter96c, 4, freeScatter96d, 5, freeScatter96e)
	freeScatter96.Describe(freeScatter96ID, "multi-select generate scatters - RTP 96")

	// generate 2x wilds on free spins
	freeWild2a := comp.NewGenerateSymbolAction(wild2, freeWild2Chances1, scatterReels...)
	freeWild2a.Describe(freeWild2aID, "generate 2x wild - free spin band 1")
	freeWild2b := comp.NewGenerateSymbolAction(wild2, freeWild2Chances2, scatterReels...)
	freeWild2b.Describe(freeWild2bID, "generate 2x wild - free spin band 2")
	freeWild2c := comp.NewGenerateSymbolAction(wild2, freeWild2Chances3, scatterReels...)
	freeWild2c.Describe(freeWild2cID, "generate 2x wild - free spin band 3")
	freeWild2d := comp.NewGenerateSymbolAction(wild2, freeWild2Chances4, scatterReels...)
	freeWild2d.Describe(freeWild2dID, "generate 2x wild - free spin band 4")
	freeWild2e := comp.NewGenerateSymbolAction(wild2, freeWild2Chances5, scatterReels...)
	freeWild2e.Describe(freeWild2eID, "generate 2x wild - free spin band 5")

	freeWild2 := comp.NewMultiActionFlagValue(flagFreeSpinBands, 1, freeWild2a, 2, freeWild2b, 3, freeWild2c, 4, freeWild2d, 5, freeWild2e)
	freeWild2.WithTriggerFilters(comp.OnRoundFlagValue(flagWildRespin, 1))
	freeWild2.Describe(freeWild2ID, "multi-select generate 2x wild")

	// generate 3x wilds on free spins
	freeWild3a := comp.NewGenerateSymbolAction(wild3, freeWild3Chances1, scatterReels...)
	freeWild3a.Describe(freeWild3aID, "generate 3x wild - free spin band 1")
	freeWild3b := comp.NewGenerateSymbolAction(wild3, freeWild3Chances2, scatterReels...)
	freeWild3b.Describe(freeWild3bID, "generate 3x wild - free spin band 2")
	freeWild3c := comp.NewGenerateSymbolAction(wild3, freeWild3Chances3, scatterReels...)
	freeWild3c.Describe(freeWild3cID, "generate 3x wild - free spin band 3")
	freeWild3d := comp.NewGenerateSymbolAction(wild3, freeWild3Chances4, scatterReels...)
	freeWild3d.Describe(freeWild3dID, "generate 3x wild - free spin band 4")
	freeWild3e := comp.NewGenerateSymbolAction(wild3, freeWild3Chances5, scatterReels...)
	freeWild3e.Describe(freeWild3eID, "generate 3x wild - free spin band 5")

	freeWild3 := comp.NewMultiActionFlagValue(flagFreeSpinBands, 1, freeWild3a, 2, freeWild3b, 3, freeWild3c, 4, freeWild3d, 5, freeWild3e)
	freeWild3.WithTriggerFilters(comp.OnRoundFlagValue(flagWildRespin, 1))
	freeWild3.Describe(freeWild3ID, "multi-select generate 3x wild")

	// generate 4x wilds on free spins
	freeWild4a := comp.NewGenerateSymbolAction(wild4, freeWild4Chances1, scatterReels...)
	freeWild4a.Describe(freeWild4aID, "generate 4x wild - free spin band 1")
	freeWild4b := comp.NewGenerateSymbolAction(wild4, freeWild4Chances2, scatterReels...)
	freeWild4b.Describe(freeWild4bID, "generate 4x wild - free spin band 2")
	freeWild4c := comp.NewGenerateSymbolAction(wild4, freeWild4Chances3, scatterReels...)
	freeWild4c.Describe(freeWild4cID, "generate 4x wild - free spin band 3")
	freeWild4d := comp.NewGenerateSymbolAction(wild4, freeWild4Chances4, scatterReels...)
	freeWild4d.Describe(freeWild4dID, "generate 4x wild - free spin band 4")
	freeWild4e := comp.NewGenerateSymbolAction(wild4, freeWild4Chances5, scatterReels...)
	freeWild4e.Describe(freeWild4eID, "generate 4x wild - free spin band 5")

	freeWild4 := comp.NewMultiActionFlagValue(flagFreeSpinBands, 1, freeWild4a, 2, freeWild4b, 3, freeWild4c, 4, freeWild4d, 5, freeWild4e)
	freeWild4.WithTriggerFilters(comp.OnRoundFlagValue(flagWildRespin, 1))
	freeWild4.Describe(freeWild4ID, "multi-select generate 4x wild")

	// all paylines.
	paylines := comp.NewAllPaylinesAction(true)
	paylines.Describe(paylinesID, "all paylines")

	// reduce payouts by band.
	reducePayouts1 := comp.NewRemovePayoutBandsAction(3, direction, true, true, reduceBandsFirst)
	reducePayouts1.WithTriggerFilters(comp.OnFirstSpin) // leave this in place as it must not be called during refills!!!
	reducePayouts1.Describe(removePayouts1ID, "remove payouts by band (first spins)")

	reducePayouts2 := comp.NewRemovePayoutBandsAction(3, direction, false, true, reduceBandsFree1)
	reducePayouts2.WithTriggerFilters(comp.OnRoundFlagBelow(flagFreeSpinCount, 2))
	reducePayouts2.Describe(removePayouts2ID, "remove payouts by band (first few free spins)")

	reducePayouts3 := comp.NewRemovePayoutBandsAction(3, direction, false, true, reduceBandsFree2)
	reducePayouts3.WithTriggerFilters(comp.OnRoundFlagAbove(flagFreeSpinCount, 2))
	reducePayouts3.Describe(removePayouts3ID, "remove payouts by band (free spins)")

	// grid teasers.
	teaser1 := comp.NewFakeSpinAction(fakeSpinChance, teasers1...)
	teaser1.Describe(teaser1ID, "teaser full page (no wilds)")

	teaser2 := comp.NewFakeSpinAction(fakeSpinChance, teasers2...)
	teaser2.Describe(teaser2ID, "teaser full page (1 wild, reel 1+6)")

	teaser3 := comp.NewFakeSpinAction(fakeSpinChance, teasers3...)
	teaser3.Describe(teaser3ID, "teaser full page (2 wilds, reel 1+6)")

	teaser4 := comp.NewFakeSpinAction(fakeSpinChance, teasers4...)
	teaser4.Describe(teaser4ID, "teaser full page (3 wilds, reel 1+6)")

	teaser5 := comp.NewFakeSpinAction(fakeSpinChance, teasers5...)
	teaser5.Describe(teaser5ID, "teaser full page (4 wilds, reel 1+6)")

	teaserFullPage := comp.NewMultiActionGridCount(wild, []int{1, 6}, 0, teaser1, 1, teaser2, 2, teaser3, 3, teaser4, 4, teaser5)
	teaserFullPage.Describe(teaserFPID, "multi-select teaser full page")

	teaserFunc := func(spin *comp.Spin) bool {
		if spin.TestChance2(2.0) { // 1 in ~50
			teaserReplace1(spin)
			return true
		}
		return false
	}

	teaser3or4 := comp.NewCustomAction(comp.AwardBonuses, comp.GridModified, teaserFunc, nil).WithAlternate(teaserFullPage)
	teaser3or4.WithTriggerFilters(comp.OnFirstSpin, comp.OnZeroPayouts(), comp.OnGridNotContains(scatter, scatterReels2...))
	teaser3or4.Describe(teaserNoScattersID, "custom teaser page (no scatters)")

	scattersWildsFunc := func(spin *comp.Spin) bool {
		if spin.TestChance2(10.0) { // 1 in ~10
			teaserReplace2(spin)
			return true
		}
		return false
	}

	teaserScattersWilds := comp.NewCustomAction(comp.AwardBonuses, comp.GridModified, scattersWildsFunc, nil)
	teaserScattersWilds.WithTriggerFilters(comp.OnFirstSpin, comp.OnZeroPayouts(), comp.OnGridCounts(scatter, []int{3, 4}, scatterReels2...), comp.OnGridNotContains(wild))
	teaserScattersWilds.Describe(teaserScattersWildsID, "custom teaser wilds (3 or 4 scatters)")

	// inject wilds teaser
	teaserZeroPayWilds := comp.NewCustomAction(comp.AwardBonuses, comp.GridModified, teaserReplace3, nil)
	teaserZeroPayWilds.WithTriggerFilters(comp.OnFirstSpin, comp.OnZeroPayouts(), comp.OnGridNotContains(scatter, scatterReels2...), comp.OnGridNotContains(wild))
	teaserZeroPayWilds.Describe(teaserZeroPayWildsID, "custom teaser wilds (no payouts, no scatters, no wilds)")

	// award free spins.
	award5 := comp.NewScatterFreeSpinsAction(5, false, scatter, 3, false)
	award5.Describe(award5ID, "award 5 free spins - first spin - 3 scatters")
	award8 := comp.NewScatterFreeSpinsAction(8, false, scatter, 4, false).WithAlternate(award5)
	award8.WithTriggerFilters(comp.OnFirstSpin) // leave this in place as it must not be called during refills!!!
	award8.Describe(award8ID, "award 8 free spins - first spin- 4 scatters")

	retrigger1 := comp.NewScatterFreeSpinsAction(1, false, scatter, 1, false)
	retrigger1.Describe(retrigger1ID, "award 1 free spin - free spin - 1 scatter")
	retrigger2 := comp.NewScatterFreeSpinsAction(2, false, scatter, 2, false).WithAlternate(retrigger1)
	retrigger2.Describe(retrigger2ID, "award 2 free spins - free spin - 2 scatters")
	retrigger3 := comp.NewScatterFreeSpinsAction(3, false, scatter, 3, false).WithAlternate(retrigger2)
	retrigger3.Describe(retrigger3ID, "award 3 free spins - free spin - 3 scatters")
	retrigger4 := comp.NewScatterFreeSpinsAction(4, false, scatter, 4, false).WithAlternate(retrigger3)
	retrigger4.Describe(retrigger4ID, "award 4 free spins - free spin - 4 scatters")

	// mark refill spin.
	refill := comp.NewShapeRefillAction(wild, wildRespinGrid, wildRespinCenter, false)
	refill.WithTriggerFilters(comp.OnFirstSpin) // leave this in place as it must not be called during refills!!!
	refill.WithClearFilters(comp.OnNoFreeSpins)
	refill.Describe(refillID, "wild respin - refill")

	// update round flag 2 if the wild respin triggered.
	wildRespinFlag := comp.NewRoundFlagShapeDetect(flagWildRespin, wild, wildRespinGrid)
	wildRespinFlag.WithTriggerFilters(comp.OnFirstSpin) // leave this in place as it must not be called during refills!!!
	wildRespinFlag.Describe(wildRespinFlagID, "mark wild respin (flag 2)")

	// mark sticky wilds.
	stickies1 := comp.NewStickySymbolAction(wild, 1, 6)
	stickies1.WithTriggerFilters(comp.OnFirstSpin, comp.OnRoundFlagValue(flagWildRespin, 1))
	stickies1.Describe(stickies1ID, "make wild symbols sticky (first spin)")
	stickies2 := comp.NewStickySymbolsAction(wild, wild2, wild3, wild4)
	stickies2.Describe(stickies2ID, "make wild symbols sticky (free spins)")

	// round flag 0 for payout bands.
	payoutBands := comp.NewRoundFlagWeightedAction(flagFreeSpinBands, freeSpinBands)
	payoutBands.WithTriggerFilters(comp.OnFirstSpin) // leave this in place as it must not be called during refills!!!
	payoutBands.Describe(payoutBandsID, "weighted payout band (flag 0)")

	// update round flag 3 marking sequence of free spin.
	freeCounter = comp.NewRoundFlagIncreaseAction(flagFreeSpinCount)
	freeCounter.Describe(freeCounterID, "count number of free spins (flag 1)")

	// force wild respin for scripted rounds.
	forceWildRespin := comp.NewGenerateShapeAction(100, wildRespinGrid, wildRespinCenter, wildRespin)
	forceWildRespin.Describe(forceWildRespinID, "force wild respin (scripted round)")

	// force payout band for scripted rounds.
	forceBand5 := comp.NewRoundFlagWeightedAction(flagFreeSpinBands, freeSpinBand5)
	forceBand5.Describe(forceBand5ID, "force payout band 5 (flag 0; scripted round)")

	// force no payouts for scripted rounds.
	forceNoPayout := comp.NewRemovePayoutsAction(0, math.MaxFloat64, 100, direction, 3, false, true)
	forceNoPayout.WithTriggerFilters(comp.OnFirstSpin) // leave this in place as it must not be called during refills!!!
	forceNoPayout.Describe(forceNoPayoutsID, "remove all payouts (first spin; scripted rounds)")

	// force scrolls for scripted rounds.
	forceScrolls4 := comp.NewGenerateSymbolAction(scatter, []float64{100, 100, 100, 100}, scatterReels...).GenerateNoDupes()
	forceScrolls4.Describe(forceScrolls4ID, "force 4 scatters (scripted round)")

	forceScrolls4S6 := comp.NewGenerateSymbolAction(scatter, []float64{100, 100, 100, 100}, scatterReels...).GenerateNoDupes()
	forceScrolls4S6.WithTriggerFilters(comp.OnSpinSequence(6))
	forceScrolls4S6.Describe(forceScrolls4S6ID, "force 4 scatters (spin 6; scripted round)")

	forceScrollsS10Up := comp.NewGenerateSymbolAction(scatter, []float64{20, 90, 100, 40}, scatterReels...).GenerateNoDupes()
	forceScrollsS10Up.WithTriggerFilters(comp.OnSpinSequenceAbove(9))
	forceScrollsS10Up.Describe(forceScrollsS10UpID, "force 3+ scatters (spin 10+; scripted round)")

	forceScrollsS2to8 := comp.NewGenerateSymbolAction(scatter, []float64{75, 15, 2}, scatterReels...).GenerateNoDupes()
	forceScrollsS2to8.WithTriggerFilters(comp.OnSpinSequenceBelow(9))
	forceScrollsS2to8.Describe(forceScrollsS2to8ID, "force 1 or 2 scatters (spin 2-8; scripted round)")

	forceScrollsS9Up := comp.NewGenerateSymbolAction(scatter, []float64{100, 10}, scatterReels...).GenerateNoDupes()
	forceScrollsS9Up.WithTriggerFilters(comp.OnSpinSequenceAbove(8))
	forceScrollsS9Up.Describe(forceScrollsS9UpID, "force scatters (spin 9+; scripted round)")

	// force full page of chalice
	forceChaliceFP := comp.NewGenerateShapeAction(100, fullPageR1to4shape, fullPageR1offset, fullPageChaliceWeights).AllowPrevious()
	forceChaliceFP.WithTriggerFilters(comp.OnSpinSequence(2))
	forceChaliceFP.Describe(forceChaliceFPID, "force full page of chalice (spin 2, scripted round)")

	// force full page of potion.
	forcePotionFP := comp.NewGenerateShapeAction(100, fullPageR1to4shape, fullPageR1offset, fullPagePotionWeights).AllowPrevious()
	forcePotionFP.WithTriggerFilters(comp.OnSpinSequence(3))
	forcePotionFP.Describe(forcePotionFPID, "force full page of potion (spin 3, scripted round)")

	// force full page of sword.
	forceDaggerFP := comp.NewGenerateShapeAction(100, fullPageR1to4shape, fullPageR1offset, fullPageDaggerWeights).AllowPrevious()
	forceDaggerFP.WithTriggerFilters(comp.OnSpinSequence(4))
	forceDaggerFP.Describe(forceDaggerFPID, "force full page of dagger (spin 4, scripted round)")

	// force full page of book.
	forceBookFP := comp.NewGenerateShapeAction(100, fullPageR1to4shape, fullPageR1offset, fullPageBookWeights).AllowPrevious()
	forceBookFP.WithTriggerFilters(comp.OnSpinSequence(5))
	forceBookFP.Describe(forceBookFPID, "force full page of book (spin 5, scripted round)")

	// force full page of medusa.
	forceMedusaFP := comp.NewGenerateShapeAction(100, fullPageR1to4shape, fullPageR1offset, fullPageMedusaWeights).AllowPrevious()
	forceMedusaFP.WithTriggerFilters(comp.OnSpinSequence(7))
	forceMedusaFP.Describe(forceMedusaFPID, "force full page of medusa (spin 7, scripted round)")

	// force full page of death.
	forceDeathFP := comp.NewGenerateShapeAction(100, fullPageR1to4shape, fullPageR1offset, fullPageDeathWeights).AllowPrevious()
	forceDeathFP.WithTriggerFilters(comp.OnSpinSequence(8))
	forceDeathFP.Describe(forceDeathFPID, "force full page of death (spin 8, scripted round)")

	// force full page of succubus.
	forceSuccubusFP := comp.NewGenerateShapeAction(100, fullPageR1to4shape, fullPageR1offset, fullPageSuccubusWeights).AllowPrevious()
	forceSuccubusFP.WithTriggerFilters(comp.OnSpinSequence(9))
	forceSuccubusFP.Describe(forceSuccubusFPID, "force full page of succubus (spin 9, scripted round)")

	// force full page of demon goat.
	forceDemonGoatFP := comp.NewGenerateShapeAction(100, fullPageR1to4shape, fullPageR1offset, fullPageDemonGoatWeights).AllowPrevious()
	forceDemonGoatFP.WithTriggerFilters(comp.OnSpinSequence(11))
	forceDemonGoatFP.Describe(forceDemonGoatFPID, "force full page of demon goat (spin 11, scripted round)")

	// force many 1x wilds on first spin for scripted rounds.
	forceW1S1 := comp.NewGenerateSymbolAction(wild, []float64{100, 100, 100, 100, 100, 90, 70, 5}, 2, 3, 4, 5).AllowPrevious()
	forceW1S1.Describe(forceW1S1ID, "force many x1 wilds (spin 1; scripted round)")

	// force 3-4 random multiplier wilds for scripted rounds.
	weights4RW := util.AcquireWeighting().AddWeights(util.Indexes{wild, wild2, wild3, wild4}, []float64{12, 6, 2, 8})
	force4RW := comp.NewGenerateSymbolsAction(weights4RW, []float64{100, 100, 60, 30}, scatterReels...)
	force4RW.Describe(force4RWID, "force 2-4 wilds (scripted round)")

	weights4RW8Up := util.AcquireWeighting().AddWeights(util.Indexes{wild, wild2, wild3, wild4}, []float64{6, 4, 8, 2})
	force4RWS8Up := comp.NewGenerateSymbolsAction(weights4RW8Up, []float64{100, 100, 60, 30}, scatterReels...)
	force4RWS8Up.WithTriggerFilters(comp.OnSpinSequenceAbove(7))
	force4RWS8Up.Describe(force4RWS8UpID, "force 2-4 wilds (spin 8+; scripted round)")

	// force x4 wilds for scripted rounds
	forceW4S2R5 := comp.NewGenerateSymbolAction(wild4, []float64{100}, 5)
	forceW4S2R5.WithTriggerFilters(comp.OnSpinSequence(2))
	forceW4S2R5.Describe(forceW4S2R5ID, "force x4 wild (spin 2; reel 5; scripted round)")

	forceW4S3R4 := comp.NewGenerateSymbolAction(wild4, []float64{100}, 4)
	forceW4S3R4.WithTriggerFilters(comp.OnSpinSequence(3))
	forceW4S3R4.Describe(forceW4S3R4ID, "force x4 wild (spin 3; reel 4; scripted round)")

	forceW4S4R5 := comp.NewGenerateSymbolAction(wild4, []float64{100}, 5)
	forceW4S4R5.WithTriggerFilters(comp.OnSpinSequence(4))
	forceW4S4R5.Describe(forceW4S4R5ID, "force x4 wild (spin 4; reel 5; scripted round)")

	forceW4S5R4 := comp.NewGenerateSymbolAction(wild4, []float64{100}, 4)
	forceW4S5R4.WithTriggerFilters(comp.OnSpinSequence(5))
	forceW4S5R4.Describe(forceW4S5R4ID, "force x4 wild (spin 5; reel 4; scripted round)")

	forceW4S7R2 := comp.NewGenerateSymbolAction(wild4, []float64{100}, 2)
	forceW4S7R2.WithTriggerFilters(comp.OnSpinSequence(7))
	forceW4S7R2.Describe(forceW4S7R2ID, "force x4 wild (spin 7; reel 2; scripted round)")

	forceW4S8R3 := comp.NewGenerateSymbolAction(wild4, []float64{100}, 3)
	forceW4S8R3.WithTriggerFilters(comp.OnSpinSequence(8))
	forceW4S8R3.Describe(forceW4S8R3ID, "force x4 wild (spin 8; reel 3; scripted round)")

	forceW4S9UpRandom := comp.NewGenerateSymbolAction(wild4, []float64{90, 30, 8, 1}, scatterReels...)
	forceW4S9UpRandom.WithTriggerFilters(comp.OnSpinSequenceAbove(8))
	forceW4S9UpRandom.Describe(forceW4S9UpRandomID, "force x4 wild (spin 9+ random; scripted round)")

	// force full page of 4x wilds for scripted rounds.
	forceFullW4S9UpRandom := comp.NewGenerateShapeAction(30, fullPageR1to4shape, fullPageR1offset, fullPage4Xweights)
	forceFullW4S9UpRandom.WithTriggerFilters(comp.OnSpinSequenceAbove(8))
	forceFullW4S9UpRandom.Describe(forceW4S9UpRandomID, "force full page x4 wild (spin 9+; scripted round)")

	forceFullW4S12UpRandom := comp.NewGenerateShapeAction(70, fullPageR1to4shape, fullPageR1offset, fullPage4Xweights)
	forceFullW4S12UpRandom.WithTriggerFilters(comp.OnSpinSequenceAbove(11))
	forceFullW4S12UpRandom.Describe(forceW4S12UpRandomID, "force full page x4 wild (spin 12+; scripted round)")

	// scripted rounds - scenario 1 configuration (max payout).
	script1 := comp.NewScriptedRound(1, script1weight,
		comp.SpinActions{forceWildRespin, forceScrolls4, paylines, forceNoPayout, award8, wildRespinFlag, refill, stickies1, forceBand5},
		comp.SpinActions{forceW4S2R5, forceW4S3R4, forceW4S4R5, forceW4S5R4, forceScrolls4S6,
			forceW4S7R2, forceW4S8R3, forceW4S9UpRandom, forceScrollsS10Up,
			paylines, retrigger4, stickies2},
	)

	// scripted rounds - scenario 2 configuration (max payout).
	script2 := comp.NewScriptedRound(2, script2weight,
		comp.SpinActions{forceWildRespin, forceScrolls4, paylines, award8, wildRespinFlag, refill, stickies1, forceBand5},
		comp.SpinActions{force4RW, paylines, retrigger4, stickies2})

	// scripted rounds - scenario 3 configuration (max payout).
	script3 := comp.NewScriptedRound(3, script3weight,
		comp.SpinActions{forceWildRespin, forceScrolls4, paylines, award8, wildRespinFlag, refill, stickies1, forceBand5},
		comp.SpinActions{forceScrollsS2to8, force4RWS8Up, forceScrollsS9Up, paylines, retrigger4, stickies2})

	// scripted rounds - scenario 4 configuration (max payout).
	script4 := comp.NewScriptedRound(4, script4weight,
		comp.SpinActions{forceWildRespin, forceW1S1, forceScrolls4, paylines, award8, wildRespinFlag, refill, stickies1, forceBand5},
		comp.SpinActions{forceScrollsS2to8, force4RWS8Up, forceScrollsS9Up, paylines, retrigger4, stickies2})

	// scripted rounds - scenario 5 configuration (max payout).
	script5 := comp.NewScriptedRound(5, script5weight,
		comp.SpinActions{forceWildRespin, forceScrolls4, paylines, forceNoPayout, award8, wildRespinFlag, refill, stickies1, forceBand5},
		comp.SpinActions{forceScrollsS2to8, forceScrollsS9Up, forceFullW4S9UpRandom, paylines, retrigger4, stickies2})

	// scripted rounds - scenario 6 configuration (max payout).
	script6 := comp.NewScriptedRound(6, script6weight,
		comp.SpinActions{forceWildRespin, forceScrolls4, paylines, forceNoPayout, award8, wildRespinFlag, refill, stickies1, forceBand5},
		comp.SpinActions{forceChaliceFP, forcePotionFP, forceDaggerFP, forceBookFP, forceMedusaFP, forceDeathFP, forceSuccubusFP,
			forceScrolls4S6, forceScrollsS10Up, forceDemonGoatFP, forceFullW4S12UpRandom, paylines, retrigger4, stickies2})

	// scripted rounds - scenario 9 configuration (dummy to fill up the weighting).
	script9 := comp.NewScriptedRound(9, script9weight, nil, nil)

	// scripted rounds accumulated.
	scriptedRounds = comp.NewScriptedRoundSelector(0.1, script1, script2, script3, script4, script5, script6, script9)

	// gather it all together.
	actionsA92all := comp.SpinActions{freeCounter, wildRespin92g, firstWild92, freeWild92, firstScatter92, freeScatter92a, freeScatter92b, freeScatter92c, freeScatter92d, freeScatter92e, freeScatter92}
	actionsA92first := comp.SpinActions{wildRespin92g, firstWild92, freeWild2, freeWild3, freeWild4, firstScatter92}
	actionsA92free := comp.SpinActions{freeCounter, freeWild92, freeScatter92}

	actionsA94all := comp.SpinActions{freeCounter, wildRespin94g, firstWild94, freeWild94, firstScatter94, freeScatter94a, freeScatter94b, freeScatter94c, freeScatter94d, freeScatter94e, freeScatter94}
	actionsA94first := comp.SpinActions{wildRespin94g, firstWild94, freeWild2, freeWild3, freeWild4, firstScatter94}
	actionsA94free := comp.SpinActions{freeCounter, freeWild94, freeScatter94}

	actionsA96all := comp.SpinActions{freeCounter, wildRespin96g, firstWild96, freeWild96, firstScatter96, freeScatter96a, freeScatter96b, freeScatter96c, freeScatter96d, freeScatter96e, freeScatter96}
	actionsA96first := comp.SpinActions{wildRespin96g, firstWild96, freeWild2, freeWild3, freeWild4, firstScatter96}
	actionsA96free := comp.SpinActions{freeCounter, freeWild96, freeScatter96}

	actionsBall := comp.SpinActions{
		freeWild2a, freeWild2b, freeWild2c, freeWild2d, freeWild2e, freeWild2,
		freeWild3a, freeWild3b, freeWild3c, freeWild3d, freeWild3e, freeWild3,
		freeWild4a, freeWild4b, freeWild4c, freeWild4d, freeWild4e, freeWild4,
		paylines, reducePayouts1, reducePayouts2, reducePayouts3,
		teaser1, teaser2, teaser3, teaser4, teaser5, teaserFullPage, teaser3or4, teaserScattersWilds, teaserZeroPayWilds,
		award8, retrigger4, wildRespinFlag, refill, stickies1, stickies2, payoutBands,
	}
	actionsBfirst := comp.SpinActions{
		paylines, reducePayouts1, teaser3or4, teaserScattersWilds, teaserZeroPayWilds,
		award8, wildRespinFlag, refill, stickies1, payoutBands,
	}
	actionsBfree := comp.SpinActions{
		freeWild2, freeWild3, freeWild4,
		paylines, reducePayouts2, reducePayouts3,
		retrigger4, stickies2,
	}

	forceAll := comp.SpinActions{
		forceWildRespin, forceBand5, forceNoPayout,
		forceScrolls4, forceScrolls4S6, forceScrollsS10Up, forceScrollsS2to8, forceScrollsS9Up,
		forceChaliceFP, forcePotionFP, forceDaggerFP, forceBookFP, forceMedusaFP, forceDeathFP, forceSuccubusFP, forceDemonGoatFP,
		forceW1S1, force4RW, force4RWS8Up,
		forceW4S2R5, forceW4S3R4, forceW4S4R5, forceW4S5R4, forceW4S7R2, forceW4S8R3, forceW4S9UpRandom,
		forceFullW4S9UpRandom, forceFullW4S12UpRandom}

	actions92all = append(append(actionsA92all, actionsBall...), forceAll...)
	actions92first = append(actionsA92first, actionsBfirst...)
	actions92free = append(actionsA92free, actionsBfree...)
	actions94all = append(append(actionsA94all, actionsBall...), forceAll...)
	actions94first = append(actionsA94first, actionsBfirst...)
	actions94free = append(actionsA94free, actionsBfree...)
	actions96all = append(append(actionsA96all, actionsBall...), forceAll...)
	actions96first = append(actionsA96first, actionsBfirst...)
	actions96free = append(actionsA96free, actionsBfree...)
}

var (
	validOffsets = []int{0, 1, 4, 5, 6, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 20, 21}
	symWeighting = util.AcquireWeightingDedup4().AddWeights(util.Indexes{1, 2, 3, 4, 5, 6, 7, 8}, []float64{1, 1, 1, 1, 2, 3, 4, 5})
)

// special teaser with almost full page of 3 or 4 different symbols.
func teaserReplace1(spin *comp.Spin) {
	indexes, prng := spin.Indexes(), spin.PRNG()
	clearReel2 := spin.TestChance2(66.0)
	threeSymbols := spin.TestChance2(66.0)

	ss := make(util.Indexes, 3, 4)
	if !threeSymbols {
		ss = ss[:4]
	}
	symWeighting.FillRandom(prng, len(ss), ss)

	var min, max int
	if clearReel2 {
		min, max = 4, 7
	} else {
		min, max = 8, 12
	}
	for offset := min; offset < max; offset++ {
		indexes[offset] = ss[0]
	}

	if clearReel2 {
		min, max = 8, 12
	} else {
		min, max = 4, 7
	}
	for offset := min; offset < max; offset++ {
		indexes[offset] = ss[1+prng.IntN(len(ss)-1)]
	}

	for _, offset := range validOffsets {
		if offset < 4 || offset >= 12 {
			if indexes[offset] != wild {
				ix := prng.IntN(len(ss))
				indexes[offset] = ss[ix]
			}
		}
	}
}

// special teaser inject wilds on reel 1 on 3 or 4 scatters, 0 wilds and no paylines.
func teaserReplace2(spin *comp.Spin) {
	indexes, reel2, prng := spin.Indexes(), spin.Reels()[1], spin.PRNG()

	indexes[0] = wild
	indexes[1] = wild

	dupe := func(symbol util.Index) bool {
		for offs := 8; offs < 12; offs++ {
			if indexes[offs] == symbol {
				return true
			}
		}
		return false
	}

	for offs := 4; offs < 7; offs++ {
		for dupe(indexes[offs]) {
			indexes[offs] = reel2.RandomIndex(prng)
		}
	}
}

// special teaser inject 1 or 3-6 wilds on reel 2, 3, 4 or 5, making sure to prevent new paylines.
// this should only be called for the first spin when there are no wilds and no scatters!
// chance 1 in ~10 on a single wild.
// change 1 in ~10 on 3-6 wilds.
// which means 1 in ~5 will inject at least 1 wild!
func teaserReplace3(spin *comp.Spin) bool {
	indexes, prng := spin.Indexes(), spin.PRNG()

	chance := prng.IntN(10000)
	if chance >= 3000 {
		return false
	}

	allowed := []byte{12, 13, 14, 15, 16, 17, 18}
	switch {
	case !dupesOnReel(indexes[0:2], indexes[8:12]):
		allowed = append(allowed, 4, 5, 6)
	case !dupesOnReel(indexes[0:2], indexes[4:7]):
		allowed = append(allowed, 8, 9, 10, 11)
	}

	count := 1
	if chance >= 2000 {
		count = 3 + prng.IntN(4)
	}

	l := len(allowed)
	for count > 0 {
		ix := allowed[prng.IntN(l)]
		for indexes[ix] == wild {
			ix = allowed[prng.IntN(l)]
		}

		indexes[ix] = wild
		count--
	}

	return true
}

func dupesOnReel(a, b util.Indexes) bool {
	for ix := range a {
		for iy := range b {
			if b[iy] == a[ix] {
				return true
			}
		}
	}
	return false
}

func initSlots(target float64, weights [symbolCount]comp.SymbolOption, actions1, actions2 []comp.SpinActioner) *comp.Slots {
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
		comp.WithMask(gridMask...),
		comp.WithSymbols(s1),
		comp.PayDirections(direction),
		comp.MaxPayout(maxPayout),
		comp.WithRTP(target),
		comp.WithRoundFlags(flags...),
		comp.WithActions(actions1, actions2, nil, nil),
		comp.WithScriptedRoundSelector(scriptedRounds),
	)
}

func init() {
	initActions()

	slots92 = initSlots(92.0, weights92, actions92first, actions92free)
	slots94 = initSlots(94.0, weights94, actions94first, actions94free)
	slots96 = initSlots(96.0, weights96, actions96first, actions96free)

	slots92params = game.RegularParams{Slots: slots92}
	slots94params = game.RegularParams{Slots: slots94}
	slots96params = game.RegularParams{Slots: slots96}
}

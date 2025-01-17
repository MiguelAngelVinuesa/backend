package fpr

import (
	"fmt"
	"math"

	comp "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	game "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/slots"
	util "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

const (
	symbolCount   = 13
	reels         = 6
	rows          = 4
	direction     = comp.PayLTR
	maxPayout     = 5000.0
	bonusBuyCost1 = 100.0
	bonusBuyCost2 = 500.0

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

	bonusKind1 = 1
	bonusKind2 = 2

	// always add new ID's with a new unique number!
	wildRespin92aID     = 101
	wildRespin92bID     = 102
	wildRespin92cID     = 103
	wildRespin92dID     = 104
	wildRespin92eID     = 105
	wildRespin92fID     = 106
	wildRespin92gID     = 107
	firstWild92ID       = 111
	freeWild92aID       = 112
	freeWild92bID       = 113
	freeWild92bb1ID     = 115
	freeWild92bb2ID     = 116
	firstScatter92aID   = 121
	firstScatter92bID   = 122
	freeScatter92a1ID   = 123 // 124,125,126,127
	freeScatter92aID    = 129
	freeScatter92b1ID   = 133 // 134,135,136,137
	freeScatter92bID    = 139
	freeScatter92bb1aID = 143 // 144,145,146,147
	freeScatter92bb1ID  = 149
	freeScatter92bb2aID = 153 // 154,155,156,157
	freeScatter92bb2ID  = 159

	wildRespin94aID     = 201
	wildRespin94bID     = 202
	wildRespin94cID     = 203
	wildRespin94dID     = 204
	wildRespin94eID     = 205
	wildRespin94fID     = 206
	wildRespin94gID     = 207
	firstWild94ID       = 211
	freeWild94aID       = 212
	freeWild94bID       = 213
	freeWild94bb1ID     = 215
	freeWild94bb2ID     = 216
	firstScatter94aID   = 221
	firstScatter94bID   = 222
	freeScatter94a1ID   = 223 // 224,225,226,227
	freeScatter94aID    = 229
	freeScatter94b1ID   = 233 // 234,235,236,237
	freeScatter94bID    = 239
	freeScatter94bb1aID = 243 // 244,245,246,247
	freeScatter94bb1ID  = 249
	freeScatter94bb2aID = 253 // 254,255,256,257
	freeScatter94bb2ID  = 259

	wildRespin96aID     = 301
	wildRespin96bID     = 302
	wildRespin96cID     = 303
	wildRespin96dID     = 304
	wildRespin96eID     = 305
	wildRespin96fID     = 306
	wildRespin96gID     = 307
	firstWild96ID       = 311
	freeWild96aID       = 312
	freeWild96bID       = 313
	freeWild96bb1ID     = 315
	freeWild96bb2ID     = 316
	firstScatter96aID   = 321
	firstScatter96bID   = 322
	freeScatter96a1ID   = 323 // 324,325,326,327
	freeScatter96aID    = 329
	freeScatter96b1ID   = 333 // 334,335,336,337
	freeScatter96bID    = 339
	freeScatter96bb1aID = 343 // 344,345,347,347
	freeScatter96bb1ID  = 349
	freeScatter96bb2aID = 353 // 354,355,356,357
	freeScatter96bb2ID  = 359

	freeWildRespin1ID = 501 // 502,503,504,505
	freeWildRespin2ID = 509
	freeWildSuper1ID  = 511 // 532,533,534,535
	freeWildSuper2ID  = 519

	paylinesID            = 601
	removePayoutsFirstID  = 602
	removePayoutsRespinID = 603
	removePayoutsFree1ID  = 604
	removePayoutsFree2ID  = 605
	teaser1ID             = 611
	teaser2ID             = 612
	teaser3ID             = 613
	teaser4ID             = 614
	teaser5ID             = 615
	teaserFPID            = 618
	teaserNoScattersID    = 619
	teaserScattersWildsID = 620
	teaserZeroPayWildsID  = 621
	award5ID              = 631
	award8ID              = 632
	retrigger1ID          = 633
	retrigger2ID          = 634
	retrigger3ID          = 635
	retrigger4ID          = 636
	refillID              = 641
	stickies1ID           = 642
	stickies2ID           = 643
	payoutBandsID         = 651
	freeCounterID         = 652
	wildRespinFlagID      = 653

	bonusBuy1ID      = 801
	bonusBuy2ID      = 802
	bonusExtra1ID    = 803
	bonusExtra2ID    = 804
	bonusBuy2wildsID = 805
	payoutBandsBB1ID = 811
	payoutBandsBB2ID = 812

	forceWildRespinID        = 900
	forceBand5ID             = 905
	forceNoPayoutsID         = 906
	forceScrolls4ID          = 911
	forceScrolls4S6ID        = 912
	forceScrollsS10UpID      = 913
	forceScrollsS9UpID       = 914
	forceScrollsS2to8ID      = 915
	forceDiamondsFPID        = 920
	forceClubsFPID           = 921
	forceHeartsFPID          = 922
	forceSpadesFPID          = 923
	forceMedusaFPID          = 924
	forceDeathFPID           = 925
	forceSuccubusFPID        = 926
	forceDemonGoatFPID       = 927
	forceW1S1ID              = 950
	force4RWID               = 951
	force4RWS8UpID           = 952
	forceW4S2R5ID            = 961
	forceW4S3R4ID            = 962
	forceW4S4R5ID            = 963
	forceW4S5R4ID            = 964
	forceW4S7R2ID            = 965
	forceW4S8R3ID            = 966
	forceW4S6UpRandomID      = 991
	forceW4S9UpRandomID      = 992
	forceFullW4S6UpRandomID  = 995
	forceFullW4S9UpRandomID  = 996
	forceFullW4S12UpRandomID = 997

	flagFreeSpinBands = 0
	flagFreeSpinCount = 1
	flagWildRespin    = 2
	flagBonusBuy      = 3
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
	wildRespinChance92     = 1.69
	fakeRespin1Chance92    = 0.51
	fakeRespin2Chance92    = 0.52
	fakeRespin3Chance92    = 0.53
	fakeRespin4Chance92    = 0.39
	fakeRespin5Chance92    = 0.4
	fakeRespin6Chance92    = 0.41
	firstWildChances92     = []float64{20, 14, 1}
	freeWildChances92a     = []float64{30.7, 21.4, 18, 9, 1}
	freeWildChances92b     = []float64{22.8, 6.4}
	freeWildChances92bb1   = []float64{30.8, 21.8, 18, 9, 1}
	freeWildChances92bb2   = []float64{22.7, 6.2}
	firstScatterChances92a = []float64{25, 18, 11.5, 27}
	firstScatterChances92b = []float64{35, 26, 21, 7}
	freeScatterChances92a  = [5][]float64{
		{12},
		{18, 3},
		{22, 10, 2},
		{26, 14, 3, 1},
		{30, 17, 12, 5, 1}}
	freeScatterChances92b = [5][]float64{
		{1},
		{4},
		{7, 0.5},
		{10, 1.3},
		{14, 3, 0.5}}
	freeScatterChances92bb1 = [5][]float64{
		{12},
		{18, 3},
		{22, 10, 2},
		{26, 14, 3, 1},
		{30, 17, 12, 5, 1}}
	freeScatterChances92bb2 = [5][]float64{
		{1},
		{4},
		{7, 0.5},
		{10, 1.3},
		{14, 3, 0.5}}

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
	wildRespinChance94     = 1.7
	fakeRespin1Chance94    = 0.51
	fakeRespin2Chance94    = 0.52
	fakeRespin3Chance94    = 0.53
	fakeRespin4Chance94    = 0.39
	fakeRespin5Chance94    = 0.4
	fakeRespin6Chance94    = 0.41
	firstWildChances94     = []float64{20, 15, 1}
	freeWildChances94a     = []float64{31, 23, 18, 9, 1}
	freeWildChances94b     = []float64{23, 7}
	freeWildChances94bb1   = []float64{31, 23, 18, 9, 1}
	freeWildChances94bb2   = []float64{23, 7}
	firstScatterChances94a = []float64{25, 18, 11.5, 27}
	firstScatterChances94b = []float64{35, 26, 21, 7}
	freeScatterChances94a  = [5][]float64{
		{12},
		{18, 3},
		{22, 10, 2},
		{26, 14, 3, 1},
		{30, 17, 12, 5, 1}}
	freeScatterChances94b = [5][]float64{
		{1},
		{4},
		{7, 0.5},
		{10, 1.3},
		{14, 3, 0.5}}
	freeScatterChances94bb1 = [5][]float64{
		{12},
		{18, 3},
		{22, 10, 2},
		{26, 14, 3, 1},
		{30, 17, 12, 5, 1}}
	freeScatterChances94bb2 = [5][]float64{
		{1},
		{4},
		{7, 0.5},
		{10, 1.3},
		{14, 3, 0.5}}

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
	wildRespinChance96     = 1.72
	fakeRespin1Chance96    = 0.51
	fakeRespin2Chance96    = 0.52
	fakeRespin3Chance96    = 0.53
	fakeRespin4Chance96    = 0.39
	fakeRespin5Chance96    = 0.4
	fakeRespin6Chance96    = 0.41
	firstWildChances96     = []float64{21, 15, 1}
	freeWildChances96a     = []float64{32, 23.3, 18, 9, 1}
	freeWildChances96b     = []float64{25, 7.3}
	freeWildChances96bb1   = []float64{31.7, 23.4, 18, 9, 1}
	freeWildChances96bb2   = []float64{25.4, 7.6}
	firstScatterChances96a = []float64{25, 18, 11.5, 27}
	firstScatterChances96b = []float64{35, 26, 21, 7}
	freeScatterChances96a  = [5][]float64{
		{12},
		{18, 3},
		{22, 10, 2},
		{26, 15, 3, 1},
		{30, 17, 12, 5, 1}}
	freeScatterChances96b = [5][]float64{
		{1},
		{4},
		{7, 0.5},
		{11, 1.3},
		{14, 3.3, 0.5}}
	freeScatterChances96bb1 = [5][]float64{
		{12},
		{18, 3},
		{22, 10, 2},
		{26, 14, 3, 1},
		{30, 17, 12, 5, 1}}
	freeScatterChances96bb2 = [5][]float64{
		{1},
		{4},
		{7, 0.5},
		{10, 1.3},
		{14, 3, 0.5}}
)

var (
	// free spin payout band weights.
	freeSpinBandWeights    = []float64{20, 150, 350, 30, 4}
	freeSpinBandWeightsBB1 = []float64{1, 150, 350, 30, 5}
	freeSpinBandWeightsBB2 = []float64{1, 250, 350, 20, 6}

	// weights to remove payouts based on total payout factor during first spins.
	reduceBandsFirst = []comp.RemovePayoutBand{
		{MinPayout: 0, MaxPayout: 1, RemoveChance: 10},
		{MinPayout: 1, MaxPayout: 2.5, RemoveChance: 45},
		{MinPayout: 2.5, MaxPayout: 5, RemoveChance: 65},
		{MinPayout: 5, MaxPayout: 10, RemoveChance: 50},
		{MinPayout: 10, MaxPayout: 40, RemoveChance: 70},
		{MinPayout: 40, MaxPayout: 125, RemoveChance: 75},
		{MinPayout: 125, MaxPayout: 200, RemoveChance: 85},
		{MinPayout: 200, MaxPayout: maxPayout, RemoveChance: 95},
	}

	// weights to remove payouts based on total payout factor during wild respin.
	reduceBandsRespin = []comp.RemovePayoutBand{
		{MinPayout: 0.0, MaxPayout: 5, RemoveChance: 35},
		{MinPayout: 5, MaxPayout: 10, RemoveChance: 42},
		{MinPayout: 10, MaxPayout: maxPayout, RemoveChance: 71},
		{MinPayout: 30, MaxPayout: maxPayout, RemoveChance: 88},
	}

	// weights to remove payouts based on total payout factor during first few free spins.
	reduceBandsFree1 = []comp.RemovePayoutBand{
		{MinPayout: 0.0, MaxPayout: 10, RemoveChance: 12},
		{MinPayout: 10, MaxPayout: 30, RemoveChance: 20},
		{MinPayout: 30, MaxPayout: maxPayout, RemoveChance: 40},
	}

	// weights to remove payouts based on total payout factor during remaining free spins.
	reduceBandsFree2 = []comp.RemovePayoutBand{
		{MinPayout: 0, MaxPayout: 10, RemoveChance: 12},
		{MinPayout: 10, MaxPayout: 30, RemoveChance: 15},
		{MinPayout: 30, MaxPayout: maxPayout, RemoveChance: 1},
	}

	// chance to get any wild in wild respins per band.
	freeWildRespinChances = [5][]float64{
		{14, 1},
		{21, 4},
		{22.5, 9, 1},
		{35, 17, 6, 1},
		{55, 38, 12, 3.5, 0.5},
	}

	// chance to get any wild in super bonus spins per band.
	freeWildSuperChances = [5][]float64{
		{1},
		{2},
		{3.1},
		{5, 0.2},
		{8, 0.5},
	}

	modBonus = comp.NewMultiFunc(
		comp.NewPowerFunc(2.4, 0, func(spin *comp.Spin) float64 {
			return float64(spin.SpinSeq()-1) / 3.3
		}),
		comp.NewDivideFunc(10, 0, func(spin *comp.Spin) float64 {
			return 10 + (2.3 * spin.TotalPayout())
		}),
	)

	band5f = comp.OnRoundFlagValue(flagFreeSpinBands, 5)
	bb2    = comp.OnRoundFlagValue(flagBonusBuy, bonusKind2)

	modSuperBonus = comp.NewMultiFunc(
		comp.NewPowerFunc(2.75, 0, func(spin *comp.Spin) float64 {
			if bb2(spin) {
				if band5f(spin) {
					return 1
				}

				if spin.TotalPayout() >= 90 {
					return -1
				}
				return float64(spin.SpinSeq()-1) / 1.61
			}

			if spin.TotalPayout() >= 90 {
				return -1
			}
			return float64(spin.SpinSeq()-1) / 1.5
		}),
		comp.NewDivideFunc(6, 0, func(spin *comp.Spin) float64 {
			if bb2(spin) && band5f(spin) {
				return 1
			}
			return math.Pow(float64(spin.FreeSpins()+1), 1.1)
		}),
		comp.NewDivideFunc(10, 0, func(spin *comp.Spin) float64 {
			if bb2(spin) && band5f(spin) {
				return 1
			}

			p := spin.TotalPayout()
			switch {
			case p > 90:
				return 300
			case p > 150:
				return 700
			case p > 500:
				return 5000
			default:
				return 10 + (2.4 * spin.TotalPayout())
			}
		}),
	)

	script1Weight = 1.4
	script2Weight = 2.2
	script3Weight = 2.2
	script4Weight = 2.2
	script5Weight = 1.0
	script6Weight = 1.0
	script9Weight = 29990.0
)

var (
	gridMask = []uint8{2, 3, 4, 4, 3, 2}

	n01 = comp.WithName("Diamond")
	n02 = comp.WithName("Club")
	n03 = comp.WithName("Heart")
	n04 = comp.WithName("Spade")
	n05 = comp.WithName("Knight")
	n06 = comp.WithName("Rogue")
	n07 = comp.WithName("Magician")
	n08 = comp.WithName("Barbarian")
	n09 = comp.WithName("Magic Orb")
	n10 = comp.WithName("Ice Princess")
	n11 = comp.WithName("Ice Princess x2")
	n12 = comp.WithName("Ice Princess x3")
	n13 = comp.WithName("Ice Princess x4")

	r01 = comp.WithResource("l4")
	r02 = comp.WithResource("l3")
	r03 = comp.WithResource("l2")
	r04 = comp.WithResource("l1")
	r05 = comp.WithResource("h4")
	r06 = comp.WithResource("h3")
	r07 = comp.WithResource("h2")
	r08 = comp.WithResource("h1")
	r09 = comp.WithResource("scatter")
	r10 = comp.WithResource("wild")
	r11 = comp.WithResource("wild2")
	r12 = comp.WithResource("wild3")
	r13 = comp.WithResource("wild4")

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

	wildRespinGrid    = comp.GridOffsets{{0, 0}, {0, 1}, {5, 0}, {5, 1}}
	fakeRespinGrid1   = comp.GridOffsets{{0, 0}, {0, 1}}
	fakeRespinGrid2   = comp.GridOffsets{{0, 0}, {0, 1}, {5, 0}}
	fakeRespinGrid3   = comp.GridOffsets{{0, 0}, {0, 1}, {5, 1}}
	fakeRespinGrid4   = comp.GridOffsets{{5, 0}, {5, 1}}
	fakeRespinGrid5   = comp.GridOffsets{{0, 0}, {5, 0}, {5, 1}}
	fakeRespinGrid6   = comp.GridOffsets{{0, 1}, {5, 0}, {5, 1}}
	wildRespinCenter  = comp.GridOffsets{{0, 0}}
	wildRespin        = util.AcquireWeighting().AddWeights(util.Indexes{wild}, []float64{1})
	wildRespinOffsets = util.UInt8s{0, 1, 20, 21}
	scatterReels      = util.UInt8s{2, 3, 4, 5}
	scatterReels2     = []int{2, 3, 4, 5}
	freeSpinBands     = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5}, freeSpinBandWeights)
	freeSpinBandsBB1  = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5}, freeSpinBandWeightsBB1)
	freeSpinBandsBB2  = util.AcquireWeighting().AddWeights(util.Indexes{1, 2, 3, 4, 5}, freeSpinBandWeightsBB2)
	freeSpinBand5     = util.AcquireWeighting().AddWeights(util.Indexes{1}, []float64{1})

	flag0 = comp.NewRoundFlag(flagFreeSpinBands, "free spin band")
	flag1 = comp.NewRoundFlag(flagFreeSpinCount, "free spin count")
	flag2 = comp.NewRoundFlag(flagWildRespin, "wild respin landed")
	flag3 = comp.NewRoundFlag(flagBonusBuy, "bonus buy feature")
	flags = comp.RoundFlags{flag0, flag1, flag2, flag3}

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
	fullPage4XWeights        = util.AcquireWeighting().AddWeights(util.Indexes{wild4}, []float64{100})
	fullPageDiamondsWeights  = util.AcquireWeighting().AddWeights(util.Indexes{1}, []float64{100})
	fullPageClubsWeights     = util.AcquireWeighting().AddWeights(util.Indexes{2}, []float64{100})
	fullPageHeartsWeights    = util.AcquireWeighting().AddWeights(util.Indexes{3}, []float64{100})
	fullPageSpadesWeights    = util.AcquireWeighting().AddWeights(util.Indexes{4}, []float64{100})
	fullPageMedusaWeights    = util.AcquireWeighting().AddWeights(util.Indexes{5}, []float64{100})
	fullPageDeathWeights     = util.AcquireWeighting().AddWeights(util.Indexes{6}, []float64{100})
	fullPageSuccubusWeights  = util.AcquireWeighting().AddWeights(util.Indexes{7}, []float64{100})
	fullPageDemonGoatWeights = util.AcquireWeighting().AddWeights(util.Indexes{8}, []float64{100})

	actions92All        comp.SpinActions
	actions92First      comp.SpinActions
	actions92Free       comp.SpinActions
	actions92FirstBonus comp.SpinActions
	actions92FreeBonus  comp.SpinActions
	actions94All        comp.SpinActions
	actions94First      comp.SpinActions
	actions94Free       comp.SpinActions
	actions94FirstBonus comp.SpinActions
	actions94FreeBonus  comp.SpinActions
	actions96All        comp.SpinActions
	actions96First      comp.SpinActions
	actions96Free       comp.SpinActions
	actions96FirstBonus comp.SpinActions
	actions96FreeBonus  comp.SpinActions
	actionsAllB         comp.SpinActions

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

	// bonus buy features.
	bonusBuy1 := comp.NewPaidAction(comp.FreeSpins, 0, bonusBuyCost1, scatter, 3, scatterReels...).WithFlag(flagBonusBuy, bonusKind1).WithBonusKind(bonusKind1)
	bonusBuy1.Describe(bonusBuy1ID, "princess spin bonus buy")

	bonusBuy2 := comp.NewPaidAction(comp.FreeSpins, 0, bonusBuyCost2, scatter, 3, scatterReels...).WithFlag(flagBonusBuy, bonusKind2).WithBonusKind(bonusKind2)
	bonusBuy2.Describe(bonusBuy2ID, "frost princess spin bonus buy")

	bonusExtraScatter1 := comp.NewGenerateSymbolAction(scatter, []float64{25.8}, scatterReels...).AllowPrevious().GenerateNoDupes()
	bonusExtraScatter1.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusBuy, bonusKind1))
	bonusExtraScatter1.Describe(bonusExtra1ID, "generate 4th scatter - bonus buy 1")

	bonusExtraScatter2 := comp.NewGenerateSymbolAction(scatter, []float64{11}, scatterReels...).AllowPrevious().GenerateNoDupes()
	bonusExtraScatter2.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusBuy, bonusKind2))
	bonusExtraScatter2.Describe(bonusExtra2ID, "generate 4th scatter - bonus buy 2")

	bonusBuy2wilds := comp.NewGenerateShapeAction(100, wildRespinGrid, wildRespinCenter, wildRespin)
	bonusBuy2wilds.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusBuy, bonusKind2))
	bonusBuy2wilds.Describe(bonusBuy2wildsID, "generate wilds - frost princess spin bonus buy")

	// generate (fake) wild respins RTP 92.
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

	actions92All = append(actions92All, wildRespin92g)

	// generate wilds RTP 92.
	firstWild92 := comp.NewGenerateSymbolAction(wild, firstWildChances92).GenerateNoDupes()
	firstWild92.WithTriggerFilters(comp.OnFirstSpin) // leave this in place as it must not be called during refills!!!
	firstWild92.Describe(firstWild92ID, "generate wilds - first spin - RTP 92")

	freeWild92a := comp.NewGenerateSymbolAction(wild, freeWildChances92a, scatterReels...)
	freeWild92a.WithChanceModifier(modBonus)
	freeWild92a.WithTriggerFilters(comp.OnRoundFlagValue(flagWildRespin, 0))
	freeWild92a.Describe(freeWild92aID, "generate wilds - free spins - no devil - RTP 92")

	freeWild92b := comp.NewGenerateSymbolAction(wild, freeWildChances92b, scatterReels...)
	freeWild92b.WithChanceModifier(modSuperBonus)
	freeWild92b.WithTriggerFilters(comp.OnRoundFlagValue(flagWildRespin, 1))
	freeWild92b.Describe(freeWild92bID, "generate wilds - free spins - with devil - RTP 92")

	freeWild92bb1 := comp.NewGenerateSymbolAction(wild, freeWildChances92bb1, scatterReels...)
	freeWild92bb1.WithChanceModifier(modBonus)
	freeWild92bb1.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusBuy, bonusKind1))
	freeWild92bb1.Describe(freeWild92bb1ID, "generate wilds - free spins - RTP 92 - bonus buy 1")

	freeWild92bb2 := comp.NewGenerateSymbolAction(wild, freeWildChances92bb2, scatterReels...)
	freeWild92bb2.WithChanceModifier(modSuperBonus)
	freeWild92bb2.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusBuy, bonusKind2))
	freeWild92bb2.Describe(freeWild92bb2ID, "generate wilds - free spins - RTP 92 - bonus buy 2")

	actions92All = append(actions92All, firstWild92, freeWild92a, freeWild92b, freeWild92bb1, freeWild92bb2)

	// generate scatters RTP 92.
	firstScatter92a := comp.NewGenerateSymbolAction(scatter, firstScatterChances92a, scatterReels...).GenerateNoDupes()
	firstScatter92a.WithTriggerFilters(comp.OnFirstSpin) // leave this in place as it must not be called during refills!!!
	firstScatter92a.Describe(firstScatter92aID, "generate scatters - first spin - no wild respin - RTP 92")

	firstScatter92b := comp.NewGenerateSymbolAction(scatter, firstScatterChances92b, scatterReels...).GenerateNoDupes()
	firstScatter92b.WithTriggerFilters(comp.OnFirstSpin, comp.OnGridShape(wild, wildRespinOffsets)) // leave this in place as it must not be called during refills!!!
	firstScatter92b.Describe(firstScatter92bID, "generate scatters - first spin - wild respin - RTP 92")

	freeScatter92a := buildFreeScatterAction(92, freeScatter92a1ID, freeScatterChances92a, "no devil")
	freeScatter92a.WithChanceModifier(modBonus)
	freeScatter92a.WithTriggerFilters(comp.OnRoundFlagValue(flagWildRespin, 0))
	freeScatter92a.Describe(freeScatter92aID, "multi-select generate scatters - no devil - RTP 92")

	freeScatter92b := buildFreeScatterAction(92, freeScatter92b1ID, freeScatterChances92b, "with devil")
	freeScatter92b.WithChanceModifier(modSuperBonus)
	freeScatter92b.WithTriggerFilters(comp.OnRoundFlagValue(flagWildRespin, 1))
	freeScatter92b.Describe(freeScatter92bID, "multi-select generate scatters - with devil - RTP 92")

	freeScatter92bb1 := buildFreeScatterAction(92, freeScatter92bb1aID, freeScatterChances92bb1, "bonus buy 1")
	freeScatter92bb1.WithChanceModifier(modBonus)
	freeScatter92bb1.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusBuy, bonusKind1))
	freeScatter92bb1.Describe(freeScatter92bb1ID, "multi-select generate scatters - RTP 92 - bonus buy 1")

	freeScatter92bb2 := buildFreeScatterAction(92, freeScatter92bb2aID, freeScatterChances92bb2, "bonus buy 2")
	freeScatter92bb2.WithChanceModifier(modSuperBonus)
	freeScatter92bb2.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusBuy, bonusKind2))
	freeScatter92bb2.Describe(freeScatter92bb2ID, "multi-select generate scatters - RTP 92 - bonus buy 2")

	actions92All = append(actions92All, firstScatter92a, firstScatter92b, freeScatter92a, freeScatter92b, freeScatter92bb1, freeScatter92bb2)

	// generate (fake) wild respins RTP 94.
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

	actions94All = append(actions94All, wildRespin94g)

	// generate wilds RTP 94.
	firstWild94 := comp.NewGenerateSymbolAction(wild, firstWildChances94).GenerateNoDupes()
	firstWild94.WithTriggerFilters(comp.OnFirstSpin) // leave this in place as it must not be called during refills!!!
	firstWild94.Describe(firstWild94ID, "generate wilds - first spin - RTP 94")

	freeWild94a := comp.NewGenerateSymbolAction(wild, freeWildChances94a, scatterReels...)
	freeWild94a.WithChanceModifier(modBonus)
	freeWild94a.WithTriggerFilters(comp.OnRoundFlagValue(flagWildRespin, 0))
	freeWild94a.Describe(freeWild94aID, "generate wilds - free spins - no devil - RTP 94")

	freeWild94b := comp.NewGenerateSymbolAction(wild, freeWildChances94b, scatterReels...)
	freeWild94b.WithChanceModifier(modSuperBonus)
	freeWild94b.WithTriggerFilters(comp.OnRoundFlagValue(flagWildRespin, 1))
	freeWild94b.Describe(freeWild94bID, "generate wilds - free spins - with devil - RTP 94")

	freeWild94bb1 := comp.NewGenerateSymbolAction(wild, freeWildChances94bb1, scatterReels...)
	freeWild94bb1.WithChanceModifier(modBonus)
	freeWild94bb1.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusBuy, bonusKind1))
	freeWild94bb1.Describe(freeWild94bb1ID, "generate wilds - free spins - RTP 94 - bonus buy 1")

	freeWild94bb2 := comp.NewGenerateSymbolAction(wild, freeWildChances94bb2, scatterReels...)
	freeWild94bb2.WithChanceModifier(modSuperBonus)
	freeWild94bb2.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusBuy, bonusKind2))
	freeWild94bb2.Describe(freeWild94bb2ID, "generate wilds - free spins - RTP 94 - bonus buy 2")

	actions94All = append(actions94All, firstWild94, freeWild94a, freeWild94b, freeWild94bb1, freeWild94bb2)

	// generate scatters RTP 94.
	firstScatter94a := comp.NewGenerateSymbolAction(scatter, firstScatterChances94a, scatterReels...).GenerateNoDupes()
	firstScatter94a.WithTriggerFilters(comp.OnFirstSpin) // leave this in place as it must not be called during refills!!!
	firstScatter94a.Describe(firstScatter94aID, "generate scatters - first spin - no wild respin - RTP 94")

	firstScatter94b := comp.NewGenerateSymbolAction(scatter, firstScatterChances94b, scatterReels...).GenerateNoDupes()
	firstScatter94b.WithTriggerFilters(comp.OnFirstSpin, comp.OnGridShape(wild, wildRespinOffsets)) // leave this in place as it must not be called during refills!!!
	firstScatter94b.Describe(firstScatter94bID, "generate scatters - first spin - wild respin - RTP 94")

	freeScatter94a := buildFreeScatterAction(94, freeScatter94a1ID, freeScatterChances94a, "no devil")
	freeScatter94a.WithChanceModifier(modBonus)
	freeScatter94a.WithTriggerFilters(comp.OnRoundFlagValue(flagWildRespin, 0))
	freeScatter94a.Describe(freeScatter94aID, "multi-select generate scatters - no devil - RTP 94")

	freeScatter94b := buildFreeScatterAction(94, freeScatter94b1ID, freeScatterChances94b, "with devil")
	freeScatter94b.WithChanceModifier(modSuperBonus)
	freeScatter94b.WithTriggerFilters(comp.OnRoundFlagValue(flagWildRespin, 1))
	freeScatter94b.Describe(freeScatter94bID, "multi-select generate scatters - with devil - RTP 94")

	freeScatter94bb1 := buildFreeScatterAction(94, freeScatter94bb1aID, freeScatterChances94bb1, "bonus buy 1")
	freeScatter94bb1.WithChanceModifier(modBonus)
	freeScatter94bb1.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusBuy, bonusKind1))
	freeScatter94bb1.Describe(freeScatter94bb1ID, "multi-select generate scatters - RTP 94 - bonus buy 1")

	freeScatter94bb2 := buildFreeScatterAction(94, freeScatter94bb2aID, freeScatterChances94bb2, "bonus buy 2")
	freeScatter94bb2.WithChanceModifier(modSuperBonus)
	freeScatter94bb2.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusBuy, bonusKind2))
	freeScatter94bb2.Describe(freeScatter94bb2ID, "multi-select generate scatters - RTP 94 - bonus buy 2")

	actions94All = append(actions94All, firstScatter94a, firstScatter94b, freeScatter94a, freeScatter94b, freeScatter94bb1, freeScatter94bb2)

	// generate (fake) wild respins RTP 96.
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

	actions96All = append(actions96All, wildRespin96g)

	// generate wilds RTP 96.
	firstWild96 := comp.NewGenerateSymbolAction(wild, firstWildChances96).GenerateNoDupes()
	firstWild96.WithTriggerFilters(comp.OnFirstSpin) // leave this in place as it must not be called during refills!!!
	firstWild96.Describe(firstWild96ID, "generate wilds - first spin - RTP 96")

	freeWild96a := comp.NewGenerateSymbolAction(wild, freeWildChances96a, scatterReels...)
	freeWild96a.WithChanceModifier(modBonus)
	freeWild96a.WithTriggerFilters(comp.OnRoundFlagValue(flagWildRespin, 0))
	freeWild96a.Describe(freeWild96aID, "generate wilds - free spins - no devil - RTP 96")

	freeWild96b := comp.NewGenerateSymbolAction(wild, freeWildChances96b, scatterReels...)
	freeWild96b.WithChanceModifier(modSuperBonus)
	freeWild96b.WithTriggerFilters(comp.OnRoundFlagValue(flagWildRespin, 1))
	freeWild96b.Describe(freeWild96bID, "generate wilds - free spins - with devil - RTP 96")

	freeWild96bb1 := comp.NewGenerateSymbolAction(wild, freeWildChances96bb1, scatterReels...)
	freeWild96bb1.WithChanceModifier(modBonus)
	freeWild96bb1.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusBuy, bonusKind1))
	freeWild96bb1.Describe(freeWild96bb1ID, "generate wilds - free spins - RTP 96 - bonus buy 1")

	freeWild96bb2 := comp.NewGenerateSymbolAction(wild, freeWildChances96bb2, scatterReels...)
	freeWild96bb2.WithChanceModifier(modSuperBonus)
	freeWild96bb2.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusBuy, bonusKind2))
	freeWild96bb2.Describe(freeWild96bb2ID, "generate wilds - free spins - RTP 96 - bonus buy 2")

	actions96All = append(actions96All, firstWild96, freeWild96a, freeWild96b, freeWild96bb1, freeWild96bb2)

	// generate scatters RTP 96.
	firstScatter96a := comp.NewGenerateSymbolAction(scatter, firstScatterChances96a, scatterReels...).GenerateNoDupes()
	firstScatter96a.WithTriggerFilters(comp.OnFirstSpin) // leave this in place as it must not be called during refills!!!
	firstScatter96a.Describe(firstScatter96aID, "generate scatters - first spin - no wild respin - RTP 96")

	firstScatter96b := comp.NewGenerateSymbolAction(scatter, firstScatterChances96b, scatterReels...).GenerateNoDupes()
	firstScatter96b.WithTriggerFilters(comp.OnFirstSpin, comp.OnGridShape(wild, wildRespinOffsets)) // leave this in place as it must not be called during refills!!!
	firstScatter96b.Describe(firstScatter96bID, "generate scatters - first spin - wild respin - RTP 96")

	freeScatter96a := buildFreeScatterAction(96, freeScatter96a1ID, freeScatterChances96a, "no devil")
	freeScatter96a.WithChanceModifier(modBonus)
	freeScatter96a.WithTriggerFilters(comp.OnRoundFlagValue(flagWildRespin, 0))
	freeScatter96a.Describe(freeScatter96aID, "multi-select generate scatters - no devil - RTP 96")

	freeScatter96b := buildFreeScatterAction(96, freeScatter96b1ID, freeScatterChances96b, "with devil")
	freeScatter96b.WithChanceModifier(modSuperBonus)
	freeScatter96b.WithTriggerFilters(comp.OnRoundFlagValue(flagWildRespin, 1))
	freeScatter96b.Describe(freeScatter96bID, "multi-select generate scatters - with devil - RTP 96")

	freeScatter96bb1 := buildFreeScatterAction(96, freeScatter96bb1aID, freeScatterChances96bb1, "bonus buy 1")
	freeScatter96bb1.WithChanceModifier(modBonus)
	freeScatter96bb1.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusBuy, bonusKind1))
	freeScatter96bb1.Describe(freeScatter96bb1ID, "multi-select generate scatters - RTP 96 - bonus buy 1")

	freeScatter96bb2 := buildFreeScatterAction(96, freeScatter96bb2aID, freeScatterChances96bb2, "bonus buy 2")
	freeScatter96bb2.WithChanceModifier(modSuperBonus)
	freeScatter96bb2.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusBuy, bonusKind2))
	freeScatter96bb2.Describe(freeScatter96bb2ID, "multi-select generate scatters - RTP 96 - bonus buy 2")

	actions96All = append(actions96All, firstScatter96a, firstScatter96b, freeScatter96a, freeScatter96b, freeScatter96bb1, freeScatter96bb2)

	// generate 2x/3x/4x wilds on free wild respins.
	freeWildRespin := buildFreeWildAction(freeWildRespin1ID, freeWildRespinChances, nil, "wild respin")
	freeWildRespin.WithTriggerFilters(comp.OnRefillSpin, comp.OnRoundFlagValue(flagWildRespin, 1))
	freeWildRespin.Describe(freeWildRespin2ID, "multi-select generate any wild - wild respin")

	// generate 2x/3x/4x wilds on free super bonus spins.
	freeWildSuper := buildFreeWildAction(freeWildSuper1ID, freeWildSuperChances, modSuperBonus, "super bonus")
	freeWildSuper.WithTriggerFilters(comp.OnFreeSpin, comp.OnRoundFlagValue(flagWildRespin, 1))
	freeWildSuper.Describe(freeWildSuper2ID, "multi-select generate any wild - super bonus")

	// all paylines.
	paylines := comp.NewAllPaylinesAction(true)
	paylines.Describe(paylinesID, "all paylines")

	// reduce payouts by band - first spin.
	reducePayoutsFirst := comp.NewRemovePayoutBandsAction(3, direction, true, true, reduceBandsFirst)
	reducePayoutsFirst.WithTriggerFilters(comp.OnFirstSpin) // leave this in place as it must not be called during refills!!!
	reducePayoutsFirst.Describe(removePayoutsFirstID, "remove payouts by band (first spin)")

	// reduce payouts by band - free wild respin.
	reducePayoutsRespin := comp.NewRemovePayoutBandsAction(3, direction, true, true, reduceBandsRespin)
	reducePayoutsRespin.WithTriggerFilters(comp.OnRefillSpin)
	reducePayoutsRespin.Describe(removePayoutsRespinID, "remove payouts by band (wild respin)")

	// reduce payouts by band - first few free bonus spins.
	reducePayoutsFree1 := comp.NewRemovePayoutBandsAction(3, direction, false, true, reduceBandsFree1)
	reducePayoutsFree1.WithTriggerFilters(comp.OnRoundFlagBelow(flagFreeSpinCount, 4))
	reducePayoutsFree1.Describe(removePayoutsFree1ID, "remove payouts by band (first few free spins)")

	// reduce payouts by band - remaining free bonus spins.
	reducePayoutsFree2 := comp.NewRemovePayoutBandsAction(3, direction, false, true, reduceBandsFree2)
	reducePayoutsFree2.WithTriggerFilters(comp.OnRoundFlagAbove(flagFreeSpinCount, 3))
	reducePayoutsFree2.Describe(removePayoutsFree2ID, "remove payouts by band (remaining free spins)")

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

	payoutBandsBB1 := comp.NewRoundFlagWeightedAction(flagFreeSpinBands, freeSpinBandsBB1)
	payoutBandsBB1.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusBuy, bonusKind1))
	payoutBandsBB1.Describe(payoutBandsBB1ID, "weighted payout band (flag 0)")

	payoutBandsBB2 := comp.NewRoundFlagWeightedAction(flagFreeSpinBands, freeSpinBandsBB2)
	payoutBandsBB2.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusBuy, bonusKind2))
	payoutBandsBB2.Describe(payoutBandsBB2ID, "weighted payout band (flag 0)")

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

	// force full page of diamonds
	forceDiamondsFP := comp.NewGenerateShapeAction(100, fullPageR1to4shape, fullPageR1offset, fullPageDiamondsWeights).AllowPrevious()
	forceDiamondsFP.WithTriggerFilters(comp.OnSpinSequence(2))
	forceDiamondsFP.Describe(forceDiamondsFPID, "force full page of diamonds (spin 2, scripted round)")

	// force full page of clubs.
	forceClubsFP := comp.NewGenerateShapeAction(100, fullPageR1to4shape, fullPageR1offset, fullPageClubsWeights).AllowPrevious()
	forceClubsFP.WithTriggerFilters(comp.OnSpinSequence(3))
	forceClubsFP.Describe(forceClubsFPID, "force full page of clubs (spin 3, scripted round)")

	// force full page of sword.
	forceHeartsFP := comp.NewGenerateShapeAction(100, fullPageR1to4shape, fullPageR1offset, fullPageHeartsWeights).AllowPrevious()
	forceHeartsFP.WithTriggerFilters(comp.OnSpinSequence(4))
	forceHeartsFP.Describe(forceHeartsFPID, "force full page of hearts (spin 4, scripted round)")

	// force full page of spades.
	forceSpadesFP := comp.NewGenerateShapeAction(100, fullPageR1to4shape, fullPageR1offset, fullPageSpadesWeights).AllowPrevious()
	forceSpadesFP.WithTriggerFilters(comp.OnSpinSequence(5))
	forceSpadesFP.Describe(forceSpadesFPID, "force full page of spades (spin 5, scripted round)")

	// force full page of medusa.
	forceKnightFP := comp.NewGenerateShapeAction(100, fullPageR1to4shape, fullPageR1offset, fullPageMedusaWeights).AllowPrevious()
	forceKnightFP.WithTriggerFilters(comp.OnSpinSequence(7))
	forceKnightFP.Describe(forceMedusaFPID, "force full page of knights (spin 7, scripted round)")

	// force full page of death.
	forceRogueFP := comp.NewGenerateShapeAction(100, fullPageR1to4shape, fullPageR1offset, fullPageDeathWeights).AllowPrevious()
	forceRogueFP.WithTriggerFilters(comp.OnSpinSequence(8))
	forceRogueFP.Describe(forceDeathFPID, "force full page of rogues (spin 8, scripted round)")

	// force full page of succubus.
	forceMagicianFP := comp.NewGenerateShapeAction(100, fullPageR1to4shape, fullPageR1offset, fullPageSuccubusWeights).AllowPrevious()
	forceMagicianFP.WithTriggerFilters(comp.OnSpinSequence(9))
	forceMagicianFP.Describe(forceSuccubusFPID, "force full page of magicians (spin 9, scripted round)")

	// force full page of demon goat.
	forceBarbarianFP := comp.NewGenerateShapeAction(100, fullPageR1to4shape, fullPageR1offset, fullPageDemonGoatWeights).AllowPrevious()
	forceBarbarianFP.WithTriggerFilters(comp.OnSpinSequence(11))
	forceBarbarianFP.Describe(forceDemonGoatFPID, "force full page of barbarians (spin 11, scripted round)")

	// force many 1x wilds on first spin for scripted rounds.
	forceW1S1 := comp.NewGenerateSymbolAction(wild, []float64{100, 100, 100, 100, 100, 90, 70, 5}, 2, 3, 4, 5).AllowPrevious()
	forceW1S1.Describe(forceW1S1ID, "force many x1 wilds (spin 1; scripted round)")

	// force 3-4 random multiplier wilds for scripted rounds.
	weights4RW := util.AcquireWeighting().AddWeights(util.Indexes{wild, wild2, wild3, wild4}, []float64{12, 6, 2, 8})
	force4RW := comp.NewGenerateSymbolsAction(weights4RW, []float64{100, 100, 95, 95, 70, 60, 30, 3}, scatterReels...)
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

	forceW4S6UpRandom := comp.NewGenerateSymbolAction(wild4, []float64{100, 100, 95, 70, 60, 15, 1}, scatterReels...)
	forceW4S6UpRandom.WithTriggerFilters(comp.OnSpinSequenceAbove(5))
	forceW4S6UpRandom.Describe(forceW4S6UpRandomID, "force x4 wilds (spin 6+ random; scripted round)")

	forceW4S9UpRandom := comp.NewGenerateSymbolAction(wild4, []float64{100, 100, 95, 70, 60, 15, 1}, scatterReels...)
	forceW4S9UpRandom.WithTriggerFilters(comp.OnSpinSequenceAbove(8))
	forceW4S9UpRandom.Describe(forceW4S9UpRandomID, "force x4 wilds (spin 9+ random; scripted round)")

	// force full page of 4x wilds for scripted rounds.
	forceFullW4S6UpRandom := comp.NewGenerateShapeAction(30, fullPageR1to4shape, fullPageR1offset, fullPage4XWeights)
	forceFullW4S6UpRandom.WithTriggerFilters(comp.OnSpinSequenceAbove(5))
	forceFullW4S6UpRandom.Describe(forceFullW4S6UpRandomID, "force full page x4 wild (spin 6+; scripted round)")

	forceFullW4S9UpRandom := comp.NewGenerateShapeAction(30, fullPageR1to4shape, fullPageR1offset, fullPage4XWeights)
	forceFullW4S9UpRandom.WithTriggerFilters(comp.OnSpinSequenceAbove(8))
	forceFullW4S9UpRandom.Describe(forceFullW4S9UpRandomID, "force full page x4 wild (spin 9+; scripted round)")

	forceFullW4S12UpRandom := comp.NewGenerateShapeAction(70, fullPageR1to4shape, fullPageR1offset, fullPage4XWeights)
	forceFullW4S12UpRandom.WithTriggerFilters(comp.OnSpinSequenceAbove(11))
	forceFullW4S12UpRandom.Describe(forceFullW4S12UpRandomID, "force full page x4 wild (spin 12+; scripted round)")

	// scripted rounds - scenario 1 configuration (max payout).
	script1 := comp.NewScriptedRound(1, script1Weight,
		comp.SpinActions{forceWildRespin, forceScrolls4, paylines, forceNoPayout, award8, wildRespinFlag, refill, stickies1, forceBand5},
		comp.SpinActions{forceW4S2R5, forceW4S3R4, forceW4S4R5, forceW4S5R4, forceScrolls4S6,
			forceW4S6UpRandom, forceScrollsS10Up,
			paylines, retrigger4, stickies2},
	).WithBonusBuys(bonusKind2)

	// scripted rounds - scenario 2 configuration (max payout).
	script2 := comp.NewScriptedRound(2, script2Weight,
		comp.SpinActions{forceWildRespin, forceScrolls4, paylines, award8, wildRespinFlag, refill, stickies1, forceBand5},
		comp.SpinActions{force4RW, paylines, retrigger4, stickies2},
	).WithBonusBuys(bonusKind2)

	// scripted rounds - scenario 3 configuration (max payout).
	script3 := comp.NewScriptedRound(3, script3Weight,
		comp.SpinActions{forceWildRespin, forceScrolls4, paylines, award8, wildRespinFlag, refill, stickies1, forceBand5},
		comp.SpinActions{forceScrollsS2to8, force4RWS8Up, forceScrollsS9Up, paylines, retrigger4, stickies2},
	).WithBonusBuys(bonusKind2)

	// scripted rounds - scenario 4 configuration (max payout).
	script4 := comp.NewScriptedRound(4, script4Weight,
		comp.SpinActions{forceWildRespin, forceW1S1, forceScrolls4, paylines, award8, wildRespinFlag, refill, stickies1, forceBand5},
		comp.SpinActions{forceScrollsS2to8, force4RWS8Up, forceScrollsS9Up, paylines, retrigger4, stickies2},
	).WithBonusBuys(bonusKind2)

	// scripted rounds - scenario 5 configuration (max payout).
	script5 := comp.NewScriptedRound(5, script5Weight,
		comp.SpinActions{forceWildRespin, forceScrolls4, paylines, forceNoPayout, award8, wildRespinFlag, refill, stickies1, forceBand5},
		comp.SpinActions{forceScrollsS2to8, forceScrollsS9Up, forceFullW4S9UpRandom, paylines, retrigger4, stickies2},
	).WithBonusBuys(bonusKind2)

	// scripted rounds - scenario 6 configuration (max payout).
	script6 := comp.NewScriptedRound(6, script6Weight,
		comp.SpinActions{forceWildRespin, forceScrolls4, paylines, forceNoPayout, award8, wildRespinFlag, refill, stickies1, forceBand5},
		comp.SpinActions{forceDiamondsFP, forceClubsFP, forceHeartsFP, forceSpadesFP, forceFullW4S6UpRandom, paylines, retrigger4, stickies2},
	).WithBonusBuys(bonusKind2)

	// scripted rounds - scenario 9 configuration (dummy to fill up the weighting).
	script9 := comp.NewScriptedRound(9, script9Weight, nil, nil)

	// scripted rounds accumulated.
	scriptedRounds = comp.NewScriptedRoundSelector(0.1, script1, script2, script3, script4, script5, script6, script9)

	// gather it all together.
	actions92All = append(comp.SpinActions{bonusBuy1, bonusBuy2, bonusExtraScatter1, bonusExtraScatter2, bonusBuy2wilds, freeCounter}, actions92All...)
	actions92FirstA := comp.SpinActions{wildRespin92g, firstWild92, freeWildRespin, firstScatter92a, firstScatter92b}
	actions92FreeA := comp.SpinActions{freeCounter, freeWild92a, freeWild92b, freeScatter92a, freeScatter92b}
	actions92FirstBonusA := comp.SpinActions{bonusBuy1, bonusBuy2, bonusExtraScatter1, bonusExtraScatter2, bonusBuy2wilds}
	actions92FreeBonusA := comp.SpinActions{freeCounter, freeWild92bb1, freeWild92bb2, freeScatter92bb1, freeScatter92bb2}

	actions94All = append(comp.SpinActions{bonusBuy1, bonusBuy2, bonusExtraScatter1, bonusExtraScatter2, bonusBuy2wilds, freeCounter}, actions94All...)
	actions94FirstA := comp.SpinActions{wildRespin94g, firstWild94, freeWildRespin, firstScatter94a, firstScatter94b}
	actions94FreeA := comp.SpinActions{freeCounter, freeWild94a, freeWild94b, freeScatter94a, freeScatter94b}
	actions94FirstBonusA := comp.SpinActions{bonusBuy1, bonusBuy2, bonusExtraScatter1, bonusExtraScatter2, bonusBuy2wilds}
	actions94FreeBonusA := comp.SpinActions{freeCounter, freeWild94bb1, freeWild94bb2, freeScatter94bb1, freeScatter94bb2}

	actions96All = append(comp.SpinActions{bonusBuy1, bonusBuy2, bonusExtraScatter1, bonusExtraScatter2, bonusBuy2wilds, freeCounter}, actions96All...)
	actions96FirstA := comp.SpinActions{wildRespin96g, firstWild96, freeWildRespin, firstScatter96a, firstScatter96b}
	actions96FreeA := comp.SpinActions{freeCounter, freeWild96a, freeWild96b, freeScatter96a, freeScatter96b}
	actions96FirstBonusA := comp.SpinActions{bonusBuy1, bonusBuy2, bonusExtraScatter1, bonusExtraScatter2, bonusBuy2wilds}
	actions96FreeBonusA := comp.SpinActions{freeCounter, freeWild96bb1, freeWild96bb2, freeScatter96bb1, freeScatter96bb2}

	actionsAllB = append(actionsAllB,
		paylines, reducePayoutsFirst, reducePayoutsRespin, reducePayoutsFree1, reducePayoutsFree2,
		teaser1, teaser2, teaser3, teaser4, teaser5, teaserFullPage, teaser3or4, teaserScattersWilds, teaserZeroPayWilds,
		award8, retrigger4, wildRespinFlag, refill, stickies1, stickies2, payoutBands, payoutBandsBB1, payoutBandsBB2,
	)
	actionsFirstB := comp.SpinActions{
		paylines, reducePayoutsFirst, reducePayoutsRespin, teaser3or4, teaserScattersWilds, teaserZeroPayWilds,
		award8, wildRespinFlag, refill, stickies1, payoutBands,
	}
	actionsFreeB := comp.SpinActions{
		freeWildSuper, paylines, reducePayoutsFree1, reducePayoutsFree2, retrigger4, stickies2,
	}
	actionsFirstBonusB := comp.SpinActions{
		paylines, award8, wildRespinFlag, stickies1, payoutBandsBB1, payoutBandsBB2,
	}
	actionsFreeBonusB := comp.SpinActions{
		freeWildSuper, paylines, reducePayoutsFree1, reducePayoutsFree2, retrigger4, stickies2,
	}

	forceAll := comp.SpinActions{
		forceWildRespin, forceBand5, forceNoPayout,
		forceScrolls4, forceScrolls4S6, forceScrollsS10Up, forceScrollsS2to8, forceScrollsS9Up,
		forceDiamondsFP, forceClubsFP, forceHeartsFP, forceSpadesFP, forceKnightFP, forceRogueFP, forceMagicianFP, forceBarbarianFP,
		forceW1S1, force4RW, force4RWS8Up,
		forceW4S2R5, forceW4S3R4, forceW4S4R5, forceW4S5R4, forceW4S7R2, forceW4S8R3,
		forceW4S6UpRandom, forceW4S9UpRandom,
		forceFullW4S6UpRandom, forceFullW4S9UpRandom, forceFullW4S12UpRandom}

	actions92All = append(append(actions92All, actionsAllB...), forceAll...)
	actions92First = append(actions92FirstA, actionsFirstB...)
	actions92Free = append(actions92FreeA, actionsFreeB...)
	actions92FirstBonus = append(actions92FirstBonusA, actionsFirstBonusB...)
	actions92FreeBonus = append(actions92FreeBonusA, actionsFreeBonusB...)

	actions94All = append(append(actions94All, actionsAllB...), forceAll...)
	actions94First = append(actions94FirstA, actionsFirstB...)
	actions94Free = append(actions94FreeA, actionsFreeB...)
	actions94FirstBonus = append(actions94FirstBonusA, actionsFirstBonusB...)
	actions94FreeBonus = append(actions94FreeBonusA, actionsFreeBonusB...)

	actions96All = append(append(actions96All, actionsAllB...), forceAll...)
	actions96First = append(actions96FirstA, actionsFirstB...)
	actions96Free = append(actions96FreeA, actionsFreeB...)
	actions96FirstBonus = append(actions96FirstBonusA, actionsFirstBonusB...)
	actions96FreeBonus = append(actions96FreeBonusA, actionsFreeBonusB...)
}

func buildFreeScatterAction(rtp, id int, chances [5][]float64, suffix string) *comp.MultiActionFlagValue {
	if suffix != "" {
		suffix = " - " + suffix
	}

	var list [5]*comp.ReviseAction
	for i := range list {
		a := comp.NewGenerateSymbolAction(scatter, chances[i], scatterReels...).GenerateNoDupes()
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

func buildFreeWildAction(id int, chances [5][]float64, mod comp.ChanceModifier, suffix string) *comp.MultiActionFlagValue {
	if suffix != "" {
		suffix = " - " + suffix
	}

	weights := util.AcquireWeighting().AddWeights(util.Indexes{wild, wild2, wild3, wild4}, []float64{15, 8, 3, 0.5})

	var list [5]*comp.ReviseAction
	for i := range list {
		a := comp.NewGenerateSymbolsAction(weights, chances[i], scatterReels...)
		if mod != nil {
			a.WithChanceModifier(mod)
		}
		a.Describe(id+i, fmt.Sprintf("generate any wild - free spin band %d%s", i+1, suffix))
		list[i] = a
	}

	a := comp.NewMultiActionFlagValue(flagFreeSpinBands, 1, list[0], 2, list[1], 3, list[2], 4, list[3], 5, list[4])
	actionsAllB = append(actionsAllB, list[0], list[1], list[2], list[3], list[4], a)
	return a
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

	var m1, m2 int
	if clearReel2 {
		m1, m2 = 4, 7
	} else {
		m1, m2 = 8, 12
	}
	for offset := m1; offset < m2; offset++ {
		indexes[offset] = ss[0]
	}

	if clearReel2 {
		m1, m2 = 8, 12
	} else {
		m1, m2 = 4, 7
	}
	for offset := m1; offset < m2; offset++ {
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
		comp.WithMask(gridMask...),
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

	slots92 = initSlots(92.0, weights92, actions92First, actions92Free, actions92FirstBonus, actions92FreeBonus)
	slots94 = initSlots(94.0, weights94, actions94First, actions94Free, actions94FirstBonus, actions94FreeBonus)
	slots96 = initSlots(96.0, weights96, actions96First, actions96Free, actions96FirstBonus, actions96FreeBonus)

	slots92params = game.RegularParams{Slots: slots92}
	slots94params = game.RegularParams{Slots: slots94}
	slots96params = game.RegularParams{Slots: slots96}
}

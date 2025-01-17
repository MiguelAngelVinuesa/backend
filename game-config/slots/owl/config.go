package owl

import (
	"math"

	comp "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	game "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/slots"
	util "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

const (
	symbolCount            = 11
	reels                  = 5
	rows                   = 3
	direction              = comp.PayLTR
	noRepeat               = 2
	maxPayout              = 2000.0
	bonusBuy               = 2.0
	freeSpinsGameType4     = 8
	freeSpinsGameType5     = 8
	freeSpinsGameType6     = 8
	freeSpinsGameType6free = 2

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
	id16 = 16
	id17 = 17
	id18 = 18
	id19 = 19
	id20 = 20

	wild    = id10
	scatter = id11

	bonusIncreaseID    = 1
	instantBonusID     = 10
	bonusTeaserID      = 11
	instantBonusHighID = 15
	bonusTeaserHighID  = 16
	bonusChoiceID      = 20
	bonusSelectorID    = 22
	removePayouts3ID   = 25
	scatters3ID        = 26
	scatters92aID      = 30
	scatters92bID      = 31
	scatters92cID      = 32
	scatters92dID      = 33
	wilds92aID         = 35
	wilds92bID         = 36
	wilds92cID         = 37
	wilds92dID         = 38
	scatters94aID      = 40
	scatters94bID      = 41
	scatters94cID      = 42
	scatters94dID      = 43
	wilds94aID         = 45
	wilds94bID         = 46
	wilds94cID         = 47
	wilds94dID         = 48
	scatters96aID      = 50
	scatters96bID      = 51
	scatters96cID      = 52
	scatters96dID      = 53
	wilds96aID         = 55
	wilds96bID         = 56
	wilds96cID         = 57
	wilds96dID         = 58
	reelNudgeID        = 70
	wildExpandID       = 72
	multiplierID       = 76
	winlinesID         = 80
	removePayouts1ID   = 81
	bonusWheelID       = 85
	freeSpins1ID       = 90
	freeSpins2ID       = 91
	freeSpins3ID       = 92
	freeSpins4aID      = 100
	freeSpins4bID      = 101
	freeSpins4cID      = 102
	freeSpins4dID      = 103
	freeSpins4eID      = 104
	freeSpins5aID      = 110
	freeSpins5bID      = 111
	freeSpins5cID      = 112
	freeSpins5dID      = 113
	freeSpins5eID      = 114
	freeSpinsFlagID    = 120

	flagBonusIncrease = 0 // 50% higher chance for instant bonus.
	flagBonusSelector = 1 // player choice from instant bonus.
	flagBonusGameType = 2 // bonus game type.
	flagFreeSpins     = 3 // number of free spins played.

	bonusChoice1 = 1 // instant bonus left.
	bonusChoice2 = 2 // instant bonus middle.
	bonusChoice3 = 3 // instant bonus right.

	bonusGameType1 = 1 // instant: 1 spin; 4-8 random wilds.
	bonusGameType2 = 2 // instant: 1 spin; 2-5 expanded wild reels.
	bonusGameType3 = 3 // instant: 1 spin; force scatters for bonus wheel.
	bonusGameType4 = 4 // wheel: 8 free spins with 4-8 random wilds.
	bonusGameType5 = 5 // wheel: 8 free spins with 2-5 expanded wild reels.
	bonusGameType6 = 6 // wheel: scatter progressive; 8 free spins; 2 per retrigger; 2x/5x/10x multiplier every 3 scatters.
)

var (
	// WIP WIP WIP
	weights92 = [symbolCount]comp.SymbolOption{
		comp.WithWeights(40, 50, 20, 50, 40),
		comp.WithWeights(40, 20, 50, 20, 40),
		comp.WithWeights(40, 30, 20, 30, 40),
		comp.WithWeights(30, 40, 25, 40, 30),
		comp.WithWeights(30, 40, 25, 40, 30),
		comp.WithWeights(12, 20, 12, 20, 12),
		comp.WithWeights(10, 20, 12, 20, 10),
		comp.WithWeights(15, 12, 6, 12, 15),
		comp.WithWeights(10, 8, 12, 8, 10),
		comp.WithWeights(0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0),
	}
	scatterChances92a = []float64{25, 10, 20}
	scatterChances92b = []float64{50, 10, 20} // double the base chance.
	scatterChances92c = []float64{15, 20, 8}
	scatterChances92d = []float64{15, 20, 8}
	wildChances92a    = []float64{15, 10, 5, 1}
	wildChances92b    = []float64{100, 100, 100, 100, 12.5, 5, 1, 0.5}
	wildChances92c    = []float64{100, 100, 5, 1, 0.5}
	wildChances92d    = []float64{20, 15, 10, 5, 1}

	// WIP WIP WIP
	weights94 = [symbolCount]comp.SymbolOption{
		comp.WithWeights(40, 50, 20, 50, 40),
		comp.WithWeights(40, 20, 50, 20, 40),
		comp.WithWeights(40, 30, 20, 30, 40),
		comp.WithWeights(30, 40, 25, 40, 30),
		comp.WithWeights(30, 40, 25, 40, 30),
		comp.WithWeights(12, 20, 12, 20, 12),
		comp.WithWeights(10, 20, 12, 20, 10),
		comp.WithWeights(15, 12, 6, 12, 15),
		comp.WithWeights(10, 8, 12, 8, 10),
		comp.WithWeights(0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0),
	}
	scatterChances94a = []float64{25, 10, 20}
	scatterChances94b = []float64{50, 10, 20} // double the base chance.
	scatterChances94c = []float64{15, 20, 8}
	scatterChances94d = []float64{15, 20, 8}
	wildChances94a    = []float64{15, 10, 5, 1}
	wildChances94b    = []float64{100, 100, 100, 100, 12.5, 5, 1, 0.5}
	wildChances94c    = []float64{100, 100, 5, 1, 0.5}
	wildChances94d    = []float64{20, 15, 10, 5, 1}

	// WIP WIP WIP
	weights96 = [symbolCount]comp.SymbolOption{
		comp.WithWeights(50, 70, 30, 70, 50),
		comp.WithWeights(40, 30, 70, 30, 40),
		comp.WithWeights(40, 50, 20, 50, 40),
		comp.WithWeights(30, 40, 50, 40, 30),
		comp.WithWeights(40, 25, 40, 25, 40),
		comp.WithWeights(12, 20, 10, 20, 12),
		comp.WithWeights(10, 18, 12, 18, 10),
		comp.WithWeights(15, 12, 6, 12, 15),
		comp.WithWeights(10, 8, 12, 8, 10),
		comp.WithWeights(0, 0, 0, 0, 0),
		comp.WithWeights(0, 0, 0, 0, 0),
	}
	scatterChances96a = []float64{25, 50, 50}
	scatterChances96b = []float64{50, 50, 50} // double the base chance.
	scatterChances96c = []float64{20, 20, 15}
	scatterChances96d = []float64{20, 20, 15}
	wildChances96a    = []float64{15, 10, 5, 1}
	wildChances96b    = []float64{100, 100, 100, 100, 12.5, 5, 1, 0.5}
	wildChances96c    = []float64{100, 100, 5, 1, 0.5}
	wildChances96d    = []float64{20, 15, 10, 5, 1}
)

var (
	bonusTeaserChance     = 0.2
	bonusChance           = 10.0
	bonusTeaserHighChance = 0.05
	n1                    = 100 - bonusTeaserChance
	n2                    = 100 - bonusTeaserHighChance
	bonusHighChance       = math.Round((bonusBuy*100*bonusChance*n1)/n2) / 100 // double the base chance.
	bonusSelIndexes       = util.Indexes{bonusGameType1, bonusGameType2, bonusGameType3}
	bonusSelWeights       = []float64{65, 25, 10}
	bonusSelWeighting     = util.AcquireWeightingDedup3().AddWeights(bonusSelIndexes, bonusSelWeights)
	bonusSelCount         = 3
	bonusWheelIndexes     = util.Indexes{bonusGameType4, bonusGameType5, bonusGameType4, bonusGameType6, bonusGameType4, bonusGameType5}
	bonusWheelWeights     = []float64{1, 1, 1, 1, 1, 1}
	bonusWheelWeighting   = util.AcquireWeighting().AddWeights(bonusWheelIndexes, bonusWheelWeights)

	scatterNudgeCount  = uint8(2)
	scatterNudgeChance = 20.0
	scatterNudgeTease  = 15.0

	multiplierScale = []float64{1, 1, 1, 2, 2, 2, 5, 5, 5, 10}

	choice1Key     = "bonus-bet"
	choice1Values  = []string{"active"}
	choice1Results = []int{bonusIncreaseID}
	choice2Key     = "selection"
	choice2Values  = []string{"left", "middle", "right"}
	choice2Results = []int{bonusChoice1, bonusChoice2, bonusChoice3}

	// weights to remove payouts based on total payout factor during first spins.
	removeBandsFirst = []comp.RemovePayoutBand{
		{MinPayout: 0, MaxPayout: 0.5, RemoveChance: 70},
		{MinPayout: 0.5, MaxPayout: 2, RemoveChance: 75},
		{MinPayout: 2, MaxPayout: 5, RemoveChance: 45},
		{MinPayout: 5, MaxPayout: 10, RemoveChance: 35},
		{MinPayout: 10, MaxPayout: 25, RemoveChance: 75},
		{MinPayout: 25, MaxPayout: maxPayout, RemoveChance: 50},
	}

	n01 = comp.WithName("10")
	n02 = comp.WithName("Jack")
	n03 = comp.WithName("Queen")
	n04 = comp.WithName("King")
	n05 = comp.WithName("Ace")
	n06 = comp.WithName("Robin Hood")
	n07 = comp.WithName("Jester")
	n08 = comp.WithName("Queen Owl")
	n09 = comp.WithName("King Owl")
	n10 = comp.WithName("Wild")
	n11 = comp.WithName("Bonus")

	r01 = comp.WithResource("l5")
	r02 = comp.WithResource("l4")
	r03 = comp.WithResource("l3")
	r04 = comp.WithResource("l2")
	r05 = comp.WithResource("l1")
	r06 = comp.WithResource("h4")
	r07 = comp.WithResource("h3")
	r08 = comp.WithResource("h2")
	r09 = comp.WithResource("h1")
	r10 = comp.WithResource("wild")
	r11 = comp.WithResource("bonus")

	p01 = comp.WithPayouts(0, 0, 0.2, 2, 4)
	p02 = comp.WithPayouts(0, 0, 0.3, 2, 5)
	p03 = comp.WithPayouts(0, 0, 0.5, 4, 8)
	p04 = comp.WithPayouts(0, 0, 1, 6, 12)
	p05 = comp.WithPayouts(0, 0, 1, 6, 12)
	p06 = comp.WithPayouts(0, 0, 2.5, 8, 17.5)
	p07 = comp.WithPayouts(0, 0, 2.5, 8, 17.5)
	p08 = comp.WithPayouts(0, 0, 5, 12.5, 25)
	p09 = comp.WithPayouts(0, 0, 6, 15, 30)
	p10 = comp.WithPayouts(0, 0, 7.5, 20, 40)
	p11 = comp.WithPayouts()

	p11111 = comp.NewPayline(id01, rows, 1, 1, 1, 1, 1)
	p00000 = comp.NewPayline(id02, rows, 0, 0, 0, 0, 0)
	p22222 = comp.NewPayline(id03, rows, 2, 2, 2, 2, 2)
	p01210 = comp.NewPayline(id04, rows, 0, 1, 2, 1, 0)
	p21012 = comp.NewPayline(id05, rows, 2, 1, 0, 1, 2)
	p10001 = comp.NewPayline(id06, rows, 1, 0, 0, 0, 1)
	p12221 = comp.NewPayline(id07, rows, 1, 2, 2, 2, 1)
	p00122 = comp.NewPayline(id08, rows, 0, 0, 1, 2, 2)
	p22100 = comp.NewPayline(id09, rows, 2, 2, 1, 0, 0)
	p12101 = comp.NewPayline(id10, rows, 1, 2, 1, 0, 1)
	p10121 = comp.NewPayline(id11, rows, 1, 0, 1, 2, 1)
	p01110 = comp.NewPayline(id12, rows, 0, 1, 1, 1, 0)
	p21112 = comp.NewPayline(id13, rows, 2, 1, 1, 1, 2)
	p01010 = comp.NewPayline(id14, rows, 0, 1, 0, 1, 0)
	p21212 = comp.NewPayline(id15, rows, 2, 1, 2, 1, 2)
	p11011 = comp.NewPayline(id16, rows, 1, 1, 0, 1, 1)
	p11211 = comp.NewPayline(id17, rows, 1, 1, 2, 1, 1)
	p00200 = comp.NewPayline(id18, rows, 0, 0, 2, 0, 0)
	p22022 = comp.NewPayline(id19, rows, 2, 2, 0, 2, 2)
	p02220 = comp.NewPayline(id20, rows, 0, 2, 2, 2, 0)

	ids       = [symbolCount]util.Index{id01, id02, id03, id04, id05, id06, id07, id08, id09, id10, id11}
	names     = [symbolCount]comp.SymbolOption{n01, n02, n03, n04, n05, n06, n07, n08, n09, n10, n11}
	resources = [symbolCount]comp.SymbolOption{r01, r02, r03, r04, r05, r06, r07, r08, r09, r10, r11}
	payouts   = [symbolCount]comp.SymbolOption{p01, p02, p03, p04, p05, p06, p07, p08, p09, p10, p11}

	paylines = comp.Paylines{
		p11111, p00000, p22222, p01210, p21012, p12221, p10001, p00122, p22100, p12101,
		p10121, p01110, p21112, p01010, p21212, p11011, p11211, p00200, p22022, p02220}

	flag0 = comp.NewRoundFlag(flagBonusIncrease, "bonus buy")
	flag1 = comp.NewRoundFlag(flagBonusSelector, "instant bonus choice")
	flag2 = comp.NewRoundFlag(flagBonusGameType, "bonus game type")
	flag3 = comp.NewRoundFlag(flagFreeSpins, "free spins count")
	flags = comp.RoundFlags{flag0, flag1, flag2, flag3}
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
	// bonus bet feature (flag 0).
	bonusIncrease := comp.NewPlayerChoiceAction(flagBonusIncrease, choice1Key, choice1Values, choice1Results)
	bonusIncrease.WithTestChoicesFilters(comp.OnRound)
	bonusIncrease.Describe(bonusIncreaseID, "bonus bet feature (flag 0)")

	// record instant bonus player choice (flag 1).
	bonusChoice := comp.NewPlayerChoiceAction(flagBonusSelector, choice2Key, choice2Values, choice2Results)
	bonusChoice.WithTestChoicesFilters(comp.OnRoundResume)
	bonusChoice.Describe(bonusChoiceID, "bonus choice (flag 1)")

	// generate new instant bonus outcome (flag 2).
	bonusSelector := comp.NewBonusSelectorAction(bonusSelWeighting, bonusSelCount, flagBonusSelector, flagBonusGameType)
	bonusSelector.WithTriggerFilters(comp.OnRoundFlagAbove(flagBonusSelector, 0))
	bonusSelector.Describe(bonusSelectorID, "bonus selector (flag 2)")

	// trigger instant bonus feature (without bonus buy).
	instantBonus := comp.NewInstantBonusAction(bonusChance).WithPlayerChoice(choice2Key, choice2Values...)
	instantBonus.Describe(instantBonusID, "instant bonus - base")
	bonusTeaser := comp.NewInstantBonusAction(bonusTeaserChance).WithTease().WithAlternate(instantBonus)
	bonusTeaser.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusSelector, 0), comp.OnRoundFlagValue(flagBonusIncrease, 0))
	bonusTeaser.Describe(bonusTeaserID, "instant bonus tease - base")

	// trigger instant bonus feature (with bonus buy).
	instantBonusHigh := comp.NewInstantBonusAction(bonusHighChance).WithPlayerChoice(choice2Key, choice2Values...)
	instantBonusHigh.Describe(instantBonusHighID, "instant bonus - bonus bet")
	bonusTeaserHigh := comp.NewInstantBonusAction(bonusTeaserHighChance).WithTease().WithAlternate(instantBonusHigh)
	bonusTeaserHigh.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusSelector, 0), comp.OnRoundFlagValue(flagBonusIncrease, bonusIncreaseID))
	bonusTeaserHigh.Describe(bonusTeaserHighID, "instant bonus tease - bonus bet")

	// remove payouts for instant bonus type 3.
	removePayouts3 := comp.NewRemovePayoutsAction(0, math.MaxFloat64, 100, direction, 1, false, false)
	removePayouts3.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusGameType, bonusGameType3))
	removePayouts3.Describe(removePayouts3ID, "remove payouts - bonus type 3")

	// force scatters for instant bonus type 3 to collect the bonus wheel.
	scatter3 := comp.NewGenerateSymbolAction(scatter, []float64{100, 100, 100}, 1, 3, 5).GenerateNoDupes()
	scatter3.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusGameType, bonusGameType3))
	scatter3.Describe(scatters3ID, "force scatters - bonus type 3")

	// generate scatter symbols RTP 92.
	scatters92a := comp.NewGenerateSymbolAction(scatter, scatterChances92a, 1, 3, 5).GenerateNoDupes()
	scatters92a.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusGameType, 0))
	scatters92a.Describe(scatters92aID, "generate scatters - first spin base - RTP 92")
	scatters92b := comp.NewGenerateSymbolAction(scatter, scatterChances92b, 1, 3, 5).GenerateNoDupes()
	scatters92b.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusGameType, 0))
	scatters92b.Describe(scatters92bID, "generate scatters - first spin bonus - RTP 92")
	scatters92c := comp.NewGenerateSymbolAction(scatter, scatterChances92c, 1, 3, 5).GenerateNoDupes()
	scatters92c.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusGameType, bonusGameType6))
	scatters92c.Describe(scatters92cID, "generate scatters - bonus type 6 - RTP 92")
	scatters92d := comp.NewGenerateSymbolAction(scatter, scatterChances92d, 1, 3, 5).GenerateNoDupes()
	scatters92d.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusGameType, bonusGameType6))
	scatters92d.Describe(scatters92dID, "generate scatters - bonus type 6 - RTP 92")

	// generate scatter symbols RTP 94.
	scatters94a := comp.NewGenerateSymbolAction(scatter, scatterChances94a, 1, 3, 5).GenerateNoDupes()
	scatters94a.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusGameType, 0))
	scatters94a.Describe(scatters94aID, "generate scatters - first spin base - RTP 94")
	scatters94b := comp.NewGenerateSymbolAction(scatter, scatterChances94b, 1, 3, 5).GenerateNoDupes()
	scatters94b.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusGameType, 0))
	scatters94b.Describe(scatters94bID, "generate scatters - first spin bonus - RTP 94")
	scatters94c := comp.NewGenerateSymbolAction(scatter, scatterChances94c, 1, 3, 5).GenerateNoDupes()
	scatters94c.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusGameType, bonusGameType6))
	scatters94c.Describe(scatters94cID, "generate scatters - bonus type 6 - RTP 94")
	scatters94d := comp.NewGenerateSymbolAction(scatter, scatterChances94d, 1, 3, 5).GenerateNoDupes()
	scatters94d.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusGameType, bonusGameType6))
	scatters94d.Describe(scatters94dID, "generate scatters - bonus type 6 - RTP 94")

	// generate scatter symbols RTP 96.
	scatters96a := comp.NewGenerateSymbolAction(scatter, scatterChances96a, 1, 3, 5).GenerateNoDupes()
	scatters96a.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusGameType, 0))
	scatters96a.Describe(scatters96aID, "generate scatters - first spin base - RTP 96")
	scatters96b := comp.NewGenerateSymbolAction(scatter, scatterChances96b, 1, 3, 5).GenerateNoDupes()
	scatters96b.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusGameType, 0))
	scatters96b.Describe(scatters96bID, "generate scatters - first spin bonus - RTP 96")
	scatters96c := comp.NewGenerateSymbolAction(scatter, scatterChances96c, 1, 3, 5).GenerateNoDupes()
	scatters96c.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusGameType, bonusGameType6))
	scatters96c.Describe(scatters96cID, "generate scatters - bonus type 6 - RTP 96")
	scatters96d := comp.NewGenerateSymbolAction(scatter, scatterChances96d, 1, 3, 5).GenerateNoDupes()
	scatters96d.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusGameType, bonusGameType6))
	scatters96d.Describe(scatters96dID, "generate scatters - bonus type 6 - RTP 96")

	// generate wild symbols RTP 92.
	wilds92a := comp.NewGenerateSymbolAction(wild, wildChances92a)
	wilds92a.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusGameType, 0))
	wilds92a.Describe(wilds92aID, "generate wilds - first spin - RTP 92")
	wilds92b := comp.NewGenerateSymbolAction(wild, wildChances92b)
	wilds92b.WithTriggerFilters(comp.OnRoundFlagValues(flagBonusGameType, bonusGameType1, bonusGameType4))
	wilds92b.Describe(wilds92bID, "generate wilds - bonus type 1/4 - RTP 92")
	wilds92c := comp.NewGenerateSymbolAction(wild, wildChances92c).GenerateNoDupes()
	wilds92c.WithTriggerFilters(comp.OnRoundFlagValues(flagBonusGameType, bonusGameType2, bonusGameType5))
	wilds92c.Describe(wilds92cID, "generate wilds - bonus type 2/5 - RTP 92")
	wilds92d := comp.NewGenerateSymbolAction(wild, wildChances92d)
	wilds92d.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusGameType, bonusGameType6))
	wilds92d.Describe(wilds92dID, "generate wilds - bonus type 6 - RTP 92")

	// generate wild symbols RTP 94.
	wilds94a := comp.NewGenerateSymbolAction(wild, wildChances94a)
	wilds94a.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusGameType, 0))
	wilds94a.Describe(wilds94aID, "generate wilds - first spin - RTP 94")
	wilds94b := comp.NewGenerateSymbolAction(wild, wildChances94b)
	wilds94b.WithTriggerFilters(comp.OnRoundFlagValues(flagBonusGameType, bonusGameType1, bonusGameType4))
	wilds94b.Describe(wilds94bID, "generate wilds - bonus type 1/4 - RTP 94")
	wilds94c := comp.NewGenerateSymbolAction(wild, wildChances94c).GenerateNoDupes()
	wilds94c.WithTriggerFilters(comp.OnRoundFlagValues(flagBonusGameType, bonusGameType2, bonusGameType5))
	wilds94c.Describe(wilds94cID, "generate wilds - bonus type 2/5 - RTP 94")
	wilds94d := comp.NewGenerateSymbolAction(wild, wildChances94d)
	wilds94d.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusGameType, bonusGameType6))
	wilds94d.Describe(wilds94dID, "generate wilds - bonus type 6 - RTP 94")

	// generate wild symbols RTP 96.
	wilds96a := comp.NewGenerateSymbolAction(wild, wildChances96a)
	wilds96a.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusGameType, 0))
	wilds96a.Describe(wilds96aID, "generate wilds - first spin - RTP 96")
	wilds96b := comp.NewGenerateSymbolAction(wild, wildChances96b)
	wilds96b.WithTriggerFilters(comp.OnRoundFlagValues(flagBonusGameType, bonusGameType1, bonusGameType4))
	wilds96b.Describe(wilds96bID, "generate wilds - bonus type 1/4 - RTP 96")
	wilds96c := comp.NewGenerateSymbolAction(wild, wildChances96c).GenerateNoDupes()
	wilds96c.WithTriggerFilters(comp.OnRoundFlagValues(flagBonusGameType, bonusGameType2, bonusGameType5))
	wilds96c.Describe(wilds96cID, "generate wilds - bonus type 2/5 - RTP 96")
	wilds96d := comp.NewGenerateSymbolAction(wild, wildChances96d)
	wilds96d.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusGameType, bonusGameType6))
	wilds96d.Describe(wilds96dID, "generate wilds - bonus type 6 - RTP 96")

	// reel nudge for scatters.
	scatterNudge := comp.NewNudgeAction(scatter, scatterNudgeCount, comp.NudgeVertical, scatterNudgeChance).GenerateNoDupes()
	scatterNudge.WithTriggerFilters(comp.OnNotRoundFlagValue(flagBonusGameType, bonusGameType3))
	scatterNudge.WithTease(scatterNudgeTease)
	scatterNudge.Describe(reelNudgeID, "scatter nudge")

	// expand wilds across reels for bonus game types 2 & 5.
	wildExpand := comp.NewWildExpansion(0, false, wild, 1, false, true, true)
	wildExpand.WithTriggerFilters(comp.OnRoundFlagValues(flagBonusGameType, bonusGameType2, bonusGameType5))
	wildExpand.Describe(wildExpandID, "expand wilds - boonus type 2/5")

	// round multiplier.
	multiplier := comp.NewMultiplierScaleAction(scatter, 1, multiplierScale...)
	multiplier.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusGameType, bonusGameType6))
	multiplier.Describe(multiplierID, "multiplier scale")

	// calculate winlines.
	winlines := comp.NewPaylinesAction()
	winlines.Describe(winlinesID, "winlines")

	// remove payouts by band.
	removePayoutsFirst := comp.NewRemovePayoutBandsAction(1, direction, true, false, removeBandsFirst)
	removePayoutsFirst.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusGameType, 0))
	removePayoutsFirst.Describe(removePayouts1ID, "remove payouts by band (first spins)")

	// bonus wheel game updating flag 2.
	bonusWheel := comp.NewScatterBonusWheelAction(scatter, 3, flagBonusGameType, bonusWheelWeighting)
	bonusWheel.WithTriggerFilters(comp.OnRoundFlagValues(flagBonusGameType, 0, bonusGameType3))
	bonusWheel.Describe(bonusWheelID, "bonus wheel game")

	// award free spins.
	freeSpins1 := comp.NewScatterFreeSpinsAction(freeSpinsGameType4, false, scatter, 0, false)
	freeSpins1.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusGameType, bonusGameType4))
	freeSpins1.Describe(freeSpins1ID, "free spins - bonus type 4")
	freeSpins2 := comp.NewScatterFreeSpinsAction(freeSpinsGameType5, false, scatter, 0, false)
	freeSpins2.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusGameType, bonusGameType5))
	freeSpins2.Describe(freeSpins2ID, "free spins - bonus type 5")
	freeSpins3 := comp.NewScatterFreeSpinsAction(freeSpinsGameType6, false, scatter, 0, false)
	freeSpins3.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusGameType, bonusGameType6))
	freeSpins3.Describe(freeSpins3ID, "free spins - bonus type 6 - first")

	freeSpins4a := comp.NewScatterFreeSpinsAction(freeSpinsGameType6free, false, scatter, 1, false)
	freeSpins4a.Describe(freeSpins4aID, "free spins - bonus type 6 - free - 1 scatter")
	freeSpins4b := comp.NewScatterFreeSpinsAction(freeSpinsGameType6free*2, false, scatter, 2, false).WithAlternate(freeSpins4a)
	freeSpins4b.Describe(freeSpins4bID, "free spins - bonus type 6 - free - 2 scatters")
	freeSpins4c := comp.NewScatterFreeSpinsAction(freeSpinsGameType6free*3, false, scatter, 3, false).WithAlternate(freeSpins4b)
	freeSpins4c.Describe(freeSpins4cID, "free spins - bonus type 6 - free - 3 scatters")
	freeSpins4d := comp.NewScatterFreeSpinsAction(freeSpinsGameType6free*4, false, scatter, 4, false).WithAlternate(freeSpins4c)
	freeSpins4d.Describe(freeSpins4dID, "free spins - bonus type 6 - free - 4 scatters")
	freeSpins4e := comp.NewScatterFreeSpinsAction(freeSpinsGameType6free*5, false, scatter, 5, false).WithAlternate(freeSpins4d)
	freeSpins4e.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusGameType, bonusGameType6))
	freeSpins4e.Describe(freeSpins4eID, "free spins - bonus type 6 - free - 5 scatters")

	freeSpins5a := comp.NewScatterFreeSpinsAction(freeSpinsGameType6free, false, wild, 1, false)
	freeSpins5a.Describe(freeSpins5aID, "free spins - bonus type 6 - free - 1 wild")
	freeSpins5b := comp.NewScatterFreeSpinsAction(freeSpinsGameType6free*2, false, wild, 2, false).WithAlternate(freeSpins5a)
	freeSpins5b.Describe(freeSpins5bID, "free spins - bonus type 6 - free - 2 wilds")
	freeSpins5c := comp.NewScatterFreeSpinsAction(freeSpinsGameType6free*3, false, wild, 3, false).WithAlternate(freeSpins5b)
	freeSpins5c.Describe(freeSpins5cID, "free spins - bonus type 6 - free - 3 wilds")
	freeSpins5d := comp.NewScatterFreeSpinsAction(freeSpinsGameType6free*4, false, wild, 4, false).WithAlternate(freeSpins5c)
	freeSpins5d.Describe(freeSpins5dID, "free spins - bonus type 6 - free - 4 wilds")
	freeSpins5e := comp.NewScatterFreeSpinsAction(freeSpinsGameType6free*5, false, wild, 5, false).WithAlternate(freeSpins5d)
	freeSpins5e.WithTriggerFilters(comp.OnRoundFlagValue(flagBonusGameType, bonusGameType6))
	freeSpins5e.Describe(freeSpins5eID, "free spins - bonus type 6 - free - 5 wilds")

	// update round flag 3 marking sequence of free spin.
	freeSpinsFlag := comp.NewRoundFlagIncreaseAction(flagFreeSpins)
	freeSpinsFlag.Describe(freeSpinsFlagID, "count number of free spins (flag 3)")

	actionsAall := comp.SpinActions{bonusIncrease, bonusChoice, bonusSelector, bonusTeaser, bonusTeaserHigh}
	actionsAfirst := comp.SpinActions{bonusIncrease, bonusChoice, bonusSelector, bonusTeaser, bonusTeaserHigh}
	actionsAfree := comp.SpinActions{bonusIncrease, bonusChoice}

	actionsBall := comp.SpinActions{removePayouts3, scatter3, scatterNudge, wildExpand, multiplier, winlines, removePayoutsFirst, bonusWheel,
		freeSpins1, freeSpins2, freeSpins3, freeSpins4e, freeSpins5e, freeSpinsFlag}
	actionsBfirst := comp.SpinActions{removePayouts3, scatter3, scatterNudge, wildExpand, winlines, removePayoutsFirst, bonusWheel,
		freeSpins1, freeSpins2, freeSpins3}
	actionsBfree := comp.SpinActions{scatterNudge, wildExpand, multiplier, winlines,
		freeSpins4e, freeSpins5e, freeSpinsFlag}

	actions92all = append(append(actionsAall, scatters92a, scatters92b, scatters92c, scatters92d, wilds92a, wilds92b, wilds92c, wilds92d), actionsBall...)
	actions92first = append(append(actionsAfirst, scatters92a, wilds92a, wilds92b, wilds92c, wilds92d), actionsBfirst...)
	actions92free = append(append(actionsAfree, scatters92c, wilds92b, wilds92c, wilds92d), actionsBfree...)
	actions92firstBB = append(append(actionsAfirst, scatters92b, wilds92a, wilds92b, wilds92c, wilds92d), actionsBfirst...)
	actions92freeBB = append(append(actionsAfree, scatters92d, wilds92b, wilds92c, wilds92d), actionsBfree...)
	actions94all = append(append(actionsAall, scatters94a, scatters94b, scatters94c, scatters94d, wilds94a, wilds94b, wilds94c, wilds94d), actionsBall...)
	actions94first = append(append(actionsAfirst, scatters94a, wilds94a, wilds94b, wilds94c, wilds94d), actionsBfirst...)
	actions94free = append(append(actionsAfree, scatters94c, wilds94b, wilds94c, wilds94d), actionsBfree...)
	actions94firstBB = append(append(actionsAfirst, scatters94b, wilds94a, wilds94b, wilds94c, wilds94d), actionsBfirst...)
	actions94freeBB = append(append(actionsAfree, scatters94d, wilds94b, wilds94c, wilds94d), actionsBfree...)
	actions96all = append(append(actionsAall, scatters96a, scatters96b, scatters96c, scatters96d, wilds96a, wilds96b, wilds96c, wilds96d), actionsBall...)
	actions96first = append(append(actionsAfirst, scatters96a, wilds96a, wilds96b, wilds96c, wilds96d), actionsBfirst...)
	actions96free = append(append(actionsAfree, scatters96c, wilds96b, wilds96c, wilds96d), actionsBfree...)
	actions96firstBB = append(append(actionsAfirst, scatters96b, wilds96a, wilds96b, wilds96c, wilds96d), actionsBfirst...)
	actions96freeBB = append(append(actionsAfree, scatters96d, wilds96b, wilds96c, wilds96d), actionsBfree...)
}

func initSlots(target float64, weights [symbolCount]comp.SymbolOption, actions1, actions2, actions3, actions4 []comp.SpinActioner) *comp.Slots {
	ss := make([]*comp.Symbol, symbolCount)
	for ix := range ss {
		id := ids[ix]
		switch id {
		case wild:
			ss[ix] = comp.NewSymbol(id, names[ix], resources[ix], payouts[ix], weights[ix], comp.WithKind(comp.Wild))
		case scatter:
			ss[ix] = comp.NewSymbol(id, names[ix], resources[ix], payouts[ix], weights[ix], comp.WithKind(comp.Scatter))
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
		comp.NoRepeat(noRepeat),
		comp.WithSymbols(s),
		comp.WithPaylines(direction, true, paylines...),
		comp.WithPlayerChoice(),
		comp.WithRoundMultiplier(),
		comp.WithBonusBuy(),
		comp.MaxPayout(maxPayout),
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

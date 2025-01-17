package slots

import (
	"sort"

	comp "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	rslt "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func (r *Regular) BuildAnimations(result *rslt.Result) rslt.Animations {
	b := animationsBuilder{game: r, result: result, payouts: result.Payouts}
	var ok bool
	if b.data, ok = result.Data.(*comp.SpinResult); ok {
		b.build(result.Animations)
	}
	return b.out
}

type animationsBuilder struct {
	maxID   utils.Index
	game    *Regular
	slots   *comp.Slots
	symbols *comp.SymbolSet
	data    *comp.SpinResult
	result  *rslt.Result
	payouts rslt.Payouts
	out     rslt.Animations
}

func (b *animationsBuilder) build(in rslt.Animations) {
	b.out = rslt.PurgeAnimations(in, 16)
	b.slots = b.game.slots
	b.symbols = b.slots.Symbols()
	if s := b.slots.AltSymbols(); s != nil {
		b.symbols = s
	}
	b.maxID = b.symbols.GetMaxSymbolID()

	b.anticipations()
	b.bombs()
	b.superShapes()
	b.shooters()
	b.reels()

	if b.game.havePaylines {
		b.paylines()
		b.bonusSymbolPayout()
	}
	if b.game.haveAllPaylines {
		b.allPaylines()
	}
	if b.game.haveScatterPays {
		b.scatterPayouts()
	}
	if b.game.haveClusterPays {
		b.clusterPayouts()
	}

	b.freeGames()
	b.awards()
	b.stickies()
	b.clearings()
	b.cascades()
	b.refills()
	b.playerChoice()
}

func (b *animationsBuilder) anticipations() {

	// TODO: FUTURE

}

var bombGrid = []comp.Offsets{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 0}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}

func (b *animationsBuilder) bombs() {
	bomb := b.symbols.GetBombSymbol()
	if bomb == nil || len(b.data.AfterExpand()) == 0 {
		return
	}

	reels, rows := b.slots.ReelCount(), b.slots.RowCount()
	symbols := b.data.Initial()
	bombID := bomb.ID()

	var offset int
	for reel := 0; reel < reels; reel++ {
		for row := 0; row < rows; row++ {
			if symbols[offset] == bombID {
				b.out = append(b.out, comp.AcquireBombEvent(b.data, reels, rows, comp.Offsets{reel, row}, bombGrid))
			}
			offset++
		}
	}
}

func (b *animationsBuilder) superShapes() {
	if b.data.Kind() == comp.SuperSpin {
		b.out = append(b.out, comp.AcquireSuperEvent(b.data))
	}
}

func (b *animationsBuilder) shooters() {

	// TODO: FUTURE

}

func (b *animationsBuilder) reels() {
	hot := b.data.Hot()
	for _, reel := range hot {
		b.out = append(b.out, comp.AcquireHotReelEvent(reel, true))
	}

	// TODO: FUTURE
	// locked := b.data.Locked()
	// for _, reel := range locked {
	// 	b.out = append(b.out, comp.AcquireLockedReel(reel, true))
	// }
}

func (b *animationsBuilder) paylines() {
	for ix := range b.payouts {
		if p := b.payouts[ix]; p.Kind() == rslt.SlotWinline {
			if sp, ok := p.(*comp.SpinPayout); ok {
				b.out = append(b.out, comp.AcquirePayoutEvent(sp))
			}
		}
	}
}

func (b *animationsBuilder) allPaylines() {
	reels, rows := b.slots.ReelCount(), b.slots.RowCount()

	events := make(rslt.Animations, 0, 32)
	for id := utils.Index(1); id <= b.maxID; id++ {
		for ix := range b.payouts {
			if p, ok := b.payouts[ix].(*comp.SpinPayout); ok && p.Symbol() == id {
				events = append(events, comp.AcquireAllPaylinesEvents(id, reels, rows, b.payouts)...)
				break
			}
		}
	}

	switch len(events) {
	case 0:
		return
	case 1:
		b.out = append(b.out, events[0])
		return
	}

	sort.Slice(events, func(i, j int) bool {
		o1, ok1 := events[i].(*comp.AllPaylinesEvent)
		o2, ok2 := events[j].(*comp.AllPaylinesEvent)
		if ok1 && ok2 {
			return o1.Factor() > o2.Factor()
		}
		return false
	})
	b.out = append(b.out, events...)
}

func (b *animationsBuilder) scatterPayouts() {
	for ix := range b.payouts {
		if p := b.payouts[ix]; p.Kind() == rslt.SlotBombScatters {
			if sp, ok := p.(*comp.SpinPayout); ok {
				b.out = append(b.out, comp.AcquirePayoutEvent(sp))
			}
		}
	}
	for ix := range b.payouts {
		if p := b.payouts[ix]; p.Kind() == rslt.SlotScatters {
			if sp, ok := p.(*comp.SpinPayout); ok {
				b.out = append(b.out, comp.AcquirePayoutEvent(sp))
			}
		}
	}
	for ix := range b.payouts {
		if p := b.payouts[ix]; p.Kind() == rslt.SlotWilds {
			if sp, ok := p.(*comp.SpinPayout); ok {
				b.out = append(b.out, comp.AcquirePayoutEvent(sp))
			}
		}
	}
}

func (b *animationsBuilder) clusterPayouts() {
	for ix := range b.payouts {
		if p := b.payouts[ix]; p.Kind() == rslt.SlotCluster {
			if sp, ok := p.(*comp.SpinPayout); ok {
				b.out = append(b.out, comp.AcquirePayoutEvent(sp))
			}
		}
	}
}

func (b *animationsBuilder) bonusSymbolPayout() {
	for ix := range b.payouts {
		if p := b.payouts[ix]; p.Kind() == rslt.SlotBonusSymbol {
			if sp, ok := p.(*comp.SpinPayout); ok {
				e := comp.AcquirePayoutEvent(sp).(*comp.PayoutEvent)
				if b.game.havePaylines {
					e.SetPaylines(b.game.slots.Paylines())
				}
				b.out = append(b.out, e)
			}
		}
	}
}

func (b *animationsBuilder) freeGames() {
	if n := b.result.AwardedFreeGames; n > 0 {
		b.out = append(b.out, comp.AcquireFreeGameEvent(n))
	}
}

func (b *animationsBuilder) awards() {

	// TODO: FUTURE

}

func (b *animationsBuilder) stickies() {
	if b.data.HaveSticky() {
		b.out = append(b.out, comp.AcquireStickyEvent(b.data))
	}
}

func (b *animationsBuilder) clearings() {
	if b.data.HaveClearing() {
		b.out = append(b.out, comp.AcquireClearEvent(b.data))
	}
}

func (b *animationsBuilder) cascades() {

	// TODO: FUTURE

}

func (b *animationsBuilder) refills() {

	// TODO: FUTURE

}

func (b *animationsBuilder) playerChoice() {
	if b.game.makeChoice {
		b.out = append(b.out, comp.AcquireChoiceRequest())
	}
}

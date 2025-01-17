package slots

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/tg"

	analyse "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/analysis/metrics/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/consts"
	wheel2 "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/wheel"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

// AcquireRounds instantiates a new player rounds metrics from the memory pool.
func AcquireRounds(gameNR tg.GameNR, playerID string, balance int64, reelCount, rowCount int, doubleSpin bool, maxPayout float64,
	symbols *slots.SymbolSet, actions slots.SpinActions, paylines []*slots.Payline, flags slots.RoundFlags) *Rounds {
	maxSymbol := symbols.GetMaxSymbolID() + 1
	maxAction := actions.GetMaxID() + 1

	r := roundsProvider.Acquire().(*Rounds)
	r.doubleSpin = doubleSpin
	r.maxSymbol = maxSymbol
	r.gameNR = gameNR
	r.startBalance = balance
	r.reelCount = reelCount
	r.rowCount = rowCount
	r.maxPayout = maxPayout
	r.ss = symbols
	r.pl = paylines

	r.LowestBalance = balance
	r.HighestBalance = balance
	r.Balance = balance
	r.PlayerID = playerID

	r.AllRounds = analyse.NewRounds(symbols)
	r.FirstPayouts = NewPayouts(rowCount, maxSymbol, paylines, false)

	if r.BonusRounds == nil {
		r.BonusRounds = make(map[analyse.BonusKind]*analyse.Rounds, 16)
	}
	if r.BonusPayouts == nil {
		r.BonusPayouts = make(map[analyse.BonusKind]*Payouts, 16)
	}

	r.Symbols = analyse.PurgeSymbols(r.Symbols, int(maxSymbol))[:maxSymbol]
	for id := utils.Index(0); id < maxSymbol; id++ {
		if s := symbols.GetSymbol(id); s != nil {
			r.Symbols[id] = analyse.NewSymbol(id, s.Name(), s.Resource(), reelCount)
		}
	}

	r.Actions = analyse.PurgeActions(r.Actions, maxAction)[:maxAction]
	for ix := range actions {
		a := actions[ix]
		r.Actions[a.ID()] = analyse.NewAction(a.ID(), a.Name(), a.Kind(), a.Config())

		alt := a.Alternate()
		for alt != nil {
			r.Actions[alt.ID()] = analyse.NewAction(alt.ID(), alt.Name(), alt.Kind(), alt.Config())
			alt = alt.Alternate()
		}
	}

	if len(flags) > 0 {
		m := flags.GetMaxID() + 1
		r.RoundFlags = analyse.PurgeRoundFlags(r.RoundFlags, m)[:m]

		for ix := range flags {
			f := flags[ix]
			r.RoundFlags[f.ID()] = analyse.NewRoundFlag(f.ID(), f.Name())
		}
	}

	return r
}

func (r *Rounds) WithOptions(noPaylines, noPayouts, noSymbols, noBest, noSpins, noCounts, noBalance bool) *Rounds {
	r.noPaylines = noPaylines
	r.noPayouts = noPayouts
	r.noSymbols = noSymbols
	r.noBest = noBest
	r.noSpins = noSpins
	r.noCounts = noCounts
	r.noBalance = noBalance

	r.FirstPayouts.noPaylines = noPaylines

	return r
}

// SetBestThreshold sets the payout factor at which rounds are considered best.
func (r *Rounds) SetBestThreshold(threshold float64, maxBest int) {
	r.bestThreshold = threshold
	r.maxBest = maxBest
}

// BestThreshold returns the current threshold for best rounds.
func (r *Rounds) BestThreshold() float64 {
	return r.bestThreshold
}

// BestNoFreeThreshold returns the current threshold for best rounds with no free spins.
func (r *Rounds) BestNoFreeThreshold() float64 {
	return r.bestNoFreeThreshold
}

// RTP returns the RTP (return-to-player) for the analyzed results.
func (r *Rounds) RTP() float64 {
	if r.AllRounds.Bets.Total > 0 {
		return float64(r.AllRounds.Wins.Total) * 100.0 / float64(r.AllRounds.Bets.Total)
	}
	return 0
}

// RTPnoFree returns the RTP (return-to-player) for the analyzed results with no free spins.
func (r *Rounds) RTPnoFree() float64 {
	if r.AllRounds.BetsNoFree.Total > 0 {
		return float64(r.AllRounds.WinsNoFree.Total) * 100.0 / float64(r.AllRounds.BetsNoFree.Total)
	}
	return 0
}

// RTPfree returns the RTP (return-to-player) for the analyzed results with free spins.
func (r *Rounds) RTPfree() float64 {
	if r.AllRounds.BetsFree.Total > 0 {
		return float64(r.AllRounds.WinsFree.Total) * 100.0 / float64(r.AllRounds.BetsFree.Total)
	}
	return 0
}

// HitRate returns the hit rate (number of times win over number of times played) for the analyzed results.
func (r *Rounds) HitRate() float64 {
	if r.AllRounds.Count > 0 {
		return float64(r.WinCount) * 100.0 / float64(r.AllRounds.Count)
	}
	return 0
}

// WinningProbability returns the winning probability.
func (r *Rounds) WinningProbability() float64 {
	if r.AllRounds.Count > 0 {
		return float64(r.WinCount) / float64(r.AllRounds.Count)
	}
	return 0
}

// Merge merges the given metrics into the player rounds metrics.
func (r *Rounds) Merge(other *Rounds) {
	if other.reelCount != r.reelCount || other.rowCount != r.rowCount ||
		len(r.Symbols) != len(other.Symbols) || len(r.Actions) != len(other.Actions) {
		panic(consts.MsgAnalysisNonMatchingRounds)
	}

	r.RoundCount += other.RoundCount
	r.WinCount += other.WinCount
	r.TotalSpins += other.TotalSpins
	r.RegularSpins += other.RegularSpins
	r.FirstSpins += other.FirstSpins
	r.SecondSpins += other.SecondSpins
	r.RefillSpins += other.RefillSpins
	r.WildRespins += other.WildRespins
	r.SuperSpins += other.SuperSpins
	r.SuperRefills += other.SuperRefills
	r.FreeTimes += other.FreeTimes
	r.FirstTimes += other.FirstTimes
	r.SecondTimes += other.SecondTimes
	r.FreeTimesSuper += other.FreeTimesSuper
	r.FirstTimesSuper += other.FirstTimesSuper
	r.SecondTimesSuper += other.SecondTimesSuper
	r.FreeAwarded += other.FreeAwarded
	r.FirstAwarded += other.FirstAwarded
	r.SecondAwarded += other.SecondAwarded
	r.FreeAwardedSuper += other.FreeAwardedSuper
	r.FirstAwardedSuper += other.FirstAwardedSuper
	r.SecondAwardedSuper += other.SecondAwardedSuper
	r.FreeSpins += other.FreeSpins
	r.FirstFreeSpins += other.FirstFreeSpins
	r.SecondFreeSpins += other.SecondFreeSpins
	r.SuperSpinsFree += other.SuperSpinsFree
	r.SuperRefillsFree += other.SuperRefillsFree
	r.BadSpins += other.BadSpins
	r.MaxPayouts += other.MaxPayouts
	r.PositiveBal += other.PositiveBal
	r.NegativeBal += other.NegativeBal
	r.Balance += other.Balance

	if other.HighestPayout > r.HighestPayout {
		r.HighestPayout = other.HighestPayout
	}
	if other.LowestBalance < r.LowestBalance {
		r.LowestBalance = other.LowestBalance
	}
	if other.HighestBalance > r.HighestBalance {
		r.HighestBalance = other.HighestBalance
	}

	r.AllRounds.Merge(other.AllRounds)
	r.FirstPayouts.Merge(other.FirstPayouts)
	r.SpinsTo25x.Merge(other.SpinsTo25x)
	r.SpinsTo100x.Merge(other.SpinsTo100x)
	r.SpinsTo250x.Merge(other.SpinsTo250x)
	r.SpinsTo1000x.Merge(other.SpinsTo1000x)
	r.SpinsTo2500x.Merge(other.SpinsTo2500x)
	r.SpinsToPlusBal.Merge(other.SpinsToPlusBal)
	r.Count25x.Merge(other.Count25x)
	r.Count100x.Merge(other.Count100x)
	r.Count250x.Merge(other.Count250x)
	r.Count1000x.Merge(other.Count1000x)
	r.Count2500x.Merge(other.Count2500x)
	r.CountPlusBal.Merge(other.CountPlusBal)
	r.BonusWheel.Merge(other.BonusWheel)
	r.MultiplierMarks.Merge(other.MultiplierMarks)
	r.Multipliers.Merge(other.Multipliers)

	for k, v := range other.BonusRounds {
		br, ok := r.BonusRounds[k]
		if !ok {
			br = analyse.NewRounds(r.ss)
		}
		br.Merge(v)
		r.BonusRounds[k] = br
	}

	for k, v := range other.BonusPayouts {
		br, ok := r.BonusPayouts[k]
		if !ok {
			br = NewPayouts(r.rowCount, r.maxSymbol, r.pl, r.noPaylines)
		}
		br.Merge(v)
		r.BonusPayouts[k] = br
	}

	for key, count := range other.InstantBonus {
		r.InstantBonus[key] = r.InstantBonus[key] + count
	}
	for key, count := range other.PlayerChoice {
		r.PlayerChoice[key] = r.PlayerChoice[key] + count
	}
	for key, count := range other.Scripts {
		r.Scripts[key] = r.Scripts[key] + count
	}

	for id := range r.Symbols {
		if s := r.Symbols[id]; s != nil {
			s.Merge(other.Symbols[id])
		}
	}

	for id := range other.Actions {
		a, b := r.Actions[id], other.Actions[id]
		if b != nil {
			if a == nil {
				a = analyse.NewAction(b.ID, b.Name, b.Kind, b.Config)
				r.Actions[id] = a
			}
			a.Merge(b)
		}
	}

	for ix := range other.Best {
		r.Best = r.addBestX(r.Best, other.Best[ix], 0)
	}
	for ix := range other.BestNoFree {
		r.BestNoFree = r.addBestX(r.BestNoFree, other.BestNoFree[ix], 0)
	}

	for ix := range other.RoundFlags {
		f := other.RoundFlags[ix]
		id, m := f.ID, len(r.RoundFlags)
		for id >= m {
			r.RoundFlags = append(r.RoundFlags, nil)
			m++
		}
		if r.RoundFlags[id] == nil {
			r.RoundFlags[id] = analyse.NewRoundFlag(f.ID, f.Name)
		}
		r.RoundFlags[id].Merge(f)
	}

	r.startBalance += other.startBalance
}

// Analyse analyses the result and updates the player rounds metrics accordingly.
func (r *Rounds) Analyse(bet, bonusBet int64, res results.Results) {
	grandTotal := results.GrandTotal(res)
	if grandTotal >= r.maxPayout {
		r.MaxPayouts++
		grandTotal = r.maxPayout
	}

	if !r.noSpins || !r.noCounts {
		switch {
		case grandTotal >= 2500:
			if r.SpinsTo2500x.First() {
				r.SpinsTo2500x.Increase(r.RoundCount)
			}
			r.Count2500x.IncreaseOne()

		case grandTotal >= 1000:
			if r.SpinsTo1000x.First() {
				r.SpinsTo1000x.Increase(r.RoundCount)
			}
			r.Count1000x.IncreaseOne()

		case grandTotal >= 250:
			if r.SpinsTo250x.First() {
				r.SpinsTo250x.Increase(r.RoundCount)
			}
			r.Count250x.IncreaseOne()

		case grandTotal >= 100:
			if r.SpinsTo100x.First() {
				r.SpinsTo100x.Increase(r.RoundCount)
			}
			r.Count100x.IncreaseOne()

		case grandTotal >= 25:
			if r.SpinsTo25x.First() {
				r.SpinsTo25x.Increase(r.RoundCount)
			}
			r.Count25x.IncreaseOne()
		}
	}

	win := int64(math.Round(float64(bet) * grandTotal))
	if win > 0 {
		r.WinCount++
		if win > r.HighestPayout {
			r.HighestPayout = win
		}
	}

	if !r.noBalance {
		r.Balance -= bonusBet
		r.Balance += win

		if r.Balance < r.LowestBalance {
			r.LowestBalance = r.Balance
		}
		if r.Balance > r.HighestBalance {
			r.HighestBalance = r.Balance
		}

		if r.Balance > r.startBalance {
			if r.SpinsToPlusBal.First() {
				r.SpinsToPlusBal.Increase(r.RoundCount)
			}
			r.CountPlusBal.IncreaseOne()
			r.PositiveBal = 1
			r.NegativeBal = 0
		} else {
			r.PositiveBal = 0
			r.NegativeBal = 1
		}
	}

	m := len(res) - 1

	var wildRespin bool
	if m > 0 && (r.gameNR == tg.MGDnr || r.gameNR == tg.FPRnr) {
		if s, ok := res[0].Data.(*slots.SpinResult); ok {
			wildRespin = r.isWildRespin(s.Initial())
		}
	}

	free, refill, super := r.analyseSpins(res, wildRespin)
	bonusKind := analyse.DetermineBonusKind(r.gameNR, res)

	var aggregateBonus analyse.BonusKind
	switch bonusKind {
	case analyse.MGDBonus3, analyse.MGDBonus4:
		aggregateBonus = analyse.MGDBonus
	case analyse.MGDSuperBonus3, analyse.MGDSuperBonus4:
		aggregateBonus = analyse.MGDSuperBonus
	case analyse.FPRBonus3, analyse.FPRBonus4:
		aggregateBonus = analyse.FPRBonus
	case analyse.FPRSuperBonus3, analyse.FPRSuperBonus4:
		aggregateBonus = analyse.FPRSuperBonus
	case analyse.BTRFreeSpins3, analyse.BTRFreeSpins4, analyse.BTRFreeSpins5, analyse.BTRFreeSpins6, analyse.BTRFreeSpins33:
		aggregateBonus = analyse.BTRFreeSpins
	case analyse.OFGLevel1, analyse.OFGLevel2, analyse.OFGLevel3, analyse.OFGLevel4:
		aggregateBonus = analyse.OFGFreeSpins
	case analyse.OFGLevel1BB, analyse.OFGLevel2BB, analyse.OFGLevel3BB, analyse.OFGLevel4BB:
		aggregateBonus = analyse.OFGFreeSpinsBB
	case analyse.OFGFreeSpins, analyse.OFGFreeSpinsBB:
		panic("ohoh")
	default:
	}

	r.RoundCount++
	r.AllRounds.NewRound(bet, win, free, refill, super, res)

	br, ok := r.BonusRounds[bonusKind]
	if !ok {
		br = analyse.NewRounds(r.ss)
		r.BonusRounds[bonusKind] = br
	}
	br.NewRound(bet, win, free, refill, super, res)

	if aggregateBonus > 0 {
		br, ok = r.BonusRounds[aggregateBonus]
		if !ok {
			br = analyse.NewRounds(r.ss)
			r.BonusRounds[aggregateBonus] = br
		}
		br.NewRound(bet, win, free, refill, super, res)
	}

	if !r.noPayouts {
		if m == 0 {
			r.FirstPayouts.analyseRound(win)
		} else {
			bp, ok2 := r.BonusPayouts[bonusKind]
			if !ok2 {
				bp = NewPayouts(r.rowCount, r.maxSymbol, r.pl, r.noPaylines)
				r.BonusPayouts[bonusKind] = bp
			}
			bp.analyseRound(win)

			if aggregateBonus > 0 {
				bp, ok2 = r.BonusPayouts[aggregateBonus]
				if !ok2 {
					bp = NewPayouts(r.rowCount, r.maxSymbol, r.pl, r.noPaylines)
					r.BonusPayouts[aggregateBonus] = bp
				}
				bp.analyseRound(win)
			}
		}
	}

	for ix := range res {
		r.analyseResult(res[ix], bonusKind, aggregateBonus, ix == 0, ix == m)
	}

	r.analyseMultiplier(res[len(res)-1])

	if !r.noBest {
		if grandTotal >= r.bestThreshold {
			r.addBest(res, grandTotal)
		}

		if free == 0 && grandTotal > r.bestNoFreeThreshold {
			r.addBestNoFree(res, grandTotal)
		}
	}
}

func (r *Rounds) analyseResult(result *results.Result, bonusKind, aggregateBonus analyse.BonusKind, first, last bool) {
	data := result.Data

	if spin, ok := data.(*slots.SpinResult); ok {
		if !r.noSymbols {
			r.analyseSymbols(spin, nil, result.Payouts)
		}

		r.analyseSpinActions(spin)
		r.analyseRoundFlags(spin.RoundFlags(), last)
		r.analyseScript(spin, first)

		if !r.noPayouts {
			if bonusKind == analyse.NoFreeSpins {
				r.FirstPayouts.analyse(result)
			} else {
				bp, ok2 := r.BonusPayouts[bonusKind]
				if !ok2 {
					bp = NewPayouts(r.rowCount, r.maxSymbol, r.pl, r.noPaylines)
					r.BonusPayouts[bonusKind] = bp
				}
				bp.analyse(result)

				if aggregateBonus > 0 {
					bp, ok2 = r.BonusPayouts[aggregateBonus]
					if !ok2 {
						bp = NewPayouts(r.rowCount, r.maxSymbol, r.pl, r.noPaylines)
						r.BonusPayouts[aggregateBonus] = bp
					}
					bp.analyse(result)
				}
			}
		}

		r.analyseChoices(spin.Choices())
		return
	}

	if bonus, ok := data.(*results.InstantBonus); ok {
		key := bonus.String()
		r.InstantBonus[key] = r.InstantBonus[key] + 1

		events, _, _ := bonus.Log()
		r.analyseActions(slots.FirstSpin, events)
		return
	}

	if wheel, ok := data.(*wheel2.BonusWheelResult); ok {
		id := wheel.Result()
		r.BonusWheel.Increase(uint64(id))

		events, _, _ := wheel.Log()
		r.analyseActions(slots.FirstSpin, events)
		return
	}

	if selector, ok := data.(*results.BonusSelector); ok {
		// TODO: analyse selector?

		events, _, _ := selector.Log()
		r.analyseActions(slots.FirstSpin, events)
		return
	}

	panic(fmt.Sprintf("bad input; not a supported result: %+v", data))
}

func (r *Rounds) analyseSpins(result results.Results, wildRespin bool) (uint64, uint64, uint64) {
	var free, refill, super uint64

	for _, res := range result {
		if s, ok := res.Data.(*slots.SpinResult); ok {
			r.TotalSpins++

			if res.AwardedFreeGames > 0 {
				if wildRespin {
					r.FreeTimesSuper++
					r.FreeAwardedSuper += res.AwardedFreeGames
				} else {
					r.FreeTimes++
					r.FreeAwarded += res.AwardedFreeGames
				}

				switch s.Kind() {
				case slots.RegularSpin, slots.FirstSpin, slots.SecondSpin:
					if wildRespin {
						r.FirstTimesSuper++
						r.FirstAwardedSuper += res.AwardedFreeGames
					} else {
						r.FirstTimes++
						r.FirstAwarded += res.AwardedFreeGames
					}

				case slots.FreeSpin, slots.FirstFreeSpin, slots.SecondFreeSpin:
					if wildRespin {
						r.SecondTimesSuper++
						r.SecondAwardedSuper += res.AwardedFreeGames
					} else {
						r.SecondTimes++
						r.SecondAwarded += res.AwardedFreeGames
					}

				default:
				}
			}

			switch s.Kind() {
			case slots.RegularSpin:
				r.RegularSpins++
			case slots.FirstSpin:
				if s.SuperSymbol() > 0 && s.SuperSymbol() < utils.MaxIndex {
					r.SuperSpins++
					super++
				} else {
					r.FirstSpins++
				}
			case slots.SecondSpin:
				r.SecondSpins++
			case slots.RefillSpin:
				if wildRespin {
					r.WildRespins++
				} else {
					r.RefillSpins++
				}
				refill++
			case slots.SuperSpin:
				if res.FreeGames > 0 {
					r.SuperRefillsFree++
				} else {
					r.SuperRefills++
				}
			case slots.FreeSpin:
				r.FreeSpins++
				free++
			case slots.FirstFreeSpin:
				if s.SuperSymbol() > 0 && s.SuperSymbol() < utils.MaxIndex {
					r.SuperSpinsFree++
				} else {
					r.FirstFreeSpins++
				}
			case slots.SecondFreeSpin:
				r.SecondFreeSpins++
				free++
			default:
				r.BadSpins++
			}
		}
	}

	return free, refill, super
}

func (r *Rounds) isWildRespin(initial utils.Indexes) bool {
	// MGD - Wild Respin detector
	if len(initial) < 24 {
		return false
	}
	return initial[0] == 10 && initial[1] == 10 && initial[20] == 10 && initial[21] == 10
}

func (r *Rounds) analyseSymbols(spin *slots.SpinResult, prevResult *results.Result, payouts results.Payouts) {
	kind := spin.Kind()
	indexes := spin.Initial()
	var offset int

	var prevSpin *slots.SpinResult
	if prevResult != nil {
		prevSpin = prevResult.Data.(*slots.SpinResult)
		_ = prevSpin
	}

	switch kind {
	case slots.RegularSpin:
		if id := spin.BonusSymbol(); id != utils.MaxIndex {
			if s := r.Symbols[id]; s != nil {
				s.IncreaseBonus()
			}
		}

		for reel := 0; reel < r.reelCount; reel++ {
			if !spin.IsLocked(uint8(reel + 1)) {
				for row := 0; row < r.rowCount; row++ {
					if id := indexes[offset]; id > 0 {
						r.Symbols[id].IncreaseFirst(reel)
					}
					offset++
				}
			}
		}

	case slots.FirstSpin, slots.SuperSpin:
		if id := spin.StickySymbol(); id != utils.MaxIndex {
			if s := r.Symbols[id]; s != nil {
				s.IncreaseSticky()
			}
		}

		if id := spin.SuperSymbol(); id != utils.MaxIndex {
			if s := r.Symbols[id]; s != nil {
				s.IncreaseSuper()
			}
		}

		for reel := 0; reel < r.reelCount; reel++ {
			if !spin.IsLocked(uint8(reel + 1)) {
				for row := 0; row < r.rowCount; row++ {
					if id := indexes[offset]; id > 0 {
						r.Symbols[id].IncreaseFirst(reel)
					}
					offset++
				}
			}
		}

	case slots.SecondSpin:
		for reel := 0; reel < r.reelCount; reel++ {
			if !spin.IsLocked(uint8(reel + 1)) {
				for row := 0; row < r.rowCount; row++ {
					if id := indexes[offset]; id > 0 {
						r.Symbols[id].IncreaseSecond(reel)
					}
					offset++
				}
			}
		}

	case slots.SecondFreeSpin:
		for reel := 0; reel < r.reelCount; reel++ {
			if !spin.IsLocked(uint8(reel + 1)) {
				for row := 0; row < r.rowCount; row++ {
					if id := indexes[offset]; id > 0 {
						r.Symbols[id].IncreaseSecondFree(reel)
					}
					offset++
				}
			}
		}

	default:
		for reel := 0; reel < r.reelCount; reel++ {
			if !spin.IsLocked(uint8(reel + 1)) {
				for row := 0; row < r.rowCount; row++ {
					if id := indexes[offset]; id > 0 {
						r.Symbols[id].IncreaseFree(reel)
					}
					offset++
				}
			}
		}
	}

	for ix := range payouts {
		payout := payouts[ix].(*slots.SpinPayout)
		r.Symbols[payout.Symbol()].AddPayout(payout.Count(), 1)
	}
}

func (r *Rounds) analyseSpinActions(spin *slots.SpinResult) {
	events, _, _ := spin.Log()
	r.analyseActions(spin.Kind(), events)
}

func (r *Rounds) analyseActions(kind slots.SpinKind, events results.Events) {
	for ix := range events {
		e := events[ix]
		id := e.ID()

		a := r.Actions[id]
		if a == nil {
			a = analyse.NewAction(id, "?", "?", "?")
			r.Actions[id] = a
		}

		switch kind {
		case slots.RegularSpin, slots.FirstSpin:
			a.IncreaseFirst(e.Triggered())
		case slots.SecondSpin:
			a.IncreaseSecond(e.Triggered())
		case slots.FreeSpin, slots.FirstFreeSpin:
			a.IncreaseFree(e.Triggered())
		case slots.SecondFreeSpin:
			a.IncreaseSecondFree(e.Triggered())
		case slots.SuperSpin:
			a.IncreaseSuper(e.Triggered())
		case slots.RefillSpin:
			a.IncreaseRefill(e.Triggered())
		}
	}
}

func (r *Rounds) analyseRoundFlags(flags []int, final bool) {
	for id := range flags {
		value := flags[id]

		m := len(r.RoundFlags)
		if id >= m && value == 0 {
			continue
		}

		for id >= m {
			r.RoundFlags = append(r.RoundFlags, nil)
			m++
		}

		f := r.RoundFlags[id]
		if f == nil {
			f = analyse.NewRoundFlag(id, "???")
			r.RoundFlags[id] = f
		}

		f.Increase(value, final)
	}
}

func (r *Rounds) analyseScript(spin *slots.SpinResult, first bool) {
	if !first {
		return
	}
	id := spin.PrngLog.ScriptID()
	r.Scripts[id] = r.Scripts[id] + 1
}

func (r *Rounds) analyseChoices(choices map[string]string) {
	for k, v := range choices {
		key := k + ":" + v
		r.PlayerChoice[key] = r.PlayerChoice[key] + 1
	}
}

func (r *Rounds) analyseMultiplier(result *results.Result) {
	if spin, ok := result.Data.(*slots.SpinResult); ok {
		if m := spin.ProgressLevel(); m > 0 {
			r.MultiplierMarks.Increase(uint64(m))
		}
		if m := spin.Multiplier(); m > 0 {
			r.Multipliers.Increase(m)
		}
	}
}

func (r *Rounds) addBest(result results.Results, total float64) {
	r.Best = r.addBestX(r.Best, result, total)
}

func (r *Rounds) addBestNoFree(result results.Results, total float64) {
	r.BestNoFree = r.addBestX(r.BestNoFree, result, total)
}

// addBestX clones the given result, and inserts it into the given slice in a GrandTotal descending order.
func (r *Rounds) addBestX(best []results.Results, result results.Results, total float64) []results.Results {
	if total == 0 {
		total = results.GrandTotal2(result, r.maxPayout)
	}

	m := len(best)

	// check if the new result needs to be inserted, appended or discarded.
	var ix int
	for ix < m {
		if bt := results.GrandTotal2(best[ix], r.maxPayout); total > bt {
			break
		}
		ix++
	}

	if ix >= m && m >= r.maxBest {
		return best
	}

	n := results.CloneResults(result)

	if ix < m {
		// found a spot inside the slice.
		if m >= r.maxBest {
			// release lowest result.
			results.FreeResults(best[m-1])
			best[m-1] = nil
		} else {
			best = append(best, nil)
			m++
		}

		// move lower results down.
		for iy := m - 1; iy > ix; iy-- {
			best[iy] = best[iy-1]
		}

		// insert the new result.
		best[ix] = n
	} else if m < r.maxBest {
		// still space so just append.
		best = append(best, n)
	}

	return best
}

// BestSorted sorts the best results and returns at most maxBest entries from the top.
// Since the slice is already sorted, this is very easy :)
func (r *Rounds) BestSorted() []results.Results {
	return r.Best
}

// BestNoFreeSorted sorts the best results and returns at most maxBest entries from the top.
// Since the slice is already sorted, this is very easy :)
func (r *Rounds) BestNoFreeSorted() []results.Results {
	return r.BestNoFree
}

// Rounds contains all metrics for a set of spin rounds by the same player.
type Rounds struct {
	noPaylines          bool
	noPayouts           bool
	noSymbols           bool
	noBest              bool
	noSpins             bool
	noCounts            bool
	noBalance           bool
	doubleSpin          bool
	gameNR              tg.GameNR
	maxSymbol           utils.Index
	reelCount           int
	rowCount            int
	maxBest             int
	startBalance        int64
	RoundCount          uint64 `json:"roundCount,omitempty"`
	WinCount            uint64 `json:"winCount,omitempty"`
	TotalSpins          uint64 `json:"totalSpins,omitempty"`
	RegularSpins        uint64 `json:"regularSpins,omitempty"`
	FirstSpins          uint64 `json:"firstSpins,omitempty"`
	SecondSpins         uint64 `json:"secondSpins,omitempty"`
	RefillSpins         uint64 `json:"refillSpins,omitempty"`
	SuperSpins          uint64 `json:"superSpins,omitempty"`
	SuperRefills        uint64 `json:"superRefills,omitempty"`
	WildRespins         uint64 `json:"wildRespins,omitempty"`
	FreeTimes           uint64 `json:"freeTimes,omitempty"`
	FirstTimes          uint64 `json:"firstTimes,omitempty"`
	SecondTimes         uint64 `json:"secondTimes,omitempty"`
	FreeTimesSuper      uint64 `json:"freeTimesSuper,omitempty"`
	FirstTimesSuper     uint64 `json:"firstTimesSuper,omitempty"`
	SecondTimesSuper    uint64 `json:"secondTimesSuper,omitempty"`
	FreeAwarded         uint64 `json:"freeAwarded,omitempty"`
	FirstAwarded        uint64 `json:"firstAwarded,omitempty"`
	SecondAwarded       uint64 `json:"secondAwarded,omitempty"`
	FreeAwardedSuper    uint64 `json:"freeAwardedSuper,omitempty"`
	FirstAwardedSuper   uint64 `json:"firstAwardedSuper,omitempty"`
	SecondAwardedSuper  uint64 `json:"secondAwardedSuper,omitempty"`
	FreeSpins           uint64 `json:"freeSpins,omitempty"`
	FirstFreeSpins      uint64 `json:"firstFreeSpins,omitempty"`
	SecondFreeSpins     uint64 `json:"secondFreeSpins,omitempty"`
	SuperSpinsFree      uint64 `json:"superSpinsFree,omitempty"`
	SuperRefillsFree    uint64 `json:"superRefillsFree,omitempty"`
	BadSpins            uint64 `json:"badSpins,omitempty"`
	MaxPayouts          uint64 `json:"maxPayouts,omitempty"`
	PositiveBal         uint64 `json:"positiveBal,omitempty"`
	NegativeBal         uint64 `json:"negativeBal,omitempty"`
	Balance             int64  `json:"balance,omitempty"`
	HighestPayout       int64  `json:"highestPayout,omitempty"`
	LowestBalance       int64  `json:"lowestBalance,omitempty"`
	HighestBalance      int64  `json:"highestBalance,omitempty"`
	maxPayout           float64
	bestThreshold       float64
	bestNoFreeThreshold float64
	ss                  *slots.SymbolSet
	AllRounds           *analyse.Rounds                       `json:"allRounds,omitempty"`
	BonusRounds         map[analyse.BonusKind]*analyse.Rounds `json:"bonusRounds,omitempty"`
	FirstPayouts        *Payouts                              `json:"firstPayouts,omitempty"`
	BonusPayouts        map[analyse.BonusKind]*Payouts        `json:"bonusPayouts,omitempty"`
	SpinsTo25x          *analyse.MinMaxUInt64                 `json:"spinsTo25x,omitempty"`
	SpinsTo100x         *analyse.MinMaxUInt64                 `json:"spinsTo100x,omitempty"`
	SpinsTo250x         *analyse.MinMaxUInt64                 `json:"spinsTo250x,omitempty"`
	SpinsTo1000x        *analyse.MinMaxUInt64                 `json:"spinsTo1000x,omitempty"`
	SpinsTo2500x        *analyse.MinMaxUInt64                 `json:"spinsTo2500x,omitempty"`
	SpinsToPlusBal      *analyse.MinMaxUInt64                 `json:"spinsToPlusBal,omitempty"`
	Count25x            *analyse.MinMaxUInt64                 `json:"count25x,omitempty"`
	Count100x           *analyse.MinMaxUInt64                 `json:"count100x,omitempty"`
	Count250x           *analyse.MinMaxUInt64                 `json:"count250x,omitempty"`
	Count1000x          *analyse.MinMaxUInt64                 `json:"count1000x,omitempty"`
	Count2500x          *analyse.MinMaxUInt64                 `json:"count2500x,omitempty"`
	CountPlusBal        *analyse.MinMaxUInt64                 `json:"countPlusBal,omitempty"`
	BonusWheel          *analyse.MinMaxUInt64                 `json:"bonusWheel,omitempty"`
	MultiplierMarks     *analyse.MinMaxUInt64                 `json:"multiplierMarks,omitempty"`
	Multipliers         *analyse.MinMaxFloat64                `json:"multipliers,omitempty"`
	Best                []results.Results                     `json:"best,omitempty"`
	BestNoFree          []results.Results                     `json:"bestNoFree,omitempty"`
	RoundFlags          []*analyse.RoundFlag                  `json:"roundFlags,omitempty"`
	InstantBonus        map[string]uint64                     `json:"instantBonus,omitempty"`
	PlayerChoice        map[string]uint64                     `json:"playerChoice,omitempty"`
	Scripts             map[int]uint64                        `json:"scripts,omitempty"`
	PlayerID            string                                `json:"playerID,omitempty"`
	Symbols             analyse.Symbols                       `json:"symbols,omitempty"`
	Actions             analyse.Actions                       `json:"actions,omitempty"`
	pl                  []*slots.Payline
	pool.Object         `json:"-"`
}

var roundsProvider = pool.NewProducer(func() (pool.Objecter, func()) {
	r := &Rounds{
		bestThreshold:       bestThreshold,
		bestNoFreeThreshold: bestNoFreeThreshold,
		maxBest:             bestMax,
		Symbols:             make(analyse.Symbols, 0, 16),
		Actions:             make(analyse.Actions, 0, 64),
		Best:                make([]results.Results, 0, bestMax+1),
		BestNoFree:          make([]results.Results, 0, bestMax+1),
		SpinsTo25x:          analyse.AcquireMinMaxUInt64(),
		SpinsTo100x:         analyse.AcquireMinMaxUInt64(),
		SpinsTo250x:         analyse.AcquireMinMaxUInt64(),
		SpinsTo1000x:        analyse.AcquireMinMaxUInt64(),
		SpinsTo2500x:        analyse.AcquireMinMaxUInt64(),
		SpinsToPlusBal:      analyse.AcquireMinMaxUInt64(),
		Count25x:            analyse.AcquireMinMaxUInt64(),
		Count100x:           analyse.AcquireMinMaxUInt64(),
		Count250x:           analyse.AcquireMinMaxUInt64(),
		Count1000x:          analyse.AcquireMinMaxUInt64(),
		Count2500x:          analyse.AcquireMinMaxUInt64(),
		CountPlusBal:        analyse.AcquireMinMaxUInt64(),
		BonusWheel:          analyse.AcquireMinMaxUInt64(),
		MultiplierMarks:     analyse.AcquireMinMaxUInt64(),
		Multipliers:         analyse.AcquireMinMaxFloat64(1),
		InstantBonus:        make(map[string]uint64, 8),
		PlayerChoice:        make(map[string]uint64, 8),
		Scripts:             make(map[int]uint64, 32),
		RoundFlags:          make([]*analyse.RoundFlag, 0, 16),
	}
	return r, r.reset
})

// reset clears the rounds metric.
func (r *Rounds) reset() {
	if r != nil {
		r.ResetData()

		r.doubleSpin = false
		r.maxSymbol = 0
		r.reelCount = 0
		r.rowCount = 0
		r.startBalance = 0
		r.LowestBalance = 0
		r.HighestBalance = 0
		r.Balance = 0
		r.maxPayout = 0
		r.PlayerID = ""

		if r.AllRounds != nil {
			r.AllRounds.Release()
			r.AllRounds = nil
		}
		if r.FirstPayouts != nil {
			r.FirstPayouts.Release()
			r.FirstPayouts = nil
		}

		if r.BonusRounds != nil {
			for ix := range r.BonusRounds {
				r.BonusRounds[ix].Release()
				r.BonusRounds[ix] = nil
			}
			clear(r.BonusRounds)
		}
		if r.BonusPayouts != nil {
			for ix := range r.BonusPayouts {
				r.BonusPayouts[ix].Release()
				r.BonusPayouts[ix] = nil
			}
			clear(r.BonusPayouts)
		}

		r.Symbols = analyse.ReleaseSymbols(r.Symbols)
		r.Actions = analyse.ReleaseActions(r.Actions)
		r.ss = nil
		r.pl = nil
	}
}

// ResetData resets the player rounds metrics to initial state.
func (r *Rounds) ResetData() {
	r.RoundCount = 0
	r.WinCount = 0
	r.TotalSpins = 0
	r.RegularSpins = 0
	r.FirstSpins = 0
	r.SecondSpins = 0
	r.RefillSpins = 0
	r.WildRespins = 0
	r.SuperSpins = 0
	r.SuperRefills = 0
	r.FreeTimes = 0
	r.FirstTimes = 0
	r.SecondTimes = 0
	r.FreeTimesSuper = 0
	r.FirstTimesSuper = 0
	r.SecondTimesSuper = 0
	r.FreeAwarded = 0
	r.FirstAwarded = 0
	r.SecondAwarded = 0
	r.FreeAwardedSuper = 0
	r.FirstAwardedSuper = 0
	r.SecondAwardedSuper = 0
	r.FreeSpins = 0
	r.FirstFreeSpins = 0
	r.SecondFreeSpins = 0
	r.SuperSpinsFree = 0
	r.SuperRefillsFree = 0
	r.BadSpins = 0
	r.MaxPayouts = 0
	r.PositiveBal = 0
	r.NegativeBal = 0
	r.HighestPayout = 0
	r.LowestBalance = r.startBalance
	r.HighestBalance = r.startBalance
	r.Balance = r.startBalance

	if r.AllRounds != nil {
		r.AllRounds.ResetData()
	}
	if r.FirstPayouts != nil {
		r.FirstPayouts.ResetData()
	}
	if r.BonusRounds != nil {
		for ix := range r.BonusRounds {
			r.BonusRounds[ix].ResetData()
		}
	}
	if r.BonusPayouts != nil {
		for ix := range r.BonusPayouts {
			r.BonusPayouts[ix].ResetData()
		}
	}

	r.SpinsTo25x.ResetData()
	r.SpinsTo100x.ResetData()
	r.SpinsTo250x.ResetData()
	r.SpinsTo1000x.ResetData()
	r.SpinsTo2500x.ResetData()
	r.SpinsToPlusBal.ResetData()

	r.Count25x.ResetData()
	r.Count100x.ResetData()
	r.Count250x.ResetData()
	r.Count1000x.ResetData()
	r.Count2500x.ResetData()
	r.CountPlusBal.ResetData()

	r.BonusWheel.ResetData()
	r.MultiplierMarks.ResetData()
	r.Multipliers.ResetData()

	clear(r.InstantBonus)
	clear(r.PlayerChoice)
	clear(r.Scripts)

	for id := range r.Symbols {
		if s := r.Symbols[id]; s != nil {
			s.ResetData()
		}
	}

	for id := range r.Actions {
		if a := r.Actions[id]; a != nil {
			a.ResetData()
		}
	}

	for ix := range r.Best {
		results.ReleaseResults(r.Best[ix])
		r.Best[ix] = nil
	}
	r.Best = r.Best[:0]

	for ix := range r.BestNoFree {
		results.ReleaseResults(r.BestNoFree[ix])
		r.BestNoFree[ix] = nil
	}
	r.BestNoFree = r.BestNoFree[:0]

	for ix := range r.RoundFlags {
		r.RoundFlags[ix].Release()
		r.RoundFlags[ix] = nil
	}
	r.RoundFlags = r.RoundFlags[:0]
}

func (r *Rounds) RoundsToCSV(rtp int) {
	// Get the home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}

	// Build the path to the Desktop
	fname := fmt.Sprintf("%s-rounds-%d-%s.csv", r.gameNR.String(), rtp, time.Now().String())
	desktopPath := filepath.Join(homeDir, "Desktop", fname)

	// Create the CSV file on the Desktop
	file, err := os.Create(desktopPath)
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{
		"NoPayouts", "NoBest", "NoSpins", "NoCounts", "GameNR",
		"StartBalance", "RoundCount", "WinCount", "TotalSpins", "RegularSpins",
		"Balance", "FreeSpins", "MaxPayouts", "MaxPayout", "AllRoundsCount",
		"HighestPayout", "BonusRounds", "SpinsTo25x", "SpinsTo100x", "SpinsTo250x",
		"SpinsTo1000x", "SpinsTo2500x", "SpinsToPlusBal", "Count25x", "Count100x",
		"Count250x", "Count1000x", "Count2500x", "CountPlusBal", "Multipliers",
		"LowestBalance", "HighestBalance", "PlayerID",
	}
	writer.Write(header)

	var bonusRounds []string
	for k, v := range r.BonusRounds {
		br := fmt.Sprintf("%s:%d", k.String(), v.Count)
		bonusRounds = append(bonusRounds, br)
	}

	row := []string{
		strconv.FormatBool(r.noPayouts),
		strconv.FormatBool(r.noBest),
		strconv.FormatBool(r.noSpins),
		strconv.FormatBool(r.noCounts),
		r.gameNR.String(),
		strconv.FormatInt(r.startBalance, 10),
		strconv.FormatUint(r.RoundCount, 10),
		strconv.FormatUint(r.WinCount, 10),
		strconv.FormatUint(r.TotalSpins, 10),
		strconv.FormatUint(r.RegularSpins, 10),
		strconv.FormatUint(r.FreeSpins, 10),
		strconv.FormatUint(r.MaxPayouts, 10),
		strconv.FormatFloat(r.maxPayout, 'f', -1, 64),
		strconv.FormatUint(r.AllRounds.Count, 10),
		strconv.FormatInt(r.HighestPayout, 10),
		strings.Join(bonusRounds, ";"),
		strconv.FormatUint(r.SpinsTo25x.Count, 10),
		strconv.FormatUint(r.SpinsTo100x.Count, 10),
		strconv.FormatUint(r.SpinsTo250x.Count, 10),
		strconv.FormatUint(r.SpinsTo1000x.Count, 10),
		strconv.FormatUint(r.SpinsTo2500x.Count, 10),
		strconv.FormatUint(r.SpinsToPlusBal.Count, 10),
		strconv.FormatUint(r.Count25x.Count, 10),
		strconv.FormatUint(r.Count100x.Count, 10),
		strconv.FormatUint(r.Count250x.Count, 10),
		strconv.FormatUint(r.Count1000x.Count, 10),
		strconv.FormatUint(r.Count2500x.Count, 10),
		strconv.FormatUint(r.CountPlusBal.Count, 10),
		strconv.FormatUint(r.MultiplierMarks.Count, 10),
		strconv.FormatInt(r.LowestBalance, 10),
		strconv.FormatInt(r.HighestBalance, 10),
		r.PlayerID,
	}

	err = writer.Write(row)
	if err != nil {
		fmt.Println("Error writing to CSV:", err)
	}

	fmt.Printf("CSV file successfully saved to %s\n", desktopPath)
}

const (
	bestMax             = 25   // max number of "best" rounds.
	bestThreshold       = 2500 // the payout factor threshold to determine "best" rounds.
	bestNoFreeThreshold = 200  // the payout factor threshold to determine "best" rounds with no free spins.
)

package slots

import (
	"math"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/object-pool/pool"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git/consts"
)

// RoundParams contains the parameters to set up a new bet round.
type RoundParams struct {
	Paid         bool
	BuyFeature   uint8
	StartBalance int64 // only if no validator is set up!
	Bet          int64
	TotalBet     int64
	TotalWin     int64 // totalWin and results are mutually exclusive!
	MaxPayout    float64
	GameState    *GameState
	PlayerState  *GamePrefs
	SessionID    string
	RoundID      string
	Results      results.Results // results and totalWin are mutually exclusive!
}

// AcquireRound instantiates a new round from the memory pool.
// It clones the given objects, so it's safe to call Release() on them.
func AcquireRound(v RoundManager, params RoundParams) *Round {
	r := roundsPool.Acquire().(*Round)

	r.validator = v
	r.sessionID = params.SessionID
	r.roundID = params.RoundID
	r.startBalance = params.StartBalance
	r.paid = params.Paid
	r.buyFeature = params.BuyFeature
	r.bet = params.Bet
	r.totalBet = params.TotalBet
	r.totalWin = params.TotalWin
	r.maxPayout = params.MaxPayout

	if params.GameState != nil {
		r.gameState = params.GameState.Clone().(*GameState)
	}
	if params.PlayerState != nil {
		r.gamePrefs = params.PlayerState.Clone().(*GamePrefs)
	}

	if len(params.Results) > 0 {
		for ix := range params.Results {
			res := params.Results[ix]
			r.results = append(r.results, res.Clone().(*results.Result))
			r.roundResults = append(r.roundResults, AcquireRoundResult(ix+1, res))
		}
	}

	return r
}

// Validate calls the D-store API to validate the round with the casino and to store the results in the DB.
// It does nothing if no validator has been set up for the round.
func (r *Round) Validate(debug, reverse bool) *Round {
	if r.validator == nil {
		return r
	}
	r.calculate(reverse)
	r.roundID, r.playerBalance, r.valid = r.validator.PostRound(r, debug)
	return r
}

// ValidateInit calls the D-store API to validate the initial bet of the round with the casino and to store the initial result in the DB.
// It does nothing if no validator has been set up for the round.
func (r *Round) ValidateInit(debug, reverse bool) *Round {
	if r.validator == nil {
		return r
	}
	r.calculate(reverse)
	r.roundID, r.playerBalance, r.valid = r.validator.PostInitRound(r, debug)
	return r
}

// ValidateComplete calls the D-store API to validate the final win of the round with the casino and to store the final results in the DB.
// It does nothing if no validator has been set up for the round.
func (r *Round) ValidateComplete(rs *RoundState, debug, reverse bool) *Round {
	if r.validator == nil {
		return r
	}
	r.calculate(reverse)
	r.roundID, r.playerBalance, r.valid = r.validator.PostCompleteRound(r, rs, debug)
	return r
}

// IsValid returns true if the round was accepted by D-store.
// This will always return false if no validator was set up.
func (r *Round) IsValid() bool {
	return r.valid == nil
}

// Error returns the last error if the round is not valid.
// This will always return nil if no validator was set up.
func (r *Round) Error() error {
	return r.valid
}

// SessionID returns the stored session id.
func (r *Round) SessionID() string {
	return r.sessionID
}

// RoundID returns the stored round id.
func (r *Round) RoundID() string {
	return r.roundID
}

// Bet returns the stored bet.
func (r *Round) Bet() int64 {
	return r.bet
}

// TotalBet returns the stored total bet amount.
func (r *Round) TotalBet() int64 {
	return r.totalBet
}

// MaxPayout returns if the maximum payout was reached.
func (r *Round) MaxPayout() bool {
	return r.maxPayout > 0.0
}

// TotalWin returns the (calculated) win amount for the bet round.
func (r *Round) TotalWin() int64 {
	return r.totalWin
}

// Results returns the stored results.
func (r *Round) Results() results.Results {
	return r.results
}

// RoundResults returns the stored round results.
func (r *Round) RoundResults() RoundResults {
	return r.roundResults
}

// HasSuperShape returns if the round started with a super shape.
func (r *Round) HasSuperShape() bool {
	if len(r.results) == 1 {
		return false
	}
	if res := r.roundResults[0].SpinData; res != nil {
		sticky := res.Sticky()
		if l := len(sticky); l > 0 {
			for ix := 0; ix < l; ix++ {
				if sticky[ix] == 2 {
					return true
				}
			}
		}
	}
	return false
}

// GameState returns the game state.
func (r *Round) GameState() *GameState {
	return r.gameState
}

// SetGameState sets the game state.
func (r *Round) SetGameState(state *GameState) {
	r.gameState = state.Clone().(*GameState)
}

// PlayerState returns the player state.
func (r *Round) PlayerState() *GamePrefs {
	return r.gamePrefs
}

// SetPlayerState sets the player state.
func (r *Round) SetPlayerState(state *GamePrefs) {
	r.gamePrefs = state.Clone().(*GamePrefs)
}

// PlayerBalance returns the new balance for the player if the round was validated by D-store.
// It returns zero if the round was declined, or if the validator is configured not to call D-store.
func (r *Round) PlayerBalance() int64 {
	return r.playerBalance
}

// StartBalance returns the balance given when the round was created. It is usually zero.
func (r *Round) StartBalance() int64 {
	return r.startBalance
}

// NewBalance returns the calculated new balance, based on the initial start balance for the round.
// If the round was set up with a start balance of zero, it should always match the total win.
func (r *Round) NewBalance() int64 {
	return r.newBalance
}

// Balances returns the balance before and after the requested result, based on the initial starting balance for the round.
// It returns zeroes if the index is not a valid result index.
func (r *Round) Balances(index int) (int64, int64) {
	if index < 0 || index >= len(r.roundResults) {
		return 0, 0
	}
	return r.roundResults[index].BalanceBefore, r.roundResults[index].BalanceAfter
}

// ProgressiveWin returns the progressive win for the requested result.
// It returns zero if the index is not a valid result index.
func (r *Round) ProgressiveWin(index int) int64 {
	if index < 0 || index >= len(r.roundResults) {
		return 0
	}
	return r.roundResults[index].ProgressiveWin
}

// PrependFirstSpin prepends the first spin of a double-spin feature to the results.
func (r *Round) PrependFirstSpin(first *RoundResult) {
	r.results = append([]*results.Result{first.Result()}, r.results...)
	r.roundResults = append([]*RoundResult{first}, r.roundResults...)
}

func (r *Round) calculate(reverse bool) {
	switch {
	case r.maxPayout > 0.0:
		r.totalWin = int64(math.Round(float64(r.bet) * r.maxPayout))
	case len(r.results) > 0:
		r.totalWin = int64(math.Round(float64(r.bet) * results.GrandTotal(r.results)))
	}

	before := r.startBalance
	if r.newBalance > 0 && before == 0 {
		before = r.newBalance - r.totalWin
	}

	var progressive, bonus, spin int64
	var bonusActive bool

	for ix := range r.roundResults {
		rr := r.roundResults[ix]
		rr.Bet = r.bet
		rr.TotalWin = r.totalWin
		rr.BalanceBefore = before

		win := int64(math.Round(float64(r.bet) * rr.TotalPayout))

		if !reverse && progressive+win > r.totalWin {
			// progressive cannot exceed the totalWin, so we must've hit maxPayout!
			// adjust the spin win accordingly and record as a maxPayout.
			win = r.totalWin - progressive
			rr.MaxPayout = r.maxPayout
		}

		progressive += win
		if reverse && progressive < 0 {
			// in a reverse win game, we can have penalties that push the total below zero.
			// we need to adjust the actual "win" and reset the progressive total.
			win += progressive
			progressive = 0
		}

		if rr.SpinData != nil {
			switch rr.SpinData.Kind() {
			case slots.RefillSpin:
				spin += win
			case slots.FreeSpin, slots.FirstFreeSpin, slots.SecondFreeSpin:
				spin = win
				bonusActive = true
			default:
				spin = win
			}
		}

		if bonusActive {
			bonus += win
		}

		rr.Win = win
		rr.ProgressiveWin = progressive
		rr.BonusWin = bonus
		rr.SpinWin = spin

		before += win
		rr.BalanceAfter = before
	}

	if r.newBalance == 0 {
		if r.maxPayout > 0.0 {
			r.newBalance = r.startBalance + r.totalWin
		} else {
			r.newBalance = before
			if len(r.roundResults) == 0 {
				r.newBalance += r.totalWin
			}
		}
	}
}

// Round contains all details for a bet round.
// A round should only be kept in memory for the duration of a round, as its details are fleeting.
type Round struct {
	paid          bool            // indicates the results were generated from a bonus buy feature.
	buyFeature    uint8           // indicates the unique id of the bonus buy feature.
	startBalance  int64           // balance at the start of the round.
	bet           int64           // bet amount for the spins in the round (e.g. the stake).
	totalBet      int64           // total bet for the round, based on an optional bonus buy.
	totalWin      int64           // total win for the round.
	newBalance    int64           // balance at the end of the round.
	playerBalance int64           // player balance reported by D-store from a validation.
	maxPayout     float64         // zero or the max payout if it was reached.
	gameState     *GameState      // game state for the round.
	gamePrefs     *GamePrefs      // game preferences for the round.
	sessionID     string          // related session id for the round.
	roundID       string          // unique id for the round.
	valid         error           // indicates if the complete round is valid or not.
	validator     RoundManager    // validation interface.
	results       results.Results // slice of results from the game engine.
	roundResults  RoundResults    // slice of enriched results (with start/end balance, progressive win, etc).
	pool.Object
}

// roundsPool is the memory pool for rounds.
var roundsPool = pool.NewProducer(func() (pool.Objecter, func()) {
	r := &Round{
		valid:        consts.ErrNotValidated,
		results:      make(results.Results, 0, 16),
		roundResults: make(RoundResults, 0, 16),
	}
	return r, r.reset
})

// reset clears the round.
func (r *Round) reset() {
	if r != nil {
		if r.gameState != nil {
			r.gameState.Release()
			r.gameState = nil
		}
		if r.gamePrefs != nil {
			r.gamePrefs.Release()
			r.gamePrefs = nil
		}

		for ix := range r.results {
			r.results[ix].Release()
			r.results[ix] = nil
		}

		for ix := range r.roundResults {
			r.roundResults[ix].Release()
			r.roundResults[ix] = nil
		}

		r.paid = false
		r.buyFeature = 0
		r.startBalance = 0
		r.bet = 0
		r.totalBet = 0
		r.totalWin = 0
		r.newBalance = 0
		r.playerBalance = 0
		r.maxPayout = 0
		r.sessionID = ""
		r.roundID = ""
		r.valid = consts.ErrNotValidated
		r.validator = nil
		r.results = r.results[:0]
		r.roundResults = r.roundResults[:0]
	}
}

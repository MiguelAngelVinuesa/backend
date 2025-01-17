package handlers

import (
	"math"
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	slot "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/games/slots"
	rslt "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
	util "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
	mngr "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git/state/slots"
	log2 "git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/log"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/game"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/log"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/metrics"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/state"
)

func initGame(params *roundParams) *slot.Regular {
	g := game.NewGame(params.gameNR, int(params.rtp))
	if g == nil {
		return nil
	}

	target := g.Slots().RTP()
	if math.Abs(target-float64(params.rtp)) > 0.005 {
		// RTP mismatch; we got a bug so bail out!!!
		g.Release()
		return nil
	}

	if g.IsDoubleSpin() {
		// e.g. ChaCha Bomb
		if params.state == nil {
			started := time.Now()
			params.state = state.GetGameState(params.sessionID)
			metrics.Metrics.AddDuration(metrics.DsSessionGet, started)
		}

		if params.state != nil {
			if params.state.Bet() != params.bet {
				// The bet changed so reload or reset the symbol flags!
				started2 := time.Now()
				_, _, gp, err := state.Manager.GetGamePrefs(params.sessionID)
				metrics.Metrics.AddDuration(metrics.DsGamePrefsGet, started2)

				var flags *slots.SymbolsState
				if err == nil && gp != nil {
					if flags = gp.GetStateCCB(params.bet); flags != nil {
						// we got stored flags for this bet, so use those!
						params.state.SetSymbolsState(flags)
					}
					gp.Release()
				}
				if flags == nil {
					// no stored flags for this bet, so reset!
					if ss := params.state.SymbolsState(); ss != nil {
						ss.ResetState()
					}
				}
			}

			// SUPERVISED-BUILD-REMOVE-START
			if config.DebugMode {
				if l := len(params.flagged); l > 0 && params.state.SymbolsState() != nil {
					for ix := 1; ix < l; ix++ {
						params.state.SymbolsState().SetFlagged(util.Index(ix), params.flagged[ix])
					}
				}
			}
			// SUPERVISED-BUILD-REMOVE-END
		}
	}

	if params.state != nil && (params.state.SpinState() != nil || params.state.SymbolsState() != nil) {
		// SUPERVISED-BUILD-REMOVE-START
		if config.DebugMode && params.state.SpinState() != nil && len(params.initial) > 0 {
			params.state.SpinState().SetIndexes(params.initial)
		}
		// SUPERVISED-BUILD-REMOVE-END
		g.RestoreState(params.state.SpinState(), params.state.SymbolsState())
	}

	return g
}

func playRound(params *roundParams, g *slot.Regular) *mngr.Round {
	started := time.Now()
	haveInitial, haveCache, haveScript := len(params.initial) > 0, len(params.prngCache) > 0, params.scriptID > 0

	var betMultiplier int64 = 1
	if params.paid {
		if paid := g.ForSale(params.bonusKind); paid == nil {
			params.paid = false
		} else {
			betMultiplier = int64(paid.BetMultiplier())
		}
	}

	var results rslt.Results
	// SUPERVISED-BUILD-REMOVE-START
	if config.DebugMode && (haveInitial || haveCache || haveScript) {
		switch {
		case haveScript:
			results = g.Scripted(int(params.scriptID), params.bonusKind, params.resume, params.choices)
		case haveCache:
			results = g.PrngCache(params.prngCache, params.bonusKind, params.resume, params.choices)
		default:
			results = g.Debug(params.initial, params.bonusKind, params.resume, params.choices)
		}
		metrics.Metrics.AddDuration(metrics.GeRoundDebug, started)
	} else {
		// SUPERVISED-BUILD-REMOVE-END
		// else body remains
		if params.resume {
			results = g.RoundResume(params.choices)
			metrics.Metrics.AddDuration(metrics.GeRoundResume, started)
		} else {
			results = g.Round(params.bonusKind)
			metrics.Metrics.AddDuration(metrics.GeRound, started)
		}
		// SUPERVISED-BUILD-REMOVE-START
	}
	// SUPERVISED-BUILD-REMOVE-END

	if len(results) == 0 {
		return nil
	}

	var roundSeq int64
	if params.state != nil {
		roundSeq = params.state.RoundSeq()
		params.state.Release()
		params.state = nil
	}
	if !g.IsDoubleSpin() || g.SpinState() != nil {
		roundSeq++
	}

	params.state = mngr.AcquireGameState(g.SpinState(), g.SymbolsState(), params.bet)
	params.state.SetRoundSeq(roundSeq)
	if params.roundID != "" {
		params.state.SetRoundID(params.roundID)
	}

	params2 := mngr.RoundParams{
		SessionID:  params.sessionID,
		RoundID:    params.roundID,
		Paid:       params.paid,
		BuyFeature: params.bonusKind,
		Bet:        params.bet,
		TotalBet:   params.bet * betMultiplier,
		Results:    results,
		GameState:  params.state,
	}

	if g.MaxPayoutReached() {
		params2.MaxPayout = g.TotalPayout()
	}

	return mngr.AcquireRound(state.Manager, params2)
}

func validateRound(params *roundParams, g *slot.Regular, round *mngr.Round) error {
	if s := round.GameState(); s != nil {
		s.SetNextOffset(0)
		if (params.second || params.resume) && params.state.SpinState() == nil {
			s.SetNextOffset(1)
		}
	}

	var err error
	switch {
	case params.second || params.resume:
		// complete the round
		started := time.Now()
		rs, _ := state.Manager.GetRoundState(params.sessionID, params.roundID)
		metrics.Metrics.AddDuration(metrics.DsRoundState, started)

		started = time.Now()
		err = round.ValidateComplete(rs, params.debug, g.IsReverseWin()).Error()
		metrics.Metrics.AddDuration(metrics.DsRoundComplete, started)

		if log.DsTime && log.Logger.Enabled(log2.DebugLevel) {
			log.Logger.Debug(consts.MsgDsRoundCompleteTrip, consts.FieldSession, params.sessionID, consts.FieldTime, time.Since(started))
		}

	case g.IsDoubleSpin() && !round.HasSuperShape():
		// first spin of a DSF.
		started := time.Now()
		err = round.ValidateInit(params.debug, g.IsReverseWin()).Error()
		metrics.Metrics.AddDuration(metrics.DsRoundInit, started)

		if log.DsTime && log.Logger.Enabled(log2.DebugLevel) {
			log.Logger.Debug(consts.MsgDsRoundInitTrip, consts.FieldSession, params.sessionID, consts.FieldTime, time.Since(started))
		}

	case g.AllowPlayerChoices() && g.NeedPlayerChoice():
		// requires player choice; e.g. we only have the first result(s).
		started := time.Now()
		err = round.ValidateInit(params.debug, g.IsReverseWin()).Error()
		metrics.Metrics.AddDuration(metrics.DsRoundInit, started)

		if log.DsTime && log.Logger.Enabled(log2.DebugLevel) {
			log.Logger.Debug(consts.MsgDsRoundInitTrip, consts.FieldSession, params.sessionID, consts.FieldTime, time.Since(started))
		}

	default:
		// normal spin; e.g. all in one complete round.
		started := time.Now()
		err = round.Validate(params.debug, g.IsReverseWin()).Error()
		metrics.Metrics.AddDuration(metrics.DsRound, started)

		if log.DsTime && log.Logger.Enabled(log2.DebugLevel) {
			log.Logger.Debug(consts.MsgDsRoundTrip, consts.FieldSession, params.sessionID, consts.FieldTime, time.Since(started))
		}
	}

	return err
}

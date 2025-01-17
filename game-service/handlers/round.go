package handlers

import (
	"runtime/debug"
	"time"

	"github.com/gofiber/fiber/v2"

	util "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
	mngr "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git/state/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/tg"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/utils/conv"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/api/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/encode"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/log"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/metrics"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/state"
)

func PostRound(req *fiber.Ctx) (err error) {
	started := time.Now()
	defer func() { metrics.Metrics.AddDuration(metrics.ApiRound, started) }()

	var params *models.RoundStartRequest
	defer func() {
		if e := recover(); e != nil {
			err = sendErrorStack(req, consts.PathRound, e, params, fiber.StatusInternalServerError, BodyInternalError(consts.ErrCdPanic, consts.ErrLvlFatal), debug.Stack())
		}
	}()

	// decode request.
	params = roundStartRequestPool.Get().(*models.RoundStartRequest)
	defer func() {
		params.SessionID = ""
		params.Bet = 0
		params.I18n = nil
		roundStartRequestPool.Put(params)
	}()

	if err = req.BodyParser(&params); err != nil {
		return sendError(req, consts.PathRound, err, params, fiber.StatusBadRequest, BodyBadRequest(consts.ErrCdParams, consts.ErrLvlFatal))
	}

	sess, err2 := tg.VerifySessionID(params.SessionID)
	bet := params.Bet
	if params.SessionID == "" || bet <= 0 || err2 != nil {
		return sendError(req, consts.PathRound, FmtInvalidSession("verification", err2), params, fiber.StatusBadRequest, BodyInvalidSession(consts.ErrCdSessionInvalid, consts.ErrLvlFatal))
	}

	// mark active session.
	metrics.MarkSession(params.SessionID)

	// retrieve game state.
	started2 := time.Now()
	gs := state.GetGameState(params.SessionID)
	metrics.Metrics.AddDuration(metrics.DsSessionGet, started2)

	if sess.DSF() && gs != nil && gs.SpinState() != nil {
		return sendError(req, consts.PathRound, consts.ErrorInvalidStatus, params, fiber.StatusBadRequest, BodyInvalidStatus(consts.ErrCdSpinStateInvalid, consts.ErrLvlFatal))
	}

	// play game round.
	return execRound(req, &roundParams{
		label:     consts.PathRound,
		req:       params,
		i18n:      params.I18n,
		sessionID: params.SessionID,
		gameNR:    sess.GameNr(),
		rtp:       sess.RTP(),
		bet:       bet,
		state:     gs,
	})
}

func PostRoundPaid(req *fiber.Ctx) (err error) {
	started := time.Now()
	defer func() { metrics.Metrics.AddDuration(metrics.ApiRoundPaid, started) }()

	var params *models.RoundPaidRequest
	defer func() {
		if e := recover(); e != nil {
			err = sendErrorStack(req, consts.PathRoundPaid, e, params, fiber.StatusInternalServerError, BodyInternalError(consts.ErrCdPanic, consts.ErrLvlFatal), debug.Stack())
		}
	}()

	// decode request.
	params = roundPaidRequestPool.Get().(*models.RoundPaidRequest)
	defer func() {
		params.SessionID = ""
		params.Feature = 0
		params.Bet = 0
		params.PlayerChoice = nil
		params.I18n = nil
		roundPaidRequestPool.Put(params)
	}()

	if err = req.BodyParser(params); err != nil {
		return sendError(req, consts.PathRoundPaid, err, params, fiber.StatusBadRequest, BodyBadRequest(consts.ErrCdParams, consts.ErrLvlFatal))
	}

	sess, err2 := tg.VerifySessionID(params.SessionID)
	if params.SessionID == "" || err2 != nil {
		return sendError(req, consts.PathRoundPaid, FmtInvalidSession("verification", err2), params, fiber.StatusBadRequest, BodyInvalidSession(consts.ErrCdSessionInvalid, consts.ErrLvlFatal))
	}

	// mark active session.
	metrics.MarkSession(params.SessionID)

	// decode player choices.
	var choices map[string]string
	if params.PlayerChoice != nil {
		m, _ := params.PlayerChoice.(map[string]any)
		choices = make(map[string]string, len(m))
		for k, v := range m {
			choices[k] = util.StringFromAny(v)
		}
	}

	// play game round.
	return execRound(req, &roundParams{
		paid:      true,
		bonusKind: uint8(params.Feature),
		label:     consts.PathRoundPaid,
		req:       params,
		i18n:      params.I18n,
		sessionID: params.SessionID,
		choices:   choices,
		bet:       params.Bet,
		gameNR:    sess.GameNr(),
		rtp:       sess.RTP(),
	})
}

func PostRoundSecond(req *fiber.Ctx) (err error) {
	started := time.Now()
	defer func() { metrics.Metrics.AddDuration(metrics.ApiRoundSecond, started) }()

	var params *models.RoundSecondRequest
	defer func() {
		if e := recover(); e != nil {
			err = sendErrorStack(req, consts.PathRoundSecond, e, params, fiber.StatusInternalServerError, BodyInternalError(consts.ErrCdPanic, consts.ErrLvlFatal), debug.Stack())
		}
	}()

	// decode request.
	params = roundSecondRequestPool.Get().(*models.RoundSecondRequest)
	defer func() {
		params.SessionID = ""
		params.RoundID = ""
		params.PlayerChoice = nil
		params.I18n = nil
		roundSecondRequestPool.Put(params)
	}()

	if err = req.BodyParser(params); err != nil {
		return sendError(req, consts.PathRoundSecond, err, params, fiber.StatusBadRequest, BodyBadRequest(consts.ErrCdParams, consts.ErrLvlFatal))
	}

	sess, err2 := tg.VerifySessionID(params.SessionID)
	if params.SessionID == "" || err2 != nil || !sess.DSF() {
		return sendError(req, consts.PathRoundSecond, FmtInvalidSession("verification", err2), params, fiber.StatusBadRequest, BodyInvalidSession(consts.ErrCdSessionInvalid, consts.ErrLvlFatal))
	}

	// mark active session.
	metrics.MarkSession(params.SessionID)

	// retrieve game state.
	started2 := time.Now()
	gs := state.GetGameState(params.SessionID)
	metrics.Metrics.AddDuration(metrics.DsSessionGet, started2)

	if gs == nil || gs.Bet() <= 0 || gs.SpinState() == nil || params.RoundID == "" {
		return sendError(req, consts.PathRoundSecond, consts.ErrorInvalidStatus, params, fiber.StatusBadRequest, BodyInvalidStatus(consts.ErrCdSpinStateInvalid, consts.ErrLvlFatal))
	}

	// decode player choices.
	var choices map[string]string
	if params.PlayerChoice != nil {
		m, _ := params.PlayerChoice.(map[string]any)
		choices = make(map[string]string, len(m))
		for k, v := range m {
			choices[k] = util.StringFromAny(v)
		}

		// test for player choice of the sticky symbol.
		if symbol := conv.IntFromAny(choices["stickySymbol"]); symbol > 0 && symbol < 100 {
			gs.SpinState().SetStickySymbol(util.Index(symbol))
		}
	}

	// play game round.
	return execRound(req, &roundParams{
		second:    true,
		label:     consts.PathRoundSecond,
		req:       params,
		i18n:      params.I18n,
		sessionID: params.SessionID,
		roundID:   params.RoundID,
		choices:   choices,
		bet:       gs.Bet(),
		gameNR:    sess.GameNr(),
		rtp:       sess.RTP(),
		state:     gs,
	})
}

func PostRoundResume(req *fiber.Ctx) (err error) {
	started := time.Now()
	defer func() { metrics.Metrics.AddDuration(metrics.ApiRoundResume, started) }()

	var params *models.RoundResumeRequest
	defer func() {
		if e := recover(); e != nil {
			err = sendErrorStack(req, consts.PathRoundResume, e, params, fiber.StatusInternalServerError, BodyInternalError(consts.ErrCdPanic, consts.ErrLvlFatal), debug.Stack())
		}
	}()

	// decode request.
	params = roundResumeRequestPool.Get().(*models.RoundResumeRequest)
	defer func() {
		params.SessionID = ""
		params.RoundID = ""
		params.PlayerChoice = nil
		params.I18n = nil
		roundResumeRequestPool.Put(params)
	}()

	if err = req.BodyParser(params); err != nil {
		return sendError(req, consts.PathRoundResume, err, params, fiber.StatusBadRequest, BodyBadRequest(consts.ErrCdParams, consts.ErrLvlFatal))
	}

	sess, err2 := tg.VerifySessionID(params.SessionID)
	if params.SessionID == "" || err2 != nil {
		return sendError(req, consts.PathRoundResume, FmtInvalidSession("verification", err2), params, fiber.StatusBadRequest, BodyInvalidSession(consts.ErrCdSessionInvalid, consts.ErrLvlFatal))
	}

	// mark active session.
	metrics.MarkSession(params.SessionID)

	// retrieve game state.
	started2 := time.Now()
	gs := state.GetGameState(params.SessionID)
	metrics.Metrics.AddDuration(metrics.DsSessionGet, started2)

	if gs == nil || gs.Bet() <= 0 || params.RoundID == "" {
		return sendError(req, consts.PathRoundResume, consts.ErrorInvalidStatus, params, fiber.StatusBadRequest, BodyInvalidStatus(consts.ErrCdSpinStateInvalid, consts.ErrLvlFatal))
	}

	// decode player choices.
	var choices map[string]string
	if params.PlayerChoice != nil {
		m, _ := params.PlayerChoice.(map[string]any)
		choices = make(map[string]string, len(m))
		for k, v := range m {
			choices[k] = util.StringFromAny(v)
		}
	}

	// play game round.
	return execRound(req, &roundParams{
		resume:    true,
		label:     consts.PathRoundResume,
		req:       params,
		i18n:      params.I18n,
		sessionID: params.SessionID,
		roundID:   params.RoundID,
		choices:   choices,
		bet:       gs.Bet(),
		gameNR:    sess.GameNr(),
		rtp:       sess.RTP(),
		state:     gs,
	})
}

type roundParams struct {
	second    bool
	resume    bool
	debug     bool
	paid      bool
	rtp       uint8
	bonusKind uint8
	gameNR    tg.GameNR
	scriptID  int32
	bet       int64
	req       any
	i18n      *models.PrefetchI18n
	state     *mngr.GameState
	prefs     *mngr.GamePrefs
	initial   util.Indexes
	prngCache []int
	flagged   []bool
	choices   map[string]string
	label     string
	sessionID string
	roundID   string
}

func execRound(req *fiber.Ctx, params *roundParams) error {
	defer func() {
		if params.state != nil {
			params.state.Release()
		}
	}()

	// init the game.
	g := initGame(params)
	if g == nil {
		return sendError(req, params.label, FmtInvalidSession("game", nil), params.req, fiber.StatusBadRequest, BodyInvalidSession(consts.ErrCdGameInvalid, consts.ErrLvlFatal))
	}
	defer g.Release()

	// play and validate a game round.
	round := playRound(params, g)
	if round == nil {
		return sendError(req, params.label, consts.ErrorInternalError, params.req, fiber.StatusInternalServerError, BodyInternalError(consts.ErrCdGameFailed, consts.ErrLvlRetry))
	}
	defer round.Release()

	if err := validateRound(params, g, round); err != nil {
		status := fiber.StatusBadRequest
		return sendError(req, params.label, FmtInvalidSession("validation", err), params.req, status, BodyDstoreError(err))
	}

	// For double-spin feature we need to remember the roundID for the second spin!
	if g.IsDoubleSpin() && params.state.SpinState() != nil {
		saveRoundID(params, round.RoundID())
	}

	// For CCB we need to update the symbol flags in the game prefs.
	if params.gameNR == tg.CCBnr {
		fixCCBflags(params)
	}

	// generate & send the response.
	return sendResponse(params.label, req, params.req, encode.BuildRoundResponse(g, round, params.i18n))
}

func saveRoundID(params *roundParams, id string) {
	params.state.SetRoundID(id)

	started := time.Now()
	err := state.Manager.PutGameState(params.sessionID, params.state)
	metrics.Metrics.AddDuration(metrics.DsSessionPut, started)

	if err != nil {
		log.Logger.Error(FmtDstoreError(err))
	}
}

func fixCCBflags(params *roundParams) {
	flags := params.state.SymbolsState()
	if flags == nil {
		return
	}

	if params.prefs == nil {
		_, _, params.prefs, _ = state.Manager.GetGamePrefs(params.sessionID)
	}

	var modified bool

	if params.prefs == nil {
		params.prefs = mngr.AcquireGamePrefs(params.bet, flags)
		modified = true
	} else {
		s := params.prefs.GetStateCCB(params.bet)
		if s == nil {
			modified = true
		} else {
			f1, f2 := s.Flagged(), flags.Flagged()
			if len(f1) != len(f2) {
				modified = true
			} else {
				for ix := range f1 {
					if f1[ix] != f2[ix] {
						modified = true
						break
					}
				}
			}
		}
		if modified {
			params.prefs.AddStateCCB(params.bet, flags)
		}
	}

	if modified {
		started2 := time.Now()
		err := state.Manager.PutGamePrefs(params.sessionID, params.prefs)
		metrics.Metrics.AddDuration(metrics.DsGamePrefsPut, started2)

		if err != nil {
			log.Logger.Error(FmtDstoreError(err))
		}
	}
}

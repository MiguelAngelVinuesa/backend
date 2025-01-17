//go:build DEBUG

// SUPERVISED-BUILD-REMOVE-START
// REMOVE ENTIRE FILE
// SUPERVISED-BUILD-REMOVE-END

package handlers

import (
	"runtime/debug"
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/utils/conv"
	"github.com/gofiber/fiber/v2"

	util "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
	mngr "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git/state/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/tg"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/api/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/metrics"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/state"
)

func PostRoundDebug(req *fiber.Ctx) (err error) {
	started := time.Now()
	defer func() { metrics.Metrics.AddDuration(metrics.ApiRoundDebug, started) }()

	return doRoundDebug(req, consts.PathRoundDebug, false, false)
}

func PostRoundDebugSecond(req *fiber.Ctx) (err error) {
	started := time.Now()
	defer func() { metrics.Metrics.AddDuration(metrics.ApiRoundDebugSecond, started) }()

	return doRoundDebug(req, consts.PathRoundDebugSecond, true, false)
}

func PostRoundDebugResume(req *fiber.Ctx) (err error) {
	started := time.Now()
	defer func() { metrics.Metrics.AddDuration(metrics.ApiRoundDebugResume, started) }()

	return doRoundDebug(req, consts.PathRoundDebugResume, false, true)
}

func doRoundDebug(req *fiber.Ctx, label string, second, resume bool) (err error) {
	var params *models.RoundDebugRequest
	defer func() {
		if e := recover(); e != nil {
			err = sendErrorStack(req, label, e, params, fiber.StatusInternalServerError, BodyInternalError(100, "F"), debug.Stack())
		}
	}()
	// decode request.
	params = roundDebugRequestPool.Get().(*models.RoundDebugRequest)
	defer func() {
		params.SessionID = ""
		params.RoundID = ""
		params.Bet = 0
		params.Initial = nil
		params.Rngmagic = nil
		params.State = nil
		params.ScriptID = 0
		params.PlayerChoice = nil
		params.I18n = nil
		roundDebugRequestPool.Put(params)
	}()

	if err = req.BodyParser(params); err != nil {
		return sendError(req, label, err, params, fiber.StatusBadRequest, BodyBadRequest(consts.ErrCdPanic, consts.ErrLvlFatal))
	}

	sess, err2 := tg.VerifySessionID(params.SessionID)

	initial := make(util.Indexes, len(params.Initial))
	for ix, v := range params.Initial {
		initial[ix] = util.Index(v)
	}

	prngCache := make([]int, len(params.Rngmagic))
	for ix, v := range params.Rngmagic {
		prngCache[ix] = int(v)
	}

	if params.SessionID == "" || params.Bet < 0 || err2 != nil || (len(initial) == 0 && len(prngCache) == 0 && params.ScriptID == 0) {
		return sendError(req, label, FmtInvalidSession("invalid parameters", nil), params, fiber.StatusBadRequest, BodyInvalidSession(consts.ErrCdParams, consts.ErrLvlFatal))
	}

	// mark active session.
	metrics.MarkSession(params.SessionID)

	// retrieve game state.
	var gs *mngr.GameState
	if second || resume || params.RoundID == "" || params.Bet == 0 {
		started2 := time.Now()
		gs = state.GetGameState(params.SessionID)
		metrics.Metrics.AddDuration(metrics.DsSessionGet, started2)

		if gs != nil {
			if params.RoundID == "" {
				params.RoundID = gs.RoundID()
			}
			if params.Bet == 0 {
				params.Bet = gs.Bet()
			}
		}
	}

	if resume && (gs == nil || gs.Bet() <= 0 || params.RoundID == "") {
		return sendError(req, label, consts.ErrorInvalidStatus, params, fiber.StatusBadRequest, BodyInvalidStatus(consts.ErrCdSpinStateInvalid, consts.ErrLvlFatal))
	}

	// decode initial flags.
	var flagged []bool
	if m, ok := params.State.(map[string]any); ok {
		if f, ok2 := m["flagged"].([]any); ok2 {
			flagged = make([]bool, len(f))
			for ix := range flagged {
				if i, ok3 := f[ix].(float64); ok3 {
					flagged[ix] = int(i) == 1
				}
			}
		}
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
		if gs.SpinState() != nil {
			if symbol := conv.IntFromAny(choices["stickySymbol"]); symbol > 0 && symbol < 100 {
				gs.SpinState().SetStickySymbol(util.Index(symbol))
			}
		}
	}

	// play game round.
	return execRound(req, &roundParams{
		second:    second,
		resume:    resume,
		gameNR:    sess.GameNr(),
		debug:     true,
		label:     label,
		req:       params,
		i18n:      params.I18n,
		sessionID: params.SessionID,
		roundID:   params.RoundID,
		bet:       params.Bet,
		rtp:       sess.RTP(),
		initial:   initial,
		prngCache: prngCache,
		scriptID:  params.ScriptID,
		flagged:   flagged,
		state:     gs,
		choices:   choices,
	})
}

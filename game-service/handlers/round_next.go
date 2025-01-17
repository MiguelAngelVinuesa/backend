package handlers

import (
	"runtime/debug"
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/tg"
	"github.com/gofiber/fiber/v2"

	mngr "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git/state/slots"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/api/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/encode"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/game"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/metrics"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/state"
)

func PostRoundNext(req *fiber.Ctx) (err error) {
	started := time.Now()
	defer func() { metrics.Metrics.AddDuration(metrics.ApiRoundNext, started) }()

	var params *models.RoundNextRequest
	defer func() {
		if e := recover(); e != nil {
			err = sendErrorStack(req, consts.PathRoundNext, e, params, fiber.StatusInternalServerError, BodyInternalError(consts.ErrCdPanic, consts.ErrLvlFatal), debug.Stack())
		}
	}()

	// decode request.
	params = roundNextRequestPool.Get().(*models.RoundNextRequest)
	defer func() {
		params.SessionID = ""
		params.RoundID = ""
		params.SpinSeq = 0
		params.I18n = nil
		roundNextRequestPool.Put(params)
	}()

	if err = req.BodyParser(params); err != nil {
		return sendError(req, consts.PathRoundNext, err, params, fiber.StatusBadRequest, BodyBadRequest(consts.ErrCdParams, consts.ErrLvlFatal))
	}

	if params.SessionID == "" || params.RoundID == "" || params.SpinSeq <= 0 {
		return sendError(req, consts.PathRoundNext, FmtInvalidSession("parameters", nil), params, fiber.StatusBadRequest, BodyInvalidSession(consts.ErrCdParams, consts.ErrLvlFatal))
	}

	sess, err2 := tg.VerifySessionID(params.SessionID)
	if err2 != nil {
		return sendError(req, consts.PathRoundNext, FmtInvalidSession("verification", err2), params, fiber.StatusBadRequest, BodyInvalidSession(consts.ErrCdParams, consts.ErrLvlFatal))
	}

	g := game.NewGame(sess.GameNr(), int(sess.RTP()))
	if g == nil {
		return sendError(req, consts.PathRoundNext, FmtInvalidSession("game", nil), params, fiber.StatusBadRequest, BodyInvalidSession(consts.ErrCdGameInvalid, consts.ErrLvlFatal))
	}
	defer g.Release()

	var r *mngr.RoundState
	if !sess.SharedRound() {
		// mark active session.
		metrics.MarkSession(params.SessionID)
	}

	// retrieve round state.
	var r2 *mngr.RoundState
	started2 := time.Now()
	r, err = state.Manager.GetRoundState(params.SessionID, params.RoundID)
	metrics.Metrics.AddDuration(metrics.DsRoundState, started2)
	if err == nil && r != nil && !sess.SharedRound() && r.Spins() > 0 {
		r2 = r
		if seq := int(params.SpinSeq) - 1; seq > 0 {
			now := time.Now().UTC()
			if r2.PlayedFull() {
				r2.SpinReplayed(seq, now)
			} else {
				r2.SpinPlayed(seq, now)
			}
		}
	}

	started2 = time.Now()
	roundResult, endBalance, err3 := state.Manager.PostRoundNext(params.SessionID, params.RoundID, r2, int(params.SpinSeq))
	metrics.Metrics.AddDuration(metrics.DsRoundNext, started2)

	if roundResult != nil {
		defer roundResult.Release()
	}
	if err3 != nil {
		return sendError(req, consts.PathRoundNext, FmtDstoreError(err3), params, fiber.StatusInternalServerError, BodyDstoreError(err3))
	}
	if roundResult == nil {
		return sendError(req, consts.PathRoundNext, consts.ErrorNotFound, params, fiber.StatusNotFound, BodyNotFound(consts.ErrCdNotFound, consts.ErrLvlRetry))
	}

	// generate & send response.
	var c int
	var id string
	if params.SpinSeq == 1 && r != nil {
		c = r.Spins()
		id = params.RoundID
	}

	return sendResponse(consts.PathRoundNext, req, params, encode.BuildRoundNextResponse(id, endBalance, g, c, roundResult, params.I18n))
}

package handlers

import (
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gofiber/fiber/v2"

	log2 "git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/log"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/api/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/log"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/metrics"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/state"
)

func PostRoundFinish(req *fiber.Ctx) (err error) {
	started := time.Now()
	defer func() { metrics.Metrics.AddDuration(metrics.ApiRoundFinish, started) }()

	var params *models.RoundFinishRequest
	defer func() {
		if e := recover(); e != nil {
			err = sendErrorStack(req, consts.PathRoundFinish, e, params, http.StatusInternalServerError, BodyInternalError(consts.ErrCdPanic, consts.ErrLvlFatal), debug.Stack())
		}
	}()

	// decode request.
	params = roundFinishRequestPool.Get().(*models.RoundFinishRequest)
	defer func() {
		params.SessionID = ""
		params.RoundID = ""
		roundFinishRequestPool.Put(params)
	}()

	if err = req.BodyParser(params); err != nil {
		return sendError(req, consts.PathRoundFinish, err, params, http.StatusBadRequest, BodyBadRequest(consts.ErrCdParams, consts.ErrLvlFatal))
	}

	if params.SessionID == "" || params.RoundID == "" {
		return sendError(req, consts.PathRoundFinish, FmtInvalidSession("invalid parameters", nil), params, http.StatusBadRequest, BodyInvalidSession(consts.ErrCdParams, consts.ErrLvlFatal))
	}

	// mark active session.
	metrics.MarkSession(params.SessionID)

	// retrieve & update round state.
	started2 := time.Now()
	r, err2 := state.Manager.GetRoundState(params.SessionID, params.RoundID)
	metrics.Metrics.AddDuration(metrics.DsRoundState, started2)
	if err2 == nil && r != nil && r.Spins() > 0 {
		seq := r.Spins()
		now := time.Now().UTC()
		if r.PlayedFull() {
			r.SpinReplayed(seq, now)
		} else {
			r.SpinPlayed(seq, now)
		}
	}

	// post round/next; e.g. forced call with spinSeq==1 to update the round state.
	started2 = time.Now()
	_, _, err2 = state.Manager.PostRoundNext(params.SessionID, params.RoundID, r, 1)
	metrics.Metrics.AddDuration(metrics.DsRoundNext, started2)
	if err2 != nil {
		return sendError(req, consts.PathRoundFinish, FmtDstoreError(err), params, http.StatusNotFound, BodyDstoreError(err))
	}

	// send response.
	if log.API && log.Logger.Enabled(log2.DebugLevel) {
		log.Logger.Debug(consts.PathRoundFinish, consts.FieldRequest, params, consts.FieldResponse, string(consts.SuccessResponse))
	}

	req.Set(consts.ContentType, consts.ApplicationJSON)
	_, err = req.Write(consts.SuccessResponse)
	return err
}

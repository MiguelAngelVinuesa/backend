package handlers

import (
	"runtime/debug"
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/utils/conv"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/clients/bo_backend"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/metrics"
)

func GetSessionInfo(req *fiber.Ctx) (err error) {
	started := time.Now()
	defer func() { metrics.Metrics.AddDuration(metrics.ApiSessionInfo, started) }()

	params := struct {
		SessionID   string `json:"sessionID,omitempty"`
		CreatedFrom string `json:"createFrom,omitempty"`
		CreatedTo   string `json:"createdTo,omitempty"`
		Limit       int    `json:"limit,omitempty"`
		// hidden
		createdFrom time.Time
		createdTo   time.Time
	}{
		SessionID:   req.Params(consts.FieldSession),
		CreatedFrom: req.Query(consts.FieldCreatedFrom),
		CreatedTo:   req.Query(consts.FieldCreatedTo),
		Limit:       req.QueryInt(consts.FieldLimit),
	}

	defer func() {
		if e := recover(); e != nil {
			err = sendErrorStack(req, consts.PathSessionInfo, e, params, fiber.StatusInternalServerError, BodyInternalError(consts.ErrCdPanic, consts.ErrLvlFatal), debug.Stack())
		}
	}()

	params.createdFrom = conv.TimestampFromAny(params.CreatedFrom)
	params.createdTo = conv.TimestampFromAny(params.CreatedTo)

	// load session info.
	resp, err2 := bo_backend.GetSessionInfo(params.SessionID, params.createdFrom, params.createdTo, params.Limit)
	if err2 != nil {
		return sendError(req, consts.PathSessionInfo, err2, params, fiber.StatusBadRequest, BodyBadRequest(consts.ErrCdParams, consts.ErrLvlFatal))
	}

	// encode and send the response.
	b, _ := json.Marshal(resp)
	req.Set(consts.ContentType, consts.ApplicationJSON)
	_, err = req.Write(b)
	return err
}

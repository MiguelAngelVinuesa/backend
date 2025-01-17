package handlers

// SUPERVISED-BUILD-REMOVE-START
// REMOVE ENTIRE FILE
// SUPERVISED-BUILD-REMOVE-END

import (
	"fmt"
	"runtime/debug"
	"sort"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/tg"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/api/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/game"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/metrics"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/rng_magic"
)

func GetRngConditionsLU(req *fiber.Ctx) (err error) {
	started := time.Now()
	defer func() { metrics.Metrics.AddDuration(metrics.ApiRngConditions, started) }()

	gameID := req.Params("game")
	defer func() {
		if e := recover(); e != nil {
			err = sendErrorStack(req, consts.PathRngConditionsLU, e, gameID, fiber.StatusInternalServerError, BodyInternalError(consts.ErrCdPanic, consts.ErrLvlFatal), debug.Stack())
		}
	}()

	gameNR, err2 := tg.VerifyGameID(gameID)
	if err2 != nil {
		return sendError(req, consts.PathRngConditionsLU, err2, gameID, fiber.StatusBadRequest, BodyBadRequest(consts.ErrCdParams, consts.ErrLvlFatal))
	}

	list, _ := rng_magic.GameData(gameNR)
	if len(list) == 0 {
		return sendError(req, consts.PathRngConditionsLU, fmt.Errorf("no conditions found"), gameID, fiber.StatusBadRequest, BodyBadRequest(consts.ErrCdParams, consts.ErrLvlFatal))
	}

	resp := struct {
		Success bool
		List    []string
	}{
		Success: true,
	}

	for _, c := range list {
		f := c.Name + "("

		for ix, p := range c.Parameters {
			if ix > 0 {
				f += ","
			}
			f += p.Name
		}

		f += ")"

		resp.List = append(resp.List, f)
	}

	sort.Strings(resp.List)

	b, _ := json.Marshal(resp)
	req.Set(consts.ContentType, consts.ApplicationJSON)
	_, err = req.Write(b)
	return err
}

func PostRngMagicTest(req *fiber.Ctx) (err error) {
	started := time.Now()
	defer func() { metrics.Metrics.AddDuration(metrics.ApiRngMagicTest, started) }()

	params := struct {
		Game      string `json:"gameID"`
		Statement string `json:"statement"`
	}{}

	defer func() {
		if e := recover(); e != nil {
			err = sendErrorStack(req, consts.PathRngMagicTest, e, params, fiber.StatusInternalServerError, BodyInternalError(consts.ErrCdPanic, consts.ErrLvlFatal), debug.Stack())
		}
	}()

	if err = req.BodyParser(&params); err != nil {
		return sendError(req, consts.PathRngMagicTest, err, params, fiber.StatusBadRequest, BodyBadRequest(consts.ErrCdParams, consts.ErrLvlFatal))
	}

	gameNR, err2 := tg.VerifyGameID(params.Game)
	if err2 != nil || params.Statement == "" {
		return sendError(req, consts.PathRngMagicTest, err2, params, fiber.StatusBadRequest, BodyBadRequest(consts.ErrCdParams, consts.ErrLvlFatal))
	}

	// test the statement and send outcome.
	if err = rng_magic.TestRngStatement(gameNR, params.Statement); err != nil {
		j, _ := json.Marshal(&models.ErrorResponse{Message: err.Error(), ErrorCode: int64(consts.ErrCdRngFunctionInvalid), ErrorLevel: consts.ErrLvlFatal})
		return sendError(req, consts.PathRngMagicTest, err, params, fiber.StatusBadRequest, j)
	}

	req.Set(consts.ContentType, consts.ApplicationJSON)
	_, err = req.Write(consts.SuccessResponse)
	return err
}

func PostRngMagic(req *fiber.Ctx) (err error) {
	started := time.Now()
	defer func() { metrics.Metrics.AddDuration(metrics.ApiRngMagic, started) }()

	params := struct {
		Game      string `json:"gameID"`
		RTP       int    `json:"rtp"`
		Statement string `json:"statement"`
		Count     int    `json:"count,omitempty"`
		Timeout   int    `json:"timeout,omitempty"`
	}{}

	defer func() {
		if e := recover(); e != nil {
			err = sendErrorStack(req, consts.PathRngMagicTest, e, params, fiber.StatusInternalServerError, BodyInternalError(consts.ErrCdPanic, consts.ErrLvlFatal), debug.Stack())
		}
	}()

	if err = req.BodyParser(&params); err != nil {
		return sendError(req, consts.PathRngMagicTest, err, params, fiber.StatusBadRequest, BodyBadRequest(consts.ErrCdParams, consts.ErrLvlFatal))
	}

	gameNR, err2 := tg.VerifyGameID(params.Game)
	if err2 != nil || params.Statement == "" {
		return sendError(req, consts.PathRngMagicTest, err2, params, fiber.StatusBadRequest, BodyBadRequest(consts.ErrCdParams, consts.ErrLvlFatal))
	}

	g := game.NewGame(gameNR, params.RTP)
	if g == nil {
		return sendError(req, consts.PathRngMagicTest, err2, params, fiber.StatusBadRequest, BodyBadRequest(consts.ErrCdParams, consts.ErrLvlFatal))
	}

	switch {
	case params.Count <= 0:
		params.Count = 1
	case params.Count > 10:
		params.Count = 10
	}

	switch {
	case params.Timeout <= 0:
		params.Timeout = 10
	case params.Timeout > 60:
		params.Timeout = 60
	}

	rngs, err3 := rng_magic.RunRngStatement(g, gameNR, params.Statement, params.Count, params.Timeout)
	if err3 != nil {
		j, _ := json.Marshal(&models.ErrorResponse{Message: err3.Error(), ErrorCode: int64(consts.ErrCdRngFunctionInvalid), ErrorLevel: consts.ErrLvlFatal})
		return sendError(req, consts.PathRngMagicTest, err3, params, fiber.StatusBadRequest, j)
	}

	resp := rng_magic.RngMagics{
		Success: true,
		Magics:  rngs,
	}
	if len(rngs) < params.Count {
		resp.Message = "timeout exceeded; not all results found in time"
		resp.Success = len(rngs) > 0
	}

	b, _ := json.Marshal(resp)
	req.Set(consts.ContentType, consts.ApplicationJSON)
	_, err = req.Write(b)
	return err
}

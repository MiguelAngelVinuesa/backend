package handlers

import (
	"bytes"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/encode"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/hashes"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/metrics"
)

func GetBinHashes(req *fiber.Ctx) error {
	started := time.Now()
	defer func() { metrics.Metrics.AddDuration(metrics.ApiBinHashes, started) }()

	// check api key.
	if !bytes.Equal(req.Request().Header.Peek("X-API-KEY"), config.ApiKeyBytes) {
		return sendError(req, consts.PathBinHashes, consts.ErrorInvalidApiKey, nil, http.StatusBadRequest, BodyBadRequest(consts.ErrCdApiKey, consts.ErrLvlFatal))
	}

	// generate & send response.
	return sendResponse(consts.PathBinHashes, req, nil, encode.Hashes())
}

func GameHash(req *fiber.Ctx) error {
	started := time.Now()
	defer func() { metrics.Metrics.AddDuration(metrics.ApiGameHash, started) }()

	code := req.Params("game")
	hash := hashes.GameHashes[code]
	if hash == "" {
		return sendError(req, consts.PathGameHash, consts.ErrorBadRequest, code, http.StatusBadRequest, BodyBadRequest(consts.ErrCdParams, consts.ErrLvlFatal))
	}

	req.Set(consts.ContentType, consts.PlainText)
	return req.Send([]byte(hash))
}

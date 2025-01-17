package handlers

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	log2 "git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/log"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/clients/i18n"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/log"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/metrics"
)

func PostStrings(req *fiber.Ctx) error {
	started := time.Now()
	defer func() { metrics.Metrics.AddDuration(metrics.ApiStrings, started) }()

	b, err := i18n.PostStrings(req.Body())
	if err != nil {
		return sendError(req, consts.PathStrings, err, string(req.Body()), http.StatusBadRequest, BodyBadRequest(consts.ErrCdNotFound, consts.ErrLvlFatal))
	}
	if log.API && log.Logger.Enabled(log2.DebugLevel) {
		log.Logger.Debug(consts.PathStrings, consts.FieldResponse, string(b))
	}

	req.Set(consts.ContentType, consts.ApplicationJSON)
	_, err = req.Write(b)
	return err
}

func GetPluralString(req *fiber.Ctx) error {
	started := time.Now()
	defer func() { metrics.Metrics.AddDuration(metrics.ApiPlural, started) }()

	loc := req.Params("loc")
	key := req.Params("key")
	if loc == "" || key == "" {
		return sendError(req, consts.PathPlural, nil, "", http.StatusBadRequest, BodyBadRequest(consts.ErrCdParams, consts.ErrLvlFatal))
	}

	b, ct, err := i18n.GetPluralString(loc, key, req.Request().URI().QueryArgs())
	if err != nil {
		return sendError(req, consts.PathPlural, err, string(req.Body()), http.StatusBadRequest, BodyBadRequest(consts.ErrCdNotFound, consts.ErrLvlFatal))
	}

	if log.API && log.Logger.Enabled(log2.DebugLevel) {
		log.Logger.Debug(consts.PathPlural, consts.FieldResponse, string(b))
	}

	req.Set(consts.ContentType, ct)
	_, err = req.Write(b)
	return err
}

func PostPluralStrings(req *fiber.Ctx) error {
	started := time.Now()
	defer func() { metrics.Metrics.AddDuration(metrics.ApiPlurals, started) }()

	b, err := i18n.PostPluralStrings(req.Body())
	if err != nil {
		return sendError(req, consts.PathPlurals, err, string(req.Body()), http.StatusBadRequest, BodyBadRequest(consts.ErrCdNotFound, consts.ErrLvlFatal))
	}

	if log.API && log.Logger.Enabled(log2.DebugLevel) {
		log.Logger.Debug(consts.PathPlurals, consts.FieldResponse, string(b))
	}

	req.Set(consts.ContentType, consts.ApplicationJSON)
	_, err = req.Write(b)
	return err
}

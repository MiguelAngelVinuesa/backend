//go:build !DEBUG

// SUPERVISED-BUILD-REMOVE-START
// remove entire file
// SUPERVISED-BUILD-REMOVE-END

package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts"
)

func RoundDebug(ctx *fiber.Ctx) error {
	return sendError(ctx, consts.PathRoundDebug, consts.ErrorNotFound, nil, http.StatusNotFound, BodyNotFound(consts.ErrCdParams, consts.ErrLvlFatal))
}

func RoundDebugSecond(ctx *fiber.Ctx) error {
	return sendError(ctx, consts.PathRoundDebugSecond, consts.ErrorNotFound, nil, http.StatusNotFound, BodyNotFound(consts.ErrCdParams, consts.ErrLvlFatal))
}

func RoundDebugResume(ctx *fiber.Ctx) error {
	return sendError(ctx, consts.PathRoundDebugResume, consts.ErrorNotFound, nil, http.StatusNotFound, BodyNotFound(consts.ErrCdParams, consts.ErrLvlFatal))
}

const (
	DebugEnabled = false
)

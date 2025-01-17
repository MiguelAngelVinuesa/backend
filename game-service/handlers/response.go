package handlers

import (
	"fmt"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git/state/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	log2 "git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/log"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/api/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/log"
)

func sendError(ctx *fiber.Ctx, label string, msg, req any, status int, data []byte) error {
	log.Logger.Error(label, consts.FieldRequest, req, consts.FieldError, msg)

	m := fmt.Sprintf(consts.ErrorCallFailed, label, msg, req)
	go log.ReportError(m)

	ctx.Set(consts.ContentType, consts.ApplicationJSON)
	return ctx.Status(status).Send(data)
}

func sendErrorStack(ctx *fiber.Ctx, label string, msg, req any, status int, data, stack []byte) error {
	log.Logger.Error(label, consts.FieldRequest, req, consts.FieldError, msg, consts.FieldStack, string(stack))

	m := fmt.Sprintf(consts.ErrorCallFailed, label, msg, req)
	go log.ReportError(m)

	ctx.Set(consts.ContentType, consts.ApplicationJSON)
	return ctx.Status(status).Send(data)
}

func sendNotModified(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusNotModified).Send([]byte{})
}

func sendResponse(label string, ctx *fiber.Ctx, req any, resp *zjson.Encoder) error {
	if log.API && log.Logger.Enabled(log2.DebugLevel) {
		log.Logger.Debug(label, consts.FieldRequest, req, consts.FieldResponse, string(resp.Bytes()))
	}

	ctx.Set(consts.ContentType, consts.ApplicationJSON)
	_, err := ctx.Write(resp.Bytes())
	resp.Release()

	return err
}

var (
	FmtInvalidSession = func(reason string, err error) string {
		if err != nil {
			return fmt.Sprintf("%s; %s; %v", consts.ErrorInvalidSession, reason, err)
		}
		return fmt.Sprintf("%s; %s", consts.ErrorInvalidSession, reason)
	}
	FmtDstoreError = func(err error) string {
		if e, ok := err.(*slots.APIerror); ok {
			return fmt.Sprintf(consts.ErrorDstoreError, e.Code, e.Message)
		}
		return fmt.Sprintf(consts.ErrorDstoreError, consts.ErrCdPanic, "unknown error")
	}

	BodyInternalError = func(code consts.ErrorCode, level string) []byte {
		j, _ := json.Marshal(&models.ErrorResponse{Message: consts.ErrorInternalError, ErrorCode: int64(code), ErrorLevel: level})
		return j
	}
	BodyBadRequest = func(code consts.ErrorCode, level string) []byte {
		j, _ := json.Marshal(&models.ErrorResponse{Message: consts.ErrorBadRequest, ErrorCode: int64(code), ErrorLevel: level})
		return j
	}
	BodyNotFound = func(code consts.ErrorCode, level string) []byte {
		j, _ := json.Marshal(&models.ErrorResponse{Message: consts.ErrorNotFound, ErrorCode: int64(code), ErrorLevel: level})
		return j
	}
	BodyInvalidSession = func(code consts.ErrorCode, level string) []byte {
		j, _ := json.Marshal(&models.ErrorResponse{Message: consts.ErrorInvalidSession, ErrorCode: int64(code), ErrorLevel: level})
		return j
	}
	BodyInvalidStatus = func(code consts.ErrorCode, level string) []byte {
		j, _ := json.Marshal(&models.ErrorResponse{Message: consts.ErrorInvalidStatus, ErrorCode: int64(code), ErrorLevel: level})
		return j
	}
	BodyDstoreError = func(err error) []byte {
		if e, ok := err.(*slots.APIerror); ok {
			j, _ := json.Marshal(&models.ErrorResponse{Message: FmtDstoreError(err), ErrorCode: int64(e.Code), ErrorLevel: e.Level})
			return j
		}

		j, _ := json.Marshal(&models.ErrorResponse{Message: consts.ErrorInvalidStatus, ErrorCode: int64(consts.ErrCdPanic), ErrorLevel: consts.ErrLvlFatal})
		return j
	}
)

package handlers

import (
	"runtime/debug"
	"strings"
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"github.com/gofiber/fiber/v2"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/events"
)

func GetMessages(req *fiber.Ctx) (err error) {
	// no metrics as this is a long-poll endpoint!

	params := struct {
		SessionID string
		Locale    string
	}{
		SessionID: strings.Clone(req.Query("sessionId")),
		Locale:    strings.Clone(req.Query(consts.FieldLocale)),
	}

	defer func() {
		if e := recover(); e != nil {
			err = sendErrorStack(req, consts.PathRound, e, params, fiber.StatusInternalServerError, BodyInternalError(consts.ErrCdPanic, consts.ErrLvlFatal), debug.Stack())
		}
	}()

	msgs := events.GetMessages(params.SessionID)
	if len(msgs) == 0 {
		// no messages, so long poll... wait for http context cancelled or timed out, local poll timeout or new message(s).
		trigger := make(chan struct{})
		timeout := time.Tick(config.LongPollTimeout)

		events.AddSession(params.SessionID, trigger)
		defer events.RemoveSession(params.SessionID)

		select {
		case <-req.Context().Done():
			return sendNotModified(req)
		case <-timeout:
			return sendNotModified(req)
		case <-trigger:
		}

		msgs = events.GetMessages(params.SessionID)
	}

	// remove messages after response was sent!
	defer events.CommitMessages(params.SessionID, msgs)

	enc := zjson.AcquireEncoder(1024)
	enc.StartArray()
	for ix := range msgs {
		msgs[ix].Encode(enc, params.Locale)
	}
	enc.EndArray()

	return sendResponse(consts.PathMessages, req, params.SessionID, enc)
}

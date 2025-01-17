package handlers

import (
	"net/http"
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/tg"
	"github.com/gofiber/fiber/v2"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/game"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/metrics"
)

func Ping(ctx *fiber.Ctx) error {
	started := time.Now()
	defer metrics.Metrics.AddDuration(metrics.ApiPing, started)

	// test our games to make sure the configs can load.
	ok := true
	for ix := range testGames {
		g := game.NewGame(testGames[ix].code, testGames[ix].rtp)
		if g == nil {
			ok = false
			break
		}
		g.Release()
	}

	if !ok {
		// k8s will restart us.
		return sendError(ctx, consts.PathPing, consts.ErrorInternalError, nil, http.StatusInternalServerError, BodyInternalError(consts.ErrCdPing, consts.ErrLvlFatal))
	}

	// send the 10-4!
	ctx.Set(consts.ContentType, consts.ApplicationJSON)
	return ctx.Send(consts.PingResponse)
}

type testGame struct {
	code tg.GameNR
	rtp  int
}

var testGames = []testGame{
	{code: tg.BOTnr, rtp: 92},
	{code: tg.BOTnr, rtp: 94},
	{code: tg.BOTnr, rtp: 96},
	// {code: tg.BOTnr, rtp: 42},
	{code: tg.CCBnr, rtp: 92},
	{code: tg.CCBnr, rtp: 94},
	{code: tg.CCBnr, rtp: 96},
	// {code: tg.CCBnr, rtp: 42},
	{code: tg.MGDnr, rtp: 92},
	{code: tg.MGDnr, rtp: 94},
	{code: tg.MGDnr, rtp: 96},
	// {code: tg.MGDnr, rtp: 42},
	{code: tg.LAMnr, rtp: 92},
	{code: tg.LAMnr, rtp: 94},
	{code: tg.LAMnr, rtp: 96},
	// {code: tg.LAMnr, rtp: 42},
	{code: tg.OWLnr, rtp: 92},
	{code: tg.OWLnr, rtp: 94},
	{code: tg.OWLnr, rtp: 96},
	// {code: tg.OWLnr, rtp: 42},
	{code: tg.FRMnr, rtp: 92},
	{code: tg.FRMnr, rtp: 94},
	{code: tg.FRMnr, rtp: 96},
	// {code: tg.FRMnr, rtp: 42},
	{code: tg.OFGnr, rtp: 92},
	{code: tg.OFGnr, rtp: 94},
	{code: tg.OFGnr, rtp: 96},
	// {code: tg.OFGnr, rtp: 42},
	{code: tg.FPRnr, rtp: 92},
	{code: tg.FPRnr, rtp: 94},
	{code: tg.FPRnr, rtp: 96},
	// {code: tg.OFGnr, rtp: 42},
}

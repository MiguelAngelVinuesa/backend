package handlers

import (
	"fmt"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git/state/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/i18n/localization/data"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/tg"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/api/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/clients/i18n"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/encode"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/game"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/metrics"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/state"
)

func GetGameInfo(req *fiber.Ctx) (err error) {
	started := time.Now()
	defer func() { metrics.Metrics.AddDuration(metrics.ApiGameInfo, started) }()

	var sessionID string
	defer func() {
		if e := recover(); e != nil {
			err = sendErrorStack(req, consts.PathGameInfo, e, sessionID, fiber.StatusInternalServerError, BodyInternalError(consts.ErrCdPanic, consts.ErrLvlFatal), debug.Stack())
		}
	}()

	// decode request.
	if sessionID = req.Query("sessionId"); sessionID == "" {
		return sendError(req, consts.PathGameInfo, err, nil, fiber.StatusBadRequest, BodyBadRequest(consts.ErrCdSessionEmpty, consts.ErrLvlFatal))
	}

	loc := data.DefaultLocale
	if prefered := string(req.Request().Header.Peek(consts.AcceptLanguage)); prefered != "" {
		loc = i18n.LocaleFromAcceptLanguage(prefered)
	}

	var sess *tg.SessionKey
	if sess, err = tg.VerifySessionID(sessionID); err != nil {
		return sendError(req, consts.PathGameInfo, FmtInvalidSession("verification", err), sessionID, fiber.StatusBadRequest, BodyInvalidSession(consts.ErrCdSessionInvalid, consts.ErrLvlFatal))
	}

	// mark active session.
	if len(sessionID) > 8 {
		metrics.MarkSession(sessionID)
	}

	// load game details.
	g := game.NewGame(sess.GameNr(), int(sess.RTP()))
	if g == nil {
		return sendError(req, consts.PathGameInfo, FmtInvalidSession("game", nil), sessionID, fiber.StatusBadRequest, BodyInvalidSession(consts.ErrCdGameInvalid, consts.ErrLvlFatal))
	}
	defer g.Release()

	return sendResponse(consts.PathGameInfo, req, sessionID, encode.GameInfo(loc, sess, g, config.DebugMode))
}

func PostPreferences(req *fiber.Ctx) (err error) {
	started := time.Now()
	defer func() { metrics.Metrics.AddDuration(metrics.ApiPreferences, started) }()

	params := &models.PreferencesRequest{}
	defer func() {
		if e := recover(); e != nil {
			err = sendErrorStack(req, consts.PathPreferences, e, params, fiber.StatusInternalServerError, BodyInternalError(consts.ErrCdPanic, consts.ErrLvlFatal), debug.Stack())
		}
	}()

	if err = req.BodyParser(params); err != nil {
		return sendError(req, consts.PathPreferences, err, params, fiber.StatusBadRequest, BodyBadRequest(consts.ErrCdParams, consts.ErrLvlFatal))
	}

	sess, err2 := tg.VerifySessionID(params.SessionID)
	if params.SessionID == "" || err2 != nil {
		return sendError(req, consts.PathPreferences, FmtInvalidSession("verification", err2), params, fiber.StatusBadRequest, BodyInvalidSession(consts.ErrCdSessionInvalid, consts.ErrLvlFatal))
	}
	if (params.Locale != nil && len(*params.Locale) > 5) ||
		(params.Music != nil && (*params.Music < 0 || *params.Music > 100)) ||
		(params.Effects != nil && (*params.Effects < 0 || *params.Effects > 100)) ||
		(params.Volume != nil && (*params.Volume < 0 || *params.Volume > 100)) ||
		(params.Bet != nil && (*params.Bet < 10 || *params.Bet > 50000 || *params.Bet%10 != 0)) {
		return sendError(req, consts.PathPreferences, consts.ErrorBadRequest, params, fiber.StatusBadRequest, BodyBadRequest(consts.ErrCdParams, consts.ErrLvlFatal))
	}

	// mark active session.
	metrics.MarkSession(params.SessionID)

	if params.Locale != nil {
		started2 := time.Now()
		pp, err3 := state.Manager.GetPlayerPrefs(params.SessionID)
		metrics.Metrics.AddDuration(metrics.DsPlayerPrefsGet, started2)

		if err3 != nil {
			pp = make(map[string]string)
		} else {
			delete(pp, consts.PrefLanguage)
		}

		loc := data.FixLocale(*params.Locale)
		if pp[consts.PrefLocale] != loc {
			pp[consts.PrefLocale] = loc

			started2 = time.Now()
			err3 = state.Manager.PutPlayerPrefs(params.SessionID, pp)
			metrics.Metrics.AddDuration(metrics.DsPlayerPrefsPut, started2)

			if err3 != nil {
				return sendError(req, consts.PathPreferences, fmt.Errorf("failed to put player preferences"), params, fiber.StatusInternalServerError, BodyInternalError(consts.ErrCdUpdateFailed, consts.ErrLvlRetry))
			}
		}
	}

	var gp *slots.GamePrefs

	if params.Music != nil || params.Effects != nil || params.Volume != nil || params.Bet != nil {
		started2 := time.Now()
		_, _, gp, err = state.Manager.GetGamePrefs(params.SessionID)
		metrics.Metrics.AddDuration(metrics.DsGamePrefsGet, started2)

		if err == nil {
			defer gp.Release()
		} else {
			gp = &slots.GamePrefs{}
		}

		var modified bool
		if params.Music != nil {
			s := strconv.Itoa(int(*params.Music))
			if gp.GamePref(consts.PrefMusic) != s {
				gp.SetGamePref(consts.PrefMusic, s)
				modified = true
			}
		}
		if params.Effects != nil {
			s := strconv.Itoa(int(*params.Effects))
			if gp.GamePref(consts.PrefEffects) != s {
				gp.SetGamePref(consts.PrefEffects, s)
				modified = true
			}
		}
		if params.Volume != nil {
			s := strconv.Itoa(int(*params.Volume))
			if gp.GamePref(consts.PrefVolume) != s {
				gp.SetGamePref(consts.PrefVolume, s)
				modified = true
			}
		}
		if params.Bet != nil {
			s := strconv.FormatInt(*params.Bet, 10)
			if gp.GamePref(consts.PrefBet) != s {
				gp.SetGamePref(consts.PrefBet, s)
				modified = true
			}
		}

		if modified {
			started2 = time.Now()
			err = state.Manager.PutGamePrefs(params.SessionID, gp)
			metrics.Metrics.AddDuration(metrics.DsGamePrefsPut, started2)
			if err != nil {
				return sendError(req, consts.PathPreferences, FmtDstoreError(err), params, fiber.StatusInternalServerError, BodyDstoreError(err))
			}
		}
	}

	if sess.GameNr() == tg.CCBnr && params.Bet != nil {
		if gp != nil {
			if flags := gp.GetStateCCB(*params.Bet); flags != nil {
				return sendResponse(consts.PathGameInfo, req, params, encode.CcbPreferences(flags))
			}
		}
	}

	req.Set(consts.ContentType, consts.ApplicationJSON)
	_, err = req.Write(consts.SuccessResponse)
	return err
}

func GetCcbFlags(req *fiber.Ctx) (err error) {
	started := time.Now()
	defer func() { metrics.Metrics.AddDuration(metrics.ApiCcbFlags, started) }()

	params := struct {
		SessionID string `query:"sessionId" json:"sessionID"`
		Bet       int64  `query:"bet" json:"bet"`
	}{}
	defer func() {
		if e := recover(); e != nil {
			err = sendErrorStack(req, consts.PathCcbFlags, e, params, fiber.StatusInternalServerError, BodyInternalError(consts.ErrCdPanic, consts.ErrLvlFatal), debug.Stack())
		}
	}()

	// decode request.
	if err = req.QueryParser(&params); err != nil {
		return sendError(req, consts.PathCcbFlags, err, params, fiber.StatusBadRequest, BodyBadRequest(consts.ErrCdParams, consts.ErrLvlFatal))
	}

	sess, err2 := tg.VerifySessionID(params.SessionID)
	if params.SessionID == "" || err2 != nil || sess.GameNr() != tg.CCBnr || params.Bet < 0 || params.Bet > 50000 {
		return sendError(req, consts.PathCcbFlags, FmtInvalidSession("verification", err2), params, fiber.StatusBadRequest, BodyInvalidSession(consts.ErrCdSessionInvalid, consts.ErrLvlFatal))
	}

	// mark active session.
	metrics.MarkSession(params.SessionID)

	started2 := time.Now()
	_, _, gp, err3 := state.Manager.GetGamePrefs(params.SessionID)
	metrics.Metrics.AddDuration(metrics.DsGamePrefsGet, started2)

	if err3 == nil {
		if flags := gp.GetStateCCB(params.Bet); flags != nil {
			return sendResponse(consts.PathCcbFlags, req, params, encode.CcbPreferences(flags))
		}
	}

	req.Set(consts.ContentType, consts.ApplicationJSON)
	_, err = req.Write(consts.SuccessResponse)
	return err
}

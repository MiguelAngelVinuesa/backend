package bo_backend

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git/state/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/i18n/localization/data"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/tg"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/utils/conv"
	"github.com/goccy/go-json"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/log"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/metrics"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/state"
)

// GetGamePrefs returns the preferences for the player, casino & jurisdiction based on the given session.
func GetGamePrefs(loc string, sess *tg.SessionKey) (*slots.GamePrefs, map[string]any, map[string]any) {
	sessionID := sess.SessionID()

	if len(sessionID) <= 5 {
		prefs := models.EmptyGamePrefs()
		return prefs, defaultCasino(sess), defaultJurisdiction(prefs)
	}

	started2 := time.Now()
	casinoID, playerID, prefs, err := state.Manager.GetGamePrefs(sessionID)
	metrics.Metrics.AddDuration(metrics.DsGamePrefsGet, started2)
	if err != nil || prefs == nil {
		prefs = models.EmptyGamePrefs()
	}

	addPlayerPrefs(sessionID, prefs)
	fixDefaultPrefs(loc, prefs)

	if sess.GameNr() == tg.CCBnr {
		addCcbFlags(sessionID, prefs)
	}

	casino, juris := loadBackOffice(casinoID, playerID, sess, prefs)
	return prefs, casino, juris
}

func addPlayerPrefs(sessionID string, prefs *slots.GamePrefs) {
	started2 := time.Now()
	m, err := state.Manager.GetPlayerPrefs(sessionID)
	metrics.Metrics.AddDuration(metrics.DsPlayerPrefsGet, started2)

	if err == nil && m != nil {
		for k, v := range m {
			if k == consts.PrefLocale {
				prefs.SetGamePref(consts.PrefLocale, data.FixLocale(v))
			} else {
				prefs.SetGamePref(k, v)
			}
		}
	}
}

func addCcbFlags(sessionID string, prefs *slots.GamePrefs) {
	s := prefs.GamePref(consts.PrefBet)
	b, err2 := strconv.ParseInt(s, 10, 64)
	if err2 != nil || b == 0 {
		b = 100
		prefs.SetGamePref(consts.PrefBet, dfltBet)
	}

	flags := prefs.GetStateCCB(b)

	var hasFlag bool
	if flags != nil {
		for _, f := range flags.Flagged() {
			if f {
				hasFlag = true
				break
			}
		}
	}

	if hasFlag {
		gs := slots.AcquireGameState(nil, flags, b)
		started := time.Now()
		state.Manager.PutGameState(sessionID, gs)
		metrics.Metrics.AddDuration(metrics.DsSessionPut, started)
		gs.Release()
	}
}

func defaultCasino(sess *tg.SessionKey) map[string]any {
	m := make(map[string]any)
	switch sess.GameNr() {
	case tg.OFGnr:
		m[keyBets] = betsOFG
	case tg.FRMnr:
		m[keyBets] = betsFRM
	default:
		m[keyBets] = betsGeneric
	}
	return m
}

func defaultJurisdiction(prefs *slots.GamePrefs) map[string]any {
	m := make(map[string]any)

	jurisdiction := codeUKGC
	if strings.HasPrefix(prefs.GamePref(consts.PrefLocale), "it") {
		jurisdiction = codeADM
	}
	m[keyCode] = jurisdiction

	switch jurisdiction {
	case codeADM:
		m[keySpinWait] = waitADM
	case codeUKGC:
		m[keySpinWait] = waitUKGC
	default:
		m[keySpinWait] = waitDefault
	}

	return m
}

func fixDefaultPrefs(loc string, prefs *slots.GamePrefs) {
	if prefs.GamePref(consts.PrefLocale) == "" {
		// find appropriate default locale
		prefs.SetGamePref(consts.PrefLocale, data.FixLocale(loc))
	}

	if prefs.GamePref(consts.PrefMusic) == "" {
		// default to 50% music volume
		prefs.SetGamePref(consts.PrefMusic, dfltMusic)
	}

	if prefs.GamePref(consts.PrefEffects) == "" {
		// default to 50% effects volume
		prefs.SetGamePref(consts.PrefEffects, dfltEffects)
	}

	if prefs.GamePref(consts.PrefVolume) == "" {
		// default to 50% overall volume
		prefs.SetGamePref(consts.PrefVolume, dfltVolume)
	}

	if prefs.GamePref(consts.PrefBet) == "" {
		// default to 100 cents in whatever currency
		prefs.SetGamePref(consts.PrefBet, dfltBet)
	}
}

func loadBackOffice(casinoID, playerID string, sess *tg.SessionKey, prefs *slots.GamePrefs) (casino map[string]any, juris map[string]any) {
	casino, juris = defaultCasino(sess), defaultJurisdiction(prefs)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	u, _ := url.Parse(config.BBhost)
	u.Path = "/v1/session-prefs/" + sess.SessionID()

	q := url.Values{}
	if casinoID != "" {
		q.Add("casino", casinoID)
	}
	if playerID != "" {
		q.Add("player", playerID)
	}
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		log.Logger.Error("failed to create session-prefs request", consts.FieldError, err)
		return
	}

	req.WithContext(ctx)
	req.Header.Set(consts.XApiKey, config.BBkey)

	started := time.Now()
	resp, err2 := http.DefaultClient.Do(req)
	metrics.Metrics.AddDuration(metrics.BoGamePrefs, started)
	cancel()

	if err2 != nil {
		log.Logger.Error("failed to retrieve session-prefs", consts.FieldError, err2)
		return
	}
	defer resp.Body.Close()

	b, err3 := io.ReadAll(resp.Body)
	if err3 != nil {
		log.Logger.Error("failed to read session-prefs", consts.FieldError, err3)
		return
	}

	m := make(map[string]any)
	if err = json.Unmarshal(b, &m); err != nil {
		log.Logger.Error("failed to decode session-prefs", consts.FieldError, err)
		return
	}

	if m2, ok := m["object"].(map[string]any); ok {
		for k, v := range m2 {
			switch k {
			case "bets":
				if l, ok2 := v.([]any); ok2 {
					bets := make([]int64, len(l))
					for ix := range l {
						bets[ix] = int64(conv.IntFromAny(l[ix]))
					}
					casino[keyBets] = bets
				}

			case "jurisdictionId":
				if s, ok2 := v.(string); ok2 {
					juris[keyCode] = s
				}

			case "spinWait":
				if f, ok2 := v.(float64); ok2 {
					juris[keySpinWait] = int64(f)
				}

			case "prefs":
				if m3, ok3 := v.(map[string]any); ok3 {
					for k3, v3 := range m3 {
						juris[k3] = v3
					}
				}
			}
		}
	}

	return
}

var (
	betsGeneric = []int64{
		10, 20, 30, 40, 50, 80, 100,
		// 150, 200, 250, 300, 400, 500, 600, 800, 1000,
		// 1500, 2000, 2500, 3000, 4000, 5000, 6000, 8000, 10000,
		// 15000, 20000, 25000, 30000, 40000, 50000,
	}
	betsOFG = []int64{
		25, 50, 75, 100,
		// 125, 250, 500, 625,
		// 1250, 1500, 2000, 2500, 3000, 5000,
	}
	betsFRM = []int64{
		20, 40, 60, 80, 100,
		// 200, 300, 400, 500, 600, 800,
		// 1000, 1500, 2000, 2500, 3000, 4000, 5000, 6000, 8000,
		// 10000, 15000, 20000, 25000, 30000, 40000, 50000,
	}
)

const (
	dfltMusic         = "50"  // 50% volume/toggled on.
	dfltEffects       = "50"  // 50% volume/toggled on.
	dfltVolume        = "50"  // 50% audio volume.
	dfltBet           = "100" // 100 cents.
	waitADM     int64 = 1000  // 1 second
	waitUKGC    int64 = 2500  // 2.5 seconds.
	waitDefault int64 = 1000  // 1 second.
	keyBets           = "bets"
	keyCode           = "code"
	keySpinWait       = "spinWait"
	codeADM           = "ADM"
	codeUKGC          = "UKGC"
)

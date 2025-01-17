package bo_backend

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/tg"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/utils/conv"
	"github.com/goccy/go-json"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/log"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/metrics"
)

// GetCasinoGameRTP returns the live RTP for the casino & game based on the given session.
func GetCasinoGameRTP(sess *tg.SessionKey) (float64, int, int) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	u, _ := url.Parse(config.BBhost)
	u.Path = "/v1/casino-game-rtp/" + sess.SessionID()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		log.Logger.Error("failed to create casino-game-rtp request", consts.FieldError, err)
		return 0, 0, 0
	}

	req.WithContext(ctx)
	req.Header.Set(consts.XApiKey, config.BBkey)

	started := time.Now()
	resp, err2 := http.DefaultClient.Do(req)
	metrics.Metrics.AddDuration(metrics.BoCasinoGameRTP, started)
	cancel()

	if err2 != nil {
		log.Logger.Error("failed to retrieve casino-game-rtp", consts.FieldError, err2)
		return 0, 0, 0
	}
	defer resp.Body.Close()

	b, err3 := io.ReadAll(resp.Body)
	if err3 != nil {
		log.Logger.Error("failed to read casino-game-rtp", consts.FieldError, err3)
		return 0, 0, 0
	}

	m := make(map[string]any)
	if err = json.Unmarshal(b, &m); err != nil {
		log.Logger.Error("failed to decode casino-game-rtp", consts.FieldError, err)
		return 0, 0, 0
	}

	m = m["object"].(map[string]any)
	return conv.FloatFromAny(m["rtp"]), conv.IntFromAny(m["count"]), conv.IntFromAny(m["hours"])
}

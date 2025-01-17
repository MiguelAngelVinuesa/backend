package bo_backend

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/utils/conv"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/api/models"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/log"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/metrics"
)

func GetSessionRounds(sessionID string, from, to time.Time, limit int) ([]*models.SessionInfoResponseRoundsItems0, error) {
	if limit <= 0 || limit > 250 {
		limit = 100
	}

	u, _ := url.Parse(config.BBhost)
	u.Path = "/v1/rounds/" + sessionID

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		log.Logger.Error("failed to create session-rounds request", consts.FieldError, err)
		return nil, err
	}

	q := url.Values{}
	if !from.IsZero() {
		q.Add("createdFrom", from.Format("2006-01-02T15:04:05Z"))
	}
	if !to.IsZero() {
		q.Add("createdTo", to.Format("2006-01-02T15:04:05Z"))
	}
	q.Add("limit", strconv.Itoa(limit))
	u.RawQuery = q.Encode()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req.WithContext(ctx)
	req.Header.Set(consts.XApiKey, config.BBkey)

	started := time.Now()
	resp, err2 := http.DefaultClient.Do(req)
	metrics.Metrics.AddDuration(metrics.BoGamePrefs, started)
	cancel()

	if err2 != nil {
		log.Logger.Error("failed to retrieve session-rounds", consts.FieldError, err2)
		return nil, err2
	}
	defer resp.Body.Close()

	b, err3 := io.ReadAll(resp.Body)
	resp.Body.Close()

	if err3 != nil || resp.StatusCode != fiber.StatusOK {
		log.Logger.Error("failed to read session-rounds", consts.FieldError, err3, consts.FieldResponse, string(b))
		return nil, err3
	}

	m := make(map[string]any)
	if err = json.Unmarshal(b, &m); err != nil {
		log.Logger.Error("failed to decode session-rounds", consts.FieldError, err, consts.FieldResponse, string(b))
		return nil, err
	}

	out := make([]*models.SessionInfoResponseRoundsItems0, 0, limit)
	if m2, ok2 := m["list"].([]any); ok2 {
		for ix := len(m2) - 1; ix >= 0; ix-- { // reverse the order so latest round is last!
			if m3, ok3 := m2[ix].(map[string]any); ok3 {
				out = append(out, &models.SessionInfoResponseRoundsItems0{
					RoundID: conv.StringFromAny(m3["id"]),
					Started: conv.StringFromAny(m3["started"]),
				})
			}
		}
	}
	return out, nil
}

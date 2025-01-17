package bo_backend

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

func GetSessionInfo(sessionID string, from, to time.Time, limit int) (*models.SessionInfoResponse, error) {
	u, _ := url.Parse(config.BBhost)
	u.Path = "/v1/session/" + sessionID

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		log.Logger.Error("failed to create session request", consts.FieldError, err)
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req.WithContext(ctx)
	req.Header.Set(consts.XApiKey, config.BBkey)

	started := time.Now()
	resp, err2 := http.DefaultClient.Do(req)
	metrics.Metrics.AddDuration(metrics.BoGamePrefs, started)
	cancel()

	if err2 != nil {
		log.Logger.Error("failed to retrieve session", consts.FieldError, err2)
		return nil, err2
	}
	defer resp.Body.Close()

	b, err3 := io.ReadAll(resp.Body)
	resp.Body.Close()

	if err3 != nil || resp.StatusCode != fiber.StatusOK {
		log.Logger.Error("failed to read session", consts.FieldError, err3, consts.FieldResponse, string(b))
		return nil, err3
	}

	m := make(map[string]any)
	if err = json.Unmarshal(b, &m); err != nil {
		log.Logger.Error("failed to decode session", consts.FieldError, err, consts.FieldResponse, string(b))
		return nil, err
	}

	m2, ok2 := m["object"].(map[string]any)
	if !ok2 {
		err = fmt.Errorf("no 'object' element in response data")
		log.Logger.Error("failed to decode session", consts.FieldError, err, consts.FieldResponse, string(b))
		return nil, err
	}

	out := &models.SessionInfoResponse{
		GameID:  conv.StringFromAny(m2["game"]),
		Rtp:     int64(conv.IntFromAny(m2["rtp"])),
		Started: conv.StringFromAny(m2["started"]),
		Updated: conv.StringFromAny(m2["updated"]),
		Success: true,
	}

	if out.Rounds, err = GetSessionRounds(sessionID, from, to, limit); err != nil {
		log.Logger.Error("failed to reetrieve session rounds", consts.FieldError, err)
		return nil, err
	}

	if m3, ok3 := m2["stats"].(map[string]any); ok3 {
		if m4, ok4 := m3["17"].(map[string]any); ok4 {
			if m5, ok5 := m4["."].(map[string]any); ok5 {
				out.Total = int64(conv.IntFromAny(m5["rounds"]))
			}
		}
	}
	return out, nil
}

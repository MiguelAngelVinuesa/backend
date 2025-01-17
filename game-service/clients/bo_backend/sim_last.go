package bo_backend

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/utils/conv"
	"github.com/goccy/go-json"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/config"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/log"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/metrics"
)

func ListSimulations(f func(game string, rtp1, rtp2 float64, count int)) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	u, _ := url.Parse(config.BBhost)
	u.Path = "/v1/sim-last/"

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		log.Logger.Error("failed to create sim-last request", consts.FieldError, err)
		return
	}

	req.WithContext(ctx)
	req.Header.Set(consts.XApiKey, config.BBkey)

	started := time.Now()
	resp, err2 := http.DefaultClient.Do(req)
	metrics.Metrics.AddDuration(metrics.BoSimulatorRTP, started)
	cancel()

	if err2 != nil {
		log.Logger.Error("failed to retrieve sim-last", consts.FieldError, err2)
		return
	}
	defer resp.Body.Close()

	b, err3 := io.ReadAll(resp.Body)
	if err3 != nil {
		log.Logger.Error("failed to read sim-last", consts.FieldError, err3)
		return
	}

	m := make(map[string]any)
	if err = json.Unmarshal(b, &m); err != nil {
		log.Logger.Error("failed to decode sim-last", consts.FieldError, err)
		return
	}

	if list, ok := m["list"].([]any); ok {
		for ix := range list {
			if m, ok = list[ix].(map[string]any); ok {
				game := conv.StringFromAny(m["game"])
				rtp := conv.IntFromAny(m["rtp"])
				if rtp > 80 && rtp < 99 {
					game += strconv.Itoa(rtp)
					choices := conv.StringFromAny(m["choices"])
					if choices == "" || choices == "RANDOM" {
						if m, ok = m["stats"].(map[string]any); ok {
							rtp1 := conv.FloatFromAny(m["rtp"])
							count := conv.IntFromAny(m["totalRounds"])
							rtp2 := 0.0 // TODO: ???

							if count <= 0 {
								count = 1000 * 1000 * 1000
							}

							f(game, rtp1, rtp2, count)
						}
					}
				}
			}
		}
	}
}

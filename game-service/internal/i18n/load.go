package i18n

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/utils/conv"
	"github.com/goccy/go-json"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/metrics"
)

func getTranslations(code string) (map[string]string, error) {
	path := fmt.Sprintf("%s/%s", stringAPI, code)
	started2 := time.Now()
	resp, err := http.DefaultClient.Get(path)
	metrics.Metrics.AddDuration(metrics.IsString, started2)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// unmarshall response.
	b, _ := io.ReadAll(resp.Body)
	m := make(map[string]any)
	if err = json.Unmarshal(b, &m); err != nil {
		return nil, err
	}

	var ok bool
	m, ok = m["string"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid response")
	}

	defaultLocale := conv.StringFromAny(m["defaultLocale"])
	defaultText := conv.StringFromAny(m["defaultText"])

	m, ok = m["translations"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid response")
	}

	out := make(map[string]string)
	suffix := "." + code

	if defaultLocale != "" && defaultText != "" {
		out[defaultLocale+suffix] = defaultText
	}

	for locale, v := range m {
		if m2, ok2 := v.(map[string]any); ok2 {
			if trans := conv.StringFromAny(m2["text"]); trans != "" {
				out[locale+suffix] = trans
			} else {
				out[locale+suffix] = defaultText
			}
		}
	}

	return out, nil
}

var stringAPI = "http://localhost:8002/v1/string"

func init() {
	if s := os.Getenv(consts.EnvI18nHost); s != "" {
		if !strings.HasPrefix(s, "http") {
			s = "http://" + s
		}
		uri, _ := url.Parse(s)
		stringAPI = fmt.Sprintf("%s://%s/v1/string", uri.Scheme, uri.Host)
	}
}

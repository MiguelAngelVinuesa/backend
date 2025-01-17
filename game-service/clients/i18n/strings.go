package i18n

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/valyala/fasthttp"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/metrics"
)

func GetStrings(locale string, keys []string) ([]string, error) {
	reqData := struct {
		Locale      string   `json:"locale"`
		Identifiers []string `json:"identifiers"`
	}{
		Locale:      locale,
		Identifiers: keys,
	}

	b, _ := json.Marshal(reqData)
	b2, err := PostStrings(b)
	if err != nil {
		return nil, err
	}

	respData := struct {
		Success bool              `json:"success,omitempty"`
		Strings map[string]string `json:"strings,omitempty"`
	}{}

	if err = json.Unmarshal(b2, &respData); err != nil {
		return nil, err
	}

	out := make([]string, len(keys))
	for i := range keys {
		out[i] = respData.Strings[keys[i]]
	}

	return out, nil
}

func PostStrings(b []byte) ([]byte, error) {
	req, err := http.NewRequest(http.MethodPost, stringsAPI, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	req.Header.Add(consts.ContentType, consts.ApplicationJSON)

	started2 := time.Now()
	resp, err2 := http.DefaultClient.Do(req)
	metrics.Metrics.AddDuration(metrics.IsStrings, started2)
	if err2 != nil {
		return nil, err2
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to post strings; status: %s", resp.Status)
	}

	return io.ReadAll(resp.Body)
}

func GetPluralString(loc, key string, args *fasthttp.Args) ([]byte, string, error) {
	path := fmt.Sprintf("%s/%s/%s", pluralAPI, loc, key)
	if args != nil && args.Len() > 0 {
		path += "?"
		path += args.String()
	}

	started2 := time.Now()
	resp, err := http.DefaultClient.Get(path)
	metrics.Metrics.AddDuration(metrics.IsPlural, started2)
	if err != nil {
		return nil, "", err
	}

	defer resp.Body.Close()

	b, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		return nil, "", err2
	}
	return b, resp.Header.Get(consts.ContentType), nil
}

func PostPluralStrings(b []byte) ([]byte, error) {
	started2 := time.Now()
	resp, err := http.DefaultClient.Post(pluralsAPI, "application/json", bytes.NewBuffer(b))
	metrics.Metrics.AddDuration(metrics.IsPlurals, started2)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to post plurals; status: %s", resp.Status)
	}
	return io.ReadAll(resp.Body)
}

func LocaleFromAcceptLanguage(preferred string) string {
	if req, err := http.NewRequest(http.MethodGet, acceptLangAPI, nil); err == nil {
		req.Header.Add(consts.AcceptLanguage, preferred)

		var resp *http.Response
		if resp, err = http.DefaultClient.Do(req); err == nil {
			defer resp.Body.Close()

			var buf []byte
			if buf, err = io.ReadAll(resp.Body); err == nil {
				var m map[string]any
				if json.Unmarshal(buf, &m) == nil {
					if m2, ok := m["locale"].(map[string]any); ok {
						if loc := m2["code"]; loc != "" {
							return loc.(string)
						}
					}
				}
			}
		}
	}

	return "en-GB"
}

var (
	stringsAPI    = "http://localhost:8003/v1/strings"
	pluralAPI     = "http://localhost:8003/v1/plural"
	pluralsAPI    = "http://localhost:8003/v1/plurals"
	acceptLangAPI = "http://localhost:8003/v1/accept-language"
)

func init() {
	if s := os.Getenv(consts.EnvI18nHost); s != "" {
		if !strings.HasPrefix(s, "http") {
			s = "http://" + s
		}

		uri, _ := url.Parse(s)
		stringsAPI = fmt.Sprintf("%s://%s/v1/strings", uri.Scheme, uri.Host)
		pluralAPI = fmt.Sprintf("%s://%s/v1/plural", uri.Scheme, uri.Host)
		pluralsAPI = fmt.Sprintf("%s://%s/v1/plurals", uri.Scheme, uri.Host)
		acceptLangAPI = fmt.Sprintf("%s://%s/v1/accept-language", uri.Scheme, uri.Host)
	}
}

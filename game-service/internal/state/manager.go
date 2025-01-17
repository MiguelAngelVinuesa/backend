package state

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	store "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git/repository/slots"
	manager "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git/state/slots"
	transport "git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/net/http"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/log"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/metrics"
)

var (
	DstoreHost     = ""
	ApiKey         = "[none]"
	DefaultTimeout = 10 * time.Second
	Manager        manager.RoundManager
)

func GetGameState(sessionID string) *manager.GameState {
	if gs, err := Manager.GetGameState(sessionID); err == nil {
		return gs
	}
	return nil
}

func Setup() {
	maxConn := 1000
	if s := os.Getenv(consts.EnvMaxClientConn); s != "" {
		if i, err := strconv.ParseInt(s, 10, 64); err == nil {
			maxConn = int(i)
		}
	}

	rt := http.DefaultTransport

	if s := os.Getenv(consts.EnvDstoreHost); s != "" {
		if strings.HasPrefix(s, "http") {
			DstoreHost = s
		} else {
			DstoreHost = fmt.Sprintf("http://%s/", s)
		}

		if resolver := os.Getenv(consts.EnvCustomDNS); resolver != "" {
			protocol := "tcp"
			resolveTimeout := 1 * time.Second
			cacheRefresh := time.Duration(0)
			clientRetries := 10
			retryBackoff := 10 * time.Millisecond

			if s = os.Getenv(consts.EnvResolveProtocol); s != "" {
				protocol = s
			}

			if s = os.Getenv(consts.EnvResolveTimeout); s != "" {
				if t, err := time.ParseDuration(s); err == nil {
					resolveTimeout = t
				}
			}

			if s = os.Getenv(consts.EnvDnsCacheRefresh); s != "" {
				if t, err := time.ParseDuration(s); err == nil {
					cacheRefresh = t
				}
			}

			if s = os.Getenv(consts.EnvClientRetries); s != "" {
				if i, err := strconv.ParseInt(s, 10, 64); err == nil {
					clientRetries = int(i)
				}
			}

			if s = os.Getenv(consts.EnvClientRetryBackoff); s != "" {
				if t, err := time.ParseDuration(s); err == nil {
					retryBackoff = t
				}
			}

			u, _ := url.Parse(DstoreHost)
			log.Logger.Info("custom resolver setup", "resolver", resolver, "protocol", protocol, "hosts", []string{u.Host})
			rt = transport.NewRetryTransport(
				clientRetries,
				retryBackoff,
				transport.RetryLogger(log.Logger),
				transport.RetryIdleConn(maxConn, maxConn, maxConn),
				transport.RetryResolver(resolver, protocol, resolveTimeout, cacheRefresh, u.Host))
		} else {
			c := http.DefaultTransport.(*http.Transport).Clone()
			c.MaxIdleConns = maxConn
			c.MaxConnsPerHost = maxConn
			c.MaxIdleConnsPerHost = maxConn
			rt = c
		}
	}

	if s := os.Getenv(consts.EnvDstoreApiKey); s != "" {
		ApiKey = s
	}

	if s := os.Getenv(consts.EnvClientTimeout); s != "" {
		if t, err := time.ParseDuration(s); err == nil {
			DefaultTimeout = t
		}
	}

	if DstoreHost == "" {
		Manager = store.NewMemory()
		log.Logger.Warn(consts.MsgNoValidation)
		metrics.PrintMetrics = 1 * time.Minute
	} else {
		Manager = store.NewDStore(DstoreHost, ApiKey, DefaultTimeout, rt, log.Logger, log.DsReq)
		log.Logger.Info(consts.MsgValidationSetup, consts.FieldURI, DstoreHost)
	}
}

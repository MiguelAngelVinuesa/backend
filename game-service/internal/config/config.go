package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/tg"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts"
)

var (
	ClientID         string
	Server           = "0.0.0.0"
	Port             = 8080
	Environment      = consts.ValueDev
	BBhost           string
	BBkey            string
	MonitorHost      string
	MonitorKey       string
	ConnCleanup      = 1 * time.Minute
	LongPollTimeout  = 1 * time.Minute
	MqBrokers        []string
	EventsTopic      string
	MessagesTopic    string
	HashesTopic      string
	ApiKey           = "abcdefghijklmn"
	ApiKeyBytes      = []byte(ApiKey)
	NoDefaultHeaders bool
	NoCors           bool
	NoCompression    bool
	DebugMode        = false // modified by compiler mode!
)

func init() {
	ClientID = tg.RandomBase42(8)

	if s := os.Getenv(consts.EnvSvcHost); s != "" {
		Server = s
	}
	if s := os.Getenv(consts.EnvSvcPort); s != "" {
		if i, err := strconv.ParseInt(s, 10, 64); err == nil && i >= 80 && i < 65536 {
			Port = int(i)
		}
	}

	if s := os.Getenv(consts.EnvEnvironment); s != "" {
		Environment = s
	}

	if s := os.Getenv(consts.EnvBBHost); s != "" {
		if !strings.HasPrefix(s, "http") {
			s = "http://" + s
		}
		BBhost = s
	}
	if s := os.Getenv(consts.EnvBBkey); s != "" {
		BBkey = s
	}

	if s := os.Getenv(consts.EnvMShost); s != "" {
		MonitorHost = s
	}
	if s := os.Getenv(consts.EnvMSkey); s != "" {
		MonitorKey = s
	}

	if s := os.Getenv(consts.EnvConnCleanup); s != "" {
		if t, err := time.ParseDuration(s); err == nil {
			ConnCleanup = t
		}
	}
	if s := os.Getenv(consts.EnvLongPoll); s != "" {
		if t, err := time.ParseDuration(s); err == nil {
			LongPollTimeout = t
		}
	}

	if s := os.Getenv(consts.EnvMsgBrokers); s != "" {
		MqBrokers = strings.Split(s, ",")
	}
	if s := os.Getenv(consts.EnvEventsTopic); s != "" {
		EventsTopic = s
	}
	if s := os.Getenv(consts.EnvMessagesTopic); s != "" {
		MessagesTopic = s
	}
	if s := os.Getenv(consts.EnvHashesTopic); s != "" {
		HashesTopic = s
	}

	if s := os.Getenv(consts.EnvApiKey); s != "" {
		ApiKey = s
		ApiKeyBytes = []byte(s)
	}

	if s := os.Getenv(consts.EnvNoDefaultHeaders); s != "" {
		NoDefaultHeaders = s == consts.ValueTrue
	}
	if s := os.Getenv(consts.EnvNoCors); s != "" {
		NoCors = s == consts.ValueTrue
	}
	if s := os.Getenv(consts.EnvNoCompression); s != "" {
		NoCompression = s == consts.ValueTrue
	}
}

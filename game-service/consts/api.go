package consts

import (
	"github.com/goccy/go-json"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/api/models"
)

const (
	PathPing        = "/v1/ping"
	PathStrings     = "/v1/strings"
	PathPlural      = "/v1/plural/:loc/:key"
	PathPlurals     = "/v1/plurals"
	PathBinHashes   = "/v1/bin-hashes"
	PathGameHash    = "/v1/game-hash/:game"
	PathGameInfo    = "/v1/game-info"
	PathPreferences = "/v1/preferences"
	PathCcbFlags    = "/v1/ccb-flags"
	PathMessages    = "/v1/messages"
	PathRound       = "/v1/round"
	PathRoundPaid   = "/v1/round/paid"
	PathRoundSecond = "/v1/round/second"
	PathRoundResume = "/v1/round/resume"
	// SUPERVISED-BUILD-REMOVE-START
	PathRoundDebug       = "/v1/round/debug"
	PathRoundDebugSecond = "/v1/round/debug-second"
	PathRoundDebugResume = "/v1/round/debug-resume"
	// SUPERVISED-BUILD-REMOVE-END
	PathRoundNext       = "/v1/round/next"
	PathRoundFinish     = "/v1/round/finish"
	PathSessionInfo     = "/v1/session/:session"
	PathRngConditionsLU = "/v1/rng-conditions-lu/:game"
	PathRngMagicTest    = "/v1/rng-magic/test"
	PathRngMagic        = "/v1/rng-magic"

	AcceptLanguage  = "Accept-Language"
	ContentType     = "Content-Type"
	XApiKey         = "X-API-Key"
	ApplicationJSON = "application/json"
	PlainText       = "text/plain; charset=utf-8"

	SharedRound = "shared-round"

	ErrLvlRetry = "R"
	ErrLvlFatal = "F"
)

type ErrorCode int

// always add new errors at the end!!
const (
	ErrCdPing ErrorCode = iota + 1000
	ErrCdApiKey
	ErrCdPanic
	ErrCdParams
	ErrCdSessionEmpty
	ErrCdSessionInvalid
	ErrCdGameInvalid
	ErrCdGameFailed
	ErrCdSpinStateInvalid
	ErrCdNotFound
	ErrCdUpdateFailed
	ErrCdRngFunctionInvalid
)

const (
	ErrorInternalError  = "internal server error"
	ErrorBadRequest     = "bad request"
	ErrorNotFound       = "not found"
	ErrorInvalidSession = "invalid sessionID, code or RTP"
	ErrorInvalidApiKey  = "invalid API key"
	ErrorInvalidStatus  = "invalid status or roundID"
	ErrorCallFailed     = "%s failed: %v; request: %v"
	ErrorDstoreError    = "D-store returned error: [%d] %s"
)

var (
	PingResponse    []byte
	SuccessResponse []byte
)

func init() {
	PingResponse, _ = json.Marshal(&models.PingResponse{Success: true})
	SuccessResponse, _ = json.Marshal(&models.RoundFinishResponse{Success: true})
}

package log

import (
	"fmt"
	"os"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/log"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/config"
)

var (
	Logger    log.Logger
	Output    = consts.LogSTDOUT
	Level     = consts.LogINFO
	Format    = consts.LogTEXT
	API       = true
	BigWin    = true
	FreeSpins = true
	DsTime    = false
	DsReq     = false
)

func Init() {
	if s := strings.ToUpper(os.Getenv(consts.EnvLogOutput)); s != "" {
		Output = s
	}
	if s := strings.ToUpper(os.Getenv(consts.EnvLogLevel)); s != "" {
		Level = s
	}
	if s := strings.ToUpper(os.Getenv(consts.EnvLogFormat)); s != "" {
		Format = s
	}
	if os.Getenv(consts.EnvLogAPI) == consts.LogONE {
		API = true
	}
	if os.Getenv(consts.EnvLogBigWin) == consts.LogONE {
		BigWin = true
	}
	if os.Getenv(consts.EnvLogFreeSpins) == consts.LogONE {
		FreeSpins = true
	}
	if os.Getenv(consts.EnvLogDsTiming) == consts.LogONE {
		DsTime = true
	}
	if os.Getenv(consts.EnvLogDsReqResp) == consts.LogONE {
		DsReq = true
	}

	Logger = log.InitZAP(Output, Level, Format, false)
	// Logger = log.InitZAP(Output, Level, Format, config.DebugMode)

	Logger.Info(consts.MsgSemVer, "version", consts.SemVerFull, "build", consts.SemDate)

	if config.DebugMode {
		Logger.Info(consts.MsgRunningDevMode)
	} else {
		Logger.Info(consts.MsgRunningProdMode)
	}
}

type InfoPrinter struct{}

func (p *InfoPrinter) Printf(format string, args ...any) {
	Logger.Info(fmt.Sprintf(format, args...))
}

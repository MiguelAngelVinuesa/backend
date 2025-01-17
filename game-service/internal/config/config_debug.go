//go:build DEBUG

// SUPERVISED-BUILD-REMOVE-START
// REMOVE ENTIRE FILE
// SUPERVISED-BUILD-REMOVE-END

package config

import (
	"os"
	"strings"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/consts"
)

func init() {
	if strings.EqualFold(os.Getenv(consts.EnvRunEnv), consts.ValueDev) {
		DebugMode = true
	}
}

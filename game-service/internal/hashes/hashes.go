package hashes

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

var (
	MainFile       string
	MainHash       = errorHash
	MainModDate    time.Time
	MainFileSize   int64
	RngLibHash     = errorHash
	RngLibModDate  time.Time
	RngLibFileSize int64
	RngIncludeHash = errorHash
	GameEngine     = errorHash
	GameConfig     = errorHash
	GameManager    = errorHash
	GameService    = errorHash
	Modules        = make([]*module, 0, 32)

	modFilter = []string{
		"amazonaws.com",
		"go-openapi/runtime",
		"goccy/go-json",
		"jessevdk/go-flags",
		"go.uber.org/zap",
		"gofiber/fiber",
		"valyala/fasthttp",
	}
)

const (
	RngLib         = "libprng.so"
	RngInclude     = "libprng.h"
	rngLibPath     = "/usr/local/lib/" + RngLib
	rngIncludePath = "/usr/local/include/" + RngInclude
	errorHash      = "!!ERROR!!"
	codeCommit     = "git-codecommit.eu-central-1.amazonaws.com/v1/repos"
)

func init() {
	RngLibHash = getHash(rngLibPath)

	if RngIncludeHash = getHash(rngIncludePath); RngIncludeHash == errorHash {
		RngIncludeHash = "[none]"
	}

	RngLibModDate, RngLibFileSize = getFileMeta(rngLibPath)

	mainPath := os.Args[0]
	if s, err := os.Stat(mainPath); err == nil {
		MainFile = s.Name()
		MainHash = getHash(mainPath)
		MainModDate, MainFileSize = getFileMeta(mainPath)
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Path == "" {
			Modules = append(Modules, &module{Path: "TopGaming/game-service", Version: GameService})
		} else {
			GameService = info.Main.Version
			Modules = append(Modules, &module{Path: replaceCodeCommit(info.Main.Path), Version: GameService})
		}

		for _, m := range info.Deps {
			for ix := range modFilter {
				if strings.Contains(m.Path, modFilter[ix]) && !strings.Contains(m.Path, "game-service") {
					Modules = append(Modules, &module{Path: replaceCodeCommit(m.Path), Version: m.Version})
					break
				}
			}

			switch {
			case strings.Contains(m.Path, "game-engine"):
				GameEngine = m.Version
			case strings.Contains(m.Path, "game-config"):
				GameConfig = m.Version
			case strings.Contains(m.Path, "game-manager"):
				GameManager = m.Version
			}
		}
	}
}

func replaceCodeCommit(input string) string {
	return strings.ReplaceAll(input, codeCommit, "TopGaming")
}

func getHash(path string) string {
	h := sha1.New()

	f, err := os.Open(path)
	if err != nil {
		return errorHash
	}
	defer f.Close()

	buf := make([]byte, 4096)
	var r int
	if r, err = f.Read(buf); err != nil {
		return errorHash
	}

	for r > 0 {
		h.Write(buf[:r])
		r, err = f.Read(buf)
		if err != nil && err != io.EOF {
			return errorHash
		}
	}

	return hex.EncodeToString(h.Sum(nil))
}

func getFileMeta(path string) (time.Time, int64) {
	s, err := os.Stat(path)
	if err != nil {
		return time.Time{}, 0
	}
	return s.ModTime().UTC(), s.Size()
}

type module struct {
	Path    string `json:"path,omitempty"`
	Version string `json:"version,omitempty"`
}

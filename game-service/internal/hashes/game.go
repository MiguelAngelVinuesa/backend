package hashes

import (
	"strconv"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/tg"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-service/internal/game"
)

var (
	GameHashes map[string]string
)

var (
	rtps  = []int{92, 94, 96}
	games = []string{"bot", "ccb", "mgd", "lam", "owl", "frm", "ofg", "fpr"}
)

func InitGameHashes() {
	GameHashes = make(map[string]string, len(rtps)*len(games))
	for ix := range games {
		id := games[ix]
		for iy := range rtps {
			rtp := rtps[iy]
			if nr, err := tg.VerifyGameID(id); err == nil {
				if g := game.NewGame(nr, rtp); g != nil {
					key := id + strconv.Itoa(rtp)
					GameHashes[key] = g.ConfigHash()
				}
			}
		}
	}
}

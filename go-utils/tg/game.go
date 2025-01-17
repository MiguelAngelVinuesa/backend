package tg

import (
	"fmt"
	"strings"
)

// GameNR is the unique numeric identifier for a game.
type GameNR uint32

// List of game numbers.
// ALWAYS MAKE SURE TO ADD NEW GAMES AT THE END OF THE LIST!
const (
	BOTnr GameNR = iota + 1
	CCBnr
	MGDnr
	LAMnr
	OWLnr
	FRMnr
	OFGnr
	FPRnr
	HOGnr
	MOGnr
	BBSnr
	BTRnr
	BERnr
	CASnr
	ANAnr
	CRWnr
	YYLnr
	FRJnr
)

// String implements the Stringer interface.
func (n GameNR) String() string {
	if s, ok := gameIds[n]; ok {
		return s
	}
	return "???"
}

// DSF returns if the game has the double-spin feature.
func (n GameNR) DSF() bool {
	switch n {
	case CCBnr:
		return true
	default:
		return false
	}
}

// VerifyGameID verifies if the gameID is supported by our game engine.
// This verification is irrespective of the RTP.
func VerifyGameID(gameID string) (GameNR, error) {
	gameID = strings.ToLower(gameID)
	if nr, ok := gameNrs[gameID]; ok {
		return nr, nil
	}

	// fallback on old codes.
	if newID := oldCodes[gameID]; newID != "" {
		if nr, ok := gameNrs[newID]; ok {
			return nr, nil
		}
	}

	return 0, fmt.Errorf("VerifyGameID: invalid GameID [%s]", gameID)
}

func fixOldCode(code string) string {
	if newID := oldCodes[code]; newID != "" {
		return newID
	}
	return code
}

var (
	gameNrs = map[string]GameNR{
		"bot": BOTnr,
		"ccb": CCBnr,
		"mgd": MGDnr,
		"lam": LAMnr,
		"owl": OWLnr,
		"frm": FRMnr,
		"ofg": OFGnr,
		"fpr": FPRnr,
		"hog": HOGnr,
		"mog": MOGnr,
		"bbs": BBSnr,
		"btr": BTRnr,
		"ber": BERnr,
		"cas": CASnr,
		"ana": ANAnr,
		"crw": CRWnr,
		"yyl": YYLnr,
		"frj": FRJnr,
	}

	gameIds = make(map[GameNR]string, len(gameNrs)) // Reverse lookup.

	oldCodes = map[string]string{
		"bob": "bot", // DEPRECATED; REMOVE IN FUTURE (just here as an example)
	}
)

func init() {
	for k, v := range gameNrs {
		gameIds[v] = k
	}
}

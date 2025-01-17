package tg

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
	"strings"
)

// MakeSessionID encodes the game info into a random sessionID and returns the result.
func MakeSessionID(gameID string, rtp uint8, serviceID uint32) (string, error) {
	gameNr, err := VerifyGameID(gameID)
	if err != nil {
		return "", err
	}
	if err = VerifyRTP(rtp); err != nil {
		return "", err
	}

	buf := make([]byte, 20)

	if _, err = rand.Read(buf[0:20]); err != nil {
		return "", fmt.Errorf("MakeSessionID - failed to read /dev/urandom: %v", err)
	}

	buf[3] = uint8(gameNr >> 16 & 0xff)
	binary.LittleEndian.PutUint16(buf[12:14], uint16(gameNr&0xffff))
	buf[1] = rtp
	binary.LittleEndian.PutUint16(buf[7:9], uint16(serviceID>>16))
	binary.LittleEndian.PutUint16(buf[16:18], uint16(serviceID&0xffff))

	return ToBase42(buf), nil
}

// VerifySessionID verifies a sessionID and returns the decoded game info.
func VerifySessionID(sessionID string) (*SessionKey, error) {
	if len(sessionID) < 6 {
		return nil, fmt.Errorf("VerifySessionID: invalid session [%s]", sessionID)
	}

	var gameID string
	var rtp uint8
	var serviceID uint32
	var gameNr GameNR
	var sharedRound bool

	if len(sessionID) > 15 {
		buf, err := FromBase42(sessionID)
		if err == nil {
			rtp = buf[1]
			serviceID = uint32(binary.LittleEndian.Uint16(buf[7:9]))<<16 | uint32(binary.LittleEndian.Uint16(buf[16:18]))
			gameNr = GameNR(uint32(buf[3])<<16 | uint32(binary.LittleEndian.Uint16(buf[12:14])))
			gameID = gameIds[gameNr]
		}
	}

	if rtp == 0 {
		if _, err := VerifyGameID(sessionID[:3]); err != nil {
			return nil, fmt.Errorf("VerifySessionID: invalid session [%s]; %v", sessionID, err)
		}

		gameID = fixOldCode(sessionID[:3])
		gameNr = gameNrs[gameID]

		if i, err2 := strconv.Atoi(sessionID[3:5]); err2 != nil || i >= math.MaxUint8 || VerifyRTP(uint8(i)) != nil {
			return nil, fmt.Errorf("VerifySessionID: invalid session [%s]; %v", sessionID, err2)
		} else {
			rtp = uint8(i)
		}

		if sessionID[5:] == "-shared-round" {
			sharedRound = true
		}
	}

	if _, err := VerifyGameID(gameID); err != nil {
		return nil, fmt.Errorf("VerifySessionID: invalid session [%s]; %v", sessionID, err)
	}
	if err := VerifyRTP(rtp); err != nil {
		return nil, fmt.Errorf("VerifySessionID: invalid session [%s]; %v", sessionID, err)
	}

	return &SessionKey{
		sharedRound: sharedRound,
		rtp:         rtp,
		serviceID:   serviceID,
		gameNr:      gameNr,
		gameID:      gameID,
		sessionID:   sessionID,
	}, nil
}

// SessionKey contains the decoded details of a sessionID.
type SessionKey struct {
	sharedRound bool
	rtp         uint8
	serviceID   uint32
	gameNr      GameNR
	gameID      string
	sessionID   string
}

// SharedRound returns if the session indicates a shared round.
func (s *SessionKey) SharedRound() bool {
	return s.sharedRound
}

// RTP returns the game RTP.
func (s *SessionKey) RTP() uint8 {
	return s.rtp
}

// GameID returns the gameID.
func (s *SessionKey) GameID() string {
	return s.gameID
}

// GameNr returns the gameNr.
func (s *SessionKey) GameNr() GameNR {
	return s.gameNr
}

// DSF returns if the game has the double-spin feature.
func (s *SessionKey) DSF() bool {
	return s.gameNr.DSF()
}

// GameUpper returns the gameID in uppercase.
func (s *SessionKey) GameUpper() string {
	return strings.ToUpper(s.gameID)
}

// ServiceID returns the serviceID.
func (s *SessionKey) ServiceID() uint32 {
	return s.serviceID
}

// SessionID returns the sessionID.
func (s *SessionKey) SessionID() string {
	return s.sessionID
}

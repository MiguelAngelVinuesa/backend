package tg

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMakeSessionID(t *testing.T) {
	testCases := []struct {
		name    string
		game    string
		rtp     uint8
		service uint32
		fail    bool
	}{
		{name: "bad game", game: "qqq", rtp: 92, service: 1, fail: true},
		{name: "bad rtp", game: "bot", rtp: 80, service: 2, fail: true},
		{name: "bob 92 (old)", game: "bot", rtp: 92, service: 11},
		{name: "bot 92", game: "bot", rtp: 92, service: 11},
		{name: "bot 94", game: "bot", rtp: 94, service: 12},
		{name: "bot 96", game: "bot", rtp: 96, service: 13},
		{name: "bot 42", game: "bot", rtp: 42, service: 14},
		{name: "ccb 92", game: "ccb", rtp: 92, service: 21},
		{name: "ccb 94", game: "ccb", rtp: 94, service: 22},
		{name: "ccb 96", game: "ccb", rtp: 96, service: 23},
		{name: "ccb 42", game: "ccb", rtp: 42, service: 24},
		{name: "mgd 92", game: "mgd", rtp: 92, service: 31},
		{name: "mgd 94", game: "mgd", rtp: 94, service: 32},
		{name: "mgd 96", game: "mgd", rtp: 96, service: 33},
		{name: "mgd 42", game: "mgd", rtp: 42, service: 34},
		{name: "lam 92", game: "lam", rtp: 92, service: 41},
		{name: "lam 94", game: "lam", rtp: 94, service: 42},
		{name: "lam 96", game: "lam", rtp: 96, service: 43},
		{name: "lam 42", game: "lam", rtp: 42, service: 44},
		{name: "owl 92", game: "owl", rtp: 92, service: 51},
		{name: "owl 94", game: "owl", rtp: 94, service: 52},
		{name: "owl 96", game: "owl", rtp: 96, service: 53},
		{name: "owl 42", game: "owl", rtp: 42, service: 54},
		{name: "frm 92", game: "frm", rtp: 92, service: 61},
		{name: "frm 94", game: "frm", rtp: 94, service: 62},
		{name: "frm 96", game: "frm", rtp: 96, service: 63},
		{name: "frm 42", game: "frm", rtp: 42, service: 64},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			id, err := MakeSessionID(tc.game, tc.rtp, tc.service)
			if tc.fail {
				require.Error(t, err)
				require.Empty(t, id)
			} else {
				require.NoError(t, err)
				require.NotEmpty(t, id)
				assert.Equal(t, 30, len(id))

				_, err = FromBase42(id)
				require.NoError(t, err)

				var s *SessionKey
				s, err = VerifySessionID(id)
				require.NoError(t, err)
				require.NotNil(t, s)
				assert.Equal(t, tc.game, s.GameID())
				assert.Equal(t, strings.ToUpper(tc.game), s.GameUpper())
				assert.Equal(t, tc.rtp, s.RTP())
				assert.Equal(t, tc.service, s.ServiceID())
				assert.Equal(t, id, s.SessionID())
			}
		})
	}
}

func TestVerifySessionID(t *testing.T) {
	testCases := []struct {
		name    string
		session string
		fail    bool
		rtp     uint8
		game    string
		service uint32
		shared  bool
	}{
		{name: "empty", fail: true},
		{name: "short", session: "bot96", fail: true},
		{name: "old bob96", session: "bob96x", game: "bot", rtp: 96},
		{name: "old bot94", session: "bot94x", game: "bot", rtp: 94},
		{name: "old ccb42", session: "ccb42x", game: "ccb", rtp: 42},
		{name: "old mgd92", session: "mgd92x", game: "mgd", rtp: 92},
		{name: "fpr96 shared round", session: "fpr96-shared-round", game: "fpr", rtp: 96, shared: true},
		{name: "new bot42", session: "7BAY43Ea7a33s4GYsY433XffJ3363g", game: "bot", rtp: 42, service: 14},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s, err := VerifySessionID(tc.session)
			if tc.fail {
				require.Error(t, err)
				assert.Nil(t, s)
			} else {
				require.NoError(t, err)
				require.NotNil(t, s)
				assert.Equal(t, tc.rtp, s.RTP())
				assert.Equal(t, tc.game, s.GameID())
				assert.Equal(t, tc.service, s.ServiceID())
				assert.Equal(t, tc.session, s.SessionID())
				assert.Equal(t, tc.shared, s.SharedRound())
			}
		})
	}
}

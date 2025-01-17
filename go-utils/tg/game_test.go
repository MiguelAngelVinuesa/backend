package tg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerifyGameID(t *testing.T) {
	testCases := []struct {
		name string
		game string
		fail bool
		want GameNR
	}{
		{name: "empty", game: "", fail: true},
		{name: "qqq", game: "qqq", fail: true},
		{name: "bob96", game: "bob96", fail: true},
		{name: "bob", game: "bob", want: BOTnr},
		{name: "bot", game: "bot", want: BOTnr},
		{name: "ccb", game: "ccb", want: CCBnr},
		{name: "mgd", game: "mgd", want: MGDnr},
		{name: "lam", game: "lam", want: LAMnr},
		{name: "owl", game: "owl", want: OWLnr},
		{name: "frm", game: "frm", want: FRMnr},
		{name: "ofg", game: "ofg", want: OFGnr},
		{name: "fpr", game: "fpr", want: FPRnr},
		{name: "hog", game: "hog", want: HOGnr},
		{name: "mog", game: "mog", want: MOGnr},
		{name: "bbs", game: "bbs", want: BBSnr},
		{name: "btr", game: "btr", want: BTRnr},
		{name: "ber", game: "ber", want: BERnr},
		{name: "cas", game: "cas", want: CASnr},
		{name: "ana", game: "ana", want: ANAnr},
		{name: "crw", game: "crw", want: CRWnr},
		{name: "yyl", game: "yyl", want: YYLnr},
		{name: "frj", game: "frj", want: FRJnr},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := VerifyGameID(tc.game)
			if tc.fail {
				require.Error(t, err)
				assert.Zero(t, got)
				assert.Equal(t, "???", got.String())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tc.want, got)

				if tc.game == "bob" {
					assert.Equal(t, "bot", got.String())
				} else {
					assert.Equal(t, tc.game, got.String())
				}

				if tc.game == "ccb" {
					assert.True(t, got.DSF())
				} else {
					assert.False(t, got.DSF())
				}
			}
		})
	}
}

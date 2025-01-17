package slots

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRoundState(t *testing.T) {
	testCases := []struct {
		name         string
		sessionID    string
		roundID      string
		playedFull   bool
		replayedFull bool
		spins        int
		playSeq      int
		replaySeq    int
		want         string
	}{
		{
			name: "empty",
			want: `{}`,
		},
		{
			name:       "played full",
			playedFull: true,
			want:       `{"playedFull":1}`,
		},
		{
			name:         "replayed full",
			replayedFull: true,
			want:         `{"replayedFull":1}`,
		},
		{
			name:  "spins",
			spins: 55,
			want:  `{"spins":55}`,
		},
		{
			name:    "play seq",
			playSeq: 3,
			want:    `{"playSeq":3}`,
		},
		{
			name:      "replay seq",
			replaySeq: 1,
			want:      `{"replaySeq":1}`,
		},
	}

	enc := zjson.AcquireEncoder(512)
	defer enc.Release()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s1 := AcquireRoundState(tc.sessionID, tc.roundID, tc.spins)
			require.NotNil(t, s1)
			assert.Equal(t, tc.sessionID, s1.sessionID)
			assert.Equal(t, tc.roundID, s1.roundID)
			assert.Equal(t, tc.spins, s1.spins)
			assert.False(t, s1.playedFull)
			assert.False(t, s1.replayedFull)
			assert.Zero(t, s1.playSeq)
			assert.Zero(t, s1.replaySeq)
			assert.Nil(t, s1.flags)
			assert.Nil(t, s1.played)
			assert.Nil(t, s1.replayed)

			s1.playedFull = tc.playedFull
			s1.replayedFull = tc.replayedFull
			s1.playSeq = tc.playSeq
			s1.replaySeq = tc.replaySeq

			enc.Reset()
			enc.Object(s1)
			got := enc.Bytes()
			assert.EqualValues(t, tc.want, string(got))

			s2, err := AcquireRoundStateFromJSON(tc.sessionID, tc.roundID, got)
			require.NoError(t, err)
			require.NotNil(t, s2)

			assert.Equal(t, s1.sessionID, s2.sessionID)
			assert.Equal(t, s1.roundID, s2.roundID)
			assert.Equal(t, s1.spins, s2.spins)
			assert.Equal(t, s1.playedFull, s2.playedFull)
			assert.Equal(t, s1.replayedFull, s2.replayedFull)
			assert.Equal(t, s1.playSeq, s2.playSeq)
			assert.Equal(t, s1.replaySeq, s2.replaySeq)
		})
	}
}

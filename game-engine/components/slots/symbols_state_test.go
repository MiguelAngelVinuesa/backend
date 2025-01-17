package slots

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
)

func TestNewSymbolsState(t *testing.T) {
	symbols := NewSymbolSet(s1, s2, s3, wf1, w1)

	testCases := []struct {
		name string
		flag []utils.Index
		want []bool
		j    string
		all  bool
	}{
		{
			name: "none flagged",
			want: []bool{false, false, false, false, false, false},
			j:    `{"flagged":[0,0,0,0,0,0],"valid":[0,1,1,1,1,1]}`,
		},
		{
			name: "1 flagged",
			flag: []utils.Index{2},
			want: []bool{false, false, true, false, false, false},
			j:    `{"flagged":[0,0,1,0,0,0],"valid":[0,1,1,1,1,1]}`,
		},
		{
			name: "3 flagged",
			flag: []utils.Index{2, 3, 5},
			want: []bool{false, false, true, true, false, true},
			j:    `{"flagged":[0,0,1,1,0,1],"valid":[0,1,1,1,1,1]}`,
		},
		{
			name: "all flagged",
			flag: []utils.Index{5, 3, 1, 2, 4},
			want: []bool{false, true, true, true, true, true},
			all:  true,
			j:    `{"flagged":[0,1,1,1,1,1],"valid":[0,1,1,1,1,1]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			s := AcquireSymbolsState(symbols)
			require.NotNil(t, s)
			defer s.Release()

			assert.Equal(t, int(symbols.maxID+1), len(s.flagged))
			assert.False(t, s.AllFlagged())
			for ix := range s.flagged {
				assert.False(t, s.flagged[ix])
			}

			for ix := range tc.flag {
				s.SetFlagged(tc.flag[ix], true)
			}

			assert.EqualValues(t, tc.want, s.flagged)
			assert.Equal(t, tc.all, s.AllFlagged())

			for ix := range tc.flag {
				assert.True(t, s.IsFlagged(tc.flag[ix]))
			}

			enc := zjson.AcquireEncoder(256)
			defer enc.Release()

			enc.Object(s)
			assert.Equal(t, tc.j, string(enc.Bytes()))

			s.ResetState()
			assert.False(t, s.AllFlagged())
			for ix := range s.flagged {
				assert.False(t, s.flagged[ix])
			}
		})
	}
}

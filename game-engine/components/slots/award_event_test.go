package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
)

func TestAcquireBonusGameEvent(t *testing.T) {
	testCases := []struct {
		name string
		want string
	}{
		{"bonus game", `{"kind":8,"bonusGame":1}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := AcquireBonusGameEvent().(*AwardEvent)
			require.NotNil(t, e)
			defer e.Release()

			assert.True(t, e.bonusGame)

			enc := zjson.AcquireEncoder(1024)
			defer enc.Release()

			e.Encode(enc)
			got := enc.Bytes()
			assert.Equal(t, tc.want, string(got))
		})
	}
}

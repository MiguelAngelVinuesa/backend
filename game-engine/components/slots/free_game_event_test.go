package slots

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
)

func TestAcquireFreeGameEvent(t *testing.T) {
	testCases := []struct {
		name  string
		count uint64
		want  string
	}{
		{name: "1", count: 1, want: `{"kind":7,"count":1}`},
		{name: "5", count: 5, want: `{"kind":7,"count":5}`},
		{name: "13", count: 13, want: `{"kind":7,"count":13}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			e := AcquireFreeGameEvent(tc.count).(*FreeGameEvent)
			require.NotNil(t, e)
			defer e.Release()

			assert.Equal(t, tc.count, e.count)

			enc := zjson.AcquireEncoder(1024)
			defer enc.Release()

			e.Encode(enc)
			got := enc.Bytes()
			assert.Equal(t, tc.want, string(got))
		})
	}
}

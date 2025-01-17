package slots

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"

	state "git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-manager.git/state/slots"
)

var (
	testPath   = "https://ds.dev.topgaming.team"
	testApiKey = "xyz"
)

func TestNewValidator(t *testing.T) {
	t.Run("new validator", func(t *testing.T) {
		m := NewDStore(testPath, testApiKey, 10*time.Second, http.DefaultTransport, nil, false)
		require.NotNil(t, m)

		d := m.(*dstore)
		assert.Equal(t, testPath, d.path)
		assert.Equal(t, testApiKey, d.apikey)

		m = NewDStore("nothing", "haha", 0, http.DefaultTransport, nil, false)
		require.NotNil(t, m)

		d = m.(*dstore)
		assert.Equal(t, "nothing", d.path)
		assert.Equal(t, "haha", d.apikey)
	})
}

func TestValidator_NewRound(t *testing.T) {
	v := NewDStore("xyz", "abc", 0, http.DefaultTransport, nil, false)
	require.NotNil(t, v)

	r1 := results.AcquireResult(nil, 0,
		slots.WinlinePayoutFromData(10, 0, 1, 3, slots.PayLTR, 1, utils.UInt8s{1, 1, 1, 1, 1}),
	)
	r2 := results.AcquireResult(nil, 0,
		slots.WinlinePayoutFromData(50, 0, 1, 3, slots.PayLTR, 1, utils.UInt8s{1, 1, 1, 1, 1}),
		slots.WinlinePayoutFromData(200, 0, 2, 4, slots.PayLTR, 1, utils.UInt8s{1, 1, 1, 1, 1}),
		slots.WinlinePayoutFromData(750, 0, 3, 5, slots.PayLTR, 1, utils.UInt8s{1, 1, 1, 1, 1}),
	)
	r3 := results.AcquireResult(nil, 0,
		slots.WinlinePayoutFromData(7.5, 0, 1, 3, slots.PayLTR, 1, utils.UInt8s{1, 1, 1, 1, 1}),
		slots.WinlinePayoutFromData(20, 0, 2, 4, slots.PayLTR, 1, utils.UInt8s{1, 1, 1, 1, 1}),
		slots.WinlinePayoutFromData(100, 0, 3, 5, slots.PayLTR, 1, utils.UInt8s{1, 1, 1, 1, 1}),
	)

	testCases := []struct {
		name    string
		session string
		bet     int64
		win     int64
		results results.Results
	}{
		{
			name:    "empty",
			results: results.Results{},
		},
		{
			name:    "session",
			session: "haha",
			results: results.Results{},
		},
		{
			name:    "bet",
			bet:     100,
			results: results.Results{},
		},
		{
			name:    "win",
			win:     200,
			results: results.Results{},
		},
		{
			name:    "1 payout",
			bet:     5,
			results: results.Results{r1},
		},
		{
			name:    "3 payouts",
			bet:     1,
			results: results.Results{r2},
		},
		{
			name:    "all - bet/win",
			session: "haha",
			bet:     100,
			win:     250,
			results: results.Results{},
		},
		{
			name:    "all - multiple results",
			session: "haha",
			bet:     10,
			results: results.Results{r1, r3},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			params := state.RoundParams{
				SessionID: tc.session,
				Bet:       tc.bet,
				TotalBet:  tc.bet,
				TotalWin:  tc.win,
				Results:   tc.results,
			}

			r := state.AcquireRound(v, params)
			require.NotNil(t, r)
			defer r.Release()

			assert.Equal(t, tc.session, r.SessionID())
			assert.Equal(t, tc.bet, r.Bet())
			assert.Equal(t, tc.results, r.Results())
		})
	}
}

func TestRound_Validate(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		in := make(map[string]interface{})
		buf, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		err = json.Unmarshal(buf, &in)
		require.NoError(t, err)

		assert.Equal(t, "hoho", in["sessionId"])
		assert.Equal(t, float64(100), in["bet"])
		assert.Equal(t, float64(200), in["win"])

		w.Header().Add("Content-Type", "application/json")
		_, err = w.Write([]byte(`{"success":true,"playerData":{"balance":1000}}`))
		require.NoError(t, err)
	}))
	require.NotNil(t, s)

	v := NewDStore(s.URL, testApiKey, time.Second, http.DefaultTransport, nil, false)
	require.NotNil(t, v)

	t.Run("round validate", func(t *testing.T) {
		r := state.AcquireRound(v, state.RoundParams{SessionID: "hoho", Bet: 100, TotalBet: 100, TotalWin: 200})
		require.NotNil(t, r)
		defer r.Release()

		assert.False(t, r.IsValid())

		r.Validate(false, false)
		assert.True(t, r.IsValid())
		assert.Equal(t, int64(1000), r.PlayerBalance())
	})
}

func TestRound_ValidateFail1(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		_, err := w.Write([]byte(`{"success":false}`))
		require.NoError(t, err)
	}))
	require.NotNil(t, s)

	v := NewDStore(s.URL, testApiKey, time.Second, http.DefaultTransport, nil, false)
	require.NotNil(t, v)

	t.Run("round validate", func(t *testing.T) {
		r := state.AcquireRound(v, state.RoundParams{SessionID: "hoho", Bet: 100, TotalBet: 100, TotalWin: 200})
		require.NotNil(t, r)
		defer r.Release()

		assert.False(t, r.IsValid())

		r.Validate(false, false)
		assert.False(t, r.IsValid())
		assert.Equal(t, int64(200), r.NewBalance())
	})
}

func TestRound_ValidateFail2(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		_, err := w.Write([]byte(`[]`))
		require.NoError(t, err)
	}))
	require.NotNil(t, s)

	v := NewDStore(s.URL, testApiKey, time.Second, http.DefaultTransport, nil, false)
	require.NotNil(t, v)

	t.Run("round validate", func(t *testing.T) {
		r := state.AcquireRound(v, state.RoundParams{SessionID: "hoho", Bet: 100, TotalBet: 100, TotalWin: 200})
		require.NotNil(t, r)
		defer r.Release()

		assert.False(t, r.IsValid())

		r.Validate(false, false)
		assert.False(t, r.IsValid())
		assert.Equal(t, int64(200), r.NewBalance())
	})
}

func TestRound_ValidateFail3(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
	}))
	require.NotNil(t, s)

	v := NewDStore(s.URL, testApiKey, time.Second, http.DefaultTransport, nil, false)
	require.NotNil(t, v)

	t.Run("round validate", func(t *testing.T) {
		r := state.AcquireRound(v, state.RoundParams{SessionID: "hoho", Bet: 100, TotalBet: 100, TotalWin: 200})
		require.NotNil(t, r)
		defer r.Release()

		assert.False(t, r.IsValid())

		r.Validate(false, false)
		assert.False(t, r.IsValid())
		assert.Equal(t, int64(200), r.NewBalance())
	})
}

package slots

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/components/slots"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/game-engine.git/results"
)

type dm struct{}

func (m *dm) PostRound(*Round, bool) (string, int64, error)                      { return "", 0, nil }
func (m *dm) PostInitRound(*Round, bool) (string, int64, error)                  { return "", 0, nil }
func (m *dm) PostCompleteRound(*Round, *RoundState, bool) (string, int64, error) { return "", 0, nil }
func (m *dm) PostRoundNext(string, string, *RoundState, int) (*RoundResult, int64, error) {
	return nil, 0, nil
}
func (m *dm) GetRoundState(string, string) (*RoundState, error)       { return nil, nil }
func (m *dm) PutGameState(string, *GameState) error                   { return nil }
func (m *dm) GetGameState(string) (*GameState, error)                 { return nil, nil }
func (m *dm) GetGamePrefs(string) (string, string, *GamePrefs, error) { return "", "", nil, nil }
func (m *dm) PutGamePrefs(string, *GamePrefs) error                   { return nil }
func (m *dm) GetPlayerPrefs(string) (map[string]string, error)        { return nil, nil }
func (m *dm) PutPlayerPrefs(string, map[string]string) error          { return nil }

func TestNewRound(t *testing.T) {
	testCases := []struct {
		name         string
		sessionID    string
		roundID      string
		bet          int64
		win          int64
		results      results.Results
		state        *GameState
		startBalance int64
		balances     []int64
		wantWin      int64
		newBalance   int64
		progressive  []int64
		maxPayout    float64
	}{
		{
			name:     "empty",
			balances: []int64{0, 0},
		},
		{
			name:      "session",
			sessionID: "haha",
			balances:  []int64{0, 0},
		},
		{
			name:     "round",
			roundID:  "hi ho",
			balances: []int64{0, 0},
		},
		{
			name:         "start balance",
			startBalance: 12500,
			balances:     []int64{12500, 12500},
			newBalance:   12500,
		},
		{
			name:     "bet",
			bet:      100,
			balances: []int64{0, 0},
		},
		{
			name:       "win",
			win:        200,
			balances:   []int64{0, 200},
			wantWin:    200,
			newBalance: 200,
		},
		{
			name:     "nil results",
			results:  results.Results{results.AcquireResult(nil, 0)},
			balances: []int64{0, 0},
		},
		{
			name:     "valid results",
			results:  results.Results{r1},
			balances: []int64{0, 0},
		},
		{
			name:     "state",
			state:    AcquireGameState(nil, nil, 0),
			balances: []int64{0, 0},
		},
		{
			name:      "max payout",
			maxPayout: 20000.0,
			balances:  []int64{0, 0},
		},
		{
			name:         "start balance + bet + win",
			startBalance: 1000,
			bet:          10,
			win:          50,
			wantWin:      50,
			newBalance:   1050,
		},
		{
			name:        "bet + valid result",
			bet:         25,
			results:     results.Results{r2},
			balances:    []int64{0, 250},
			wantWin:     250,
			newBalance:  250,
			progressive: []int64{250},
		},
		{
			name:         "all + win",
			sessionID:    "hihi",
			roundID:      "jolly",
			startBalance: 2000,
			bet:          20,
			win:          60,
			wantWin:      60,
			newBalance:   2060,
		},
		{
			name:         "all + valid result",
			sessionID:    "hoho",
			roundID:      "best round",
			startBalance: 10000,
			bet:          30,
			results:      results.Results{r3},
			balances:     []int64{10000, 10135},
			wantWin:      135,
			newBalance:   10135,
			progressive:  []int64{135},
		},
		{
			name:         "all + valid results",
			sessionID:    "nono",
			roundID:      "nice one",
			startBalance: 23575,
			bet:          5,
			results:      results.Results{r4, r5},
			balances:     []int64{23575, 23600, 23650},
			wantWin:      75,
			newBalance:   23650,
			progressive:  []int64{25, 75},
		},
		{
			name:         "all + all results",
			sessionID:    "yoyo",
			roundID:      "totally deserved",
			startBalance: 23575,
			bet:          10,
			results:      results.Results{r0, r2, r3, r0, r4, r5, r1},
			balances:     []int64{23575, 23575, 23675, 23720, 23720, 23770, 23870, 23970},
			wantWin:      395,
			newBalance:   23970,
			progressive:  []int64{0, 100, 145, 145, 195, 295, 395},
		},
		{
			name:         "all + all results + max payout",
			sessionID:    "baba",
			roundID:      "not bad",
			startBalance: 23575,
			bet:          10,
			results:      results.Results{r0, r2, r3, r0, r4, r5, r1},
			maxPayout:    3000.0,
			balances:     []int64{23575, 23575, 23675, 23720, 23720, 23770, 23870, 23970},
			wantWin:      30000,
			newBalance:   53575,
			progressive:  []int64{0, 100, 145, 145, 195, 295, 395},
		},
	}

	v := &dm{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			params := RoundParams{
				SessionID:    tc.sessionID,
				RoundID:      tc.roundID,
				StartBalance: tc.startBalance,
				Bet:          tc.bet,
				TotalWin:     tc.win,
				Results:      tc.results,
				GameState:    tc.state,
				MaxPayout:    tc.maxPayout,
			}

			r := AcquireRound(v, params)
			require.NotNil(t, r)
			defer r.Release()

			assert.Equal(t, tc.sessionID, r.SessionID())
			assert.Equal(t, tc.roundID, r.RoundID())
			assert.Equal(t, tc.startBalance, r.StartBalance())
			assert.Equal(t, tc.bet, r.Bet())
			assert.Equal(t, tc.state, r.GameState())
			assert.Equal(t, tc.maxPayout, r.maxPayout)
			assert.False(t, r.IsValid())

			n := r.Validate(false, false)
			assert.Equal(t, r, n)
			assert.True(t, r.IsValid())

			if tc.win > 0 {
				assert.Equal(t, tc.win, r.TotalWin())
				assert.Empty(t, r.results)
			} else {
				if tc.results != nil {
					assert.EqualValues(t, tc.results, r.Results())
				} else {
					assert.Empty(t, r.results)
				}
			}

			if len(tc.balances) == 0 {
				assert.Zero(t, len(r.roundResults))
			} else {
				if len(tc.results) > 0 {
					for ix := range tc.results {
						before, after := r.Balances(ix)
						assert.EqualValues(t, tc.balances[ix], before, ix)
						assert.EqualValues(t, tc.balances[ix+1], after, ix)
					}
				}
			}

			assert.Equal(t, tc.wantWin, r.TotalWin())
			assert.Equal(t, tc.newBalance-tc.wantWin, r.StartBalance())
			assert.Equal(t, tc.newBalance, r.NewBalance())

			if len(tc.progressive) > 0 {
				for ix := range tc.progressive {
					assert.Equal(t, tc.progressive[ix], r.ProgressiveWin(ix), ix)
				}
			}
		})
	}
}

var (
	p1 = slots.WinlinePayoutFromData(10, 1, 4, 3, 0, 0, nil)
	p2 = slots.WinlinePayoutFromData(10, 1, 3, 4, 0, 0, nil)
	p3 = slots.WinlinePayoutFromData(4.5, 1, 2, 2, 0, 0, nil)
	p4 = slots.WinlinePayoutFromData(5, 1, 1, 5, 0, 0, nil)
	p5 = slots.WinlinePayoutFromData(10, 1, 8, 3, 0, 0, nil)

	r0, r1, r2, r3, r4, r5 *results.Result
)

func init() {
	r0 = results.AcquireResult(nil, 0)
	r1 = results.AcquireResult(nil, 0, p1)
	r2 = results.AcquireResult(nil, 0, p2)
	r3 = results.AcquireResult(nil, 0, p3)
	r4 = results.AcquireResult(nil, 0, p4)
	r5 = results.AcquireResult(nil, 0, p5)
}

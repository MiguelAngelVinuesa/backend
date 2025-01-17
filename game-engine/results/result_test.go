package results

import (
	"testing"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/encode/zjson"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewResult(t *testing.T) {
	testCases := []struct {
		name    string
		payouts Payouts
		total   float64
		j       string
		j2      string
	}{
		{
			name:    "none",
			payouts: Payouts{},
			j:       `{"freeGames":0,"total":0,"payouts":[]}`,
			j2:      `{"freeGames":1,"total":0,"payouts":[]}`,
		},
		{
			name:    "one",
			payouts: Payouts{p2.Clone().(Payout)},
			total:   2.5,
			j:       `{"freeGames":0,"total":2.5,"payouts":[{"kind":10,"payout":2.5}]}`,
			j2:      `{"freeGames":1,"total":2.5,"payouts":[{"kind":10,"payout":2.5}]}`,
		},
		{
			name:    "two",
			payouts: Payouts{p1.Clone().(Payout), p3.Clone().(Payout)},
			total:   2,
			j:       `{"freeGames":0,"total":2,"payouts":[{"kind":10,"payout":1.5},{"kind":10,"payout":0.5}]}`,
			j2:      `{"freeGames":1,"total":2,"payouts":[{"kind":10,"payout":1.5},{"kind":10,"payout":0.5}]}`,
		},
		{
			name:    "three",
			payouts: Payouts{p4.Clone().(Payout), p5.Clone().(Payout), p6.Clone().(Payout)},
			total:   4.5,
			j:       `{"freeGames":0,"total":4.5,"payouts":[{"kind":10,"payout":1.5},{"kind":10,"payout":2.5},{"kind":10,"payout":0.5}]}`,
			j2:      `{"freeGames":1,"total":4.5,"payouts":[{"kind":10,"payout":1.5},{"kind":10,"payout":2.5},{"kind":10,"payout":0.5}]}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := AcquireResult(nil, 0, tc.payouts...)
			require.NotNil(t, r)
			defer r.Release()

			assert.Nil(t, r.Data)
			assert.Nil(t, r.State)
			assert.EqualValues(t, tc.payouts, r.Payouts)
			assert.Equal(t, tc.total, r.Total)
			assert.Zero(t, r.FreeGames)

			enc := zjson.AcquireEncoder(512)
			defer enc.Release()
			enc.Object(r)
			assert.Equal(t, tc.j, string(enc.Bytes()))

			r.SetFreeGames(5)
			assert.Equal(t, uint64(5), r.FreeGames)
			r.SetFreeGames(7)
			assert.Equal(t, uint64(7), r.FreeGames)
			r.SetFreeGames(1)
			assert.Equal(t, uint64(1), r.FreeGames)

			enc.Reset()
			enc.Object(r)
			assert.Equal(t, tc.j2, string(enc.Bytes()))
		})
	}
}

// func TestResult_ClonePayouts(t *testing.T) {
// 	testCases := []struct {
// 		name    string
// 		payouts Payouts
// 		total   float64
// 		json    string
// 	}{
// 		{
// 			name:    "none",
// 			payouts: Payouts{},
// 			json:    `{"freeGames":0,"total":0,"payouts":[]}`,
// 		},
// 		{
// 			name:    "one",
// 			payouts: Payouts{p2},
// 			total:   2.5,
// 			json:    `{"freeGames":0,"total":2.5,"payouts":[{"kind":10,"payout":2.5}]}`,
// 		},
// 		{
// 			name:    "two",
// 			payouts: Payouts{p1, p3},
// 			total:   2,
// 			json:    `{"freeGames":0,"total":2,"payouts":[{"kind":10,"payout":1.5},{"kind":10,"payout":0.5}]}`,
// 		},
// 		{
// 			name:    "three",
// 			payouts: Payouts{p4, p5, p6},
// 			total:   4.5,
// 			json:    `{"freeGames":0,"total":4.5,"payouts":[{"kind":10,"payout":1.5},{"kind":10,"payout":2.5},{"kind":10,"payout":0.5}]}`,
// 		},
// 	}
//
// 	for _, tc := range testCases {
// 		t.Run(tc.name, func(t *testing.T) {
// 			r := AcquireResult(nil)
// 			require.NotNil(t, r)
// 			defer r.Release()
//
// 			assert.Nil(t, r.Data)
// 			assert.Nil(t, r.State)
// 			assert.Zero(t, len(r.Payouts))
// 			assert.Zero(t, r.Total)
// 			assert.Zero(t, r.FreeGames)
//
// 			for ix := range tc.payouts {
// 				r.ClonePayouts(0, tc.payouts[ix])
// 			}
//
// 			assert.EqualValues(t, tc.payouts, r.Payouts)
// 			assert.Equal(t, tc.total, r.Total)
// 			assert.Zero(t, r.FreeGames)
//
// 			enc := zjson.AcquireEncoder(256)
// 			defer enc.Release()
// 			enc.Object(r)
// 			assert.Equal(t, tc.json, string(enc.Bytes()))
// 		})
// 	}
// }

func TestResult_AddAnimations(t *testing.T) {
	t.Run("add animations", func(t *testing.T) {
		r := AcquireResult(nil, 0)
		require.NotNil(t, r)
		defer r.Release()

		e1 := dummyAnimationProducer.Acquire().(Animator)
		e2 := dummyAnimationProducer.Acquire().(Animator)
		e3 := dummyAnimationProducer.Acquire().(Animator)
		e4 := dummyAnimationProducer.Acquire().(Animator)
		e5 := dummyAnimationProducer.Acquire().(Animator)
		e6 := dummyAnimationProducer.Acquire().(Animator)

		r.AddAnimations(e1)
		assert.Equal(t, 1, len(r.Animations))

		r.AddAnimations(e2, e3)
		assert.Equal(t, 3, len(r.Animations))

		r.AddAnimations(e4, e5, e6)
		assert.Equal(t, 6, len(r.Animations))

		enc := zjson.AcquireEncoder(1024)
		defer enc.Release()

		enc.Object(r)
		got := string(enc.Bytes())
		want := `{"freeGames":0,"total":0,"payouts":[],"animations":[{"kind":99},{"kind":99},{"kind":99},{"kind":99},{"kind":99},{"kind":99}]}`
		assert.Equal(t, want, got)
	})
}

func TestNewResultMany(t *testing.T) {
	t.Run("new results many", func(t *testing.T) {
		r := AcquireResult(nil, 0,
			p1.Clone().(Payout), p2.Clone().(Payout), p3.Clone().(Payout), p4.Clone().(Payout), p5.Clone().(Payout),
			p6.Clone().(Payout), p7.Clone().(Payout), p8.Clone().(Payout), p9.Clone().(Payout), p10.Clone().(Payout))
		require.NotNil(t, r)
		defer r.Release()

		assert.Nil(t, r.Data)
		assert.Nil(t, r.State)
		assert.EqualValues(t, 10, len(r.Payouts))
		assert.Equal(t, 14.0, r.Total)
		assert.Zero(t, r.FreeGames)
	})
}

func TestPurgeResults(t *testing.T) {
	r1 := AcquireResult(nil, 0)
	r2 := AcquireResult(nil, 0)
	r3 := AcquireResult(nil, 0)
	r4 := AcquireResult(nil, 0)
	r5 := AcquireResult(nil, 0)
	r6 := AcquireResult(nil, 0)
	r7 := AcquireResult(nil, 0)
	r8 := AcquireResult(nil, 0)
	r9 := AcquireResult(nil, 0)
	r10 := AcquireResult(nil, 0)
	r11 := AcquireResult(nil, 0)
	r12 := AcquireResult(nil, 0)
	r13 := AcquireResult(nil, 0)
	r14 := AcquireResult(nil, 0)

	defer func() {
		r1.Release()
		r2.Release()
		r3.Release()
		r4.Release()
		r5.Release()
		r6.Release()
		r7.Release()
		r8.Release()
		r9.Release()
		r10.Release()
		r11.Release()
		r12.Release()
		r13.Release()
		r14.Release()
	}()

	testCases := []struct {
		name    string
		results Results
		cap     int
		want    int
	}{
		{name: "empty", results: Results{}, cap: 4, want: 16},
		{name: "single", results: Results{r1}, cap: 4, want: 16},
		{name: "four", results: Results{r1, r2, r3, r4}, cap: 20, want: 32},
		{name: "six", results: Results{r1, r2, r3, r4, r5, r6}, cap: 33, want: 48},
		{name: "eight", results: Results{r1, r2, r3, r4, r5, r6, r7, r8}, cap: 6, want: 8},
		{name: "all", results: Results{r1, r2, r3, r4, r5, r6, r7, r8, r9, r10, r11, r12, r13, r14}, cap: 10, want: 14},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			results := PurgeResults(tc.results, tc.cap)
			require.NotNil(t, results)
			assert.Equal(t, tc.want, cap(results))
		})
	}
}

func TestGrandTotal(t *testing.T) {
	t.Run("new results many", func(t *testing.T) {
		r1 := AcquireResult(nil, 0, p1.Clone().(Payout), p2.Clone().(Payout), p3.Clone().(Payout))
		require.NotNil(t, r1)
		defer r1.Release()

		r2 := AcquireResult(nil, 0, p4.Clone().(Payout), p5.Clone().(Payout))
		require.NotNil(t, r2)
		defer r2.Release()

		r3 := AcquireResult(nil, 0, p6.Clone().(Payout))
		require.NotNil(t, r3)
		defer r3.Release()

		r4 := AcquireResult(nil, 0, p7.Clone().(Payout), p8.Clone().(Payout), p9.Clone().(Payout), p10.Clone().(Payout))
		require.NotNil(t, r3)
		defer r4.Release()

		gt := GrandTotal(Results{r1, r2, r3, r4})
		assert.Equal(t, 14.0, gt)
	})
}

var (
	p1  = AcquirePlayerChoice(1.5)
	p2  = AcquirePlayerChoice(2.5)
	p3  = AcquirePlayerChoice(0.5)
	p4  = AcquirePlayerChoice(1.5)
	p5  = AcquirePlayerChoice(2.5)
	p6  = AcquirePlayerChoice(0.5)
	p7  = AcquirePlayerChoice(0.5)
	p8  = AcquirePlayerChoice(1.5)
	p9  = AcquirePlayerChoice(2.5)
	p10 = AcquirePlayerChoice(0.5)
)

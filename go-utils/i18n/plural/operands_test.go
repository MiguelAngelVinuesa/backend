package plural

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseInt64(t *testing.T) {
	testCases := []struct {
		name  string
		value int64
		wantN float64
		wantI int64
		wantV int64
		wantW int64
		wantF int64
		wantT int64
	}{
		{name: "1", value: 1, wantN: 1, wantI: 1},
		{name: "1200000", value: 1200000, wantN: 1200000, wantI: 1200000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotN := Int64N(tc.value)
			gotI := Int64I(tc.value)
			gotV := Int64V(tc.value)
			gotW := Int64W(tc.value)
			gotF := Int64F(tc.value)
			gotT := Int64T(tc.value)

			assert.Equal(t, tc.wantN, gotN)
			assert.Equal(t, tc.wantI, gotI)
			assert.Equal(t, tc.wantV, gotV)
			assert.Equal(t, tc.wantW, gotW)
			assert.Equal(t, tc.wantF, gotF)
			assert.Equal(t, tc.wantT, gotT)
		})
	}
}

func TestBaseFloat64(t *testing.T) {
	testCases := []struct {
		name  string
		value float64
		wantN float64
		wantI int64
		wantV int64
		wantW int64
		wantF int64
		wantT int64
	}{
		{name: "1", value: 1, wantN: 1, wantI: 1},
		{name: "1.0", value: 1.0, wantN: 1, wantI: 1},
		{name: "1.00", value: 1.00, wantN: 1, wantI: 1},
		{name: "1.3", value: 1.3, wantN: 1.3, wantI: 1, wantV: 1, wantW: 1, wantF: 3, wantT: 3},
		{name: "1.30", value: 1.30, wantN: 1.3, wantI: 1, wantV: 1, wantW: 1, wantF: 3, wantT: 3},
		{name: "1.03", value: 1.03, wantN: 1.03, wantI: 1, wantV: 2, wantW: 2, wantF: 3, wantT: 3},
		{name: "1.230", value: 1.230, wantN: 1.23, wantI: 1, wantV: 2, wantW: 2, wantF: 23, wantT: 23},
		{name: "1200000", value: 1200000, wantN: 1200000, wantI: 1200000},
		{name: "1200.50", value: 1200.50, wantN: 1200.5, wantI: 1200, wantV: 1, wantW: 1, wantF: 5, wantT: 5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotN := FloatN(tc.value)
			gotI := FloatI(tc.value)
			gotV := FloatV(tc.value)
			gotW := FloatW(tc.value)
			gotF := FloatF(tc.value)
			gotT := FloatT(tc.value)

			assert.Equal(t, tc.wantN, gotN)
			assert.Equal(t, tc.wantI, gotI)
			assert.Equal(t, tc.wantV, gotV)
			assert.Equal(t, tc.wantW, gotW)
			assert.Equal(t, tc.wantF, gotF)
			assert.Equal(t, tc.wantT, gotT)
		})
	}
}

func TestBaseString(t *testing.T) {
	testCases := []struct {
		name  string
		value string
		wantN float64
		wantI int64
		wantV int64
		wantW int64
		wantF int64
		wantT int64
	}{
		{name: "1", value: "1", wantN: 1, wantI: 1},
		{name: "1.0", value: "1.0", wantN: 1, wantI: 1, wantV: 1},
		{name: "1.00", value: "1.00", wantN: 1, wantI: 1, wantV: 2},
		{name: "1.3", value: "1.3", wantN: 1.3, wantI: 1, wantV: 1, wantW: 1, wantF: 3, wantT: 3},
		{name: "1.30", value: "1.30", wantN: 1.3, wantI: 1, wantV: 2, wantW: 1, wantF: 30, wantT: 3},
		{name: "1.03", value: "1.03", wantN: 1.03, wantI: 1, wantV: 2, wantW: 2, wantF: 3, wantT: 3},
		{name: "1.230", value: "1.230", wantN: 1.23, wantI: 1, wantV: 3, wantW: 2, wantF: 230, wantT: 23},
		{name: "1200000", value: "1200000", wantN: 1200000, wantI: 1200000},
		{name: "1200.50", value: "1200.50", wantN: 1200.5, wantI: 1200, wantV: 2, wantW: 1, wantF: 50, wantT: 5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotN := StringN(tc.value)
			gotI := StringI(tc.value)
			gotV := StringV(tc.value)
			gotW := StringW(tc.value)
			gotF := StringF(tc.value)
			gotT := StringT(tc.value)

			assert.Equal(t, tc.wantN, gotN)
			assert.Equal(t, tc.wantI, gotI)
			assert.Equal(t, tc.wantV, gotV)
			assert.Equal(t, tc.wantW, gotW)
			assert.Equal(t, tc.wantF, gotF)
			assert.Equal(t, tc.wantT, gotT)
		})
	}
}

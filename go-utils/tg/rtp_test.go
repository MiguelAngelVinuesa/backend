package tg

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyRTP(t *testing.T) {
	testCases := []struct {
		name string
		rtp  uint8
		fail bool
	}{
		{name: "0", fail: true},
		{name: "90", rtp: 90, fail: true},
		{name: "91", rtp: 91, fail: true},
		{name: "92", rtp: 92},
		{name: "93", rtp: 93, fail: true},
		{name: "94", rtp: 94},
		{name: "95", rtp: 95, fail: true},
		{name: "96", rtp: 96},
		{name: "97", rtp: 97, fail: true},
		{name: "98", rtp: 98, fail: true},
		{name: "99", rtp: 99, fail: true},
		{name: "42", rtp: 42},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := VerifyRTP(tc.rtp)
			if tc.fail {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

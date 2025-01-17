package net

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitHostPort(t *testing.T) {
	testCases := []struct {
		name    string
		address string
		host    string
		port    int
	}{
		{name: "empty"},
		{name: "host only", address: "topgaming.team", host: "topgaming.team"},
		{name: "port only", address: ":8000", port: 8000},
		{name: "host+port", address: "topgaming.team:8000", host: "topgaming.team", port: 8000},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			host, port := SplitHostPort(tc.address)
			assert.Equal(t, tc.host, host)
			assert.Equal(t, tc.port, port)
		})
	}
}

func TestBuildAddress(t *testing.T) {
	testCases := []struct {
		name string
		host string
		port int
		want string
	}{
		{name: "empty"},
		{name: "host only", host: "xs4all.nl", want: "xs4all.nl"},
		{name: "port only", port: 8000},
		{name: "host+port", host: "help.me", port: 7777, want: "help.me:7777"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := BuildAddress(tc.host, tc.port)
			assert.Equal(t, tc.want, got)
		})
	}
}

package dns

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	testCases := []struct {
		name        string
		resolver    string
		timeout     time.Duration
		addresses   []string
		wantTimeout time.Duration
		failLookup  bool
	}{
		{
			name:        "xs4all.nl, no timeout",
			resolver:    "1.1.1.1:53",
			addresses:   []string{"xs4all.nl"},
			wantTimeout: DefaultResolveTimeout,
		},
		{
			name:        "gw.dev.topgaming.team, timeout",
			resolver:    "1.1.1.1:53",
			timeout:     5 * time.Second,
			addresses:   []string{"gw.dev.topgaming.team"},
			wantTimeout: 5 * time.Second,
		},
		{
			name:        "multiple, no timeout",
			resolver:    "1.1.1.1:53",
			addresses:   []string{"xs4all.nl", "google.com:80", "aws.com:443", "gw.dev.topgaming.team"},
			wantTimeout: DefaultResolveTimeout,
		},
		{
			name:        "fail",
			resolver:    "1.1.1.1:53",
			addresses:   []string{"i.dont.exist.now"},
			wantTimeout: DefaultResolveTimeout,
			failLookup:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			d := net.Dialer{Timeout: tc.timeout}
			r := &net.Resolver{
				PreferGo: true,
				Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
					return d.DialContext(ctx, DefaultResolveProtocol, tc.resolver)
				},
			}

			c := New(r, tc.timeout, 0, tc.addresses...)
			require.NotNil(t, c)
			assert.Equal(t, r, c.resolver)
			assert.Equal(t, tc.wantTimeout, c.timeout)

			require.Nil(t, c.cache["haha"])

			if tc.failLookup {
				for _, address := range tc.addresses {
					l := c.cache[address]
					require.NotNil(t, l)
					assert.Zero(t, len(l.ips))
					assert.Zero(t, l.max)
					assert.Zero(t, l.last)
					assert.Empty(t, c.NextIP(address))
				}
			} else {
				for _, address := range tc.addresses {
					l := c.cache[address]
					require.NotNil(t, l)
					assert.NotZero(t, len(l.ips))
					assert.Equal(t, l.max, len(l.ips))
					assert.GreaterOrEqual(t, l.last, 0)
					assert.LessOrEqual(t, l.last, l.max)

					counts := make(map[string]int)
					for ix := 0; ix < l.max; ix++ {
						a := c.NextIP(address)
						counts[a] = counts[a] + 1
					}
					assert.Equal(t, l.max, len(counts))
				}
			}
		})
	}
}

func TestCache_Refresh(t *testing.T) {
	t.Run("refresh", func(t *testing.T) {
		d := net.Dialer{Timeout: DefaultResolveTimeout}
		r := &net.Resolver{
			PreferGo: true,
			Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				return d.DialContext(ctx, DefaultResolveProtocol, "1.1.1.1:53")
			},
		}

		address := "xs4all.nl"
		c := New(r, DefaultResolveTimeout, 0, address)
		require.NotNil(t, c)

		c.Refresh(address)

		l := c.cache[address]
		require.NotNil(t, l)
		assert.NotZero(t, len(l.ips))
		assert.Equal(t, l.max, len(l.ips))
		assert.GreaterOrEqual(t, l.last, 0)
		assert.LessOrEqual(t, l.last, l.max)

		counts := make(map[string]int)
		for ix := 0; ix < l.max; ix++ {
			a := c.NextIP(address)
			counts[a] = counts[a] + 1
		}
		assert.Equal(t, l.max, len(counts))
	})
}

package dns

import (
	"context"
	"net"
	"strings"
	"time"
)

const (
	DefaultResolveTimeout  = 2 * time.Second
	DefaultResolveProtocol = "tcp"
)

// NewDialerWithResolver returns a dial function with a custom DNS resolver.
// If the given resolver address is empty it does nothing.
// If the given timeout is zero the default timeout will be used.
func NewDialerWithResolver(resolver, protocol string, timeout time.Duration) func(context.Context, string, string) (net.Conn, error) {
	if resolver == "" {
		return nil
	}

	if protocol == "" {
		protocol = DefaultResolveProtocol
	}
	if timeout == 0 {
		timeout = DefaultResolveTimeout
	}

	d1 := net.Dialer{Timeout: timeout}
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return d1.DialContext(ctx, protocol, resolver)
		},
	}

	d2 := &net.Dialer{Resolver: r}
	return d2.DialContext
}

// NewDialerWithCache returns a dial function with a DNS cache and custom resolver.
// If the given resolver address is empty the function does nothing.
// If the given timeout is zero the default timeout will be used.
func NewDialerWithCache(resolver, protocol string, timeout, autoRefresh time.Duration, addresses ...string) func(context.Context, string, string) (net.Conn, error) {
	if resolver == "" {
		return nil
	}
	_, dialer := NewCacheAndDialer(resolver, protocol, timeout, autoRefresh, addresses)
	return dialer
}

func NewCacheAndDialer(resolver, protocol string, timeout, autoRefresh time.Duration, addresses []string) (*Cache, func(context.Context, string, string) (net.Conn, error)) {
	if protocol == "" {
		protocol = DefaultResolveProtocol
	}
	if timeout == 0 {
		timeout = DefaultResolveTimeout
	}

	// Remove scheme from addresses.
	for i := 0; i < len(addresses); i++ {
		addresses[i] = strings.Replace(addresses[i], "http://", "", 1)
		addresses[i] = strings.Replace(addresses[i], "https://", "", 1)
	}

	d1 := net.Dialer{Timeout: timeout}
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			return d1.DialContext(ctx, protocol, resolver)
		},
	}

	c := New(r, timeout, autoRefresh, addresses...)

	d2 := &net.Dialer{Resolver: r}
	return c, func(ctx context.Context, network, address string) (net.Conn, error) {
		if ip := c.NextIP(address); ip != "" {
			if conn, err := d2.DialContext(ctx, network, ip); err == nil {
				return conn, nil
			}
			c.Refresh(address)
		}
		return d2.DialContext(ctx, network, address)
	}
}

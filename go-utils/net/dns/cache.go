package dns

import (
	"context"
	"math/rand"
	"net"
	"sync"
	"time"

	"github.com/goccy/go-json"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/log"
	addr "git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/net"
)

// Cache represents an in-memory DNS lookup cache.
type Cache struct {
	mu       sync.Mutex
	resolver *net.Resolver
	timeout  time.Duration
	cache    map[string]*lookup
	logger   log.Logger
}

// New instantiates a new DNS lookup cache for the given addresses.
// addresses can consist of just the host name or the host name plus a port; the cache will handle both options fine.
// autoRefresh represents the interval after which the IP addresses will be refreshed automatically.
func New(resolver *net.Resolver, timeout, autoRefresh time.Duration, addresses ...string) *Cache {
	if timeout == 0 {
		timeout = DefaultResolveTimeout
	}

	c := &Cache{
		resolver: resolver,
		timeout:  timeout,
		cache:    make(map[string]*lookup),
	}

	for _, address := range addresses {
		c.cache[address] = newLookup(c, address)
	}

	if autoRefresh != 0 {
		go func(interval time.Duration) {
			ticker := time.NewTicker(interval)
			for {
				select {
				case <-ticker.C:
					c.mu.Lock()
					for _, l := range c.cache {
						go l.refresh()
					}
					c.mu.Unlock()
				}
			}
		}(autoRefresh)
	}

	return c
}

// NextIP retrieves the next IP from a list of cached IP addresses for the given address.
// If the host doesn't exist, or there are no IP addresses for the host, the function returns the empty string.
func (c *Cache) NextIP(address string) string {
	c.mu.Lock()
	l := c.cache[address]
	c.mu.Unlock()

	if l == nil {
		return ""
	}

	next := l.next()
	_, port := addr.SplitHostPort(address)
	if port == 0 {
		return next
	}
	return addr.BuildAddress(next, port)
}

// Refresh refreshes the cached IP addresses for the given address.
// It should be called if a dialer cannot connect to one of the cached IP addresses.
// The function does nothing if the address was not cached previously.
func (c *Cache) Refresh(address string) {
	c.mu.Lock()
	l := c.cache[address]
	c.mu.Unlock()

	if l != nil {
		go l.refresh()
	}
}

// Invalidate invalidates an IP from the cache if the connection failed because the service died.
// If it is the last IP, the cache will automatically refresh.
func (c *Cache) Invalidate(address string, ip string) {
	c.mu.Lock()
	l := c.cache[address]
	c.mu.Unlock()

	if l != nil {
		l.invalidate(ip)
	}
}

// SetLogger sets the logging interface for the cache.
func (c *Cache) SetLogger(logger log.Logger) {
	c.logger = logger
}

// LogAll logs the current cache values.
func (c *Cache) LogAll() {
	if c.logger == nil {
		return
	}

	c.mu.Lock()
	c.logger.Debug("resolver", "cached", c.cache)
	c.mu.Unlock()
}

func newLookup(c *Cache, address string) *lookup {
	host, _ := addr.SplitHostPort(address)
	l := &lookup{c: c, host: host, ips: make([]string, 0)}
	l.refresh()
	return l
}

type lookup struct {
	mu   sync.RWMutex
	c    *Cache
	host string
	ips  []string
	max  int
	last int
}

func (l *lookup) refresh() {
	l.mu.Lock()
	l.ips = l.ips[:0]
	l.max = 0
	l.last = 0
	l.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), l.c.timeout)
	list, err := l.c.resolver.LookupIPAddr(ctx, l.host)
	cancel()
	if err != nil {
		return
	}

	l.mu.Lock()
	l.ips = l.ips[:0]
	for _, ip := range list {
		l.ips = append(l.ips, ip.String())
	}
	l.max = len(l.ips)
	l.last = rand.Intn(l.max)
	l.mu.Unlock()

	if l.c.logger != nil {
		l.mu.RLock()
		l.c.logger.Debug("resolver cache refreshed", "address", l.host, "ips", l.ips)
		l.mu.RUnlock()
	}
}

func (l *lookup) next() string {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if l.max == 0 {
		return ""
	}

	if l.last++; l.last >= l.max {
		l.last = 0
	}
	return l.ips[l.last]
}

func (l *lookup) invalidate(ip string) {
	l.mu.Lock()
	max := len(l.ips)
	for ix, old := range l.ips {
		if old == ip {
			switch ix {
			case 0:
				l.ips = l.ips[1:]
			case max - 1:
				l.ips = l.ips[:max-1]
			default:
				l.ips = append(l.ips[:ix], l.ips[ix+1:]...)
			}
			max--
			break
		}
	}
	l.mu.Unlock()

	if l.c.logger != nil {
		l.mu.RLock()
		l.c.logger.Debug("resolver ip invalidated", "address", l.host, "ip", ip)
		l.mu.RUnlock()
	}

	if max == 0 {
		l.refresh()
	}
}

// MarshalJSON implements the JSON marshaller interface.
func (l *lookup) MarshalJSON() ([]byte, error) {
	return json.Marshal(&jsonLookup{Ips: l.ips})
}

type jsonLookup struct {
	Ips []string `json:"ips"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

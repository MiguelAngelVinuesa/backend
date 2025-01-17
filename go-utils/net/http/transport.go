package http

import (
	"io"
	"net"
	"net/http"
	"net/url"
	"time"

	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/log"
	"git-codecommit.eu-central-1.amazonaws.com/v1/repos/go-utils.git/net/dns"
)

// NewRetryTransport instantiates a http.Transport (RoundTripper) with a retry mechanism and optional custom DNS resolver.
func NewRetryTransport(maxRetries int, retryBackoff time.Duration, opts ...RetryOption) http.RoundTripper {
	if maxRetries < 1 {
		maxRetries = 10
	}

	r := &retryTransport{
		maxIdleConn:    100,
		maxConnPerHost: 100,
		maxIdlePerHost: 100,
		maxRetries:     maxRetries,
		retryBackoff:   retryBackoff,
	}

	for ix := range opts {
		opts[ix](r)
	}

	r.transport = http.DefaultTransport.(*http.Transport).Clone()
	r.transport.MaxIdleConns = r.maxIdleConn
	r.transport.MaxConnsPerHost = r.maxConnPerHost
	r.transport.MaxIdleConnsPerHost = r.maxIdlePerHost

	if r.resolver != "" {
		r.cache, r.transport.DialContext = dns.NewCacheAndDialer(r.resolver, r.protocol, r.timeout, r.refresh, r.addresses)
		if r.logger != nil {
			r.cache.SetLogger(r.logger)
			r.cache.LogAll()
		}
	}

	return r
}

// RetryOption is the function signature for RetryTransport options.
type RetryOption func(r *retryTransport)

// RetryResolver is the RetryTransport option to set up a custom cached DNS resolver.
func RetryResolver(resolver, protocol string, timeout, refresh time.Duration, addresses ...string) RetryOption {
	return func(r *retryTransport) {
		r.resolver = resolver
		r.protocol = protocol
		r.timeout = timeout
		r.refresh = refresh
		r.addresses = addresses
	}
}

// RetryIdleConn is the RetryTransport option to set the idle connections limits.
func RetryIdleConn(maxIdleConn, maxConnPerHost, maxIdlePerHost int) RetryOption {
	return func(r *retryTransport) {
		r.maxIdleConn = maxIdleConn
		r.maxConnPerHost = maxConnPerHost
		r.maxIdlePerHost = maxIdlePerHost
	}
}

// RetryLogger is the RetryTransport option to set the logger.
func RetryLogger(logger log.Logger) RetryOption {
	return func(r *retryTransport) {
		r.logger = logger
	}
}

type retryTransport struct {
	maxIdleConn    int
	maxConnPerHost int
	maxIdlePerHost int
	maxRetries     int
	retryBackoff   time.Duration
	resolver       string
	protocol       string
	timeout        time.Duration
	refresh        time.Duration
	addresses      []string
	logger         log.Logger
	cache          *dns.Cache
	dialer         *net.Dialer
	transport      *http.Transport
}

// RoundTrip implements the RoundTripper interface.
func (r *retryTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	retries := r.maxRetries
	backoff := r.retryBackoff

	var resp *http.Response
	var err error
	for {
		if resp, err = r.transport.RoundTrip(req); err == nil || retries <= 0 {
			break
		}

		// not "EOF/request cancelled/context deadline" means a serious error, so we can't retry!
		fatal := err != io.EOF && err.Error() != "net/http: request canceled" && err.Error() != "context deadline exceeded"

		if r.logger != nil {
			j := &jsonRequest{
				Method:        req.Method,
				Proto:         req.Proto,
				Host:          req.Host,
				Close:         req.Close,
				ContentLength: req.ContentLength,
				URL:           req.URL,
				Header:        req.Header,
			}
			r.logger.Error("retryTransport round trip failed", "request", j, "error", err, "fatal", fatal)
		}

		if !fatal {
			break
		}

		retries--
		time.Sleep(backoff)
		backoff *= 2
	}

	return resp, err
}

type jsonRequest struct {
	Method        string      `json:"method"`
	Proto         string      `json:"proto"`
	Host          string      `json:"host"`
	Close         bool        `json:"close"`
	ContentLength int64       `json:"contentLength"`
	URL           *url.URL    `json:"URL"`
	Header        http.Header `json:"header"`
}

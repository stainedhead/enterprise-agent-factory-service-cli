package client

import (
	"context"
	"net"
	"net/http"
	"time"
)

// Default timeouts (NFR reliability).
const (
	DefaultConnectTimeout = 30 * time.Second
	DefaultReadTimeout    = 60 * time.Second
)

// TransportOptions configure NewTransport.
type TransportOptions struct {
	// ConnectTimeout bounds dial plus TLS handshake; 0 disables it.
	ConnectTimeout time.Duration
	// DialContext replaces the net dialer (tests).
	DialContext func(ctx context.Context, network, addr string) (net.Conn, error)
}

// NewTransport returns the HTTP transport: proxies from the environment, a connect timeout.
func NewTransport(o TransportOptions) *http.Transport {
	dial := o.DialContext
	if dial == nil {
		d := &net.Dialer{KeepAlive: 30 * time.Second}
		dial = d.DialContext
	}
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.TLSHandshakeTimeout = o.ConnectTimeout
	t.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		if o.ConnectTimeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(ctx, o.ConnectTimeout)
			defer cancel()
		}
		return dial(ctx, network, addr)
	}
	return t
}

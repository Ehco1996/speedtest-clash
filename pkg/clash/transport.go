package clash

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/Dreamacro/clash/constant"
)

type ClashTransport struct {
	proxy constant.Proxy
	tp    *http.Transport

	// change with every request
	currentURL *url.URL
}

func NewClashTransport(p constant.Proxy) *ClashTransport {
	c := &ClashTransport{proxy: p}

	tp := &http.Transport{
		// from http.DefaultTransport
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	tp.DialContext = c.DialContext
	c.tp = tp
	return c
}

func (c *ClashTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	c.currentURL = req.URL
	return c.tp.RoundTrip(req)
}

func (c *ClashTransport) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	meta, err := URLToMetadata(c.currentURL)
	if err != nil {
		return nil, err
	}
	conn, err := c.proxy.DialContext(ctx, meta)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

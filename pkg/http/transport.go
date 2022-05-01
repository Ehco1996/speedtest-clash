package http

import (
	"context"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/Dreamacro/clash/constant"

	"github.com/Ehco1996/clash-speed/pkg/clash"
)

type ClashTransport struct {
	tp http.Transport

	currentURL *url.URL

	proxy constant.Proxy
}

func NewClashTransport() *ClashTransport {
	c := &ClashTransport{}
	tp := http.Transport{
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
	meta, err := clash.URLToMetadata(c.currentURL)
	conn, err := c.proxy.DialContext(ctx, meta)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

package xrpc

import (
	"net"
	"time"
)

const (
	defaultConnectTimeout = 100 * time.Millisecond
	defaultReadTimeout    = 500 * time.Millisecond
	defaultWriteTimeout   = 500 * time.Millisecond
)

type Client struct {
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration

	selector ServerSelector
}

func NewClient(servers []string) *Client {
	ss := &RoundRobinSelector{}
	ss.SetServer(servers)
	return &Client{
		selector: ss,
	}
}

func (c *Client) connectTimeout() time.Duration {
	if c.ConnectTimeout == 0 {
		return defaultConnectTimeout
	}

	return c.ConnectTimeout
}

func (c *Client) readTimeout() time.Duration {
	if c.ReadTimeout == 0 {
		return defaultReadTimeout
	}

	return c.ReadTimeout
}

func (c *Client) writeTimeout() time.Duration {
	if c.WriteTimeout == 0 {
		return defaultWriteTimeout
	}

	return c.WriteTimeout
}

func (c *Client) Do(req *Request) (*Response, error) {
	addr, err := c.selector.PickServer()
	if err != nil {
		return nil, err
	}

	nc, err := net.DialTimeout(addr.Network(), addr.String(), c.ConnectTimeout)
	if err != nil {
		return nil, err
	}

	nc.SetReadDeadline(time.Now().Add(c.readTimeout()))
	nc.SetWriteDeadline(time.Now().Add(c.writeTimeout()))

	return nil, nil
}
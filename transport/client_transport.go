package transport

import (
	"context"
	"fmt"
)

type clientTransport struct {
	opts *ClientTransportOptions
}

var (
	clientTransportMap     = make(map[string]ClientTransport)
	DefaultClientTransport = New()
)

func init() {
	clientTransportMap["default"] = DefaultClientTransport
}

var New = func() ClientTransport {
	return &clientTransport{
		opts: &ClientTransportOptions{},
	}
}

func (c *clientTransport) Send(ctx context.Context, req []byte, opts ...ClientTransportOption) ([]byte, error) {
	for _, o := range opts {
		o(c.opts)
	}
	if c.opts.Network == "tcp" {
		return c.SendTcpReq(ctx, req)
	}
	if c.opts.Network == "udp" {
		return c.sendUdpReq(ctx, req)
	}
	return nil, fmt.Errorf("network type not supported")
}

func (c *clientTransport) SendTcpReq(ctx context.Context, req []byte) ([]byte, error) {

	// service discovery
	addr, err := c.opts.Selector.Select(c.opts.ServiceName)
	if err != nil {
		return nil, err
	}

	// defaultSelector returns "", use the target as address
	if addr == "" {
		addr = c.opts.Target
	}

	conn, err := c.opts.Pool.Get(ctx, c.opts.Network, addr)
	//	conn, err := net.DialTimeout("tcp", addr, c.opts.Timeout);
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	sendNum := 0
	num := 0
	for sendNum < len(req) {
		num, err = conn.Write(req[sendNum:])
		if err != nil {
			return nil, err
		}
		sendNum += num

		if err = isDone(ctx); err != nil {
			return nil, err
		}
	}

	// parse frame

	wrapperConn := wrapConn(conn)
	frame, err := wrapperConn.framer.ReadFrame(conn)
	if err != nil {
		return nil, err
	}

	return frame, err
}

func (c *clientTransport) sendUdpReq(ctx context.Context, bytes []byte) ([]byte, error) {

}

func isDone(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	return nil
}

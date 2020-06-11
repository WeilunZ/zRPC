package transport

import (
	"context"
	"errors"
	"io"
	"net"
	"strings"
	"time"

	"github.com/WeilunZ/zRPC/components/codec"
	"github.com/WeilunZ/zRPC/components/log"
	"github.com/WeilunZ/zRPC/components/protocol"
	"github.com/WeilunZ/zRPC/components/state"
	"github.com/golang/protobuf/proto"
)

type serverTransport struct {
	opts *ServerTransportOptions
}

var serverTransportMap = make(map[string]ServerTransport)

func init() {
	serverTransportMap["default"] = DefaultServerTransport
}

var DefaultServerTransport = NewServerTransport()

var NewServerTransport = func() ServerTransport {
	return &serverTransport{
		opts: &ServerTransportOptions{},
	}
}

func RegisterServerTransport(name string, serverTransport ServerTransport) {
	if serverTransportMap == nil {
		serverTransportMap = make(map[string]ServerTransport)
	}
	serverTransportMap[name] = serverTransport
}

func GetServerTransport(transport string) ServerTransport {
	if v, ok := serverTransportMap[transport]; ok {
		return v
	}
	return DefaultServerTransport
}

func (s *serverTransport) ListenAndServe(ctx context.Context, opts ...ServerTransportOption) error {
	for _, o := range opts {
		o(s.opts)
	}
	if strings.Index(s.opts.Network, "tcp") != -1 {
		return s.ListenAndServeTcp(ctx, opts...)
	}
	return errors.New("network protocol not supported")
}

func (s *serverTransport) ListenAndServeTcp(ctx context.Context, opts ...ServerTransportOption) error {

	lis, err := net.Listen(s.opts.Network, s.opts.Address)
	if err != nil {
		return err
	}

	go func() {
		if err = s.serve(ctx, lis); err != nil {
			log.Errorf("transport serve error, %v", err)
		}
	}()

	return nil
}

func (s *serverTransport) serve(ctx context.Context, lis net.Listener) error {

	var tempDelay time.Duration

	tl, ok := lis.(*net.TCPListener)
	if !ok {
		return errors.New("network not supported")
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		conn, err := tl.AcceptTCP()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				time.Sleep(tempDelay)
				continue
			}
			return err
		}

		if err = conn.SetKeepAlive(true); err != nil {
			return err
		}

		if s.opts.KeepAlivePeriod != 0 {
			_ = conn.SetKeepAlivePeriod(s.opts.KeepAlivePeriod)
		}

		go func() {
			if err := s.handleConn(ctx, wrapConn(conn)); err != nil {
				log.Errorf("gorpc handle tcp conn error, %v", err)
			}
		}()

	}
}

func (s *serverTransport) handleConn(ctx context.Context, conn *connWrapper) error {

	// close the connection before return
	// the connection closes only if a network read or write fails
	defer conn.Close()

	for {
		// check upstream ctx is done
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		frame, err := s.read(ctx, conn)
		if err == io.EOF {
			// read compeleted
			return nil
		}

		if err != nil {
			return err
		}

		rsp, err := s.handle(ctx, frame)
		if err != nil {
			log.Errorf("s.handle err is not nil, %v", err)
		}

		if err = s.write(ctx, conn, rsp); err != nil {
			return err
		}
	}

}

func (s *serverTransport) read(ctx context.Context, conn *connWrapper) ([]byte, error) {
	frame, err := conn.framer.ReadFrame(conn)
	if err != nil {
		return nil, err
	}
	return frame, nil
}

func (s *serverTransport) handle(ctx context.Context, frame []byte) ([]byte, error) {
	cdc := codec.GetCodec(s.opts.Protocol)
	reqb, err := cdc.Decode(frame)
	if err != nil {
		log.Errorf("decode error: %v", err)
		return nil, err
	}
	rspb, err := s.opts.Handler.Handle(ctx, reqb)
	if err != nil {
		log.Errorf("handler error: %v", err)
	}
	response := wrapResponse(rspb, err)
	rspPb, err := proto.Marshal(response)
	if err != nil {
		log.Errorf("proto marshal error: %v", err)
		return nil, err
	}
	rspBody, err := cdc.Encode(rspPb)
	if err != nil {
		log.Errorf("server encode error, response : %v, error: %v", response, err)
		return nil, err
	}
	return rspBody, nil
}

func (s *serverTransport) write(ctx context.Context, conn net.Conn, rsp []byte) error {
	if _, err := conn.Write(rsp); err != nil {
		log.Errorf("conn Write err: %v", err)
	}
	return nil
}

func wrapResponse(payload []byte, err error) *protocol.Response {
	response := &protocol.Response{
		Payload: payload,
		RetCode: state.OK,
		RetMsg:  state.SUCCESS,
	}
	if err != nil {
		if e, ok := err.(*state.Error); ok {
			response.RetCode = e.Code
			response.RetMsg = e.Message
		} else {
			response.RetCode = state.InternalError
			response.RetMsg = state.InternalErrorMessage
		}
	}
	return response
}

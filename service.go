package tinyRPC

import (
	"context"
	"errors"

	"github.com/WeilunZ/myrpc/components/log"

	"github.com/WeilunZ/myrpc/components/utils"

	"github.com/WeilunZ/myrpc/components/codec"
	"github.com/WeilunZ/myrpc/components/protocol"
	"github.com/golang/protobuf/proto"

	"github.com/WeilunZ/myrpc/transport"

	"github.com/WeilunZ/myrpc/components/interceptor"

	"github.com/WeilunZ/myrpc/server"
)

type Service interface {
	Register(string, Handler)
	Serve(*server.ServerOptions)
	Close()
}

type service struct {
	svr         interface{}        // server
	ctx         context.Context    // 每一个 service 一个上下文进行管理
	cancel      context.CancelFunc // context 的控制器
	serviceName string             // 服务名
	handlers    map[string]Handler
	opts        *server.ServerOptions // 参数选项
	closing     bool                  // 服务停止中？
}

type ServiceDesc struct {
	Svr         interface{}
	ServiceName string
	Methods     []*Method
	HandlerType interface{}
}

type Method struct {
	MethodName string
	Handler    Handler
}
type Handler func(interface{}, context.Context, func(interface{}) error, []interceptor.ServerInterceptor) (interface{}, error)

func (s *service) Register(handlerName string, handler Handler) {
	if s.handlers == nil {
		s.handlers = make(map[string]Handler)
	}
	s.handlers[handlerName] = handler
}

func (s *service) Serve(opts *server.ServerOptions) {
	s.opts = opts
	transportOpts := []transport.ServerTransportOption{
		transport.WithServerAddress(s.opts.Address),
		transport.WithServerNetwork(s.opts.Network),
		transport.WithHandler(s),
		transport.WithServerTimeout(s.opts.Timeout),
		transport.WithSerialization(s.opts.SerializationType),
		transport.WithProtocol(s.opts.Protocol),
	}

	serverTransport := transport.GetServerTransport(s.opts.Protocol)

	s.ctx, s.cancel = context.WithCancel(context.Background())

	if err := serverTransport.ListenAndServe(s.ctx, transportOpts...); err != nil {
		log.Errorf("%s serve error, %v", s.opts.Network, err)
		return
	}

	log.Infof("%s service serving started at %s ... \n", s.serviceName, s.opts.Address)
}

func (s *service) Handle(ctx context.Context, reqbuf []byte) ([]byte, error) {
	request := &protocol.Request{}
	if err := proto.Unmarshal(reqbuf, request); err != nil {
		return nil, err
	}

	serverSerialization := codec.GetSerialization(s.opts.SerializationType)

	dec := func(req interface{}) error {
		if err := serverSerialization.Deserialize(request.Payload, req); err != nil {
			return err
		}
		return nil
	}

	if s.opts.Timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, s.opts.Timeout)
		defer cancel()
	}

	_, method, err := utils.ParseServicePath(request.ServicePath)
	if err != nil {
		return nil, errors.New("invalid method")
	}

	handler := s.handlers[method]
	if handler == nil {
		return nil, errors.New("handler unregisterd")
	}
	rsp, err := handler(s.svr, ctx, dec, s.opts.Interceptors)
	if err != nil {
		return nil, err
	}
	rspb, err := serverSerialization.Serialize(rsp)
	if err != nil {
		return nil, err
	}
	return rspb, nil
}

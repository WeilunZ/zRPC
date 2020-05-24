package server

import "context"

type Service interface {
	Register(string, Handler)
	Serve(options *ServerOptions)
	Close()
}

type service struct {
	svr         interface{}        // server
	ctx         context.Context    // 每一个 service 一个上下文进行管理
	cancel      context.CancelFunc // context 的控制器
	serviceName string             // 服务名
	handlers    map[string]Handler
	opts        *ServerOptions // 参数选项
}

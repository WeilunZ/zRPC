package interceptor

import "context"

type ClientInvoker func(ctx context.Context, req, resp interface{}) error

type ServerHandler func(ctx context.Context, req interface{}) (interface{}, error)

type ClientInterceptor func(ctx context.Context, req, resp interface{}, ivk ClientInvoker) error

type ServerInterceptor func(ctx context.Context, req interface{}, handler ServerHandler) (interface{}, error)

func ClientIntercept(ctx context.Context, req, resp interface{}, interceptors []ClientInterceptor, ivk ClientInvoker) error {
	if len(interceptors) == 0 {
		return ivk(ctx, req, resp)
	}
	return interceptors[0](ctx, req, resp, getClientInvoker(0, interceptors, ivk))
}

func getClientInvoker(i int, interceptors []ClientInterceptor, ivk ClientInvoker) ClientInvoker {
	if i == len(interceptors)-1 {
		return ivk
	}
	return func(ctx context.Context, req, resp interface{}) error {
		return interceptors[i+1](ctx, req, resp, getClientInvoker(i+1, interceptors, ivk))
	}
}

func ServerIntercept(ctx context.Context, req interface{}, interceptors []ServerInterceptor, handler ServerHandler) (interface{}, error) {
	if len(interceptors) == 0 {
		return handler(ctx, req)
	}
	return interceptors[0](ctx, req, getServerHandler(0, interceptors, handler))
}

func getServerHandler(i int, interceptors []ServerInterceptor, handler ServerHandler) ServerHandler {
	if i == len(interceptors)-1 {
		return handler
	}
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return interceptors[i+1](ctx, req, getServerHandler(i+1, interceptors, handler))
	}
}

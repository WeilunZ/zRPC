package server

import "time"

type ServerOptions struct {
	address           string //e.g. 127.0.0.1:8080/www.baidu.com
	network           string // e.g. tcp/udp
	protocol          string
	timeout           time.Duration
	serializationType string   // serialization type, default: proto
	selectorSvrAddr   string   // service discovery server address, required when using the third-party service discovery plugin
	tracingSvrAddr    string   // tracing plugin server address, required when using the third-party tracing plugin
	tracingSpanName   string   // tracing span name, required when using the third-party tracing plugin
	pluginNames       []string // plugin name
	interceptors      []interceptor.ServerInterceptor
}

type ServerOption func(*ServerOptions)

func WithAddress(address string) ServerOption {
	return func(o *ServerOptions) {
		o.address = address
	}
}

func WithNetwork(network string) ServerOption {
	return func(o *ServerOptions) {
		o.network = network
	}
}

func WithProtocol(protocol string) ServerOption {
	return func(o *ServerOptions) {
		o.protocol = protocol
	}
}

func WithTimeOut(timeout time.Duration) ServerOption {
	return func(o *ServerOptions) {
		o.timeout = timeout
	}
}

func WithSerializationType(serializationType string) ServerOption {
	return func(o *ServerOptions) {
		o.serializationType = serializationType
	}
}

func WithSelectorSvrAddr(selectorSvrAddr string) ServerOption {
	return func(o *ServerOptions) {
		o.selectorSvrAddr = selectorSvrAddr
	}
}

func WithTracingSvrAddr(tracingSvrAddr string) ServerOption {
	return func(o *ServerOptions) {
		o.tracingSvrAddr = tracingSvrAddr
	}
}

func WithTracingSpanName(tracingSpanName string) ServerOption {
	return func(o *ServerOptions) {
		o.tracingSpanName = tracingSpanName
	}
}

func WithPluginNames(pluginNames []string) ServerOption {
	return func(o *ServerOptions) {
		o.pluginNames = pluginNames
	}
}

func WithInterceptors(interceptors []interceptor.ServerInterceptor) ServerOption {
	return func(o *ServerOptions) {
		o.interceptors = interceptors
	}
}

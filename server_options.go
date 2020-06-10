package zRPC

import (
	"time"

	"github.com/WeilunZ/zRPC/components/interceptor"
)

type ServerOptions struct {
	Address           string //e.g. 127.0.0.1:8080/www.baidu.com
	Network           string // e.g. tcp/udp
	Protocol          string
	Timeout           time.Duration
	SerializationType string   // serialization type, default: proto
	SelectorSvrAddr   string   // service discovery server Address, required when using the third-party service discovery plugin
	TracingSvrAddr    string   // tracing plugin server Address, required when using the third-party tracing plugin
	TracingSpanName   string   // tracing span name, required when using the third-party tracing plugin
	PluginNames       []string // plugin name
	Interceptors      []interceptor.ServerInterceptor
}

type ServerOption func(*ServerOptions)

func WithAddress(address string) ServerOption {
	return func(o *ServerOptions) {
		o.Address = address
	}
}

func WithNetwork(network string) ServerOption {
	return func(o *ServerOptions) {
		o.Network = network
	}
}

func WithProtocol(protocol string) ServerOption {
	return func(o *ServerOptions) {
		o.Protocol = protocol
	}
}

func WithTimeOut(timeout time.Duration) ServerOption {
	return func(o *ServerOptions) {
		o.Timeout = timeout
	}
}

func WithSerializationType(serializationType string) ServerOption {
	return func(o *ServerOptions) {
		o.SerializationType = serializationType
	}
}

func WithSelectorSvrAddr(selectorSvrAddr string) ServerOption {
	return func(o *ServerOptions) {
		o.SelectorSvrAddr = selectorSvrAddr
	}
}

func WithTracingSvrAddr(tracingSvrAddr string) ServerOption {
	return func(o *ServerOptions) {
		o.TracingSvrAddr = tracingSvrAddr
	}
}

func WithTracingSpanName(tracingSpanName string) ServerOption {
	return func(o *ServerOptions) {
		o.TracingSpanName = tracingSpanName
	}
}

func WithPluginNames(pluginNames []string) ServerOption {
	return func(o *ServerOptions) {
		o.PluginNames = pluginNames
	}
}

func WithInterceptors(interceptors []interceptor.ServerInterceptor) ServerOption {
	return func(o *ServerOptions) {
		o.Interceptors = interceptors
	}
}

package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/WeilunZ/myrpc/components/metrics"

	"github.com/WeilunZ/myrpc/components/connpool"
	"github.com/WeilunZ/myrpc/components/selector"
	"github.com/WeilunZ/myrpc/components/state"
	"github.com/WeilunZ/myrpc/transport"

	"github.com/WeilunZ/myrpc/components/protocol"

	"github.com/WeilunZ/myrpc/components/codec"
	"github.com/golang/protobuf/proto"
)

type Client interface {
	Invoke(ctx context.Context, req, resp interface{}, path string, opts ...Option) error
}

type defaultClient struct {
	opts *Options
}

//单例的全局唯一client
var DefaultClient = New()

var New = func() *defaultClient {
	return &defaultClient{
		opts: &Options{
			protocol: "proto",
		},
	}
}

var (
	invokeStatusCounter = metrics.NewCounterVec("client_invoke_error_count", "status")
)

func (c *defaultClient) Call(ctx context.Context, servicePath string,
	req interface{}, rsp interface{}, opts ...Option) error {
	// reflection calls need to be serialized using msgpack
	callOpts := make([]Option, 0, len(opts)+1)
	callOpts = append(callOpts, opts...)
	callOpts = append(callOpts, WithSerializationType(codec.MsgPack))

	//servicePath example: /helloworld.Greeter/SayHello
	err := c.Invoke(ctx, req, rsp, servicePath, callOpts...)
	if err != nil {
		return err
	}
	return nil
}

func (c *defaultClient) Invoke(ctx context.Context, req, resp interface{}, path string, opts ...Option) error {
	serialization := codec.GetSerialization(c.opts.serializationType)
	payload, err := serialization.Serialize(req)
	if err != nil {
		return errors.New("client request marshal failed")
	}
	clientCodec := codec.GetCodec(c.opts.protocol)

	// assemble header
	request := addReqHeader(ctx, c, payload)
	reqbuf, err := proto.Marshal(request)
	if err != nil {
		return err
	}

	reqbody, err := clientCodec.Encode(reqbuf)
	if err != nil {
		return err
	}

	clientTransport := c.NewClientTransport()
	clientTransportOpts := []transport.ClientTransportOption{
		transport.WithServiceName(c.opts.serviceName),
		transport.WithClientTarget(c.opts.target),
		transport.WithClientNetwork(c.opts.network),
		transport.WithClientPool(connpool.GetPool("default")),
		transport.WithSelector(selector.GetSelector(c.opts.selectorName)),
		transport.WithTimeout(c.opts.timeout),
	}
	frame, err := clientTransport.Send(ctx, reqbody, clientTransportOpts...)
	if err != nil {
		invokeStatusCounter.WithLabelValues("fail").Inc()
		return err
	}

	rspbuf, err := clientCodec.Decode(frame)
	if err != nil {
		invokeStatusCounter.WithLabelValues("fail").Inc()
		return err
	}

	// parse protocol header
	response := &protocol.Response{}
	if err = proto.Unmarshal(rspbuf, response); err != nil {
		invokeStatusCounter.WithLabelValues("fail").Inc()
		return err
	}

	if response.RetCode != 0 {
		invokeStatusCounter.WithLabelValues("fail").Inc()
		return state.New(response.RetCode, response.RetMsg)
	}

	invokeStatusCounter.WithLabelValues("success").Inc()
	return serialization.Deserialize(response.Payload, resp)
}

func (c *defaultClient) NewClientTransport() transport.ClientTransport {
	return transport.GetClientTransport(c.opts.protocol)
}

func addReqHeader(ctx context.Context, client *defaultClient, payload []byte) *protocol.Request {
	servicePath := fmt.Sprintf("/%s/%s", client.opts.serviceName, client.opts.method)

	// TODO pass metadata
	// TODO add authentication info

	request := &protocol.Request{
		ServicePath: servicePath,
		Payload:     payload,
		Metadata:    nil,
	}

	return request
}

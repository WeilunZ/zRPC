package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/WeilunZ/zRPC/components/interceptor"
	"github.com/WeilunZ/zRPC/components/utils"

	"github.com/WeilunZ/zRPC/components/metrics"

	"github.com/WeilunZ/zRPC/components/connpool"
	"github.com/WeilunZ/zRPC/components/selector"
	"github.com/WeilunZ/zRPC/components/state"
	"github.com/WeilunZ/zRPC/transport"

	"github.com/WeilunZ/zRPC/components/protocol"

	"github.com/WeilunZ/zRPC/components/codec"
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
	for _, o := range opts {
		o(c.opts)
	}

	if c.opts.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.opts.timeout)
		defer cancel()
	}

	serviceName, method, err := utils.ParseServicePath(path)
	if err != nil {
		return err
	}
	c.opts.serviceName = serviceName
	c.opts.method = method

	return interceptor.ClientIntercept(ctx, req, resp, c.opts.interceptors, c.doInvoke)
}

func (c *defaultClient) doInvoke(ctx context.Context, req, rsp interface{}) error {
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
	return serialization.Deserialize(response.Payload, rsp)
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
		Metadata: map[string][]byte{
			"hashKey": []byte(ctx.Value("hashKey").(string)),
		},
	}

	return request
}

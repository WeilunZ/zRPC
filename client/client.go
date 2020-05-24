package client

import (
	"context"

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
	payload, err := serialization.Marshal(req)
	if err != nil {
		return codes.NewFrameworkError(codes.ClientMsgErrorCode, "request marshal failed ...")
	}
	// assemble header
	request := addReqHeader(ctx, payload)
	reqbuf, err := proto.Marshal(request)
	if err != nil {
		return err
	}

}

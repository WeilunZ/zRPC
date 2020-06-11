package main

import (
	"context"
	"fmt"
	"github.com/WeilunZ/zRPC/client"
	"github.com/WeilunZ/zRPC/components/codec"
	"github.com/WeilunZ/zRPC/example/helloworld"
	"time"
)

func main() {
	opts := []client.Option{
		client.WithTarget("127.0.0.1:8000"),
		client.WithNetwork("tcp"),
		client.WithTimeout(time.Millisecond * 2000),
		client.WithSerializationType(codec.MsgPack),
	}
	c := client.DefaultClient
	req := &helloworld.HelloRequest{
		Msg: "hello",
	}
	resp := &helloworld.HelloResponse{
		Msg: "world",
	}
	err := c.Call(context.Background(), "/helloworld.Greeter/SayHello", req, resp, opts...)
	fmt.Println(resp.Msg, err)
}

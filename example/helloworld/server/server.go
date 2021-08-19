package main

import (
	"github.com/WeilunZ/zRPC"
	"github.com/WeilunZ/zRPC/components/codec"
	"github.com/WeilunZ/zRPC/example/helloworld"
	"time"
)

func main() {
	opts := []zRPC.ServerOption{
		zRPC.WithNetwork("tcp"),
		zRPC.WithSerializationType(codec.MsgPack),
		zRPC.WithAddress("127.0.0.1:8000"),
		zRPC.WithPluginNames([]string{"consul"}),
		zRPC.WithSelectorSvrAddr("127.0.0.1:8500"),
		zRPC.WithTimeOut(time.Millisecond * 2000),
	}
	s := zRPC.NewServer(opts...)
	if err := s.RegisterService("helloworld.Greeter", new(helloworld.Service)); err != nil {
		panic(err)
	}
	s.Serve()
}

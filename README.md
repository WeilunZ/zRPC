# zRPC

## Usage
1.定义一个服务
```go
type Service struct {

}

type HelloRequest struct {
	Msg string
}

type HelloResponse struct {
	Msg string
}

func (s *Service) SayHello(ctx context.Context, req *HelloRequest) (*HelloResponse, error){
	rsp := &HelloResponse{Msg:"world"}
	return rsp, nil
}
```
2.启动并发布服务
```go
func main(){
	opts := []zRPC.ServerOption{
		zRPC.WithNetwork("tcp"),
		zRPC.WithSerializationType(codec.MsgPack),
		zRPC.WithAddress("127.0.0.1:8000"),
		zRPC.WithTimeOut(time.Millisecond * 2000),
	}
	s := zRPC.NewServer(opts...)
	if err := s.RegisterService("/helloworld.Greeter", new(helloworld.Service)); err != nil{
		panic(err)
	}
	s.Serve()
}
```
3.调用服务
```go
func main(){
	opts := []client.Option{
		client.WithTarget("127.0.0.1:8000"),
		client.WithNetwork("tcp"),
		client.WithTimeout(time.Millisecond * 2000),
		client.WithSerializationType(codec.MsgPack),
	}
	c := client.DefaultClient
	req := &helloworld.HelloRequest{
		Msg:"hello",
	}
	resp := &helloworld.HelloResponse{
		Msg: "world",
	}
	err := c.Call(context.Background(), "/helloworld.Greeter/SayHello", req, resp, opts...)
	fmt.Println(resp.Msg, err)
}
```
4.运行并打印
```go
➜ cd example/helloworld
➜ go run server/server.go
➜ go run client/client.go

world <nil>
```

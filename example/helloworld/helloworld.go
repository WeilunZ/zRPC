package helloworld

import "context"

type Service struct {
}

type HelloRequest struct {
	Msg string
}

type HelloResponse struct {
	Msg string
}

func (s *Service) SayHello(ctx context.Context, req *HelloRequest) (*HelloResponse, error) {
	rsp := &HelloResponse{Msg: "world"}
	return rsp, nil
}

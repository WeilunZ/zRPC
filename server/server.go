package server

import (
	"os"
	"os/signal"
	"reflect"
	"syscall"

	tinyRPC "github.com/WeilunZ/myrpc"
)

type Server struct {
	opts     *ServerOptions
	services map[string]tinyRPC.Service
}

func NewServer(opt ...ServerOption) {
	s := &Server{
		opts:     &ServerOptions{},
		services: make(map[string]tinyRPC.Service),
	}
	for _, o := range opt {
		o(s.opts)
	}
}

func (*Server) RegisterService(serviceName string, svr interface{}) {
	serviceType := reflect.TypeOf(svr)
	serviceValue := reflect.ValueOf(svr)

}

func (s *Server) Serve() {
	for _, service := range s.services {
		go service.Serve(s.opts)
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGSEGV)
	<-ch
	s.Close()
}

func (s *Server) Close() {
}

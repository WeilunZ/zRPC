package zRPC

import (
	"context"
	"fmt"
	"github.com/WeilunZ/zRPC/plugin/consul"
	"github.com/WeilunZ/zRPC/plugin/tracing"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/WeilunZ/zRPC/components/log"

	"github.com/WeilunZ/zRPC/components/interceptor"

	"github.com/WeilunZ/zRPC/plugin"
)

type Server struct {
	opts     *ServerOptions
	services map[string]Service
	plugins  []plugin.Plugin
	closing  bool
}

func initSupportedPlugins(){
	consul.Init()
	tracing.Init()
}

func NewServer(opt ...ServerOption) *Server {

	initSupportedPlugins()

	s := &Server{
		opts:     &ServerOptions{},
		services: make(map[string]Service),
	}
	for _, o := range opt {
		o(s.opts)
	}
	for name, plugin := range plugin.PluginMap {
		if !s.containPlugin(name) {
			continue
		}
		s.plugins = append(s.plugins, plugin)
	}
	return s
}

func (s *Server) containPlugin(name string) bool {
	for _, p := range s.opts.PluginNames {
		if p == name {
			return true
		}
	}
	return false
}

func (s *Server) RegisterService(serviceName string, svr interface{}) error {
	// 基于反射
	serviceType := reflect.TypeOf(svr)
	serviceValue := reflect.ValueOf(svr)
	sd := &ServiceDesc{
		ServiceName: serviceName,
		HandlerType: (*interface{})(nil),
		Svr:         svr,
	}
	methods, err := getServiceMethods(serviceType, serviceValue)
	if err != nil {
		return err
	}
	sd.Methods = methods
	s.Register(sd, svr)
	return nil
}

func getServiceMethods(serviceType reflect.Type, serviceValue reflect.Value) ([]*Method, error) {
	methods := make([]*Method, 0)
	for i := 0; i < serviceType.NumMethod(); i++ {
		m := serviceType.Method(i)
		if err := validateMethod(m.Type); err != nil {
			return nil, err
		}
		method := &Method{
			MethodName: m.Name,
			Handler: func(service interface{}, ctx context.Context, deserialize func(interface{}) error, interceptors []interceptor.ServerInterceptor) (interface{}, error) {
				reqType := m.Type.In(2)
				req := reflect.New(reqType.Elem()).Interface()
				if err := deserialize(req); err != nil {
					return nil, err
				}
				if len(interceptors) == 0 {
					values := m.Func.Call([]reflect.Value{serviceValue, reflect.ValueOf(ctx), reflect.ValueOf(req)})
					return values[0].Interface(), nil
				}
				handler := func(ctx context.Context, reqbody interface{}) (interface{}, error) {
					values := m.Func.Call([]reflect.Value{serviceValue, reflect.ValueOf(ctx), reflect.ValueOf(req)})
					return values[0].Interface(), nil
				}
				return interceptor.ServerIntercept(ctx, req, interceptors, handler)
			},
		}
		methods = append(methods, method)
	}
	return methods, nil
}

func validateMethod(m reflect.Type) error {
	if m.NumIn() < 3 {
		return fmt.Errorf("method %s invalid, must have at least 2 params", m.Name())
	}
	if m.NumOut() != 2 {
		return fmt.Errorf("method %s invalid, must return 2 values", m.Name())
	}
	//parameter1: context
	if !m.In(1).Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
		return fmt.Errorf("method %s invalid, first param must be context type", m.Name())
	}
	//parameter2: pointer
	if m.In(2).Kind() != reflect.Ptr {
		return fmt.Errorf("method %s invalid, second param must be pointer", m.Name())
	}
	// return type 1 : pointer
	// return type 2 : error
	replyType := m.Out(0)
	if replyType.Kind() != reflect.Ptr {
		return fmt.Errorf("method %s invalid, reply type must be pointer", m.Name())
	}

	errType := m.Out(1)
	if !errType.Implements(reflect.TypeOf((*error)(nil)).Elem()) {
		return fmt.Errorf("method %s invalid, second param must implements error", m.Name())
	}
	return nil

}

func (s *Server) Register(sd *ServiceDesc, svr interface{}) {
	if sd == nil || svr == nil {
		return
	}
	ht := reflect.TypeOf(sd.HandlerType).Elem()
	st := reflect.TypeOf(svr)
	if !st.Implements(ht) {
		log.Fatalf("handlerType %v not match service : %v ", ht, st)
	}

	ser := &service{
		svr:         svr,
		serviceName: sd.ServiceName,
		handlers:    make(map[string]Handler),
	}

	for _, method := range sd.Methods {
		ser.handlers[method.MethodName] = method.Handler
	}

	s.services[sd.ServiceName] = ser
}

func (s *Server) Serve() {
	err := s.InitPlugins()
	if err != nil {
		panic(err)
	}


	for _, service := range s.services {
		go service.Serve(s.opts)
	}

	// gracefully shutdown
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGSEGV)
	<-ch
	s.Close()
}

func (s *Server) Close() {
	s.closing = false

	for _, service := range s.services {
		service.Close()
	}
}

func (s *Server) InitPlugins() error {
	// init plugins
	for _, p := range s.plugins {

		switch val := p.(type) {

		case plugin.ResolverPlugin:
			var services []string
			for _, ss := range s.services{
				services = append(services, ss.Name())
			}
			pluginOpts := []plugin.Option{
				plugin.WithSelectorSvrAddr(s.opts.SelectorSvrAddr),
				plugin.WithSvrAddr(s.opts.Address),
				plugin.WithServices(services),
			}
			if err := val.Init(pluginOpts...); err != nil {
				log.Errorf("resolver init error, %v", err)
				return err
			}

		default:

		}

	}

	return nil
}

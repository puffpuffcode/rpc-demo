package main

import (
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Args struct {
	X,Y int
}

type ServiceA struct {}

func (s *ServiceA) Serve(args *Args, reply *int) error {
	*reply = args.X + args.Y
	return nil
}

func main(){
	service := new(ServiceA)
	// register a service
	rpc.Register(service)
	// based on tcp protocol
	l, _ := net.Listen("tcp", ":19090")
	// serve
	for {
		c, _ := l.Accept()
		rpc.ServeCodec(jsonrpc.NewServerCodec(c))
	}
}
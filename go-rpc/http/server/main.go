package main

import (
	"net"
	"net/http"
	"net/rpc"
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
	// based on http protocol
	rpc.HandleHTTP()
	// listen
	l, _ := net.Listen("tcp", ":19090")
	// serve
	http.Serve(l, nil)
}
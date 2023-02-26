package main

import (
	"add_server/proto"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
)

type addServer struct {
	proto.UnimplementedAdderServer
}

func (s *addServer) Add(ctx context.Context, req *proto.AddRequest) (*proto.AddResponse, error) {
	res := &proto.AddResponse{
		Res: req.GetA() + req.GetB(),
	}
	return res, nil
}

func main() {
	l, err := net.Listen("tcp", ":19090")
	if err != nil {
		log.Fatalln("net.Listen failed, err:", err.Error())
	}
	s := grpc.NewServer()
	proto.RegisterAdderServer(s, &addServer{})
	err = s.Serve(l)
	if err != nil {
		log.Fatalln("s.Serve(l) failed, err:",err.Error())
	}
}

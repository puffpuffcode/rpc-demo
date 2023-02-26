package main

import (
	"context"
	"hello_server/pb"
	"log"
	"net"

	"google.golang.org/grpc"
)

// grpc server
type server struct {
	pb.UnimplementedGreeterServer
}

// to be complement
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	reply := "Hello " + in.GetName()
	return &pb.HelloResponse{Reply: reply}, nil
}

func main() {
	l, err := net.Listen("tcp", ":19090")
	if err != nil {
		log.Fatalln(err)
	}
	// create server
	s := grpc.NewServer()
	// register
	pb.RegisterGreeterServer(s, &server{})
	// run 
	err = s.Serve(l)
	if err != nil {
		log.Fatalln(err)
	}
}

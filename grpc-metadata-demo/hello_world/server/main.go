package main

import (
	"context"
	"fmt"
	"hello_server/pb"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// grpc server
type server struct {
	pb.UnimplementedGreeterServer
}

// to be complement
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	reply := "Hello " + in.GetName()
	// get metadata
	m, b := metadata.FromIncomingContext(ctx)
	if !b {
		fmt.Println("not carries data.") // deny
		return nil, status.Error(codes.DataLoss, "invalid request")
	} 
	vl := m.Get("token")
	if len(vl) < 1 || vl[0] != `app-test-makito` {
		return nil, status.Error(codes.Unauthenticated, "invalid request")
	}
	fmt.Println("got token:",vl[0])
	// resp with metadata
	grpc.SendHeader(ctx,metadata.MD{
		"resp":[]string{"i got it !"},
	})
	defer func ()  {
		t := metadata.Pairs("k1", "v1")
		grpc.SetTrailer(ctx,t)
	}()
	
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

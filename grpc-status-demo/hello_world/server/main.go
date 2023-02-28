package main

import (
	"context"
	"fmt"
	"hello_server/pb"
	"log"
	"net"
	"sync"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func main() {
	l, err := net.Listen("tcp", ":19090")
	if err != nil {
		log.Fatalln(err)
	}
	// create server
	s := grpc.NewServer()
	// register
	pb.RegisterGreeterServer(s, &server{nameCall: make(map[string]struct{}), mapMu: sync.Mutex{}})
	// run
	err = s.Serve(l)
	if err != nil {
		log.Fatalln(err)
	}
}

// grpc server
type server struct {
	pb.UnimplementedGreeterServer
	nameCall map[string]struct{}
	mapMu    sync.Mutex
}

// to be complement
// one name only can be called once or return err

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	// check whether name was duplicated
	s.mapMu.Lock()
	defer s.mapMu.Unlock()
	if _, ok := s.nameCall[in.GetName()]; ok {
		s := status.New(codes.AlreadyExists, "调用次数超过一次了！")
		ss, err := s.WithDetails(&errdetails.QuotaFailure{
			Violations: []*errdetails.QuotaFailure_Violation{
				{
					Subject:     fmt.Sprintf("name:%s", in.GetName()),
					Description: "每个 name 只能调用一次！",
				},
			},
		})
		if err != nil {
			return nil, s.Err()
		}
		return nil, ss.Err() // err with details
	}
	s.nameCall[in.GetName()] = struct{}{}

	return &pb.HelloResponse{Reply: "Hello " + in.GetName()}, nil
}

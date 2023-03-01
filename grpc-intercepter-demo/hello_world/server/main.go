package main

import (
	"context"
	"fmt"
	"hello_server/pb"
	"log"
	"net"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
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
	return &pb.HelloResponse{Reply: reply}, nil
}

// intercepter
func unaryServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	// auth
	m, b := metadata.FromIncomingContext(ctx)
	if !b {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata.")
	}
	if !valid(m["authorization"]) {
		return nil, status.Errorf(codes.Unauthenticated, "invaild token.")
	}
	i, err := handler(ctx, req)
	if err != nil {
		fmt.Printf("RPC failed with error %v\n", err)
	}
	return i, err
}

func valid(arr []string) bool {
	if len(arr) < 1 {
		return false
	}
	token := strings.Trim(arr[0], "Bearer ")
	return token == "some-secret-token"
}

// wrappedServerStream
type wrappedServerStream struct {
	grpc.ServerStream
}

func NewWarppedServerStream(ss grpc.ServerStream) *wrappedServerStream {
	wss := &wrappedServerStream{ss}
	return wss
}

func (s *wrappedServerStream) SendMsg(m interface{}) error {
	log.Printf("Send a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	return s.ServerStream.SendMsg(m)
}

func (s *wrappedServerStream) RecvMsg(m interface{}) error {
	log.Printf("Recv a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	return s.ServerStream.RecvMsg(m)
}

// stream intercepter
func streamIntercepter(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	md, b := metadata.FromIncomingContext(ss.Context())
	if !b {
		return status.Errorf(codes.InvalidArgument, "missing metadata.")
	}
	if !valid(md["authorization"]) {
		return status.Errorf(codes.Unauthenticated, "invaild token.")
	}
	err := handler(srv, NewWarppedServerStream(ss))
	return err
}

func main() {
	l, err := net.Listen("tcp", ":19090")
	if err != nil {
		log.Fatalln(err)
	}
	tc, err := credentials.NewServerTLSFromFile("certs/server.crt", "certs/server.key")
	if err != nil {
		log.Fatalln("credentials.NewServerTLSFromFile err:", err)
	}
	// create server
	s := grpc.NewServer(
		grpc.Creds(tc),
		grpc.UnaryInterceptor(unaryServerInterceptor),
		grpc.StreamInterceptor(streamIntercepter),
	)
	// register
	pb.RegisterGreeterServer(s, &server{})
	// run
	err = s.Serve(l)
	if err != nil {
		log.Fatalln(err)
	}
}

package main

import (
	"context"
	"errors"
	"fmt"
	"hello_server/pb"
	"io"
	"log"
	"net"
	"strings"

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

// stream rep
func (s *server) MutiReplies(req *pb.HelloRequest, repSever pb.Greeter_MutiRepliesServer) error {
	words := []string{
		"你好",
		"Hello",
		"こんにちは",
	}

	for _, hello := range words {
		rep := &pb.HelloResponse{
			Reply: strings.Join([]string{hello, req.GetName() + "!"}, " "),
		}
		if err := repSever.Send(rep); err != nil {
			fmt.Println("repSever.Send(rep) err:", err.Error())
			return err
		}
	}
	return nil
}

// recv stream
func (s *server) MutiRequests(stream pb.Greeter_MutiRequestsServer) error {
	reply := "Hello"
	var names []string
	for {
		req, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Fatalln("stream.Recv() err:", err)
		}
		names = append(names, req.GetName())
	}
	// 发送结果到 cli
	if err := stream.SendAndClose(&pb.HelloResponse{
		Reply: strings.Join(append([]string{reply}, names...), ", "),
	}); err != nil {
		log.Fatalln("stream.SendAndClose err:", err)
	}
	return nil
}

// 双向
func (s *server) BidiHello(stream pb.Greeter_BidiHelloServer) error {
	for {
		hr, err := stream.Recv()
		if errors.Is(err,io.EOF) {
			break
		}
		if err != nil {
			return err
		}
		msg := hr.GetName()
		// do something for msg
		msg = magic(msg)
		// send back to cli
		if err = stream.Send(&pb.HelloResponse{
			Reply: msg,
		}); err != nil {
			return err
		}
	}
	return nil
}

// magic 一段价值连城的“人工智能”代码
func magic(s string) string {
	s = strings.ReplaceAll(s, "吗", "")
	s = strings.ReplaceAll(s, "吧", "")
	s = strings.ReplaceAll(s, "你", "我")
	s = strings.ReplaceAll(s, "？", "!")
	s = strings.ReplaceAll(s, "?", "!")
	return s
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

package main

import (
	"context"
	"gateway_demo/proto"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"


	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type server struct {
	proto.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {
	name := req.GetName()
	if name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "无效的名称")
	}
	return &proto.HelloResponse{
		Message: "Hello " + name + "!",
	}, nil
}

func main() {
	go StartRPCServer()
	go StartHttpServer()
	RPCRequest()
	HTTPRequest()
}

func StartRPCServer() {
	l, err := net.Listen("tcp", ":19090")
	if err != nil {
		log.Fatalf("net.Listen err: %v\n", err)
	}
	s := grpc.NewServer()
	proto.RegisterGreeterServer(s, &server{})
	log.Fatalln(s.Serve(l))
}

// 创建一个连接到以上 RPC 服务
func StartHttpServer() {
	conn, err := grpc.DialContext(
		context.Background(),
		":19090",
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("grpc.DialContext err: %v\n", err)
	}
	defer conn.Close()
	
	sm := runtime.NewServeMux()
	// 绑定服务
	err = proto.RegisterGreeterHandler(context.Background(), sm, conn)
	if err != nil {
		log.Fatalf("proto.RegisterGreeterHandler err: %v\n", err)
	}

	gwServer := &http.Server{
		Addr:    ":19191",
		Handler: sm,
	}

	log.Fatalln(gwServer.ListenAndServe())
}

func RPCRequest() {
	cc, err := grpc.Dial(":19090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("grpc.Dial err: %v\n", err)
	}
	defer cc.Close()
	greaterCli := proto.NewGreeterClient(cc)

	ctx, cf := context.WithCancel(context.Background())
	defer cf()

	resp, err := greaterCli.SayHello(ctx, &proto.HelloRequest{
		Name: "Makito",
	})
	if err != nil {
		log.Fatalf("greaterCli.SayHello err: %v\n", err)
	}

	log.Printf("resp: %v\n", resp.Message)
}

func HTTPRequest() {
	resp, err := http.Post("http://127.0.0.1:19191/v1/hello", "application/json", strings.NewReader(
		`{
			"name": "Makitooooo"
		}`,
	))
	if err != nil {
		log.Fatalf("http.Post err: %v\n", err)
	}
	defer resp.Body.Close()
	bytes, _ := ioutil.ReadAll(resp.Body)
	log.Printf("get http post resp: %s\n", bytes)
}

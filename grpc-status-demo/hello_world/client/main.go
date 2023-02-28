package main

import (
	"context"
	"log"
	"time"

	"code.xxx.com/client/hello/proto"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

// grpc client
func main() {
	// conn to server
	cc, err := grpc.Dial(":19090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("grpc.Dial failed, err:", err.Error())
	}
	defer cc.Close()

	// create client
	c := proto.NewGreeterClient(cc)
	ctx, cf := context.WithTimeout(context.Background(), time.Second)
	defer cf()

	// req
	hr, err := c.SayHello(
		ctx,
		&proto.HelloRequest{Name: "makito"},
	)
	if err != nil {
		// 解析带 details 的 err 信息
		s := status.Convert(err)
		log.Printf("code:%s | msg:%s\n", s.Code(), s.Message())
		for _, v := range s.Details() {
			switch info := v.(type) {
			case *errdetails.QuotaFailure:
				log.Printf("QuotaFailure: %v\n", info.GetViolations())
			default:
				log.Printf("Unexpected: %v\n", info)
			}
		}
		return
	}

	//resp
	log.Println("get resp:", hr.GetReply())
}

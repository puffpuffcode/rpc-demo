package main

import (
	"context"
	"log"
	"time"

	"code.xxx.com/client/hello/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	// resp
	hr, err := c.SayHello(ctx, &proto.HelloRequest{Name: "makito"})
	if err != nil {
		log.Fatalln("c.SayHello failed, err:", err.Error())
	}
	log.Println("get resp:", hr.GetReply())
}

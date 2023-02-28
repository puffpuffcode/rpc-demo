package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"code.xxx.com/client/hello/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
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
	// req with metadata
	md := metadata.Pairs(
		"token", "app-test-makito",
	)
	// add metadata
	// if expected to get header or trailer, should previously assign
	
	ctx = metadata.NewOutgoingContext(ctx, md)
	var header, trailer metadata.MD
	hr, err := c.SayHello(
		ctx,
		&proto.HelloRequest{Name: "makito"},
		grpc.Header(&header),
		grpc.Trailer(&trailer),
	)
	if err != nil {
		log.Fatalln("c.SayHello failed, err:", err.Error())
	}

	// before resp
	fmt.Printf("%#v\n",header)

	//resp
	log.Println("get resp:", hr.GetReply())
	
	// after resp
	fmt.Printf("%#v\n",trailer)
}

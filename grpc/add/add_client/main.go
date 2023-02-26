package main

import (
	"add_client/proto"
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	cc, err := grpc.Dial(":19090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("grpc.Dial failed, err:", err.Error())
	}
	defer cc.Close()

	c := proto.NewAdderClient(cc)
	ctx, cf := context.WithTimeout(context.Background(), time.Second)
	defer cf()
	resp, err := c.Add(ctx, &proto.AddRequest{A: 100000, B: 14514})
	if err != nil {
		log.Fatalln("c.Add failed, err:", err.Error())
	}
	log.Println("resp:", resp.GetRes())
}

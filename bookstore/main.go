package main

import (
	"bookstore/pb"
	"bookstore/server"
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	go SomeClient()
	server.StartServer()
	// server.StartHttpProxy()
}

func SomeClient() {
	cc, _ := grpc.Dial(":19090", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	defer cc.Close()
	cli := pb.NewBookStoreClient(cc)
	resp, err := cli.ListBooks(
		context.Background(),
		&pb.ListBooksRequest{
			Shelf:     3,
			PageToken: "",
		},
	)
	if err != nil {
		log.Fatalf("%#v\n",err)
	}
	log.Printf("resp: %v\n",resp)
}

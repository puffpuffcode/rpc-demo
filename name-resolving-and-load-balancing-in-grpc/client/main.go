package main

import (
	"client/pb"
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	name = flag.String("name", "makito", "input your name.")
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(
		"makito:///resolver.makito.com",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithResolvers(&myResolverBuilder{}),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`), // 设置负载均衡策略
	)
	if err != nil {
		log.Fatalln(err)
	}

	c := pb.NewGreeterClient(conn)
	// 调用RPC方法
	for i := 0; i < 10; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		resp, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
		if err != nil {
			fmt.Printf("c.SayHello failed, err:%v\n", err)
			return
		}
		// 拿到了RPC响应
		fmt.Printf("resp:%v\n", resp.GetMessage())
	}
}

package main

import (
	"context"
	"fmt"
	"log"

	"code.xxx.com/client/hello/proto"
	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	_ "github.com/mbobakov/grpc-consul-resolver"
)

// grpc client
func main() {
	// SimpleWayToInvoke()
	ElegantWayToInvoce()
}

func SimpleWayToInvoke() {
	// connect to consul
	c, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatalln(err)
	}
	// find service which are avilible
	m, err := c.Agent().ServicesWithFilter("Service == hello_server")
	if err != nil {
		log.Fatalln(err)
	}
	var addr string
	for _, v := range m {
		addr = fmt.Sprintf("%s:%d", v.Address, v.Port)
		break
	}

	// conn to service
	cc, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln(err)
	}
	gc := proto.NewGreeterClient(cc)
	ctx, cf := context.WithCancel(context.Background())
	defer cf()
	resp, err := gc.SayHello(ctx, &proto.HelloRequest{
		Name: "makito",
	})
	if err != nil {
		log.Fatalln("sayHello err:", err.Error())
	}
	log.Printf("resp from %s: %s", addr, resp.GetReply())
}

func ElegantWayToInvoce() {
	c, err := grpc.Dial(
		`consul://localhost:8500/hello_server?healthy=true`,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
	)
	if err != nil {
		log.Fatalln("grpc dial failed:", err.Error())
	}
	defer c.Close()
	gc := proto.NewGreeterClient(c)

	for i := 0; i < 10; i++ {
		ctx, cf := context.WithCancel(context.Background())
		defer cf()
		resp, err := gc.SayHello(ctx, &proto.HelloRequest{
			Name: "makito",
		})
		if err != nil {
			log.Fatalln("sayHello err:", err.Error())
		}
		log.Printf("resp: %s", resp.GetReply())
	}

}

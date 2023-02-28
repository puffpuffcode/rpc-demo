package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
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

	replyDemo(c)
	fmt.Println("----------------------")
	streamReplyDemo(c)
	fmt.Println("----------------------")
	mutiGreetingDemo(c)
	fmt.Println("----------------------")
	bidiHelloDemo(c)
}

// simple cli rpc
func replyDemo(c proto.GreeterClient) {
	ctx, cf := context.WithTimeout(context.Background(), time.Second)
	defer cf()
	// resp
	hr, err := c.SayHello(ctx, &proto.HelloRequest{Name: "makito"})
	if err != nil {
		log.Fatalln("c.SayHello failed, err:", err.Error())
	}
	log.Println("get resp:", hr.GetReply())
}

// recv server stream
func streamReplyDemo(c proto.GreeterClient) {
	ctx1, cf := context.WithTimeout(context.Background(), time.Second)
	defer cf()
	req := proto.HelloRequest{
		Name: "makito",
	}
	repCli, err := c.MutiReplies(ctx1, &req)
	if err != nil {
		log.Fatalln("c.MutiReplies(ctx1, &req) err:", err.Error())
	}

	for {
		resp, err := repCli.Recv()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Println("err:", err.Error())
			return
		}
		log.Println("recv data:", resp.GetReply())
	}
}

// send to server stream
func mutiGreetingDemo(c proto.GreeterClient) {
	ctx, cf := context.WithTimeout(context.Background(), time.Second)
	defer cf()
	cli, err := c.MutiRequests(ctx)
	if err != nil {
		log.Fatalln("c.MutiRequests(ctx) err:", err)
	}
	// req arrays
	a := []string{
		"cb",
		"makito",
		"ceri",
	}
	// send to server
	for _, name := range a {
		cli.Send(&proto.HelloRequest{
			Name: name,
		})
	}

	// attention!!!
	hr, err := cli.CloseAndRecv()
	if err != nil {
		log.Fatalln("err:", err)
	}
	fmt.Println("get resp:", hr.GetReply())
}

// run bidi hello
func bidiHelloDemo(c proto.GreeterClient) {
	ctx, cf := context.WithCancel(context.Background())
	defer cf()
	stream, err := c.BidiHello(ctx)
	if err != nil {
		log.Fatalln("c.BidiHello(ctx) err:", err)
	}
	// run a goroutine to listen req
	waitC := make(chan struct{})
	go func() {
		for {
			hr, err := stream.Recv()
			if errors.Is(err, io.EOF) {
				waitC <- struct{}{}
				return
			}
			if err != nil {
				log.Fatalln("err:", err)
			}
			fmt.Printf("%s\n", hr.GetReply())
		}
	}()
	// listen std input
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			s, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("read from stdin err:", err)
			}
			s = strings.TrimSpace(s)
			if len(s) == 0 {
				continue
			}
			if s == "QUIT" {
				waitC <- struct{}{}
				break
			}
			// send to server
			if err = stream.Send(&proto.HelloRequest{
				Name: s,
			}); err != nil {
				fmt.Println("stream.Send err:", err)
			}
		}
	}()

	<-waitC

	err = stream.CloseSend()
	if err != nil {
		log.Fatalln(err)
	}
}

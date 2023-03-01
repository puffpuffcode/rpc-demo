package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"code.xxx.com/client/hello/proto"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

func unaryClientIntercepter(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	var credsConfigured bool
	for _, o := range opts {
		_, ok := o.(grpc.PerRPCCredsCallOption)
		if ok {
			credsConfigured = true
			break
		}
	}
	if !credsConfigured {
		opts = append(opts, grpc.PerRPCCredentials(oauth.NewOauthAccess(&oauth2.Token{
			AccessToken: "some-secret-token",
		})))
	}
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	end := time.Now()
	fmt.Printf("RPC: %s, start time: %s, end time: %s, err: %v\n", method, start.Format("Basic"), end.Format(time.RFC3339), err)
	return err
}

type wrappedStream struct {
	grpc.ClientStream
}

func (w *wrappedStream) RecvMsg(m interface{}) error {
	log.Printf("Receive a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	return w.ClientStream.RecvMsg(m)
}

func (w *wrappedStream) SendMsg(m interface{}) error {
	log.Printf("Send a message (Type: %T) at %v", m, time.Now().Format(time.RFC3339))
	return w.ClientStream.SendMsg(m)
}

func newWrappedStream(s grpc.ClientStream) grpc.ClientStream {
	return &wrappedStream{s}
}

// streamInterceptor 流式拦截器
func streamInterceptor(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	var credsConfigured bool
	for _, o := range opts {
		_, ok := o.(*grpc.PerRPCCredsCallOption)
		if ok {
			credsConfigured = true
			break
		}
	}
	if !credsConfigured {
		opts = append(opts, grpc.PerRPCCredentials(oauth.NewOauthAccess(&oauth2.Token{
			AccessToken: "some-secret-token",
		})))
	}
	s, err := streamer(ctx, desc, cc, method, opts...)
	if err != nil {
		return nil, err
	}
	return newWrappedStream(s), nil
}

// grpc client
func main() {
	// conn to server
	tc, err := credentials.NewClientTLSFromFile("certs/client.crt", "atelier.icu")
	if err != nil {
		log.Fatalf("credentials.NewClientTLSFromFile err: %v\n", err)
	}
	cc, err := grpc.Dial(
		":19090",
		grpc.WithTransportCredentials(tc),
		grpc.WithUnaryInterceptor(unaryClientIntercepter),
		grpc.WithStreamInterceptor(streamInterceptor),
	)
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

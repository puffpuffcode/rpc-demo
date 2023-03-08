package main

import (
	"context"
	"flag"
	"fmt"
	"hello_server/pb"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc"

	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	port = flag.Int64("port", 10000, "to export")
)

func init() {
	flag.Parse()
}

// grpc server
type server struct {
	pb.UnimplementedGreeterServer
	addr string
}

// to be complement
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	reply := "Hello " + in.GetName() + " from " + s.addr
	return &pb.HelloResponse{Reply: reply}, nil
}

// consul
type consul struct {
	client *api.Client
}

// new consul
func NewConsul(addr string) (*consul, error) {
	ip, _ := GetOutboundIP()
	config := api.Config{
		Address: fmt.Sprintf("%s:%d", ip, 8500),
	}
	cli, err := api.NewClient(&config)
	return &consul{cli}, err
}

// get export ip
func GetOutboundIP() (net.IP, error) {
	c, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer c.Close()
	udpAddr := c.LocalAddr().(*net.UDPAddr)
	return udpAddr.IP, nil
}

func main() {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalln(err)
	}
	// create server
	s := grpc.NewServer()
	// register hello server
	ip, err := GetOutboundIP()
	pb.RegisterGreeterServer(s, &server{addr: fmt.Sprintf("%s:%d",ip.String(),*port)})
	// register health check server
	healthpb.RegisterHealthServer(s, health.NewServer())

	// register service to consul
	
	if err != nil {
		log.Fatalln(err)
	}
	consul, err := NewConsul(ip.String())
	if err != nil {
		log.Fatalln(err)
	}
	srv1 := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s-%d", "hello_server", ip.String(), *port),
		Name:    "hello_server",
		Tags:    []string{"hello", "makito"},
		Address: ip.String(),
		Port:    int(*port),
		// add health check
		Check: &api.AgentServiceCheck{
			Name:                           "hello_server-check",
			Interval:                       "5s",
			Timeout:                        "1m",
			GRPC:                           fmt.Sprintf("%s:%d", ip.String(), *port),
			DeregisterCriticalServiceAfter: "10s",
		},
	}
	if err = consul.client.Agent().ServiceRegister(srv1); err != nil {
		log.Fatalln("ServiceRegister err:", err.Error())
	}
	// run
	log.Printf("starting server on %s:%d", ip.String(), *port)
	go func() {
		err = s.Serve(l)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	log.Println("waiting quit signal ...")

	<-quit
	log.Println("shutdown server ...")

	// exit with server deregistion
	if err = consul.client.Agent().ServiceDeregister(srv1.ID); err != nil {
		log.Printf("deregister failed with err: %s\n", err.Error())
	}
	log.Println("deregister completed ...")
}

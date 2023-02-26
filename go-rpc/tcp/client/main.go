package main

import (
	"fmt"
	"log"
	"net/rpc"
)

type Args struct {
	X, Y int
}

func main() {
	// establish a tcp conn
	cli, err := rpc.Dial("tcp","localhost:19090")
	if err != nil {
		log.Fatalf("rpc.Dial failed...\nerr: %v\n",err)
	}

	// claim args
	args := &Args{100000, 14514}
	// claim res
	var res int // not (res *int) !
	// sync call remote func
	if err := cli.Call("ServiceA.Serve", args, &res); err != nil {
		log.Fatalf("call remote func failed...\nerr: %v\n", err)
	}
	fmt.Println("ServiceA.Serve return:", res)

	fmt.Println("------------")

	// async call
	divCall := cli.Go("ServiceA.Serve", args, &res, nil) 
	replyCall := <- divCall.Done // when the call is complete, done channel will signal
	fmt.Println("async call err:",replyCall.Error)
	fmt.Println("async ServiceA.Serve return:", res)
}

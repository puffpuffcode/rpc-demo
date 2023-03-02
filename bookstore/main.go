package main

import (
	"bookstore/server"
)

func main() {
	go server.StartServer()
	server.StartHttpProxy()
}

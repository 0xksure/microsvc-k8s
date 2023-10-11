package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/err/m/shared"
)

func main() {
	logger := log.Default()
	logger.Printf("Starting service1. Waiting for requests")

	rpcServer := rpc.NewServer()
	message := new(shared.MessageServer)
	rpcServer.RegisterName("getMessage", message)
	rpcServer.HandleHTTP("/", "/debug")

	listener, err := net.Listen("tcp", ":1122")

	if err != nil {
		logger.Fatal("listen error: ", err)
	}
	defer listener.Close()
	logger.Printf("Listening on 1122/tcp")
	http.Serve(listener, nil)
}

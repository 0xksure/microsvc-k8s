package main

import (
	"log"
	"net"
	"net/http/httptest"
	"net/rpc"
	"sync"
	"testing"

	"github.com/err/shared"
)

var (
	serverAddr, newServerAddr string
	httpServerAddr            string
	once, newOnce, httpOnce   sync.Once
)

func startServer() {
	server := rpc.NewServer()
	message := new(shared.MessageServer)
	server.Register(message)

	var l net.Listener
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatalf("net.listen tcp :0: %v", err)
	}
	serverAddr = l.Addr().String()

	go server.Accept(l)
	rpc.HandleHTTP()
	httpOnce.Do(startHttpServer)
}

func startHttpServer() {
	server := httptest.NewServer(nil)
	httpServerAddr = server.Listener.Addr().String()
	log.Println("Test HTTP RPC server listening on", httpServerAddr)
}

// TestRPC is the entrypoint for the tests
// it is heavily inspired by the net/rpc package server tests
func TestRPC(t *testing.T) {
	once.Do(startServer)
	testRPC(t, serverAddr)
}

// testRPC tests the current rpc
func testRPC(t *testing.T, addr string) {
	t.Log("Connect to addr: ", addr)
	client, err := rpc.Dial("tcp", addr)
	if err != nil {
		t.Fatal("dialing", err)
	}
	defer client.Close()

	// calls
	args := shared.Args{}
	reply := new(string)
	err = client.Call("MessageServer.GetMessage", args, reply)
	if err != nil {
		t.Error("Failed to getMessage", err)
	}
	if *reply != "hello from your server" {
		t.Errorf("GetMessage: Expected `hello from your server` but got `%s`", *reply)
	}
}

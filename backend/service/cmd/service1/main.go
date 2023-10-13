package main

import (
	"database/sql"
	"errors"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"

	"github.com/err/m/shared"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func listenRPC(logger *log.Logger, err_chan chan error) {
	psqlConnStr := "postgres://user:pwd@postgres:5432/user?sslmode=disable"
	db, err := sql.Open("postgres", psqlConnStr)
	if err != nil {
		logger.Print("db open error: ", err)
		err_chan <- err
		return
	}
	defer db.Close()

	rpcServer := rpc.NewServer()

	messageServer := &shared.MessageServer{
		Db: db,
	}
	rpcServer.Register(messageServer)
	rpcServer.HandleHTTP("/", "/debug")

	listener, err := net.Listen("tcp", ":1122")

	if err != nil {
		logger.Print("listen error: ", err)
		err_chan <- err
		return
	}
	defer listener.Close()
	http.Serve(listener, nil)
	err_chan <- errors.New("listener closed")
}
func listenHttp(logger *log.Logger, err_chan chan error) {
	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	err_chan <- srv.ListenAndServe()
}

func main() {
	logger := log.Default()
	logger.Printf("Starting service1. Waiting for requests")
	rpcErr := make(chan error)
	httpErr := make(chan error)
	go listenRPC(logger, rpcErr)
	go listenHttp(logger, httpErr)

	select {
	case err := <-rpcErr:
		logger.Fatalf("existed rpc with error: %s", err.Error())
	case err := <-httpErr:
		logger.Fatalf("existed http with error: %s", err.Error())
	}

	// NOTE: in env

	logger.Printf("Listening on 1122/tcp")

}

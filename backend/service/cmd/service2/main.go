package main

import (
	"log"
	"net/http"
	"net/rpc"
	"time"

	"github.com/err/m/shared"
	"github.com/gorilla/mux"
)

type Service2 struct {
	client *rpc.Client
}

func (s *Service2) GetMessage() (*shared.Reply, error) {
	reply := new(shared.Reply)
	err := s.client.Call("MessageServer.getMessage", &shared.Args{}, reply)
	if err != nil {
		return reply, err
	}
	return reply, nil
}

// Service 2 exposes a server on port 8080
// it also uses other microservices to perform its tasks
func main() {
	logger := log.Default()
	logger.Println("Starting service2...")

	client, err := rpc.DialHTTP("tcp", "http://micro1-service:1122")
	if err != nil {
		logger.Fatal("dialing:", err)
	}
	svc2 := &Service2{
		client: client,
	}

	r := mux.NewRouter()
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	r.HandleFunc("/getMessage", func(w http.ResponseWriter, r *http.Request) {
		logger.Println("Received request for /getMessage")
		reply, err := svc2.GetMessage()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		logger.Println("Reply: ", *reply)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(*reply))
	})

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

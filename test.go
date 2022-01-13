package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type HTTPServer struct {
	http.Server
}

func main() {

	router := mux.NewRouter()
	server := &HTTPServer{
		Server: http.Server{
			Addr:    "8000",
			Handler: router,
		},
	}

	router.HandleFunc("/", handleMain)
	router.HandleFunc("/counter", handleMain)
	router.HandleFunc("/counter/get", handleMain)
	log.Println("Main Server is running!")
	done := make(chan bool)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("Listen and serve: %v", err)
		}
		done <- true
	}()

	//wait shutdown
	server.WaitShutdown()
	log.Println("Shutting down")

	<-done
	log.Printf("DONE!")
}

func handleMain(rw http.ResponseWriter, r *http.Request) {
	log.Println("main.handleMain")
	response := map[string]string{
		"message": "Welcome to test-pvc - Main",
	}
	json.NewEncoder(rw).Encode(response)
}

func (s *HTTPServer) WaitShutdown() {
	irqSig := make(chan os.Signal, 1)
	signal.Notify(irqSig, syscall.SIGINT, syscall.SIGTERM)

	//Wait interrupt or shutdown request through /shutdown
	select {
	case sig := <-irqSig:
		log.Printf("Shutdown request (signal: %v)", sig)
	case <-time.After(2 * time.Second):
		log.Printf("Its been 2 seconds")
	}
	log.Printf("Stopping http server ...")

	//shutdown the server
	err := s.Shutdown(context.Background())
	if err != nil {
		log.Printf("Shutdown request error: %v", err)
	}
}

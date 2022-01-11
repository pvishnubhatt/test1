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

type Test struct {
	name string
	zip  string
}
type HTTPServer struct {
	http.Server
	shutdownReq chan bool
	Test
}

func main() {

	router := mux.NewRouter()
	server := &HTTPServer{
		Server: http.Server{
			Addr:    ":8000",
			Handler: router,
		},
		shutdownReq: make(chan bool),
		Test: Test{
			name: "Paddy",
			zip:  "94024",
		},
	}

	router.HandleFunc("/", handleMain)
	router.HandleFunc("/counter", handleMain)
	router.HandleFunc("/counter/get", handleMain)
	router.HandleFunc("/shutdown", server.ShutdownHandler)
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
	log.Println("Shutting down ", server.name, server.zip)

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
	case sig := <-s.shutdownReq:
		log.Printf("Shutdown request (/shutdown %v)", sig)
	}

	log.Printf("Stoping http server ...")

	//Create shutdown context with 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//shutdown the server
	err := s.Shutdown(ctx)
	if err != nil {
		log.Printf("Shutdown request error: %v", err)
	}
}

func (s *HTTPServer) ShutdownHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Shutdown server"))
	go func() {
		s.shutdownReq <- true
	}()
}

package main

import (
	"httpTest/internal/request"
	"httpTest/internal/response"
	"httpTest/internal/server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const port = 42069

func main() {
	server, err := server.Serve(port, ResponseHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func ResponseHandler(w *response.Writer, req *request.Request) {
	hErr := &server.HandlerError{}

	if req.RequestLine.RequestTarget == "/yourproblem" {
		hErr.StatusCode = 400
		hErr.StatusMsg = "Your problem is not my problem\n"
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		hErr.StatusCode = 500
		hErr.StatusMsg = "Woopsie, my bad\n"
		return
	}

	w.Write([]byte("All good, frfr\n"))

}

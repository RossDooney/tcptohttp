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
	hRsp := &server.HandlerResponse{}

	if req.RequestLine.RequestTarget == "/yourproblem" {
		hRsp.StatusCode = 400
		hRsp.StatusMsg = "Your problem is not my problem\n"
		hRsp.Write(w)
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		hRsp.StatusCode = 500
		hRsp.StatusMsg = "Woopsie, my bad\n"
		hRsp.Write(w)
		return
	}

	hRsp.StatusCode = 200
	hRsp.StatusMsg = "All good, frfr\n"
	hRsp.Write(w)

}

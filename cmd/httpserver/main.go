package main

import (
	"fmt"
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

	if req.RequestLine.RequestTarget == "/yourproblem" {

		err := w.WriteStatusLine(400)
		if err != nil {
			fmt.Println()
		}

		body := []byte("Your problem is not my problem\n")
		headers := response.GetDefaultHeaders(len(body))
		headers.Set("Content-Type", "text/html")
		err = w.WriteHeaders(headers)
		if err != nil {
			fmt.Println()
		}
		w.Write(body)
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		err := w.WriteStatusLine(500)
		if err != nil {
			fmt.Println()
		}

		body := []byte("Woopsie, my bad\nn")
		headers := response.GetDefaultHeaders(len(body))
		headers.Set("Content-Type", "text/html")
		err = w.WriteHeaders(headers)
		if err != nil {
			fmt.Println()
		}
		w.Write(body)
		return
	}

	err := w.WriteStatusLine(200)
	if err != nil {
		fmt.Println()
	}

	body := []byte("All good, frfr\n")
	headers := response.GetDefaultHeaders(len(body))
	headers.Set("Content-Type", "text/html")
	err = w.WriteHeaders(headers)
	if err != nil {
		fmt.Println()
	}
	w.Write(body)

}

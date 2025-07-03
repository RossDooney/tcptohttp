package main

import (
	"fmt"
	"httpTest/internal/request"
	"httpTest/internal/response"
	"httpTest/internal/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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

	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {

		url := "https://httpbin.org/" + strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")

		fmt.Println(url)

		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		buf := make([]byte, 32)
		i, err := res.Body.Read(buf)

		fmt.Println("number of bytes ", i)
		fmt.Println("What was read: ", buf)
		res.Body.Close()

		if err != nil {
			log.Fatal(err)
		}

		// fmt.Printf("%s", body)

		err = w.WriteStatusLine(200)
		if err != nil {
			fmt.Println()
		}
		headers := response.GetDefaultHeaders(len(buf))
		delete(headers, "Content-Length")
		headers.Set("Transfer-Encoding", "chunked")
		headers.Set("Content-Type", "text/html")
		err = w.WriteHeaders(headers)
		if err != nil {
			fmt.Println()
		}
		w.Write(buf)
		return

	}

	if req.RequestLine.RequestTarget == "/yourproblem" {

		err := w.WriteStatusLine(400)
		if err != nil {
			fmt.Println()
		}

		body := []byte(`<html>
  <head>
    <title>400 Bad Request</title>
  </head>
  <body>
    <h1>Bad Request</h1>
    <p>Your request honestly kinda sucked.</p>
  </body>
</html>`)
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

		body := []byte(`<html>
<head>
	<title>500 Internal Server Error</title>
</head>
<body>
	<h1>Internal Server Error</h1>
	<p>Okay, you know what? This one is on me.</p>
</body>
</html>`)
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

	body := []byte(`<html>
<head>
	<title>200 OK</title>
</head>
<body>
	<h1>Success!</h1>
	<p>Your request was an absolute banger.</p>
</body>
</html>`)
	headers := response.GetDefaultHeaders(len(body))
	headers.Set("Content-Type", "text/html")
	err = w.WriteHeaders(headers)
	if err != nil {
		fmt.Println()
	}
	w.Write(body)

}

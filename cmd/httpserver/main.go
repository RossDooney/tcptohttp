package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"httpTest/internal/headers"
	"httpTest/internal/request"
	"httpTest/internal/response"
	"httpTest/internal/server"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

func main() {

	port := flag.Int("port", 8080, "Server port")

	flag.Parse()

	server, err := server.Serve(*port, ResponseHandler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", *port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func ResponseHandler(w *response.Writer, req *request.Request) {

	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		proxyHandler(w, req)
		return
	}

	if req.RequestLine.RequestTarget == "/yourproblem" {
		handler400(w)
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		handler500(w)
		return
	}
	if req.RequestLine.RequestTarget == "/video" {
		videoHandler(w)
		return
	}

	if req.RequestLine.RequestTarget == "/" {
		fmt.Println("Index page")
		handler200(w, "static/index.html")
		return
	}

}

func handler200(w *response.Writer, file string) {
	err := w.WriteStatusLine(200)
	if err != nil {
		fmt.Println(err)
	}

	body, err := os.ReadFile(file)
	if err != nil {
		fmt.Println(err)
	}

	headers := response.GetDefaultHeaders(len(body))
	headers.Set("Content-Type", "text/html")
	err = w.WriteHeaders(headers)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(body)
}

func handler400(w *response.Writer) {

	err := w.WriteStatusLine(400)
	if err != nil {
		fmt.Println(err)
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
	headers.Override("Content-Type", "text/html")
	err = w.WriteHeaders(headers)
	if err != nil {
		fmt.Println(err)
	}
	w.Write(body)
}

func handler500(w *response.Writer) {
	err := w.WriteStatusLine(500)
	if err != nil {
		fmt.Println(err)
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
}

func proxyHandler(w *response.Writer, req *request.Request) {
	url := "https://httpbin.org/" + strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin/")
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("error getting URL")
		return
	}
	defer resp.Body.Close()

	w.WriteStatusLine(200)
	h := response.GetDefaultHeaders(0)
	h.Set("Transfer-Encoding", "chunked")
	h.Set("Trailers", "X-Content-SHA256 X-Content-Length")
	h.Remove("Content-Length")

	w.WriteHeaders(h)

	bodySize := 0
	var fullBody []byte
	buf := make([]byte, 1024)
	for {
		i, err := resp.Body.Read(buf)
		bodySize += i
		if i > 0 {
			_, err = w.WriteChunkedBody(buf[:i])
			fullBody = append(fullBody, buf[:i]...)
			if err != nil {
				fmt.Println("Error writing chunked body:", err)
				break
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading response body:", err)
			break
		}
	}
	_, err = w.WriteChunkedBodyDone()
	if err != nil {
		fmt.Println("Error writing chunked body done:", err)
	}

	trailer := headers.NewHeaders()

	sum := sha256.Sum256(fullBody)
	strSum := hex.EncodeToString(sum[:])

	trailer.Set("X-Content-SHA256", strSum)
	trailer.Set("X-Content-Length", strconv.Itoa(bodySize))

	w.WriteTrailers(trailer)

}

func videoHandler(w *response.Writer) {

	data, err := os.ReadFile("assets/vim.mp4")
	if err != nil {
		fmt.Println(err)
		handler500(w)
		return
	}

	w.WriteStatusLine(200)
	h := response.GetDefaultHeaders(0)
	h.Override("Content-Type", "video/mp4")
	h.Remove("Content-Length")
	w.WriteHeaders(h)

	w.WriteBody(data)

}

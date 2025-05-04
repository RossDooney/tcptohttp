package main

import (
	"httpTest/internal/request"
	"strings"
)

func main() {
	_, _ = request.RequestFromReader(strings.NewReader("GET / HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n"))
}

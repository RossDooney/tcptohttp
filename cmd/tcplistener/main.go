package main

import (
	"fmt"
	"httpTest/internal/request"
	"net"
)

func main() {

	listener, err := net.Listen("tcp", ":42069")

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	defer listener.Close()

	for {
		conn, err := listener.Accept()

		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		req, err := request.RequestFromReader(conn)

		if err != nil {
			fmt.Println("Error: ", err)
			return
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)
		fmt.Println("Headers: ")
		for key, value := range req.Headers {
			fmt.Printf("- %s: %s\n", key, value)
		}
	}
}

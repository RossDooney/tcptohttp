package request

import (
	"fmt"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var httpMethods = []string{"GET", "HEAD", "OPTIONS", "PUT", "DELETE", "POST"}

func RequestFromReader(reader io.Reader) (*Request, error) {

	httpRequest, err := io.ReadAll(reader)
	if err != nil {

		return nil, err
	}

	header, err := parseRequestLine(string(httpRequest))
	if err != nil {
		return nil, err
	}

	return header, nil
}

func parseRequestLine(httpRequest string) (*Request, error) {

	var header Request

	if !strings.Contains(httpRequest, "\r\n") {
		return nil, fmt.Errorf("incirrect request line format")
	}

	requestLine := strings.Split(httpRequest, "\r\n")[0]

	if len(strings.Split(requestLine, " ")) != 3 {
		return nil, fmt.Errorf("incorrect request line format: %s", requestLine)
	}

	parts := strings.Split(requestLine, " ")

	httpVersion := strings.Split(parts[2], "/")

	if httpVersion[0] != "HTTP" {
		return nil, fmt.Errorf("invalid HTTP version: %s", httpVersion[0])
	}
	if httpVersion[1] != "1.1" {
		return nil, fmt.Errorf("unsupported HTPP version: %s", httpVersion[1])
	}

	header.RequestLine.HttpVersion = httpVersion[1]
	header.RequestLine.RequestTarget = parts[1]

	for _, method := range httpMethods {
		if method == parts[0] {
			header.RequestLine.Method = parts[0]
			return &header, nil
		}
	}

	return nil, fmt.Errorf("invalid method: %s", parts[0])
}

// func parseRequestString(str string) (RequestLine, error){

// }

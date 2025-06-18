package response

import (
	"httpTest/internal/headers"
	"io"
)

type ServerStatusCode int

const (
	StatusOK          ServerStatusCode = 200
	StatusBadRequest  ServerStatusCode = 400
	StatusServerError ServerStatusCode = 500
)

var statusCode = map[ServerStatusCode]string{
	StatusOK:          "Ok",
	StatusBadRequest:  "Bad Request",
	StatusServerError: "Internal Server Error",
}

func WriteStatusLine(w io.Writer, statusCode ServerStatusCode) error {

	return nil
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {

	return nil
}

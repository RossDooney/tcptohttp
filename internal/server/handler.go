package server

import (
	"httpTest/internal/request"
	"httpTest/internal/response"
	"io"
	"net"
)

type HandlerError struct {
	StatusCode int
	StatusMsg  string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	err := response.WriteStatusLine(conn, 200)

	if err != nil {
		return
	}

	headers := response.GetDefaultHeaders(0)
	err = response.WriteHeaders(conn, headers)

	if err != nil {
		return
	}

}

package server

import (
	"bytes"
	"fmt"
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

func (handleErr HandlerError) Write(w io.Writer) {
	errMsg := fmt.Sprintf("HTTP/1.1 %d %s\r\n", handleErr.StatusCode, handleErr.StatusMsg)

	w.Write([]byte(errMsg))
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)

	if err != nil {
		handleErr := &HandlerError{
			StatusCode: 400,
			StatusMsg:  err.Error(),
		}
		handleErr.Write(conn)
		return
	}

	buf := bytes.NewBuffer([]byte{})
	handleErr := s.handler(buf, req)

	if handleErr != nil {
		handleErr.Write(conn)
	}

	b := buf.Bytes()
	err = response.WriteStatusLine(conn, 200)

	if err != nil {
		return
	}

	headers := response.GetDefaultHeaders(len(b))
	err = response.WriteHeaders(conn, headers)
	if err != nil {
		return
	}

}

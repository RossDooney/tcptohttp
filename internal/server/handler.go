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

func (handleErr HandlerError) Write(w io.Writer, body []byte) {
	errMsg := fmt.Sprintf("HTTP/1.1 %d %s\r\n", handleErr.StatusCode, handleErr.StatusMsg)
	w.Write([]byte(errMsg))
	headers := response.GetDefaultHeaders(len(body))
	err := response.WriteHeaders(w, headers)
	if err != nil {
		return
	}

	w.Write(body)
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)

	if err != nil {
		handleErr := &HandlerError{
			StatusCode: 400,
			StatusMsg:  err.Error(),
		}
		handleErr.Write(conn, nil)
		return
	}

	buffer := bytes.NewBuffer([]byte{})
	handleErr := s.handler(buffer, req)
	body := buffer.Bytes()
	if handleErr != nil {
		handleErr.Write(conn, body)
		return
	}

	err = response.WriteStatusLine(conn, 200)
	if err != nil {
		return
	}

	headers := response.GetDefaultHeaders(len(body))
	err = response.WriteHeaders(conn, headers)
	if err != nil {
		return
	}

	conn.Write(body)
}

package server

import (
	"bytes"
	"httpTest/internal/request"
	"httpTest/internal/response"
	"io"
	"net"
)

type HandlerError struct {
	StatusCode response.ServerStatusCode
	StatusMsg  string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func (handleErr HandlerError) Write(w io.Writer) {
	response.WriteStatusLine(w, handleErr.StatusCode)
	body := []byte(handleErr.StatusMsg)
	headers := response.GetDefaultHeaders(len(body))
	response.WriteHeaders(w, headers)
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
		handleErr.Write(conn)
		return
	}

	buffer := bytes.NewBuffer([]byte{})
	handleErr := s.handler(buffer, req)
	if handleErr != nil {
		handleErr.Write(conn)
		return
	}
	body := buffer.Bytes()
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

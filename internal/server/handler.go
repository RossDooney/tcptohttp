package server

import (
	"bytes"
	"httpTest/internal/request"
	"httpTest/internal/response"
	"net"
)

type HandlerError struct {
	StatusCode response.ServerStatusCode
	StatusMsg  string
}

type Handler func(w *response.Writer, req *request.Request)

func (handleErr HandlerError) Write(w *response.Writer) {
	w.WriteStatusLine(handleErr.StatusCode)
	body := []byte(handleErr.StatusMsg)
	headers := response.GetDefaultHeaders(len(body))
	w.WriteHeaders(headers)
	w.Write(body)
}

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()

	req, err := request.RequestFromReader(conn)

	w := &response.Writer{
		Writer: conn,
	}

	if err != nil {
		handleErr := &HandlerError{
			StatusCode: 400,
			StatusMsg:  err.Error(),
		}
		handleErr.Write(w)
		return
	}

	buffer := bytes.NewBuffer([]byte{})
	s.handler(w, req)

	body := buffer.Bytes()
	err = w.WriteStatusLine(210)
	if err != nil {
		return
	}

	headers := response.GetDefaultHeaders(len(body))
	err = w.WriteHeaders(headers)
	if err != nil {
		return
	}

	conn.Write(body)
}

package server

import (
	"fmt"
	"httpTest/internal/request"
	"httpTest/internal/response"
	"net"
)

type HandlerResponse struct {
	StatusCode response.ServerStatusCode
	StatusMsg  string
}

type Handler func(w *response.Writer, req *request.Request)

func (handleRsp HandlerResponse) Write(w *response.Writer) {
	err := w.WriteStatusLine(handleRsp.StatusCode)
	if err != nil {
		fmt.Println(err)
		return
	}
	body := []byte(handleRsp.StatusMsg)
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
		handleRsp := &HandlerResponse{
			StatusCode: 400,
			StatusMsg:  err.Error(),
		}
		handleRsp.Write(w)
		return
	}

	s.handler(w, req)
}

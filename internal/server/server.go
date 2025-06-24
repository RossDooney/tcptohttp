package server

import (
	"fmt"
	"httpTest/internal/response"
	"log"
	"net"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	state    atomic.Bool
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		listener: listener,
	}
	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	s.state.Store(false)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if !s.state.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}

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

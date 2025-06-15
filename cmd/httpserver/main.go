package server

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	Listener net.Listener
}

const port = 42069

func main() {
	server, err := server.Serve(port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

func Serve(port int) (*Server, error) {

	listener, err := net.Listen("tcp", ":"+string(port))
	if err != nil {
		return nil, err
	}

	s := &Server{
		Listener: listener,
	}
	return s, nil
}

func (s *Server) Close() error {

}

func (s *Server) listen() {

}

func (s *Server) handle(conn net.Conn) {

}

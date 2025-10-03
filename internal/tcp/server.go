package tcp

import (
	"github.com/Lekuruu/go-chat/internal/logging"

	"fmt"
	"net"
)

type Server struct {
	Host   string
	Port   int
	Logger logging.Logger

	listener       net.Listener
	requestHandler func(net.Conn)
}

func NewServer(host string, port int, logger logging.Logger, handler func(net.Conn)) *Server {
	return &Server{
		Host:           host,
		Port:           port,
		Logger:         logger,
		requestHandler: handler,
	}
}

func (server *Server) Bind() string {
	return fmt.Sprintf("%s:%d", server.Host, server.Port)
}

func (server *Server) Run() error {
	listener, err := net.Listen("tcp", server.Bind())
	if err != nil {
		return err
	}
	server.listener = listener
	defer server.listener.Close()

	for {
		conn, err := server.listener.Accept()
		if err != nil {
			return err
		}
		go server.requestHandler(conn)
	}
}

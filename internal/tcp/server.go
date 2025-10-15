package tcp

import (
	"github.com/Lekuruu/go-chat/internal/logging"

	"fmt"
	"net"
)

type Server struct {
	Name   string
	Host   string
	Port   int
	Logger *logging.Logger

	listener       net.Listener
	requestHandler func(net.Conn)
}

func NewServer(name string, host string, port int, handler func(net.Conn)) *Server {
	return &Server{
		Name:           name,
		Host:           host,
		Port:           port,
		Logger:         logging.CreateLogger(name, logging.DEBUG),
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
	defer listener.Close()
	server.listener = listener
	server.Logger.Infof("Listening on '%s' ...", server.Bind())

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go server.requestHandler(conn)
	}
}

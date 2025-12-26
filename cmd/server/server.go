package main

import (
	"net"

	"github.com/Lekuruu/go-chat/internal/tcp"
)

type ChatServer struct {
	*tcp.Server
	Clients           map[string]*Client
	EncryptionKey     []byte
	Version           uint8
	RequireEncryption bool
}

func NewChatServer(host string, port int, key []byte, handler func(net.Conn)) *ChatServer {
	// Create base server from tcp package
	tcpServer := tcp.NewServer("chat-server", host, port, handler)

	return &ChatServer{
		Clients:           make(map[string]*Client),
		Server:            tcpServer,
		EncryptionKey:     key,
		RequireEncryption: true,
		Version:           1,
	}
}

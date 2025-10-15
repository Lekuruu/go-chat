package main

import (
	"net"

	"github.com/Lekuruu/go-chat/internal/logging"
	"github.com/Lekuruu/go-chat/internal/protocol"
)

type Client struct {
	Name            string
	Conn            net.Conn
	Server          *ChatServer
	Logger          *logging.Logger
	Encryption      protocol.EncryptionType
	IsAuthenticated bool
}

func (c *Client) Close() error {
	return c.Conn.Close()
}

func (c *Client) Address() string {
	return c.Conn.RemoteAddr().String()
}

func (c *Client) ReadPacket() (*protocol.Packet, error) {
	return protocol.DeserializePacket(c.Conn)
}

func (c *Client) SendPacket(packet *protocol.Packet) error {
	packet.Version = c.Server.Version
	packet.Encryption = c.Encryption
	return packet.Serialize(c.Conn, c.Server.EncryptionKey)
}

func (c *Client) SendError(e *ChatError) error {
	packet, err := e.Packet()
	if err != nil {
		return err
	}
	return c.SendPacket(packet)
}

func NewClient(conn net.Conn, server *ChatServer) *Client {
	address := conn.RemoteAddr().String()
	logger := logging.CreateLogger(address, server.Logger.GetLevel())

	return &Client{
		Name:            "",
		Conn:            conn,
		Server:          server,
		Logger:          logger,
		Encryption:      protocol.EncryptionTypeAES,
		IsAuthenticated: false,
	}
}

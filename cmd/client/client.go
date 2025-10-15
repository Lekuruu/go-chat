package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"net"

	"github.com/Lekuruu/go-chat/internal/logging"
	"github.com/Lekuruu/go-chat/internal/protocol"
)

type ChatClient struct {
	Host          string
	Port          int
	EncryptionKey []byte
	Version       uint8
	Name          string
	Encryption    protocol.EncryptionType

	Conn   net.Conn
	Logger *logging.Logger
	UI     *ChatUI

	IsAuthenticated bool
	Challenge       []byte
}

func NewChatClient(host string, port int, key []byte) *ChatClient {
	logger := logging.CreateLogger("client", logging.INFO)

	return &ChatClient{
		Host:            host,
		Port:            port,
		Logger:          logger,
		Encryption:      protocol.EncryptionTypeNone,
		EncryptionKey:   key,
		Version:         1,
		IsAuthenticated: false,
	}
}

func (c *ChatClient) Bind() string {
	return net.JoinHostPort(c.Host, fmt.Sprintf("%d", c.Port))
}

func (c *ChatClient) Address() string {
	return c.Conn.RemoteAddr().String()
}

func (c *ChatClient) AddSystemMessage(format string, args ...interface{}) {
	if c.UI != nil {
		c.UI.AddSystemMessage(format, args...)
	} else {
		c.Logger.Infof(format, args...)
	}
}

func (c *ChatClient) ShowDisconnectMessage() {
	if c.UI != nil {
		c.UI.ShowDisconnectMessage("Connection lost!")
	}
}

func (c *ChatClient) ReadPacket() (*protocol.Packet, error) {
	return protocol.DeserializePacket(c.Conn, c.EncryptionKey)
}

func (c *ChatClient) SendPacket(packet *protocol.Packet) error {
	packet.Version = c.Version
	packet.Encryption = c.Encryption
	return packet.Serialize(c.Conn, c.EncryptionKey)
}

func (c *ChatClient) SendChallenge() error {
	challenge := make([]byte, 16)
	if _, err := rand.Read(challenge); err != nil {
		return err
	}

	c.Challenge = challenge
	challengeData := protocol.Challenge{Data: challenge}

	data, err := challengeData.ToBytes()
	if err != nil {
		return err
	}

	packet := &protocol.Packet{
		Id:   protocol.PacketIdChallenge,
		Data: data,
	}

	// Send unencrypted challenge packet to server
	return c.SendPacket(packet)
}

func (c *ChatClient) SendNickname(nickname string) error {
	nicknameString := protocol.String{Value: nickname}
	data, err := nicknameString.ToBytes()
	if err != nil {
		return err
	}

	packet := &protocol.Packet{
		Id:   protocol.PacketIdNickname,
		Data: data,
	}

	return c.SendPacket(packet)
}

func (c *ChatClient) SendMessage(content string) error {
	message := protocol.Message{
		Sender:  c.Name,
		Content: content,
	}

	buffer := new(bytes.Buffer)
	if err := message.Serialize(buffer); err != nil {
		return err
	}

	packet := &protocol.Packet{
		Id:   protocol.PacketIdMessage,
		Data: buffer.Bytes(),
	}

	return c.SendPacket(packet)
}

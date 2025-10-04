package main

import (
	"bytes"

	"github.com/Lekuruu/go-chat/internal/protocol"
)

type ChatError struct {
	Code    uint16
	Message string
}

func (e *ChatError) Error() string {
	return e.Message
}

func (e *ChatError) Packet() (*protocol.Packet, error) {
	data := protocol.Error{
		Code:    e.Code,
		Message: e.Message,
	}

	buffer := new(bytes.Buffer)
	err := data.Serialize(buffer)
	if err != nil {
		return nil, err
	}

	return &protocol.Packet{
		Id:   protocol.PacketIdError,
		Data: buffer.Bytes(),
	}, nil
}

func NewChatError(code uint16, message string) *ChatError {
	return &ChatError{
		Code:    code,
		Message: message,
	}
}

var (
	ErrInvalidPacket = NewChatError(1, "Received an invalid packet. Please try again!")
	ErrUnknownPacket = NewChatError(2, "Received an unknown packet. Ensure that the target server is up to date.")
)

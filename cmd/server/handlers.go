package main

import (
	"github.com/Lekuruu/go-chat/internal/protocol"
)

var AuthHandlers = make(map[protocol.PacketId]func(*protocol.Packet, *Client))
var MainHandlers = make(map[protocol.PacketId]func(*protocol.Packet, *Client))

// TODO: Add handlers

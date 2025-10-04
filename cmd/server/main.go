package main

import (
	"net"
)

func main() {
	var server *ChatServer
	var key []byte

	// Example key, we should replace this later
	key = []byte("A0KWJW3qRCiYcEj3")

	connectionHandler := func(conn net.Conn) { handleConnection(conn, server) }
	server = NewChatServer("localhost", 8080, key, connectionHandler)
	server.Run()
}

func handleConnection(conn net.Conn, server *ChatServer) {
	defer conn.Close()

	// Create client instance
	client := NewClient(conn, server)

	// Authentication stage
	for {
		packet, err := client.ReadPacket()
		if err != nil {
			client.SendError(ErrInvalidPacket)
			return
		}

		handler, ok := AuthHandlers[packet.Id]
		if !ok {
			client.SendError(ErrUnknownPacket)
			return
		}

		handler(packet, client)
		if client.IsAuthenticated {
			break
		}
	}

	// Add client to server map & remove on disconnect
	server.Clients[client.Name] = client
	defer delete(server.Clients, client.Name)

	// TODO: broadcast join
	// TODO: defer quit broadcast

	// Main communication loop
	for {
		packet, err := client.ReadPacket()
		if err != nil {
			client.SendError(ErrInvalidPacket)
			return
		}

		handler, ok := MainHandlers[packet.Id]
		if !ok {
			client.SendError(ErrUnknownPacket)
			return
		}

		handler(packet, client)
	}
}

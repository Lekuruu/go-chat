package main

import (
	"fmt"
	"net"
)

func main() {
	var server *ChatServer
	var config *Config

	// Ensure config file exists
	created, err := EnsureConfig(DefaultConfigFilename)
	if err != nil {
		fmt.Printf("Failed to create config file: %v\n", err)
		return
	}
	if created {
		fmt.Println("Created default config file 'server.json'.")
	}

	config, err = ReadConfig(DefaultConfigFilename)
	if err != nil {
		fmt.Printf("Failed to read config file: %v\n", err)
		return
	}

	connectionHandler := func(conn net.Conn) { handleConnection(conn, server) }
	server = NewChatServer(config.ServerHost, config.ServerPort, config.SecretKey, connectionHandler)
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
			client.Logger.Errorf("Failed to read authentication packet: %v", err)
			client.SendError(ErrInvalidPacket)
			return
		}

		handler, ok := AuthHandlers[packet.Id]
		if !ok {
			client.Logger.Errorf("Unknown authentication packet ID: %v", packet.Id)
			client.SendError(ErrUnknownPacket)
			return
		}

		handler(packet, client)
		if client.IsAuthenticated {
			break
		}
	}

	// Change logger name to client's username
	client.Logger.Infof("Client authenticated with username: '%s'", client.Name)
	client.Logger.SetName(client.Name)

	// Add client to server map & remove on disconnect
	server.Clients[client.Name] = client
	defer delete(server.Clients, client.Name)

	// Broadcast join
	broadcastJoin(client)
	defer broadcastQuit(client)

	// Main communication loop
	for {
		packet, err := client.ReadPacket()
		if err != nil && err.Error() == "EOF" {
			client.Logger.Infof("Client disconnected")
			return
		}
		if err != nil {
			client.Logger.Errorf("Failed to read packet: %v", err)
			client.SendError(ErrInvalidPacket)
			return
		}

		handler, ok := MainHandlers[packet.Id]
		if !ok {
			client.Logger.Errorf("Unknown packet ID: %v", packet.Id)
			client.SendError(ErrUnknownPacket)
			return
		}

		handler(packet, client)
	}
}

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// Example key, we should replace this later
var key = []byte("A0KWJW3qRCiYcEj3")

func main() {
	client := NewChatClient("localhost", 8080, key)
	conn, err := net.Dial("tcp", client.Bind())
	if err != nil {
		client.Logger.Errorf("Failed to connect to server: %v", err)
		return
	}
	defer conn.Close()

	client.Conn = conn
	client.Logger.Infof("Connected to %s", client.Address())

	if err := handleAuthentication(client); err != nil {
		client.Logger.Errorf("Authentication failed: %v", err)
		client.Logger.WaitForInput()
		return
	}

	if !client.IsAuthenticated {
		// Client failed the authentication and has
		// already seen an error message most likely
		client.Logger.WaitForInput()
		return
	}

	client.UI = NewChatUI(func(content string) {
		if err := client.SendMessage(content); err != nil {
			client.Logger.Errorf("Failed to send message: %v", err)
		}
	})

	// Handle all incoming packets in the background
	go handlePackets(client)

	if err := client.UI.Run(); err != nil {
		client.Logger.Errorf("UI error: %v", err)
		client.Logger.WaitForInput()
	}
}

func handleAuthentication(client *ChatClient) error {
	reader := bufio.NewReader(os.Stdin)

	if err := client.SendChallenge(); err != nil {
		return fmt.Errorf("failed to send challenge: %w", err)
	}

	packet, err := client.ReadPacket()
	if err != nil {
		return fmt.Errorf("failed to read challenge response: %w", err)
	}

	handler, ok := AuthHandlers[packet.Id]
	if !ok {
		return fmt.Errorf("received unexpected packet during authentication")
	}
	handler(packet, client)

	fmt.Print("Enter your nickname: ")
	nickname, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read nickname: %w", err)
	}

	// Trim spaces & newline characters from nickname
	nickname = strings.TrimSpace(nickname)
	nickname = strings.Trim(nickname, "\n")

	// Handle possible carriage return for Windows
	if len(nickname) > 1 && strings.HasSuffix(nickname, "\r") {
		nickname = strings.TrimSuffix(nickname, "\r")
	}
	if len(nickname) == 0 {
		return fmt.Errorf("nickname cannot be empty")
	}

	client.Name = nickname
	if err := client.SendNickname(nickname); err != nil {
		return fmt.Errorf("failed to send nickname: %w", err)
	}

	// We now expect either an acknowledgment or an error packet
	packet, err = client.ReadPacket()
	if err != nil {
		return fmt.Errorf("failed to read authentication response: %w", err)
	}

	handler, ok = AuthHandlers[packet.Id]
	if !ok {
		return fmt.Errorf("unexpected packet during authentication")
	}
	handler(packet, client)
	return nil
}

func handlePackets(client *ChatClient) {
	for {
		packet, err := client.ReadPacket()
		if err != nil {
			client.Logger.Errorf("Failed to read packet: %v", err)
			return
		}

		handler, ok := MainHandlers[packet.Id]
		if !ok {
			client.Logger.Warningf("Unknown packet ID: %v", packet.Id)
			continue
		}

		handler(packet, client)
	}
}

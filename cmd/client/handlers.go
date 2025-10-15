package main

import (
	"bytes"

	"github.com/Lekuruu/go-chat/internal/protocol"
)

var AuthHandlers = make(map[protocol.PacketId]func(*protocol.Packet, *ChatClient))
var MainHandlers = make(map[protocol.PacketId]func(*protocol.Packet, *ChatClient))

func init() {
	AuthHandlers[protocol.PacketIdError] = handleError
	AuthHandlers[protocol.PacketIdChallenge] = handleChallenge
	AuthHandlers[protocol.PacketIdNicknameAck] = handleNicknameAck

	MainHandlers[protocol.PacketIdError] = handleError
	MainHandlers[protocol.PacketIdNames] = handleNames
	MainHandlers[protocol.PacketIdJoin] = handleJoin
	MainHandlers[protocol.PacketIdQuit] = handleQuit
	MainHandlers[protocol.PacketIdMessage] = handleMessage
}

func handleError(packet *protocol.Packet, client *ChatClient) {
	var err protocol.Error
	buffer := bytes.NewBuffer(packet.Data)

	if deserializeErr := err.Deserialize(buffer); deserializeErr != nil {
		client.Logger.Errorf("Failed to deserialize error: %v", deserializeErr)
		return
	}

	client.AddSystemMessage("Error [%d]: %s", err.Code, err.Message)
}

func handleChallenge(packet *protocol.Packet, client *ChatClient) {
	var challenge protocol.Challenge
	buffer := bytes.NewBuffer(packet.Data)

	if err := challenge.Deserialize(buffer); err != nil {
		client.Logger.Errorf("Failed to deserialize challenge: %v", err)
		return
	}

	if !bytes.Equal(challenge.Data, client.Challenge) {
		client.AddSystemMessage("Authentication failed: Challenge mismatch")
		return
	}

	client.Logger.Info("Challenge verified successfully")
	client.Encryption = protocol.EncryptionTypeAES
}

func handleNicknameAck(packet *protocol.Packet, client *ChatClient) {
	client.Logger.Info("Nickname acknowledged")
	client.IsAuthenticated = true
}

func handleNames(packet *protocol.Packet, client *ChatClient) {
	var userList protocol.UserList
	buffer := bytes.NewBuffer(packet.Data)

	if err := userList.Deserialize(buffer); err != nil {
		client.Logger.Errorf("Failed to deserialize user list: %v", err)
		return
	}

	users := make([]string, 0, len(userList.Users))
	for _, user := range userList.Users {
		users = append(users, user.Name)
	}

	if client.UI == nil {
		client.Logger.Warning("UI is not initialized, cannot update user list")
		return
	}

	client.UI.SetUsers(users)
}

func handleJoin(packet *protocol.Packet, client *ChatClient) {
	var user protocol.User
	buffer := bytes.NewBuffer(packet.Data)

	if err := user.Deserialize(buffer); err != nil {
		client.Logger.Errorf("Failed to deserialize user: %v", err)
		return
	}

	client.AddSystemMessage("%s joined the chat", user.Name)

	if client.UI != nil {
		client.UI.AddUser(user.Name)
	}
}

func handleQuit(packet *protocol.Packet, client *ChatClient) {
	var user protocol.User
	buffer := bytes.NewBuffer(packet.Data)

	if err := user.Deserialize(buffer); err != nil {
		client.Logger.Errorf("Failed to deserialize user: %v", err)
		return
	}

	client.AddSystemMessage("%s left the chat", user.Name)

	if client.UI != nil {
		client.UI.RemoveUser(user.Name)
	}
}

func handleMessage(packet *protocol.Packet, client *ChatClient) {
	var message protocol.Message
	buffer := bytes.NewBuffer(packet.Data)

	if err := message.Deserialize(buffer); err != nil {
		client.Logger.Errorf("Failed to deserialize message: %v", err)
		return
	}

	if client.UI == nil {
		client.Logger.Warning("UI is not initialized, cannot display message")
		return
	}

	client.UI.AddMessage(message.Sender, message.Content)
}

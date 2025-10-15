package main

import (
	"bytes"

	"github.com/Lekuruu/go-chat/internal/protocol"
)

var AuthHandlers = make(map[protocol.PacketId]func(*protocol.Packet, *Client))
var MainHandlers = make(map[protocol.PacketId]func(*protocol.Packet, *Client))

func init() {
	AuthHandlers[protocol.PacketIdChallenge] = handleAuthChallenge
	AuthHandlers[protocol.PacketIdNickname] = handleNickname
	MainHandlers[protocol.PacketIdMessage] = handleMessage
}

func handleAuthChallenge(packet *protocol.Packet, client *Client) {
	var challenge protocol.Challenge

	if err := challenge.FromBytes(packet.Data); err != nil {
		client.Logger.Errorf("Failed to deserialize challenge: %v", err)
		client.SendError(ErrInvalidPacket)
		return
	}

	client.Logger.Debugf("Received challenge with %d bytes", len(challenge.Data))

	encryptedChallenge := protocol.Challenge{Data: challenge.Data}
	responseBuffer := new(bytes.Buffer)

	if err := encryptedChallenge.Serialize(responseBuffer); err != nil {
		client.Logger.Errorf("Failed to serialize challenge response: %v", err)
		client.SendError(ErrInvalidPacket)
		return
	}

	response := &protocol.Packet{
		Id:   protocol.PacketIdChallenge,
		Data: responseBuffer.Bytes(),
	}
	client.Encryption = protocol.EncryptionTypeAES

	if err := client.SendPacket(response); err != nil {
		client.Logger.Errorf("Failed to send challenge response: %v", err)
	}
}

func handleNickname(packet *protocol.Packet, client *Client) {
	var nicknameString protocol.String

	if err := nicknameString.FromBytes(packet.Data); err != nil {
		client.Logger.Errorf("Failed to read nickname: %v", err)
		client.SendError(ErrInvalidPacket)
		return
	}
	nickname := nicknameString.Value

	if _, exists := client.Server.Clients[nickname]; exists {
		client.Logger.Warningf("Nickname already in use: %s", nickname)
		client.SendError(ErrNicknameInUse)
		return
	}

	client.Name = nickname
	client.IsAuthenticated = true

	nicknameAck := &protocol.Packet{Id: protocol.PacketIdNicknameAck}
	if err := client.SendPacket(nicknameAck); err != nil {
		client.Logger.Errorf("Failed to send nickname acknowledgement: %v", err)
		// Try to continue either way - stuff may go wrong though
	}

	// Send list of existing users to client
	users := make([]protocol.User, 0, len(client.Server.Clients))
	for name := range client.Server.Clients {
		users = append(users, protocol.User{Name: name})
	}

	userList := protocol.UserList{Users: users}
	data, err := userList.ToBytes()
	if err != nil {
		client.Logger.Errorf("Failed to serialize user list: %v", err)
		return
	}

	namePacket := &protocol.Packet{
		Id:   protocol.PacketIdNames,
		Data: data,
	}

	if err := client.SendPacket(namePacket); err != nil {
		client.Logger.Errorf("Failed to send user list: %v", err)
	}
}

func handleMessage(packet *protocol.Packet, client *Client) {
	var message protocol.Message

	if err := message.FromBytes(packet.Data); err != nil {
		client.Logger.Errorf("Failed to deserialize message: %v", err)
		client.SendError(ErrInvalidPacket)
		return
	}

	client.Logger.Infof(
		"Message from %s: %s",
		message.Sender, message.Content,
	)

	messageBuffer := new(bytes.Buffer)
	if err := message.Serialize(messageBuffer); err != nil {
		client.Logger.Errorf("Failed to serialize message: %v", err)
		return
	}

	broadcastPacket := &protocol.Packet{
		Id:   protocol.PacketIdMessage,
		Data: messageBuffer.Bytes(),
	}

	for _, targetClient := range client.Server.Clients {
		if err := targetClient.SendPacket(broadcastPacket); err != nil {
			client.Logger.Errorf("Failed to send message to %s: %v", targetClient.Name, err)
		}
	}
}

func broadcastJoin(client *Client) {
	user := protocol.User{Name: client.Name}
	data, err := user.ToBytes()
	if err != nil {
		client.Logger.Errorf("Failed to serialize join: %v", err)
		return
	}

	packet := &protocol.Packet{
		Id:   protocol.PacketIdJoin,
		Data: data,
	}

	for _, targetClient := range client.Server.Clients {
		if targetClient.Name == client.Name {
			continue
		}
		if err := targetClient.SendPacket(packet); err != nil {
			client.Logger.Errorf("Failed to send join to %s: %v", targetClient.Name, err)
		}
	}
}

func broadcastQuit(client *Client) {
	user := protocol.User{Name: client.Name}
	data, err := user.ToBytes()
	if err != nil {
		client.Logger.Errorf("Failed to serialize quit: %v", err)
		return
	}

	packet := &protocol.Packet{
		Id:   protocol.PacketIdQuit,
		Data: data,
	}

	for _, targetClient := range client.Server.Clients {
		if targetClient.Name == client.Name {
			continue
		}
		if err := targetClient.SendPacket(packet); err != nil {
			client.Logger.Errorf("Failed to send quit to %s: %v", targetClient.Name, err)
		}
	}
}

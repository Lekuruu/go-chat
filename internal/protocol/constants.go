package protocol

type PacketId uint16
type EncryptionType uint8

const (
	PacketIdError PacketId = iota
	PacketIdChallenge
	PacketIdNickname
	PacketIdNicknameAck
	PacketIdNames
	PacketIdJoin
	PacketIdQuit
	PacketIdMessage
)

const (
	EncryptionTypeNone EncryptionType = iota
	EncryptionTypeAES
)

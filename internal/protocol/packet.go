package protocol

import (
	"io"
)

type Packet struct {
	Version    uint8
	Id         PacketId
	Encryption EncryptionType
	Data       []byte
}

func NewPacket(version uint8, id PacketId, encryption EncryptionType, data []byte) *Packet {
	return &Packet{
		Version:    version,
		Id:         id,
		Encryption: encryption,
		Data:       data,
	}
}

func (packet *Packet) Serialize(writer io.Writer, key []byte) error {
	// Write packet header
	if err := writeUint8(writer, packet.Version); err != nil {
		return err
	}
	if err := writeUint16(writer, uint16(packet.Id)); err != nil {
		return err
	}
	if err := writeUint8(writer, uint8(packet.Encryption)); err != nil {
		return err
	}

	// Encrypt data if needed
	outgoing := packet.outgoingData(key)

	// Write the packet body
	if err := writeUint32(writer, uint32(len(outgoing))); err != nil {
		return err
	}
	_, err := writer.Write(outgoing)
	return err
}

func (packet *Packet) outgoingData(key []byte) []byte {
	switch packet.Encryption {
	case EncryptionTypeAES:
		// We only support AES encryption for now
		encryptedData, err := Encrypt(packet.Data, key)
		if err != nil {
			return nil
		}
		return encryptedData
	default:
		// Use no encryption by default
		return packet.Data
	}
}

func DeserializePacket(reader io.Reader, key []byte) (*Packet, error) {
	version, err := readUint8(reader)
	if err != nil {
		return nil, err
	}

	packetId, err := readUint16(reader)
	if err != nil {
		return nil, err
	}

	encryptionType, err := readUint8(reader)
	if err != nil {
		return nil, err
	}

	dataLength, err := readUint32(reader)
	if err != nil {
		return nil, err
	}

	data := make([]byte, dataLength)
	_, err = io.ReadFull(reader, data)
	if err != nil {
		return nil, err
	}

	packet := &Packet{
		Version:    version,
		Id:         PacketId(packetId),
		Encryption: EncryptionType(encryptionType),
		Data:       data,
	}

	// Handle data decryption
	processedData, err := handleIncomingData(packet, key)
	if err != nil {
		return nil, err
	}

	packet.Data = processedData
	return packet, nil
}

func handleIncomingData(packet *Packet, key []byte) ([]byte, error) {
	switch packet.Encryption {
	case EncryptionTypeAES:
		// We only support AES encryption for now
		decryptedData, err := Decrypt(packet.Data, key)
		if err != nil {
			return nil, err
		}
		return decryptedData, nil
	default:
		// Use no encryption by default
		return packet.Data, nil
	}
}

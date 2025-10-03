package protocol

import (
	"encoding/binary"
	"io"
)

func writeUint64(w io.Writer, v uint64) error {
	return binary.Write(w, binary.LittleEndian, v)
}

func writeInt64(w io.Writer, v int64) error {
	return writeUint64(w, uint64(v))
}

func writeUint32(w io.Writer, v uint32) error {
	return binary.Write(w, binary.LittleEndian, v)
}

func writeInt32(w io.Writer, v int32) error {
	return writeUint32(w, uint32(v))
}

func writeUint16(w io.Writer, v uint16) error {
	return binary.Write(w, binary.LittleEndian, v)
}

func writeInt16(w io.Writer, v int16) error {
	return writeUint16(w, uint16(v))
}

func writeUint8(w io.Writer, v uint8) error {
	return binary.Write(w, binary.LittleEndian, v)
}

func writeInt8(w io.Writer, v int8) error {
	return writeUint8(w, uint8(v))
}

func writeBoolean(w io.Writer, v bool) error {
	return binary.Write(w, binary.LittleEndian, v)
}

func writeFloat32(w io.Writer, v float32) error {
	return binary.Write(w, binary.LittleEndian, v)
}

func writeFloat64(w io.Writer, v float64) error {
	return binary.Write(w, binary.LittleEndian, v)
}

func writeString(w io.Writer, v string) error {
	if err := writeUint16(w, uint16(len(v))); err != nil {
		return err
	}
	_, err := w.Write([]byte(v))
	return err
}

func readUint64(r io.Reader) (v uint64, err error) {
	err = binary.Read(r, binary.LittleEndian, &v)
	return v, err
}

func readInt64(r io.Reader) (v int64, err error) {
	uv, err := readUint64(r)
	return int64(uv), err
}

func readUint32(r io.Reader) (v uint32, err error) {
	err = binary.Read(r, binary.LittleEndian, &v)
	return v, err
}

func readInt32(r io.Reader) (v int32, err error) {
	uv, err := readUint32(r)
	return int32(uv), err
}

func readUint16(r io.Reader) (v uint16, err error) {
	err = binary.Read(r, binary.LittleEndian, &v)
	return v, err
}

func readInt16(r io.Reader) (v int16, err error) {
	uv, err := readUint16(r)
	return int16(uv), err
}

func readUint8(r io.Reader) (v uint8, err error) {
	err = binary.Read(r, binary.LittleEndian, &v)
	return v, err
}

func readInt8(r io.Reader) (v int8, err error) {
	uv, err := readUint8(r)
	return int8(uv), err
}

func readBoolean(r io.Reader) (v bool, err error) {
	err = binary.Read(r, binary.LittleEndian, &v)
	return v, err
}

func readFloat32(r io.Reader) (v float32, err error) {
	err = binary.Read(r, binary.LittleEndian, &v)
	return v, err
}

func readFloat64(r io.Reader) (v float64, err error) {
	err = binary.Read(r, binary.LittleEndian, &v)
	return v, err
}

func readString(r io.Reader) (string, error) {
	length, err := readUint16(r)
	if err != nil {
		return "", err
	}

	buf := make([]byte, length)
	_, err = io.ReadFull(r, buf)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

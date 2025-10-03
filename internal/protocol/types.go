package protocol

import "io"

// Serializable is an interface for types that can
// be serialized to and deserialized from a reader/writer.
type Serializable interface {
	Serialize(io.Writer) error
	Deserialize(io.Reader) error
}

type Message struct {
	Serializable
	Sender  string
	Content string
}

func (m *Message) Serialize(w io.Writer) error {
	if err := writeString(w, m.Sender); err != nil {
		return err
	}
	if err := writeString(w, m.Content); err != nil {
		return err
	}
	return nil
}

func (m *Message) Deserialize(r io.Reader) (err error) {
	m.Sender, err = readString(r)
	if err != nil {
		return err
	}
	m.Content, err = readString(r)
	if err != nil {
		return err
	}
	return nil
}

type User struct {
	Serializable
	Name string
}

func (u *User) Serialize(w io.Writer) error {
	return writeString(w, u.Name)
}

func (u *User) Deserialize(r io.Reader) (err error) {
	u.Name, err = readString(r)
	return err
}

type UserList struct {
	Serializable
	Users []User
}

func (ul *UserList) Serialize(w io.Writer) error {
	if err := writeUint32(w, uint32(len(ul.Users))); err != nil {
		return err
	}
	for _, user := range ul.Users {
		if err := user.Serialize(w); err != nil {
			return err
		}
	}
	return nil
}

func (ul *UserList) Deserialize(r io.Reader) (err error) {
	var length uint32
	length, err = readUint32(r)
	if err != nil {
		return err
	}

	ul.Users = make([]User, 0, length)
	for i := uint32(0); i < length; i++ {
		var user User
		if err = user.Deserialize(r); err != nil {
			return err
		}
		ul.Users = append(ul.Users, user)
	}
	return nil
}

type Challenge struct {
	Serializable
	Data []byte
}

func (c *Challenge) Serialize(w io.Writer) error {
	if err := writeUint16(w, uint16(len(c.Data))); err != nil {
		return err
	}
	if _, err := w.Write(c.Data); err != nil {
		return err
	}
	return nil
}

func (c *Challenge) Deserialize(r io.Reader) (err error) {
	length, err := readUint16(r)
	if err != nil {
		return err
	}

	c.Data = make([]byte, length)
	if _, err := io.ReadFull(r, c.Data); err != nil {
		return err
	}
	return nil
}

type Error struct {
	Serializable
	Code    uint16
	Message string
}

func (e *Error) Serialize(w io.Writer) error {
	if err := writeUint16(w, e.Code); err != nil {
		return err
	}
	if err := writeString(w, e.Message); err != nil {
		return err
	}
	return nil
}

func (e *Error) Deserialize(r io.Reader) (err error) {
	if e.Code, err = readUint16(r); err != nil {
		return err
	}
	if e.Message, err = readString(r); err != nil {
		return err
	}
	return nil
}

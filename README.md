# go-chat

A terminal-based chat application written in Go, featuring encrypted communication using AES-GCM symmetric encryption. This project was created for a school research paper exploring symmetric encryption in real-time communication systems.

## Features

- **Secure Communication**: All messages are encrypted using AES-GCM
- **Modern Terminal UI**: Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) for an interactive experience
- **Real-time Messaging**: Instant message broadcasting to all connected clients
- **User Management**: Join/quit notifications and live user list

## Requirements

- Go 1.23 or higher
- A terminal with support for ANSI escape codes (modern terminals on Windows, macOS, Linux)

## Installation

Clone the repository:

```bash
git clone https://github.com/Lekuruu/go-chat.git
cd go-chat
```

### Running the Server

Start the chat server on `localhost:8080`:

```bash
go run ./cmd/server
```

The server will start listening for incoming connections and log all activity to the console.

### Running the Client

In a separate terminal, start one or more clients:

```bash
go run ./cmd/client
```

When prompted:
1. Enter your desired nickname
2. Start chatting! Type your message and press enter to send
3. Press `Ctrl+C` or `Esc` to quit

### Building Executables

To build standalone executables:

```bash
# Build server
go build -o server ./cmd/server

# Build client
go build -o client ./cmd/client
```

Then run them directly:
```bash
./server
./client
```

Note that if you are on windows, you will receive `.exe` files.

## Configuration

Currently, the encryption key is hardcoded in both the client and server (`A0KWJW3qRCiYcEj3`). Both must use the same 16-byte key for successful authentication.

To change the key, modify the `key` variable in:
- `cmd/server/main.go`
- `cmd/client/main.go`

This will later on be changed into configuration files.

## Protocol

Every packet that is sent between the client and the server uses the following structure:

| Description    | Type |
|:-------------- | :--- |
| Version        | u8   |
| PacketId       | u16  |
| EncryptionType | u8   |
| CipherLength   | u32  |
| CipherData     | x    |

The decrypted cipher will contain variable data depending on the given packet ID. Depending on the encryption type, the data will either be fully encrypted or not at all. The encryption type should also ensure that different encryption standards could be used in the future.  
Right now, only AES-GCM will be supported. This means that both the client and the server will have to use a shared secret key, which will be specified inside a configuration file.

### Authentication

When a client wants to connect to a remote server, it will send a challenge request packet, to ensure that the server is using the same key.  
The challenge packet will be unencrypted with a random set of data. The server will then proceed to encrypt the data with its secret key, and send it back, which the client can then use to validate the data by decrypting it and comparing it to the previously sent data.

Once that is done, the client will be prompted for a nickname, which is then sent to the server. If the username is already taken, the server will send back an error, indicating that the name is already taken by someone else. If the nickname is available, the client is now successfully authenticated and ready to start messaging other users.

### Messaging

Clients can send message requests to the server, which will then be validated and broadcasted back to other users.
The message type contains the sender itself and the message content.

### User Listing

Similar to a regular IRC server, the server will send a list of users who are currently online, once a client authenticates. Including that, it will also send a join & quit packet to each authenticated client, if a join/quit event occurs.

## Types

### String

| Description    | Type |
|:-------------- | :--- |
| Length         | u16  |
| Content        |      |

---

### Challenge

| Description    | Type   |
|:-------------- | :----- |
| Length         | u16    |
| Random Bytes   |        |

---

### Error

| Description    | Type   |
|:-------------- | :----- |
| Code           | u16    |
| Message        | String |

---

### User

| Description    | Type   |
|:-------------- | :----- |
| Username       | String |

---

### Message

| Description    | Type   |
|:-------------- | :----- |
| Username       | String |
| Content        | String |

---

### User Listing

| Description    | Type   |
|:-------------- | :----- |
| Length         | u32    |
| Users          | User[] |

package protocol

import (
	"flag"
	"testing"
)

var encryptionKey = flag.String("encryption_key", "A0KWJW3qRCiYcEj3", "Key used for encryption/decryption")

func TestEncryption(t *testing.T) {
	encodedKey := []byte(*encryptionKey)
	plaintext := []byte("Hello, World!")

	ciphertext, err := Encrypt(plaintext, encodedKey)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
		return
	}

	decryptedText, err := Decrypt(ciphertext, encodedKey)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
		return
	}

	if string(decryptedText) != string(plaintext) {
		t.Fatalf(
			"Decrypted text does not match original: got %q, want %q",
			string(decryptedText), string(plaintext),
		)
	}
}

package crypto

import "testing"

func TestEncryptDecrypt(t *testing.T) {
	key := make([]byte, 32)
	for i := range key {
		key[i] = byte(i)
	}

	plaintext := []byte("super secret data")

	nonce, cipher, err := Encrypt(key, plaintext, 32)
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	out, err := Decrypt(key, nonce, cipher, 32)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if string(out) != string(plaintext) {
		t.Fatalf("Expected %q, got %q", plaintext, out)
	}
}

func TestDecryptWithWrongKeyFails(t *testing.T) {
	key := make([]byte, 32)
	badKey := make([]byte, 32)

	// change the key
	key[0] = 2

	nonce, cipher, _ := Encrypt(key, []byte("secret"), 32)

	if _, err := Decrypt(badKey, nonce, cipher, 32); err == nil {
		t.Fatal("Decrypt should fail with wrong key")
	}
}


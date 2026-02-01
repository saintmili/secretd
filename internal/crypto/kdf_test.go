package crypto

import "testing"

var (
	time       = uint32(3)
	memory     = uint32(65536)
	threads    = uint8(4)
	keyLength  = uint32(32)
	saltLength = int(16)
)

func TestDeriveKeys(t *testing.T) {
	password := []byte("correct horse battery staple")
	salt, _ := GenerateSalt(16)

	encKey, macKey, err := DeriveKeys(password, salt, time, memory, threads, keyLength, saltLength)
	if err != nil {
		t.Fatalf("DeriveKeys failed: %v", err)
	}

	if len(encKey) != 32 {
		t.Fatalf("Expected encKey length 32, got %d", len(encKey))
	}
	if len(macKey) != 32 {
		t.Fatalf("Expected macKey length 32, got %d", len(macKey))
	}

	if string(encKey) == string(macKey) {
		t.Fatal("encKey and macKey must differ")
	}
}

func TestDeriveKeysEmptyPasswordFails(t *testing.T) {
	salt, _ := GenerateSalt(16)

	if _, _, err := DeriveKeys([]byte{}, salt, time, memory, threads, keyLength, saltLength); err == nil {
		t.Fatal("Expected error for empty password")
	}
}

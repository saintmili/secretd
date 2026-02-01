package crypto

import (
	"crypto/rand"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

// create a cryptographically secure random salt
func GenerateSalt(saltLength int) ([]byte, error) {
	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return nil, ErrFailedGenerateSalt
	}
	return salt, nil
}

// derives a secret key from a password and salt using Argon2id
func DeriveKeys(password []byte, salt []byte, time, memory uint32, threads uint8, keyLength uint32, saltLength int) (encKey, macKey []byte, err error) {
	if len(password) == 0 {
		return nil, nil, ErrEmptyPassword
	}
	if len(salt) < saltLength {
		return nil, nil, ErrInvalidSalt
	}

	// we have two keys, enc and mac
	totalKeyLength := 2 * keyLength

	key := argon2.IDKey(
		password,
		salt,
		time,
		memory,
		threads,
		totalKeyLength,
	)

	encKey = make([]byte, keyLength)
	macKey = make([]byte, keyLength)

	copy(encKey, key[:keyLength])
	copy(macKey, key[keyLength:totalKeyLength])

	// wipe master key
	for i := range key {
		key[i] = 0
	}

	return encKey, macKey, nil
}

// encode salt for storage
func EncodeSalt(salt []byte) string {
	return base64.StdEncoding.EncodeToString(salt)
}

// decode salt from storage
func DecodeSalt(encoded string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(encoded)
}

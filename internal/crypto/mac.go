package crypto

import (
	"crypto/hmac"
	"crypto/sha256"
)

// compute HMAC over vault critical fields
func ComputeMAC(macKey, nonce, ciphertext []byte) []byte {
	h := hmac.New(sha256.New, macKey)
	h.Write(nonce)
	h.Write(ciphertext)
	return h.Sum(nil)
}

// verify vault integrity
func VerifyMAC(macKey, expectedMAC, nonce, ciphertext []byte) bool {
	actual := ComputeMAC(macKey, nonce, ciphertext)
	return hmac.Equal(actual, expectedMAC)
}


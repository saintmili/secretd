package crypto

import "testing"

func TestMACVerify(t *testing.T) {
	key := make([]byte, 32)
	data := []byte("ciphertext")
	nonce := []byte("nonce")

	mac := ComputeMAC(key, nonce, data)

	if !VerifyMAC(key, mac, nonce, data) {
		t.Fatal("MAC verification failed")
	}
}

func TestMACDetectsTampering(t *testing.T) {
	key := make([]byte, 32)
	data := []byte("ciphertext")
	nonce := []byte("nonce")

	mac := ComputeMAC(key, nonce, data)

	data[0] ^= 0xff // tamper

	if VerifyMAC(key, mac, nonce, data) {
		t.Fatal("MAC verification should fail on tampering")
	}
}


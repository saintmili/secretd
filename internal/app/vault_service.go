package app

import (
	"github.com/saintmili/secretd/internal/config"
	"github.com/saintmili/secretd/internal/crypto"
)

type VaultService struct {
	ArgonMemory      uint32
	ArgonIterations  uint32
	ArgonParallelism uint8
	KeyLength        uint32
	SaltLength       int
}

func NewVaultService(cfg *config.Config) *VaultService {
	return &VaultService{
		ArgonMemory:      uint32(cfg.Security.Argon2Memory),
		ArgonIterations:  uint32(cfg.Security.Argon2Time),
		ArgonParallelism: uint8(cfg.Security.Argon2Threads),
		KeyLength:        uint32(cfg.Security.KeyLength),
		SaltLength:       cfg.Security.SaltLength,
	}
}

func (v *VaultService) DeriveKeys(password, salt []byte) ([]byte, []byte, error) {
	return crypto.DeriveKeys(
		password,
		salt,
		v.ArgonIterations,
		v.ArgonMemory,
		v.ArgonParallelism,
		v.KeyLength,
		v.SaltLength,
	)
}

func (v *VaultService) Decrypt(key, nonce, cipherText []byte) ([]byte, error) {
	return crypto.Decrypt(key, nonce, cipherText, int(v.KeyLength))

}
func (v *VaultService) Encrypt(key, plainText []byte) ([]byte, []byte, error) {
	return crypto.Encrypt(key, plainText, int(v.KeyLength))
}

func (v *VaultService) ReadPassword(prompt string) ([]byte, error) {
	return crypto.ReadPassword(prompt)
}

func (v *VaultService) GenerateSalt() ([]byte, error) {
	return crypto.GenerateSalt(v.SaltLength)
}

func (v *VaultService) EncodeSalt(salt []byte) string {
	return crypto.EncodeSalt(salt)
}

// decode salt from storage
func (v *VaultService) DecodeSalt(encoded string) ([]byte, error) {
	return crypto.DecodeSalt(encoded)
}

func (v *VaultService) ComputeMAC(macKey, nonce, cipherText []byte) []byte {
	return crypto.ComputeMAC(macKey, nonce, cipherText)
}

func (v *VaultService) VerifyMAC(macKey, expectedMAC, nonce, cipherText []byte) bool {
	return crypto.VerifyMAC(macKey, expectedMAC, nonce, cipherText)
}

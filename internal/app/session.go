package app

import (
	"github.com/saintmili/secretd/internal/storage"
	"github.com/saintmili/secretd/internal/vault"
)

type Session struct {
	Vault  *vault.Vault
	EncKey []byte
	MacKey []byte
	File   *storage.VaultFile
}

// wipe keys from memory
func (s *Session) Close() {
	if s == nil {
		return
	}

	zeroBytes(s.EncKey)
	zeroBytes(s.MacKey)
	s.Vault.Wipe()
}


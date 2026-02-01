package app

import "errors"

var (
	ErrMismatchPassword     = errors.New("Mismatch password")
	ErrFailedReadPassword   = errors.New("Failed to read password")
	ErrVaultExists          = errors.New("Vault already exists")
	ErrInvalidPassword      = errors.New("Invalid password")
	ErrEmptyPassword        = errors.New("Password can not be empty")
	ErrFailedSerializeVault = errors.New("Failed to serialize vault")
	ErrFailedVaultIntegrity = errors.New("❌ Vault integrity check failed (tampered)")
	ErrWrongMasterPassword  = errors.New("❌ Wrong master password")
	ErrVaultCorrupted       = errors.New("Vault corrupted")
)

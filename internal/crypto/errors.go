package crypto

import "errors"

var (
	ErrFailedGenerateSalt = errors.New("Failed to generate salt")
	ErrEmptyPassword      = errors.New("Password can not be empty")
	ErrInvalidSalt        = errors.New("Invalid salt")
	ErrInvalidKeyLength   = errors.New("Invalid key length")
	ErrInvalidNonceSize   = errors.New("InvalidNonceSize")
)

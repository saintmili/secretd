package doctor

import "errors"

var (
	ErrFailedReadPassword = errors.New("Failed to read password")
)

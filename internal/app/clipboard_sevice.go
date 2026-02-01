package app

import (
	"github.com/saintmili/secretd/internal/clipboard"
)

type ClipboardService struct {
	Timeout int
}

func NewClipboardService(timeout int) *ClipboardService {
	return &ClipboardService{
		Timeout: timeout,
	}
}

func (c *ClipboardService) Copy(data string) error {
	return clipboard.Copy(data, c.Timeout)
}

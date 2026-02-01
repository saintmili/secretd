//go:build windows

package clipboard

import (
	"fmt"
	"os/exec"
	"strings"
)

func copy(data string, timeout int) error {
	if timeout <= 0 {
		timeout = 15 // safety fallback
	}

	cmd := exec.Command("cmd", "/C", "clip")
	cmd.Stdin = strings.NewReader(data)
	if err := cmd.Run(); err == nil {
		cmd := exec.Command("cmd", "/C", fmt.Sprintf("timeout /T %d >nul && echo off | clip", timeout))
		cmd.Start()
		return nil
	}
	return fmt.Errorf("clipboard failed: clip not available")
}

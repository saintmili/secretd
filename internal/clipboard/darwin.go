//go:build darwin

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

	cmd := exec.Command("pbcopy")
	cmd.Stdin = strings.NewReader(data)
	if err := cmd.Run(); err == nil {
		cmd := exec.Command("bash", "-c", fmt.Sprintf("sleep %d && echo -n | pbcopy", timeout))
		cmd.Start()
		return nil
	}
	return fmt.Errorf("clipboard failed: pbcopy not available")
}

//go:build linux

package clipboard

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func copy(data string, timeout int) error {
	if timeout <= 0 {
		timeout = 15 // safety fallback
	}

	// Wayland
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		cmd := exec.Command("wl-copy")
		cmd.Stdin = strings.NewReader(data)
		if err := cmd.Run(); err == nil {
			cmd := exec.Command("sh", "-c", fmt.Sprintf("sleep %d && wl-copy --clear", timeout))
			cmd.Start() // detached, ignore errors
			return nil
		}
	} else {
		// X11
		cmd := exec.Command("xclip", "-selection", "clipboard")
		cmd.Stdin = strings.NewReader(data)
		if err := cmd.Run(); err == nil {
			cmd := exec.Command("sh", "-c", fmt.Sprintf("sleep %d && xclip -selection clipboard /dev/null && xclip -selection primary /dev/null", timeout))
			cmd.Start() // detached, ignore errors
			return nil
		}
	}
	return fmt.Errorf("clipboard failed: install wl-clipboard (Wayland) or xclip (X11)")

}

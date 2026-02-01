package crypto

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// ReadPassword securely reads a password from TTY without echo
func ReadPassword(prompt string) ([]byte, error) {
	fmt.Fprint(os.Stderr, prompt)

	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr) // newline after input

	if err != nil {
		return nil, err
	}

	return passwordBytes, nil
}


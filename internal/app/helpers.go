package app

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
)

func readOptional(prompt, current string) string {
	val := readLine(fmt.Sprintf("%s [%s]: ", prompt, current))
	if val == "" {
		return current
	}
	return val
}

func readLine(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

// wipe data from memory
func zeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

func parseGenerateArgs(args []string) (generate bool, length int) {
	generate = false
	length = 16

	for i := range args {
		if args[i] == "--generate" {
			generate = true

			// i+a is the next arg after --generate
			if i+1 < len(args) {
				if l, err := strconv.Atoi(args[i+1]); err == nil && l > 0 {
					length = l
				}
			}
			break
		}
	}

	return
}

func generatePassword(length int) ([]byte, error) {
	pw := make([]byte, length)
	for i := range pw {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return nil, err
		}
		pw[i] = charset[n.Int64()]
	}
	return pw, nil
}

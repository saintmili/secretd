package doctor

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/saintmili/secretd/internal/config"
	"github.com/saintmili/secretd/internal/crypto"
	"github.com/saintmili/secretd/internal/storage"
)

func Run(cfg *config.Config) error {
	fmt.Println("🩺 secretd doctor")
	fmt.Println(strings.Repeat("-", 40))

	cfg, warnings, _ := config.LoadConfig()
	if len(warnings) == 0 {
		fmt.Println("✅ Config validation OK")
	} else {
		fmt.Println("⚠️  Config warnings:")
		for _, w := range warnings {
			fmt.Printf("   - %s: %s\n", w.Field, w.Message)
		}
	}

	// 1️⃣ Vault presence
	vf, err := storage.Load(cfg.Vault.Path)
	if err != nil {
		fmt.Println("❌ Vault file not found or unreadable")
		fmt.Println("   → Run: secretd init")
		return err
	}
	fmt.Println("✅ Vault file found")

	// 2️⃣ Vault structure sanity
	salt, err := crypto.DecodeSalt(vf.Salt)
	if err != nil || len(salt) < 16 {
		fmt.Println("❌ Invalid vault salt")
		return err
	}
	defer zeroBytes(salt)
	fmt.Println("✅ Salt OK")

	if len(vf.Nonce) == 0 {
		fmt.Println("❌ Missing nonce")
		return err
	}
	fmt.Println("✅ Nonce OK")

	if len(vf.Ciphertext) == 0 {
		fmt.Println("❌ Ciphertext is empty")
		return err
	}
	fmt.Println("✅ Ciphertext present")

	if len(vf.MAC) == 0 {
		fmt.Println("⚠️  Vault has no MAC (integrity protection missing)")
		fmt.Println("   → Unlock once to auto-upgrade")
	} else {
		fmt.Println("✅ MAC present")
	}

	// 3️⃣ Ask for password to verify MAC
	password, err := crypto.ReadPassword("Enter master password to verify integrity: ")
	if err != nil {
		return ErrFailedReadPassword
	}
	defer zeroBytes(password)

	encKey, macKey, err := crypto.DeriveKeys(
		password,
		salt,
		uint32(cfg.Security.Argon2Time),
		uint32(cfg.Security.Argon2Memory),
		uint8(cfg.Security.Argon2Threads),
		uint32(cfg.Security.KeyLength),
		cfg.Security.SaltLength,
	)
	if err != nil {
		fmt.Println("❌ Key derivation failed")
		return err
	}
	defer zeroBytes(encKey)
	defer zeroBytes(macKey)

	if len(vf.MAC) > 0 {
		if !crypto.VerifyMAC(macKey, vf.MAC, vf.Nonce, vf.Ciphertext) {
			fmt.Println("❌ Vault integrity check FAILED")
			fmt.Println("   → Possible tampering or wrong password")
			return err
		}
		fmt.Println("✅ Vault integrity verified")
	}

	// 4️⃣ Clipboard backend check
	fmt.Println("🔎 Clipboard check:")
	switch runtime.GOOS {
	case "linux":
		if os.Getenv("WAYLAND_DISPLAY") != "" {
			if _, err := exec.LookPath("wl-copy"); err != nil {
				fmt.Println("⚠️  Wayland detected but wl-copy not found")
				fmt.Println("   → Install: wl-clipboard")
			} else {
				fmt.Println("✅ wl-copy available (Wayland)")
			}
		} else {
			if _, err := exec.LookPath("xclip"); err != nil {
				fmt.Println("⚠️  X11 detected but xclip not found")
				fmt.Println("   → Install: xclip")
			} else {
				fmt.Println("✅ xclip available (X11)")
			}
		}
	case "darwin":
		if _, err := exec.LookPath("pbcopy"); err != nil {
			fmt.Println("⚠️  pbcopy not available")
		} else {
			fmt.Println("✅ pbcopy available")
		}
	case "windows":
		if _, err := exec.LookPath("clip"); err != nil {
			fmt.Println("⚠️  clip not available")
		} else {
			fmt.Println("✅ clip available")
		}
	}

	// 5️⃣ File permission check
	vPath := cfg.Vault.Path
	info, err := os.Stat(vPath)
	if err == nil {
		mode := info.Mode().Perm()
		if mode&0o077 != 0 {
			fmt.Printf("⚠️  Vault file permissions too open (%o)\n", mode)
			fmt.Println("   → Recommended: chmod 600")
		} else {
			fmt.Println("✅ Vault file permissions OK")
		}
	}

	fmt.Println(strings.Repeat("-", 40))
	fmt.Println("✅ Doctor check completed")
	return nil
}

// wipe data from memory
func zeroBytes(b []byte) {
	for i := range b {
		b[i] = 0
	}
}

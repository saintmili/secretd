package config

import (
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Version   int       `toml:"version"`
	Clipboard Clipboard `toml:"clipboard"`
	Vault     Vault     `toml:"vault"`
	Security  Security  `toml:"security"`
	Logging   Logging   `toml:"logging"`
}

type Clipboard struct {
	ClearAfterSeconds int    `toml:"clear_after_seconds"`
	WaylandBackend    string `toml:"wayland_backend"`
	X11Backend        string `toml:"x11_backend"`
}
type Vault struct {
	Path       string `toml:"path"`
	AutoBackup bool   `toml:"auto_backup"`
}
type Security struct {
	WipeMemory       bool `toml:"wipe_memory"`
	Argon2Time       int  `toml:"argon2_time"`
	Argon2Memory     int  `toml:"argon2_memory"`
	Argon2Threads    int  `toml:"argon2_threads"`
	SaltLength       int  `toml:"salt_length"`
	KeyLength        int  `toml:"key_length"`
	MaxFailedUnlocks int  `toml:"max_failed_unlocks"`
	LockoutSeconds   int  `toml:"lockout_seconds"`
}

type Logging struct {
	Enabled bool
	File    string
}

// returns a Config struct with default values
func DefaultConfig() *Config {
	var cfg Config

	cfg.Version = CurrentVersion

	cfg.Clipboard.ClearAfterSeconds = 15
	cfg.Clipboard.WaylandBackend = "wl-copy"
	cfg.Clipboard.X11Backend = "xclip"

	cfg.Vault.Path = filepath.Join(os.Getenv("HOME"), ".local", "share", "secretd", "vault.json")
	cfg.Vault.AutoBackup = true

	cfg.Security.WipeMemory = true
	cfg.Security.Argon2Time = 3
	cfg.Security.Argon2Memory = 65536
	cfg.Security.Argon2Threads = 4
	cfg.Security.SaltLength = 16
	cfg.Security.KeyLength = 32
	cfg.Security.MaxFailedUnlocks = 5
	cfg.Security.LockoutSeconds = 300

	cfg.Logging.Enabled = true
	cfg.Logging.File = "~/.local/share/secretd/secretd.log"

	return &cfg
}

func ConfigPath() (string, error) {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		base = filepath.Join(home, ".config")
	}

	dir := filepath.Join(base, "secretd")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}

	return filepath.Join(dir, "config.toml"), nil
}

func (c Config) MarshalPretty() (string, error) {
	data, err := toml.Marshal(c)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

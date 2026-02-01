package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

type VaultFile struct {
	Version    int    `json:"version"`
	Salt       string `json:"salt"`       // base64
	Nonce      []byte `json:"nonce"`      // raw bytes
	Ciphertext []byte `json:"ciphertext"` // raw bytes
	MAC        []byte `json:"mac"`        // raw bytes
}

// writes the encrypted vault to disk with safe permissions
func Save(vf *VaultFile, path string) error {
	tmp := path + ".tmp"

	data, err := json.MarshalIndent(vf, "", "  ")
	if err != nil {
		return err
	}

	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	if _, err := f.Write(data); err != nil {
		f.Close()
		return err
	}

	if err := f.Sync(); err != nil {
		f.Close()
		return err
	}

	defer f.Close()
	return os.Rename(tmp, path)
}

// reads the encrypted vault from disk
func Load(path string) (*VaultFile, error) {
	var vf *VaultFile

	info, err := os.Stat(path)
	if err != nil {
		return vf, err
	}
	if info.Mode().Perm()&0o077 != 0 {
		return vf, fmt.Errorf("vault file has insecure permissions")
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return vf, err
	}

	err = json.Unmarshal(data, &vf)

	return vf, err
}

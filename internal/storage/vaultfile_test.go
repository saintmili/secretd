package storage

import (
	"testing"
)

func TestSaveLoadVault(t *testing.T) {
	tmp := t.TempDir()
	filePath := tmp + "/vault.json"
	// hijack HOME so VaultPath uses temp dir
	t.Setenv("HOME", tmp)

	vf := VaultFile{
		Version:    1,
		Salt:       "salt",
		Nonce:      []byte{1, 2, 3},
		Ciphertext: []byte{4, 5, 6},
		MAC:        []byte{7, 8, 9},
	}

	if err := Save(vf, filePath); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load(filePath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Version != vf.Version {
		t.Fatal("Version mismatch")
	}
}

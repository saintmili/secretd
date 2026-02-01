package vault

import "testing"

func TestNewVault(t *testing.T) {
	v := New()

	if v.Version != 1 {
		t.Fatalf("Expected version 1, got %d", v.Version)
	}

	if len(v.Entries) != 0 {
		t.Fatal("New vault should have zero entries")
	}
}


package app

import "testing"

func TestGeneratePasswordLength(t *testing.T) {
	pw, _ := generatePassword(24)

	if len(pw) != 24 {
		t.Fatalf("Expected password length 24, got %d", len(pw))
	}
}

func TestParseGenerateArgs(t *testing.T) {
	gen, length := parseGenerateArgs([]string{"--generate", "32"})

	if !gen {
		t.Fatal("Expected generate=true")
	}
	if length != 32 {
		t.Fatalf("Expected length 32, got %d", length)
	}
}

func TestParseGenerateArgsDefault(t *testing.T) {
	gen, length := parseGenerateArgs([]string{})

	if gen {
		t.Fatal("Expected generate=false")
	}
	if length != 16 {
		t.Fatalf("Expected default length 16, got %d", length)
	}
}

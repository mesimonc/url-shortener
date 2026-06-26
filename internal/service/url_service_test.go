package service

import (
	"testing"
)

func TestGenerateCode(t *testing.T) {
	code, err := generateCode(6)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(code) != 6 {
		t.Errorf("expected code length 6, got %d", len(code))
	}
}

func TestGenerateCodeUnique(t *testing.T) {
	codes := make(map[string]bool)
	for i := 0; i < 100; i++ {
		code, err := generateCode(6)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if codes[code] {
			t.Errorf("duplicate code generated: %s", code)
		}
		codes[code] = true
	}
}

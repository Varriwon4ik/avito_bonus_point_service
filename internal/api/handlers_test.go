package api

import (
	"strings"
	"testing"
)

func TestResolveTTLDays(t *testing.T) {
	s := &Server{DefaultTTLDays: 365, MinTTLDays: 1, MaxTTLDays: 1825}

	got, err := s.resolveTTLDays(nil)
	if err != nil {
		t.Fatalf("expected nil ttl to use default, got error: %v", err)
	}
	if got != 365 {
		t.Fatalf("expected default ttl 365, got %d", got)
	}

	for _, ttl := range []int{0, 1826} {
		_, err := s.resolveTTLDays(&ttl)
		if err == nil {
			t.Fatalf("expected error for ttl %d", ttl)
		}
		if !strings.Contains(err.Error(), "ttl_days must be between 1 and 1825 days") {
			t.Fatalf("expected clear range error, got %v", err)
		}
	}
}

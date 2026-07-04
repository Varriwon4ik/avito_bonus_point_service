package data

import "testing"

func TestNormalizeTransactionLabel(t *testing.T) {
	t.Run("empty becomes no label", func(t *testing.T) {
		got, err := NormalizeTransactionLabel("   ")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "" {
			t.Fatalf("expected empty label, got %q", got)
		}
	})

	t.Run("trimmed label is preserved", func(t *testing.T) {
		got, err := NormalizeTransactionLabel("  qa-run / sprint-2  ")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "qa-run / sprint-2" {
			t.Fatalf("expected trimmed label, got %q", got)
		}
	})

	t.Run("too long label is rejected", func(t *testing.T) {
		_, err := NormalizeTransactionLabel("123456789012345678901234567890123")
		if err == nil || err.Error() != "label must be at most 32 characters" {
			t.Fatalf("expected length validation error, got %v", err)
		}
	})

	t.Run("control characters are rejected", func(t *testing.T) {
		_, err := NormalizeTransactionLabel("qa\nrun")
		if err == nil || err.Error() != "label must not contain control characters" {
			t.Fatalf("expected control-character validation error, got %v", err)
		}
	})

	t.Run("unsupported characters are rejected", func(t *testing.T) {
		_, err := NormalizeTransactionLabel("<script>")
		if err == nil || err.Error() != "label contains unsupported characters; use letters, numbers, spaces, '.', '-', '_', '/', or ':'" {
			t.Fatalf("expected unsupported-character validation error, got %v", err)
		}
	})
}

package data

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

const TransactionLabelMaxLength = 32

type InvalidLabelError struct {
	message string
}

func (e *InvalidLabelError) Error() string {
	return e.message
}

// NormalizeTransactionLabel trims a user-supplied label and validates that it
// is short, printable, and safe to surface in the UI and API responses.
func NormalizeTransactionLabel(raw string) (string, error) {
	label := strings.TrimSpace(raw)
	if label == "" {
		return "", nil
	}
	if !utf8.ValidString(label) {
		return "", &InvalidLabelError{message: "label must be valid UTF-8 text"}
	}
	if utf8.RuneCountInString(label) > TransactionLabelMaxLength {
		return "", &InvalidLabelError{message: fmt.Sprintf("label must be at most %d characters", TransactionLabelMaxLength)}
	}

	for _, r := range label {
		switch {
		case unicode.IsControl(r):
			return "", &InvalidLabelError{message: "label must not contain control characters"}
		case unicode.IsLetter(r), unicode.IsNumber(r):
		case r == ' ' || r == '-' || r == '_' || r == '.' || r == '/' || r == ':':
		default:
			return "", &InvalidLabelError{message: "label contains unsupported characters; use letters, numbers, spaces, '.', '-', '_', '/', or ':'"}
		}
	}

	return label, nil
}

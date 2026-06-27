package api

import (
	"crypto/subtle"
	"net/http"
	"strings"
)

// authorizeAdmin enforces US-07's "only authenticated admins can perform this
// operation" constraint for privileged endpoints such as manual accrual.
//
// Authentication is a static bearer token supplied in the Authorization header
// ("Authorization: Bearer <token>"). The expected token is configured out of
// band (ADMIN_API_TOKEN / -admin-token). When no token is configured the check
// is disabled and the request is allowed through; this keeps local development
// and the existing test suite working, while production deployments enable
// authentication simply by setting the token.
//
// It returns true when the request may proceed. On failure it writes a 401
// response (including a WWW-Authenticate challenge) and returns false, so
// callers should stop processing.
func (s *Server) authorizeAdmin(w http.ResponseWriter, r *http.Request) bool {
	if s.AdminToken == "" {
		return true
	}

	token, ok := bearerToken(r.Header.Get("Authorization"))
	if !ok {
		w.Header().Set("WWW-Authenticate", `Bearer realm="admin"`)
		unauthorized(w, "admin authentication required")
		return false
	}

	// Constant-time comparison avoids leaking the token via timing.
	if subtle.ConstantTimeCompare([]byte(token), []byte(s.AdminToken)) != 1 {
		w.Header().Set("WWW-Authenticate", `Bearer realm="admin"`)
		unauthorized(w, "invalid admin credentials")
		return false
	}

	return true
}

// bearerToken extracts the token from an "Authorization: Bearer <token>"
// header value. The scheme match is case-insensitive per RFC 7235.
func bearerToken(header string) (string, bool) {
	const prefix = "bearer "
	if len(header) < len(prefix) || !strings.EqualFold(header[:len(prefix)], prefix) {
		return "", false
	}
	token := strings.TrimSpace(header[len(prefix):])
	if token == "" {
		return "", false
	}
	return token, true
}

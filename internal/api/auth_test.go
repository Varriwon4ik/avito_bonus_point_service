package api

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestAccrueHandlerRequiresAdmin verifies the accrual route rejects an
// unauthenticated request before reaching the store. The store is nil here on
// purpose: the auth check must short-circuit, so it is never dereferenced.
func TestAccrueHandlerRequiresAdmin(t *testing.T) {
	s := &Server{AdminToken: "s3cret", Mux: http.NewServeMux()}
	s.routes()

	rec := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/v1/users/u1/accruals", nil)
	s.ServeHTTP(rec, r)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
	if got := rec.Header().Get("WWW-Authenticate"); got == "" {
		t.Errorf("expected WWW-Authenticate challenge header on 401")
	}
}

func TestBearerToken(t *testing.T) {
	cases := []struct {
		header    string
		wantToken string
		wantOK    bool
	}{
		{"Bearer secret", "secret", true},
		{"bearer secret", "secret", true},   // scheme is case-insensitive
		{"BEARER  secret ", "secret", true}, // surrounding whitespace trimmed
		{"Bearer ", "", false},              // empty token
		{"Basic secret", "", false},         // wrong scheme
		{"secret", "", false},               // no scheme
		{"", "", false},                     // missing header
	}
	for _, c := range cases {
		token, ok := bearerToken(c.header)
		if ok != c.wantOK || token != c.wantToken {
			t.Errorf("bearerToken(%q) = (%q, %v), want (%q, %v)", c.header, token, ok, c.wantToken, c.wantOK)
		}
	}
}

func TestAuthorizeAdmin(t *testing.T) {
	cases := []struct {
		name       string
		adminToken string
		authHeader string
		wantOK     bool
		wantStatus int
	}{
		{"auth disabled allows request", "", "", true, http.StatusOK},
		{"auth disabled ignores header", "", "Bearer whatever", true, http.StatusOK},
		{"valid token allowed", "s3cret", "Bearer s3cret", true, http.StatusOK},
		{"missing header rejected", "s3cret", "", false, http.StatusUnauthorized},
		{"wrong token rejected", "s3cret", "Bearer nope", false, http.StatusUnauthorized},
		{"wrong scheme rejected", "s3cret", "Basic s3cret", false, http.StatusUnauthorized},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := &Server{AdminToken: c.adminToken}
			rec := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/v1/users/u1/accruals", nil)
			if c.authHeader != "" {
				r.Header.Set("Authorization", c.authHeader)
			}

			ok := s.authorizeAdmin(rec, r)
			if ok != c.wantOK {
				t.Fatalf("authorizeAdmin ok = %v, want %v", ok, c.wantOK)
			}
			if !c.wantOK {
				if rec.Code != c.wantStatus {
					t.Fatalf("status = %d, want %d", rec.Code, c.wantStatus)
				}
				if got := rec.Header().Get("WWW-Authenticate"); got == "" {
					t.Errorf("expected WWW-Authenticate challenge header on 401")
				}
			}
		})
	}
}

package mailtd

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestDeriveAuthKey_KnownVectors locks the Argon2id parameters against fixed
// outputs computed with golang.org/x/crypto/argon2 (the same library the
// Mail.td backend uses). If any of (iterations, memory, parallelism,
// hashLength, salt derivation) drifts, this test fails.
func TestDeriveAuthKey_KnownVectors(t *testing.T) {
	cases := []struct {
		address  string
		password string
		want     string
	}{
		{"alice@mail.td", "password123", "2d0b5b1cd63138ba6e5b13777000e55b5dcd8ab4286f16d2fdd3aae8948c6bcf"},
		// Verifies salt = SHA256(lower(trim(address))).
		{"BOB@mail.td  ", "P@ssw0rd!", "5c35127b2175a8aadd1fbb16ccca66701d34b78f1f96e7caa51774159ac41060"},
	}
	for _, tc := range cases {
		got := DeriveAuthKey(tc.address, tc.password)
		if got != tc.want {
			t.Errorf("DeriveAuthKey(%q, %q) = %s, want %s", tc.address, tc.password, got, tc.want)
		}
		if len(got) != 64 {
			t.Errorf("expected 64-char hex, got %d chars", len(got))
		}
	}
}

// captureBody wraps the test server's last received JSON body.
type capture struct {
	body map[string]any
}

func newCaptureServer(t *testing.T, c *capture) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		c.body = nil
		if len(raw) > 0 {
			if err := json.Unmarshal(raw, &c.body); err != nil {
				t.Fatalf("server received non-JSON body: %s", raw)
			}
		}
		// Minimal well-formed responses for each endpoint we exercise.
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == "POST" && r.URL.Path == "/api/accounts":
			_, _ = w.Write([]byte(`{"id":"00000000-0000-0000-0000-000000000000","address":"alice@mail.td","auth_key":"x"}`))
		default:
			w.WriteHeader(http.StatusNoContent)
		}
	}))
}

func TestAccountsCreate_DerivesPasswordLocally(t *testing.T) {
	cap := &capture{}
	srv := newCaptureServer(t, cap)
	defer srv.Close()

	c := NewClient("test-token", WithBaseURL(srv.URL))
	pw := "password123"
	_, err := c.Accounts.Create(context.Background(), "alice@mail.td", &CreateOptions{Password: &pw})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if _, has := cap.body["password"]; has {
		t.Errorf("body must not include password field, got: %v", cap.body)
	}
	want := DeriveAuthKey("alice@mail.td", pw)
	if got, _ := cap.body["auth_key"].(string); got != want {
		t.Errorf("auth_key = %q, want %q", got, want)
	}
}

func TestAccountsCreate_AuthKeyTakesPrecedence(t *testing.T) {
	cap := &capture{}
	srv := newCaptureServer(t, cap)
	defer srv.Close()

	c := NewClient("t", WithBaseURL(srv.URL))
	ak := strings.Repeat("a", 64)
	pw := "ignored"
	_, err := c.Accounts.Create(context.Background(), "alice@mail.td", &CreateOptions{AuthKey: &ak, Password: &pw})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if got, _ := cap.body["auth_key"].(string); got != ak {
		t.Errorf("auth_key = %q, want %q", got, ak)
	}
	if _, has := cap.body["password"]; has {
		t.Errorf("body must not include password field")
	}
}

func TestAccountsResetPassword_PasswordWithEmailID(t *testing.T) {
	cap := &capture{}
	srv := newCaptureServer(t, cap)
	defer srv.Close()

	c := NewClient("t", WithBaseURL(srv.URL))
	pw := "password123"
	if err := c.Accounts.ResetPassword(context.Background(), "alice@mail.td", &ResetPasswordOptions{Password: &pw}); err != nil {
		t.Fatalf("ResetPassword: %v", err)
	}
	if _, has := cap.body["password"]; has {
		t.Errorf("body must not include password field")
	}
	want := DeriveAuthKey("alice@mail.td", pw)
	if got, _ := cap.body["auth_key"].(string); got != want {
		t.Errorf("auth_key = %q, want %q", got, want)
	}
}

func TestAccountsResetPassword_PasswordWithUUIDRequiresAddress(t *testing.T) {
	c := NewClient("t", WithBaseURL("http://unused"))
	pw := "password123"
	err := c.Accounts.ResetPassword(context.Background(), "11111111-1111-1111-1111-111111111111", &ResetPasswordOptions{Password: &pw})
	if err == nil {
		t.Fatal("expected error when password is used with UUID id and no Address")
	}
	if !strings.Contains(err.Error(), "Address is required") {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestAccountsResetPassword_PasswordWithUUIDAndAddress(t *testing.T) {
	cap := &capture{}
	srv := newCaptureServer(t, cap)
	defer srv.Close()

	c := NewClient("t", WithBaseURL(srv.URL))
	pw := "password123"
	addr := "alice@mail.td"
	err := c.Accounts.ResetPassword(context.Background(),
		"11111111-1111-1111-1111-111111111111",
		&ResetPasswordOptions{Password: &pw, Address: &addr})
	if err != nil {
		t.Fatalf("ResetPassword: %v", err)
	}
	want := DeriveAuthKey(addr, pw)
	if got, _ := cap.body["auth_key"].(string); got != want {
		t.Errorf("auth_key = %q, want %q", got, want)
	}
}

func TestAccountsLogin_DerivesPasswordLocally(t *testing.T) {
	cap := &capture{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &cap.body)
		if r.URL.Path != "/api/token" || r.Method != "POST" {
			t.Errorf("unexpected %s %s", r.Method, r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"00000000-0000-0000-0000-000000000000","address":"alice@mail.td","token":"jwt.x.y"}`))
	}))
	defer srv.Close()

	c := NewClient("", WithBaseURL(srv.URL))
	pw := "password123"
	res, err := c.Accounts.Login(context.Background(), "alice@mail.td", &LoginOptions{Password: &pw})
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if res.Token != "jwt.x.y" {
		t.Errorf("unexpected token: %q", res.Token)
	}
	if _, has := cap.body["password"]; has {
		t.Errorf("body must not include password field")
	}
	want := DeriveAuthKey("alice@mail.td", pw)
	if got, _ := cap.body["auth_key"].(string); got != want {
		t.Errorf("auth_key = %q, want %q", got, want)
	}
}

func TestAccountsLogin_AuthKeyTakesPrecedence(t *testing.T) {
	cap := &capture{}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &cap.body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"x","address":"alice@mail.td","token":"jwt"}`))
	}))
	defer srv.Close()

	c := NewClient("", WithBaseURL(srv.URL))
	ak := strings.Repeat("a", 64)
	pw := "ignored"
	_, err := c.Accounts.Login(context.Background(), "alice@mail.td", &LoginOptions{AuthKey: &ak, Password: &pw})
	if err != nil {
		t.Fatalf("Login: %v", err)
	}
	if got, _ := cap.body["auth_key"].(string); got != ak {
		t.Errorf("auth_key = %q, want %q", got, ak)
	}
}

func TestAccountsLogin_RequiresCredentials(t *testing.T) {
	c := NewClient("", WithBaseURL("http://unused"))
	_, err := c.Accounts.Login(context.Background(), "alice@mail.td", nil)
	if err == nil {
		t.Fatal("expected error when opts is nil")
	}
	_, err = c.Accounts.Login(context.Background(), "alice@mail.td", &LoginOptions{})
	if err == nil {
		t.Fatal("expected error when neither Password nor AuthKey is set")
	}
}

func TestUserResetAccountPassword_PasswordDerivedLocally(t *testing.T) {
	cap := &capture{}
	srv := newCaptureServer(t, cap)
	defer srv.Close()

	c := NewClient("t", WithBaseURL(srv.URL))
	pw := "password123"
	addr := "alice@mail.td"
	err := c.User.ResetAccountPassword(context.Background(),
		"11111111-1111-1111-1111-111111111111",
		&ResetPasswordOptions{Password: &pw, Address: &addr})
	if err != nil {
		t.Fatalf("ResetAccountPassword: %v", err)
	}
	if _, has := cap.body["password"]; has {
		t.Errorf("body must not include password field")
	}
	want := DeriveAuthKey(addr, pw)
	if got, _ := cap.body["auth_key"].(string); got != want {
		t.Errorf("auth_key = %q, want %q", got, want)
	}
}

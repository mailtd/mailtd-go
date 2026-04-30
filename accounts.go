package mailtd

import (
	"context"
	"fmt"
	"strings"
)

// AccountsResource handles account-related API calls.
type AccountsResource struct {
	client *Client
}

// ListDomains returns the available public email domains.
func (r *AccountsResource) ListDomains(ctx context.Context) ([]Domain, error) {
	var wrapper struct {
		Domains []Domain `json:"domains"`
	}
	err := r.client.request(ctx, "GET", "/api/domains", nil, &wrapper)
	return wrapper.Domains, err
}

// CreateOptions are optional parameters for creating an account.
type CreateOptions struct {
	Password *string `json:"password,omitempty"`
	AuthKey  *string `json:"auth_key,omitempty"`
}

// Create creates a new email account.
//
// If opts.AuthKey is set it is sent as-is. Otherwise, if opts.Password is set
// the SDK derives the auth_key locally via Argon2id (see DeriveAuthKey) and
// the password never leaves the client.
func (r *AccountsResource) Create(ctx context.Context, address string, opts *CreateOptions) (*CreateAccountResult, error) {
	body := map[string]any{"address": address}
	if opts != nil {
		switch {
		case opts.AuthKey != nil:
			body["auth_key"] = *opts.AuthKey
		case opts.Password != nil:
			body["auth_key"] = DeriveAuthKey(address, *opts.Password)
		}
	}
	var result CreateAccountResult
	err := r.client.request(ctx, "POST", "/api/accounts", body, &result)
	return &result, err
}

// LoginOptions are credentials for AccountsResource.Login.
//
// Exactly one of Password or AuthKey must be set. AuthKey takes precedence
// when both are supplied.
type LoginOptions struct {
	Password *string
	AuthKey  *string
}

// LoginResult is the response from AccountsResource.Login.
type LoginResult struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	Token   string `json:"token"`
}

// Login authenticates a mailbox and returns a mailbox-scoped JWT.
//
// When opts.Password is supplied the SDK derives the auth_key locally via
// Argon2id (see DeriveAuthKey); the password never leaves the client.
//
// The returned Token grants access to /api/accounts/{id}/* endpoints when
// addressed by UUID. Use the Token with a fresh Client:
//
//	res, _ := client.Accounts.Login(ctx, addr, &mailtd.LoginOptions{Password: &pw})
//	mb := mailtd.NewClient(res.Token)
//	msgs, _ := mb.Messages.List(ctx, res.ID, nil)
func (r *AccountsResource) Login(ctx context.Context, address string, opts *LoginOptions) (*LoginResult, error) {
	if opts == nil || (opts.AuthKey == nil && opts.Password == nil) {
		return nil, fmt.Errorf("mailtd: Login requires Password or AuthKey")
	}
	body := map[string]any{"address": address}
	if opts.AuthKey != nil {
		body["auth_key"] = *opts.AuthKey
	} else {
		body["auth_key"] = DeriveAuthKey(address, *opts.Password)
	}
	var result LoginResult
	err := r.client.request(ctx, "POST", "/api/token", body, &result)
	return &result, err
}

// Get returns an account by ID. id accepts a UUID or email address.
func (r *AccountsResource) Get(ctx context.Context, id string) (*AccountInfo, error) {
	var result AccountInfo
	err := r.client.request(ctx, "GET", fmt.Sprintf("/api/accounts/%s", id), nil, &result)
	return &result, err
}

// Delete removes an account by ID. id accepts a UUID or email address.
func (r *AccountsResource) Delete(ctx context.Context, id string) error {
	return r.client.request(ctx, "DELETE", fmt.Sprintf("/api/accounts/%s", id), nil, nil)
}

// ResetPasswordOptions specifies the new credentials.
//
// When Password is supplied the SDK derives the auth_key locally; the
// derivation needs the mailbox's email address. If id is already an email
// address it is used directly; otherwise (UUID id) Address must be set.
type ResetPasswordOptions struct {
	Password *string `json:"password,omitempty"`
	AuthKey  *string `json:"auth_key,omitempty"`
	Address  *string `json:"-"`
}

// ResetPassword resets an account's password. id accepts a UUID or email address.
func (r *AccountsResource) ResetPassword(ctx context.Context, id string, opts *ResetPasswordOptions) error {
	body, err := buildResetPasswordBody(id, opts)
	if err != nil {
		return err
	}
	return r.client.request(ctx, "PUT", fmt.Sprintf("/api/accounts/%s/reset-password", id), body, nil)
}

// buildResetPasswordBody is shared by AccountsResource.ResetPassword and
// UserResource.ResetAccountPassword. It enforces the local-derivation rule.
func buildResetPasswordBody(id string, opts *ResetPasswordOptions) (map[string]any, error) {
	body := map[string]any{}
	if opts == nil {
		return body, nil
	}
	if opts.AuthKey != nil {
		body["auth_key"] = *opts.AuthKey
		return body, nil
	}
	if opts.Password != nil {
		address := ""
		switch {
		case opts.Address != nil && *opts.Address != "":
			address = *opts.Address
		case strings.Contains(id, "@"):
			address = id
		default:
			return nil, fmt.Errorf("mailtd: ResetPasswordOptions.Address is required when id is a UUID and Password is used")
		}
		body["auth_key"] = DeriveAuthKey(address, *opts.Password)
	}
	return body, nil
}

package mailtd

import (
	"context"
	"fmt"
)

// AccountsResource handles account-related API calls.
type AccountsResource struct {
	client *Client
}

// ListDomains returns the available public email domains.
func (r *AccountsResource) ListDomains(ctx context.Context) ([]Domain, error) {
	var result []Domain
	err := r.client.request(ctx, "GET", "/api/domains", nil, &result)
	return result, err
}

// CreateOptions are optional parameters for creating an account.
type CreateOptions struct {
	Password *string `json:"password,omitempty"`
	AuthKey  *string `json:"auth_key,omitempty"`
}

// Create creates a new email account.
func (r *AccountsResource) Create(ctx context.Context, address string, opts *CreateOptions) (*CreateAccountResult, error) {
	body := map[string]any{"address": address}
	if opts != nil {
		if opts.Password != nil {
			body["password"] = *opts.Password
		}
		if opts.AuthKey != nil {
			body["auth_key"] = *opts.AuthKey
		}
	}
	var result CreateAccountResult
	err := r.client.request(ctx, "POST", "/api/accounts", body, &result)
	return &result, err
}

// LoginOptions specifies credentials for logging in.
type LoginOptions struct {
	Password *string `json:"password,omitempty"`
	AuthKey  *string `json:"auth_key,omitempty"`
}

// Login authenticates with an existing account.
func (r *AccountsResource) Login(ctx context.Context, address string, opts *LoginOptions) (*LoginResult, error) {
	body := map[string]any{"address": address}
	if opts != nil {
		if opts.Password != nil {
			body["password"] = *opts.Password
		}
		if opts.AuthKey != nil {
			body["auth_key"] = *opts.AuthKey
		}
	}
	var result LoginResult
	err := r.client.request(ctx, "POST", "/api/token", body, &result)
	return &result, err
}

// Get returns an account by ID.
func (r *AccountsResource) Get(ctx context.Context, id string) (*AccountInfo, error) {
	var result AccountInfo
	err := r.client.request(ctx, "GET", fmt.Sprintf("/api/accounts/%s", id), nil, &result)
	return &result, err
}

// Delete removes an account by ID.
func (r *AccountsResource) Delete(ctx context.Context, id string) error {
	return r.client.request(ctx, "DELETE", fmt.Sprintf("/api/accounts/%s", id), nil, nil)
}

// ResetPasswordOptions specifies the new credentials.
type ResetPasswordOptions struct {
	Password *string `json:"password,omitempty"`
	AuthKey  *string `json:"auth_key,omitempty"`
}

// ResetPassword resets an account's password.
func (r *AccountsResource) ResetPassword(ctx context.Context, id string, opts *ResetPasswordOptions) error {
	body := map[string]any{}
	if opts != nil {
		if opts.Password != nil {
			body["password"] = *opts.Password
		}
		if opts.AuthKey != nil {
			body["auth_key"] = *opts.AuthKey
		}
	}
	return r.client.request(ctx, "PUT", fmt.Sprintf("/api/accounts/%s/reset-password", id), body, nil)
}

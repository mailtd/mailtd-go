package mailtd

import (
	"context"
	"fmt"
	"net/url"
)

// UserResource handles Pro user management API calls.
type UserResource struct {
	client *Client
}

// GetMe returns the authenticated user's profile.
func (r *UserResource) GetMe(ctx context.Context) (*ProUser, error) {
	var result ProUser
	err := r.client.request(ctx, "GET", "/api/user/me", nil, &result)
	return &result, err
}

// ListAccountsPage returns a page of accounts with cursor-based pagination.
func (r *UserResource) ListAccountsPage(ctx context.Context, cursor string) (*AccountListResult, error) {
	path := "/api/user/accounts"
	if cursor != "" {
		path += "?cursor=" + url.QueryEscape(cursor)
	}
	var result AccountListResult
	err := r.client.request(ctx, "GET", path, nil, &result)
	return &result, err
}

// ListAccounts returns the first page of accounts belonging to the user.
//
// Deprecated: Use ListAccountsPage for cursor-based pagination.
func (r *UserResource) ListAccounts(ctx context.Context) (*AccountListResult, error) {
	return r.ListAccountsPage(ctx, "")
}

// DeleteAccount removes a user account by ID. id accepts a UUID or email address.
func (r *UserResource) DeleteAccount(ctx context.Context, id string) error {
	return r.client.request(ctx, "DELETE", fmt.Sprintf("/api/user/accounts/%s", id), nil, nil)
}

// ResetAccountPassword resets a user account's password. id accepts a UUID or email address.
func (r *UserResource) ResetAccountPassword(ctx context.Context, id string, opts *ResetPasswordOptions) error {
	body := map[string]any{}
	if opts != nil {
		if opts.Password != nil {
			body["password"] = *opts.Password
		}
		if opts.AuthKey != nil {
			body["auth_key"] = *opts.AuthKey
		}
	}
	return r.client.request(ctx, "PUT", fmt.Sprintf("/api/user/accounts/%s/reset-password", id), body, nil)
}

// UserListMessagesOptions are optional parameters for listing user account messages.
type UserListMessagesOptions struct {
	Page int
}

// ListAccountMessages returns messages for a user account. id accepts a UUID or email address.
func (r *UserResource) ListAccountMessages(ctx context.Context, id string, opts *UserListMessagesOptions) (*MessageListResult, error) {
	path := fmt.Sprintf("/api/user/accounts/%s/messages", id)
	if opts != nil {
		path = addPageParam(path, opts.Page)
	}
	var result MessageListResult
	err := r.client.request(ctx, "GET", path, nil, &result)
	return &result, err
}

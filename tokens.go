package mailtd

import (
	"context"
	"fmt"
)

// TokensResource handles Pro API token management.
type TokensResource struct {
	client *Client
}

// List returns all API tokens.
func (r *TokensResource) List(ctx context.Context) ([]Token, error) {
	var wrapper struct {
		Tokens []Token `json:"tokens"`
	}
	err := r.client.request(ctx, "GET", "/api/user/tokens", nil, &wrapper)
	return wrapper.Tokens, err
}

// Create generates a new API token.
func (r *TokensResource) Create(ctx context.Context, name string) (*Token, error) {
	body := map[string]string{"name": name}
	var result Token
	err := r.client.request(ctx, "POST", "/api/user/tokens", body, &result)
	return &result, err
}

// Revoke disables an API token.
func (r *TokensResource) Revoke(ctx context.Context, id string) error {
	return r.client.request(ctx, "DELETE", fmt.Sprintf("/api/user/tokens/%s", id), nil, nil)
}

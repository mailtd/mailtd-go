package mailtd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// AccountsResource handles account-related API calls.
type AccountsResource struct {
	client          *Client
	cachedDifficulty int
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
// For free users (no token set), proof-of-work is computed automatically.
// If the server requests a higher difficulty, the client retries once.
func (r *AccountsResource) Create(ctx context.Context, address string, opts *CreateOptions) (*CreateAccountResult, error) {
	// Pro users with a token skip PoW entirely.
	if r.client.token != "" {
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

	// Normalize address for PoW — server verifies against lowercased form.
	powAddress := strings.ToLower(strings.TrimSpace(address))

	// Free user: solve PoW locally, starting from cached difficulty.
	difficulty := r.cachedDifficulty
	if difficulty < defaultDifficulty {
		difficulty = defaultDifficulty
	}
	pow := SolvePow(powAddress, difficulty)
	result, retry, err := r.createWithPow(ctx, address, opts, &pow)
	if err != nil {
		return nil, err
	}
	if retry == nil {
		if result.SuggestedNextDifficulty > 0 {
			r.cachedDifficulty = result.SuggestedNextDifficulty
		}
		return result, nil
	}

	// Server asked for higher difficulty — re-solve once.
	r.cachedDifficulty = retry.RequiredDifficulty
	pow2 := SolvePow(powAddress, retry.RequiredDifficulty)
	pow2.Token = retry.Token
	result, _, err = r.createWithPow(ctx, address, opts, &pow2)
	if err == nil && result.SuggestedNextDifficulty > 0 {
		r.cachedDifficulty = result.SuggestedNextDifficulty
	}
	return result, err
}

// createWithPow posts a create-account request with a PoW solution.
// It returns the account result on success, or a retry response if the server
// demands a higher difficulty.
func (r *AccountsResource) createWithPow(ctx context.Context, address string, opts *CreateOptions, pow *PoWSolution) (*CreateAccountResult, *powRetryResponse, error) {
	body := map[string]any{"address": address}
	if opts != nil {
		if opts.Password != nil {
			body["password"] = *opts.Password
		}
		if opts.AuthKey != nil {
			body["auth_key"] = *opts.AuthKey
		}
	}
	body["pow"] = pow

	b, err := json.Marshal(body)
	if err != nil {
		return nil, nil, fmt.Errorf("mailtd: marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", r.client.baseURL+"/api/accounts", bytes.NewReader(b))
	if err != nil {
		return nil, nil, fmt.Errorf("mailtd: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := r.client.httpClient.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("mailtd: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		apiErr := &APIError{Status: resp.StatusCode}
		if err := json.NewDecoder(resp.Body).Decode(apiErr); err != nil {
			return nil, nil, &APIError{Status: resp.StatusCode, Code: "unknown", Message: resp.Status}
		}
		return nil, nil, apiErr
	}

	// Read the response body to check for retry.
	var raw json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, nil, fmt.Errorf("mailtd: decode response: %w", err)
	}

	// Check if this is a retry response.
	var retry powRetryResponse
	if err := json.Unmarshal(raw, &retry); err == nil && retry.Status == "retry" {
		return nil, &retry, nil
	}

	var result CreateAccountResult
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, nil, fmt.Errorf("mailtd: decode response: %w", err)
	}
	return &result, nil, nil
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

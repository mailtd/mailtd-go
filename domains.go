package mailtd

import (
	"context"
	"fmt"
)

// DomainsResource handles Pro custom domain API calls.
type DomainsResource struct {
	client *Client
}

// List returns all custom domains.
func (r *DomainsResource) List(ctx context.Context) ([]ProDomain, error) {
	var result []ProDomain
	err := r.client.request(ctx, "GET", "/api/user/domains", nil, &result)
	return result, err
}

// Create adds a new custom domain.
func (r *DomainsResource) Create(ctx context.Context, domain string) (*CreateDomainResult, error) {
	body := map[string]string{"domain": domain}
	var result CreateDomainResult
	err := r.client.request(ctx, "POST", "/api/user/domains", body, &result)
	return &result, err
}

// Verify triggers verification for a domain.
func (r *DomainsResource) Verify(ctx context.Context, id string) (*VerifyDomainResult, error) {
	var result VerifyDomainResult
	err := r.client.request(ctx, "POST", fmt.Sprintf("/api/user/domains/%s/verify", id), nil, &result)
	return &result, err
}

// Delete removes a custom domain.
func (r *DomainsResource) Delete(ctx context.Context, id string) error {
	return r.client.request(ctx, "DELETE", fmt.Sprintf("/api/user/domains/%s", id), nil, nil)
}

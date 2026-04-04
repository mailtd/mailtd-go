package mailtd

import "context"

// BillingResource handles Pro subscription and billing API calls.
type BillingResource struct {
	client *Client
}

// GetStatus returns the current subscription status.
func (r *BillingResource) GetStatus(ctx context.Context) (*SubscriptionStatus, error) {
	var result SubscriptionStatus
	err := r.client.request(ctx, "GET", "/api/user/subscription/status", nil, &result)
	return &result, err
}

// Cancel cancels the current subscription.
func (r *BillingResource) Cancel(ctx context.Context) error {
	return r.client.request(ctx, "POST", "/api/user/subscription/cancel", nil, nil)
}

// Resume resumes a cancelled subscription.
func (r *BillingResource) Resume(ctx context.Context) error {
	return r.client.request(ctx, "POST", "/api/user/subscription/resume", nil, nil)
}

// GetPortalURL returns a Stripe billing portal URL.
func (r *BillingResource) GetPortalURL(ctx context.Context) (*PortalURL, error) {
	var result PortalURL
	err := r.client.request(ctx, "POST", "/api/user/billing/portal", nil, &result)
	return &result, err
}

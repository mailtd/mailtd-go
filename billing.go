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

// CancelWithResult cancels the current subscription and returns the cancel mode and operation ID.
func (r *BillingResource) CancelWithResult(ctx context.Context) (*CancelResult, error) {
	var result CancelResult
	err := r.client.request(ctx, "POST", "/api/user/subscription/cancel", nil, &result)
	return &result, err
}

// Cancel cancels the current subscription.
//
// Deprecated: Use CancelWithResult to receive the cancel mode and operation ID.
func (r *BillingResource) Cancel(ctx context.Context) error {
	_, err := r.CancelWithResult(ctx)
	return err
}

// ResumeWithResult resumes a cancelled subscription and returns the status and operation ID.
func (r *BillingResource) ResumeWithResult(ctx context.Context) (*ResumeResult, error) {
	var result ResumeResult
	err := r.client.request(ctx, "POST", "/api/user/subscription/resume", nil, &result)
	return &result, err
}

// Resume resumes a cancelled subscription.
//
// Deprecated: Use ResumeWithResult to receive the status and operation ID.
func (r *BillingResource) Resume(ctx context.Context) error {
	_, err := r.ResumeWithResult(ctx)
	return err
}

// GetPortalURL returns a Stripe billing portal URL.
func (r *BillingResource) GetPortalURL(ctx context.Context) (*PortalURL, error) {
	var result PortalURL
	err := r.client.request(ctx, "POST", "/api/user/billing/portal", nil, &result)
	return &result, err
}

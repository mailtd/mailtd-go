package mailtd

import (
	"context"
	"fmt"
)

// WebhooksResource handles Pro webhook API calls.
type WebhooksResource struct {
	client *Client
}

// List returns all webhooks.
func (r *WebhooksResource) List(ctx context.Context) ([]Webhook, error) {
	var wrapper struct {
		Webhooks []Webhook `json:"webhooks"`
	}
	err := r.client.request(ctx, "GET", "/api/user/webhooks", nil, &wrapper)
	return wrapper.Webhooks, err
}

// Create adds a new webhook.
func (r *WebhooksResource) Create(ctx context.Context, url string, events []string) (*Webhook, error) {
	body := map[string]any{"url": url, "events": events}
	var result Webhook
	err := r.client.request(ctx, "POST", "/api/user/webhooks", body, &result)
	return &result, err
}

// Delete removes a webhook.
func (r *WebhooksResource) Delete(ctx context.Context, id string) error {
	return r.client.request(ctx, "DELETE", fmt.Sprintf("/api/user/webhooks/%s", id), nil, nil)
}

// RotateSecret generates a new secret for a webhook.
func (r *WebhooksResource) RotateSecret(ctx context.Context, id string) (*Webhook, error) {
	var result Webhook
	err := r.client.request(ctx, "POST", fmt.Sprintf("/api/user/webhooks/%s/rotate", id), nil, &result)
	return &result, err
}

// ListDeliveries returns delivery attempts for a webhook.
func (r *WebhooksResource) ListDeliveries(ctx context.Context, id string) ([]WebhookDelivery, error) {
	var wrapper struct {
		Deliveries []WebhookDelivery `json:"deliveries"`
	}
	err := r.client.request(ctx, "GET", fmt.Sprintf("/api/user/webhooks/%s/deliveries", id), nil, &wrapper)
	return wrapper.Deliveries, err
}

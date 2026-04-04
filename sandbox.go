package mailtd

import (
	"context"
	"fmt"
)

// SandboxResource handles Pro sandbox API calls.
type SandboxResource struct {
	client *Client
}

// GetInfo returns sandbox environment details.
func (r *SandboxResource) GetInfo(ctx context.Context) (*SandboxInfo, error) {
	var result SandboxInfo
	err := r.client.request(ctx, "GET", "/api/user/sandbox", nil, &result)
	return &result, err
}

// SandboxListOptions are optional parameters for listing sandbox messages.
type SandboxListOptions struct {
	Page int
}

// ListMessages returns sandbox messages.
func (r *SandboxResource) ListMessages(ctx context.Context, opts *SandboxListOptions) (*PaginatedResponse[SandboxEmailSummary], error) {
	path := "/api/user/sandbox/messages"
	if opts != nil {
		path = addPageParam(path, opts.Page)
	}
	var result PaginatedResponse[SandboxEmailSummary]
	err := r.client.request(ctx, "GET", path, nil, &result)
	return &result, err
}

// GetMessage returns a single sandbox message.
func (r *SandboxResource) GetMessage(ctx context.Context, id string) (*EmailDetail, error) {
	var result EmailDetail
	err := r.client.request(ctx, "GET", fmt.Sprintf("/api/user/sandbox/messages/%s", id), nil, &result)
	return &result, err
}

// DeleteMessage removes a sandbox message.
func (r *SandboxResource) DeleteMessage(ctx context.Context, id string) error {
	return r.client.request(ctx, "DELETE", fmt.Sprintf("/api/user/sandbox/messages/%s", id), nil, nil)
}

// PurgeMessages removes all sandbox messages.
func (r *SandboxResource) PurgeMessages(ctx context.Context) error {
	return r.client.request(ctx, "DELETE", "/api/user/sandbox/messages", nil, nil)
}

// GetMessageSource returns the raw MIME source of a sandbox message.
func (r *SandboxResource) GetMessageSource(ctx context.Context, id string) ([]byte, error) {
	return r.client.requestRaw(ctx, "GET", fmt.Sprintf("/api/user/sandbox/messages/%s/source", id))
}

// GetAttachment returns the raw bytes of a sandbox message attachment.
func (r *SandboxResource) GetAttachment(ctx context.Context, id string, index int) ([]byte, error) {
	return r.client.requestRaw(ctx, "GET", fmt.Sprintf("/api/user/sandbox/messages/%s/attachments/%d", id, index))
}

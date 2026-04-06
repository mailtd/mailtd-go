package mailtd

import (
	"context"
	"fmt"
)

// MessagesResource handles message-related API calls.
type MessagesResource struct {
	client *Client
}

// ListOptions are optional parameters for listing messages.
type ListOptions struct {
	Page int
}

// List returns messages for an account. accountID accepts a UUID or email address.
func (r *MessagesResource) List(ctx context.Context, accountID string, opts *ListOptions) (*PaginatedResponse[EmailSummary], error) {
	path := fmt.Sprintf("/api/accounts/%s/messages", accountID)
	if opts != nil {
		path = addPageParam(path, opts.Page)
	}
	var result PaginatedResponse[EmailSummary]
	err := r.client.request(ctx, "GET", path, nil, &result)
	return &result, err
}

// Get returns a single message. accountID accepts a UUID or email address.
func (r *MessagesResource) Get(ctx context.Context, accountID, messageID string) (*EmailDetail, error) {
	var result EmailDetail
	err := r.client.request(ctx, "GET", fmt.Sprintf("/api/accounts/%s/messages/%s", accountID, messageID), nil, &result)
	return &result, err
}

// Delete removes a message. accountID accepts a UUID or email address.
func (r *MessagesResource) Delete(ctx context.Context, accountID, messageID string) error {
	return r.client.request(ctx, "DELETE", fmt.Sprintf("/api/accounts/%s/messages/%s", accountID, messageID), nil, nil)
}

// GetSource returns the raw MIME source of a message. accountID accepts a UUID or email address.
func (r *MessagesResource) GetSource(ctx context.Context, accountID, messageID string) ([]byte, error) {
	return r.client.requestRaw(ctx, "GET", fmt.Sprintf("/api/accounts/%s/messages/%s/source", accountID, messageID))
}

// MarkAsRead marks a single message as read. accountID accepts a UUID or email address.
func (r *MessagesResource) MarkAsRead(ctx context.Context, accountID, messageID string) error {
	return r.client.request(ctx, "PUT", fmt.Sprintf("/api/accounts/%s/messages/%s/read", accountID, messageID), nil, nil)
}

// BatchMarkAsReadOptions are optional parameters for batch marking messages as read.
type BatchMarkAsReadOptions struct {
	MessageIDs []string `json:"message_ids,omitempty"`
}

// BatchMarkAsRead marks multiple messages as read. accountID accepts a UUID or email address.
func (r *MessagesResource) BatchMarkAsRead(ctx context.Context, accountID string, opts *BatchMarkAsReadOptions) error {
	var body any
	if opts != nil {
		body = opts
	}
	return r.client.request(ctx, "PUT", fmt.Sprintf("/api/accounts/%s/messages/read", accountID), body, nil)
}

// GetAttachment returns the raw bytes of an attachment. accountID accepts a UUID or email address.
func (r *MessagesResource) GetAttachment(ctx context.Context, accountID, messageID string, index int) ([]byte, error) {
	return r.client.requestRaw(ctx, "GET", fmt.Sprintf("/api/accounts/%s/messages/%s/attachments/%d", accountID, messageID, index))
}

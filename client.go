// Package mailtd provides a Go client for the Mail.td API.
//
// Usage:
//
//	client := mailtd.NewClient("your-api-token")
//	domains, err := client.Accounts.ListDomains(context.Background())
package mailtd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

const defaultBaseURL = "https://api.mail.td"

// Client is the Mail.td API client.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client

	Accounts *AccountsResource
	Messages *MessagesResource
	Domains  *DomainsResource
	Webhooks *WebhooksResource
	Tokens   *TokensResource
	Sandbox  *SandboxResource
	Billing  *BillingResource
	User     *UserResource
}

// Option configures a Client.
type Option func(*Client)

// WithBaseURL sets a custom base URL.
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// NewClient creates a new Mail.td API client.
func NewClient(token string, opts ...Option) *Client {
	c := &Client{
		baseURL:    defaultBaseURL,
		token:      token,
		httpClient: http.DefaultClient,
	}
	for _, opt := range opts {
		opt(c)
	}
	c.Accounts = &AccountsResource{client: c}
	c.Messages = &MessagesResource{client: c}
	c.Domains = &DomainsResource{client: c}
	c.Webhooks = &WebhooksResource{client: c}
	c.Tokens = &TokensResource{client: c}
	c.Sandbox = &SandboxResource{client: c}
	c.Billing = &BillingResource{client: c}
	c.User = &UserResource{client: c}
	return c
}

// APIError represents an error response from the API.
type APIError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("mailtd: %d %s: %s", e.Status, e.Code, e.Message)
}

func (c *Client) request(ctx context.Context, method, path string, body any, result any) error {
	u := c.baseURL + path

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("mailtd: marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, u, bodyReader)
	if err != nil {
		return fmt.Errorf("mailtd: create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("mailtd: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		apiErr := &APIError{Status: resp.StatusCode}
		var errBody struct {
			Error   string `json:"error"`
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errBody); err != nil {
			return &APIError{Status: resp.StatusCode, Code: "unknown", Message: resp.Status}
		}
		apiErr.Code = errBody.Error
		if apiErr.Code == "" {
			apiErr.Code = errBody.Code
		}
		apiErr.Message = errBody.Message
		if apiErr.Message == "" {
			apiErr.Message = errBody.Error
		}
		return apiErr
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("mailtd: decode response: %w", err)
		}
	}

	return nil
}

func (c *Client) requestRaw(ctx context.Context, method, path string) ([]byte, error) {
	u := c.baseURL + path

	req, err := http.NewRequestWithContext(ctx, method, u, nil)
	if err != nil {
		return nil, fmt.Errorf("mailtd: create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("mailtd: send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		apiErr := &APIError{Status: resp.StatusCode}
		var errBody struct {
			Error   string `json:"error"`
			Code    string `json:"code"`
			Message string `json:"message"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errBody); err != nil {
			return nil, &APIError{Status: resp.StatusCode, Code: "unknown", Message: resp.Status}
		}
		apiErr.Code = errBody.Error
		if apiErr.Code == "" {
			apiErr.Code = errBody.Code
		}
		apiErr.Message = errBody.Message
		if apiErr.Message == "" {
			apiErr.Message = errBody.Error
		}
		return nil, apiErr
	}

	return io.ReadAll(resp.Body)
}

func addPageParam(path string, page int) string {
	if page <= 0 {
		return path
	}
	u, err := url.Parse(path)
	if err != nil {
		return path
	}
	q := u.Query()
	q.Set("page", strconv.Itoa(page))
	u.RawQuery = q.Encode()
	return u.String()
}

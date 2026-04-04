package mailtd

import "time"

// Domain represents a public email domain.
type Domain struct {
	ID        string `json:"id"`
	Domain    string `json:"domain"`
	Default   bool   `json:"default"`
	SortOrder int    `json:"sort_order"`
}

// AccountInfo represents an email account.
type AccountInfo struct {
	ID        string    `json:"id"`
	Address   string    `json:"address"`
	Role      string    `json:"role"`
	Quota     int64     `json:"quota"`
	Used      int64     `json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

// CreateAccountResult is the response from creating an account.
type CreateAccountResult struct {
	ID                      string `json:"id"`
	Address                 string `json:"address"`
	Token                   string `json:"token"`
	SuggestedNextDifficulty int    `json:"suggested_next_difficulty,omitempty"`
}

// LoginResult is the response from logging in.
type LoginResult struct {
	ID      string `json:"id"`
	Address string `json:"address"`
	Token   string `json:"token"`
}

// EmailSummary represents a message in a list.
type EmailSummary struct {
	ID          string    `json:"id"`
	Sender      string    `json:"sender"`
	From        string    `json:"from"`
	Subject     string    `json:"subject"`
	PreviewText string    `json:"preview_text"`
	Size        int64     `json:"size"`
	IsRead      bool      `json:"is_read"`
	CreatedAt   time.Time `json:"created_at"`
}

// EmailDetail represents a full email message.
type EmailDetail struct {
	ID          string       `json:"id"`
	Sender      string       `json:"sender"`
	From        string       `json:"from"`
	Subject     string       `json:"subject"`
	Address     string       `json:"address"`
	Size        int64        `json:"size"`
	CreatedAt   time.Time    `json:"created_at"`
	TextBody    *string      `json:"text_body"`
	HTMLBody    *string      `json:"html_body"`
	Attachments []Attachment `json:"attachments"`
}

// Attachment represents an email attachment.
type Attachment struct {
	Index       int    `json:"index"`
	Filename    string `json:"filename"`
	ContentType string `json:"content_type"`
	Size        int64  `json:"size"`
}

// ProUser represents a Pro user profile.
type ProUser struct {
	ID           string     `json:"id"`
	Email        string     `json:"email"`
	Plan         string     `json:"plan"`
	Role         string     `json:"role"`
	Status       string     `json:"status"`
	MaxAccounts  int        `json:"max_accounts"`
	MaxDomains   int        `json:"max_domains"`
	AccountCount int        `json:"account_count"`
	DomainCount  int        `json:"domain_count"`
	CreatedAt    time.Time  `json:"created_at"`
	Downgraded   *time.Time `json:"downgraded"`
}

// ProDomain represents a custom domain.
type ProDomain struct {
	ID           string     `json:"id"`
	Domain       string     `json:"domain"`
	VerifyStatus string     `json:"verify_status"`
	VerifyToken  string     `json:"verify_token"`
	VerifiedAt   *time.Time `json:"verified_at"`
	MXConfigured bool       `json:"mx_configured"`
	CreatedAt    time.Time  `json:"created_at"`
}

// DNSRecord represents a DNS record for domain verification.
type DNSRecord struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Value    string `json:"value"`
	Priority *int   `json:"priority"`
	OK       bool   `json:"ok"`
}

// CreateDomainResult is the response from creating a domain.
type CreateDomainResult struct {
	ID          string      `json:"id"`
	Domain      string      `json:"domain"`
	VerifyToken string      `json:"verify_token"`
	DNSRecords  []DNSRecord `json:"dns_records"`
}

// VerifyDomainResult is the response from verifying a domain.
type VerifyDomainResult struct {
	VerifyStatus string      `json:"verify_status"`
	TXTRecord    *DNSRecord  `json:"txt_record"`
	MXRecord     *DNSRecord  `json:"mx_record"`
	DNSRecords   []DNSRecord `json:"dns_records"`
	Message      *string     `json:"message"`
}

// Webhook represents a webhook configuration.
type Webhook struct {
	ID              string     `json:"id"`
	URL             string     `json:"url"`
	Events          []string   `json:"events"`
	Secret          string     `json:"secret"`
	Status          string     `json:"status"`
	FailureCount    int        `json:"failure_count"`
	LastTriggeredAt *time.Time `json:"last_triggered_at"`
	CreatedAt       time.Time  `json:"created_at"`
}

// WebhookDelivery represents a webhook delivery attempt.
type WebhookDelivery struct {
	ID         string    `json:"id"`
	EventType  string    `json:"event_type"`
	EventID    string    `json:"event_id"`
	StatusCode *int      `json:"status_code"`
	Error      *string   `json:"error"`
	Attempt    int       `json:"attempt"`
	DurationMs int       `json:"duration_ms"`
	CreatedAt  time.Time `json:"created_at"`
}

// Token represents an API token.
type Token struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	TokenValue *string    `json:"token_value"`
	LastUsedAt *time.Time `json:"last_used_at"`
	CreatedAt  time.Time  `json:"created_at"`
	RevokedAt  *time.Time `json:"revoked_at"`
}

// SandboxInfo represents sandbox environment details.
type SandboxInfo struct {
	Enabled    bool   `json:"enabled"`
	AccountID  string `json:"account_id"`
	Address    string `json:"address"`
	SMTPHost   string `json:"smtp_host"`
	SMTPPort   int    `json:"smtp_port"`
	AuthMethod string `json:"auth_method"`
	Username   string `json:"username"`
	Note       string `json:"note"`
	Quota      int64  `json:"quota"`
	Used       int64  `json:"used"`
}

// SandboxEmailSummary represents a sandbox message in a list.
type SandboxEmailSummary struct {
	ID          string    `json:"id"`
	Sender      string    `json:"sender"`
	From        string    `json:"from"`
	Subject     string    `json:"subject"`
	PreviewText string    `json:"preview_text"`
	Size        int64     `json:"size"`
	CreatedAt   time.Time `json:"created_at"`
}

// SubscriptionStatus represents the current subscription state.
type SubscriptionStatus struct {
	Status            string     `json:"status"`
	CancelMode        *string    `json:"cancel_mode"`
	ScheduledCancelAt *time.Time `json:"scheduled_cancel_at"`
}

// PaginatedResponse wraps a paginated list of items.
type PaginatedResponse[T any] struct {
	Data       []T `json:"data"`
	Page       int `json:"page"`
	TotalPages int `json:"total_pages"`
	Total      int `json:"total"`
}

// PortalURL is the response from GetPortalURL.
type PortalURL struct {
	URL string `json:"url"`
}

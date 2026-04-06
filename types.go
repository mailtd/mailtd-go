package mailtd

import (
	"encoding/json"
	"time"
)

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
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	Plan        string    `json:"plan"`
	Role        string    `json:"role"`
	Status      string    `json:"status"`
	MaxDomains  int       `json:"max_domains"`
	OpsUsed     int       `json:"ops_used"`
	OpsLimit    int       `json:"ops_limit"`
	DomainCount int       `json:"domain_count"`
	CreatedAt   time.Time `json:"created_at"`
	Downgraded  bool      `json:"downgraded"`
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
	TXTRecordOK  bool        `json:"-"`
	MXRecordOK   bool        `json:"-"`
	DNSRecords   []DNSRecord `json:"dns_records"`
	Message      *string     `json:"message"`

	// Deprecated: TXTRecord was never successfully decoded; use TXTRecordOK instead.
	TXTRecord *DNSRecord `json:"-"`
	// Deprecated: MXRecord was never successfully decoded; use MXRecordOK instead.
	MXRecord *DNSRecord `json:"-"`
}

func (v *VerifyDomainResult) UnmarshalJSON(data []byte) error {
	type alias VerifyDomainResult
	var raw struct {
		alias
		TXTRecord json.RawMessage `json:"txt_record"`
		MXRecord  json.RawMessage `json:"mx_record"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*v = VerifyDomainResult(raw.alias)
	if len(raw.TXTRecord) > 0 {
		_ = json.Unmarshal(raw.TXTRecord, &v.TXTRecordOK)
	}
	if len(raw.MXRecord) > 0 {
		_ = json.Unmarshal(raw.MXRecord, &v.MXRecordOK)
	}
	return nil
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
	Token      string     `json:"token,omitempty"`
	LastUsedAt *time.Time `json:"last_used_at"`
	CreatedAt  time.Time  `json:"created_at"`
	RevokedAt  *time.Time `json:"revoked_at"`

	// Deprecated: TokenValue is unused; the backend returns the field as "token".
	// Use Token instead.
	TokenValue *string `json:"token_value,omitempty"`
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

// ScheduledChange represents a scheduled subscription change.
type ScheduledChange struct {
	Action      string `json:"action"`
	EffectiveAt string `json:"effective_at"`
}

// SubscriptionStatus represents the current subscription state.
type SubscriptionStatus struct {
	Status          string           `json:"status"`
	CancelMode      *string          `json:"cancel_mode"`
	ScheduledCancel *ScheduledChange `json:"-"`

	// Deprecated: ScheduledCancelAt only captures the timestamp.
	// Use ScheduledCancel for the full object (action + effective_at).
	ScheduledCancelAt *time.Time `json:"-"`
}

func (s *SubscriptionStatus) UnmarshalJSON(data []byte) error {
	type alias SubscriptionStatus
	var raw struct {
		alias
		ScheduledCancelAt json.RawMessage `json:"scheduled_cancel_at"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	*s = SubscriptionStatus(raw.alias)
	if len(raw.ScheduledCancelAt) > 0 && string(raw.ScheduledCancelAt) != "null" {
		var sc ScheduledChange
		if err := json.Unmarshal(raw.ScheduledCancelAt, &sc); err == nil {
			s.ScheduledCancel = &sc
			if t, err := time.Parse(time.RFC3339, sc.EffectiveAt); err == nil {
				s.ScheduledCancelAt = &t
			}
		}
	}
	return nil
}

// MessageListResult is the response from listing messages.
type MessageListResult struct {
	Messages []EmailSummary `json:"messages"`
	Page     int            `json:"page"`
}

// SandboxMessageListResult is the response from listing sandbox messages.
type SandboxMessageListResult struct {
	Messages []SandboxEmailSummary `json:"messages"`
	Page     int                   `json:"page"`
}

// AccountListResult is the response from listing user accounts.
type AccountListResult struct {
	Accounts   []AccountInfo `json:"accounts"`
	NextCursor string        `json:"next_cursor"`
}

// PortalURL is the response from GetPortalURL.
type PortalURL struct {
	URL string `json:"url"`
}

// BatchMarkAsReadResult is the response from batch marking messages as read.
type BatchMarkAsReadResult struct {
	Updated int `json:"updated"`
}

// PurgeMessagesResult is the response from purging all sandbox messages.
type PurgeMessagesResult struct {
	Deleted int `json:"deleted"`
}

// CancelResult is the response from cancelling a subscription.
type CancelResult struct {
	CancelMode  string `json:"cancel_mode"`
	OperationID string `json:"operation_id,omitempty"`
}

// ResumeResult is the response from resuming a subscription.
type ResumeResult struct {
	Status      string `json:"status"`
	OperationID string `json:"operation_id,omitempty"`
}

# mailtd-go

[![Go Reference](https://pkg.go.dev/badge/github.com/mailtd/mailtd-go.svg)](https://pkg.go.dev/github.com/mailtd/mailtd-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Official Go SDK for [Mail.td](https://mail.td) — the developer email platform for **temp mail**, **email testing**, and **SMTP sandbox**.

- **Temp Mail API** — Create and manage temporary email addresses programmatically
- **Email Testing** — Receive, inspect, and verify emails in your test suite
- **SMTP Sandbox** — Capture outbound emails in a safe sandbox environment without sending to real inboxes
- **Webhooks** — Get notified in real-time when emails arrive
- **Custom Domains** — Use your own domain for branded temporary mailboxes

## Installation

```bash
go get github.com/mailtd/mailtd-go
```

## Quick Start

```go
package main

import (
	"context"
	"fmt"
	"log"

	mailtd "github.com/mailtd/mailtd-go"
)

func main() {
	client := mailtd.NewClient("td_...")
	ctx := context.Background()

	// Create a temporary email address
	pw := "mypassword"
	account, err := client.Accounts.Create(ctx, "test@mail.td", &mailtd.CreateOptions{
		Password: &pw,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created account: %s\n", account.Address)

	// List messages
	result, err := client.Messages.List(ctx, account.ID, nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, m := range result.Data {
		fmt.Printf("%s: %s\n", m.From, m.Subject)
	}
}
```

## Use Cases

- **Automated testing** — Create temp mail addresses in CI/CD to test signup flows, OTP verification, and transactional emails
- **Email verification testing** — Validate that your app sends the right emails with the right content
- **SMTP sandbox** — Route your app's outbound SMTP to Mail.td sandbox to inspect emails without spamming real users
- **QA environments** — Give each test run its own mailbox, then tear it down

## Authentication

All API calls require a Pro API Token (`td_...`). Pass it when creating the client:

```go
client := mailtd.NewClient("td_...")
```

### Options

```go
// Custom base URL
client := mailtd.NewClient("td_...", mailtd.WithBaseURL("https://custom.api.url"))

// Custom HTTP client
client := mailtd.NewClient("td_...", mailtd.WithHTTPClient(&http.Client{
    Timeout: 30 * time.Second,
}))
```

## Resources

| Resource | Description |
|----------|-------------|
| `client.Accounts` | Create, get, delete accounts; reset password; list domains |
| `client.Messages` | List, get, delete messages; attachments; mark as read |
| `client.Domains` | Pro: manage custom domains |
| `client.Webhooks` | Pro: manage webhooks |
| `client.Tokens` | Pro: manage API tokens |
| `client.Sandbox` | Pro: sandbox email testing |
| `client.Billing` | Pro: subscription management |
| `client.User` | Pro: user profile and account management |

## Error Handling

API errors are returned as `*mailtd.APIError`:

```go
import "errors"

_, err := client.Accounts.Get(ctx, "invalid-id")
if err != nil {
    var apiErr *mailtd.APIError
    if errors.As(err, &apiErr) {
        fmt.Printf("Status: %d, Code: %s, Message: %s\n", apiErr.Status, apiErr.Code, apiErr.Message)
    }
}
```

## Links

- [Website](https://mail.td) — Create temp mail, email testing, SMTP sandbox
- [API Documentation](https://docs.mail.td) — Full API reference
- [Node.js SDK](https://www.npmjs.com/package/mailtd) — `npm install mailtd`
- [Python SDK](https://pypi.org/project/mailtd/) — `pip install mailtd`
- [CLI](https://github.com/mailtd/mailcx-cli) — Command-line tool

## License

MIT

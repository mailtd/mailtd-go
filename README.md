# mailtd-go

Official Go SDK for the [Mail.td](https://mail.td) API.

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
	client := mailtd.NewClient("your-api-token")
	ctx := context.Background()

	// List available domains
	domains, err := client.Accounts.ListDomains(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, d := range domains {
		fmt.Println(d.Domain)
	}

	// Create an account on a system domain
	result, err := client.Accounts.Create(ctx, "demo@sugtbt.com", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Created account: %s (token: %s)\n", result.Address, result.Token)

	// List messages
	messages, err := client.Messages.List(ctx, result.ID, nil)
	if err != nil {
		log.Fatal(err)
	}
	for _, m := range messages.Messages {
		fmt.Printf("%s: %s\n", m.From, m.Subject)
	}
}
```

## Authentication

All API calls require a bearer token. Pass it when creating the client:

```go
client := mailtd.NewClient("your-token")
```

### Options

```go
// Custom base URL
client := mailtd.NewClient("token", mailtd.WithBaseURL("https://custom.api.url"))

// Custom HTTP client
client := mailtd.NewClient("token", mailtd.WithHTTPClient(&http.Client{
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

## License

MIT

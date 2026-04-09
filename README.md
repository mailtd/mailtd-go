# mailtd-go

Official Go SDK for the [Mail.td](https://mail.td) developer email platform.

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

	// Create a mailbox
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

## License

MIT

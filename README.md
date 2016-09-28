# go-imap-idle

[![GoDoc](https://godoc.org/github.com/emersion/go-imap-idle?status.svg)](https://godoc.org/github.com/emersion/go-imap-idle)

[IDLE extension](https://tools.ietf.org/html/rfc2177) for [go-imap](https://github.com/emersion/go-imap).

## Usage

### Client

```go
// Select a mailbox
if _, err := c.Select("INBOX", false); err != nil {
	log.Fatal(err)
}

idleClient := idle.NewClient(c)

// Get capabilities if needed
if c.Caps == nil {
	if _, err := c.Capability(); err != nil {
		log.Fatal(err)
	}
}

// Create a channel to receive mailbox updates
statuses := make(chan *imap.MailboxStatus)
c.MailboxUpdates = statuses

// Check support for the IDLE extension
if idleClient.SupportsIdle() {
	// Start idling
	stop := make(chan struct{})
	done := make(chan error, 1)
	go func() {
		done <- idleClient.Idle(stop)
	}()

	// Listen for updates
	for {
		select {
		case status := <-statuses:
			log.Println("New mailbox status:", status)
			close(stop)
		case err := <-done:
			if err != nil {
				log.Fatal(err)
			}
			log.Println("Not idling anymore")
			return
		}
	}
} else {
	// Fallback: call periodically c.Noop()
}
```

### Server

```go
s.Enable(idle.NewExtension())
```

## License

MIT

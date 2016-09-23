# go-imap-idle

[![GoDoc](https://godoc.org/github.com/emersion/go-imap-idle?status.svg)](https://godoc.org/github.com/emersion/go-imap-idle)

[IDLE extension](https://tools.ietf.org/html/rfc2177) for [go-imap](https://github.com/emersion/go-imap).

## Usage

### Client

```go
idleClient := idle.NewClient(c)

// Get capabilities if needed
if c.Caps == nil {
	if _, err := c.Capability(); err != nil {
		log.Fatal(err)
	}
}

if idleClient.SupportsIdle() {
	stop := make(chan struct{})
	done := make(chan error, 1)
	go func() {
		done <- idleClient.Idle(stop)
	}()

	for {
		select {
		case status := <-c.MailboxUpdates:
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

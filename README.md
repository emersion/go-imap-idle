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

// Create a channel to receive mailbox updates
statuses := make(chan *imap.MailboxStatus)
c.MailboxUpdates = statuses

// Start idling
stopped := false
stop := make(chan struct{})
done := make(chan error, 1)
go func() {
	done <- idleClient.IdleWithFallback(stop, 0)
}()

// Listen for updates
for {
	select {
	case status := <-statuses:
		log.Println("New mailbox status:", status)
		if !stopped {
			close(stop)
			stopped = true
		}
	case err := <-done:
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Not idling anymore")
		return
	}
}
```

### Server

```go
s.Enable(idle.NewExtension())
```

## License

MIT

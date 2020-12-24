# go-imap-idle

[![GoDoc](https://godoc.org/github.com/emersion/go-imap-idle?status.svg)](https://godoc.org/github.com/emersion/go-imap-idle)

[IDLE extension](https://tools.ietf.org/html/rfc2177) for [go-imap](https://github.com/emersion/go-imap).

## Usage

### Client

```go
// Let's assume c is an IMAP client
var c *client.Client

// Select a mailbox
if _, err := c.Select("INBOX", false); err != nil {
	log.Fatal(err)
}

idleClient := idle.NewClient(c)

// Create a channel to receive mailbox updates
updates := make(chan client.Update)
c.Updates = updates

// Start idling
done := make(chan error, 1)
go func() {
	done <- idleClient.IdleWithFallback(nil, 0)
}()

// Listen for updates
for {
	select {
	case update := <-updates:
		log.Println("New update:", update)
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

Backend should implement Backend interface defined in this package.

## License

MIT

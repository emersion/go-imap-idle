package idle_test

import (
	"log"

	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-imap-idle"
)

func ExampleClient_Idle() {
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

	// Check support for the IDLE extension
	if ok, err := idleClient.SupportIdle(); err == nil && ok {
		// Start idling
		stopped := false
		stop := make(chan struct{})
		done := make(chan error, 1)
		go func() {
			done <- idleClient.Idle(stop)
		}()

		// Listen for updates
		for {
			select {
			case update := <-updates:
				log.Println("New update:", update)
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
	} else {
		// Fallback: call periodically c.Noop()
	}
}

func ExampleClient_IdleWithFallback() {
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
}

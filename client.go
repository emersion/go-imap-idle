package idle

import (
	"github.com/emersion/go-imap/client"
)

type Client struct {
	client *client.Client
}

// Indicate to the server that the client is ready to receive unsolicited
// mailbox update messages. When the client wants to send commands again, it
// must first close done.
func (c *Client) Idle(done <-chan struct{}) error {
	cmd := &Command{}

	res := &Response{
		Done: done,
		Writer: c.client.Writer(),
	}

	status, err := c.client.Execute(cmd, res)
	if err != nil {
		return err
	}
	return status.Err()
}

// SupportsIdle returns true if the server supports the IDLE extension.
func (c *Client) SupportsIdle() bool {
	return c.client.Caps[Capability]
}

// Create a new client.
func NewClient(c *client.Client) *Client {
	return &Client{client: c}
}

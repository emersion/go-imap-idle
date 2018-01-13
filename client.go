package idle

import (
	"errors"
	"time"

	"github.com/emersion/go-imap/client"
)

const (
	defaultLogoutTimeout = 25 * time.Minute
	defaultPollInterval = time.Minute
)

// Client is an IDLE client.
type Client struct {
	c *client.Client

	// LogoutTimeout is used to avoid being logged out by the server when idling.
	// Each LogoutTimeout, the idle command is restarted. If set to zero, this
	// behavior is disabled.
	LogoutTimeout time.Duration
}

// NewClient creates a new client.
func NewClient(c *client.Client) *Client {
	return &Client{c, defaultLogoutTimeout}
}

func (c *Client) idle(stop <-chan struct{}) error {
	cmd := &Command{}

	done := make(chan error, 1)
	res := &Response{
		Stop:   stop,
		Done:   done,
		Writer: c.c.Writer(),
	}

	if status, err := c.c.Execute(cmd, res); err != nil {
		return err
	} else if err := status.Err(); err != nil {
		return err
	} else {
		return <-done
	}
}

// Idle indicates to the server that the client is ready to receive unsolicited
// mailbox update messages. When the client wants to send commands again, it
// must first close stop.
func (c *Client) Idle(stop <-chan struct{}) error {
	if c.LogoutTimeout == 0 {
		return c.idle(stop)
	}

	t := time.NewTicker(c.LogoutTimeout)
	defer t.Stop()

	for {
		stopOrRestart := make(chan struct{})
		done := make(chan error, 1)
		go func() {
			done <- c.idle(stopOrRestart)
		}()

		select {
		case <-t.C:
			close(stopOrRestart)
			if err := <-done; err != nil {
				return err
			}
		case <-stop:
			close(stopOrRestart)
			return <-done
		case err := <-done:
			if err != nil {
				return err
			}
		}
	}
}

// SupportIdle checks if the server supports the IDLE extension.
func (c *Client) SupportIdle() (bool, error) {
	return c.c.Support(Capability)
}

// IdleWithFallback tries to idle if the server supports it. If it doesn't, it
// falls back to polling. If pollInterval is zero, a sensible default will be
// used.
func (c *Client) IdleWithFallback(stop <-chan struct{}, pollInterval time.Duration) error {
	if ok, err := c.SupportIdle(); err != nil {
		return err
	} else if ok {
		return c.Idle(stop)
	}

	if pollInterval == 0 {
		pollInterval = defaultPollInterval
	}

	t := time.NewTicker(pollInterval)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			if err := c.c.Noop(); err != nil {
				return err
			}
		case <-stop:
			return nil
		case <-c.c.LoggedOut():
			return errors.New("disconnected while idling")
		}
	}
}

// IdleClient is an alias used to compose multiple client extensions.
type IdleClient = Client

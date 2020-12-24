package idle

import (
	"bufio"
	"errors"
	"strings"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/server"
)

type Handler struct {
	Command
}

type Mailbox interface {
	/*
	Idle allows backend to send updates without explicit Poll calls or any other
	commands running.

	When called - it should block indefinitely and return immediately when
	done channel is written to.
	*/
	Idle(done <-chan struct{})
}

func (h *Handler) Handle(conn server.Conn) error {
	cont := &imap.ContinuationReq{Info: "idling"}
	if err := conn.WriteResp(cont); err != nil {
		return err
	}

	if mbox, ok := conn.Context().Mailbox.(Mailbox); ok {
		done := make(chan struct{})
		go mbox.Idle(done)
		defer func() {
			done <- struct{}{}
		}()
	}
	// TODO(foxcpp): Fallback to short-interval polling if backend does not support
	// corresponding interface?

	// Wait for DONE
	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return err
	}

	if strings.ToUpper(scanner.Text()) != doneLine {
		return errors.New("Expected DONE")
	}
	return nil
}

type extension struct{}

func (ext *extension) Capabilities(c server.Conn) []string {
	if c.Context().State&imap.AuthenticatedState != 0 {
		return []string{Capability}
	}
	return nil
}

func (ext *extension) Command(name string) server.HandlerFactory {
	if name != commandName {
		return nil
	}

	return func() server.Handler {
		return &Handler{}
	}
}

func NewExtension() server.Extension {
	return &extension{}
}

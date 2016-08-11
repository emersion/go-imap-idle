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

func (h *Handler) Handle(conn server.Conn) error {
	cont := &imap.ContinuationResp{Info: "idling"}
	if err := conn.WriteResp(cont); err != nil {
		return err
	}

	// Wait for DONE
	scanner := bufio.NewScanner(conn)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return err
	}

	if strings.ToUpper(scanner.Text()) != "DONE" {
		return errors.New("Expected DONE")
	}
	return nil
}

type extension struct{}

func (ext *extension) Capabilities(state imap.ConnState) (caps []string) {
	if state&imap.SelectedState != 0 {
		caps = append(caps, Capability)
	}
	return
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
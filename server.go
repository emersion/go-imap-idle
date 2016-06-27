package idle

import (
	"errors"
	"strings"

	"github.com/emersion/go-imap/common"
	"github.com/emersion/go-imap/server"
)

type Handler struct {
	Command
}

func (h *Handler) Handle(conn *server.Conn) error {
	cont := &common.ContinuationResp{Info: "idling"}
	if err := conn.WriteResp(cont); err != nil {
		return err
	}

	// Wait for DONE
	line, err := conn.ReadInfo()
	if err != nil {
		return err
	}

	if strings.ToUpper(line) != "DONE" {
		return errors.New("Expected DONE")
	}
	return nil
}

// Enable the IDLE extension for a server.
func NewServer(s *server.Server) {
	s.RegisterCapability(CommandName, common.SelectedState)

	s.RegisterCommand(CommandName, func() server.Handler {
		return &Handler{}
	})
}

package idle

import (
	"github.com/emersion/go-imap/common"
)

// An IDLE command.
// Se RFC 2177 section 3.
type Command struct {}

func (cmd *Command) Command() *common.Command {
	return &common.Command{Name: CommandName}
}

func (cmd *Command) Parse(fields []interface{}) error {
	return nil
}

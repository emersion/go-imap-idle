package idle

import (
	"github.com/emersion/go-imap"
)

// An IDLE response.
type Response struct {
	Done <-chan struct{}
	Writer *imap.Writer
}

func (r *Response) HandleFrom(hdlr imap.RespHandler) error {
	w := r.Writer

	for h := range hdlr {
		if _, ok := h.Resp.(*imap.ContinuationResp); !ok {
			h.Reject()
			continue
		}
		h.Accept()

		<-r.Done

		if _, err := w.Write([]byte(done+"\r\n")); err != nil {
			return err
		}
		if err := w.Flush(); err != nil {
			return err
		}
	}

	return nil
}

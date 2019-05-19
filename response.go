package idle

import (
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/responses"
)

// An IDLE response.
type Response struct {
	RepliesCh chan []byte
	Stop   <-chan struct{}

	gotContinuationReq bool
}

func (r *Response) Replies() <-chan []byte {
	return r.RepliesCh
}

func (r *Response) stop() {
	r.RepliesCh <- []byte(doneLine + "\r\n")
}

func (r *Response) Handle(resp imap.Resp) error {
	// Wait for a continuation request
	if _, ok := resp.(*imap.ContinuationReq); ok && !r.gotContinuationReq {
		r.gotContinuationReq = true

		// We got a continuation request, wait for r.Stop to be closed
		go func() {
			<-r.Stop
			r.stop()
		}()

		return nil
	}

	return responses.ErrUnhandled
}

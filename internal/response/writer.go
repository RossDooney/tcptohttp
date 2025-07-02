package response

import (
	"fmt"
	"io"
)

type Writer struct {
	Writer io.Writer
	state  responseState
}

type responseState int

const (
	responseInitialized responseState = iota
	responseStateWritingStatusLine
	responseStateWritingHeaders
	responseStateWritingBody
	responseStateDone
)

func (w *Writer) Write(p []byte) (n int, err error) {
	return w.Writer.Write(p)
}

func (w *Writer) WriteBody(p []byte) (int, error) {

	if w.state != responseStateWritingHeaders {
		return 0, fmt.Errorf("error: trying to write status line after with responseStateWritingHeaders not set, state set to: %s", w.state)
	}

	w.state = responseStateWritingBody

	w.Write(p)
	return 0, nil
}

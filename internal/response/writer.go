package response

import (
	"fmt"
	"io"
)

type Writer struct {
	Writer io.Writer
	State  responseState
}

type responseState int

const (
	respWritingStatusLine responseState = iota
	respWritingHeaders
	respWritingBody
	responseStateDone
)

func (w *Writer) Write(p []byte) (n int, err error) {
	return w.Writer.Write(p)
}

func (w *Writer) WriteBody(p []byte) (int, error) {

	if w.State != respWritingBody {
		return 0, fmt.Errorf("error: trying to write status line after with respWritingBody not set, state set to: %s", w.State)
	}

	return w.Write(p)
}
func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.State != respWritingBody {
		return 0, fmt.Errorf("error: trying to write status line after with respWritingBody not set, state set to: %s", w.State)
	}
	chunkSize := len(p)

	total := 0

	i, err := fmt.Fprintf(w, "%x\r\n", chunkSize)
	if err != nil {
		return total, err
	}
	total += i

	i, err = w.Write(p)
	if err != nil {
		return total, err
	}
	total += i

	i, err = w.Write([]byte("\r\n"))
	if err != nil {
		return total, err
	}
	total += i
	return total, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.State != respWritingBody {
		return 0, fmt.Errorf("error: trying to write status line after with respWritingBody not set, state set to: %s", w.State)
	}
	i, err := w.Write([]byte("0\r\n\r\n"))
	if err != nil {
		return i, err
	}
	w.State = responseStateDone
	return i, nil
}

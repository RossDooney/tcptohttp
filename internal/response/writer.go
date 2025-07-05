package response

import (
	"fmt"
	"httpTest/internal/headers"
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
	respWritingTrailer
)

func (w *Writer) Write(p []byte) (n int, err error) {
	return w.Writer.Write(p)
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {

	if w.State != respWritingHeaders {
		return fmt.Errorf("error: trying to write status line after with responseStateWritingStatusLine not set, state set to: %s", w.State)
	}

	var headerTxt string

	for key, value := range headers {
		headerTxt = fmt.Sprintf("%s: %s\r\n", key, value)
		_, err := w.Write([]byte(headerTxt))
		if err != nil {
			return err
		}
	}

	_, err := w.Write([]byte("\r\n"))
	if err != nil {
		return err
	}

	w.State = respWritingBody
	return nil
}

func (w *Writer) WriteStatusLine(statusCode ServerStatusCode) error {
	status, ok := statusText[statusCode]

	if w.State != respWritingStatusLine {
		return fmt.Errorf("error: trying to write status line after with responseInitialized not set, state set to: %s", w.State)
	}

	if !ok {
		status = ""
	}

	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, status)
	_, err := w.Write([]byte(statusLine))

	if err != nil {
		return err
	}
	w.State = respWritingHeaders
	return nil
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
	i, err := w.Write([]byte("0\r\n"))
	if err != nil {
		return i, err
	}
	w.State = respWritingTrailer
	return i, nil
}

func (w *Writer) WriteTrailers(h headers.Headers) error {

	if w.State != respWritingTrailer {
		return fmt.Errorf("error: trying to write status line after with respWritingTrailer not set, state set to: %s", w.State)
	}

	var trailerTxt string

	for key, value := range h {
		trailerTxt = fmt.Sprintf("%s: %s\r\n", key, value)
		_, err := w.Write([]byte(trailerTxt))
		if err != nil {
			return err
		}
	}

	_, err := w.Write([]byte("\r\n"))
	if err != nil {
		return err
	}

	return nil
}

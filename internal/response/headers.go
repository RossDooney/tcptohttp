package response

import (
	"fmt"
	"httpTest/internal/headers"
	"strconv"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()

	h.Set("Content-Length", strconv.Itoa(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return h
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {

	if w.state != responseStateWritingStatusLine {
		return fmt.Errorf("error: trying to write status line after with responseStateWritingStatusLine not set, state set to: %s", w.state)
	}

	w.state = responseStateWritingHeaders

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

	return nil
}

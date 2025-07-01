package response

import (
	"fmt"
	"httpTest/internal/headers"
	"io"
	"strconv"
)

type ServerStatusCode int

const (
	StatusOK          ServerStatusCode = 200
	StatusBadRequest  ServerStatusCode = 400
	StatusServerError ServerStatusCode = 500
)

var statusText = map[ServerStatusCode]string{
	StatusOK:          "Ok",
	StatusBadRequest:  "Bad Request",
	StatusServerError: "Internal Server Error",
}

type Writer struct {
	Writer io.Writer
}

func (w *Writer) Write(p []byte) (n int, err error) {
	return w.Writer.Write(p)
}

func (w *Writer) WriteStatusLine(statusCode ServerStatusCode) error {
	status, ok := statusText[statusCode]

	if !ok {
		status = ""
	}

	statusLine := fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, status)
	_, err := w.Write([]byte(statusLine))

	if err != nil {
		return err
	}

	return nil
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()

	h.Set("Content-Length", strconv.Itoa(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")

	return h
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {

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

func (w *Writer) WriteBody(p []byte) (int, error) {
	w.Write(p)
	return 0, nil
}

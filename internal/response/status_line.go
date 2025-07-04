package response

import "fmt"

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

func (rs responseState) String() string {
	switch rs {
	case respWritingStatusLine:
		return "respWritingStatusLine"
	case respWritingHeaders:
		return "respWritingHeaders"
	case respWritingBody:
		return "respWritingBody"
	case responseStateDone:
		return "responseStateDone"
	default:
		return "unknown responseState"
	}
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

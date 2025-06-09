package headers

import (
	"bytes"
	"fmt"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {

	done = false
	bytesRead := 0
	for !done {
		n, err = parseRequestLine(data)
		if err != nil {
			return 0, false, err
		}
		if n == 0 {
			done = true
			break
		}
		requestLineText := string(data[:n])
		key, value, err := headerFromString(requestLineText)

		if err != nil {
			return 0, false, err
		}

		h[key] = value
		// n + 2 to include \r\n, another +2 at the end to include the last escape char
		bytesRead += n + 2
		data = data[n+2:]

	}

	return bytesRead, done, nil

}

func parseRequestLine(data []byte) (int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, nil
	}
	return idx, nil
}

func headerFromString(str string) (string, string, error) {
	before, after, _ := strings.Cut(str, ":")

	start := 0
	for _, c := range before {
		if c != ' ' && c != '\t' {
			break
		}
		start++
	}

	key := before[start:]
	for _, c := range key {
		if c == ' ' || c == '\t' {
			return "", "", fmt.Errorf("Whitespace present in header key")
		}
	}

	start = 0
	end := len(after) - 1

	isASCIISpace := func(c byte) bool {
		return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '\v' || c == '\f'
	}

	for start <= end && isASCIISpace(after[start]) {
		start++
	}
	for end >= start && isASCIISpace(after[end]) {
		end--
	}

	value := after[start : end+1]

	for _, c := range value {
		if c == ' ' || c == '\t' {
			return "", "", fmt.Errorf("Whitespace present in header value")
		}
	}

	return key, value, nil
}

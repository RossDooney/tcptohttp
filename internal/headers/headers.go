package headers

import (
	"bytes"
	"fmt"
	"strings"
)

const crlf = "\r\n"

type Headers map[string]string

func NewHeaders() Headers {
	return map[string]string{}
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, false, nil
	}
	if idx == 0 {
		// headers are done, consume the CRLF
		return 2, true, nil
	}
	parts := bytes.SplitN(data[:idx], []byte(":"), 2)
	key := string(parts[0])

	if key != strings.TrimRight(key, " ") {
		return 0, false, fmt.Errorf("invalid header name: %s", key)
	}

	value := bytes.TrimSpace(parts[1])
	key = strings.TrimSpace(key)

	key, err = normalizeHead(key)

	if err != nil {
		return 0, false, fmt.Errorf("invalid characters in header: %s", key)
	}
	if val, exists := h[key]; exists {
		var valueString strings.Builder

		valueString.WriteString(val)
		valueString.WriteString(", ")
		valueString.WriteString(string(value))
		h.Set(key, valueString.String())
		return idx + 2, false, nil
	}
	h.Set(key, string(value))
	return idx + 2, false, nil
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	h[key] = value
}

func (h Headers) Override(key, value string) {
	key = strings.ToLower(key)
	h[key] = value
}

func (h Headers) Remove(key string) {
	key = strings.ToLower(key)
	delete(h, key)
}

func (h Headers) Get(key string) (string, bool) {
	key = strings.ToLower(key)
	value, exists := h[key]
	return value, exists
}

func normalizeHead(data string) (string, error) {
	allowedSpecialChars := "!#$%&'*+-.^_`|~"
	dataBytes := []byte(data)
	for i, c := range dataBytes {

		if c >= 'a' && c <= 'z' {
			continue
		}

		if c >= 'A' && c <= 'z' {
			dataBytes[i] = c + 32
			continue
		}

		if c >= '0' && c <= '9' {
			continue
		}

		if strings.ContainsRune(allowedSpecialChars, rune(c)) {
			continue
		}
		fmt.Println("incorrect char at ", i)
		return data, fmt.Errorf("in correct character found: ")

	}

	return string(dataBytes), nil
}

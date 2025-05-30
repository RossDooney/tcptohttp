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

	n, err = parseRequestLine(data)
	if err != nil {
		// something actually went wrong
		return 0, false, err
	}
	if n == 0 {
		// just need more data
		return 0, false, nil
	}

	return n, true, nil

}

func parseRequestLine(data []byte) (int, error) {
	idx := bytes.Index(data, []byte(crlf))
	if idx == -1 {
		return 0, nil
	}
	requestLineText := string(data[:idx])
	err := requestLineFromString(requestLineText)
	if err != nil {
		return 0, err
	}
	return idx + 2, nil
}

func requestLineFromString(str string) error {
	before, after, _ := strings.Cut(str, ":")

	fmt.Println(before)
	fmt.Println(after)

	return nil
}

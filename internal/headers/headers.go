package headers

import (
	"bytes"
	"fmt"
)

type parserState string

const (
	StateInit  parserState = "init"
	StateDone  parserState = "done"
	StateError parserState = "error"
)

var CLRF = []byte("\r\n")
var whiteSpac = []byte{':'}

type Headers map[string]string

func NewHeaders() (h Headers) {
	return make(Headers, 0)
}
func (h Headers) Parse(data []byte) (int, bool, error) {
	read := 0
	done := false
	for {
		idx := bytes.Index(data[read:], CLRF)
		if idx == -1 {
			break
		}
		//EmptyHeader
		if idx == 0 {
			done = true
			read += len(CLRF)
			break
		}
		name, value, err := parseHeader(data[read : read+idx])
		if err != nil {
			return 0, false, err
		}
		read += idx + len(CLRF)
		h[name] = value
	}
	return read, done, nil
}

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("Malformed RequestLine")
	}
	name := bytes.TrimPrefix(parts[0], []byte(" "))
	value := bytes.TrimSpace(parts[1])
	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", fmt.Errorf("Malformed header Name")
	}
	return string(name), string(value), nil
}

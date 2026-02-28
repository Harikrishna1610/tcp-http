package request

import (
	"bytes"
	"fmt"
	"http_frm_udp/internal/headers"
	"io"
	"strings"
)

type parserState string

const (
	StateInit    parserState = "init"
	StateDone    parserState = "done"
	StateHeaders parserState = "HeadersState"
	StateError   parserState = "error"
	StateBody    parserState = "body"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	Body        string
	state       parserState
}

func (r *RequestLine) validateHTTP() bool {
	return r.HttpVersion == "HTTP/1.1"
}

func GetRequest() *Request {
	return &Request{state: StateInit, Headers: headers.NewHeaders(), Body: ""}
}

var MALFORMED_HTTP = fmt.Errorf("Http request is malformed")
var UNSUPPORTED_VERSION = fmt.Errorf("Unsupported http version")
var SEPARATOR = []byte("\r\n")

func (r *Request) parse(data []byte) (int, error) {
	read := 0
OUTER:
	for {
		currentData := data[read:]
		if len(currentData) == 0 {
			break
		}
		switch r.state {
		case StateError:
			return 0, fmt.Errorf("Request in Error State")
		case StateInit:
			rl, n, err := ParseRequest(currentData)
			if err != nil {
				return 0, err
			}
			if n == 0 {
				break OUTER
			}
			r.RequestLine = *rl
			read += n
			r.state = StateHeaders
		case StateHeaders:
			n, done, err := r.Headers.Parse(currentData)
			if err != nil {
				r.state = StateError
				return 0, err
			}
			if n == 0 {
				break OUTER
			}
			read += n
			if done {
				r.state = StateBody
			}
		case StateBody:
			contentLen, _ := r.Headers.GetContentLen()

			if contentLen == 0 {
				r.state = StateDone
				break OUTER
			}
			remaining := min(contentLen-len(r.Body), len(currentData))
			r.Body += string(currentData[:remaining])
			read += remaining
			if len(r.Body) == contentLen {
				r.state = StateDone
			}
		case StateDone:
			break OUTER
		}
	}
	return read, nil

}

func (r *Request) done() bool {
	return r.state == StateDone || r.state == StateError
}
func ParseRequest(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPARATOR)
	if idx == -1 {
		return nil, 0, nil
	}
	startLine := b[:idx]
	read := idx + len(SEPARATOR)

	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return &RequestLine{}, 0, MALFORMED_HTTP
	}
	rl := &RequestLine{Method: string(parts[0]), RequestTarget: string(parts[1]), HttpVersion: strings.Split(string(parts[2]), "/")[1]}
	return rl, read, nil
}
func RequestFromReader(reader io.Reader) (*Request, error) {
	//Note: Buffer could get
	var b = make([]byte, 1024)
	bIdx := 0
	var req = GetRequest()
	for !req.done() {
		n, err := reader.Read(b[bIdx:])
		if err == io.EOF {
			break
		}
		if err != nil {
			return req, err
		}
		bIdx += n
		readN, err := req.parse(b[:bIdx])
		if err != nil {
			return nil, err
		}
		copy(b, b[readN:bIdx])
		bIdx -= readN
	}
	if val, _ := req.Headers.GetContentLen(); val != len(req.Body) {
		return req, fmt.Errorf("Body Content Incomplete")
	}
	return req, nil
}

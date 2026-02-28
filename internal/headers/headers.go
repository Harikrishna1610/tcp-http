package headers

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var SymbolRangeTable = &unicode.RangeTable{
	R16: []unicode.Range16{
		{Lo: '!', Hi: '!', Stride: 1},
		{Lo: '#', Hi: '#', Stride: 1},
		{Lo: '$', Hi: '$', Stride: 1},
		{Lo: '%', Hi: '%', Stride: 1},
		{Lo: '&', Hi: '&', Stride: 1},
		{Lo: '\'', Hi: '\'', Stride: 1},
		{Lo: '*', Hi: '*', Stride: 1},
		{Lo: '+', Hi: '+', Stride: 1},
		{Lo: '-', Hi: '-', Stride: 1},
		{Lo: '.', Hi: '.', Stride: 1},
		{Lo: '^', Hi: '^', Stride: 1},
		{Lo: '_', Hi: '_', Stride: 1},
		{Lo: '`', Hi: '`', Stride: 1},
		{Lo: '|', Hi: '|', Stride: 1},
		{Lo: '~', Hi: '~', Stride: 1},
	},
	LatinOffset: 15,                  // All 15 entries have Hi <= MaxLatin1
	R32:         []unicode.Range32{}, // No 32-bit code points
}

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
func (h Headers) GetContentLen() (int, error) {
	val, ok := h["content-length"]
	if ok {
		i, err := strconv.Atoi(val)
		if err != nil {
			return 0, err
		}
		return i, nil
	}
	return 0, nil
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
		if h[name] != "" {
			h[name] = strings.Join([]string{h[name], value}, ", ")
		} else {
			h[name] = value
		}

	}
	return read, done, nil
}

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("Malformed RequestLine")
	}
	name := bytes.ToLower(bytes.TrimPrefix(parts[0], []byte(" ")))
	if len(name) <= 0 || !isValidbyte(name) {
		return "", "", fmt.Errorf("Malformed header key")

	}
	value := bytes.TrimSpace(parts[1])
	if bytes.HasSuffix(name, []byte(" ")) {
		return "", "", fmt.Errorf("Malformed header Name")
	}
	return string(name), string(value), nil
}

func isValidbyte(char []byte) bool {
	for _, b := range char {
		if !(unicode.IsDigit(rune(b)) || unicode.IsLetter(rune(b)) || unicode.In(rune(b), SymbolRangeTable)) {
			return false
		}
	}
	return true
}

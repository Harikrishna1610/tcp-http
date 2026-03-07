package response

import (
	"fmt"
	"http_frm_udp/internal/headers"
	"io"
)

type Writer struct {
	Writer io.Writer
}

type StatusCode int

const (
	StatusOkay                StatusCode = 200
	StatusNotFound            StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func NewWriter(w io.Writer) *Writer {
	return &Writer{Writer: w}
}
func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h["Content-Length"] = fmt.Sprintf(": %d\r\n", contentLen)
	h["Connection"] = ": close\r\n"
	h["Content-Type"] = ": text/plain\r\n"
	return h
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	statusLine := []byte{}

	switch statusCode {
	case StatusOkay:
		statusLine = []byte("HTTP/1.1 200 OK\r\n")

	case StatusNotFound:
		statusLine = []byte("HTTP/1.1 400 Bad Request\r\n")

	case StatusInternalServerError:
		statusLine = []byte("HTTP/1.1 500 Internal Server Error\r\n")
	}
	_, err := w.Writer.Write(statusLine)
	return err
}
func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if len(headers) == 0 {
		return fmt.Errorf("No Headers")
	}
	for k, h := range headers {
		_, err := w.Writer.Write([]byte(k + h))
		if err != nil {
			return err
		}
	}
	w.Writer.Write([]byte("\r\n"))
	return nil
}
func (w *Writer) WriteBody(p []byte) (int, error) {
	n, err := w.Writer.Write(p)
	return n, err
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	return nil
}

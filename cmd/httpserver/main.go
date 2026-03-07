package main

import (
	"crypto/sha256"
	"fmt"
	"http_frm_udp/internal/headers"
	response "http_frm_udp/internal/reponse"
	"http_frm_udp/internal/request"
	"http_frm_udp/internal/server"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

const port = 42069

func getResponse200() []byte {
	return []byte(`<html>
  <head>
    <title>200 OK</title>
  </head>
  <body>
    <h1>Success!</h1>
    <p>Your request was an absolute banger.</p>
  </body>
</html>`)
}
func getResponse400() []byte {
	return []byte(`<html>
	<head>
		<title>400 Bad Request</title>
	</head>
	<body>
		<h1>Bad Request</h1>
		<p>Your request honestly kinda sucked.</p>
	</body>
	</html>`)
}
func getResponse500() []byte {
	return []byte(`<html>
  <head>
    <title>500 Internal Server Error</title>
  </head>
  <body>
    <h1>Internal Server Error</h1>
    <p>Okay, you know what? This one is on me.</p>
  </body>
</html>`)
}
func main() {
	server, err := server.Serve(port, func(w *response.Writer, req *request.Request) {
		status := response.StatusOkay
		body := getResponse200()
		h := response.GetDefaultHeaders(0)
		if req.RequestLine.RequestTarget == "/yourproblem" {
			status = response.StatusNotFound
			body = getResponse400()
		} else if req.RequestLine.RequestTarget == "/myproblem" {
			status = response.StatusInternalServerError
			body = getResponse500()
		} else if req.RequestLine.RequestTarget == "/video" {
			h["Content-Type"] = ": video/mp4\r\n"
			w.WriteStatusLine(status)
			n, err := os.ReadFile("/media/hari/Data/projects/Learning/GOLang/Http/assets/vim.mp4")
			if err != nil {
				log.Fatal(err)
			}
			h["Content-Length"] = fmt.Sprintf(": %d\r\n", len(n))
			w.WriteHeaders(h)
			w.WriteBody(n)
			return

		} else if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
			target := req.RequestLine.RequestTarget
			res, err := http.Get("https://httpbin.org/" + target[len("/httpbin/"):])
			if err != nil {
				status = response.StatusInternalServerError
				body = getResponse500()
			} else {
				h["Transfer-Encoding"] = ": chunked\r\n"
				h["Content-Type"] = ": text/plain\r\n"
				delete(h, "Content-Length")
				h["Trailer"] = ": X-Content-SHA256\r\n"
				w.WriteStatusLine(status)
				w.WriteHeaders(h)
				fullBody := []byte{}
				for {
					data := make([]byte, 32)
					n, err := res.Body.Read(data)
					if err != nil {
						break
					}
					fullBody = append(fullBody, data[:n]...)
					w.WriteBody(fmt.Appendf(nil, "%x\r\n", n))
					w.WriteBody(data[:n])
					w.WriteBody([]byte("\r\n"))
				}
				w.WriteBody([]byte("0\r\n"))
				trailers := headers.NewHeaders()
				sum := sha256.Sum256(fullBody)
				trailers["X-Content-SHA256"] = fmt.Sprintf(": %x\r\n", sum)
				trailers["X-Content-Length"] = ": " + strconv.Itoa(len(fullBody))
				w.WriteHeaders(trailers)
				return
			}

		}
		h["Content-Length"] = fmt.Sprintf(": %d\r\n", len(body))
		h["Content-Type"] = ": text/html\r\n"
		w.WriteStatusLine(status)
		w.WriteHeaders(h)
		w.WriteBody(body)
	})
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}

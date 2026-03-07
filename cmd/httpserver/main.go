package main

import (
	"fmt"
	response "http_frm_udp/internal/reponse"
	"http_frm_udp/internal/request"
	"http_frm_udp/internal/server"
	"log"
	"os"
	"os/signal"
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
		switch req.RequestLine.RequestTarget {
		case "/yourproblem":
			status = response.StatusNotFound
			body = getResponse400()
			//break
		case "/myproblem":
			status = response.StatusInternalServerError
			body = getResponse500()
			//break
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

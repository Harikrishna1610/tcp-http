package server

import (
	"fmt"
	response "http_frm_udp/internal/reponse"
	"http_frm_udp/internal/request"
	"io"
	"net"
)

type Server struct {
	isClosed bool
	handler  Handler
}
type HandlerError struct {
	StatusCode response.StatusCode
	Message    []byte
}
type Handler func(w *response.Writer, req *request.Request)

func runConnection(s *Server, conn io.ReadWriteCloser) {
	defer conn.Close()
	//writer := bytes.NewBuffer([]byte{})

	h := response.GetDefaultHeaders(0)
	res := response.NewWriter(conn)
	r, err := request.RequestFromReader(conn)
	if err != nil {
		res.WriteStatusLine(response.StatusNotFound)
		res.WriteHeaders(h)
		return
	}
	s.handler(res, r)
}
func runServer(s *Server, listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		if s.isClosed {
			return
		}
		go runConnection(s, conn)
	}
}
func Serve(port uint16, handler Handler) (*Server, error) {
	f, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	server := &Server{
		isClosed: false,
		handler:  handler}

	go runServer(server, f)
	return server, nil
}

func (s *Server) Close() error {
	s.isClosed = true
	return nil
}

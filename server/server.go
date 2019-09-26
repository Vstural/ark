package server

import (
	"ark/protocol"
	"errors"
	"net"
)

type Server struct {
	handler protocol.ServerHandler
}

func NewServer(h protocol.ServerHandler) (*Server, error) {
	if h == nil {
		return nil, errors.New("handler should not be nil")
	}
	s := &Server{
		handler: h,
	}
	return s, nil
}

func (s *Server) Serve(ln net.Listener) error {
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go s.handler.HandleServerConn(conn)
	}
}

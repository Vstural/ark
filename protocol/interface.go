package protocol

import (
	"net"

	"github.com/Vstural/socksgo"
)

type ClientHandler interface {
	HandleClientConn(shakeinfo *socksgo.HandShakeInfo, conn net.Conn) (targetAddr string, err error)
}

type ServerHandler interface {
	HandleServerConn(conn net.Conn)
}

type Handler interface {
	ClientHandler
	ServerHandler
}

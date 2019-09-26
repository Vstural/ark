package protocol

import (
	"io"
	"net"

	"github.com/Vstural/socksgo"
)

const MTUSIZE = 1400

type RawLocalHandler struct{}

func NewRawLocalHandler() *RawLocalHandler {
	return &RawLocalHandler{}
}

// HandleConn raw processer here
// connection to target server directly
func (r *RawLocalHandler) HandleClientConn(shakeinfo *socksgo.HandShakeInfo, conn net.Conn) (string, error) {
	targetConn, err := net.Dial("tcp", shakeinfo.SockRequest.GetDstAddr())
	if err != nil {
		return "", err
	}

	// optimize? mtu detect
	// a random mtu maybe helpful to pass fire wall
	go io.CopyBuffer(conn, targetConn, make([]byte, MTUSIZE))
	go io.CopyBuffer(targetConn, conn, make([]byte, MTUSIZE))
	return targetConn.RemoteAddr().String(), nil
}

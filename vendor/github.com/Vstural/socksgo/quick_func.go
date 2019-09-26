package socksgo

import (
	"errors"
	"fmt"
	"io"
	"net"
)

type HandShakeInfo struct {
	MethodSelect *MethodSelection
	SockRequest  *SockRequest
}

var VER5_NOAUTH_REPLY = (&MethodSelectionReply{
	Ver:    Socks5version,
	Method: MethodNoAuthenticationRequired,
}).ToByte()

func ReadSocksRequest(conn net.Conn) (*HandShakeInfo, error) {
	methodSelect, err := readMethodSelection(conn)
	if err != nil {
		return nil, err
	}
	if methodSelect.Ver != Socks5version {
		return nil, errors.New("only ver 5 is support")
	}
	// skip method check
	conn.Write(VER5_NOAUTH_REPLY)

	req, err := readSockRequest(conn)
	if err != nil {
		return nil, err
	}

	return &HandShakeInfo{
		MethodSelect: methodSelect,
		SockRequest:  req,
	}, nil
}

// WriteSuccReplyQuick write is connection establish success
func WriteSuccReplyQuick(writer io.Writer,
	handShakeInfo *HandShakeInfo) (n int, err error) {
	return WriteSuccReply(writer,
		handShakeInfo.SockRequest.Atyp,
		handShakeInfo.SockRequest.DstAddr,
		handShakeInfo.SockRequest.DstPort)
}

// WriteSuccReply write is connection establish success
func WriteSuccReply(writer io.Writer,
	atyp uint8,
	dstAddr []uint8,
	dstPort uint16) (n int, err error) {
	resp := &SockRequestReply{
		Ver:     Socks5version,
		Rep:     Succeeded,
		Rev:     Rsv,
		Atyp:    atyp,
		DstAddr: dstAddr,
		DstPort: dstPort,
	}
	return writer.Write(resp.ToByte())
}

// WriteFailReplyQuick write is connection establish success
// in most of case, we do not care the fail reason
func WriteFailReplyQuick(writer io.Writer,
	handShakeInfo *HandShakeInfo) (n int, err error) {
	return WriteFailReply(writer,
		handShakeInfo.SockRequest.Atyp,
		handShakeInfo.SockRequest.DstAddr,
		handShakeInfo.SockRequest.DstPort)
}

// WriteFailReply write is connection establish success
func WriteFailReply(writer io.Writer,
	atyp uint8,
	dstAddr []uint8,
	dstPort uint16) (n int, err error) {
	resp := &SockRequestReply{
		Ver:     Socks5version,
		Rep:     SockServerFailure,
		Rev:     Rsv,
		Atyp:    atyp,
		DstAddr: dstAddr,
		DstPort: dstPort,
	}
	return writer.Write(resp.ToByte())
}

// GetAtypType return atyp type name
func GetAtypType(val uint8) string {
	switch val {
	case AtypDomainName:
		return "AtypDomainName"
	case AtypIPV4:
		return "AtypIPV4"
	case AtypIPV6:
		return "AtypIPV6"
	default:
		return fmt.Sprintf("unknown atyp type %d", val)
	}
}

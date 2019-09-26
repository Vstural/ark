package socksgo

import (
	"encoding/binary"
	"fmt"
	"net"
)

const (
	MethodNoAuthenticationRequired = 0x00
	GSSAPI                         = 0x01
	UsernamePW                     = 0x03
	// IanaAssigned             = 0x04 - 0x07
	NoAcceptableMethods = 0xff

	Socks5version = 0x05

	CmdConnect      = 0x01
	CmdBind         = 0x02
	CmdUdpAssociate = 0x03

	Rsv = 0x00

	AtypIPV4       = 0x01
	AtypDomainName = 0x03
	AtypIPV6       = 0x04

	Succeeded          = 0x00
	SockServerFailure  = 0x01
	ConnectionNotAllow = 0x02
	NetworkUnreachable = 0x03
	HostUnreachable    = 0x04
	ConnectionRefused  = 0x05
	TTLexpired         = 0x06
	CmdNotSupport      = 0x07
	AddrNotSupport     = 0x08
)

type MethodSelection struct {
	Ver      uint8
	NMethods uint8
	Methods  []uint8
}

// +----+----------+----------+
// |VER | NMETHODS | METHODS  |
// +----+----------+----------+
// | 1  |    1     | 1 to 255 |
// +----+----------+----------+

func readMethodSelection(conn net.Conn) (*MethodSelection, error) {
	ms := &MethodSelection{}
	err := binary.Read(conn, binary.BigEndian, &ms.Ver)
	if err != nil {
		return nil, err
	}
	err = binary.Read(conn, binary.BigEndian, &ms.NMethods)
	if err != nil {
		return nil, err
	}
	ms.Methods = make([]uint8, ms.NMethods)
	err = binary.Read(conn, binary.BigEndian, &ms.NMethods)
	if err != nil {
		return nil, err
	}
	return ms, nil
}

// +----+--------+
// |VER | METHOD |
// +----+--------+
// | 1  |   1    |
// +----+--------+

type MethodSelectionReply struct {
	Ver    uint8
	Method uint8
}

func (m *MethodSelectionReply) ToByte() []byte {
	resp := make([]byte, 2)
	resp[0] = m.Ver
	resp[1] = m.Method
	return resp
}

// +----+-----+-------+------+----------+----------+
// |VER | CMD |  RSV  | ATYP | DST.ADDR | DST.PORT |
// +----+-----+-------+------+----------+----------+
// | 1  |  1  | X'00' |  1   | Variable |    2     |
// +----+-----+-------+------+----------+----------+

type SockRequest struct {
	Ver     uint8
	Cmd     uint8
	Rsv     uint8
	Atyp    uint8
	DstAddr []uint8
	DstPort uint16
}

func (s *SockRequest) GetDstAddr() string {
	return fmt.Sprintf("%s:%d", s.DstAddr, s.DstPort)
}

// wow so many if err
func readSockRequest(conn net.Conn) (*SockRequest, error) {
	sr := &SockRequest{}
	err := binary.Read(conn, binary.BigEndian, &sr.Ver)
	if err != nil {
		return nil, err
	}
	err = binary.Read(conn, binary.BigEndian, &sr.Cmd)
	if err != nil {
		return nil, err
	}
	err = binary.Read(conn, binary.BigEndian, &sr.Rsv)
	if err != nil {
		return nil, err
	}
	err = binary.Read(conn, binary.BigEndian, &sr.Atyp)
	if err != nil {
		return nil, err
	}
	if sr.Atyp == AtypIPV4 {
		sr.DstAddr = make([]uint8, 4)
	} else if sr.Atyp == AtypDomainName {
		var srlength uint8
		err = binary.Read(conn, binary.BigEndian, &srlength)
		if err != nil {
			return nil, err
		}
		sr.DstAddr = make([]uint8, srlength)
	} else if sr.Atyp == AtypIPV6 {
		sr.DstAddr = make([]uint8, 16)
	}
	err = binary.Read(conn, binary.BigEndian, &sr.DstAddr)
	if err != nil {
		return nil, err
	}
	err = binary.Read(conn, binary.BigEndian, &sr.DstPort)
	if err != nil {
		return nil, err
	}
	return sr, nil
}

// +----+-----+-------+------+----------+----------+
// |VER | REP |  RSV  | ATYP | BND.ADDR | BND.PORT |
// +----+-----+-------+------+----------+----------+
// | 1  |  1  | X'00' |  1   | Variable |    2     |
// +----+-----+-------+------+----------+----------+

type SockRequestReply struct {
	Ver     uint8
	Rep     uint8
	Rev     uint8
	Atyp    uint8
	DstAddr []uint8
	DstPort uint16
}

func (s *SockRequestReply) ToByte() []byte {

	var addrlen = 0

	bytePort := make([]byte, 2)
	binary.BigEndian.PutUint16(bytePort, s.DstPort)
	if s.DstAddr != nil {
		addrlen = len(s.DstAddr)
	}
	buf := make([]byte, 4+addrlen+2)

	buf[0] = s.Ver
	buf[1] = s.Rep
	buf[2] = s.Rev
	buf[3] = s.Atyp
	// buf[4 : 4+addrlen] = s.dstAddr
	if s.DstAddr != nil {
		// buf[4+addrlen:] = bytePort
		copy(buf[4:4+addrlen], s.DstAddr)
		copy(buf[4+addrlen:], bytePort)
	}

	return buf
}

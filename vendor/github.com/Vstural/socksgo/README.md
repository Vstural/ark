# A implement of socks5 server side

`go get github.com/Vstural/socksgo`

## use

see example

```go
package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"

	"github.com/Vstural/socksgo"
)

const (
	ListenAddr = "localhost:9879"
)

func main() {
	ln, err := net.Listen("tcp", ListenAddr)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Print(err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	// first we read socks5 request
	sockReq, err := socksgo.ReadSocksRequest(conn)
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}

	reqtype := sockReq.SockRequest.Atyp
	dstaddr := sockReq.SockRequest.GetDstAddr()

	fmt.Printf("reqtype: %s, dstaddr: %s\n", socksgo.GetAtypType(reqtype), dstaddr)

	// then we try to connection to target server
	// if success, reply succ
	// or reply fail and close connection
	err = handleRequestRaw(sockReq, conn)
	if err != nil {
		log.Printf("dial to targetr server fail: %s", err.Error())
		conn.Close()
		return
	}
	return
}

// here we just repost requet to target server directly
func handleRequestRaw(shakeinfo *socksgo.HandShakeInfo, conn net.Conn) error {
	targetConn, err := net.Dial("tcp", shakeinfo.SockRequest.GetDstAddr())
	if err != nil {
		socksgo.WriteFailReplyQuick(conn, shakeinfo)
		return err
	}

	ip, port := parseAddr(targetConn.LocalAddr().String())
	socksgo.WriteSuccReply(
		conn,
		socksgo.AtypIPV4,
		ip.To4(),
		uint16(port),
	)
	go io.CopyBuffer(conn, targetConn, make([]byte, 1500))
	go io.CopyBuffer(targetConn, conn, make([]byte, 1500))

	return nil
}

func parseAddr(addr string) (ip net.IP, port int) {
	var i int

	for i = 0; i < len(addr); i++ {
		if addr[i] == ':' {
			break
		}
	}
	ipstr := addr[:i]
	portstr := addr[i+1:]
	ip = net.ParseIP(ipstr)

	port, err := strconv.Atoi(portstr)

	if err != nil {
		log.Println(err)
		return nil, 0
	}
	return ip, port
}
```
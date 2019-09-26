package client

import (
	"ark/protocol"
	"ark/tools"
	"errors"
	"log"
	"net"

	"github.com/Vstural/socksgo"
)

// Client base client repost reqeust to proxy server
type Client struct {
	handler protocol.ClientHandler
}

// NewClient return a new client
func NewClient(f protocol.ClientHandler) (*Client, error) {
	if f == nil {
		return nil, errors.New("client protocol handler should not be nil")
	}
	c := &Client{
		handler: f,
	}
	return c, nil
}

// Serve start serve!
func (c *Client) Serve(ln net.Listener) error {
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go c.handleConn(conn)
	}
}

func (c *Client) handleConn(conn net.Conn) {
	sockReq, err := socksgo.ReadSocksRequest(conn)
	if err != nil {
		log.Println(err)
		conn.Close()
		return
	}

	reqtype := sockReq.SockRequest.Atyp
	dstaddr := sockReq.SockRequest.GetDstAddr()

	log.Printf("reqtype: %s, dstaddr: %s\n", socksgo.GetAtypType(reqtype), dstaddr)
	if c.handler != nil {
		targetAddr, err := c.handler.HandleClientConn(sockReq, conn)
		if err != nil {
			log.Printf("dial to target server fail: %s", err.Error())
			conn.Close()
			return
		}
		if err != nil {
			socksgo.WriteFailReplyQuick(conn, sockReq)
		} else {
			ip, port := tools.ParseAddr(targetAddr)
			socksgo.WriteSuccReply(
				conn,
				socksgo.AtypIPV4,
				ip.To4(),
				uint16(port),
			)
		}
	}
	return
}

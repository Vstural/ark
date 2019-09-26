package protocol

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/Vstural/socksgo"
)

type RequestInfo struct {
	DstAddr string `json:"dst_addr"`
}

type ResponseInfo struct {
	TargetAddr string `json:"target_addr"`
}

type RawRemoteHandlerOption struct {
	ServerAddr string
}

type RawRemoteHandler struct {
	options RawRemoteHandlerOption
}

func NewRawRemoteHandler(options RawRemoteHandlerOption) *RawRemoteHandler {
	return &RawRemoteHandler{
		options: options,
	}
}

var lf = []byte("\n")

func (r *RawRemoteHandler) HandleClientConn(shakeinfo *socksgo.HandShakeInfo,
	conn net.Conn) (string, error) {
	targetAddr := shakeinfo.SockRequest.GetDstAddr()

	request := &RequestInfo{
		DstAddr: targetAddr,
	}

	reqByte, err := json.Marshal(request)
	if err != nil {
		return "", fmt.Errorf("marshal request with addr %s fail: %s",
			targetAddr, err.Error())
	}

	serverConn, err := net.Dial("tcp", r.options.ServerAddr)
	if err != nil {
		return "", fmt.Errorf("dail to server with addr %s fail: %s",
			r.options.ServerAddr, err.Error())
	}

	_, err = serverConn.Write(append(reqByte, lf...))
	if err != nil {
		return "", err
	}

	reader := bufio.NewReader(serverConn)
	respByte, _, err := reader.ReadLine()
	if err != nil {
		return "", err
	}
	resp := &ResponseInfo{}
	err = json.Unmarshal(respByte, resp)
	if err != nil {
		return "", err
	}
	log.Printf("write request header finish and get target addr %s",
		resp.TargetAddr)

	go io.CopyBuffer(conn, reader, make([]byte, MTUSIZE))
	go io.CopyBuffer(serverConn, conn, make([]byte, MTUSIZE))
	return resp.TargetAddr, nil
}

func (r *RawRemoteHandler) HandleServerConn(conn net.Conn) {
	log.Println("got connection reqeuest")

	reader := bufio.NewReader(conn)
	l, _, err := reader.ReadLine()
	if err != nil {
		log.Printf("read line fail: %s", err.Error())
		conn.Close()
		return
	}
	request := &RequestInfo{}
	err = json.Unmarshal(l, request)
	if err != nil {
		log.Printf("unmarshal request header fail: %s", err.Error())
		conn.Close()
		return
	}
	targetConn, err := net.Dial("tcp", request.DstAddr)
	if err != nil {
		log.Printf("dial to addr %s fail : %s", request.DstAddr, err.Error())
		conn.Close()
		return
	}
	log.Printf("dial to %s success", request.DstAddr)

	resp := &ResponseInfo{
		TargetAddr: targetConn.RemoteAddr().String(),
	}
	respByte, err := json.Marshal(resp)
	if err != nil {
		log.Printf("marshal response with target addr %s fail: %s",
			resp.TargetAddr, err.Error())
		conn.Close()
		return
	}
	conn.Write(append(respByte, lf...))
	log.Printf("write response with target addr %s success",
		resp.TargetAddr)
	// process data transport now
	go io.CopyBuffer(conn, targetConn, make([]byte, MTUSIZE))
	go io.CopyBuffer(targetConn, reader, make([]byte, MTUSIZE))
	return
}

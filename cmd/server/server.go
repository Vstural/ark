package main

import (
	"ark/protocol"
	"ark/server"
	"net"
)

const ln_addr = "0.0.0.0:5001"

func main() {
	ln, err := net.Listen("tcp", ln_addr)
	if err != nil {
		panic(err)
	}
	ser, err := server.NewServer(protocol.NewRawRemoteHandler(protocol.RawRemoteHandlerOption{}))
	if err != nil {
		panic(err)
	}
	ser.Serve(ln)
	return
}

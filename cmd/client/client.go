package main

import (
	"ark/client"
	"ark/protocol"
	"net"
)

const ln_addr = "127.0.0.1:5000"

func main() {
	ln, err := net.Listen("tcp", ln_addr)
	if err != nil {
		panic(err)
	}

	// todo: config conn handler

	client, err := client.NewClient(
		protocol.NewRawRemoteHandler(
			protocol.RawRemoteHandlerOption{
				ServerAddr: "127.0.0.1:5001",
			}))
	// client, err := client.NewClient(protocol.NewRawLocalHandler())
	if err != nil {
		panic(err)
	}
	client.Serve(ln)
	return
}

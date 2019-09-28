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

	// run example code with server
	client, err := client.NewClient(
		protocol.NewRawRemoteHandler(
			protocol.RawRemoteHandlerOption{
				ServerAddr: "127.0.0.1:5001",
			}))

	// run example code in local mode
	// client, err := client.NewClient(protocol.NewRawLocalHandler())

	if err != nil {
		panic(err)
	}
	client.Serve(ln)
	return
}

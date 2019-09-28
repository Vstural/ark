# ark

a simple socks5 proxy example use [socksgo](https://github.com/Vstural/socksgo)

## run 

main func in `cmd/client/client.go` and `cmd/server/server.go`

simply modify client code to change start mode

```go
// run example code with server
client, err := client.NewClient(
    protocol.NewRawRemoteHandler(
        protocol.RawRemoteHandlerOption{
            ServerAddr: "127.0.0.1:5001",
        }))

// run example code in local mode
// client, err := client.NewClient(protocol.NewRawLocalHandler())
```

## custom protocol

write your own protocol under `protocol` folder and implement  `Handler` interface in `protocol/interface.go`

```go
type Handler interface {
	ClientHandler
	ServerHandler
}
```

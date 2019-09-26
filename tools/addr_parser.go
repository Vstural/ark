package tools

import (
	"log"
	"net"
	"strconv"
)

// tool func
func ParseAddr(addr string) (ip net.IP, port int) {
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

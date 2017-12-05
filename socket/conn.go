package socket

import "net"

type Conn struct {
	conn       net.Conn
	connID     int
	status     int32
	userData   interface{}
	localAddr  string
	remoteAddr string
}

type ConnEvent struct {
	Type int
	Conn net.Conn
	Data interface{}
}

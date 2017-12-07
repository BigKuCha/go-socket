package socket

import "net"

const (
	MSG_TYPE_ACK = iota // 响应消息类型， 用于服务器和客户端交互userid和connid
	MSG_TYPE_DATA
)

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
	Conn Conn
	Data []byte
}

type Msg struct {
	msgType int
	data    map[string]interface{}
}

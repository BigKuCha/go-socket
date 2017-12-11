package gosocket

import "net"

const (
	MSG_TYPE_ACK = iota // 响应消息类型， 用于服务器和客户端交互userid和connid
	MSG_TYPE_CHAT
)

type Conn struct {
	conn       net.Conn
	connID     uint32
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
	MsgType int
	Data    map[string]string
}

type ChatMsg struct {
	MsgType int
	FromID  int
	ToID    int
	Data    []byte
}

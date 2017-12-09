package socket

import (
	"net"
)

const (
	MSG_TYPE_ACK  = iota // 响应消息类型， 用于服务器和客户端交互userid和connid
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
	MsgType int
	Data    map[string]string
}

//func (c *Conn) SendMsg(msg Msg) (n int, err error) {
//	msgJson, _ := json.Marshal(msg)
//	msgByte := []byte(msgJson)
//	var msgHead [4]byte
//	binary.BigEndian.PutUint32(msgHead[0:], uint32(len(msgByte)))
//	msgBody := append(msgHead[0:], []byte(msgJson)...)
//	n, err = c.conn.Write(msgBody)
//	return
//}

package gosocket

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net"
)

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

type ChatMsg struct {
	MsgType int
	FromID  int
	ToID    int
	Data    []byte
}

func (c *Conn) SendMsg(msg ChatMsg) (n int, err error) {
	msgBody := SerialMsg(msg)
	n, err = c.conn.Write(msgBody)
	return
}

func (c *Conn) GetRemoteAddr() string {
	return c.remoteAddr
}

func SerialMsg(msg ChatMsg) []byte {
	msgJson, _ := json.Marshal(msg)
	msgByte := []byte(msgJson)
	var msgHead [4]byte
	binary.BigEndian.PutUint32(msgHead[0:], uint32(len(msgByte)))
	msgBody := append(msgHead[0:], []byte(msgJson)...)
	return msgBody
}

func HandleMsg(b []byte) (ChatMsg, error) {
	msgHeader := b[:4]
	msgLength := binary.BigEndian.Uint32(msgHeader)
	var msg ChatMsg
	err := json.Unmarshal(b[4:msgLength+4], &msg)
	if err != nil {
		fmt.Printf("解码错误: %+v \n", err)
		return msg, err
	}
	return msg, err
}

package socket

import (
	"encoding/json"
	"encoding/binary"
)

type NetWork struct {
	conn Conn
}

func (net *NetWork) SendMsg(msg Msg) (n int, err error) {
	msgJson, _ := json.Marshal(msg)
	msgByte := []byte(msgJson)
	var msgHead [4]byte
	binary.BigEndian.PutUint32(msgHead[0:], uint32(len(msgByte)))
	msgBody := append(msgHead[0:], []byte(msgJson)...)
	n, err = net.conn.conn.Write(msgBody)
	return
}

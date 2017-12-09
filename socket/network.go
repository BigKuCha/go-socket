package socket

import (
	"encoding/json"
	"encoding/binary"
)

type NetWork struct {
	conn Conn
}

func serialMsg(msg ChatMsg) []byte {
	msgJson, _ := json.Marshal(msg)
	msgByte := []byte(msgJson)
	var msgHead [4]byte
	binary.BigEndian.PutUint32(msgHead[0:], uint32(len(msgByte)))
	msgBody := append(msgHead[0:], []byte(msgJson)...)
	return msgBody
}

func (net *NetWork) SendMsg(msg ChatMsg) (n int, err error) {
	msgBody := serialMsg(msg)
	n, err = net.conn.conn.Write(msgBody)
	return
}

func handleMsg(b []byte) (ChatMsg, error) {
	msgHeader := b[:4]
	msgLength := binary.BigEndian.Uint32(msgHeader)
	var msg ChatMsg
	err := json.Unmarshal(b[4: msgLength+4], &msg)
	if err != nil {
		return msg, err
	}
	return msg, err
}

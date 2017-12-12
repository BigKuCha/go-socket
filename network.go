package gosocket

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
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

func handleMsg(b []byte) (ChatMsg, error) {
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

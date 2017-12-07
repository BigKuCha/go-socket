package socket

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
)

/*连接状态*/
const (
	CONN_STATUS_NONE = iota
	CONN_STATUS_CONNECTED
	CONN_STATUS_DISCONNECTED
)

//conn                net.Conn
//status              int32
//connId              int
//sendMsgQueue        chan *sendTask
//sendTimeoutSec      int
//eventQueue          IEventQueue
//streamProtocol      IStreamProtocol
//maxReadBufferLength int
//userdata            interface{}
//from                int
//readTimeoutSec      int
//fnSyncExecute       FuncSyncExecute
//unpacker            IUnpacker
//disableSend         int32
//localAddr           string
//remoteAddr          string

type client struct {
	userID       int
	connID       int64
	Conn         net.Conn
	eventQueue   chan ConnEvent
	OnConnect    func(event ConnEvent)
	OnData       func(event ConnEvent)
	OnDisconnect func(event ConnEvent)
}

func NewClient(userID int) *client {
	return &client{
		userID:     userID,
		eventQueue: make(chan ConnEvent),
	}
}

func (c *client) Connect(host string, port int) error {
	conn, err := net.Dial("tcp", net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		return err
	}
	c.Conn = conn
	go c.handleEvent()
	go c.handleClientConn()
	event := ConnEvent{
		Type: EVT_ON_CONNECT,
	}
	c.eventQueue <- event
	// 链接成功后，告知服务器自己的userID
	msg := Msg{
		msgType: MSG_TYPE_ACK,
		data: map[string]interface{}{
			"userID": c.userID,
		},
	}
	msgJson, _ := json.Marshal(msg)
	c.Conn.Write([]byte(msgJson))
	return nil
}

func (c *client) handleEvent() {
	for {
		select {
		case evt, ok := <-c.eventQueue:
			if !ok {
				return
			}
			switch evt.Type {
			case EVT_ON_CONNECT:
				if c.OnConnect != nil {
					c.OnConnect(evt)
				}
			case EVT_ON_DATA:
				if c.OnData != nil {
					c.OnData(evt)
				}
			case EVT_ON_DISCONNECT:
				if c.OnDisconnect != nil {
					c.OnDisconnect(evt)
				}
			}
		}
	}
}

func (c *client) handleClientConn() {
	buf := make([]byte, 65535)
	for {
		_, err := c.Conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("%+v", err)
			}
			eventQueue := ConnEvent{
				Type: EVT_ON_DISCONNECT,
			}
			c.eventQueue <- eventQueue
			break
		}
		eventQueue := ConnEvent{
			Type: EVT_ON_DATA,
			Data: buf,
		}
		c.eventQueue <- eventQueue
	}
}

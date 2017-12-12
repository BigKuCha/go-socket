package gosocket

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"errors"
)

/*连接状态*/
const (
	CONN_STATUS_NONE         = iota // 初始状态
	CONN_STATUS_CONNECTED           // 已连接,但未获得connID
	CONN_STATUS_READY               // 已获得connID,可以发送聊天消息
	CONN_STATUS_DISCONNECTED        // 失去连接
)

type client struct {
	NetWork
	userID       int
	connStatus   int
	eventQueue   chan ConnEvent
	OnConnect    func(event ConnEvent)
	OnData       func(msg ChatMsg)
	OnDisconnect func(event ConnEvent)
}

func NewClient(userID int) *client {
	return &client{
		connStatus: CONN_STATUS_NONE,
		userID:     userID,
		eventQueue: make(chan ConnEvent),
	}
}

func (c *client) Connect(host string, port int) error {
	netconn, err := net.Dial("tcp", net.JoinHostPort(host, strconv.Itoa(port)))
	if err != nil {
		return err
	}
	c.connStatus = CONN_STATUS_CONNECTED
	conn := Conn{
		conn: netconn,
	}
	c.conn = conn
	go c.handleEvent()
	go c.handleClientConn()
	event := ConnEvent{
		Type: EVT_ON_CONNECT,
	}
	c.eventQueue <- event
	// 链接成功后，告知服务器自己的userID
	msg := ChatMsg{
		MsgType: MSG_TYPE_ACK,
		Data:    []byte(strconv.Itoa(c.userID)),
	}
	c.SendMsg(msg)
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
					msg, err := handleMsg(evt.Data)
					if err != nil {
						fmt.Println("消息解析错误")
						continue
					}
					if msg.MsgType == MSG_TYPE_ACK {
						connID, _ := strconv.Atoi(string(msg.Data))
						c.conn.connID = uint32(connID)
						fmt.Println("我的链接ID是", connID)
						c.connStatus = CONN_STATUS_READY
						continue
					}
					c.OnData(msg)
				}
			case EVT_ON_DISCONNECT:
				if c.OnDisconnect != nil {
					c.connStatus = CONN_STATUS_DISCONNECTED
					c.OnDisconnect(evt)
				}
			}
		}
	}
}

func (c *client) handleClientConn() {
	buf := make([]byte, 65535)
	for {
		_, err := c.conn.conn.Read(buf)
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

// 发送消息
func (c *client) SendMsg(msg ChatMsg) (n int, err error) {
	msg.FromID = c.userID
	if msg.MsgType == MSG_TYPE_CHAT && c.connStatus != CONN_STATUS_READY {
		return 0, errors.New("未连接就绪")
	}
	msgBody := serialMsg(msg)
	n, err = c.conn.conn.Write(msgBody)
	return
}

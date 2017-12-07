package socket

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

const (
	EVT_ON_CONNECT = iota
	EVT_ON_DISCONNECT
	EVT_ON_DATA
	EVT_ON_CLOSE
)

type server struct {
	userConns map[int]int  // 用户ID和连接ID对应
	clients   map[int]Conn // 存放所有连接的客户端
	connID    int          // 为连接的客户端生成连接ID,自增

	addr         string
	port         int
	listener     net.Listener
	eventQueue   chan ConnEvent
	OnConnect    func(event ConnEvent)
	OnData       func(event ConnEvent)
	OnDisconnect func(event ConnEvent)
}

func NewServer(addr string, port int) *server {
	return &server{
		connID:     0,
		userConns:  make(map[int]int),
		clients:    make(map[int]Conn),
		addr:       addr,
		port:       port,
		eventQueue: make(chan ConnEvent, 10),
	}
}

func (s *server) Run() error {
	ls, err := net.Listen("tcp", net.JoinHostPort(s.addr, strconv.Itoa(s.port)))
	if err != nil {
		panic("listen 错误")
		return err
	}
	fmt.Printf("建立一个服务器，地址: %s \n", ls.Addr().String())
	s.listener = ls
	for {
		netConn, err := ls.Accept()
		if err != nil {
			panic("链接失败")
		}

		connID := s.connID + 1
		s.connID = connID
		conn := Conn{
			conn:   netConn,
			connID: connID,
			status: CONN_STATUS_CONNECTED,
		}
		s.clients[connID] = conn

		connEvent := ConnEvent{
			Type: EVT_ON_CONNECT,
			Conn: conn,
		}
		s.eventQueue <- connEvent
		go handleConn(s, conn)
		go s.handleEvent()
	}

	return nil
}

func (s *server) handleEvent() {
	for {
		select {
		case evt, ok := <-s.eventQueue:
			if !ok {
				return
			}
			switch evt.Type {
			case EVT_ON_CONNECT:
				if s.OnConnect != nil {
					s.OnConnect(evt)
				}
			case EVT_ON_DATA:
				if s.OnData != nil {
					msg, err := handleMsg(evt.Data)
					if err != nil {
						fmt.Println("消息类型错误")
						fmt.Printf("%+v", err)
						return
					}
					if msg.msgType == MSG_TYPE_ACK {
						connID := evt.Conn.connID
						fmt.Printf("链接ID是:%d\n", connID)
						if uid, ok := msg.data["userID"]; ok {
							userID, ok := uid.(int)
							if ok {
								s.userConns[connID] = userID
							}
						}
						continue
					}
					s.OnData(evt)
				}
			case EVT_ON_DISCONNECT:
				if s.OnDisconnect != nil {
					s.OnDisconnect(evt)
				}
			}
		}
	}
}

func handleConn(s *server, conn Conn) {
	buf := make([]byte, 65535)
	for {
		_, err := conn.conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("%+v", err)
				os.Exit(0)
			}
			fmt.Println(conn.conn.RemoteAddr(), "断开连接")
			eventQueue := ConnEvent{
				Type: EVT_ON_DISCONNECT,
			}
			s.eventQueue <- eventQueue
			break
		}
		eventQueue := ConnEvent{
			Type: EVT_ON_DATA,
			Data: buf,
		}
		s.eventQueue <- eventQueue
	}
}

func handleMsg(b []byte) (Msg, error) {
	var msg Msg
	err := json.Unmarshal(b, &msg)
	if err != nil {
		return msg, err
	}
	return msg, err
}

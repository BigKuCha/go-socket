package gosocket

import (
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
	conn      Conn
	userConns map[uint32]uint32 // 用户ID和连接ID对应
	clients   map[uint32]*Conn  // 存放所有连接的客户端
	connID    uint32            // 为连接的客户端生成连接ID,自增

	addr         string
	port         int
	listener     net.Listener
	eventQueue   chan ConnEvent
	OnConnect    func(event ConnEvent)
	OnData       func(msg ChatMsg)
	OnDisconnect func(event ConnEvent)
}

func NewServer(addr string, port int) *server {
	return &server{
		connID:     0,
		userConns:  make(map[uint32]uint32),
		clients:    make(map[uint32]*Conn),
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
			conn:       netConn,
			connID:     connID,
			localAddr:  netConn.LocalAddr().String(),
			remoteAddr: netConn.RemoteAddr().String(),
		}
		s.clients[connID] = &conn
		// 通知客户端其连接ID
		msg := ChatMsg{
			MsgType: MSG_TYPE_ACK,
			Data:    []byte(strconv.Itoa(int(connID))),
		}
		conn.SendMsg(msg)
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
					msg, err := HandleMsg(evt.Data)
					if err != nil {
						fmt.Println("消息类型错误")
						fmt.Printf("%+v", err)
						return
					}
					if msg.MsgType == MSG_TYPE_ACK {
						userID, _ := strconv.Atoi(string(msg.Data))
						connID := evt.Conn.connID
						s.userConns[uint32(userID)] = connID
						continue
					} else if msg.MsgType == MSG_TYPE_CHAT {
						if connID, ok := s.userConns[uint32(msg.ToID)]; ok {
							if toConn, ok := s.clients[connID]; ok {
								toConn.SendMsg(msg)
							} else {
								fmt.Println("对方未连接")
							}
						} else {
							fmt.Println("对方未连接")
						}
					}
					s.OnData(msg)
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
			eventQueue := ConnEvent{
				Conn: conn,
				Type: EVT_ON_DISCONNECT,
			}
			s.eventQueue <- eventQueue
			break
		}
		eventQueue := ConnEvent{
			Conn: conn,
			Type: EVT_ON_DATA,
			Data: buf,
		}
		s.eventQueue <- eventQueue
	}
}

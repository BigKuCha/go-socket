package socket

import (
	"net"
	"io"
	"fmt"
)

/*连接状态*/
const (
	SERVER_STATUS_NONE         = iota
	SERVER_STATUS_CONNECTED
	SERVER_STATUS_DISCONNECTED
)

const (
	EVT_ON_CONNECT    = iota
	EVT_ON_DISCONNECT
	EVT_ON_DATA
	EVT_ON_CLOSE
)

type server struct {
	addr         string
	port         int
	listener     net.Listener
	eventQueue   chan ConnEvent
	OnConnect    func()
	OnData       func()
	OnDisconnect func()
}

func NewServer(addr string, port int) *server {
	return &server{
		addr:       addr,
		port:       port,
		eventQueue: make(chan ConnEvent, 10),
	}
}

func (s *server) Run() error {
	ls, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		panic("listen 错误")
		return err
	}
	fmt.Println("建立一个服务器，监听端口 8888...")
	s.listener = ls
	for {
		conn, err := ls.Accept()
		if err != nil {
			panic("链接失败")
		}
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
				fmt.Println("获得一个连接")
				s.OnConnect()
			case EVT_ON_DATA:
				fmt.Println("收到一个消息")
				s.OnData()
			case EVT_ON_DISCONNECT:
				fmt.Println("断开链接")
				s.OnDisconnect()
			}
		}
	}
}

func handleConn(s *server, conn net.Conn) {
	buf := make([]byte, 65535)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
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

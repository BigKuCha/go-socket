package socket

import (
	"fmt"
	"io"
	"net"
	"strconv"
	"os"
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
	OnConnect    func(event ConnEvent)
	OnData       func(event ConnEvent)
	OnDisconnect func(event ConnEvent)
}

func NewServer(addr string, port int) *server {
	return &server{
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
	fmt.Printf("建立一个服务器，监听端口 %d...\n", s.port)
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
				if s.OnConnect != nil {
					s.OnConnect(evt)
				}
			case EVT_ON_DATA:
				if s.OnData != nil {
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

func handleConn(s *server, conn net.Conn) {
	buf := make([]byte, 65535)
	for {
		_, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				fmt.Printf("%+v", err)
				os.Exit(0)
			}
			fmt.Println(conn.RemoteAddr(), "断开连接")
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

package socket

import (
	"net"
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
	onConnect    func()
	onData       func()
	onDisconnect func()
}

func NewServer(addr string, port int) *server {
	return &server{
		addr: addr,
		port: port,
	}
}

func (s *server) Run() error {
	ls, err := net.Listen("tcp", s.addr+":"+string(s.port))
	if err != nil {
		return err
	}
	s.listener = ls
	return nil
}

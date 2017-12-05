package socket

import "net"

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
type Conn struct {
	conn       net.Conn
	connID     int
	status     int32
	userData   interface{}
	localAddr  string
	remoteAddr string
}

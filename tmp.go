package main

import (
	"net"
	"fmt"
	"github.com/urfave/cli"
	"io"
)

func runServer(context *cli.Context) error {
	ln, err := net.Listen("tcp", "localhost:8888")
	if err != nil {
		panic("服务器声明错误")
	}
	fmt.Println("声明了一个服务器， 监听端口8888")
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			panic("链接接受错误")
		}
		fmt.Println("获得了一个连接", conn.RemoteAddr())
		go handleConn(conn)
	}
	return nil
}

func handleConn(conn net.Conn) {
	buf := make([]byte, 0, 4096) // big buffer
	tmp := make([]byte, 256)     // using small tmo buffer for demonstrating
	for {
		n, err := conn.Read(tmp)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error:", err)
			}
			break
		}
		//fmt.Println("got", n, "bytes.")
		buf = append(buf, tmp[:n]...)

	}
	fmt.Println(string(buf))

	//ticker := time.NewTicker(6 * time.Second)
	//for {
	//	select {
	//	case <-ticker.C:
	//		fmt.Println("检测是否断开了连接====")
	//		_, err := conn.Read(make([]byte, 0))
	//		if err != nil && err != io.EOF {
	//			fmt.Println("断开了一个连接", conn.RemoteAddr())
	//			conn.Close()
	//			ticker.Stop()
	//			break
	//		} else {
	//
	//		}
	//	}
	//}
}

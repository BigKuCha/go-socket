package main

import (
	"github.com/urfave/cli"
	"os"
	"go-socket/socket"
	"fmt"
	"net"
	"bufio"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name: "server",
			Action: func(*cli.Context) {
				server := socket.NewServer("localhost", 8888)
				server.OnConnect = onConnect
				server.OnData = onData
				server.Run()
			},
		},
		{
			Name:   "client",
			Action: runClient,
		},
	}
	app.Run(os.Args)
}

func onConnect(event socket.ConnEvent) {
	fmt.Println("我收到了一个连接")
}

func onData(event socket.ConnEvent) {
	fmt.Printf("%+v", event.Data)
	//fmt.Println("我收到了一条数据")
}

func runClient(context *cli.Context) error {
	conn, err := net.Dial("tcp", "localhost:8888")
	if err != nil {
		panic(err)
		// handle error
	}
	go func() {
		status, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println(status)
	}()
	fmt.Println("----")
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("enter msg>")
		b, _, _ := r.ReadLine()
		conn.Write(b)
	}
	select {}
	return nil
}

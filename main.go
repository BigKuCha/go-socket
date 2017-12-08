package main

import (
	"bufio"
	"fmt"
	"github.com/urfave/cli"
	"go-socket/socket"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name: "server",
			Action: func(context *cli.Context) {
				server := socket.NewServer("localhost", 8888)
				server.OnConnect = onServerConnect
				server.OnData = onServerData
				server.Run()
			},
		},
		{
			Name: "client",
			Action: func(context *cli.Context) {
				client := socket.NewClient(10244)
				client.OnConnect = onClientConnect
				client.OnData = onClientData
				client.OnDisconnect = onClientDisconnect
				err := client.Connect("localhost", 8888)
				if err != nil {
					fmt.Println("客户端连接失败", err)
					os.Exit(0)
				}
				r := bufio.NewReader(os.Stdin)
				for {
					fmt.Print("enter msg>")
					b, _, _ := r.ReadLine()
					client.Conn.Write(b)
				}
			},
		},
	}
	app.Run(os.Args)
}

func onServerConnect(event socket.ConnEvent) {
	fmt.Println("我收到了一个连接")
}

func onServerData(event socket.ConnEvent) {
	fmt.Println("我收到了一个消息")
	//fmt.Printf("%+v", event.Data)
}

func onClientConnect(event socket.ConnEvent) {
	fmt.Println("客户端连接成功了！")
}

func onClientDisconnect(event socket.ConnEvent) {
	fmt.Println("客户端已经断开链接")
	os.Exit(0)
}

func onClientData(event socket.ConnEvent) {
	fmt.Println("客户端收到一条消息")
}

package main

import (
	"bufio"
	"fmt"
	"github.com/bigkucha/go-socket"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name: "server",
			Action: func(context *cli.Context) {
				server := gosocket.NewServer("localhost", 8888)
				server.OnConnect = onServerConnect
				server.OnData = onServerData
				server.Run()
			},
		},
		{
			Name: "client",
			Action: func(context *cli.Context) {
				client := gosocket.NewClient(1)
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
					msg := gosocket.ChatMsg{
						ToID:    2,
						MsgType: gosocket.MSG_TYPE_CHAT,
						Data:    b,
					}
					client.SendMsg(msg)
				}
			},
		},
	}
	app.Run(os.Args)
}

func onServerConnect(event gosocket.ConnEvent) {
	fmt.Println("我收到了一个连接")
}

func onServerData(msg gosocket.ChatMsg) {
	fmt.Println("客户端消息:", string(msg.Data))
	//if data, ok := msg.Data["data"]; ok {
	//	fmt.Println("客户端消息:", string(data))
	//}
}

func onClientConnect(event gosocket.ConnEvent) {
	fmt.Println("客户端连接成功了！")
}

func onClientDisconnect(event gosocket.ConnEvent) {
	fmt.Println("客户端已经断开链接")
	os.Exit(0)
}

func onClientData(msg gosocket.ChatMsg) {
	fmt.Println("别人对我说:", string(msg.Data))
}

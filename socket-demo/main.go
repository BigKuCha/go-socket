package main

import (
	"bufio"
	"fmt"
	"github.com/bigkucha/go-socket"
	"github.com/urfave/cli"
	"os"
	"strconv"
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
				server.OnDisconnect = onServerDisconnect
				server.Run()
			},
		},
		{
			Name: "client",
			Action: func(context *cli.Context) {
				if context.NArg() < 2 {
					fmt.Println("参数不够！ 第一个参数为用户ID，第二个参数为想要聊天的用户ID")
					return
				}

				userID, err := strconv.Atoi(context.Args()[0])
				if err != nil {
					fmt.Println("第一个参数字为整形数字")
					return
				}
				toID, err := strconv.Atoi(context.Args()[1])

				if err != nil {
					fmt.Println("第二个参数为整形数字！")
					return
				}

				client := gosocket.NewClient(userID)
				client.OnConnect = onClientConnect
				client.OnData = onClientData
				client.OnDisconnect = onClientDisconnect
				err = client.Connect("localhost", 8888)
				if err != nil {
					fmt.Println("客户端连接失败", err)
					os.Exit(0)
				}
				r := bufio.NewReader(os.Stdin)
				for {
					fmt.Print("enter msg>")
					b, _, _ := r.ReadLine()
					msg := gosocket.ChatMsg{
						ToID:    toID,
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
	fmt.Printf("客户端消息: %d 对 %d 说:%s \n", msg.FromID, msg.ToID, string(msg.Data))
}

func onServerDisconnect(event gosocket.ConnEvent) {
	fmt.Println(event.Conn.GetRemoteAddr(), "断开连接")
}

func onClientConnect(event gosocket.ConnEvent) {
	fmt.Println("客户端连接成功了！")
}

func onClientDisconnect(event gosocket.ConnEvent) {
	fmt.Println("客户端已经断开链接")
	os.Exit(0)
}

func onClientData(msg gosocket.ChatMsg) {
	fmt.Printf("%d 对我说 :%s \n", msg.FromID, string(msg.Data))
}

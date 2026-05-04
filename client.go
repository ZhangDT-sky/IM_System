package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
)

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	flag       int
}

func NewClient(serverIp string, serverPort int) *Client {
	//创建客户端对象
	client := &Client{
		ServerIp:   serverIp,
		ServerPort: serverPort,
		flag:       999,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort))
	if err != nil {
		fmt.Println("net.Dial err:", err)
		return nil
	}
	client.conn = conn
	return client
}

func (client *Client) DealResponse() {
	io.Copy(os.Stdout, client.conn)
}

func (client *Client) menu() bool {
	var flag int
	fmt.Println("1.公聊模式")
	fmt.Println("2.私聊模式")
	fmt.Println("3.更新用户名")
	fmt.Println("0.退出")
	fmt.Scanln(&flag)
	if flag >= 0 && flag <= 3 {
		client.flag = flag
		return true
	} else {
		fmt.Println(">>>>>请输入合法范维数字>>>>>")
		return false
	}
}

func (client *Client) Run() {
	for client.flag != 0 {
		for client.menu() != true {

		}
		// 根据flag处理业务
		switch client.flag {
		case 1:
			//
			break
		case 2:
			//
			break
		case 3:
			break
		}
	}
}

var serverIp string
var serverPort int

func init() {
	flag.StringVar(&serverIp, "serverIp", "127.0.0.1", "server ip (127.0.0.1) )")
	flag.IntVar(&serverPort, "serverPort", 8888, "server port (8888)")
}

func main() {
	flag.Parse()

	client := NewClient(serverIp, serverPort)
	if client == nil {
		fmt.Println(">>>>>连接服务器失败>>>>>")
		return
	}
	fmt.Println(">>>>>连接服务器成功>>>>>")

	client.Run()
}

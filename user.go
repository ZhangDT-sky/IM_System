package main

import (
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

func NewUser(conn net.Conn) *User {
	userAdder := conn.RemoteAddr().String()
	user := &User{
		Name: userAdder,
		Addr: userAdder,
		C:    make(chan string),
		conn: conn,
	}
	go user.ListenMessage()
	return user
}

// 监听user channel 方法，有消息发送客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\r\n"))
	}
}

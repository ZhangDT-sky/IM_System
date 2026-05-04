package main

import (
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn

	server *Server
}

func NewUser(conn net.Conn, server *Server) *User {
	userAdder := conn.RemoteAddr().String()
	user := &User{
		Name:   userAdder,
		Addr:   userAdder,
		C:      make(chan string),
		conn:   conn,
		server: server,
	}
	go user.ListenMessage()
	return user
}

// 用户上线
func (this *User) Online() {
	// 将用户加入列表
	this.server.mapLock.Lock()
	this.server.OnlineMap[this.Name] = this
	this.server.mapLock.Unlock()
	//广播
	this.server.BroadCast(this, "Online")
}

// 用户下线
func (this *User) Offline() {
	this.server.mapLock.Lock()
	delete(this.server.OnlineMap, this.Name)
	this.server.mapLock.Unlock()
	this.server.BroadCast(this, "Offline")
}

// 处理消息
func (this *User) DoMessage(msg string) {
	this.server.BroadCast(this, msg)
}

// 监听user channel 方法，有消息发送客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\r\n"))
	}
}

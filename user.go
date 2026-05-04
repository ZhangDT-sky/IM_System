package main

import (
	"net"
	"strings"
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

func (this *User) SendMsg(msg string) {
	this.conn.Write([]byte(msg))
}

// 处理消息
func (this *User) DoMessage(msg string) {
	if msg == "who" {
		this.server.mapLock.Lock()
		for _, user := range this.server.OnlineMap {
			onlineMsg := "[" + user.Addr + "]" + user.Name + ":" + "Ready...\n"
			this.SendMsg(onlineMsg)
		}
		this.server.mapLock.Unlock()
	} else if strings.HasPrefix(msg, "rename|") {
		newName := strings.TrimSpace(strings.TrimPrefix(msg, "rename|"))
		if newName == "" {
			this.SendMsg("username cannot be empty\r\n")
			return
		}

		this.server.mapLock.Lock()
		_, ok := this.server.OnlineMap[newName]
		if ok {
			this.server.mapLock.Unlock()
			this.SendMsg("username already use\r\n")
			return
		}

		delete(this.server.OnlineMap, this.Name)
		this.Name = newName
		this.server.OnlineMap[newName] = this
		this.server.mapLock.Unlock()

		this.SendMsg("username already change: " + newName + "\r\n")
	} else if strings.HasPrefix(msg, "to|") {
		parts := strings.SplitN(msg, "|", 3)
		if len(parts) != 3 {
			this.SendMsg("message format error, use: to|name|content\r\n")
			return
		}

		remoteName := strings.TrimSpace(parts[1])
		content := strings.TrimSpace(parts[2])

		if remoteName == "" {
			this.SendMsg("username cannot be empty\r\n")
			return
		}
		if content == "" {
			this.SendMsg("message cannot be empty\r\n")
			return
		}

		this.server.mapLock.RLock()
		remoteUser, ok := this.server.OnlineMap[remoteName]
		this.server.mapLock.RUnlock()

		if !ok {
			this.SendMsg("username not found\r\n")
			return
		}
		remoteUser.SendMsg(this.Name + " speak: " + content + "\r\n")
	} else {
		this.server.BroadCast(this, msg)
	}
}

// 监听user channel 方法，有消息发送客户端
func (this *User) ListenMessage() {
	for {
		msg := <-this.C
		this.conn.Write([]byte(msg + "\r\n"))
	}
}

package main

import (
	"fmt"
	"net"
	"strings"
)

type User struct {
	Name       string
	Addr       string
	C          chan string
	connection net.Conn
	server     *Server
}

func NewUser(con net.Conn, server *Server) *User {
	name := con.RemoteAddr().String()
	name = strings.Split(name, ":")[1]
	user := &User{
		Name:       name,
		Addr:       name,
		C:          make(chan string),
		connection: con,
		server:     server,
	}

	go user.ListenC()

	return user
}

func (u *User) ListenC() {
	for {
		msg := <-u.C
		fmt.Println("发送消息", msg, "给客户端")
		u.connection.Write([]byte(msg))
	}
}

func (u *User) Online() {
	// 加入server的map中
	u.server.lock.Lock()
	u.server.OnlineMap[u.Name] = u
	u.server.lock.Unlock()
	u.server.Broadcast(u, "上线")
}

func (u *User) OffLine() {
	// 从server的map中移除
	u.server.lock.Lock()
	delete(u.server.OnlineMap, u.Name)
	u.server.lock.Unlock()
	u.server.Broadcast(u, "下线")
}

func (u *User) SendMessage(msg string) {
	// 发送消息
	u.server.Broadcast(u, msg)
}

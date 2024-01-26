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
}

func NewUser(con net.Conn) *User {
	name := con.RemoteAddr().String()
	name = strings.Split(name, ":")[1]
	user := &User{
		Name:       name,
		Addr:       name,
		C:          make(chan string),
		connection: con,
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

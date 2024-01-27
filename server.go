package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

type Server struct {
	IP        string
	Port      int
	Message   chan string
	OnlineMap map[string]*User
	lock      sync.RWMutex
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		IP:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}

	return server
}

func (s *Server) Start() {
	// 创建套接字
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		err.Error()
	}
	defer listener.Close()
	go s.ListenMessage()
	// 建立链接
	for {
		// 循环监听链接建了
		conn, err := listener.Accept()
		if err != nil {
			err.Error()
		}
		// 有一个客户端请求建立了链接
		go s.Handler(conn)
	}
}

func (s *Server) Broadcast(user *User, msg string) {
	message := fmt.Sprintf("%s:%s\n", user.Name, msg)
	s.Message <- message
}

func (s *Server) Handler(conn net.Conn) {
	// 把conn写入onlineMap
	user := NewUser(conn, s)
	// 用户上线
	user.Online()
	go func() {
		buf := make([]byte, 10000)
		for {
			// 收到来自客户端的消息
			n, err := conn.Read(buf)
			if n == 0 {
				user.OffLine()
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println(err)
				return
			}
			// 广播给所有在线的用户
			msg := string(buf[:n-1])
			user.SendMessage(msg)
		}
	}()

	select {}
}

func (s *Server) ListenMessage() {
	for {
		msg := <-s.Message
		s.lock.Lock()
		for _, cli := range s.OnlineMap {
			cli.C <- msg
		}
		s.lock.Unlock()
	}
}

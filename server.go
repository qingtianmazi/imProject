package main

import (
	"fmt"
	"net"
)

type Server struct {
	IP   string
	Port int
}

func NewServer(ip string, port int) *Server {
	return &Server{
		IP:   ip,
		Port: port,
	}
}

func (s *Server) Start() {
	// 创建套接字
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		err.Error()
	}
	defer listener.Close()
	// 建立链接
	for {
		// 循环监听链接建了
		conn, err := listener.Accept()
		if err != nil {
			err.Error()
		}
		// 处理方法
		go s.Handler(conn)
	}
}

func (s *Server) Handler(conn net.Conn) {
	fmt.Println("连接建立成功...")
}

package main

import (
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*User
	mapLock   sync.RWMutex
	Message   chan string
}

func NewServer(ip string, port int) *Server {
	return &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*User),
		Message:   make(chan string),
	}
}

func (this *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip, this.Port))
	if err != nil {
		fmt.Println("start err:", err)
	}
	defer listener.Close()
	go this.ListenMessage()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(" listener Accept err:", err)
			continue
		}
		this.Handler(conn)
	}

}

func (this *Server) Handler(conn net.Conn) {
	fmt.Println(conn, "连接成功")
	user := NewUser(conn)
	//	当前用户上线了
	this.mapLock.Lock()
	this.OnlineMap[user.Name] = user
	this.mapLock.Unlock()
	this.BroadCast(user, "已上线")
	select {}
}

func (this *Server) BroadCast(user *User, msg string) {
	sendMsg := fmt.Sprintf("[%s] %s: %s", user.Addr, user.Name, msg)
	this.Message <- sendMsg
}
func (this *Server) ListenMessage() {
	for {
		msg := <-this.Message
		this.mapLock.Lock()
		for _, cli := range this.OnlineMap {
			cli.C <- msg
		}
		this.mapLock.Unlock()

	}
}

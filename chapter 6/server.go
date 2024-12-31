package main

import (
	"fmt"
	"net"
)

type Server struct {
	Ip string
	Port int
}

func NewServer(ip string, port int) *Server{
	server :=&Server{
		Ip: ip,
		Port: port,
	}
	return server
}

func (this *Server) Handler(conn net.Conn){
	fmt.Println("连接建立成功")
}

func (this *Server) start(){
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", this.Ip,this.Port))
	if err != nil{
		fmt.Println("net.Listen err:", err)

	}
	defer listener.Close()
	for {
	//accept
	conn, err := listener.Accept()
	if err != nil{
		fmt.Println("listener accept err:", err)
		continue
	}


	//do handler
	go this.Handler(conn)
	}


	//close listen socket



}
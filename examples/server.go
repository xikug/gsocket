package main

import (
	"bufio"
	"log"
	"os"

	"github.com/xikug/gsocket"
)

type demoServer struct{}

// OnConnect 客户端连接事件
func (server demoServer) OnConnect(c *gsocket.Connection) {
	log.Printf("CONNECTED: %s\n", c.RemoteAddr())
}

// OnDisconnect 客户端断开连接事件
func (server demoServer) OnDisconnect(c *gsocket.Connection) {
	log.Printf("DISCONNECTED: %s\n", c.RemoteAddr())
}

// OnRecv 收到客户端发来的数据
func (server demoServer) OnRecv(c *gsocket.Connection, data []byte) {
	log.Printf("DATA RECVED: %s %d - %v\n", c.RemoteAddr(), len(data), data)
	c.Send(data)
}

// OnError 有错误发生
func (server demoServer) OnError(c *gsocket.Connection, err error) {
	log.Printf("ERROR: %s - %s\n", c.RemoteAddr(), err.Error())
}

func main() {
	demoServer := &demoServer{}
	//CreateTCPServer 的handler可以传nil
	server := gsocket.CreateTCPServer("0.0.0.0", 9595,
		demoServer.OnConnect, demoServer.OnDisconnect, demoServer.OnRecv, demoServer.OnError)

	err := server.Start()
	if err != nil {
		log.Printf("Start Server Error: %s\n", err.Error())
		return
	}

	log.Printf("Listening %s...\n", server.Addr())

	pause()
}

func pause() {
	println("按回车键退出...\n")
	r := bufio.NewReader(os.Stdin)
	r.ReadByte()
}

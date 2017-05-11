package main

import (
	"bufio"
	"log"
	"os"

	"github.com/xikug/gsocket"
)

type demoServer struct{}

// OnConnect 客户端连接事件
func (server demoServer) OnConnect(session *gsocket.Session) {
	log.Printf("CONNECTED: %s\n", session.RemoteAddr())
}

// OnDisconnect 客户端断开连接事件
func (server demoServer) OnDisconnect(session *gsocket.Session) {
	log.Printf("DISCONNECTED: %s\n", session.RemoteAddr())
}

// OnRecv 收到客户端发来的数据
func (server demoServer) OnRecv(session *gsocket.Session, data []byte) {
	log.Printf("DATA RECVED: %s %d - %v\n", session.RemoteAddr(), len(data), data)
}

// OnError 有错误发生
func (server demoServer) OnError(session *gsocket.Session, err error) {
	log.Printf("ERROR: %s - %s\n", session.RemoteAddr(), err.Error())
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

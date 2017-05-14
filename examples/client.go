package main

import (
	"bufio"
	"log"
	"os"

	"github.com/xikug/gsocket"
)

type demoClient struct{}

func (client *demoClient) OnConnect(c *gsocket.Connection) {
	log.Printf("CONNECTED: %s\n", c.RemoteAddr())
}

func (client *demoClient) OnDisconnect(c *gsocket.Connection) {
	log.Printf("DISCONNECTED: %s\n", c.RemoteAddr())
}

func (client *demoClient) OnRecv(c *gsocket.Connection, data []byte) {
	log.Printf("DATA RECVED: %s %d - %v\n", c.RemoteAddr(), len(data), data)
}

func (client *demoClient) OnError(c *gsocket.Connection, err error) {
	log.Printf("ERROR: %s - %s\n", c.RemoteAddr(), err.Error())
}

func main() {
	demoClient := &demoClient{}

	client := gsocket.CreateTCPClient(demoClient.OnConnect, demoClient.OnDisconnect, demoClient.OnRecv, demoClient.OnError)

	err := client.Connect("127.0.0.1", 9595)
	if err != nil {
		log.Printf("Coneect Server Error: %s\n", err.Error())
		return
	}

	log.Printf("Connect Server %s Success\n", client.RemoteAddr())

	client.Send([]byte("Hello World!!!"))

	pause()
	client.Close()
}

func pause() {
	println("按回车键退出...\n")
	r := bufio.NewReader(os.Stdin)
	r.ReadByte()
}

package gsocket

import (
	"fmt"
	"net"
	"sync"
)

// TCPClient TCP客户端描述
type TCPClient struct {
	tcpClientState
	c           *Connection
	userHandler tcpEventHandler
	wg          sync.WaitGroup
}

type tcpClientState struct {
	remoteAddr string
	remotePort uint16
	connected  bool
}

// CreateTCPClient 创建一个TCPClient实例
func CreateTCPClient(handlerConnect TCPConnectHandler, handlerDisconnect TCPDisconnectHandler,
	handlerRecv TCPRecvHandler, handlerError TCPErrorHandler) *TCPClient {
	client := &TCPClient{
		userHandler: tcpEventHandler{
			handlerConnect:    handlerConnect,
			handlerDisconnect: handlerDisconnect,
			handlerRecv:       handlerRecv,
			handlerError:      handlerError,
		},
	}

	return client
}

// Connect 连接到服务器
func (client *TCPClient) Connect(addr string, port uint16) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", addr, port))
	if err != nil {
		return err
	}
	client.c = newConnection(conn)

	client.tcpClientState = tcpClientState{
		remoteAddr: addr,
		remotePort: port,
		connected:  true,
	}

	client.wg.Add(2)
	go client.c.recvThread(&client.wg, client.userHandler)
	go client.c.sendThread(&client.wg)
	return nil
}

// Send 发送数据
func (client *TCPClient) Send(data []byte) {
	client.c.Send(data)
}

// Close 关闭连接
func (client *TCPClient) Close() {
	client.c.Close()
	client.wg.Wait()
	if client.userHandler.handlerDisconnect != nil {
		client.userHandler.handlerDisconnect(client.c)
	}
}

// RemoteAddr 返回服务器地址
func (client *TCPClient) RemoteAddr() string {
	return fmt.Sprintf("%s:%d", client.tcpClientState.remoteAddr, client.tcpClientState.remotePort)
}

// LocalAddr 返回本机的连接地址
func (client *TCPClient) LocalAddr() string {
	return client.c.conn.LocalAddr().String()
}

package gsocket

import (
	"fmt"
	"net"
	"sync"
)

// TCPServer 描述一个TCP服务器的结构
type TCPServer struct {
	tcpServerState
	userHandler   eventHandler   // 用户的事件处理Handler
	connectionMax int            // 最大连接数，为0则不限制服务器最大连接数
	listener      net.Listener   // 监听句柄
	terminated    bool           // 通知是否停止Service
	wg            sync.WaitGroup // 等待所有goroutine结束
}

type eventHandler struct {
	handlerConnect    TCPConnectHandler
	handlerDisconnect TCPDisconnectHandler
	handlerRecv       TCPRecvHandler
	handlerError      TCPErrorHandler
}

type tcpServerState struct {
	listenAddr      string // 监听地址
	listenPort      uint16 // 监听端口
	connectionCount uint32 // 当前连接数
}

// CreateTCPServer 创建一个TCPServer, 返回*TCPServer
func CreateTCPServer(addr string, port uint16, handlerConnect TCPConnectHandler, handlerDisconnect TCPDisconnectHandler,
	handlerRecv TCPRecvHandler, handlerError TCPErrorHandler) *TCPServer {
	server := &TCPServer{
		tcpServerState: tcpServerState{
			listenAddr:      addr,
			listenPort:      port,
			connectionCount: 0,
		},
		userHandler: eventHandler{
			handlerConnect:    handlerConnect,
			handlerDisconnect: handlerDisconnect,
			handlerRecv:       handlerRecv,
			handlerError:      handlerError,
		},
	}

	return server
}

// Start 开始服务
func (server TCPServer) Start() error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.listenAddr, server.listenPort))
	if err != nil {
		return err
	}

	server.wg.Add(1)

	go func() {
		for {
			if server.terminated {
				server.wg.Done()
				break
			}

			conn, err := listener.Accept()
			if err != nil {
				server.ProcessError(nil, err)
				continue
			}
		}
	}()

	return nil
}

// Stop 停止服务
func (server TCPServer) Stop() {
	server.terminated = true
	server.wg.Wait() // 等待结束
}

// ConnectionCount 返回服务器当前连接数
func (server TCPServer) ConnectionCount() uint32 {
	return server.tcpServerState.connectionCount
}

// SetMaxConnection 设置服务器最大连接数
func (server TCPServer) SetMaxConnection(maxCount int) {
	server.connectionMax = maxCount
}

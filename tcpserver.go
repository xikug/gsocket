package gsocket

import (
	"fmt"
	"log"
	"net"
	"sync"
)

// TCPServer 描述一个TCP服务器的结构
type TCPServer struct {
	tcpServerState
	userHandler   eventHandler       // 用户的事件处理Handler
	connectionMax int                // 最大连接数，为0则不限制服务器最大连接数
	listener      net.Listener       // 监听句柄
	terminated    bool               // 通知是否停止Service
	wg            sync.WaitGroup     // 等待所有goroutine结束
	sessions      map[uint64]Session // 会话
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
func (server *TCPServer) Start() error {
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
				server.processError(nil, err)
				continue
			}

			go func(conn net.Conn) {
				session := server.makeSession(conn)

				server.processConnect(session)
			}(conn)
		}
	}()

	return nil
}

// Stop 停止服务
func (server *TCPServer) Stop() {
	server.terminated = true
	server.wg.Wait() // 等待结束
}

// ConnectionCount 返回服务器当前连接数
func (server *TCPServer) ConnectionCount() uint32 {
	return server.tcpServerState.connectionCount
}

// SetMaxConnection 设置服务器最大连接数
func (server *TCPServer) SetMaxConnection(maxCount int) {
	server.connectionMax = maxCount
}

// Addr 返回服务器监听的地址
func (server *TCPServer) Addr() string {
	return fmt.Sprintf("%s:%d", server.listenAddr, server.listenPort)
}

func (server *TCPServer) makeSession(conn net.Conn) (session *Session) {
	session = newSession(conn)

	server.wg.Add(2)
	go session.recvThread(server)
	go session.sendThread(server)

	return session
}

func (server *TCPServer) processConnect(session *Session) {
	log.Printf("ACCEPTED: %s\n", session.RemoteAddr())
	if server.userHandler.handlerConnect != nil {
		server.userHandler.handlerConnect.OnConnect(session)
	}
}

func (server *TCPServer) processDisconnect(session *Session) {
	log.Printf("CONNECTION CLOSED: %s\n", session.RemoteAddr())
	if server.userHandler.handlerDisconnect != nil {
		server.userHandler.handlerDisconnect.OnDisconnect(session)
	}
}

func (server *TCPServer) processRecv(session *Session, data []byte) {
	log.Printf("DATA RECVED: %x\n", data)
	if server.userHandler.handlerRecv != nil {
		server.userHandler.handlerRecv.OnRecv(session, data)
	}
}

func (server *TCPServer) processError(session *Session, err error) {
	log.Printf("ERROR: %s\n", err.Error())
	if server.userHandler.handlerError != nil {
		server.userHandler.handlerError.OnError(session, err)
	}
}

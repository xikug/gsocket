package gsocket

// TCPServer 描述一个TCP服务器的结构
type TCPServer struct {
	tcpServerState
	userHandler   eventHandler // 用户的事件处理Handler
	connectionMax int          // 最大连接数，为0则不限制服务器最大连接数
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
func (server TCPServer) Start() {

}

// Stop 停止服务
func (server TCPServer) Stop() {

}

// ConnectionCount 返回服务器当前连接数
func (server TCPServer) ConnectionCount() uint32 {
	return server.tcpServerState.connectionCount
}

// SetMaxConnection 设置服务器最大连接数
func (server TCPServer) SetMaxConnection(maxCount int) {
	server.connectionMax = maxCount
}

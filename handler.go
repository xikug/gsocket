package gsocket

// TCPConnectHandler 连接事件
type TCPConnectHandler func(session *Connection)

// TCPDisconnectHandler 断开连接事件
type TCPDisconnectHandler func(session *Connection)

// TCPRecvHandler 收到数据事件
type TCPRecvHandler func(session *Connection, data []byte)

// TCPErrorHandler 有错误发生
type TCPErrorHandler func(session *Connection, err error)

type tcpEventHandler struct {
	handlerConnect    TCPConnectHandler
	handlerDisconnect TCPDisconnectHandler
	handlerRecv       TCPRecvHandler
	handlerError      TCPErrorHandler
}

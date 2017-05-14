package gsocket

// TCPConnectHandler 连接事件
type TCPConnectHandler func(c *Connection)

// TCPDisconnectHandler 断开连接事件
type TCPDisconnectHandler func(c *Connection)

// TCPRecvHandler 收到数据事件
type TCPRecvHandler func(c *Connection, data []byte)

// TCPErrorHandler 有错误发生
type TCPErrorHandler func(c *Connection, err error)

type tcpEventHandler struct {
	handlerConnect    TCPConnectHandler
	handlerDisconnect TCPDisconnectHandler
	handlerRecv       TCPRecvHandler
	handlerError      TCPErrorHandler
}

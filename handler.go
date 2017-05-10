package gsocket

// TCPConnectHandler 连接事件
type TCPConnectHandler func(session *Session)

// TCPDisconnectHandler 断开连接事件
type TCPDisconnectHandler func(session *Session)

// TCPRecvHandler 收到数据事件
type TCPRecvHandler func(session *Session, data []byte)

// TCPErrorHandler 有错误发生
type TCPErrorHandler func(session *Session, err error)

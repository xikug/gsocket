package gsocket

type TCPConnectHandler interface {
	OnConnect(session *Session)
}

type TCPDisconnectHandler interface {
	OnDisconnect(session *Session)
}

type TCPRecvHandler interface {
	OnRecv(session *Session, data []byte)
}

type TCPErrorHandler interface {
	OnError(session *Session, err error)
}

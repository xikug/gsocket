package gsocket

type TCPConnectHandler interface {
	OnConnect()
}

type TCPDisconnectHandler interface {
	OnDisconnect()
}

type TCPRecvHandler interface {
	OnRecv()
}

type TCPErrorHandler interface {
	OnError()
}

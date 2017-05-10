package gsocket

import (
	"log"
	"net"
)

// Session 代表一个连接会话
type Session struct {
	connection net.Conn
	sendBuffer chan []byte
	terminated bool
}

// NewSession 生成一个新的Session
func newSession(conn net.Conn) (session *Session) {
	session = &Session{
		connection: conn,
		sendBuffer: make(chan []byte, 10),
		terminated: false,
	}

	return session
}

// RemoteAddr 返回客户端的地址和端口
func (session *Session) RemoteAddr() string {
	return session.connection.RemoteAddr().String()
}

// Close 关闭Session
func (session *Session) Close() {
	close(session.sendBuffer)
	session.connection.Close()
}

func (session *Session) recvThread(server *TCPServer) {
	defer server.wg.Done()
	buffer := make([]byte, 4096)
	for {
		n, err := session.connection.Read(buffer)
		if err != nil {
			break
		}

		//session.RecvedPackets = append(session.RecvedPackets, buffer[:n]...)
		if server.userHandler.handlerRecv != nil {
			server.userHandler.OnRecv(session, buffer[:n])
		}
	}

	log.Printf("session %s recvThread Exit", session.RemoteAddr())
}

func (session *Session) sendThread(server *TCPServer) {
	defer server.wg.Done()

	for {
		packet, ok := <-session.sendBuffer
		if !ok {
			// 意味着道通已经空了，并且已被关闭
			break
		}
		_, err := session.connection.Write(packet)
		if err != nil {
			break
		}
	}

	log.Printf("session %s sendThread Exit", session.RemoteAddr())
}

// Send 发送数据
func (session *Session) Send(packet []byte) {
	session.sendBuffer <- packet
}

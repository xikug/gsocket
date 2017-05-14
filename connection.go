package gsocket

import (
	"io"
	"log"
	"net"
	"sync"
)

// Connection 代表一个连接会话
type Connection struct {
	conn       net.Conn
	sendBuffer chan []byte
	terminated bool
}

// newConnection 生成一个新的Session
func newConnection(conn net.Conn) (session *Connection) {
	session = &Connection{
		conn:       conn,
		sendBuffer: make(chan []byte, 10),
		terminated: false,
	}

	return session
}

// RemoteAddr 返回客户端的地址和端口
func (session *Connection) RemoteAddr() string {
	return session.conn.RemoteAddr().String()
}

// Close 关闭Session
func (session *Connection) Close() {
	session.terminated = true
	close(session.sendBuffer)
	session.conn.Close()
}

func (session *Connection) recvThread(wg *sync.WaitGroup, handler tcpEventHandler) {
	defer wg.Done()
	buffer := make([]byte, 4096)
	for {
		n, err := session.conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				if handler.handlerError != nil {
					handler.handlerError(session, err)
				}

				break
			}

			if handler.handlerDisconnect != nil {
				handler.handlerDisconnect(session)
			}
			break
		}

		//session.RecvedPackets = append(session.RecvedPackets, buffer[:n]...)
		if handler.handlerRecv != nil {
			handler.handlerRecv(session, buffer[:n])
		}
	}

	if session.terminated == false {
		session.Close()
	}
	log.Printf("session %s recvThread Exit", session.RemoteAddr())
}

func (session *Connection) sendThread(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		packet, ok := <-session.sendBuffer
		if !ok {
			// 意味着道通已经空了，并且已被关闭
			break
		}
		_, err := session.conn.Write(packet)
		if err != nil {
			break
		}
	}

	log.Printf("session %s sendThread Exit", session.RemoteAddr())
}

// Send 发送数据
func (session *Connection) Send(data []byte) {
	session.sendBuffer <- data
}

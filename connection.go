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
func newConnection(conn net.Conn) (c *Connection) {
	c = &Connection{
		conn:       conn,
		sendBuffer: make(chan []byte, 10),
		terminated: false,
	}

	return c
}

// RemoteAddr 返回客户端的地址和端口
func (c *Connection) RemoteAddr() string {
	return c.conn.RemoteAddr().String()
}

// LocalAddr 返回本机地址和端口
func (c *Connection) LocalAddr() string {
	return c.conn.LocalAddr().String()
}

// Close 关闭连接
func (c *Connection) Close() {
	c.terminated = true
	close(c.sendBuffer)
	c.conn.Close()
}

func (c *Connection) recvThread(wg *sync.WaitGroup, handler tcpEventHandler) {
	defer wg.Done()
	buffer := make([]byte, 4096)
	for {
		n, err := c.conn.Read(buffer)
		if err != nil {
			if c.terminated {
				// 直接退出
				break
			}
			if err != io.EOF {
				if handler.handlerError != nil {
					handler.handlerError(c, err)
				}
				break
			}

			if handler.handlerDisconnect != nil {
				handler.handlerDisconnect(c)
			}
			break
		}

		//session.RecvedPackets = append(session.RecvedPackets, buffer[:n]...)
		if handler.handlerRecv != nil {
			handler.handlerRecv(c, buffer[:n])
		}
	}

	if c.terminated == false {
		c.Close()
	}
	log.Printf("session %s recvThread Exit", c.RemoteAddr())
}

func (c *Connection) sendThread(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		packet, ok := <-c.sendBuffer
		if !ok {
			// 意味着道通已经空了，并且已被关闭
			break
		}
		_, err := c.conn.Write(packet)
		if err != nil {
			break
		}
	}

	log.Printf("session %s sendThread Exit", c.RemoteAddr())
}

// Send 发送数据
func (c *Connection) Send(data []byte) {
	c.sendBuffer <- data
}

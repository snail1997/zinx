package znet

import (
	"fmt"
	"net"
	"zinx/ziface"
)

type Connection struct {
	//当前连接的socket TCP套接字
	Conn *net.TCPConn
	//当前连接的ID 也可以称作为SessionID，ID全局唯一
	ConnID uint32
	//当前连接的关闭状态
	isClosed bool

	//当前的处理方法
	Router ziface.IRouter

	//告知该链接已经退出
	ExitBuffChan chan bool
}

// NewConnection  创建连接的方法
func NewConnection(conn *net.TCPConn, connID uint32, router ziface.IRouter) *Connection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		Router:       router,
		isClosed:     false,
		ExitBuffChan: make(chan bool, 1),
	}
	return c
}

// StartReader 处理conn读数据的Goroutine
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), "conn reader exit!")
	defer c.Stop()

	for {
		buf := make([]byte, 512)
		_, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("rev buf err ", err)
			c.ExitBuffChan <- true
			continue
		}
		//得到当前的conn数据的Request数据
		req := Request{
			conn: c,
			data: buf,
		}

		go func(request ziface.IRequest) {
			//从路由中,找到注册绑定的Conn到对应router调用
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

	}
}

func (c *Connection) Start() {
	//开启处理该链接读取到客户端数据之后的请求业务
	go c.StartReader()

	for {
		select {
		case <-c.ExitBuffChan:
			//得到退出消息，不再阻塞
			return
		}
	}
}

// Stop 停止连接，结束当前连接状态M
func (c *Connection) Stop() {
	//1. 如果当前链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	//TODO Connection Stop() 如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用

	//关闭socket
	err := c.Conn.Close()
	if err != nil {
		return
	}
	//通知从缓冲队列读数据的业务，该链接已经关闭
	c.ExitBuffChan <- true
	//关闭该链接全部管道
	close(c.ExitBuffChan)
}

// GetTCPConnection 从当前连接获取原始的socket TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

// GetConnID 获取当前连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

// RemoteAddr 获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

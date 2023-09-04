package ziface

import "net"

// IConnection 定义连接接口
type IConnection interface {
	// Start 启动连接，让当前连接开始工作
	Start()
	// Stop 停止连接，结束当前连接状态
	Stop()
	// GetConnID 从当前连接获取原始的socket TCPConn GetTCPConnection() *net.TCPConn //获取当前连接ID
	GetConnID() uint32
}

// HandFunc 定义一个统⼀处理链接业务的接⼝
type HandFunc func(*net.TCPConn, []byte, int) error

package main

import "zinx/znet"

func main() {
	//1 创建一个server 句柄 s
	s := znet.NewServer("[zinx V0.1]")

	//2 开启服务
	s.Serve()
}

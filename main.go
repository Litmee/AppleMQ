package main

import (
	"AppleMQ/treaty"
	"bufio"
	"fmt"
	"net"
)

func process(c net.Conn) {
	// 处理完毕之后关闭连接
	defer func(c net.Conn) {
		err := c.Close()
		if err != nil {
			return
		}
	}(c)

	for {
		s, err := treaty.Decode(bufio.NewReader(c))
		if err != nil {
			fmt.Printf("read from conn failed, err:%v\n", err)
			break
		}
		fmt.Println("接收到的数据:", s)
	}
}

func main() {
	fmt.Println("Hello AppleMQ")

	// 1. 开启服务
	listen, err := net.Listen("tcp", "127.0.0.1:9082")
	if err != nil {
		fmt.Println("listen failed, err:%V\n", err)
		return
	}

	// 2. 轮询等待客户端来建立连接
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("accept failed, err:%V\n", err)
			continue
		}
		// 3. 启动一个单独的 goroutine 处理连接
		go process(conn)
	}
}

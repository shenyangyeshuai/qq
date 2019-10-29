package serv

import (
	"fmt"
	"net"
	"qq/server/cli"
)

type Config struct {
	Addr string
}

func Run(config *Config) {
	listener, err := net.Listen("tcp", config.Addr)
	if err != nil {
		fmt.Printf("服务器监听端口失败: %v\n", err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("客户端连接失败: %v\n", err)
			continue
		}

		go process(conn)
	}
}

func process(c net.Conn) {
	defer c.Close()

	client := cli.NewClient(c)
	err := client.Process()
	if err != nil {
		fmt.Printf("客户端处理消息失败: %v\n", err)
		return
	}
}

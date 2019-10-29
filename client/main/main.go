package main

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"qq/proto"
)

func login(c net.Conn) error {
	msg := &proto.Message{}
	msg.Cmd = proto.UserLogin

	var loginCmd = &proto.LoginCmd{}
	loginCmd.Id = 1
	loginCmd.Passwd = "123123"

	data, err := json.Marshal(loginCmd)
	if err != nil {
		return err
	}

	msg.Data = string(data)
	data, err = json.Marshal(msg)
	if err != nil {
		fmt.Println(err)
		return err
	}

	var buf [4]byte
	packLen := uint32(len(data))
	binary.BigEndian.PutUint32(buf[0:4], packLen)

	n, err := c.Write(buf[:])
	if err != nil || n != 4 {
		fmt.Println(err)
		return err
	}

	_, err = c.Write([]byte(data))
	if err != nil {
		fmt.Println(err)
		return err
	}

	msg, err = ReadPack(c)
	if err != nil {
		fmt.Println("读取数据包出错", err)
		return err
	}

	return nil
}

func ReadPack(c net.Conn) (*proto.Message, error) {
	var buf [8192]byte
	n, err := c.Read(buf[0:4])
	if n != 4 {
		return nil, errors.New("abc")
	}

	var packLen uint32
	packLen = binary.BigEndian.Uint32(buf[0:4])

	n, err = c.Read(buf[0:packLen])
	if n != int(packLen) {
		return nil, errors.New("def")
	}

	var msg proto.Message
	err = json.Unmarshal(buf[0:packLen], &msg)
	if err != nil {
		return nil, fmt.Errorf("消息反序列化失败: %v", err)
	}

	return &msg, nil
}

func main() {
	c, err := net.Dial("tcp", "127.0.0.1:10000")
	if err != nil {
		fmt.Println("客户端连接出错")
		return
	}
	defer c.Close()

	err = login(c)
	if err != nil {
		fmt.Println("登录失败: %v", err)
		return
	}
}

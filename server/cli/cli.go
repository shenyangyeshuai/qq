package cli

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"qq/proto"
	"qq/server/mgr"
)

var (
	ErrClientHeader   = errors.New("客户端连接格式错误")
	ErrClientPackLen  = errors.New("读取客户端消息长度出错")
	ErrClientPackBody = errors.New("读取客户端消息内容出错")
)

type Client struct {
	c   net.Conn
	buf [8192]byte
}

func NewClient(c net.Conn) *Client {
	return &Client{c: c}
}

func (cli *Client) ReadPack() (*proto.Message, error) {
	n, err := cli.c.Read(cli.buf[0:4])
	if n != 4 {
		return nil, ErrClientHeader
	}

	var packLen uint32
	packLen = binary.BigEndian.Uint32(cli.buf[0:4])

	n, err = cli.c.Read(cli.buf[0:packLen])
	if n != int(packLen) {
		return nil, ErrClientPackBody
	}

	var msg proto.Message
	err = json.Unmarshal(cli.buf[0:packLen], &msg)
	if err != nil {
		return nil, fmt.Errorf("消息反序列化失败: %v", err)
	}

	return &msg, nil
}

func (cli *Client) WritePack(data []byte) error {
	packLen := uint32(len(data))
	binary.BigEndian.PutUint32(cli.buf[0:4], packLen)

	n, err := cli.c.Write(cli.buf[0:4])
	if err != nil {
		return err
	}

	n, err = cli.c.Write(data)
	if err != nil {
		return err
	}
	if n != int(packLen) {
		return fmt.Errorf("数据未写完")
	}

	return nil
}

func (cli *Client) Process() error {
	for {
		msg, err := cli.ReadPack()
		if err != nil {
			return err
		}

		err = cli.processMsg(msg)
		if err != nil {
			return err
		}
	}
}

func (cli *Client) processMsg(msg *proto.Message) error {
	switch msg.Cmd {
	case proto.UserLogin:
		return cli.login(msg)
	case proto.UserRegister:
		return cli.register(msg)
	default:
		return fmt.Errorf("不支持的消息")
	}

	return nil
}

func (cli *Client) loginResp(err error) {
	var msg = &proto.Message{}
	msg.Cmd = proto.UserLoginRes

	var loginCmdRes = &proto.LoginCmdResp{}
	loginCmdRes.Code = 200
	if err != nil {
		loginCmdRes.Code = 500
		loginCmdRes.Err = fmt.Sprintf("%v", err)
	}

	data, err := json.Marshal(loginCmdRes)
	if err != nil {
		fmt.Println("序列化失败: %v", err)
		return
	}

	msg.Data = string(data)

	data, err = json.Marshal(msg)
	if err != nil {
		fmt.Println("序列化失败: %v", err)
		return
	}

	err = cli.WritePack(data)
	if err != nil {
		fmt.Printf("发送失败: %v", err)
		return
	}
}

func (cli *Client) login(msg *proto.Message) error {
	var err error
	defer func() {
		cli.loginResp(err)
	}()

	fmt.Printf("服务器收到客户端登录请求: %v\n", msg)

	var cmd proto.LoginCmd
	err = json.Unmarshal([]byte(msg.Data), &cmd)
	if err != nil {
		return err
	}

	_, err = mgr.Mgr.Login(cmd.Id, cmd.Passwd)
	if err != nil {
		return err
	}

	return nil
}

func (cli *Client) register(msg *proto.Message) error {
	var cmd proto.RegisterCmd
	err := json.Unmarshal([]byte(msg.Data), &cmd)
	if err != nil {
		return err
	}

	err = mgr.Mgr.Register(cmd.User)
	if err != nil {
		return err
	}

	return nil
}

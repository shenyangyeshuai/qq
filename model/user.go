package model

import (
	"errors"
)

const (
	StatusOnline = iota
	StatusOffline
)

var (
	UserTable = "users"

	ErrUserExist           = errors.New("用户已存在")
	ErrUserNotExist        = errors.New("用户不存在")
	ErrUserLoginInvalid    = errors.New("登录用户名/密码错误")
	ErrUserRegisterInvalid = errors.New("注册信息为空")
)

type User struct {
	Id        int    `json:"id"`          // 用户 id : 数字
	Passwd    string `json:"passwd"`      // 用户密码 : 字母数字组合
	Nick      string `json:"nick" `       // 用户昵称 : 显示
	Sex       string `json:"sex"`         // 用户性别 : 字符串
	Avatar    string `json:"avatar" `     // 用户头像 : url
	LastLogin string `json:"last_login" ` // 用户上线登录时间 : 字符串
	Status    int    `json:"status"`      // 用户是否在线 ： online
}

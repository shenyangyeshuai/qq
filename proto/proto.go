package proto

import (
	"qq/model"
)

const (
	UserLogin    = "user_login"
	UserLoginRes = "user_login_res"
	UserRegister = "user_register"
)

type Message struct {
	Cmd  string `json:"cmd"`
	Data string `json:"data"`
}

type LoginCmd struct {
	Id     int    `json:"id"`
	Passwd string `json:"passwd"`
}

type RegisterCmd struct {
	User *model.User `json:"user"`
}

type LoginCmdResp struct {
	Code int    `json:"code"`
	Err  string `json:"err"`
}

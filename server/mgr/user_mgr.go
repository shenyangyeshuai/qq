package mgr

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"qq/model"
	"time"
)

var (
	Mgr *UserMgr
)

type UserMgr struct {
	pool *redis.Pool
}

func NewUserMgr(pool *redis.Pool) *UserMgr {
	return &UserMgr{pool: pool}
}

func InitUserMgr(pool *redis.Pool) {
	Mgr = NewUserMgr(pool)
}

func (mgr *UserMgr) findById(c redis.Conn, id int) (*model.User, error) {
	userInfo, err := redis.String(c.Do("HGet", model.UserTable, fmt.Sprintf("%d", id)))
	if err != nil {
		if err == redis.ErrNil {
			return nil, model.ErrUserNotExist
		}
		return nil, err
	}

	user := model.User{}
	err = json.Unmarshal([]byte(userInfo), &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (mgr *UserMgr) Login(id int, passwd string) (*model.User, error) {
	c := mgr.pool.Get()
	defer c.Close()

	// 查找用户
	user, err := mgr.findById(c, id)
	if err != nil {
		return nil, err
	}

	// 登录的真正业务逻辑
	if user.Id != id || user.Passwd != passwd {
		return nil, model.ErrUserLoginInvalid
	}
	user.Status = model.StatusOnline
	user.LastLogin = fmt.Sprintf("%v", time.Now())

	return user, nil
}

func (mgr *UserMgr) Register(u *model.User) error {
	c := mgr.pool.Get()
	defer c.Close()

	if u == nil {
		fmt.Println("登录用户不能为空")
		return model.ErrUserRegisterInvalid
	}

	// 查找用户
	_, err := mgr.findById(c, u.Id)
	// 找到了, 所以不能注册
	if err == nil {
		return model.ErrUserExist
	}
	// 发生了除没找到的别的错误
	if err != model.ErrUserNotExist {
		return err
	}

	// 能跑到这里, 说明 User Not Exist
	userInfo, err := json.Marshal(u)
	if err != nil {
		return err
	}

	_, err = c.Do("HSet", model.UserTable, fmt.Sprintf("%d", u.Id), string(userInfo))
	if err != nil {
		return err
	}

	return nil
}

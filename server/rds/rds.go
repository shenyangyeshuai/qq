package rds

import (
	"github.com/garyburd/redigo/redis"
	"time"
)

var (
	Pool *redis.Pool
)

type Config struct {
	Addr        string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

func InitRedis(config *Config) {
	Pool = &redis.Pool{
		MaxIdle:     config.MaxIdle,
		MaxActive:   config.MaxActive,
		IdleTimeout: config.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", config.Addr)
		},
	}
}

func GetConn() (redis.Conn, error) {
	return Pool.Get(), nil
}

// TODO: 没太明白这个逻辑
func PutConn(c redis.Conn) {
	c.Close()
}

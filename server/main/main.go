package main

import (
	"qq/server/mgr"
	"qq/server/rds"
	"qq/server/serv"
	"time"
)

func main() {
	// 初始化 redis (连接池 Pool)
	rds.InitRedis(&rds.Config{
		Addr:        "127.0.0.1:2379",
		MaxIdle:     16,
		MaxActive:   0,
		IdleTimeout: time.Second * 300,
	}) // rds 包暴露了一个 Pool 变量

	// 初始化 UserMgr
	mgr.InitUserMgr(rds.Pool) // mgr 包暴露了一个 Mgr

	// 运行
	serv.Run(&serv.Config{
		Addr: "127.0.0.1:10000",
	})
}

package main

import (
	"fmt"
	"minireddit/logger"
	"minireddit/settings"

	"go.uber.org/zap"
)

func main() {
	// 加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed, err:%v\n", err)
		return
	}
	// 初始化日志
	if err := logger.Init(); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	zap.L().Debug("logger init success...")
	// 初始化MySQL连接
	if err := mysql.Init(); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	// 初始化Redis连接
	if err := redis.Init(); err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}
	// 注册路由

	// 启动服务 (优雅关机)
}

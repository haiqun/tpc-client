package main

import (
	"context"
	"tcp_client/http/tcp"
	"tcp_client/providers"
	"tcp_client/routers"
	"github.com/gin-gonic/gin"
)



var ctxCannel context.CancelFunc

func init()  {
	// 初始化配置
	providers.GetConfig()
	// 初始化日志工具
	providers.GetLogger()
	// 初始化网站根目录
	providers.GetRootPath()
	// 启动tcp-server的监听
	var ctx context.Context
	ctx , ctxCannel = context.WithCancel(context.Background())
	go tcp.ClientRun(ctx)
}

func main()  {
	providers.Logger.Info("tcp_clinet_starting....")
	// 获取链接，发送数据和打印接受的数据
	app := providers.Config.GetStringMapString("app")
	// 设置对应的模式
	if app["model"] == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	r.Use(gin.Logger())

	// ===  路由声明 start ====
	api := r.Group("api")
	routers.NewApis(api)
	// ===  路由声明 end ====

	// 端口监听
	if port, exist := app["port"]; exist {
		r.Run(":" + port)
	} else {
		r.Run(":8080")
	}
	ctxCannel();
	providers.Logger.Info("tcp_clinet_end....")
}
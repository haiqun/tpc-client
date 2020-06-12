package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"tcp_client/http/tcp"
	"tcp_client/providers"
	"tcp_client/routers"
)

var CtxCannel context.CancelFunc

func init()  {
	// 初始化变量
	providers.GetRootPath()
	providers.GetLogger()
	providers.GetConfig()
	// 启动tcp - client
	var tcpCom tcp.Remote
	var ctx context.Context
	ctx,CtxCannel = context.WithCancel(context.Background())
	go tcpCom.ClientRun(ctx,2)
}

func main()  {
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

	defer CtxCannel()
}
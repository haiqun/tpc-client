package main

import (
	"github.com/gin-gonic/gin"
	"tcp-tls-project/providers"
	"tcp-tls-project/routers"
)

func init()  {
	// 初始化变量
	providers.GetRootPath()
	providers.GetLogger()
	providers.GetConfig()
	// 启动tcp

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
}
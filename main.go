package main

import (
	"demo_backend/api"
	"demo_backend/middleware"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

func main() {
	s := g.Server()
	// s.EnableHTTPS("./httpsCertification/server.crt", "./httpsCertification/server.key")

	s.Use(middleware.MiddlewareLog)          //业务日志
	s.Use(middleware.MiddlewareErrorHandler) //异常日志
	s.Use(middleware.MiddlewareCORS)         //跨域

	s.Group("/api", func(group *ghttp.RouterGroup) {
		// api文件夹下的文件，不能以api_开头，不然有引入错误
		// 哪些接口需要登录，由各组自己确定
		group.Group("/test", api.RouterGroup_Test)       // 示范
		group.Group("/account", api.RouterGroup_Account) // 账号类
		group.Group("/ws", api.RouterGroup_Websocket)

	})

	s.SetPort(8080)
	s.SetIndexFolder(true) //静态文件
	s.SetServerRoot("./static")
	s.Run()
}

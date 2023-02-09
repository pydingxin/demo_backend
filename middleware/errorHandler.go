package middleware

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

var ctx = context.Background()

func MiddlewareErrorHandler(r *ghttp.Request) {
	r.Middleware.Next()
	if err := r.GetError(); err != nil {
		// 记录到自定义错误日志文件
		g.Log("exception").Error(ctx, err)
		//返回固定的友好信息
		r.Response.ClearBuffer()
		r.Response.WriteJsonExit(g.Map{"status": false, "msg": err})
	}
}

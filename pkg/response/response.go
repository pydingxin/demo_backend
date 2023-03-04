package response

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// 所有返回数据都走这里
func DoneData(r *ghttp.Request, data interface{}) {
	r.Response.WriteJsonExit(g.Map{"status": true, "data": data})
}

func FailMsg(r *ghttp.Request, msg interface{}) {
	r.Response.WriteJsonExit(g.Map{"status": false, "msg": msg})
}

package middleware

import (
	"fmt"

	"github.com/gogf/gf/v2/net/ghttp"
)

func MiddlewareLog(r *ghttp.Request) {
	r.Middleware.Next()
	ses, _ := r.Session.Data()
	fmt.Println(r.Response.Status, r.URL.Path, ses)
}

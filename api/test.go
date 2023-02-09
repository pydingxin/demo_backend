package api

import "github.com/gogf/gf/v2/net/ghttp"

func handler_api_test_test1(r *ghttp.Request) {
	// /api/test/test1
	r.Response.Write("/api/test/test1")
}

func handler_api_test_panic(r *ghttp.Request) {
	// /api/test/panic
	r.Response.Write("/api/test/panic")
	panic("handler_api_test_panic 抛出异常")
}

func RouterGroup_ApiTest(group *ghttp.RouterGroup) {
	group.POST("/test1", handler_api_test_test1)
	group.POST("/panic", handler_api_test_panic)
}

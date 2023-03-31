package api

import (
	"fmt"

	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"demo_backend/middleware"
	"demo_backend/pkg/response"
	"demo_backend/pkg/tool"

	"github.com/gogf/gf/v2/os/gmutex"
)

//-------------------------------------------------------------------------------------------
// allws = map[usename]WebSocket 管理所有websocket
var allws map[string]*ghttp.WebSocket = make(map[string]*ghttp.WebSocket, 100)
var allwsMtx *gmutex.Mutex = gmutex.New()

func allwsSet(username string, ws *ghttp.WebSocket) {
	allwsMtx.LockFunc(func() {
		if oldWs, alreadyExists := allws[username]; alreadyExists {
			oldWs.Close()
		}
		allws[username] = ws
	})
}
func allwsGet(username string) (ws *ghttp.WebSocket) {
	// return nil if not exists
	allwsMtx.RLockFunc(func() {
		if oldWs, alreadyExists := allws[username]; alreadyExists {
			ws = oldWs
		}
	})
	return
}

//-------------------------------------------------------------------------------------------
/* 	ws不带cookie，无法使用session，必须自己搞定身份验证。
	用户登录模块时，把sessionid记录在tool里，专门给本模块验证身份用
	websocket连接时发送用户名+sessionid，只要检验一下就知道是否为该用户

 	url是 ws://ip:port/api/ws/connect?username=dingxin&sessionid=xxxxxxxx
	前端获取sessionid的方式：
	function sessionid(){
		let arr = document.cookie.split(';');
		for(let x of arr){if(x.startsWith("gfsessionid"))return(x.split("=")[1])}
	}
*/
func userOfConnection(r *ghttp.Request) (username string, valid bool) {
	username = r.Get("username").String()
	sessionid := r.Get("sessionid").String()
	if username == "" || sessionid == "" || tool.SessionIdGet(username) != sessionid {
		valid = false // 连接未加参数 或参数对不上
	} else {
		valid = true
	}
	return
}

func handler_api_ws_connect(r *ghttp.Request) {
	// 建立websocket连接
	ws, err := r.WebSocket() //http upgrade to ws
	if err != nil {
		response.FailMsg(r, err)
		return
	}
	//连接已经在上一步成功了，如果不valid，要把它关闭
	if username, valid := userOfConnection(r); valid {
		allwsSet(username, ws)
		fmt.Println("连接成功")
	} else {
		ws.WriteJSON(g.Map{"status": false, "msg": "invalid user"})
		ws.Close()
		fmt.Println("连接失败")
	}
}

//-------------------------------------------------------------------------------------------

func handler_api_ws_channel_make(r *ghttp.Request) {
	// key := r.Get("key").String()
	// val := r.Get("val").String()
	// fmt.Println(key, val)
	// var ctx = gctx.New()
	// _, err := g.Redis().Set(ctx, key, val)
	// if err != nil {
	// 	g.Log().Fatal(ctx, err)
	// }
	// value, err := g.Redis().Get(ctx, key)
	// if err != nil {
	// 	g.Log().Fatal(ctx, err)
	// }
	// response.DoneData(r, value)
}

func RouterGroup_Websocket(group *ghttp.RouterGroup) {
	// 不能登录验证，因为ws连接不带cookie，不能检验gfsessionid
	group.GET("/connect", handler_api_ws_connect) // api/ws/connect

	group.Group("/channel", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.MiddlewareAuth)
		group.POST("/make", handler_api_ws_channel_make) // api/ws/channel/make
	})
}

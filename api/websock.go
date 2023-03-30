package api

import (
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

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
	if username == "" || sessionid == "" {
		valid = false // 连接未加参数
		return
	}
	if tool.SessionIdGet(username) == sessionid {
		valid = true
	} else {
		valid = false
	}
	return
}

//-------------------------------------------------------------------------------------------

func handler_api_ws_connect(r *ghttp.Request) {
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

func RouterGroup_Websocket(group *ghttp.RouterGroup) {
	// 不能登录验证，因为ws连接不带cookie，不能检验gfsessionid
	group.GET("/connect", handler_api_ws_connect)

	// group.Group("/", func(group *ghttp.RouterGroup) {	})
}

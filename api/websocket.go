package api

import (
	"fmt"

	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"demo_backend/middleware"
	"demo_backend/model/dbmodel"
	"demo_backend/model/do"
	"demo_backend/pkg/response"
	"demo_backend/pkg/tool"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gmutex"
)

func accountId(r *ghttp.Request) uint {
	return r.Session.MustGet("accountId").Uint()
}

//-------------------------------------------------------------------------------------------
func init() {
	v, err := g.Redis().Do(gctx.New(), "ping")
	if err != nil || v.String() != "PONG" {
		panic("websocket.go init fail" + err.Error())
	}
}

//-------------------------------------------------------------------------------------------
// allws = map[usename]WebSocket 管理所有websocket
var allws map[uint]*ghttp.WebSocket = make(map[uint]*ghttp.WebSocket, 100)
var allwsMtx *gmutex.Mutex = gmutex.New()

func allwsSet(accountid uint, ws *ghttp.WebSocket) {
	allwsMtx.LockFunc(func() {
		if oldWs, alreadyExists := allws[accountid]; alreadyExists {
			oldWs.Close()
		}
		allws[accountid] = ws
	})
}
func allwsGet(accountid uint) (ws *ghttp.WebSocket) {
	// return nil if not exists
	allwsMtx.RLockFunc(func() {
		if oldWs, alreadyExists := allws[accountid]; alreadyExists {
			ws = oldWs
		} else {
			ws = nil
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
func userOfConnection(r *ghttp.Request) (accountId uint, valid bool) {
	accountId = r.Get("accountid").Uint()
	sessionid := r.Get("sessionid").String()
	if accountId == 0 || sessionid == "" || tool.SessionIdGet(accountId) != sessionid {
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
	if accountid, valid := userOfConnection(r); valid {
		allwsSet(accountid, ws)
		fmt.Println("连接成功")
	} else {
		ws.WriteJSON(g.Map{"status": false, "msg": "invalid user"})
		ws.Close()
		fmt.Println("连接失败")
	}
}

//-------------------------------------------------------------------------------------------

func handler_api_ws_channel_create(r *ghttp.Request) {
	// 创建一个自己的新频道
	data := &dbmodel.WebsocketChannel{AccountID: accountId(r)}
	do.WebsocketChannel_CreateOne(data)

	// 在redis中创建订阅者 members_id
	chmembers := fmt.Sprintf("members_%d", data.ID)
	_, err := g.Redis().Do(gctx.New(), "sadd", chmembers, accountId(r)) //把本人加入members
	if err != nil {
		g.Redis().Do(gctx.New(), "del", chmembers) //若创建失败，删除该键
		response.FailMsg(r, err)
	} else {
		response.DoneData(r, "WebsocketChannel and "+chmembers+" created")
	}
}

func handler_api_ws_channel_delete(r *ghttp.Request) {
	// 本用户删除自己的某个频道
	channelid := r.Get("channelid").Uint()
	if channelid == 0 {
		response.FailMsg(r, "need channelid") //没有输入channelid
	}

	//在库中查找该频道
	ch := do.WebsocketChannel_QueryOneById(channelid)
	if channelid != ch.ID {
		response.FailMsg(r, "channel not found") //没找到该channelid
	} else if ch.AccountID != accountId(r) {
		response.FailMsg(r, "not your channel") //该channel不属于该用户
	}

	// 从数据库删除channel，从redis中删除members_id
	do.WebsocketChannel_DeleteOneById(channelid)
	chmembers := "members_" + r.Get("channelid").String()
	if _, err := g.Redis().Do(gctx.New(), "del", chmembers); err != nil {
		response.FailMsg(r, err)
	} else {
		response.DoneData(r, "ok")
	}
}

func handler_api_ws_channel_subscribe(r *ghttp.Request) {
	// 本用户订阅某个频道
	channelid := r.Get("channelid").Uint()
	if channelid == 0 {
		response.FailMsg(r, "need channelid") //没有输入channelid
	}

	//在库中查找该频道
	ch := do.WebsocketChannel_QueryOneById(channelid)
	if channelid != ch.ID {
		response.FailMsg(r, "channel not found") //没找到该channelid
	}

	// 把 redis键"members_channelid" 作为集合，存储成员的accountid
	chmembers := "members_" + r.Get("channelid").String()
	_, err := g.Redis().Do(gctx.New(), "sadd", chmembers, accountId(r))
	if err != nil {
		response.FailMsg(r, err)
	} else {
		response.DoneData(r, "ok")
	}
}

func handler_api_ws_channel_unsubscribe(r *ghttp.Request) {
	//退订频道，所有者不应退订自己创建的频道
}
func handler_api_ws_channel_publish(r *ghttp.Request) {
	//向频道中发送信息
	// 应当已订阅该频道
}
func handler_api_ws_channel_members(r *ghttp.Request) {
	// 获取某频道当前所有订阅者
}
func RouterGroup_Websocket(group *ghttp.RouterGroup) {
	// 不能登录验证，因为ws连接不带cookie，不能检验gfsessionid.
	// 也无法用postman调试
	group.GET("/connect", handler_api_ws_connect) // api/ws/connect

	group.Group("/channel", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.MiddlewareAuth)
		// channel只需要uuid就能用，它是其他资源的附属产物，其创建销毁不是本模块该负责的事
		group.POST("/create", handler_api_ws_channel_create) // api/ws/channel/create
		group.POST("/delete", handler_api_ws_channel_delete) // api/ws/channel/delete

		group.POST("/subscribe", handler_api_ws_channel_subscribe)     // api/ws/channel/subscribe
		group.POST("/unsubscribe", handler_api_ws_channel_unsubscribe) // api/ws/channel/unsubscribe
		group.POST("/members", handler_api_ws_channel_members)         // api/ws/channel/members
		group.POST("/publish", handler_api_ws_channel_publish)         // api/ws/channel/publish
	})
}

package tool

import (
	"github.com/gogf/gf/v2/os/gmutex"
)

//------------------------------------------------------------------------------------------
// 管理所有用户登录后的sessionid，方便websocket验证状态
var allSessionId map[uint]string = make(map[uint]string, 100)
var allSessionIdMtx *gmutex.Mutex = gmutex.New()

func SessionIdSet(userid uint, sessionid string) {
	allSessionIdMtx.LockFunc(func() {
		allSessionId[userid] = sessionid
	})
}
func SessionIdGet(userid uint) (sessionid string) {
	allSessionIdMtx.RLockFunc(func() {
		if old, alreadyExists := allSessionId[userid]; alreadyExists {
			sessionid = old
		}
	})
	return
}

//------------------------------------------------------------------------------------------

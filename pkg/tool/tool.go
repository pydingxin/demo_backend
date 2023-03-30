package tool

import (
	"github.com/gogf/gf/v2/os/gmutex"
)

//------------------------------------------------------------------------------------------
// 管理所有用户登录后的sessionid，方便websocket验证状态
var allSessionId map[string]string = make(map[string]string, 100)
var allSessionIdMtx *gmutex.Mutex = gmutex.New()

func SessionIdSet(username, id string) {
	allSessionIdMtx.LockFunc(func() {
		allSessionId[username] = id
	})
}
func SessionIdGet(username string) (id string) {
	allSessionIdMtx.RLockFunc(func() {
		if old, alreadyExists := allSessionId[username]; alreadyExists {
			id = old
		}
	})
	return
}

//------------------------------------------------------------------------------------------

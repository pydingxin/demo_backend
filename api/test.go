package api

import (
	"demo_backend/model/do"

	"github.com/gogf/gf/v2/net/ghttp"
)

func RouterGroup_Test(group *ghttp.RouterGroup) {
	group.GET("/excel", handler_api_test_excelDownload)

}
func Hello() {
	print("hello")
}
func handler_api_test_excelDownload(r *ghttp.Request) {
	sac := do.Account_QueryAll_SmallAccount()
	do.SmallAccount_SaveToExcelFile("./test.xlsx", sac, nil)
	r.Response.ServeFileDownload("test.xlsx")
}

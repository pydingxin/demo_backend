package api

import (
	"errors"

	"github.com/gogf/gf/v2/net/ghttp"
)

func RouterGroup_Test(group *ghttp.RouterGroup) {
	group.GET("/excel", handler_api_test_excelDownload)

}

func handler_api_test_excelDownload(r *ghttp.Request) {
	// acs := do.Account_QueryAll()
	// do.Account_SaveToExcelFile("./static/test.xlsx", acs, nil)
	// r.Response.ServeFileDownload("./static/test.xlsx")
}
func hello() error {
	return errors.New("hello")
}

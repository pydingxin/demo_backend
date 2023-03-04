package orm

import (
	"fmt"

	"github.com/gogf/gf/os/gctx"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// 全局数据库连接
var gormConnection *gorm.DB = nil

func initialize_mssql_connection() *gorm.DB {
	// 读取配置的dsn，连接数据库
	dsn := g.Cfg().MustGet(gctx.New(), "database.dsn").String()
	conn, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		CreateBatchSize: 1000,
	})
	if err != nil {
		panic("initialize_mssql_connection(), gorm-sqlserver open error: " + err.Error())
	} else {
		fmt.Println("mssql连接成功")
	}
	conn = conn.Debug() //debug模式，打印sql语句
	return conn
}

func Conn() *gorm.DB {
	// 获取gormConnection，这个连接使用时自动生成副本，当前使用的永远是副本
	if gormConnection == nil {
		gormConnection = initialize_mssql_connection()
	}
	return gormConnection
}

func ToOrmFieldName(column string) string {
	// 把名字转为gorm库表字段的命名格式，如果首字母是大写则转换，否则不处理
	if gstr.IsLetterUpper(column[0]) {
		return gstr.CaseSnakeFirstUpper(column)
	}
	return column
}

func PanicGormResultError(errmsg string, result *gorm.DB) {
	// 检查gorm操作的result，若发生错误，则抛出异常
	if result.Error != nil {
		panic(errmsg + " gorm result error: " + result.Error.Error())
	}
}
func PanicGormResultNotOne(errmsg string, result *gorm.DB) {
	// RowsAffected数量不达预期，这应该是个业务问题。只要操作过程没有error，不应该抛出异常
}

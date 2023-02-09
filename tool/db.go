package tool

import (
	"fmt"

	"github.com/gogf/gf/os/gctx"
	"github.com/gogf/gf/v2/frame/g"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// 全局数据库连接
var gormConnection *gorm.DB = nil

func initialize_mssql_connection() *gorm.DB {
	dsn := g.Cfg().MustGet(gctx.New(), "database.dsn").String()
	conn, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	} else {
		fmt.Println("mssql连接成功")
	}
	return conn
}

func GetGormConnection() *gorm.DB {
	if gormConnection == nil {
		gormConnection = initialize_mssql_connection()
	}
	return gormConnection
}

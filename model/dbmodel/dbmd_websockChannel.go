package dbmodel

import "gorm.io/gorm"

// CMD_DBMODEL_MAKE_ORM_SEMANTIC_FILE 生成orm代码
type WebsocketChannel struct {
	ID        uint           `gorm:"primaryKey"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	AccountID uint           `gorm:"index:ownerid;not null"` //创建者的id
}

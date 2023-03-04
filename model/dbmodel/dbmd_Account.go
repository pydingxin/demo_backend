package dbmodel

import "gorm.io/gorm"

// CMD_DBMODEL_MAKE_ORM_SEMANTIC_FILE 生成orm代码
type Account struct {
	ID        uint           `gorm:"primaryKey"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	User      string         `gorm:"size:40;index:name;not null;unique"`
	Pass      string         `gorm:"size:40;not null"`
}

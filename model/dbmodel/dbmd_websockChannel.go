package dbmodel

// CMD_DBMODEL_MAKE_ORM_SEMANTIC_FILE 生成orm代码
type WebsocketChannel struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"size:50;index:name;not null"`
	Owner string `gorm:"size:40;index:owner;not null"`
}

// CMD_DBMODEL_MAKE_ORM_SEMANTIC_FILE 生成orm代码
type WebsocketChannelFollower struct {
	ChannelID   uint   `gorm:"index:channel_id"`
	ChannelName string `gorm:"size:50;index:channel_name"`
	UserID      uint   `gorm:"index:user_id"`
	UserName    string `gorm:"size:40;index:user_name"`
}

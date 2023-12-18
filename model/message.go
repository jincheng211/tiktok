package model

// Message
type Message struct {
	Content    string `gorm:"column:content;"`                     // 消息内容
	CreateTime int64  `gorm:"column:create_time;"`                 // 消息发送时间 yyyy-MM-dd HH:MM:ss
	FromUserID int64  `gorm:"column:from_user_id;"`                // 消息发送者id
	ID         int64  `gorm:"column:id;primaryKey;autoIncrement;"` // 消息id
	ToUserID   int64  `gorm:"column:to_user_id;"`                  // 消息接收者id
}

func (m *Message) TableName() string {
	return "Message"
}

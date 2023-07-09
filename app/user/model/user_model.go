package model

import (
	"time"

	channelModel "coen-chat/app/channel/model"
	messageModel "coen-chat/app/message/model"
)

type User struct {
	Email     string    `gorm:"primaryKey;column:email;not null;size:191" json:"email"`
	Password  string    `gorm:"column:password;not null" json:"password"`
	Name      string    `gorm:"column:name;not null" json:"name"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"createdAt"`

	Channels []channelModel.Channel `gorm:"foreignKey:created_by" json:"channelList"`
	Members  []channelModel.Member  `gorm:"foreignKey:email" json:"memberList"`
	Messages []messageModel.Message `gorm:"foreignKey:email" json:"messageList"`
}

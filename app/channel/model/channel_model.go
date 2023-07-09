package model

import (
	"time"

	"coen-chat/app/message/model"
)

type Channel struct {
	ChannelID   uint64    `gorm:"primaryKey;column:channel_id" json:"channelId"`
	ChannelName string    `gorm:"column:channel_name;not null" json:"channelName"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not null" json:"updatedAt"`
	CreatedBy   string    `gorm:"column:created_by;not null;size:191" json:"createdBy"`

	Messages []model.Message `gorm:"foreignKey:ChannelID;constraint:onDelete:CASCADE"`
	Members  []Member        `gorm:"foreignKey:ChannelID;constraint:onDelete:CASCADE"`
}

type Member struct {
	Email     string `gorm:"primaryKey" json:"email"`
	ChannelID uint64 `gorm:"primaryKey" json:"channelId"`
}

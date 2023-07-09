package model

import "time"

type Message struct {
	MessageID string    `gorm:"primaryKey;column:message_id;not null" json:"messageId"`
	ChannelID uint64    `gorm:"column:channel_id;not null" json:"channelId"`
	Email     string    `gorm:"column:email;not null" json:"email"`
	Content   string    `gorm:"column:content;not null" json:"content"`
	CreatedAt time.Time `gorm:"column:created_at;not null" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null" json:"updatedAt"`
}

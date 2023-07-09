package model

import "time"

type ElasticsearchModel struct {
	MessageId string    `json:"messageId,omitempty"`
	ChannelId int       `json:"channelId,omitempty"`
	Email     string    `json:"email,omitempty"`
	Content   string    `json:"content,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

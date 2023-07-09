package repository

import (
	channelModel "coen-chat/app/channel/model"
	"coen-chat/app/message/model"

	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{
		db: db,
	}
}

func (r *MessageRepository) IsJoined(member *channelModel.Member) bool {
	var count int64
	r.db.First(member).Count(&count)
	return count > 0
}

func (r *MessageRepository) FindById(messageID string) (*model.Message, error) {
	var message = &model.Message{MessageID: messageID}
	err := r.db.First(message).Error
	return message, err
}

func (r *MessageRepository) GetMessageList(channelID uint64, offset int) ([]*model.Message, error) {
	pageLimit := 10
	var messageList []*model.Message
	err := r.db.Where("channel_id = ?", channelID).
		Order("created_at desc").
		Limit(pageLimit).
		Offset(pageLimit * offset).
		Find(&messageList).Error
	return messageList, err
}

func (r *MessageRepository) CreateMessage(message *model.Message) error {
	return r.db.Create(&message).Error
}

func (r *MessageRepository) UpdateMessage(message *model.Message) error {
	return r.db.Save(&message).Error
}

func (r *MessageRepository) DeleteMessage(message *model.Message) error {
	return r.db.Delete(&message).Error
}

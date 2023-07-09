package repository

import (
	"coen-chat/app/channel/model"

	"gorm.io/gorm"
)

type ChannelRepository struct {
	db *gorm.DB
}

func NewChannelRepository(db *gorm.DB) *ChannelRepository {
	return &ChannelRepository{
		db: db,
	}
}

func (r *ChannelRepository) FindByIdPreload(channelID uint64) (*model.Channel, error) {
	pageLimit := 10
	var channel = &model.Channel{ChannelID: channelID}
	err := r.db.Preload("Members").Preload("Messages", func(db *gorm.DB) *gorm.DB {
		return db.Order("messages.created_at desc").Limit(pageLimit).Offset(pageLimit * 0)
	}).First(channel).Error
	return channel, err
}

func (r *ChannelRepository) FindById(channelID uint64) (*model.Channel, error) {
	var channel = &model.Channel{ChannelID: channelID}
	err := r.db.First(channel).Error
	return channel, err
}

func (r *ChannelRepository) GetChannelList(isJoined bool, email string) ([]*model.Channel, error) {
	var channelList []*model.Channel
	var err error
	if isJoined {
		err = r.db.Table("channels").Select("*").
			Where("EXISTS (?)",
				r.db.Table("members").
					Where("channels.channel_id = members.channel_id").
					Where("members.email = ?", email).
					Select("*"),
			).Find(&channelList).Error
	} else {
		err = r.db.Table("channels").Select("*").
			Where("NOT EXISTS (?)",
				r.db.Table("members").
					Where("channels.channel_id = members.channel_id").
					Where("members.email = ?", email).
					Select("*"),
			).Find(&channelList).Error
	}
	return channelList, err
}

func (r *ChannelRepository) CreateChannel(channel *model.Channel) error {
	return r.db.Create(&channel).Error
}

func (r *ChannelRepository) UpdateChannel(channel *model.Channel) error {
	return r.db.Save(&channel).Error
}

func (r *ChannelRepository) DeleteChannel(channel *model.Channel) error {
	return r.db.Delete(&channel).Error
}

func (r *ChannelRepository) JoinChannel(member *model.Member) error {
	return r.db.Create(&member).Error
}

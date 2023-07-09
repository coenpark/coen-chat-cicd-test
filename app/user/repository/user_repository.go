package repository

import (
	"time"

	"coen-chat/app/user/model"
	"coen-chat/utils"

	"github.com/go-redis/redis/v7"
	"gorm.io/gorm"
)

type UserRepository struct {
	db    *gorm.DB
	redis *redis.Client
}

func NewUserRepository(db *gorm.DB, redis *redis.Client) *UserRepository {
	return &UserRepository{
		db:    db,
		redis: redis,
	}
}

func (r *UserRepository) IsExist(email string) bool {
	var count int64
	r.db.Model(&model.User{}).Where("email = ?", email).
		Where("created_at != ''").
		Count(&count)
	return count > 0
}

func (r *UserRepository) CreateUser(user *model.User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *UserRepository) FindByEmail(email string) (model.User, error) {
	var user = model.User{Email: email}
	result := r.db.First(&user)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}

func (r *UserRepository) CreateAuth(email string, td *utils.TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) //converting Unix to UTC
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := r.redis.Set(td.AccessUUID, email, at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}
	errRefresh := r.redis.Set(td.RefreshUUID, email, rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func (r *UserRepository) DeleteAuth(UUID string) (int64, error) {
	deleted, err := r.redis.Del(UUID).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}

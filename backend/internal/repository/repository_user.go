package repository

import (
	"github.com/asliddinberdiev/eirsystem/config"
	"github.com/asliddinberdiev/eirsystem/internal/model"
	"github.com/asliddinberdiev/eirsystem/pkg/logger"
	"github.com/asliddinberdiev/eirsystem/pkg/redis"
	"gorm.io/gorm"
)

type User interface {
	GetAll() ([]model.User, error)
	GetByID(id string) (model.User, error)
	GetByUsername(username string) (model.User, error)
}

type userRepo struct {
	cfg    *config.Config
	logger logger.Logger
	db     *gorm.DB
	rd     *redis.RedisClient
}

func NewUserRepository(cfg *config.Config, logger logger.Logger, db *gorm.DB, rd *redis.RedisClient) User {
	return &userRepo{cfg: cfg, logger: logger, db: db, rd: rd}
}

func (r *userRepo) GetAll() ([]model.User, error) {
	var users []model.User
	return users, r.db.Find(&users).Error
}

func (r *userRepo) GetByID(id string) (model.User, error) {
	var user model.User
	return user, r.db.Where("id = ?", id).Take(&user).Error
}

func (r *userRepo) GetByUsername(username string) (model.User, error) {
	var user model.User
	return user, r.db.Where("username = ?", username).Take(&user).Error
}

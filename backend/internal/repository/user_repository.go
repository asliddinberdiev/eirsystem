package repository

import (
	"github.com/asliddinberdiev/eirsystem/config"
	"github.com/asliddinberdiev/eirsystem/internal/domain"
	"github.com/asliddinberdiev/eirsystem/pkg/logger"
	"github.com/asliddinberdiev/eirsystem/pkg/redis"
	"gorm.io/gorm"
)

type User interface {
	GetAll() ([]domain.User, error)
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

func (r *userRepo) GetAll() ([]domain.User, error) {
	var users []domain.User
	return users, r.db.Find(&users).Error
}

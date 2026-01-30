// Package repository - Repository layer
package repository

import (
	"github.com/asliddinberdiev/eirsystem/config"
	"github.com/asliddinberdiev/eirsystem/pkg/logger"
	"github.com/asliddinberdiev/eirsystem/pkg/redis"
	"gorm.io/gorm"
)

type Repository struct {
	User User
}

func New(cfg *config.Config, logger logger.Logger, db *gorm.DB, rd *redis.RedisClient) *Repository {
	return &Repository{
		User: NewUserRepository(cfg, logger, db, rd),
	}
}
// Package service - Service layer
package service

import (
	"github.com/asliddinberdiev/eirsystem/config"
	"github.com/asliddinberdiev/eirsystem/internal/repository"
	"github.com/asliddinberdiev/eirsystem/pkg/logger"
	"github.com/asliddinberdiev/eirsystem/pkg/minio"
)

type Service struct {
	User User
}

func New(cfg *config.Config, logger logger.Logger, s3 *minio.Client, repo *repository.Repository) *Service {
	return &Service{
		User: NewUserService(cfg, logger, s3, repo),
	}
}

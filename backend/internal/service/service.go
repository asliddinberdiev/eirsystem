// Package service - Service layer
package service

import (
	"github.com/asliddinberdiev/eirsystem/config"
	"github.com/asliddinberdiev/eirsystem/internal/repository"
	"github.com/asliddinberdiev/eirsystem/pkg/logger"
	"github.com/asliddinberdiev/eirsystem/pkg/minio"
	"github.com/casbin/casbin/v3"
)

type Service struct {
	User   User
	Policy Policy
}

func New(cfg *config.Config, logger logger.Logger, s3 *minio.Client, repo *repository.Repository, enforcer *casbin.Enforcer) *Service {
	return &Service{
		User:   NewUserService(cfg, logger, s3, repo),
		Policy: NewPolicyService(enforcer),
	}
}

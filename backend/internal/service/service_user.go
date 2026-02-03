package service

import (
	"github.com/asliddinberdiev/eirsystem/config"
	"github.com/asliddinberdiev/eirsystem/internal/model"
	"github.com/asliddinberdiev/eirsystem/internal/repository"
	"github.com/asliddinberdiev/eirsystem/pkg/logger"
	"github.com/asliddinberdiev/eirsystem/pkg/minio"
)

type User interface {
	GetAll() ([]model.User, error)
	GetByUsername(username string) (model.User, error)
}

type userServ struct {
	cfg    *config.Config
	logger logger.Logger
	s3     *minio.Client
	repo   *repository.Repository
}

func NewUserService(cfg *config.Config, logger logger.Logger, s3 *minio.Client, repo *repository.Repository) User {
	return &userServ{
		cfg:    cfg,
		logger: logger,
		s3:     s3,
		repo:   repo,
	}
}

func (s *userServ) GetAll() ([]model.User, error) {
	return s.repo.User.GetAll()
}

func (s *userServ) GetByUsername(username string) (model.User, error) {
	return s.repo.User.GetByUsername(username)
}

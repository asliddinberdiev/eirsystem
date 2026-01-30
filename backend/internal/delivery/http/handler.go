// Package http - HTTP delivery layer
package http

import (
	"github.com/asliddinberdiev/eirsystem/config"
	"github.com/asliddinberdiev/eirsystem/internal/delivery/http/middleware"
	v1 "github.com/asliddinberdiev/eirsystem/internal/delivery/http/v1"
	"github.com/asliddinberdiev/eirsystem/internal/service"
	"github.com/asliddinberdiev/eirsystem/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	cfg *config.Config
	log logger.Logger
	svc *service.Service
}

func New(cfg *config.Config, log logger.Logger, svc *service.Service) *Handler {
	return &Handler{
		cfg: cfg,
		log: log,
		svc: svc,
	}
}

func (h *Handler) InitRouter(cfg *config.App) *gin.Engine {
	if !cfg.IsDev() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(middleware.RequestID())
	router.Use(gin.Recovery())
	router.Use(logger.GinLogger(h.log))
	router.Use(cors.Default())

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.cfg, h.log, h.svc)

	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}

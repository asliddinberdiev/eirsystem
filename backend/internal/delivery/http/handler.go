// Package http - HTTP delivery layer
package http

import (
	"github.com/asliddinberdiev/eirsystem/config"
	"github.com/asliddinberdiev/eirsystem/internal/delivery/http/middleware"
	v1 "github.com/asliddinberdiev/eirsystem/internal/delivery/http/v1"
	"github.com/asliddinberdiev/eirsystem/internal/service"
	"github.com/asliddinberdiev/eirsystem/pkg/logger"
	"github.com/asliddinberdiev/eirsystem/pkg/validator"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	cfg         *config.Config
	log         logger.Logger
	valid       validator.Validator
	redisClient *redis.Client
	svc         *service.Service
}

func New(cfg *config.Config, log logger.Logger, redisClient *redis.Client, svc *service.Service) *Handler {
	return &Handler{
		cfg:         cfg,
		log:         log,
		valid:       validator.New(),
		redisClient: redisClient,
		svc:         svc,
	}
}

func (h *Handler) InitRouter() *gin.Engine {
	if !h.cfg.App.IsDev() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(cors.Default())
	router.Use(middleware.RequestID())
	router.Use(gin.Recovery())
	router.Use(logger.GinLogger(h.log))
	router.Use(middleware.NewRateLimiter(h.log, h.redisClient, "120-S", "app"))

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.cfg, h.log, h.valid, h.svc)

	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}

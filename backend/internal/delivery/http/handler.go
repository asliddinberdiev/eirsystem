// Package http - HTTP delivery layer
package http

import (
	"github.com/asliddinberdiev/eirsystem/config"
	"github.com/asliddinberdiev/eirsystem/internal/delivery/http/middleware"
	v1 "github.com/asliddinberdiev/eirsystem/internal/delivery/http/v1"
	"github.com/asliddinberdiev/eirsystem/internal/service"
	"github.com/asliddinberdiev/eirsystem/pkg/jwt"
	"github.com/asliddinberdiev/eirsystem/pkg/logger"
	"github.com/asliddinberdiev/eirsystem/pkg/validator"
	"github.com/casbin/casbin/v3"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type Handler struct {
	cfg         *config.Config
	log         logger.Logger
	valid       validator.Validator
	jwtManager  *jwt.Manager
	redisClient *redis.Client
	svc         *service.Service
	enforcer    *casbin.Enforcer
}

func New(cfg *config.Config, log logger.Logger, redisClient *redis.Client, svc *service.Service, enforcer *casbin.Enforcer) *Handler {
	return &Handler{
		cfg:         cfg,
		log:         log,
		valid:       validator.New(),
		jwtManager:  jwt.New(&cfg.JWT, redisClient),
		redisClient: redisClient,
		svc:         svc,
		enforcer:    enforcer,
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
	handlerV1 := v1.NewHandler(h.cfg, h.log, h.valid, h.jwtManager, h.svc, h.enforcer)

	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}

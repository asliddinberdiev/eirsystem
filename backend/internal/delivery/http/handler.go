package http

import (
	"github.com/asliddinberdiev/eirsystem/config"
	"github.com/asliddinberdiev/eirsystem/internal/delivery/http/middleware"
	v1 "github.com/asliddinberdiev/eirsystem/internal/delivery/http/v1"
	"github.com/asliddinberdiev/eirsystem/pkg/logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handler struct {
	db  *gorm.DB
	log logger.Logger
}

func NewHandler(db *gorm.DB, log logger.Logger) *Handler {
	return &Handler{
		db:  db,
		log: log,
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

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "pong",
		})
	})

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.db, h.log)

	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}

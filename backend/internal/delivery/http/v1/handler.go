// Package v1 - HTTP v1 delivery layer
package v1

import (
	"github.com/asliddinberdiev/eirsystem/config"
	_ "github.com/asliddinberdiev/eirsystem/docs"
	"github.com/asliddinberdiev/eirsystem/internal/service"
	"github.com/asliddinberdiev/eirsystem/pkg/logger"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	cfg *config.Config
	log logger.Logger
	svc *service.Service
}

// @title EIR System API
// @version v1

// @contact.name Asliddin Berdiev
// @contact.url https://github.com/asliddinberdiev

// @host localhost:8080
// @BasePath /api/v1

func NewHandler(cfg *config.Config, log logger.Logger, svc *service.Service) *Handler {
	return &Handler{
		cfg: cfg,
		log: log,
		svc: svc,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	v1.Use(gzip.Gzip(gzip.BestSpeed))
	{
		// swagger docs
		v1.GET(
			"/docs/*any",
			gin.BasicAuth(gin.Accounts{h.cfg.App.DocsUsername: h.cfg.App.DocsPassword}),
			ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/api/v1/docs/doc.json")))
	}

	{
		h.initUserRoutes(v1)
	}
}

func (h *Handler) initUserRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.GET("", h.GetAll)
	}
}

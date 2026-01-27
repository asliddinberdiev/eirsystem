package v1

import (
	"github.com/asliddinberdiev/eirsystem/pkg/logger"
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

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initUserRoutes(v1)
	}
}

func (h *Handler) initUserRoutes(api *gin.RouterGroup) {
	users := api.Group("/users")
	{
		users.GET("", h.GetUsers)
	}
}
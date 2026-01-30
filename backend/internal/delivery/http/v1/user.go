// Package v1 - HTTP v1 delivery layer
package v1

import (
	"github.com/asliddinberdiev/eirsystem/pkg/codes"
	"github.com/asliddinberdiev/eirsystem/pkg/logger"
	"github.com/asliddinberdiev/eirsystem/pkg/response"
	"github.com/gin-gonic/gin"
)

// GetAll godoc
// @Summary Get users
// @Description Fetch list of users
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users [get]
func (h *Handler) GetAll(c *gin.Context) {
	users, err := h.svc.User.GetAll()
	if err != nil {
		response.Error(c, h.log, codes.InternalError, err)
		return
	}

	h.log.Info("Users", logger.Int("count", len(users)))

	response.Success(c, codes.Ok, users)
}

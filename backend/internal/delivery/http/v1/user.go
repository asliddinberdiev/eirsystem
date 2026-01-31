// Package v1 - HTTP v1 delivery layer
package v1

import (
	"github.com/asliddinberdiev/eirsystem/pkg/codes"
	"github.com/asliddinberdiev/eirsystem/pkg/response"
	"github.com/gin-gonic/gin"
)

// GetAll godoc
// @Summary Get users
// @Description Fetch list of users
// @Tags users
// @Accept  json
// @Produce  json
// @Response 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users [get]
func (h *Handler) GetAll(c *gin.Context) {
	users := []string{"user1", "user2", "user3"}

	response.Success(c, codes.Ok, users)
}

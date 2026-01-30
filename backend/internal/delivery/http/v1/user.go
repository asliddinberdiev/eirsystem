// Package v1 - HTTP v1 delivery layer
package v1

import (
	"github.com/asliddinberdiev/eirsystem/pkg/codes"
	"github.com/asliddinberdiev/eirsystem/pkg/response"
	"github.com/gin-gonic/gin"
)

// GetUsers godoc
// @Summary Get users
// @Description Fetch list of users
// @Tags users
// @Accept  json
// @Produce  json
// @Success 200 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /users [get]
func (h *Handler) GetUsers(c *gin.Context) {
	users := map[string]string{
		"name": "Alex",
	}

	response.Success(c, codes.Ok, users)
}

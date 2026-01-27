package v1

import (
	"fmt"

	"github.com/asliddinberdiev/eirsystem/pkg/codes"
	"github.com/asliddinberdiev/eirsystem/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *Handler) GetUsers(c *gin.Context) {
	h.log.Info("Fetched users successfully")

	response.Error(c, h.log, codes.InternalError, fmt.Errorf("dasdsad"))
}

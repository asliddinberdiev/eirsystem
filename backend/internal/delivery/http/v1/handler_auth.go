// Package v1 - HTTP v1 delivery layer
package v1

import (
	"errors"

	"github.com/asliddinberdiev/eirsystem/internal/dto"
	"github.com/asliddinberdiev/eirsystem/pkg/codes"
	"github.com/asliddinberdiev/eirsystem/pkg/hasher"
	"github.com/asliddinberdiev/eirsystem/pkg/response"
	"github.com/gin-gonic/gin"
)

func (h *Handler) initAuthRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/sign-in", h.SignIn)
	}
}

// SignIn godoc
// @Summary SignIn
// @Description SignIn
// @Tags auth
// @Accept  json
// @Produce  json
// @Param request body dto.SignInRequest true "SignIn Request"
// @Response 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/sign-in [post]
func (h *Handler) SignIn(c *gin.Context) {
	var req dto.SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, h.log, codes.InvalidRequest, err)
		return
	}

	if err := h.valid.Struct(&req); err != nil {
		response.Error(c, h.log, codes.InvalidRequest, err)
		return
	}

	user, err := h.svc.User.GetByUsername(req.Username)
	if err != nil {
		response.Error(c, h.log, codes.AuthInvalidCredentials, errors.New("username or password is incorrect"))
		return
	}

	if err := hasher.Verify(req.Password, user.PasswordHash); err != nil {
		response.Error(c, h.log, codes.AuthInvalidCredentials, errors.New("username or password is incorrect"))
		return
	}

	resp := dto.SignInResponse{
		AccessToken:  "access_token",
		RefreshToken: "refresh_token",
		User: dto.User{
			ID:       user.ID,
			Username: user.Username,
			FullName: user.FullName,
			Role:     user.Role,
		},
	}

	response.Success(c, codes.Ok, resp)
}

// logout

// refresh

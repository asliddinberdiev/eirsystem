// Package v1 - HTTP v1 delivery layer
package v1

import (
	"errors"
	"strings"

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
		auth.POST("/refresh", h.Refresh)
		auth.POST("/logout", h.Logout)
	}
}

// SignIn godoc
// @Summary SignIn
// @Description Foydalanuvchi tizimga kirishi va token olishi
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
	accessToken, refreshToken, err := h.jwt.Generate(c.Request.Context(), user.ID, c.Request.UserAgent(), c.ClientIP())
	if err != nil {
		response.Error(c, h.log, codes.InternalError, err)
		return
	}

	resp := dto.SignInResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.User{
			ID:       user.ID,
			Username: user.Username,
			FullName: user.FullName,
			Role:     user.Role,
		},
	}

	response.Success(c, codes.Ok, resp)
}

// Refresh godoc
// @Summary Refresh Token
// @Description Access tokenni yangilash
// @Tags auth
// @Accept  json
// @Produce  json
// @Param request body dto.RefreshTokenRequest true "Refresh Request"
// @Response 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /auth/refresh [post]
// @Security BearerAuth
func (h *Handler) Refresh(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, h.log, codes.InvalidRequest, err)
		return
	}

	if err := h.valid.Struct(&req); err != nil {
		response.Error(c, h.log, codes.InvalidRequest, err)
		return
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		response.Error(c, h.log, codes.AuthAccessTokenRequired, errors.New("authorization header required"))
		return
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		response.Error(c, h.log, codes.AuthTokenInvalid, errors.New("invalid auth header format"))
		return
	}
	accessToken := headerParts[1]

	claims, err := h.jwt.ParseTokenUnverified(accessToken)
	if err != nil {
		response.Error(c, h.log, codes.AuthTokenInvalid, err)
		return
	}

	newAccess, newRefresh, err := h.jwt.Refresh(c.Request.Context(), claims.UserID, claims.SessionID, req.RefreshToken, c.Request.UserAgent())
	if err != nil {
		response.Error(c, h.log, codes.AuthTokenInvalid, err)
		return
	}

	response.Success(c, codes.Ok, dto.RefreshTokenResponse{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	})
}

// Logout godoc
// @Summary Logout
// @Description Tizimdan chiqish (Sessiyani o'chirish)
// @Tags auth
// @Security ApiKeyAuth
// @Response 200 {object} response.Response
// @Router /auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		response.Success(c, codes.Ok, nil)
		return
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := h.jwt.ParseTokenUnverified(tokenStr)
	if err != nil {
		response.Success(c, codes.Ok, nil)
		return
	}

	err = h.jwt.Logout(c.Request.Context(), claims.UserID, claims.SessionID)
	if err != nil {
		response.Error(c, h.log, codes.InternalError, err)
		return
	}

	response.Success(c, codes.Ok, nil)
}

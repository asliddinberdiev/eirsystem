package middleware

import (
	"errors"
	"strings"

	"github.com/asliddinberdiev/eirsystem/internal/service"
	"github.com/asliddinberdiev/eirsystem/pkg/codes"
	"github.com/asliddinberdiev/eirsystem/pkg/jwt"
	"github.com/asliddinberdiev/eirsystem/pkg/logger"
	"github.com/asliddinberdiev/eirsystem/pkg/response"
	"github.com/gin-gonic/gin"
)

func NewJWTMiddleware(log logger.Logger, jwt *jwt.Manager, svc *service.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Error("STEP 1")
			response.Error(c, log, codes.AuthAccessTokenRequired, errors.New("authorization header required"))
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			response.Error(c, log, codes.AuthTokenInvalid, errors.New("invalid auth header format"))
			return
		}
		accessToken := headerParts[1]

		claims, err := jwt.ParseTokenUnverified(accessToken)
		if err != nil {
			response.Error(c, log, codes.AuthTokenInvalid, err)
			return
		}

		user, err := svc.User.GetByID(claims.UserID)
		if err != nil {
			response.Error(c, log, codes.AuthTokenInvalid, err)
			return
		}

		log.Info("User found", logger.Any("user", user))

		c.Set("userID", user.ID)
		c.Set("userRole", user.Role)
		c.Set("tenantID", user.TenantID)
		c.Next()
	}
}

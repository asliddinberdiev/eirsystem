package middleware

import (
	"errors"
	"fmt"

	"github.com/asliddinberdiev/eirsystem/pkg/codes"
	"github.com/asliddinberdiev/eirsystem/pkg/logger"
	"github.com/asliddinberdiev/eirsystem/pkg/response"
	"github.com/casbin/casbin/v3"
	"github.com/gin-gonic/gin"
)

func Authorizer(log logger.Logger, e *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			response.Error(c, log, codes.AuthAccessTokenRequired, errors.New("user not authenticated"))
			return
		}
		sub := fmt.Sprintf("%v", userID)

		dom, exists := c.Get("tenantID")
		if !exists {
			response.Error(c, log, codes.AuthAccessTokenRequired, errors.New("tenant not found"))
			return
		}

		obj := c.Request.URL.Path
		act := c.Request.Method

		ok, err := e.Enforce(sub, dom, obj, act)
		if err != nil {
			response.Error(c, log, codes.AuthAccessTokenRequired, errors.New("authorization error"))
			return
		}

		if !ok {
			response.Error(c, log, codes.AuthAccessTokenRequired, errors.New("permission denied"))
			return
		}

		c.Next()
	}
}

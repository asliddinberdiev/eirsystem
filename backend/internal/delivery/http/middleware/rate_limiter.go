package middleware

import (
	"errors"
	"fmt"

	"github.com/asliddinberdiev/eirsystem/pkg/codes"
	"github.com/asliddinberdiev/eirsystem/pkg/logger"
	"github.com/asliddinberdiev/eirsystem/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	sredis "github.com/ulule/limiter/v3/drivers/store/redis"
)

func NewRateLimiter(log logger.Logger, client *redis.Client, limitString string, keyPrefix string) gin.HandlerFunc {
	rate, err := limiter.NewRateFromFormatted(limitString)
	if err != nil {
		log.Fatal("Rate limiter configuration error", logger.Error(err))
	}

	store, err := sredis.NewStoreWithOptions(client, limiter.StoreOptions{
		Prefix:   keyPrefix,
		MaxRetry: 3,
	})
	if err != nil {
		log.Fatal("Redis store creation failed", logger.Error(err))
	}

	instance := limiter.New(store, rate)

	middleware := mgin.NewMiddleware(
		instance,
		mgin.WithLimitReachedHandler(func(c *gin.Context) {
			response.Error(c, log, codes.TooManyRequests, errors.New("rate limit exceeded"))
		}),
		mgin.WithKeyGetter(func(c *gin.Context) string {
			if userID, exists := c.Get("userID"); exists {
				return fmt.Sprintf("user:%v", userID)
			}

			return c.ClientIP()
		}),
	)

	return middleware
}

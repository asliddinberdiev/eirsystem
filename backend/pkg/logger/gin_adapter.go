package logger

import (
	"time"

	"github.com/gin-gonic/gin"
)

func GinLogger(l Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		reqID := c.GetString("requestID")
		end := time.Now()
		latency := end.Sub(start)

		fields := []Field{
			String("request_id", reqID),
			String("ip", c.ClientIP()),
			Int("status", c.Writer.Status()),
			Duration("latency", latency),
			String("method", c.Request.Method),
			String("path", path),
			String("query", query),
		}

		if len(c.Errors) > 0 {
			fields = append(fields, String("errors", c.Errors.String()))
			l.Error("HTTP Failed", fields...)
		} else {
			l.Info("HTTP", fields...)
		}
	}
}

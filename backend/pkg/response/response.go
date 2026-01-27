package response

import (
	"fmt"
	"net/http"

	"github.com/asliddinberdiev/eirsystem/pkg/codes"
	"github.com/asliddinberdiev/eirsystem/pkg/logger"
	"github.com/asliddinberdiev/eirsystem/pkg/telegram"
	"github.com/gin-gonic/gin"
)

type Response struct {
	Success   bool   `json:"success"`
	Code      int    `json:"code"`
	Message   string `json:"message,omitempty"`
	Data      any    `json:"data,omitempty"`
	Error     any    `json:"error,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

func Success(c *gin.Context, code codes.Code, data any) {
	reqID := c.GetString("requestID")

	c.JSON(http.StatusOK, Response{
		Success:   true,
		Code:      int(code),
		Message:   code.String(),
		Data:      data,
		RequestID: reqID,
	})
}

func Error(c *gin.Context, log logger.Logger, code codes.Code, err error) {
	reqID := c.GetString("requestID")
	status := code.HTTPStatus()

	var technicalError any
	if err != nil {
		technicalError = err.Error()
		log.Error(
			code.String(),
			logger.String("request_id", reqID),
			logger.Int("code", int(code)),
			logger.Error(err),
		)

		if status >= 500 {
			msg := fmt.Sprintf(
				"🚨 <b>CRITICAL ERROR</b>\n\n"+
					"🆔 <b>RequestID:</b> <code>%s</code>\n"+
					"🔢 <b>Status:</b> %d\n"+
					"🛣 <b>Path:</b> %s %s\n"+
					"❌ <b>Error:</b> <pre>%s</pre>",
				reqID, status, c.Request.Method, c.Request.URL.Path, err.Error(),
			)

			telegram.Send(msg)
		}
	}

	if gin.Mode() == gin.ReleaseMode {
		technicalError = nil
	}

	c.AbortWithStatusJSON(code.HTTPStatus(), Response{
		Success:   false,
		Code:      int(code),
		Message:   code.String(),
		Error:     technicalError,
		RequestID: reqID,
	})
}

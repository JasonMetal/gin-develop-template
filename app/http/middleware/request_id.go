package middleware

import (
	"github.com/JasonMetal/submodule-support-go.git/helper/logger"
	"github.com/JasonMetal/submodule-support-go.git/helper/number"
	"github.com/gin-gonic/gin"
)

const (
	XRequestIDKey = "X-Traceid"
)

// RequestId is a middleware that injects a 'X-Request-ID' into the context and request/response header of each request.
func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check for incoming header, use it if exists
		rid := c.GetHeader(XRequestIDKey)
		if rid == "" {
			rid = number.GenerateTraceId()
			c.Request.Header.Set(XRequestIDKey, rid)
		}
		c.Set(XRequestIDKey, rid)

		// Set XRequestIDKey header
		c.Writer.Header().Set(XRequestIDKey, rid)
		c.Next()
	}
}
func Context() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(logger.KeyRequestID, c.GetString(XRequestIDKey))
		//c.Set(log.KeyUsername, c.GetString(UsernameKey))
		c.Next()
	}
}

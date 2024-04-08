package middleware

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

const (
	// XRequestIDKey defines X-Request-ID key string.
	XRequestIDKey = "X-Request-ID"
)

func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		rid := c.Request.Header.Get(XRequestIDKey)
		if rid == "" {
			rid = uuid.NewV4().String()
		}
		c.Set(XRequestIDKey, rid)
		c.Next()
	}
}

package middleware

import (
	"github.com/gin-gonic/gin"
)

// CheckUserAuth 校验用户权限信息
func CheckUserAuth() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Next()
	}
}

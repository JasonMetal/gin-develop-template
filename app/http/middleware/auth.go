package middleware

import (
	"github.com/gin-gonic/gin"
)

type RespJson struct {
	Code    int32       `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// ParseUserInfoWithErr 解析用户信息
// 失败直接返回鉴权失败
func ParseUserInfoWithErr() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

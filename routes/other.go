package router

import (
	"github.com/gin-gonic/gin"
)

func RegisterOther(router *gin.Engine) {
	// gin-framework 的健康检测
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(0, "Health check for Gin framework:init pong")
		return
	})

}

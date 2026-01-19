package helloRouter

import (
	"develop-template/app/http/controller/api/helloController"
	"develop-template/app/http/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterHello(router *gin.Engine) {
	router.Use(middleware.RequestId())
	router.Use(middleware.Context())

	router.GET("/test", func(ctx *gin.Context) {
		helloController.NewController(ctx).Hello()
	})
}

package helloRouter

import (
	"develop-template/app/http/controller/api/helloController"
	"github.com/gin-gonic/gin"
)

func RegisterHello(router *gin.Engine) {

	router.GET("/test", func(ctx *gin.Context) {
		helloController.NewController(ctx).Hello()
	})
}

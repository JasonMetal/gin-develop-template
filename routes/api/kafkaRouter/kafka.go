package kafkaRouter

import (
	"develop-template/app/http/controller/api/kafkaController"
	"develop-template/app/http/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterKafka(router *gin.Engine) {
	router.Use(middleware.RequestId())
	router.Use(middleware.Context())

	apiGroup := router.Group("kafka")

	// Kafka消息发送接口
	apiGroup.POST("send", func(ctx *gin.Context) {
		kafkaController.NewController(ctx).SendMessage()
	})

	apiGroup.POST("send-json", func(ctx *gin.Context) {
		kafkaController.NewController(ctx).SendJSON()
	})

	// Kafka消息查询接口
	apiGroup.GET("messages", func(ctx *gin.Context) {
		kafkaController.NewController(ctx).FetchMessages()
	})

	// Kafka主题信息接口
	apiGroup.GET("topic-info", func(ctx *gin.Context) {
		kafkaController.NewController(ctx).GetTopicInfo()
	})
}

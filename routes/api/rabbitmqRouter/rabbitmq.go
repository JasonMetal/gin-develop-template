package rabbitmqRouter

import (
	"develop-template/app/http/controller/api/rabbitmqController"
	"develop-template/app/http/middleware"
	"github.com/gin-gonic/gin"
)

func RegisterRabbitMQ(router *gin.Engine) {
	router.Use(middleware.RequestId())
	router.Use(middleware.Context())

	// RabbitMQ基础消息发送接口
	router.POST("/rabbitmq/send", func(ctx *gin.Context) {
		rabbitmqController.NewController(ctx).SendMessage()
	})

	router.POST("/rabbitmq/send-json", func(ctx *gin.Context) {
		rabbitmqController.NewController(ctx).SendJSON()
	})

	// RabbitMQ交换机模式接口
	router.POST("/rabbitmq/send-exchange", func(ctx *gin.Context) {
		rabbitmqController.NewController(ctx).SendToExchange()
	})

	router.POST("/rabbitmq/send-fanout", func(ctx *gin.Context) {
		rabbitmqController.NewController(ctx).SendFanout()
	})

	router.POST("/rabbitmq/send-direct", func(ctx *gin.Context) {
		rabbitmqController.NewController(ctx).SendDirect()
	})

	router.POST("/rabbitmq/send-topic", func(ctx *gin.Context) {
		rabbitmqController.NewController(ctx).SendTopic()
	})

	// RabbitMQ任务和批量接口
	router.POST("/rabbitmq/send-task", func(ctx *gin.Context) {
		rabbitmqController.NewController(ctx).SendTask()
	})

	router.POST("/rabbitmq/send-batch", func(ctx *gin.Context) {
		rabbitmqController.NewController(ctx).SendBatch()
	})

	// RabbitMQ健康检查接口
	router.GET("/rabbitmq/health", func(ctx *gin.Context) {
		rabbitmqController.NewController(ctx).HealthCheck()
	})

	// ============= RabbitMQ查询接口 =============

	// 获取队列信息
	router.GET("/rabbitmq/queue/info", func(ctx *gin.Context) {
		rabbitmqController.NewController(ctx).GetQueueInfo()
	})

	// 查看队列消息（不消费）
	router.GET("/rabbitmq/queue/peek", func(ctx *gin.Context) {
		rabbitmqController.NewController(ctx).PeekMessages()
	})

	// 消费队列消息（会删除）
	router.GET("/rabbitmq/queue/consume", func(ctx *gin.Context) {
		rabbitmqController.NewController(ctx).ConsumeMessages()
	})

	// 清空队列
	router.POST("/rabbitmq/queue/purge", func(ctx *gin.Context) {
		rabbitmqController.NewController(ctx).PurgeQueue()
	})

	// 删除队列
	router.POST("/rabbitmq/queue/delete", func(ctx *gin.Context) {
		rabbitmqController.NewController(ctx).DeleteQueue()
	})
}

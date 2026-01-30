package main

import (
	"develop-template/constant"
	router "develop-template/routes"
	"log"

	"github.com/JasonMetal/submodule-support-go.git/bootstrap"
	"github.com/gin-gonic/gin"
)

// ProjectName

// @title        业务应用开发API
// @description  提供业务应用开发的业务功能APIs
// @schemes      http https
func main() {
	bootstrap.SetProjectName(constant.ProjectName)
	// 初始化Web
	bootstrap.Init()

	// ✅ 设置错误处理器
	bootstrap.SetAsyncErrorHandler(func(topic, message string, err error) {
		log.Printf("❌ Kafka发送失败 - Topic: %s, Error: %v", topic, err)

		// 你的处理逻辑：
		// 1. 保存到重试队列
		// 2. 发送告警
		// 3. 记录日志
	})
	// 启动服务...
	middleFun := []gin.HandlerFunc{
		//	middleware.CheckUserAuth(),
	}
	r := bootstrap.InitWeb(middleFun)
	router.RegisterRouter(r)
	bootstrap.RunWeb(r, constant.HttpServiceHostPort)
}

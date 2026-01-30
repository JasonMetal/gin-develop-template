package main

import (
	"log"
	"time"

	"github.com/JasonMetal/submodule-support-go.git/bootstrap"
)

// 这是一个简单的 Kafka 异步错误处理示例
// 可以直接复制到你的 main.go 或初始化代码中

func main() {
	// 1. 初始化项目
	bootstrap.Init()

	// 2. ✅ 设置 Kafka 异步错误处理器
	setupKafkaErrorHandler()

	// 3. 测试发送消息
	testKafkaAsync()
}

// setupKafkaErrorHandler 设置 Kafka 错误处理器
func setupKafkaErrorHandler() {
	bootstrap.SetAsyncErrorHandler(func(topic, message string, err error) {
		// 📝 打印错误日志
		log.Printf("❌ Kafka异步发送失败")
		log.Printf("   Topic: %s", topic)
		log.Printf("   Error: %v", err)
		log.Printf("   Message: %s", truncateString(message, 100))

		// 🔄 保存到重试队列（你可以实现这个函数）
		saveToRetryQueue(topic, message, err)

		// 🚨 如果是重要主题，发送告警
		if isImportantTopic(topic) {
			sendAlert(topic, err)
		}
	})

	log.Println("✅ Kafka异步错误处理器已设置")
}

// isImportantTopic 判断是否是重要主题
func isImportantTopic(topic string) bool {
	importantTopics := []string{
		"order-created",   // 订单创建
		"payment-success", // 支付成功
		"user-registered", // 用户注册
	}

	for _, t := range importantTopics {
		if topic == t {
			return true
		}
	}
	return false
}

// saveToRetryQueue 保存到重试队列
func saveToRetryQueue(topic, message string, err error) {
	// TODO: 实现保存到数据库或 Redis 的逻辑
	// 示例：
	// db := bootstrap.GetMysqlInstance("default")
	// sql := `INSERT INTO kafka_retry_queue (topic, message, error, retry_time) VALUES (?, ?, ?, ?)`
	// db.Exec(sql, topic, message, err.Error(), time.Now().Add(5*time.Minute))

	log.Printf("💾 已保存到重试队列: Topic=%s", topic)
}

// sendAlert 发送告警
func sendAlert(topic string, err error) {
	// TODO: 实现告警逻辑（钉钉、邮件、短信等）
	log.Printf("🚨 告警：重要主题 %s 发送失败！Error: %v", topic, err)
}

// truncateString 截断字符串
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// testKafkaAsync 测试异步发送
func testKafkaAsync() {
	log.Println("\n开始测试 Kafka 异步发送...")

	// 测试1：发送普通消息
	if err := bootstrap.ProducerAsync("test-topic", "Hello Kafka!"); err != nil {
		log.Printf("发送失败: %v", err)
	}

	// 测试2：发送重要消息
	if err := bootstrap.ProducerAsync("order-created", `{"order_id":"12345","amount":999.99}`); err != nil {
		log.Printf("发送失败: %v", err)
	}

	// 等待一下，让异步消息有时间处理
	time.Sleep(2 * time.Second)

	log.Println("测试完成")
}

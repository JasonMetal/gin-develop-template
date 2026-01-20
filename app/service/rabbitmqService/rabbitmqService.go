package rabbitmqService

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/JasonMetal/submodule-support-go.git/bootstrap"
)

// RabbitMQService RabbitMQ服务
type RabbitMQService struct{}

// NewRabbitMQService 创建RabbitMQ服务
func NewRabbitMQService() *RabbitMQService {
	return &RabbitMQService{}
}

// SendMessage 发送单条消息到队列(同步)
func (s *RabbitMQService) SendMessage(queueName string, message string) error {
	return bootstrap.PublishSimple(queueName, message)
}

// SendMessageWithContext 带上下文发送消息
func (s *RabbitMQService) SendMessageWithContext(ctx context.Context, queueName string, message string) error {
	return bootstrap.PublishSimpleWithContext(ctx, queueName, message)
}

// SendJSON 发送JSON格式数据
func (s *RabbitMQService) SendJSON(queueName string, data interface{}) error {
	return bootstrap.PublishJSON(queueName, data)
}

// SendJSONWithContext 带上下文发送JSON数据
func (s *RabbitMQService) SendJSONWithContext(ctx context.Context, queueName string, data interface{}) error {
	return bootstrap.PublishJSONWithContext(ctx, queueName, data)
}

// SendToExchange 发送消息到交换机
func (s *RabbitMQService) SendToExchange(exchangeName string, exchangeType string, routingKey string, message string) error {
	return bootstrap.PublishToExchange(exchangeName, exchangeType, routingKey, message)
}

// SendToExchangeWithContext 带上下文发送消息到交换机
func (s *RabbitMQService) SendToExchangeWithContext(ctx context.Context, exchangeName string, exchangeType string, routingKey string, message string) error {
	return bootstrap.PublishToExchangeWithContext(ctx, exchangeName, exchangeType, routingKey, message)
}

// SendLog 发送日志消息
func (s *RabbitMQService) SendLog(queueName string, logLevel string, logMessage string, extra map[string]interface{}) error {
	logData := map[string]interface{}{
		"level":     logLevel,
		"message":   logMessage,
		"timestamp": time.Now().Format(time.RFC3339),
		"extra":     extra,
	}
	return s.SendJSON(queueName, logData)
}

// SendEvent 发送事件消息
func (s *RabbitMQService) SendEvent(queueName string, eventType string, eventData interface{}) error {
	event := map[string]interface{}{
		"event_type": eventType,
		"event_data": eventData,
		"timestamp":  time.Now().Format(time.RFC3339),
	}
	return s.SendJSON(queueName, event)
}

// PublishFanout 发送广播消息(fanout模式)
func (s *RabbitMQService) PublishFanout(exchangeName string, message string) error {
	return s.SendToExchange(exchangeName, "fanout", "", message)
}

// PublishDirect 发送直接消息(direct模式)
func (s *RabbitMQService) PublishDirect(exchangeName string, routingKey string, message string) error {
	return s.SendToExchange(exchangeName, "direct", routingKey, message)
}

// PublishTopic 发送主题消息(topic模式)
func (s *RabbitMQService) PublishTopic(exchangeName string, routingKey string, message string) error {
	return s.SendToExchange(exchangeName, "topic", routingKey, message)
}

// SendBatch 批量发送消息
func (s *RabbitMQService) SendBatch(queueName string, messages []string) error {
	var errors []error
	for i, message := range messages {
		if err := s.SendMessage(queueName, message); err != nil {
			errors = append(errors, fmt.Errorf("消息 %d 发送失败: %v", i, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("批量发送失败，错误数: %d, 首个错误: %v", len(errors), errors[0])
	}

	return nil
}

// SendJSONBatch 批量发送JSON消息
func (s *RabbitMQService) SendJSONBatch(queueName string, dataList []interface{}) error {
	var errors []error
	for i, data := range dataList {
		if err := s.SendJSON(queueName, data); err != nil {
			errors = append(errors, fmt.Errorf("消息 %d 发送失败: %v", i, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("批量发送失败，错误数: %d, 首个错误: %v", len(errors), errors[0])
	}

	return nil
}

// SendTask 发送任务消息(Worker模式)
func (s *RabbitMQService) SendTask(queueName string, taskName string, taskData interface{}) error {
	task := map[string]interface{}{
		"task_name":  taskName,
		"task_data":  taskData,
		"created_at": time.Now().Unix(),
		"status":     "pending",
	}
	return s.SendJSON(queueName, task)
}

// SendDelayedMessage 发送延迟消息(需要RabbitMQ延迟插件)
func (s *RabbitMQService) SendDelayedMessage(queueName string, message string, delaySeconds int) error {
	// 这需要 rabbitmq-delayed-message-exchange 插件
	// 这里提供基础实现，实际使用需要在bootstrap层添加延迟支持
	return fmt.Errorf("延迟消息功能需要RabbitMQ延迟插件支持")
}

// ValidateConnection 验证连接状态
func (s *RabbitMQService) ValidateConnection() error {
	manager := bootstrap.GetRabbitMQManager()
	if manager == nil {
		return fmt.Errorf("RabbitMQ未初始化")
	}
	return nil
}

// FormatMessage 格式化消息为JSON字符串
func (s *RabbitMQService) FormatMessage(data interface{}) (string, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("JSON序列化失败: %v", err)
	}
	return string(jsonData), nil
}

// ============= 查询相关方法 =============

// GetQueueInfo 获取队列信息
func (s *RabbitMQService) GetQueueInfo(queueName string) (*bootstrap.QueueInfo, error) {
	return bootstrap.GetQueueInfo(queueName)
}

// DeclareAndGetQueueInfo 声明并获取队列信息
func (s *RabbitMQService) DeclareAndGetQueueInfo(queueName string) (*bootstrap.QueueInfo, error) {
	return bootstrap.DeclareAndGetQueueInfo(queueName)
}

// PeekMessages 查看队列消息（不消费，查看后重新入队）
func (s *RabbitMQService) PeekMessages(queueName string, count int) ([]bootstrap.QueueMessage, error) {
	return bootstrap.PeekMessages(queueName, count)
}

// ConsumeMessages 消费队列消息（会从队列中删除）
func (s *RabbitMQService) ConsumeMessages(queueName string, count int) ([]bootstrap.QueueMessage, error) {
	return bootstrap.ConsumeMessages(queueName, count)
}

// PurgeQueue 清空队列
func (s *RabbitMQService) PurgeQueue(queueName string) (int, error) {
	return bootstrap.PurgeQueue(queueName)
}

// DeleteQueue 删除队列
func (s *RabbitMQService) DeleteQueue(queueName string, ifUnused, ifEmpty bool) (int, error) {
	return bootstrap.DeleteQueue(queueName, ifUnused, ifEmpty)
}

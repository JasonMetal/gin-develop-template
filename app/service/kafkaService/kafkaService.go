package kafkaService

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/JasonMetal/submodule-support-go.git/bootstrap"
)

// KafkaService Kafka服务
type KafkaService struct{}

// NewKafkaService 创建Kafka服务
func NewKafkaService() *KafkaService {
	return &KafkaService{}
}

// SendMessage 发送单条消息(同步)
func (s *KafkaService) SendMessage(topic string, message string) error {
	return bootstrap.ProducerSync(topic, message)
}

// SendMessageAsync 发送单条消息(异步)
func (s *KafkaService) SendMessageAsync(topic string, message string) error {
	return bootstrap.ProducerAsync(topic, message)
}

// SendMessageWithContext 带上下文发送消息
func (s *KafkaService) SendMessageWithContext(ctx context.Context, topic string, message string) error {
	return bootstrap.ProducerSyncWithContext(ctx, topic, message)
}

// SendJSON 发送JSON格式数据
func (s *KafkaService) SendJSON(topic string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("JSON序列化失败: %v", err)
	}
	return s.SendMessage(topic, string(jsonData))
}

// SendJSONAsync 异步发送JSON格式数据
func (s *KafkaService) SendJSONAsync(topic string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("JSON序列化失败: %v", err)
	}
	return s.SendMessageAsync(topic, string(jsonData))
}

// SendBatch 批量发送消息
func (s *KafkaService) SendBatch(topic string, messages []string) error {
	return bootstrap.ProducerSyncBatch(topic, messages)
}

// SendLog 发送日志消息
func (s *KafkaService) SendLog(topic string, logLevel string, logMessage string, extra map[string]interface{}) error {
	logData := map[string]interface{}{
		"level":     logLevel,
		"message":   logMessage,
		"timestamp": time.Now().Format(time.RFC3339),
		"extra":     extra,
	}
	return s.SendJSON(topic, logData)
}

// SendEvent 发送事件消息
func (s *KafkaService) SendEvent(topic string, eventType string, eventData interface{}) error {
	event := map[string]interface{}{
		"event_type": eventType,
		"event_data": eventData,
		"timestamp":  time.Now().Format(time.RFC3339),
	}
	return s.SendJSON(topic, event)
}

// SendMetric 发送指标消息
func (s *KafkaService) SendMetric(topic string, metricName string, metricValue float64, tags map[string]string) error {
	metric := map[string]interface{}{
		"metric_name":  metricName,
		"metric_value": metricValue,
		"tags":         tags,
		"timestamp":    time.Now().Unix(),
	}
	return s.SendJSON(topic, metric)
}

// FetchMessages 获取指定主题和分区的消息
func (s *KafkaService) FetchMessages(topic string, partition int32, offset int64, limit int) ([]*bootstrap.KafkaMessage, error) {
	if limit <= 0 {
		limit = 10 // 默认10条
	}
	if limit > 100 {
		limit = 100 // 最多100条
	}

	return bootstrap.FetchMessages(topic, partition, offset, limit)
}

// FetchMessagesFromAllPartitions 从所有分区获取消息
func (s *KafkaService) FetchMessagesFromAllPartitions(topic string, offset int64, limit int) ([]*bootstrap.KafkaMessage, error) {
	if limit <= 0 {
		limit = 10 // 默认10条
	}
	if limit > 100 {
		limit = 100 // 最多100条
	}

	return bootstrap.FetchMessagesFromAllPartitions(topic, offset, limit)
}

// GetTopicPartitions 获取主题的分区信息
func (s *KafkaService) GetTopicPartitions(topic string) ([]int32, error) {
	return bootstrap.GetTopicPartitions(topic)
}

// GetPartitionOffset 获取分区的偏移量信息
func (s *KafkaService) GetPartitionOffset(topic string, partition int32) (oldest int64, newest int64, err error) {
	return bootstrap.GetPartitionOffset(topic, partition)
}

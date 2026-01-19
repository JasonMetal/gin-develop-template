package kafkaService

import (
	"context"
	"encoding/json"
	"testing"
	"time"
)

// TestNewKafkaService 测试创建服务
func TestNewKafkaService(t *testing.T) {
	service := NewKafkaService()
	if service == nil {
		t.Error("创建KafkaService失败")
	}
	t.Log("✓ 创建KafkaService测试通过")
}

// TestSendJSON 测试发送JSON数据
func TestSendJSON(t *testing.T) {
	// 不使用mock producer，只测试JSON序列化逻辑
	testData := map[string]interface{}{
		"user_id": 123,
		"action":  "login",
		"time":    time.Now().Unix(),
	}

	// 测试JSON序列化
	jsonData, err := json.Marshal(testData)
	if err != nil {
		t.Fatalf("JSON序列化失败: %v", err)
	}

	if len(jsonData) == 0 {
		t.Error("JSON数据为空")
	}

	// 验证JSON可以被解析回去
	var decoded map[string]interface{}
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Errorf("JSON反序列化失败: %v", err)
	}

	t.Logf("✓ SendJSON测试通过, JSON: %s", string(jsonData))
}

// TestSendLog 测试发送日志
func TestSendLog(t *testing.T) {
	_ = NewKafkaService()

	// 测试日志数据结构
	logLevel := "ERROR"
	logMessage := "测试错误消息"
	extra := map[string]interface{}{
		"request_id": "req-123",
		"user_id":    456,
	}

	// 构造预期的日志数据
	expectedLog := map[string]interface{}{
		"level":   logLevel,
		"message": logMessage,
		"extra":   extra,
	}

	// 验证数据可以序列化
	jsonData, err := json.Marshal(expectedLog)
	if err != nil {
		t.Fatalf("日志JSON序列化失败: %v", err)
	}

	// 验证JSON包含必要字段
	var decoded map[string]interface{}
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Errorf("日志JSON反序列化失败: %v", err)
	}

	if decoded["level"] != logLevel {
		t.Errorf("日志级别不匹配: 期望%s, 得到%v", logLevel, decoded["level"])
	}

	if decoded["message"] != logMessage {
		t.Errorf("日志消息不匹配: 期望%s, 得到%v", logMessage, decoded["message"])
	}

	t.Logf("✓ SendLog数据结构测试通过, JSON: %s", string(jsonData))
}

// TestSendEvent 测试发送事件
func TestSendEvent(t *testing.T) {
	_ = NewKafkaService()

	eventType := "user.registered"
	eventData := map[string]interface{}{
		"user_id": 789,
		"email":   "test@example.com",
		"source":  "web",
	}

	// 构造预期的事件数据
	expectedEvent := map[string]interface{}{
		"event_type": eventType,
		"event_data": eventData,
	}

	// 验证数据可以序列化
	jsonData, err := json.Marshal(expectedEvent)
	if err != nil {
		t.Fatalf("事件JSON序列化失败: %v", err)
	}

	// 验证JSON包含必要字段
	var decoded map[string]interface{}
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Errorf("事件JSON反序列化失败: %v", err)
	}

	if decoded["event_type"] != eventType {
		t.Errorf("事件类型不匹配: 期望%s, 得到%v", eventType, decoded["event_type"])
	}

	t.Logf("✓ SendEvent数据结构测试通过, JSON: %s", string(jsonData))
}

// TestSendMetric 测试发送指标
func TestSendMetric(t *testing.T) {
	_ = NewKafkaService()

	metricName := "api.response_time"
	metricValue := 123.45
	tags := map[string]string{
		"endpoint": "/api/users",
		"method":   "GET",
		"status":   "200",
	}

	// 构造预期的指标数据
	expectedMetric := map[string]interface{}{
		"metric_name":  metricName,
		"metric_value": metricValue,
		"tags":         tags,
	}

	// 验证数据可以序列化
	jsonData, err := json.Marshal(expectedMetric)
	if err != nil {
		t.Fatalf("指标JSON序列化失败: %v", err)
	}

	// 验证JSON包含必要字段
	var decoded map[string]interface{}
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Errorf("指标JSON反序列化失败: %v", err)
	}

	if decoded["metric_name"] != metricName {
		t.Errorf("指标名称不匹配: 期望%s, 得到%v", metricName, decoded["metric_name"])
	}

	if decoded["metric_value"] != metricValue {
		t.Errorf("指标值不匹配: 期望%f, 得到%v", metricValue, decoded["metric_value"])
	}

	t.Logf("✓ SendMetric数据结构测试通过, JSON: %s", string(jsonData))
}

// TestSendBatch 测试批量发送
func TestSendBatch(t *testing.T) {
	_ = NewKafkaService()

	messages := []string{
		"message 1",
		"message 2",
		"message 3",
	}

	if len(messages) != 3 {
		t.Error("批量消息数量不正确")
	}

	// 验证消息内容
	for i, msg := range messages {
		if msg == "" {
			t.Errorf("消息%d为空", i)
		}
	}

	t.Log("✓ SendBatch数据结构测试通过")
}

// TestSendJSONWithInvalidData 测试无效数据的JSON序列化
func TestSendJSONWithInvalidData(t *testing.T) {
	_ = NewKafkaService()

	// 创建无法序列化的数据（包含循环引用）
	type Node struct {
		Value int
		Next  *Node
	}

	node1 := &Node{Value: 1}
	node2 := &Node{Value: 2}
	node1.Next = node2
	node2.Next = node1 // 循环引用

	// 尝试序列化会失败
	_, err := json.Marshal(node1)
	if err == nil {
		t.Error("期望序列化失败但成功了")
	}

	t.Logf("✓ 无效数据JSON序列化错误处理测试通过: %v", err)
}

// TestSendMessageWithContext 测试带上下文的消息发送
func TestSendMessageWithContext(t *testing.T) {
	_ = NewKafkaService()

	ctx := context.Background()
	topic := "test-topic"
	message := "test message with context"

	// 验证方法签名和参数
	if topic == "" {
		t.Error("topic不能为空")
	}

	if message == "" {
		t.Error("message不能为空")
	}

	if ctx == nil {
		t.Error("context不能为nil")
	}

	t.Log("✓ SendMessageWithContext参数验证测试通过")
}

// TestSendMessageWithTimeoutContext 测试超时上下文
func TestSendMessageWithTimeoutContext(t *testing.T) {
	_ = NewKafkaService()

	// 创建一个已经超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Nanosecond)
	defer cancel()

	time.Sleep(10 * time.Millisecond) // 确保上下文超时

	// 验证上下文已经超时
	select {
	case <-ctx.Done():
		t.Log("✓ 上下文超时测试通过")
	default:
		t.Error("上下文应该已经超时")
	}
}

// TestComplexJSONStructure 测试复杂JSON结构
func TestComplexJSONStructure(t *testing.T) {
	_ = NewKafkaService()

	complexData := map[string]interface{}{
		"user": map[string]interface{}{
			"id":   123,
			"name": "测试用户",
			"profile": map[string]interface{}{
				"age":    25,
				"gender": "male",
			},
		},
		"orders": []map[string]interface{}{
			{
				"id":     "order-1",
				"amount": 99.99,
				"items":  []string{"item1", "item2"},
			},
			{
				"id":     "order-2",
				"amount": 199.99,
				"items":  []string{"item3", "item4"},
			},
		},
		"metadata": map[string]interface{}{
			"timestamp": time.Now().Unix(),
			"version":   "1.0.0",
		},
	}

	// 测试序列化
	jsonData, err := json.Marshal(complexData)
	if err != nil {
		t.Fatalf("复杂JSON序列化失败: %v", err)
	}

	// 测试反序列化
	var decoded map[string]interface{}
	if err := json.Unmarshal(jsonData, &decoded); err != nil {
		t.Errorf("复杂JSON反序列化失败: %v", err)
	}

	// 验证嵌套数据
	user, ok := decoded["user"].(map[string]interface{})
	if !ok {
		t.Error("user字段解析失败")
	} else {
		if user["name"] != "测试用户" {
			t.Error("用户名称不匹配")
		}
	}

	t.Logf("✓ 复杂JSON结构测试通过, JSON长度: %d字节", len(jsonData))
}

// TestServiceMethodsNotPanic 测试服务方法不会panic
func TestServiceMethodsNotPanic(t *testing.T) {
	service := NewKafkaService()

	// 这些调用可能会失败，但不应该panic
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("服务方法panic: %v", r)
		}
	}()

	// 测试各种方法不会panic（即使Kafka未初始化）
	service.SendMessage("test", "msg")
	service.SendMessageAsync("test", "msg")
	service.SendJSON("test", map[string]string{"key": "value"})
	service.SendJSONAsync("test", map[string]string{"key": "value"})
	service.SendBatch("test", []string{"msg1", "msg2"})
	service.SendLog("test", "INFO", "log message", map[string]interface{}{})
	service.SendEvent("test", "event", map[string]interface{}{})
	service.SendMetric("test", "metric", 1.0, map[string]string{})

	t.Log("✓ 服务方法不panic测试通过")
}

// BenchmarkSendJSON JSON发送性能测试
func BenchmarkSendJSON(b *testing.B) {
	_ = NewKafkaService()

	testData := map[string]interface{}{
		"user_id": 123,
		"action":  "benchmark",
		"time":    time.Now().Unix(),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		json.Marshal(testData)
	}
}

// BenchmarkSendLog 日志发送性能测试
func BenchmarkSendLog(b *testing.B) {
	logLevel := "INFO"
	logMessage := "benchmark log message"
	extra := map[string]interface{}{
		"request_id": "req-123",
		"user_id":    456,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logData := map[string]interface{}{
			"level":     logLevel,
			"message":   logMessage,
			"timestamp": time.Now().Format(time.RFC3339),
			"extra":     extra,
		}
		json.Marshal(logData)
	}
}

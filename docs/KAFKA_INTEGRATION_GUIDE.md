# Kafka 集成使用指南

## 简介

本文档说明如何在 gin-develop-template 项目中使用 Kafka 消息队列功能。

## 快速开始

### 1. 配置 Kafka

编辑对应环境的配置文件 `config/{env}/kafka.yml`:

```yaml
brokers:
  - host: "localhost"
    port: 9092

ssl:
  enable: false

producer:
  required_acks: 1
  max_retries: 5
  return_successes: true
  return_errors: true
```

### 2. 项目启动

Kafka 会在项目启动时自动初始化，无需额外配置。

```bash
go run http-server.go -e local
```

### 3. 使用 Kafka 服务

```go
import "develop-template/app/service/kafkaService"

func main() {
    service := kafkaService.NewKafkaService()
    
    // 发送消息
    err := service.SendMessage("test-topic", "Hello Kafka!")
    if err != nil {
        log.Printf("发送失败: %v", err)
    }
}
```

## API 参考

### 基础消息发送

#### SendMessage - 同步发送消息

```go
func (s *KafkaService) SendMessage(topic string, message string) error
```

**参数**:
- `topic`: Kafka主题名称
- `message`: 要发送的消息内容

**返回**: 错误信息，成功返回 nil

**示例**:
```go
err := service.SendMessage("user-logs", "User login: userID=123")
```

#### SendMessageAsync - 异步发送消息

```go
func (s *KafkaService) SendMessageAsync(topic string, message string) error
```

异步发送，不等待确认，吞吐量更高。

**示例**:
```go
err := service.SendMessageAsync("analytics", "Page view: /home")
```

#### SendMessageWithContext - 带上下文发送

```go
func (s *KafkaService) SendMessageWithContext(ctx context.Context, topic string, message string) error
```

支持超时和取消操作。

**示例**:
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
err := service.SendMessageWithContext(ctx, "orders", "New order")
```

### 结构化数据发送

#### SendJSON - 发送JSON数据

```go
func (s *KafkaService) SendJSON(topic string, data interface{}) error
```

自动将数据序列化为JSON格式发送。

**示例**:
```go
user := map[string]interface{}{
    "id": 123,
    "name": "张三",
    "action": "login",
}
err := service.SendJSON("user-events", user)
```

#### SendJSONAsync - 异步发送JSON

```go
func (s *KafkaService) SendJSONAsync(topic string, data interface{}) error
```

**示例**:
```go
metrics := map[string]interface{}{
    "cpu": 45.5,
    "memory": 78.2,
}
err := service.SendJSONAsync("system-metrics", metrics)
```

#### SendBatch - 批量发送

```go
func (s *KafkaService) SendBatch(topic string, messages []string) error
```

批量发送多条消息，适合批处理场景。

**示例**:
```go
messages := []string{
    "log entry 1",
    "log entry 2",
    "log entry 3",
}
err := service.SendBatch("bulk-logs", messages)
```

### 业务场景方法

#### SendLog - 发送日志

```go
func (s *KafkaService) SendLog(topic string, logLevel string, logMessage string, extra map[string]interface{}) error
```

发送结构化日志消息。

**参数**:
- `topic`: 日志主题
- `logLevel`: 日志级别 (DEBUG/INFO/WARN/ERROR)
- `logMessage`: 日志消息
- `extra`: 额外的上下文信息

**示例**:
```go
extra := map[string]interface{}{
    "request_id": "req-123",
    "user_id": 456,
    "ip": "192.168.1.100",
}
err := service.SendLog(
    "app-logs",
    "ERROR",
    "数据库连接失败",
    extra,
)
```

**生成的消息格式**:
```json
{
    "level": "ERROR",
    "message": "数据库连接失败",
    "timestamp": "2026-01-19T23:00:00Z",
    "extra": {
        "request_id": "req-123",
        "user_id": 456,
        "ip": "192.168.1.100"
    }
}
```

#### SendEvent - 发送事件

```go
func (s *KafkaService) SendEvent(topic string, eventType string, eventData interface{}) error
```

发送业务事件消息。

**示例**:
```go
eventData := map[string]interface{}{
    "order_id": "ORDER-001",
    "user_id": 789,
    "amount": 999.99,
}
err := service.SendEvent("order-events", "order.created", eventData)
```

**生成的消息格式**:
```json
{
    "event_type": "order.created",
    "event_data": {
        "order_id": "ORDER-001",
        "user_id": 789,
        "amount": 999.99
    },
    "timestamp": "2026-01-19T23:00:00Z"
}
```

#### SendMetric - 发送指标

```go
func (s *KafkaService) SendMetric(topic string, metricName string, metricValue float64, tags map[string]string) error
```

发送监控指标数据。

**示例**:
```go
tags := map[string]string{
    "endpoint": "/api/users",
    "method": "GET",
    "status": "200",
}
err := service.SendMetric(
    "api-metrics",
    "api.response_time",
    123.45,
    tags,
)
```

**生成的消息格式**:
```json
{
    "metric_name": "api.response_time",
    "metric_value": 123.45,
    "tags": {
        "endpoint": "/api/users",
        "method": "GET",
        "status": "200"
    },
    "timestamp": 1737323400
}
```

## 常见使用场景

### 场景 1: 用户行为日志

```go
func logUserAction(userID int, action string) {
    service := kafkaService.NewKafkaService()
    
    logData := map[string]interface{}{
        "user_id": userID,
        "action": action,
        "timestamp": time.Now().Unix(),
        "ip": getClientIP(),
    }
    
    err := service.SendJSON("user-behaviors", logData)
    if err != nil {
        log.Printf("记录用户行为失败: %v", err)
    }
}
```

### 场景 2: 异步任务通知

```go
func notifyTaskComplete(taskID string, result interface{}) {
    service := kafkaService.NewKafkaService()
    
    eventData := map[string]interface{}{
        "task_id": taskID,
        "result": result,
        "completed_at": time.Now(),
    }
    
    // 使用异步发送，不阻塞主流程
    service.SendEventAsync("task-completed", eventData)
}
```

### 场景 3: API 性能监控

```go
func recordAPIMetrics(endpoint string, duration time.Duration, statusCode int) {
    service := kafkaService.NewKafkaService()
    
    tags := map[string]string{
        "endpoint": endpoint,
        "status": strconv.Itoa(statusCode),
    }
    
    // 记录响应时间（毫秒）
    service.SendMetric(
        "api-performance",
        "response_time_ms",
        float64(duration.Milliseconds()),
        tags,
    )
}
```

### 场景 4: 错误日志收集

```go
func logError(err error, context map[string]interface{}) {
    service := kafkaService.NewKafkaService()
    
    service.SendLog(
        "error-logs",
        "ERROR",
        err.Error(),
        context,
    )
}
```

### 场景 5: 订单状态变更事件

```go
func publishOrderEvent(orderID string, status string, details interface{}) {
    service := kafkaService.NewKafkaService()
    
    eventData := map[string]interface{}{
        "order_id": orderID,
        "status": status,
        "details": details,
        "timestamp": time.Now().Unix(),
    }
    
    service.SendEvent("order-status-change", "order."+status, eventData)
}
```

## 最佳实践

### 1. 错误处理

始终检查并处理错误：

```go
err := service.SendMessage(topic, message)
if err != nil {
    log.Printf("Kafka发送失败: %v", err)
    // 根据业务需要决定是否重试或使用备用方案
}
```

### 2. 选择合适的发送方式

- **同步发送** (`SendMessage`): 重要消息，需要确认送达
- **异步发送** (`SendMessageAsync`): 高吞吐量场景，可接受少量消息丢失
- **批量发送** (`SendBatch`): 大量消息需要一次性发送

### 3. 使用上下文控制超时

```go
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()

err := service.SendMessageWithContext(ctx, topic, message)
if err == context.DeadlineExceeded {
    log.Println("发送超时")
}
```

### 4. 合理设置 Topic 命名

建议使用以下命名规范：
- 日志: `{service}-logs`
- 事件: `{domain}-events`
- 指标: `{service}-metrics`
- 任务: `{service}-tasks`

### 5. 消息格式标准化

对于同一类型的消息，保持格式一致：

```go
// 推荐：使用统一的结构
type UserEvent struct {
    EventType string                 `json:"event_type"`
    UserID    int                     `json:"user_id"`
    Data      map[string]interface{}  `json:"data"`
    Timestamp int64                   `json:"timestamp"`
}
```

## 配置说明

### Producer 配置

| 配置项 | 说明 | 推荐值 |
|-------|------|--------|
| required_acks | 0: 不等待<br>1: 等待leader<br>-1: 等待所有副本 | 生产: -1<br>测试: 1 |
| max_retries | 最大重试次数 | 3-5 |
| return_successes | 是否返回成功确认 | true |
| return_errors | 是否返回错误信息 | true |

### SSL/TLS 配置

生产环境建议启用 SSL：

```yaml
ssl:
  enable: true
```

## 故障排查

### 问题 1: Kafka未初始化错误

**错误信息**: `Kafka未初始化，请先调用InitKafka`

**解决方案**:
- 检查配置文件是否存在
- 确认Kafka broker地址正确
- 查看启动日志中的Kafka初始化信息

### 问题 2: 连接超时

**可能原因**:
- Kafka broker 不可达
- 网络配置问题
- 防火墙阻止连接

**解决方案**:
```bash
# 测试Kafka连接
telnet {kafka_host} {kafka_port}
```

### 问题 3: 消息发送失败

**调试步骤**:
1. 检查topic是否存在
2. 确认Kafka服务是否正常
3. 查看详细错误日志
4. 尝试使用同步发送并检查返回的错误

## 性能优化建议

### 1. 批量发送

对于大量消息，使用批量发送可以提高性能：

```go
messages := make([]string, 0, 1000)
for _, item := range items {
    msg, _ := json.Marshal(item)
    messages = append(messages, string(msg))
}
service.SendBatch(topic, messages)
```

### 2. 异步发送

非关键消息使用异步发送：

```go
// 不需要立即确认的场景
service.SendMessageAsync(topic, message)
```

### 3. 合理设置重试次数

根据业务重要性设置重试次数：

```yaml
producer:
  max_retries: 3  # 一般业务
  # max_retries: 5  # 重要业务
```

## 参考资料

- [完整测试报告](../tests/kafka/KAFKA_TEST_REPORT.md)
- [IBM Sarama 文档](https://github.com/IBM/sarama)
- [Apache Kafka 文档](https://kafka.apache.org/documentation/)

## 更新日志

- **2026-01-19**: 初始版本，完成基础功能集成

---

如有问题，请提交 Issue 或联系开发团队。

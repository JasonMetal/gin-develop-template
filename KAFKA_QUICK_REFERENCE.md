# Kafka 快速参考

## 快速开始

```go
import "develop-template/app/service/kafkaService"

service := kafkaService.NewKafkaService()
```

## 常用 API

### 基础发送

```go
// 同步发送
service.SendMessage("topic", "message")

// 异步发送 (高性能)
service.SendMessageAsync("topic", "message")

// 批量发送
service.SendBatch("topic", []string{"msg1", "msg2", "msg3"})
```

### JSON 数据

```go
// 发送 JSON
data := map[string]interface{}{
    "user_id": 123,
    "action": "login",
}
service.SendJSON("user-events", data)

// 异步发送 JSON
service.SendJSONAsync("analytics", data)
```

### 业务场景

```go
// 日志
service.SendLog("app-logs", "ERROR", "错误消息", map[string]interface{}{
    "request_id": "req-123",
})

// 事件
service.SendEvent("order-events", "order.created", orderData)

// 指标
tags := map[string]string{"endpoint": "/api/users"}
service.SendMetric("api-metrics", "response_time", 123.45, tags)
```

### 上下文控制

```go
// 5秒超时
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
service.SendMessageWithContext(ctx, "topic", "message")
```

## 配置文件

`config/{env}/kafka.yml`:

```yaml
brokers:
  - host: "localhost"
    port: 9092

ssl:
  enable: false

producer:
  required_acks: 1      # 0/1/-1
  max_retries: 5
  return_successes: true
  return_errors: true
```

## 配置说明

| 参数 | 说明 | 推荐值 |
|------|------|--------|
| required_acks | 0=不等待, 1=leader, -1=全部 | 生产:-1, 测试:1 |
| max_retries | 重试次数 | 3-5 |
| ssl.enable | 是否启用SSL | 生产:true |

## 最佳实践

### ✅ DO

```go
// 检查错误
if err := service.SendMessage(topic, msg); err != nil {
    log.Printf("发送失败: %v", err)
}

// 使用上下文控制超时
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
defer cancel()

// 重要消息用同步
service.SendMessage("important-topic", msg)

// 统计数据用异步
service.SendMessageAsync("analytics", msg)
```

### ❌ DON'T

```go
// 不检查错误
service.SendMessage(topic, msg)  // 不好

// 无限超时
service.SendMessage(topic, msg)  // 考虑使用 WithContext

// 所有消息都同步发送 (性能差)
// 所有消息都异步发送 (可能丢失)
```

## 常见错误

| 错误 | 原因 | 解决 |
|------|------|------|
| Kafka未初始化 | 配置文件错误 | 检查 kafka.yml |
| 连接超时 | Broker 不可达 | 检查网络和地址 |
| Topic 不存在 | 未创建 Topic | 创建或配置自动创建 |

## 运行测试

```bash
# Windows
.\tests\kafka\run_tests.bat

# Linux/Mac  
bash tests/kafka/run_tests.sh
```

## 文档链接

- 📖 [完整使用指南](docs/KAFKA_INTEGRATION_GUIDE.md)
- 📊 [测试报告](tests/kafka/KAFKA_TEST_REPORT.md)
- 📋 [集成总结](KAFKA_INTEGRATION_SUMMARY.md)
- 🚀 [快速入门](README_KAFKA.md)

## 支持的消息格式

### 日志消息

```json
{
  "level": "ERROR",
  "message": "错误消息",
  "timestamp": "2026-01-19T23:00:00Z",
  "extra": {
    "request_id": "req-123",
    "user_id": 456
  }
}
```

### 事件消息

```json
{
  "event_type": "order.created",
  "event_data": {
    "order_id": "ORDER-001",
    "amount": 999.99
  },
  "timestamp": "2026-01-19T23:00:00Z"
}
```

### 指标消息

```json
{
  "metric_name": "api.response_time",
  "metric_value": 123.45,
  "tags": {
    "endpoint": "/api/users",
    "method": "GET"
  },
  "timestamp": 1737323400
}
```

## Topic 命名规范

```
{service}-logs        # 日志: app-logs, api-logs
{domain}-events       # 事件: order-events, user-events
{service}-metrics     # 指标: api-metrics, db-metrics
{service}-tasks       # 任务: email-tasks, sms-tasks
```

## 性能建议

| 场景 | 推荐方法 | 说明 |
|------|---------|------|
| 重要消息 | SendMessage | 等待确认 |
| 统计数据 | SendMessageAsync | 高吞吐量 |
| 大量消息 | SendBatch | 批量处理 |
| 实时日志 | SendLog | 结构化 |
| 业务事件 | SendEvent | 标准格式 |

## 依赖版本

```
github.com/IBM/sarama v1.46.3
Go 1.24+
Kafka 2.x/3.x
```

## 快速调试

```go
// 启用详细日志
import "log"

// 发送测试消息
service := kafkaService.NewKafkaService()
err := service.SendMessage("test-topic", "Hello Kafka!")
if err != nil {
    log.Printf("发送失败: %v", err)
} else {
    log.Println("发送成功!")
}
```

## 测试状态

✅ Bootstrap: 13/14 通过 (覆盖率 13.4%)  
✅ Service: 11/11 通过 (覆盖率 84.2%)  
✅ 总体: 24/25 通过 (覆盖率 48.8%)

---

**更新时间**: 2026-01-19  
**版本**: v1.0

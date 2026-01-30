# Kafka 异步错误回调修复说明

## 🔧 修复的问题

### 问题描述
代码中变量名不一致导致编译错误：

```go
// ❌ 修复前（第235行）
case err := <-producer.Errors():      // 变量名是 err
    if err != nil {
        log.Printf("Kafka消息发送失败: %v", err)
        
        if asyncErrorHandler != nil {
            asyncErrorHandler(
                errMsg.Msg.Topic,              // ❌ 使用了未定义的 errMsg
                string(errMsg.Msg.Value.(sarama.StringEncoder)),
                errMsg.Err,
            )
        }
    }
```

### 修复方案

```go
// ✅ 修复后
case errMsg := <-producer.Errors():   // 变量名改为 errMsg
    if errMsg != nil {
        log.Printf("Kafka消息发送失败: %v", errMsg.Err)
        
        if asyncErrorHandler != nil {
            asyncErrorHandler(
                errMsg.Msg.Topic,              // ✅ 现在可以正常使用
                string(errMsg.Msg.Value.(sarama.StringEncoder)),
                errMsg.Err,
            )
        }
    }
```

## 📚 errMsg 的来源

### 1. errMsg 是什么？

`errMsg` 来自 Kafka 异步生产者的 Errors channel，类型是 `*sarama.ProducerError`：

```go
// Sarama 库的定义
type ProducerError struct {
    Msg *ProducerMessage  // 发送失败的消息
    Err error             // 具体的错误信息
}

type ProducerMessage struct {
    Topic     string      // 主题名称
    Value     Encoder     // 消息内容（需要类型断言）
    Key       Encoder     // 消息Key
    Partition int32       // 分区
    Offset    int64       // 偏移量
    Timestamp time.Time   // 时间戳
}
```

### 2. 数据流程图

```
你的代码
    ↓
bootstrap.ProducerAsync("topic", "message")
    ↓
producer.Input() <- msg  (写入 Input channel)
    ↓
Kafka 内部异步处理
    ↓
如果发送失败 ❌
    ↓
producer.Errors()  (写入 Errors channel)
    ↓
case errMsg := <-producer.Errors()  (读取错误)
    ↓
调用 asyncErrorHandler(topic, message, err)
    ↓
你的错误处理逻辑
```

### 3. errMsg 包含的信息

| 字段 | 说明 | 类型 | 示例 |
|------|------|------|------|
| `errMsg.Msg.Topic` | 主题名称 | string | "order-created" |
| `errMsg.Msg.Value` | 消息内容 | Encoder | 需要类型断言为 StringEncoder |
| `errMsg.Err` | 错误信息 | error | "connection refused" |
| `errMsg.Msg.Partition` | 分区号 | int32 | 0 |
| `errMsg.Msg.Offset` | 偏移量 | int64 | -1 (发送失败时) |

## 🚀 使用方法

### 步骤1：在项目初始化时设置错误处理器

在 `main.go` 或 `http-server.go` 中：

```go
package main

import (
    "log"
    "develop-template/constant"
    "develop-template/routes"
    "github.com/JasonMetal/submodule-support-go.git/bootstrap"
    "github.com/gin-gonic/gin"
)

func main() {
    bootstrap.SetProjectName(constant.ProjectName)
    bootstrap.Init()
    
    // ✅ 设置 Kafka 异步错误处理器
    bootstrap.SetAsyncErrorHandler(func(topic, message string, err error) {
        log.Printf("❌ Kafka发送失败 - Topic: %s, Error: %v", topic, err)
        
        // 你的错误处理逻辑：
        // 1. 保存到重试队列
        // 2. 发送告警
        // 3. 记录到日志文件
    })
    
    // 启动服务
    middleFun := []gin.HandlerFunc{}
    r := bootstrap.InitWeb(middleFun)
    routes.RegisterRouter(r)
    bootstrap.RunWeb(r, constant.HttpServiceHostPort)
}
```

### 步骤2：使用异步发送

```go
package yourpackage

import (
    "develop-template/app/service/kafkaService"
)

func YourFunction() {
    service := kafkaService.NewKafkaService()
    
    // ✅ 异步发送消息
    err := service.SendMessageAsync("user-behavior", "用户登录")
    if err != nil {
        // 这个 err 只是写入 channel 的错误
        // 真正的发送错误会在 asyncErrorHandler 中处理
        log.Printf("写入失败: %v", err)
    }
    
    // ✅ 异步发送 JSON
    err = service.SendJSONAsync("analytics", map[string]interface{}{
        "user_id": 123,
        "action": "click",
    })
}
```

## 📋 完整示例

### 示例1：基础错误处理

```go
func setupKafkaErrorHandler() {
    bootstrap.SetAsyncErrorHandler(func(topic, message string, err error) {
        log.Printf("Kafka异步发送失败")
        log.Printf("  Topic: %s", topic)
        log.Printf("  Error: %v", err)
        log.Printf("  Message: %s", truncateString(message, 100))
    })
}
```

### 示例2：保存到重试队列

```go
type RetryMessage struct {
    Topic     string
    Message   string
    Error     string
    RetryTime time.Time
    Attempts  int
}

func setupKafkaErrorHandler() {
    bootstrap.SetAsyncErrorHandler(func(topic, message string, err error) {
        log.Printf("Kafka发送失败: Topic=%s, Error=%v", topic, err)
        
        // 保存到数据库重试队列
        db := bootstrap.GetMysqlInstance("default")
        sql := `INSERT INTO kafka_retry_queue 
                (topic, message, error, retry_time, attempts) 
                VALUES (?, ?, ?, ?, ?)`
        
        _, dbErr := db.Exec(sql, topic, message, err.Error(), 
                           time.Now().Add(5*time.Minute), 0)
        if dbErr != nil {
            log.Printf("保存重试队列失败: %v", dbErr)
        }
    })
}
```

### 示例3：区分重要主题

```go
func setupKafkaErrorHandler() {
    bootstrap.SetAsyncErrorHandler(func(topic, message string, err error) {
        log.Printf("Kafka发送失败: Topic=%s, Error=%v", topic, err)
        
        // 重要主题发送告警
        importantTopics := []string{
            "order-created",
            "payment-success",
        }
        
        for _, importantTopic := range importantTopics {
            if topic == importantTopic {
                sendDingTalkAlert(fmt.Sprintf(
                    "【告警】Kafka重要消息发送失败\nTopic: %s\nError: %v",
                    topic, err,
                ))
                break
            }
        }
        
        // 保存到重试队列
        saveToRetryQueue(topic, message, err)
    })
}
```

## ⚠️ 注意事项

### 1. 异步发送的特点

| 特点 | 说明 |
|------|------|
| ✅ 性能高 | 不等待 Kafka 确认，立即返回 |
| ⚠️ 可能丢失 | 如果进程突然退出，缓冲区消息会丢 |
| ⚠️ 错误延迟 | 错误在后台处理，调用方不立即知道 |

### 2. 什么时候用异步？

**✅ 适合异步的场景**：
- 用户行为日志（丢失几条可以接受）
- 页面浏览统计
- 性能监控指标
- 非关键业务数据

**❌ 不适合异步的场景**：
- 订单创建、支付成功（关键业务）
- 账户变更（涉及金钱）
- 审计日志（合规要求）

### 3. 推荐配置

```yaml
# config/prod/kafka.yml
producer:
  required_acks: -1        # 等待所有副本确认（最可靠）
  max_retries: 5           # 失败重试5次
  return_successes: true   # 返回成功确认
  return_errors: true      # 返回错误信息
```

## 📂 相关文件

- **修复的文件**: `submodule/support-go.git/bootstrap/kafka.go`
- **使用示例**: `examples/kafka_error_handler_example.go`
- **详细文档**: `docs/KAFKA_ERROR_CALLBACK_USAGE.md`
- **改进建议**: `docs/KAFKA_ASYNC_IMPROVEMENT.md`

## 🎯 总结

### 修复内容

✅ 修复了 `errMsg` 变量名不一致的问题  
✅ 现在异步错误可以被正确捕获和处理  
✅ 提供了完整的使用示例和文档  

### 使用建议

1. **一定要设置错误处理器**：
   ```go
   bootstrap.SetAsyncErrorHandler(func(topic, message string, err error) {
       // 你的错误处理逻辑
   })
   ```

2. **根据业务重要性选择发送方式**：
   - 🔴 关键业务 → `SendMessage` (同步)
   - 🟡 一般业务 → `SendBatch` (批量)
   - 🟢 非关键业务 → `SendMessageAsync` (异步)

3. **建立重试机制**：
   - 保存失败消息到数据库
   - 定时任务重新发送
   - 设置重试次数上限

4. **监控告警**：
   - 重要主题失败时发送告警
   - 统计失败次数
   - 定期检查重试队列

---

**更新时间**: 2026-01-19  
**修复状态**: ✅ 已完成

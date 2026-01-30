# Kafka 异步发送改进建议

## 当前问题

当前异步实现的主要问题：

1. **错误处理不完善**：异步错误只打印日志，调用方不知道失败
2. **Goroutine 泄漏**：错误监听的 goroutine 没有退出机制
3. **消息可能丢失**：进程退出时缓冲区消息会丢失

## 改进方案

### 方案1：添加错误回调（推荐）

```go
// 在 bootstrap/kafka.go 中添加

type AsyncErrorHandler func(topic string, message string, err error)

var asyncErrorHandler AsyncErrorHandler

// SetAsyncErrorHandler 设置异步错误处理器
func SetAsyncErrorHandler(handler AsyncErrorHandler) {
    asyncErrorHandler = handler
}

// 修改 getAsyncProducer
func (km *KafkaManager) getAsyncProducer() (sarama.AsyncProducer, error) {
    // ... 创建 producer ...

    // 启动监听goroutine
    go func() {
        for {
            select {
            case success := <-producer.Successes():
                if success != nil {
                    log.Printf("Kafka消息发送成功: Topic=%s, Partition=%d, Offset=%d",
                        success.Topic, success.Partition, success.Offset)
                }
            case errMsg := <-producer.Errors():
                if errMsg != nil {
                    log.Printf("Kafka消息发送失败: %v", errMsg.Err)
                    
                    // ✅ 调用错误处理器
                    if asyncErrorHandler != nil {
                        asyncErrorHandler(
                            errMsg.Msg.Topic,
                            string(errMsg.Msg.Value.(sarama.StringEncoder)),
                            errMsg.Err,
                        )
                    }
                }
            }
        }
    }()

    return producer, nil
}
```

**使用示例**：

```go
// 在 main.go 或初始化代码中设置
bootstrap.SetAsyncErrorHandler(func(topic, message string, err error) {
    // 将失败的消息保存到数据库或重试队列
    log.Printf("异步发送失败 Topic:%s, Error:%v", topic, err)
    
    // 保存到重试队列
    retryQueue.Push(RetryMessage{
        Topic:   topic,
        Message: message,
        Error:   err.Error(),
        Time:    time.Now(),
    })
})
```

### 方案2：添加发送确认机制

```go
// 为重要的异步消息添加确认
type AsyncResult struct {
    Topic     string
    Success   bool
    Error     error
    Partition int32
    Offset    int64
}

// ProducerAsyncWithCallback 带回调的异步发送
func ProducerAsyncWithCallback(topic string, message string, callback func(result AsyncResult)) error {
    if kafkaManager == nil {
        return fmt.Errorf("Kafka未初始化")
    }

    producer, err := kafkaManager.getAsyncProducer()
    if err != nil {
        return err
    }

    // 为这条消息创建唯一标识
    msgID := fmt.Sprintf("%s-%d", topic, time.Now().UnixNano())
    
    msg := &sarama.ProducerMessage{
        Topic: topic,
        Value: sarama.StringEncoder(message),
        Metadata: msgID, // 使用 Metadata 传递标识
    }

    // 启动监听这条消息结果的 goroutine
    go func() {
        timeout := time.After(10 * time.Second)
        for {
            select {
            case success := <-producer.Successes():
                if success.Metadata == msgID {
                    callback(AsyncResult{
                        Topic:     success.Topic,
                        Success:   true,
                        Partition: success.Partition,
                        Offset:    success.Offset,
                    })
                    return
                }
            case errMsg := <-producer.Errors():
                if errMsg.Msg.Metadata == msgID {
                    callback(AsyncResult{
                        Topic:   errMsg.Msg.Topic,
                        Success: false,
                        Error:   errMsg.Err,
                    })
                    return
                }
            case <-timeout:
                callback(AsyncResult{
                    Topic:   topic,
                    Success: false,
                    Error:   fmt.Errorf("等待确认超时"),
                })
                return
            }
        }
    }()

    producer.Input() <- msg
    return nil
}
```

**使用示例**：

```go
bootstrap.ProducerAsyncWithCallback("important-events", eventData, func(result AsyncResult) {
    if !result.Success {
        log.Printf("发送失败: %v", result.Error)
        // 重试或保存到数据库
        saveToRetryQueue(eventData)
    } else {
        log.Printf("发送成功: Offset=%d", result.Offset)
    }
})
```

### 方案3：优雅关闭改进

```go
// 改进 CloseKafka，确保异步消息全部发送完成
func CloseKafka() error {
    if kafkaManager == nil {
        return nil
    }

    var errs []error

    // 关闭异步生产者时，等待所有消息发送完成
    if kafkaManager.asyncProducer != nil {
        log.Println("等待异步消息发送完成...")
        
        // 关闭 Input channel，不再接收新消息
        kafkaManager.asyncProducer.AsyncClose()
        
        // 等待所有消息处理完成（最多等待10秒）
        timeout := time.After(10 * time.Second)
        done := make(chan bool)
        
        go func() {
            // 等待 Successes 和 Errors channel 都关闭
            for range kafkaManager.asyncProducer.Successes() {
            }
            for range kafkaManager.asyncProducer.Errors() {
            }
            done <- true
        }()
        
        select {
        case <-done:
            log.Println("所有异步消息已发送完成")
        case <-timeout:
            log.Println("异步消息发送超时，强制关闭")
        }
    }

    if kafkaManager.syncProducer != nil {
        if err := kafkaManager.syncProducer.Close(); err != nil {
            errs = append(errs, fmt.Errorf("关闭同步生产者失败: %v", err))
        }
    }

    if len(errs) > 0 {
        return fmt.Errorf("关闭Kafka连接时出现错误: %v", errs)
    }

    log.Println("Kafka连接已关闭")
    return nil
}
```

## 业务场景选择指南

### 使用同步发送的场景

| 场景 | 原因 |
|------|------|
| 订单创建/支付 | 不能丢失，必须确认 |
| 账户变更 | 涉及金钱，必须可靠 |
| 关键业务事件 | 需要立即知道是否成功 |
| 审计日志 | 合规要求，不能丢失 |

```go
// ✅ 示例：订单支付
func ProcessPayment(order Order) error {
    // 同步发送支付事件
    if err := kafkaService.SendMessage("payment-success", orderJSON); err != nil {
        log.Printf("支付事件发送失败: %v", err)
        // 回滚支付或重试
        return err
    }
    return nil
}
```

### 使用异步发送的场景

| 场景 | 原因 |
|------|------|
| 用户行为日志 | 丢失几条可以接受 |
| 页面浏览统计 | 非关键数据 |
| 性能监控指标 | 高频数据，性能优先 |
| 搜索关键词统计 | 聚合数据，少量丢失影响小 |

```go
// ✅ 示例：用户行为追踪
func TrackUserBehavior(behavior UserBehavior) {
    // 异步发送，不阻塞主流程
    kafkaService.SendMessageAsync("user-behavior", behaviorJSON)
    // 不需要检查错误，失败了也不影响业务
}
```

### 使用批量发送的场景

| 场景 | 原因 |
|------|------|
| 批量导入数据 | 大量数据，需要高性能 |
| 定时任务数据同步 | 批量处理，性能优先 |
| 日志归档 | 批量写入，减少网络开销 |

```go
// ✅ 示例：批量同步数据
func SyncDataBatch(records []DataRecord) error {
    messages := make([]string, len(records))
    for i, record := range records {
        messages[i] = record.ToJSON()
    }
    
    // 批量同步发送
    return kafkaService.SendBatch("data-sync", messages)
}
```

## 推荐的配置

### 生产环境配置

```yaml
# config/prod/kafka.yml
brokers:
  - host: "kafka1.prod.com"
    port: 9092
  - host: "kafka2.prod.com"
    port: 9092

ssl:
  enable: true

producer:
  required_acks: -1        # ✅ 等待所有副本确认（最可靠）
  max_retries: 5           # ✅ 失败重试5次
  return_successes: true
  return_errors: true

version: "3.7.2"
```

### 测试环境配置

```yaml
# config/test/kafka.yml
producer:
  required_acks: 1         # 只等待leader确认（平衡性能和可靠性）
  max_retries: 3
  return_successes: true
  return_errors: true
```

## 最佳实践总结

### ✅ DO（推荐做法）

1. **关键业务用同步**：
   ```go
   // 订单、支付、账户变更等
   err := service.SendMessage("order-created", data)
   if err != nil {
       // 立即处理错误
   }
   ```

2. **非关键数据用异步**：
   ```go
   // 日志、统计、监控等
   service.SendMessageAsync("user-behavior", data)
   ```

3. **批量数据用批量发送**：
   ```go
   service.SendBatch("logs", messages)
   ```

4. **设置异步错误处理器**：
   ```go
   bootstrap.SetAsyncErrorHandler(func(topic, msg string, err error) {
       saveToRetryQueue(topic, msg, err)
   })
   ```

### ❌ DON'T（不推荐做法）

1. **不要混用**：
   ```go
   // ❌ 同一个业务流程混用同步和异步
   service.SendMessage("order-created", data1)      // 同步
   service.SendMessageAsync("order-detail", data2)  // 异步
   // 可能导致 order-detail 先到，order-created 后到
   ```

2. **不要忽略异步错误**：
   ```go
   // ❌ 重要数据用异步，但不处理错误
   service.SendMessageAsync("payment-success", data)
   // 如果发送失败，你永远不知道
   ```

3. **不要在循环中同步发送**：
   ```go
   // ❌ 性能很差
   for _, item := range items {
       service.SendMessage("topic", item)  // 每次都等待
   }
   
   // ✅ 改用批量
   messages := convertToMessages(items)
   service.SendBatch("topic", messages)
   ```

## 总结

**你的问题**："异步会不会有问题？"

**答案**：
- 异步本身没问题，但要看 **使用场景**
- 当前实现的异步 **错误处理不完善**，建议改进
- 建议根据 **业务重要性** 选择：
  - 🔴 关键业务 → 同步
  - 🟡 一般业务 → 批量
  - 🟢 非关键业务 → 异步

**建议**：
1. 保持当前实现，**根据场景选择** 同步/异步
2. 为重要的异步消息添加 **错误回调机制**
3. 在 gracefulShutdown 中 **等待异步消息发送完成**
4. 为异步失败的消息建立 **重试队列**

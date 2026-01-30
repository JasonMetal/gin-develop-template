# Kafka 异步错误回调使用指南

## errMsg 的来源

### 1. errMsg 是什么？

`errMsg` 是从 Kafka 异步生产者的错误 channel 中获取的，类型是 `*sarama.ProducerError`：

```go
// sarama 库的定义
type ProducerError struct {
    Msg *ProducerMessage  // 发送失败的消息
    Err error             // 错误信息
}

// ProducerMessage 包含
type ProducerMessage struct {
    Topic     string        // 主题
    Key       Encoder       // 消息Key (可选)
    Value     Encoder       // 消息内容
    Headers   []RecordHeader // 消息头 (可选)
    Metadata  interface{}   // 元数据 (可选)
    Offset    int64         // 偏移量
    Partition int32         // 分区
    Timestamp time.Time     // 时间戳
}
```

### 2. 数据流向

```
你的代码发送消息
    ↓
producer.Input() <- msg  (写入 Input channel)
    ↓
Kafka 内部处理
    ↓
发送失败
    ↓
producer.Errors() (写入 Errors channel)
    ↓
errMsg := <-producer.Errors()  (从 Errors channel 读取)
    ↓
调用你的错误处理器
```

### 3. 修复后的完整代码

```go
// bootstrap/kafka.go

// 异步错误处理器类型定义
type AsyncErrorHandler func(topic string, message string, err error)

var asyncErrorHandler AsyncErrorHandler

// SetAsyncErrorHandler 设置异步错误处理器
func SetAsyncErrorHandler(handler AsyncErrorHandler) {
	asyncErrorHandler = handler
}

// getAsyncProducer 中的错误处理
go func() {
    for {
        select {
        case success := <-producer.Successes():
            if success != nil {
                log.Printf("Kafka消息发送成功: Topic=%s, Partition=%d, Offset=%d",
                    success.Topic, success.Partition, success.Offset)
            }
            
        // ✅ errMsg 从这里来
        case errMsg := <-producer.Errors():
            if errMsg != nil {
                // errMsg.Err 是错误信息
                log.Printf("Kafka消息发送失败: %v", errMsg.Err)

                // 调用用户设置的错误处理器
                if asyncErrorHandler != nil {
                    asyncErrorHandler(
                        errMsg.Msg.Topic,                              // 主题名
                        string(errMsg.Msg.Value.(sarama.StringEncoder)), // 消息内容
                        errMsg.Err,                                    // 错误信息
                    )
                }
            }
        }
    }
}()
```

## 使用示例

### 示例1：基础使用（在 main.go 中）

```go
package main

import (
    "log"
    "develop-template/constant"
    "github.com/JasonMetal/submodule-support-go.git/bootstrap"
)

func main() {
    bootstrap.SetProjectName(constant.ProjectName)
    bootstrap.Init()
    
    // ✅ 设置异步错误处理器
    bootstrap.SetAsyncErrorHandler(func(topic, message string, err error) {
        log.Printf("【Kafka异步发送失败】Topic: %s, Error: %v", topic, err)
        log.Printf("【失败消息内容】%s", message)
        
        // 可以在这里做进一步处理：
        // 1. 保存到数据库
        // 2. 写入重试队列
        // 3. 发送告警
    })
    
    // 启动服务...
    middleFun := []gin.HandlerFunc{}
    r := bootstrap.InitWeb(middleFun)
    router.RegisterRouter(r)
    bootstrap.RunWeb(r, constant.HttpServiceHostPort)
}
```

### 示例2：保存到重试队列

```go
// 定义重试队列结构
type RetryMessage struct {
    Topic     string    `json:"topic"`
    Message   string    `json:"message"`
    Error     string    `json:"error"`
    RetryTime time.Time `json:"retry_time"`
    Attempts  int       `json:"attempts"`
}

// 在初始化时设置错误处理器
func setupKafkaErrorHandler() {
    bootstrap.SetAsyncErrorHandler(func(topic, message string, err error) {
        log.Printf("Kafka异步发送失败 - Topic: %s, Error: %v", topic, err)
        
        // 保存到重试队列（数据库或 Redis）
        retryMsg := RetryMessage{
            Topic:     topic,
            Message:   message,
            Error:     err.Error(),
            RetryTime: time.Now().Add(5 * time.Minute), // 5分钟后重试
            Attempts:  0,
        }
        
        // 保存到数据库
        if err := saveToRetryQueue(retryMsg); err != nil {
            log.Printf("保存重试队列失败: %v", err)
        }
    })
}

// 保存到 MySQL
func saveToRetryQueue(msg RetryMessage) error {
    db := bootstrap.GetMysqlInstance("default")
    sql := `INSERT INTO kafka_retry_queue 
            (topic, message, error, retry_time, attempts) 
            VALUES (?, ?, ?, ?, ?)`
    
    _, err := db.Exec(sql, msg.Topic, msg.Message, msg.Error, msg.RetryTime, msg.Attempts)
    return err
}

// 定时重试任务
func retryFailedMessages() {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        // 查询需要重试的消息
        messages := queryRetryMessages()
        
        service := kafkaService.NewKafkaService()
        for _, msg := range messages {
            // 重试发送
            if err := service.SendMessage(msg.Topic, msg.Message); err != nil {
                // 重试失败，增加重试次数
                updateRetryAttempts(msg.ID, msg.Attempts+1)
            } else {
                // 重试成功，删除记录
                deleteRetryMessage(msg.ID)
            }
        }
    }
}
```

### 示例3：发送告警

```go
func setupKafkaErrorHandler() {
    bootstrap.SetAsyncErrorHandler(func(topic, message string, err error) {
        log.Printf("Kafka异步发送失败 - Topic: %s, Error: %v", topic, err)
        
        // 如果是重要主题，发送告警
        importantTopics := []string{"order-created", "payment-success"}
        for _, importantTopic := range importantTopics {
            if topic == importantTopic {
                // 发送告警（邮件、短信、钉钉等）
                sendAlert(fmt.Sprintf(
                    "Kafka重要消息发送失败！\nTopic: %s\nError: %v\nMessage: %s",
                    topic, err, truncateMessage(message, 100),
                ))
                break
            }
        }
        
        // 保存到重试队列
        saveToRetryQueue(RetryMessage{
            Topic:     topic,
            Message:   message,
            Error:     err.Error(),
            RetryTime: time.Now().Add(5 * time.Minute),
            Attempts:  0,
        })
    })
}

func sendAlert(message string) {
    // 发送钉钉告警
    // 或发送邮件
    // 或发送短信
    log.Printf("【告警】%s", message)
}

func truncateMessage(msg string, maxLen int) string {
    if len(msg) <= maxLen {
        return msg
    }
    return msg[:maxLen] + "..."
}
```

### 示例4：区分不同类型的错误

```go
func setupKafkaErrorHandler() {
    bootstrap.SetAsyncErrorHandler(func(topic, message string, err error) {
        log.Printf("Kafka异步发送失败 - Topic: %s, Error: %v", topic, err)
        
        // 根据错误类型做不同处理
        errStr := err.Error()
        
        switch {
        case strings.Contains(errStr, "connection refused"):
            // Kafka 服务不可用
            log.Printf("【严重】Kafka服务不可用: %v", err)
            sendUrgentAlert("Kafka服务不可用", err.Error())
            saveToRetryQueue(topic, message, err, 10*time.Minute) // 10分钟后重试
            
        case strings.Contains(errStr, "timeout"):
            // 超时错误
            log.Printf("【警告】Kafka发送超时: %v", err)
            saveToRetryQueue(topic, message, err, 1*time.Minute) // 1分钟后重试
            
        case strings.Contains(errStr, "message too large"):
            // 消息太大
            log.Printf("【错误】消息超过大小限制: Topic=%s, Size=%d", topic, len(message))
            // 消息太大，不重试，记录日志
            saveToErrorLog(topic, message, err)
            
        default:
            // 其他错误
            log.Printf("【错误】Kafka发送失败: %v", err)
            saveToRetryQueue(topic, message, err, 5*time.Minute) // 5分钟后重试
        }
    })
}
```

### 示例5：统计失败次数

```go
// 错误统计
type ErrorStats struct {
    sync.RWMutex
    Count    map[string]int // topic -> 失败次数
    LastTime map[string]time.Time
}

var errorStats = &ErrorStats{
    Count:    make(map[string]int),
    LastTime: make(map[string]time.Time),
}

func setupKafkaErrorHandler() {
    bootstrap.SetAsyncErrorHandler(func(topic, message string, err error) {
        log.Printf("Kafka异步发送失败 - Topic: %s, Error: %v", topic, err)
        
        // 统计失败次数
        errorStats.Lock()
        errorStats.Count[topic]++
        errorStats.LastTime[topic] = time.Now()
        failCount := errorStats.Count[topic]
        errorStats.Unlock()
        
        // 如果失败次数超过阈值，发送告警
        if failCount >= 10 {
            sendAlert(fmt.Sprintf(
                "Topic %s 在最近时间内失败了 %d 次，请检查！",
                topic, failCount,
            ))
            
            // 重置计数
            errorStats.Lock()
            errorStats.Count[topic] = 0
            errorStats.Unlock()
        }
        
        // 保存到重试队列
        saveToRetryQueue(topic, message, err, 5*time.Minute)
    })
    
    // 定期重置统计（每小时）
    go func() {
        ticker := time.NewTicker(1 * time.Hour)
        defer ticker.Stop()
        
        for range ticker.C {
            errorStats.Lock()
            errorStats.Count = make(map[string]int)
            errorStats.LastTime = make(map[string]time.Time)
            errorStats.Unlock()
            log.Println("Kafka错误统计已重置")
        }
    }()
}
```

## 数据库表结构示例

如果要保存重试队列到数据库：

```sql
CREATE TABLE `kafka_retry_queue` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `topic` varchar(255) NOT NULL COMMENT '主题',
  `message` text NOT NULL COMMENT '消息内容',
  `error` varchar(500) DEFAULT NULL COMMENT '错误信息',
  `retry_time` datetime NOT NULL COMMENT '下次重试时间',
  `attempts` int(11) DEFAULT '0' COMMENT '重试次数',
  `status` tinyint(4) DEFAULT '0' COMMENT '状态: 0-待重试, 1-重试成功, 2-重试失败',
  `created_at` datetime DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_retry_time` (`retry_time`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='Kafka重试队列';
```

## 完整的 main.go 示例

```go
package main

import (
    "log"
    "time"
    "develop-template/constant"
    "develop-template/routes"
    "github.com/JasonMetal/submodule-support-go.git/bootstrap"
    "github.com/gin-gonic/gin"
)

func main() {
    // 设置项目名称
    bootstrap.SetProjectName(constant.ProjectName)
    
    // 初始化（会自动初始化 Kafka）
    bootstrap.Init()
    
    // ✅ 设置 Kafka 异步错误处理器
    setupKafkaErrorHandler()
    
    // 启动重试任务
    go retryFailedMessages()
    
    // 启动 Web 服务
    middleFun := []gin.HandlerFunc{
        // middleware.CheckUserAuth(),
    }
    r := bootstrap.InitWeb(middleFun)
    routes.RegisterRouter(r)
    bootstrap.RunWeb(r, constant.HttpServiceHostPort)
}

// 设置 Kafka 错误处理器
func setupKafkaErrorHandler() {
    bootstrap.SetAsyncErrorHandler(func(topic, message string, err error) {
        log.Printf("【Kafka异步发送失败】Topic: %s, Error: %v", topic, err)
        
        // 保存到重试队列
        retryMsg := RetryMessage{
            Topic:     topic,
            Message:   message,
            Error:     err.Error(),
            RetryTime: time.Now().Add(5 * time.Minute),
            Attempts:  0,
        }
        
        if err := saveToRetryQueue(retryMsg); err != nil {
            log.Printf("保存重试队列失败: %v", err)
        }
        
        // 如果是重要主题，发送告警
        if isImportantTopic(topic) {
            sendAlert(fmt.Sprintf(
                "Kafka重要消息发送失败！\nTopic: %s\nError: %v",
                topic, err,
            ))
        }
    })
    
    log.Println("Kafka异步错误处理器已设置")
}

// 判断是否是重要主题
func isImportantTopic(topic string) bool {
    importantTopics := []string{
        "order-created",
        "payment-success",
        "user-registered",
    }
    
    for _, t := range importantTopics {
        if topic == t {
            return true
        }
    }
    return false
}

// 重试消息结构
type RetryMessage struct {
    ID        int64     `json:"id"`
    Topic     string    `json:"topic"`
    Message   string    `json:"message"`
    Error     string    `json:"error"`
    RetryTime time.Time `json:"retry_time"`
    Attempts  int       `json:"attempts"`
}

// 保存到重试队列
func saveToRetryQueue(msg RetryMessage) error {
    db := bootstrap.GetMysqlInstance("default")
    sql := `INSERT INTO kafka_retry_queue 
            (topic, message, error, retry_time, attempts) 
            VALUES (?, ?, ?, ?, ?)`
    
    _, err := db.Exec(sql, msg.Topic, msg.Message, msg.Error, msg.RetryTime, msg.Attempts)
    return err
}

// 查询需要重试的消息
func queryRetryMessages() []RetryMessage {
    db := bootstrap.GetMysqlInstance("default")
    sql := `SELECT id, topic, message, error, retry_time, attempts 
            FROM kafka_retry_queue 
            WHERE status = 0 AND retry_time <= ? AND attempts < 5
            LIMIT 100`
    
    rows, err := db.Query(sql, time.Now())
    if err != nil {
        log.Printf("查询重试消息失败: %v", err)
        return nil
    }
    defer rows.Close()
    
    var messages []RetryMessage
    for rows.Next() {
        var msg RetryMessage
        if err := rows.Scan(&msg.ID, &msg.Topic, &msg.Message, 
                            &msg.Error, &msg.RetryTime, &msg.Attempts); err != nil {
            log.Printf("扫描重试消息失败: %v", err)
            continue
        }
        messages = append(messages, msg)
    }
    
    return messages
}

// 更新重试次数
func updateRetryAttempts(id int64, attempts int) error {
    db := bootstrap.GetMysqlInstance("default")
    
    // 如果重试次数超过5次，标记为失败
    if attempts >= 5 {
        sql := `UPDATE kafka_retry_queue SET attempts = ?, status = 2 WHERE id = ?`
        _, err := db.Exec(sql, attempts, id)
        return err
    }
    
    sql := `UPDATE kafka_retry_queue 
            SET attempts = ?, retry_time = ? 
            WHERE id = ?`
    
    nextRetryTime := time.Now().Add(time.Duration(attempts+1) * 5 * time.Minute)
    _, err := db.Exec(sql, attempts, nextRetryTime, id)
    return err
}

// 删除重试消息
func deleteRetryMessage(id int64) error {
    db := bootstrap.GetMysqlInstance("default")
    sql := `UPDATE kafka_retry_queue SET status = 1 WHERE id = ?`
    _, err := db.Exec(sql, id)
    return err
}

// 重试失败的消息
func retryFailedMessages() {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()
    
    for range ticker.C {
        messages := queryRetryMessages()
        if len(messages) == 0 {
            continue
        }
        
        log.Printf("开始重试 %d 条失败消息", len(messages))
        
        service := kafkaService.NewKafkaService()
        successCount := 0
        failCount := 0
        
        for _, msg := range messages {
            // 重试发送（使用同步发送确保可靠性）
            if err := service.SendMessage(msg.Topic, msg.Message); err != nil {
                // 重试失败，增加重试次数
                updateRetryAttempts(msg.ID, msg.Attempts+1)
                failCount++
                log.Printf("消息重试失败 [%d/%d]: Topic=%s, Error=%v", 
                          msg.Attempts+1, 5, msg.Topic, err)
            } else {
                // 重试成功，删除记录
                deleteRetryMessage(msg.ID)
                successCount++
                log.Printf("消息重试成功: Topic=%s", msg.Topic)
            }
        }
        
        log.Printf("重试完成: 成功=%d, 失败=%d", successCount, failCount)
    }
}

// 发送告警
func sendAlert(message string) {
    // 这里可以集成：
    // 1. 钉钉机器人
    // 2. 邮件通知
    // 3. 短信通知
    // 4. 企业微信
    log.Printf("【告警】%s", message)
}
```

## 总结

### errMsg 包含的信息

| 字段 | 说明 | 示例 |
|------|------|------|
| `errMsg.Msg.Topic` | 消息主题 | "order-created" |
| `errMsg.Msg.Value` | 消息内容 | "订单数据JSON" |
| `errMsg.Err` | 错误信息 | "connection refused" |
| `errMsg.Msg.Partition` | 分区 | 0 |
| `errMsg.Msg.Offset` | 偏移量 | 1234 |

### 使用建议

1. **一定要设置错误处理器**，否则异步错误只会打日志
2. **重要消息要保存到重试队列**，确保不丢失
3. **设置重试次数上限**（比如5次），避免无限重试
4. **关键主题发送告警**，及时发现问题
5. **定期清理重试队列**，避免数据堆积

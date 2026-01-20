# RabbitMQ 查询功能实现文档

## 📝 实现概述

本文档说明了为 RabbitMQ 集成添加的查询功能，类似于 Kafka 的消息查询能力。

---

## 🎯 功能清单

### ✅ 已实现的功能

1. **获取队列信息** - 查询队列状态（消息数、消费者数等）
2. **查看队列消息（Peek）** - 查看消息但不删除
3. **消费队列消息（Consume）** - 查看并删除消息
4. **清空队列** - 删除队列中的所有消息
5. **删除队列** - 删除整个队列

---

## 📁 代码结构

### 1. Bootstrap 层 (`submodule/support-go.git/bootstrap/rabbitmq.go`)

添加了以下核心方法：

```go
// 数据结构
type QueueInfo struct {
    Name       string `json:"name"`
    Messages   int    `json:"messages"`
    Consumers  int    `json:"consumers"`
    Durable    bool   `json:"durable"`
    AutoDelete bool   `json:"auto_delete"`
    Exclusive  bool   `json:"exclusive"`
}

type QueueMessage struct {
    Body          string            `json:"body"`
    ContentType   string            `json:"content_type"`
    DeliveryMode  uint8             `json:"delivery_mode"`
    Priority      uint8             `json:"priority"`
    CorrelationId string            `json:"correlation_id"`
    ReplyTo       string            `json:"reply_to"`
    Expiration    string            `json:"expiration"`
    MessageId     string            `json:"message_id"`
    Timestamp     time.Time         `json:"timestamp"`
    Type          string            `json:"type"`
    UserId        string            `json:"user_id"`
    AppId         string            `json:"app_id"`
    Headers       map[string]string `json:"headers"`
    DeliveryTag   uint64            `json:"delivery_tag"`
    Redelivered   bool              `json:"redelivered"`
    Exchange      string            `json:"exchange"`
    RoutingKey    string            `json:"routing_key"`
}

// 查询方法
func GetQueueInfo(queueName string) (*QueueInfo, error)
func DeclareAndGetQueueInfo(queueName string) (*QueueInfo, error)
func PeekMessages(queueName string, count int) ([]QueueMessage, error)
func ConsumeMessages(queueName string, count int) ([]QueueMessage, error)
func PurgeQueue(queueName string) (int, error)
func DeleteQueue(queueName string, ifUnused, ifEmpty bool) (int, error)
```

**实现细节**：
- 使用 `rabbitmqManager.channel.Get()` 获取消息
- Peek 模式使用 `autoAck=false` + `Nack(requeue=true)` 实现消息重新入队
- Consume 模式使用 `autoAck=true` 实现自动确认并删除
- 使用读锁保护并发访问

### 2. Service 层 (`app/service/rabbitmqService/rabbitmqService.go`)

添加了服务层封装：

```go
func (s *RabbitMQService) GetQueueInfo(queueName string) (*bootstrap.QueueInfo, error)
func (s *RabbitMQService) DeclareAndGetQueueInfo(queueName string) (*bootstrap.QueueInfo, error)
func (s *RabbitMQService) PeekMessages(queueName string, count int) ([]bootstrap.QueueMessage, error)
func (s *RabbitMQService) ConsumeMessages(queueName string, count int) ([]bootstrap.QueueMessage, error)
func (s *RabbitMQService) PurgeQueue(queueName string) (int, error)
func (s *RabbitMQService) DeleteQueue(queueName string, ifUnused, ifEmpty bool) (int, error)
```

### 3. 请求/响应实体

#### 请求实体 (`app/entity/req/rabbitmqReq.go`)

```go
type GetQueueInfoReq struct {
    Queue string `form:"queue" json:"queue" binding:"required"`
}

type PeekMessagesReq struct {
    Queue string `form:"queue" json:"queue" binding:"required"`
    Limit int    `form:"limit" json:"limit"`
}

type ConsumeMessagesReq struct {
    Queue string `form:"queue" json:"queue" binding:"required"`
    Limit int    `form:"limit" json:"limit"`
}

type PurgeQueueReq struct {
    Queue string `json:"queue" binding:"required"`
}

type DeleteQueueReq struct {
    Queue    string `json:"queue" binding:"required"`
    IfUnused bool   `json:"if_unused"`
    IfEmpty  bool   `json:"if_empty"`
}
```

#### 响应实体 (`app/entity/resp/rabbitmqResp.go`)

```go
type QueueInfoResp struct {
    Name       string `json:"name"`
    Messages   int    `json:"messages"`
    Consumers  int    `json:"consumers"`
    Durable    bool   `json:"durable"`
    AutoDelete bool   `json:"auto_delete"`
    Exclusive  bool   `json:"exclusive"`
}

type QueueMessageResp struct {
    Body          string            `json:"body"`
    ContentType   string            `json:"content_type"`
    DeliveryMode  uint8             `json:"delivery_mode"`
    Priority      uint8             `json:"priority"`
    CorrelationId string            `json:"correlation_id"`
    ReplyTo       string            `json:"reply_to"`
    Expiration    string            `json:"expiration"`
    MessageId     string            `json:"message_id"`
    Timestamp     string            `json:"timestamp"`
    Type          string            `json:"type"`
    UserId        string            `json:"user_id"`
    AppId         string            `json:"app_id"`
    Headers       map[string]string `json:"headers"`
    DeliveryTag   uint64            `json:"delivery_tag"`
    Redelivered   bool              `json:"redelivered"`
    Exchange      string            `json:"exchange"`
    RoutingKey    string            `json:"routing_key"`
}

type PeekMessagesResp struct {
    Queue    string             `json:"queue"`
    Total    int                `json:"total"`
    Messages []QueueMessageResp `json:"messages"`
}

type ConsumeMessagesResp struct {
    Queue    string             `json:"queue"`
    Total    int                `json:"total"`
    Messages []QueueMessageResp `json:"messages"`
}

type PurgeQueueResp struct {
    Queue        string `json:"queue"`
    DeletedCount int    `json:"deleted_count"`
    Success      bool   `json:"success"`
}

type DeleteQueueResp struct {
    Queue        string `json:"queue"`
    DeletedCount int    `json:"deleted_count"`
    Success      bool   `json:"success"`
}
```

### 4. Controller 层 (`app/http/controller/api/rabbitmqController/rabbitmq.go`)

添加了 HTTP 处理方法：

```go
func (c *controller) GetQueueInfo()     // GET /rabbitmq/queue/info
func (c *controller) PeekMessages()     // GET /rabbitmq/queue/peek
func (c *controller) ConsumeMessages()  // GET /rabbitmq/queue/consume
func (c *controller) PurgeQueue()       // POST /rabbitmq/queue/purge
func (c *controller) DeleteQueue()      // POST /rabbitmq/queue/delete
```

**特点**：
- 参数验证和默认值处理
- 统一的错误处理
- 详细的日志记录
- Swagger 文档注释

### 5. 路由注册 (`routes/api/rabbitmqRouter/rabbitmq.go`)

```go
// 获取队列信息
router.GET("/rabbitmq/queue/info", func(ctx *gin.Context) {
    rabbitmqController.NewController(ctx).GetQueueInfo()
})

// 查看队列消息（不消费）
router.GET("/rabbitmq/queue/peek", func(ctx *gin.Context) {
    rabbitmqController.NewController(ctx).PeekMessages()
})

// 消费队列消息（会删除）
router.GET("/rabbitmq/queue/consume", func(ctx *gin.Context) {
    rabbitmqController.NewController(ctx).ConsumeMessages()
})

// 清空队列
router.POST("/rabbitmq/queue/purge", func(ctx *gin.Context) {
    rabbitmqController.NewController(ctx).PurgeQueue()
})

// 删除队列
router.POST("/rabbitmq/queue/delete", func(ctx *gin.Context) {
    rabbitmqController.NewController(ctx).DeleteQueue()
})
```

---

## 🔧 技术实现要点

### Peek 模式的实现

**关键代码**：
```go
for i := 0; i < count; i++ {
    // 1. 获取消息但不自动确认
    delivery, ok, err := rabbitmqManager.channel.Get(queueName, false)
    if err != nil {
        return messages, fmt.Errorf("获取消息失败: %v", err)
    }
    
    // 2. 没有更多消息则退出
    if !ok {
        break
    }
    
    // 3. 处理消息内容
    msg := QueueMessage{
        Body: string(delivery.Body),
        // ... 其他字段
    }
    messages = append(messages, msg)
    
    // 4. 拒绝消息并重新入队（关键步骤）
    err = rabbitmqManager.channel.Nack(delivery.DeliveryTag, false, true)
    if err != nil {
        log.Printf("拒绝消息失败: %v", err)
    }
}
```

**原理**：
1. `Get(queueName, false)` - 获取消息但不自动确认
2. `Nack(deliveryTag, false, true)` - 拒绝消息并重新入队
   - 第1个参数：投递标签
   - 第2个参数：`multiple=false` 只处理当前消息
   - 第3个参数：`requeue=true` 重新入队

### Consume 模式的实现

**关键代码**：
```go
for i := 0; i < count; i++ {
    // 获取消息并自动确认（会删除）
    delivery, ok, err := rabbitmqManager.channel.Get(queueName, true)
    if err != nil {
        return messages, fmt.Errorf("获取消息失败: %v", err)
    }
    
    if !ok {
        break
    }
    
    // 处理消息...
    // autoAck=true，消息已被自动确认并删除
}
```

### 并发安全

使用读锁保护：
```go
rabbitmqManager.mu.RLock()
defer rabbitmqManager.mu.RUnlock()

if rabbitmqManager.closed {
    return nil, fmt.Errorf("RabbitMQ连接已关闭")
}

// 执行操作...
```

---

## 📊 API 接口对比

### RabbitMQ vs Kafka 查询功能对比

| 功能 | RabbitMQ | Kafka |
|------|----------|-------|
| 获取队列/主题信息 | ✅ | ✅ |
| 查看消息（不消费） | ✅ Peek | ✅ Fetch |
| 消费消息 | ✅ Consume | ✅ Fetch (自动提交) |
| 分页查询 | ✅ limit 参数 | ✅ offset + limit |
| 消息持久化 | ❌ 默认不支持 | ✅ 持久化日志 |
| 历史消息查询 | ❌ | ✅ |
| 清空队列 | ✅ | ❌ (只能删除) |
| 删除队列/主题 | ✅ | ✅ |

---

## 📋 API 接口列表

### 查询接口

| 序号 | 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|------|
| 1 | 获取队列信息 | GET | `/rabbitmq/queue/info` | 查询队列状态 |
| 2 | 查看消息（不消费） | GET | `/rabbitmq/queue/peek` | Peek 模式 |
| 3 | 消费消息 | GET | `/rabbitmq/queue/consume` | Consume 模式 |
| 4 | 清空队列 | POST | `/rabbitmq/queue/purge` | 删除所有消息 |
| 5 | 删除队列 | POST | `/rabbitmq/queue/delete` | 删除整个队列 |

### 发送接口（已有）

| 序号 | 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|------|
| 1 | 发送简单消息 | POST | `/rabbitmq/send` | 文本消息 |
| 2 | 发送 JSON 消息 | POST | `/rabbitmq/send-json` | JSON 格式 |
| 3 | 发送广播消息 | POST | `/rabbitmq/send-fanout` | Fanout 模式 |
| 4 | 发送直接消息 | POST | `/rabbitmq/send-direct` | Direct 模式 |
| 5 | 发送主题消息 | POST | `/rabbitmq/send-topic` | Topic 模式 |
| 6 | 发送任务消息 | POST | `/rabbitmq/send-task` | Worker 模式 |
| 7 | 批量发送消息 | POST | `/rabbitmq/send-batch` | 批量操作 |
| 8 | 健康检查 | GET | `/rabbitmq/health` | 连接状态 |

---

## 📦 文件清单

### 新增/修改的文件

1. **核心实现**
   - `submodule/support-go.git/bootstrap/rabbitmq.go` - 添加查询方法
   - `app/service/rabbitmqService/rabbitmqService.go` - 添加服务层方法
   - `app/entity/req/rabbitmqReq.go` - 添加请求实体
   - `app/entity/resp/rabbitmqResp.go` - 添加响应实体
   - `app/http/controller/api/rabbitmqController/rabbitmq.go` - 添加控制器方法
   - `routes/api/rabbitmqRouter/rabbitmq.go` - 添加路由

2. **文档和工具**
   - `rabbitmq-api-collection.json` - 更新 Postman 集合
   - `RABBITMQ_QUERY_GUIDE.md` - 查询功能使用指南
   - `RABBITMQ_QUERY_IMPLEMENTATION.md` - 本文档
   - `test-rabbitmq-query.bat` - 快速测试脚本

---

## 🧪 测试方法

### 方法 1: 使用测试脚本

```bash
# 1. 启动服务
http-server.exe

# 2. 运行测试脚本
test-rabbitmq-query.bat
```

### 方法 2: 使用 Postman/Apifox

1. 导入 `rabbitmq-api-collection.json`
2. 依次执行：
   - 11. 获取队列信息
   - 12. 查看队列消息（不消费）
   - 13. 消费队列消息（会删除）
   - 14. 清空队列所有消息
   - 15. 删除队列

### 方法 3: 使用 curl

参考 `RABBITMQ_QUERY_GUIDE.md` 中的完整测试流程。

---

## ⚠️ 注意事项

### 1. Peek vs Consume

- **Peek**: 查看消息但不删除，消息会重新入队
- **Consume**: 查看消息并删除，消息不会再出现

### 2. 限制参数

- `limit` 参数范围：1-100
- 超出范围会返回错误

### 3. 并发安全

- 使用读锁保护并发访问
- 多个请求可以同时查询

### 4. 连接状态

- 查询前会检查连接状态
- 如果连接关闭会返回错误

---

## 🎯 最佳实践

### 开发环境

1. ✅ 使用 Peek 接口调试
2. ✅ 定期清理测试队列
3. ✅ 使用特殊前缀命名测试队列（如 `test-`）

### 生产环境

1. ⚠️ 谨慎使用 Consume 接口
2. ⚠️ 避免使用 Purge 和 Delete 接口
3. ✅ 配置适当的权限控制
4. ✅ 记录所有删除操作的审计日志

---

## 🔄 与 Kafka 功能对比

### 相似之处

- ✅ 都提供队列/主题信息查询
- ✅ 都支持消息内容查看
- ✅ 都支持限制返回数量

### 不同之处

| 特性 | RabbitMQ | Kafka |
|------|----------|-------|
| 消息模型 | 队列（消费后删除） | 日志（持久化存储） |
| Peek 实现 | Nack + Requeue | 不改变 offset |
| 历史查询 | ❌ | ✅ |
| 分区概念 | ❌ | ✅ |
| Offset 管理 | ❌ | ✅ |

---

## 📚 相关文档

- [RABBITMQ_QUERY_GUIDE.md](./RABBITMQ_QUERY_GUIDE.md) - 使用指南
- [RABBITMQ_API_USAGE.md](./RABBITMQ_API_USAGE.md) - API 文档
- [KAFKA_QUERY_FEATURE.md](./KAFKA_QUERY_FEATURE.md) - Kafka 查询功能（对比参考）

---

## 📝 版本历史

- **v1.0.0** (2026-01-20)
  - ✅ 实现获取队列信息
  - ✅ 实现 Peek 模式查看消息
  - ✅ 实现 Consume 模式消费消息
  - ✅ 实现清空队列
  - ✅ 实现删除队列
  - ✅ 完善文档和测试工具

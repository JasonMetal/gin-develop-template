# RabbitMQ API 使用指南

## 📋 概述

本项目提供了完整的 RabbitMQ 消息队列功能，支持：
- Simple（简单队列）
- Work Queue（工作队列）
- Fanout（广播模式）
- Direct（直接路由）
- Topic（主题模式）
- 批量发送
- 健康检查

## 🎯 与 Kafka 的对比

本项目同时集成了 Kafka 和 RabbitMQ，以下是它们的特点对比：

| 特性 | Kafka | RabbitMQ |
|------|-------|----------|
| **定位** | 分布式流平台 | 消息中间件 |
| **吞吐量** | 极高（百万级/秒） | 较高（数万级/秒） |
| **消息持久化** | 默认持久化到磁盘 | 可选持久化 |
| **消息顺序** | 分区内有序 | 队列内有序 |
| **消息路由** | 简单（Topic-Partition） | 灵活（多种交换机类型） |
| **延迟** | 较高（ms级） | 低（μs级） |
| **适用场景** | 日志收集、事件溯源、流处理 | 任务队列、RPC、微服务通信 |

**使用建议**：
- 📊 **日志收集、数据分析** → 使用 Kafka
- 📨 **任务队列、异步处理** → 使用 RabbitMQ
- 🔄 **事件驱动架构** → 两者皆可，根据场景选择

---

## 🚀 快速开始

### 1. 配置 RabbitMQ

编辑配置文件 `config/local/rabbitmq.yml`:

```yaml
host: "localhost"
port: 5672
username: "guest"
password: "guest"
vhost: "/"

pool:
  max_open: 10
  max_idle: 5
  max_lifetime: 3600

producer:
  confirm_mode: true
  mandatory: false
  immediate: false
```

### 2. 启动项目

```bash
go run http-server.go -e local
```

RabbitMQ 会在项目启动时自动初始化。

---

## 📡 API 接口

### 1. 发送简单消息

**接口**: `POST /rabbitmq/send`

**请求参数**:
```json
{
  "queue_name": "test-queue",
  "message": "Hello RabbitMQ!"
}
```

**响应示例**:
```json
{
  "code": 0,
  "data": {
    "queue_name": "test-queue",
    "message": "Hello RabbitMQ!",
    "sent_at": "2024-01-20T10:30:00Z"
  },
  "message": "消息发送成功"
}
```

**cURL 示例**:
```bash
curl -X POST http://localhost:8080/rabbitmq/send \
  -H "Content-Type: application/json" \
  -d '{"queue_name":"test-queue","message":"Hello RabbitMQ"}'
```

---

### 2. 发送 JSON 消息

**接口**: `POST /rabbitmq/send-json`

**请求参数**:
```json
{
  "queue_name": "user-events",
  "data": {
    "user_id": 12345,
    "action": "login",
    "ip": "192.168.1.100",
    "timestamp": 1768874603
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "data": {
    "queue_name": "user-events",
    "message": "map[action:login ip:192.168.1.100 timestamp:1768874603 user_id:12345]",
    "sent_at": "2024-01-20T10:30:00Z"
  },
  "message": "JSON消息发送成功"
}
```

---

### 3. 发送广播消息（Fanout）

**接口**: `POST /rabbitmq/send-fanout`

**请求参数**:
```json
{
  "exchange_name": "logs-fanout",
  "message": "System notification: Server will restart in 5 minutes"
}
```

**特点**：
- 发送到交换机的消息会被广播到所有绑定的队列
- 不需要 routing_key
- 适用于：系统公告、日志广播

---

### 4. 发送直接消息（Direct）

**接口**: `POST /rabbitmq/send-direct`

**请求参数**:
```json
{
  "exchange_name": "logs-direct",
  "routing_key": "error",
  "message": "Database connection failed"
}
```

**特点**：
- 根据 routing_key 精确匹配
- 一对一或一对多路由
- 适用于：日志分级、任务分发

---

### 5. 发送主题消息（Topic）

**接口**: `POST /rabbitmq/send-topic`

**请求参数**:
```json
{
  "exchange_name": "events-topic",
  "routing_key": "user.created.premium",
  "message": "Premium user created: user_12345"
}
```

**Routing Key 模式**：
- `*` - 匹配一个单词
- `#` - 匹配零个或多个单词

**示例**：
```
user.created.*     → 匹配: user.created.free, user.created.premium
user.#             → 匹配: user.created, user.updated.profile
*.created.*        → 匹配: user.created.free, order.created.paid
```

**适用于**：事件系统、微服务通信

---

### 6. 发送任务消息

**接口**: `POST /rabbitmq/send-task`

**请求参数**:
```json
{
  "queue_name": "task-queue",
  "task_name": "send_email",
  "task_data": {
    "to": "user@example.com",
    "subject": "Welcome!",
    "body": "Welcome to our service"
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "data": {
    "queue_name": "task-queue",
    "task_name": "send_email",
    "status": "sent",
    "created_at": "2024-01-20T10:30:00Z"
  },
  "message": "任务发送成功"
}
```

**适用于**：Worker Queue 模式、异步任务处理

---

### 7. 批量发送消息

**接口**: `POST /rabbitmq/send-batch`

**请求参数**:
```json
{
  "queue_name": "batch-queue",
  "messages": [
    "Message 1",
    "Message 2",
    "Message 3",
    "Message 4",
    "Message 5"
  ]
}
```

**响应示例**:
```json
{
  "code": 0,
  "data": {
    "queue_name": "batch-queue",
    "total_count": 5,
    "success_count": 5,
    "failed_count": 0,
    "sent_at": "2024-01-20T10:30:00Z"
  },
  "message": "批量发送成功"
}
```

---

### 8. 健康检查

**接口**: `GET /rabbitmq/health`

**响应示例**:
```json
{
  "code": 0,
  "data": {
    "status": "ok",
    "connected": true,
    "message": "RabbitMQ连接正常"
  }
}
```

---

## 💻 在代码中使用

### 示例 1: 简单发送消息

```go
import "develop-template/app/service/rabbitmqService"

func SendSimpleMessage() {
    service := rabbitmqService.NewRabbitMQService()
    
    err := service.SendMessage("test-queue", "Hello RabbitMQ")
    if err != nil {
        log.Printf("发送失败: %v", err)
    }
}
```

### 示例 2: 发送 JSON 数据

```go
func SendUserEvent() {
    service := rabbitmqService.NewRabbitMQService()
    
    data := map[string]interface{}{
        "user_id": 12345,
        "action": "login",
        "timestamp": time.Now().Unix(),
    }
    
    err := service.SendJSON("user-events", data)
    if err != nil {
        log.Printf("发送失败: %v", err)
    }
}
```

### 示例 3: 发送到交换机（Topic 模式）

```go
func SendTopicMessage() {
    service := rabbitmqService.NewRabbitMQService()
    
    err := service.PublishTopic(
        "events-topic",          // 交换机名称
        "user.created.premium",  // Routing Key
        "New premium user",      // 消息内容
    )
    if err != nil {
        log.Printf("发送失败: %v", err)
    }
}
```

### 示例 4: 批量发送

```go
func SendBatchMessages() {
    service := rabbitmqService.NewRabbitMQService()
    
    messages := []string{
        "Task 1",
        "Task 2",
        "Task 3",
    }
    
    err := service.SendBatch("task-queue", messages)
    if err != nil {
        log.Printf("批量发送失败: %v", err)
    }
}
```

---

## 🎨 使用场景示例

### 场景 1: 异步任务处理

```go
// 发送邮件任务
func SendEmailTask(to, subject, body string) error {
    service := rabbitmqService.NewRabbitMQService()
    
    taskData := map[string]interface{}{
        "to": to,
        "subject": subject,
        "body": body,
    }
    
    return service.SendTask("email-queue", "send_email", taskData)
}
```

### 场景 2: 日志收集（Direct 模式）

```go
// 发送不同级别的日志
func SendLog(level, message string) error {
    service := rabbitmqService.NewRabbitMQService()
    
    // level: "info", "warning", "error"
    return service.PublishDirect("logs-direct", level, message)
}
```

### 场景 3: 事件通知（Topic 模式）

```go
// 发送用户事件
func PublishUserEvent(eventType, userType string, data interface{}) error {
    service := rabbitmqService.NewRabbitMQService()
    
    routingKey := fmt.Sprintf("user.%s.%s", eventType, userType)
    message, _ := json.Marshal(data)
    
    return service.PublishTopic("user-events", routingKey, string(message))
}

// 使用示例
PublishUserEvent("created", "premium", userData)
// Routing Key: user.created.premium
```

### 场景 4: 系统广播（Fanout 模式）

```go
// 发送系统通知给所有在线服务
func BroadcastSystemNotification(message string) error {
    service := rabbitmqService.NewRabbitMQService()
    
    return service.PublishFanout("system-broadcast", message)
}
```

---

## 📊 RabbitMQ 模式详解

### 1. Simple（简单队列）

```
Producer → Queue → Consumer
```

**特点**：
- 一对一通信
- 最简单的模式
- 适用于：点对点消息传递

### 2. Work Queue（工作队列）

```
Producer → Queue → Consumer1
                 → Consumer2
                 → Consumer3
```

**特点**：
- 多个消费者竞争消费
- 轮询分发
- 适用于：任务分发、负载均衡

### 3. Fanout（广播）

```
Producer → Exchange(fanout) → Queue1 → Consumer1
                            → Queue2 → Consumer2
                            → Queue3 → Consumer3
```

**特点**：
- 消息广播到所有队列
- 忽略 routing key
- 适用于：日志广播、实时通知

### 4. Direct（直接路由）

```
Producer → Exchange(direct) → Queue1 (key=error)   → Consumer1
           ↓
           routing_key=error → Queue2 (key=info)    → Consumer2
                            → Queue3 (key=warning) → Consumer3
```

**特点**：
- 精确匹配 routing key
- 一对一或一对多
- 适用于：日志分级、任务路由

### 5. Topic（主题）

```
Producer → Exchange(topic) → Queue1 (key=user.*.premium) → Consumer1
           ↓
           routing_key=user.created.premium
                          → Queue2 (key=user.created.*) → Consumer2
                          → Queue3 (key=#) → Consumer3
```

**特点**：
- 模式匹配（`*` 和 `#`）
- 灵活的路由规则
- 适用于：事件系统、复杂路由

---

## 🔧 配置说明

### 连接池配置

```yaml
pool:
  max_open: 10        # 最大连接数
  max_idle: 5         # 最大空闲连接数
  max_lifetime: 3600  # 连接最大生命周期（秒）
```

### 生产者配置

```yaml
producer:
  confirm_mode: true   # 消息确认模式
  mandatory: false     # 消息无法路由时是否返回
  immediate: false     # 消息无法立即消费时是否返回
```

### 消费者配置

```yaml
consumer:
  auto_ack: false      # 是否自动确认
  prefetch_count: 10   # 预取消息数量
  prefetch_size: 0     # 预取消息大小
```

### 重连配置

```yaml
reconnect:
  max_retries: 5       # 最大重试次数
  interval: 5          # 重试间隔（秒）
```

---

## 🆚 Kafka vs RabbitMQ 选择指南

### 使用 Kafka 的场景

✅ **日志聚合** - 收集大量应用日志
✅ **事件溯源** - 需要消息回溯的场景
✅ **流处理** - 实时数据流处理
✅ **高吞吐量** - 需要处理百万级消息/秒
✅ **数据管道** - 大数据系统间的数据传输

**项目中的 Kafka 接口**：
```bash
POST /kafka/send         # 发送Kafka消息
POST /kafka/send-json    # 发送JSON到Kafka
GET  /kafka/messages     # 查询Kafka消息（支持分页）
GET  /kafka/topic-info   # 获取主题信息
```

### 使用 RabbitMQ 的场景

✅ **任务队列** - 异步任务处理
✅ **RPC** - 远程过程调用
✅ **低延迟** - 需要微秒级延迟
✅ **复杂路由** - 需要灵活的消息路由
✅ **消息优先级** - 需要消息优先级队列

**项目中的 RabbitMQ 接口**：
```bash
POST /rabbitmq/send         # 发送简单消息
POST /rabbitmq/send-json    # 发送JSON消息
POST /rabbitmq/send-fanout  # 广播消息
POST /rabbitmq/send-direct  # 直接路由
POST /rabbitmq/send-topic   # 主题路由
POST /rabbitmq/send-task    # 任务消息
POST /rabbitmq/send-batch   # 批量发送
GET  /rabbitmq/health       # 健康检查
```

---

## 🐛 常见问题

### Q1: Kafka 和 RabbitMQ 可以同时使用吗？

**A**: 当然可以！本项目已经同时集成了两者。根据不同的业务场景选择：
- 日志收集、数据分析 → Kafka
- 任务队列、即时通信 → RabbitMQ

### Q2: 消息发送失败怎么办？

**A**: 项目已实现自动重连机制：
- RabbitMQ 连接断开会自动重连（最多5次）
- Kafka 有内置的重试机制
- 建议在业务代码中也添加重试逻辑

### Q3: 如何监控消息发送情况？

**A**: 
- 使用健康检查接口：`GET /rabbitmq/health`
- 查看服务日志
- 使用 RabbitMQ Management UI (http://localhost:15672)
- 使用 Kafka Manager 或 Kafdrop

### Q4: 如何保证消息不丢失？

**A**: 
- 开启持久化：配置中 `durable: true`
- 使用消息确认：`confirm_mode: true`
- 手动确认消费：`auto_ack: false`

---

## 📚 更多资源

- [RabbitMQ 官方文档](https://www.rabbitmq.com/documentation.html)
- [Kafka API 使用指南](./KAFKA_API_USAGE.md)
- [Kafka 性能优化](./KAFKA_PERFORMANCE_OPTIMIZATION.md)

---

更新时间: 2024-01-20

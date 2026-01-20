# RabbitMQ 查询功能使用指南

## 📋 目录

1. [RabbitMQ UI 说明](#rabbitmq-ui-说明)
2. [查询 API 接口](#查询-api-接口)
3. [完整测试流程](#完整测试流程)
4. [常见问题](#常见问题)

---

## 🖥️ RabbitMQ UI 说明

### 访问管理界面

**URL**: http://localhost:15672  
**默认账号**: admin / admin123

### 如何在 UI 中查看消息

#### ⚠️ 重要说明：RabbitMQ vs Kafka 的区别

| 特性 | Kafka | RabbitMQ |
|------|-------|----------|
| 消息模型 | **消息日志** | **消息队列** |
| 消息持久化 | 持久化存储，可随时回溯 | 消费后删除（默认） |
| 查询历史 | ✅ 支持 | ❌ 不支持 |
| UI 功能 | 完整的消息浏览 | 只能查看未消费的消息 |

#### 在 RabbitMQ UI 中查看消息步骤：

1. 登录管理界面：http://localhost:15672
2. 点击顶部导航栏的 **Queues** 标签
3. 点击您想查看的 **队列名称**
4. 向下滚动到 **"Get messages"** 区域
5. 设置参数：
   - **Messages**: 要获取的消息数量（如 10）
   - **Ack Mode**: 选择 `Nack message requeue true`（查看后重新入队）或 `Automatic ack`（查看后删除）
6. 点击 **"Get Message(s)"** 按钮

#### 💡 提示

- RabbitMQ UI 只能查看**未被消费**的消息
- 如果队列是空的，说明消息已经被消费者处理了
- 如果需要查看所有消息（包括历史），建议使用本项目提供的 API 接口（需要配置消息持久化）

---

## 🔍 查询 API 接口

本项目提供了完整的 RabbitMQ 查询 API，使用更加灵活：

### 1. 获取队列信息

**接口**: `GET /rabbitmq/queue/info`

**参数**:
```
queue: 队列名称（必填）
```

**示例**:
```bash
curl "http://localhost:8080/rabbitmq/queue/info?queue=test-queue"
```

**响应**:
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "name": "test-queue",
    "messages": 15,          // 队列中的消息数
    "consumers": 0,          // 消费者数量
    "durable": true,         // 是否持久化
    "auto_delete": false,    // 是否自动删除
    "exclusive": false       // 是否排他
  }
}
```

---

### 2. 查看队列消息（不消费）🔍

**接口**: `GET /rabbitmq/queue/peek`

**参数**:
```
queue: 队列名称（必填）
limit: 查看数量（可选，默认10，最多100）
```

**示例**:
```bash
curl "http://localhost:8080/rabbitmq/queue/peek?queue=test-queue&limit=10"
```

**特点**:
- ✅ 查看消息内容
- ✅ 消息会重新入队（不删除）
- ✅ 适合调试和监控

**响应**:
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "queue": "test-queue",
    "total": 3,
    "messages": [
      {
        "body": "Hello RabbitMQ! 这是一条测试消息",
        "content_type": "text/plain",
        "delivery_mode": 2,
        "priority": 0,
        "correlation_id": "",
        "reply_to": "",
        "expiration": "",
        "message_id": "",
        "timestamp": "2026-01-20T10:30:00Z",
        "type": "",
        "user_id": "",
        "app_id": "",
        "headers": {},
        "delivery_tag": 1,
        "redelivered": false,
        "exchange": "",
        "routing_key": "test-queue"
      }
    ]
  }
}
```

---

### 3. 消费队列消息（会删除）⚠️

**接口**: `GET /rabbitmq/queue/consume`

**参数**:
```
queue: 队列名称（必填）
limit: 消费数量（可选，默认10，最多100）
```

**示例**:
```bash
curl "http://localhost:8080/rabbitmq/queue/consume?queue=test-queue&limit=5"
```

**特点**:
- ⚠️ 消息会被删除
- ✅ 适合手动处理积压消息
- ✅ 返回格式与 peek 相同

---

### 4. 清空队列 🗑️

**接口**: `POST /rabbitmq/queue/purge`

**请求体**:
```json
{
  "queue": "test-queue"
}
```

**示例**:
```bash
curl -X POST http://localhost:8080/rabbitmq/queue/purge \
  -H "Content-Type: application/json" \
  -d '{"queue":"test-queue"}'
```

**响应**:
```json
{
  "code": 200,
  "msg": "成功清空队列，删除了 15 条消息",
  "data": {
    "queue": "test-queue",
    "deleted_count": 15,
    "success": true
  }
}
```

---

### 5. 删除队列 ❌

**接口**: `POST /rabbitmq/queue/delete`

**请求体**:
```json
{
  "queue": "test-queue-to-delete",
  "if_unused": false,  // 仅在没有消费者时删除
  "if_empty": false    // 仅在队列为空时删除
}
```

**示例**:
```bash
curl -X POST http://localhost:8080/rabbitmq/queue/delete \
  -H "Content-Type: application/json" \
  -d '{
    "queue": "test-queue-to-delete",
    "if_unused": false,
    "if_empty": false
  }'
```

**响应**:
```json
{
  "code": 200,
  "msg": "成功删除队列，删除了 0 条消息",
  "data": {
    "queue": "test-queue-to-delete",
    "deleted_count": 0,
    "success": true
  }
}
```

---

## 🧪 完整测试流程

### 步骤 1: 发送测试消息

```bash
curl -X POST http://localhost:8080/rabbitmq/send \
  -H "Content-Type: application/json" \
  -d '{
    "queue_name": "my-test-queue",
    "message": "测试消息 1"
  }'
```

### 步骤 2: 查看队列信息

```bash
curl "http://localhost:8080/rabbitmq/queue/info?queue=my-test-queue"
```

**预期结果**: `messages` 字段应该 > 0

### 步骤 3: 查看消息（不消费）

```bash
curl "http://localhost:8080/rabbitmq/queue/peek?queue=my-test-queue&limit=10"
```

**预期结果**: 能看到消息内容

### 步骤 4: 再次查看队列信息

```bash
curl "http://localhost:8080/rabbitmq/queue/info?queue=my-test-queue"
```

**预期结果**: `messages` 数量不变（因为 peek 不消费）

### 步骤 5: 发送更多测试消息

```bash
# 发送 JSON 消息
curl -X POST http://localhost:8080/rabbitmq/send-json \
  -H "Content-Type: application/json" \
  -d '{
    "queue_name": "my-test-queue",
    "data": {
      "user_id": 12345,
      "action": "test",
      "timestamp": 1737363600
    }
  }'

# 批量发送
curl -X POST http://localhost:8080/rabbitmq/send-batch \
  -H "Content-Type: application/json" \
  -d '{
    "queue_name": "my-test-queue",
    "messages": ["消息1", "消息2", "消息3", "消息4", "消息5"]
  }'
```

### 步骤 6: 再次查看消息

```bash
curl "http://localhost:8080/rabbitmq/queue/peek?queue=my-test-queue&limit=20"
```

### 步骤 7: 消费部分消息

```bash
curl "http://localhost:8080/rabbitmq/queue/consume?queue=my-test-queue&limit=3"
```

**预期结果**: 返回 3 条消息，队列中的消息数减少 3

### 步骤 8: 验证消息被消费

```bash
curl "http://localhost:8080/rabbitmq/queue/info?queue=my-test-queue"
```

**预期结果**: `messages` 数量减少了 3

### 步骤 9: 清空队列

```bash
curl -X POST http://localhost:8080/rabbitmq/queue/purge \
  -H "Content-Type: application/json" \
  -d '{"queue":"my-test-queue"}'
```

### 步骤 10: 删除队列（可选）

```bash
curl -X POST http://localhost:8080/rabbitmq/queue/delete \
  -H "Content-Type: application/json" \
  -d '{"queue":"my-test-queue"}'
```

---

## ❓ 常见问题

### Q1: 为什么 RabbitMQ UI 看不到消息内容？

**A**: 可能的原因：
1. **消息已被消费**：RabbitMQ 是消息队列，消费后消息就删除了
2. **队列是空的**：检查队列信息中的 `messages` 字段
3. **选错了队列**：确认队列名称是否正确

**解决方法**:
- 使用本项目的 API 接口查询
- 在发送消息前先查看，不要让消费者自动消费

### Q2: Peek 和 Consume 有什么区别？

**A**: 

| 操作 | Peek | Consume |
|------|------|---------|
| 查看消息 | ✅ | ✅ |
| 删除消息 | ❌ | ✅ |
| 消息重新入队 | ✅ | ❌ |
| 适用场景 | 调试、监控 | 真正消费 |

### Q3: 如何保留消息历史用于查询？

**A**: RabbitMQ 本身不支持消息历史查询，建议：
1. **使用 Kafka**：如果需要消息历史查询，改用 Kafka
2. **插件方案**：使用 RabbitMQ 的 Event Exchange 插件将消息转发到其他存储
3. **应用层记录**：在应用层将消息存储到数据库或 Elasticsearch

### Q4: 队列信息显示 messages=0，但我刚发了消息？

**A**: 可能的原因：
1. **消费者太快**：有消费者正在监听队列，消息立即被消费
2. **队列不存在**：发送到了不同的队列
3. **TTL 过期**：消息设置了过期时间并已过期

**解决方法**:
```bash
# 1. 先停止消费者
# 2. 发送消息
curl -X POST http://localhost:8080/rabbitmq/send \
  -H "Content-Type: application/json" \
  -d '{"queue_name":"test-queue","message":"测试"}'

# 3. 立即查看
curl "http://localhost:8080/rabbitmq/queue/info?queue=test-queue"
```

### Q5: 为什么 Peek 后消息还会被重复消费？

**A**: 这是正常的！Peek 操作会：
1. 从队列中取出消息
2. 读取消息内容
3. 使用 `Nack(requeue=true)` 将消息重新放回队列

这样消费者仍然可以正常消费这些消息。

### Q6: 如何在 Apifox/Postman 中使用？

**A**: 
1. 导入 `rabbitmq-api-collection.json` 文件
2. 所有接口已预配置好示例数据
3. 直接点击 "Send" 即可测试

---

## 📚 相关文档

- [RabbitMQ API 使用文档](./RABBITMQ_API_USAGE.md)
- [RabbitMQ 集成总结](./RABBITMQ_INTEGRATION_SUMMARY.md)
- [RabbitMQ 快速开始](./RABBITMQ_QUICKSTART.md)
- [Kafka 查询功能](./KAFKA_QUERY_FEATURE.md) - 对比参考

---

## 🎯 最佳实践

### 开发环境

1. ✅ 使用 **Peek** 接口调试
2. ✅ 定期清理测试队列
3. ✅ 为测试队列使用特殊前缀（如 `test-`）

### 生产环境

1. ⚠️ 谨慎使用 **Consume** 接口
2. ⚠️ 避免使用 **Purge** 和 **Delete** 接口
3. ✅ 使用监控工具而不是频繁查询
4. ✅ 配置适当的消息持久化策略

---

## 🔧 技术实现

### Peek 模式实现原理

```go
// 1. 从队列获取消息（不自动确认）
delivery, ok, err := channel.Get(queueName, false)

// 2. 读取消息内容
message := QueueMessage{
    Body: string(delivery.Body),
    // ... 其他字段
}

// 3. Nack 并重新入队（requeue=true）
err = channel.Nack(delivery.DeliveryTag, false, true)
```

### Consume 模式实现原理

```go
// 从队列获取消息（自动确认）
delivery, ok, err := channel.Get(queueName, true)
// autoAck=true，消息会被自动确认并删除
```

---

## 📞 支持

如有问题，请查看：
1. RabbitMQ 管理界面: http://localhost:15672
2. 项目 README
3. RabbitMQ 官方文档: https://www.rabbitmq.com/documentation.html

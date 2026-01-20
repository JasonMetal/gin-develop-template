# RabbitMQ 集成说明

本项目已成功集成 RabbitMQ 消息队列功能，支持多种消息模式和完整的查询功能。

## 快速开始

### 1. 启动 RabbitMQ

使用 Docker Compose 快速启动：

```bash
# 启动 RabbitMQ (包含管理界面)
docker-compose -f docker-compose.yml up -d rabbitmq

# 或使用项目提供的配置
docker run -d --name rabbitmq \
  -p 5672:5672 \
  -p 15672:15672 \
  -e RABBITMQ_DEFAULT_USER=admin \
  -e RABBITMQ_DEFAULT_PASS=admin123 \
  rabbitmq:3.12-management
```

### 2. 配置 RabbitMQ

编辑 `config/{env}/rabbitmq.yml` 文件：

```yaml
# 连接配置
host: "localhost"
port: 5672
username: "admin"
password: "admin123"
vhost: "/"

# 连接池配置
pool:
  max_open: 10
  max_idle: 5
  max_lifetime: 3600
```

### 3. 使用示例

```go
import "develop-template/app/service/rabbitmqService"

service := rabbitmqService.NewRabbitMQService()

// 1. 发送简单消息
service.SendMessage("test-queue", "Hello RabbitMQ!")

// 2. 发送JSON数据
service.SendJSON("user-events", map[string]interface{}{
    "user_id": 123,
    "action": "login",
    "timestamp": time.Now().Unix(),
})

// 3. 发送到交换机（Direct模式）
service.SendDirect("logs-direct", "error", "数据库连接失败")

// 4. 发送到交换机（Topic模式）
service.SendTopic("events-topic", "user.created.premium", "新用户注册")

// 5. 查看队列消息（不消费）
messages, _ := service.PeekMessages("test-queue", 10)

// 6. 获取队列信息
info, _ := service.GetQueueInfo("test-queue")
fmt.Printf("队列消息数: %d\n", info.Messages)
```

## 功能特性

### 📤 消息发送

- ✅ **简单队列** - 点对点消息传递
- ✅ **Worker 队列** - 任务分发和负载均衡
- ✅ **发布/订阅** - Fanout 广播模式
- ✅ **路由模式** - Direct 精确路由
- ✅ **主题模式** - Topic 通配符路由
- ✅ **批量发送** - 批量消息处理
- ✅ **JSON 序列化** - 自动 JSON 编解码
- ✅ **自定义交换机** - 灵活的消息路由

### 🔍 查询功能

- ✅ **队列信息查询** - 消息数、消费者数等
- ✅ **Peek 模式** - 查看消息但不删除
- ✅ **Consume 模式** - 消费并删除消息
- ✅ **队列管理** - 清空、删除队列
- ✅ **分页支持** - limit 参数控制

### 🔧 高级特性

- ✅ **连接池** - 高效的连接管理
- ✅ **自动重连** - 断线自动恢复
- ✅ **消息持久化** - 可配置持久化策略
- ✅ **健康检查** - 连接状态监控
- ✅ **优雅关闭** - 安全的资源释放
- ✅ **并发安全** - 读写锁保护

## API 参考

### 消息发送接口

#### 基础队列

```go
// 发送简单消息
SendMessage(queueName, message string) error

// 发送 JSON 消息
SendJSON(queueName string, data interface{}) error

// 发送任务消息
SendTask(queueName, taskName string, taskData interface{}) error

// 批量发送消息
SendBatch(queueName string, messages []string) error
```

#### 交换机模式

```go
// Fanout 广播模式
SendFanout(exchangeName, message string) error

// Direct 直接路由
SendDirect(exchangeName, routingKey, message string) error

// Topic 主题模式
SendTopic(exchangeName, routingKey, message string) error

// 自定义交换机
SendExchange(exchangeName, exchangeType, routingKey, message, contentType string, headers map[string]interface{}) error
```

### 查询接口

```go
// 获取队列信息
GetQueueInfo(queueName string) (*QueueInfo, error)

// 查看消息（不消费）
PeekMessages(queueName string, count int) ([]QueueMessage, error)

// 消费消息（会删除）
ConsumeMessages(queueName string, count int) ([]QueueMessage, error)

// 清空队列
PurgeQueue(queueName string) (int, error)

// 删除队列
DeleteQueue(queueName string, ifUnused, ifEmpty bool) (int, error)

// 健康检查
HealthCheck() (bool, error)
```

## HTTP API 接口

### 消息发送

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 发送简单消息 | POST | `/rabbitmq/send` | 文本消息 |
| 发送 JSON | POST | `/rabbitmq/send-json` | JSON 格式 |
| 发送广播 | POST | `/rabbitmq/send-fanout` | Fanout 模式 |
| 发送直接消息 | POST | `/rabbitmq/send-direct` | Direct 模式 |
| 发送主题消息 | POST | `/rabbitmq/send-topic` | Topic 模式 |
| 发送任务 | POST | `/rabbitmq/send-task` | Worker 模式 |
| 批量发送 | POST | `/rabbitmq/send-batch` | 批量操作 |
| 健康检查 | GET | `/rabbitmq/health` | 连接状态 |

### 查询接口

| 接口 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 获取队列信息 | GET | `/rabbitmq/queue/info` | 队列状态 |
| 查看消息 | GET | `/rabbitmq/queue/peek` | Peek 模式 |
| 消费消息 | GET | `/rabbitmq/queue/consume` | Consume 模式 |
| 清空队列 | POST | `/rabbitmq/queue/purge` | 删除所有消息 |
| 删除队列 | POST | `/rabbitmq/queue/delete` | 删除队列 |

### HTTP 使用示例

```bash
# 1. 发送消息
curl -X POST http://localhost:8080/rabbitmq/send \
  -H "Content-Type: application/json" \
  -d '{"queue_name":"test-queue","message":"Hello RabbitMQ"}'

# 2. 查看队列信息
curl "http://localhost:8080/rabbitmq/queue/info?queue=test-queue"

# 3. 查看消息（不消费）
curl "http://localhost:8080/rabbitmq/queue/peek?queue=test-queue&limit=10"

# 4. 消费消息
curl "http://localhost:8080/rabbitmq/queue/consume?queue=test-queue&limit=5"

# 5. 清空队列
curl -X POST http://localhost:8080/rabbitmq/queue/purge \
  -H "Content-Type: application/json" \
  -d '{"queue":"test-queue"}'
```

## RabbitMQ 管理界面

### 访问地址

- **URL**: http://localhost:15672
- **用户名**: admin
- **密码**: admin123

### 管理界面功能

- 📊 **Overview** - 系统概览和统计
- 📬 **Queues** - 队列管理和监控
- 🔄 **Exchanges** - 交换机配置
- 🔗 **Connections** - 连接管理
- 👥 **Users** - 用户权限管理
- ⚙️ **Admin** - 系统配置

### 在 UI 中查看消息

1. 登录管理界面
2. 点击 **Queues** 标签
3. 点击队列名称
4. 向下滚动到 **"Get messages"** 区域
5. 设置参数并点击 **"Get Message(s)"**

## 配置示例

### 基础配置 (config/local/rabbitmq.yml)

```yaml
# 连接配置
host: "localhost"
port: 5672
username: "admin"
password: "admin123"
vhost: "/"

# 连接池配置
pool:
  max_open: 10        # 最大连接数
  max_idle: 5         # 最大空闲连接
  max_lifetime: 3600  # 连接最大生命周期(秒)

# 生产者配置
producer:
  confirm_mode: false # 消息确认模式
  mandatory: false    # 强制路由
  immediate: false    # 立即投递

# 消费者配置
consumer:
  auto_ack: false       # 自动确认
  prefetch_count: 10    # 预取数量
  prefetch_size: 0      # 预取大小

# 重连配置
reconnect:
  max_retries: 5   # 最大重试次数
  interval: 3      # 重试间隔(秒)

# 队列配置
queue:
  durable: true       # 持久化
  auto_delete: false  # 自动删除
  exclusive: false    # 排他
  no_wait: false      # 不等待

# 交换机配置
exchange:
  durable: true       # 持久化
  auto_delete: false  # 自动删除
  internal: false     # 内部使用
  no_wait: false      # 不等待
```

### 生产环境配置 (config/prod/rabbitmq.yml)

```yaml
host: "rabbitmq.example.com"
port: 5672
username: "prod_user"
password: "strong_password_here"
vhost: "/production"

pool:
  max_open: 50
  max_idle: 20
  max_lifetime: 7200

producer:
  confirm_mode: true  # 生产环境建议启用
  mandatory: true
  immediate: false

reconnect:
  max_retries: 10
  interval: 5

queue:
  durable: true       # 生产环境必须持久化
  auto_delete: false
  exclusive: false
  no_wait: false
```

## 消息模式详解

### 1. 简单队列 (Simple Queue)

最基本的点对点模式。

```go
service.SendMessage("task-queue", "处理订单#12345")
```

**适用场景**：
- 简单的任务队列
- 点对点消息传递

### 2. Worker 队列 (Work Queue)

多个消费者竞争消费，实现负载均衡。

```go
service.SendTask("email-queue", "send_email", map[string]interface{}{
    "to": "user@example.com",
    "subject": "欢迎",
})
```

**适用场景**：
- 耗时任务处理
- 负载均衡
- 并发处理

### 3. 发布/订阅 (Fanout)

消息广播到所有订阅者。

```go
service.SendFanout("logs-fanout", "系统维护通知")
```

**适用场景**：
- 系统通知
- 日志收集
- 广播消息

### 4. 路由模式 (Direct)

根据 routing key 精确匹配。

```go
service.SendDirect("logs-direct", "error", "数据库错误")
service.SendDirect("logs-direct", "info", "用户登录")
```

**适用场景**：
- 日志分级
- 消息分类
- 精确路由

### 5. 主题模式 (Topic)

支持通配符的路由匹配。

```go
// user.* 可以匹配 user.created, user.updated
service.SendTopic("events-topic", "user.created", "新用户")
service.SendTopic("events-topic", "user.updated.premium", "用户升级")
service.SendTopic("events-topic", "order.created", "新订单")
```

**通配符规则**：
- `*` - 匹配一个单词
- `#` - 匹配零个或多个单词

**适用场景**：
- 复杂的消息路由
- 事件系统
- 微服务通信

## Peek vs Consume 对比

| 特性 | Peek（查看） | Consume（消费） |
|------|-------------|----------------|
| 查看消息 | ✅ | ✅ |
| 删除消息 | ❌ | ✅ |
| 消息重新入队 | ✅ | ❌ |
| 适用场景 | 调试、监控 | 真正消费 |
| HTTP 接口 | `/queue/peek` | `/queue/consume` |

## 项目结构

```
gin-develop-template/
├── app/
│   ├── service/rabbitmqService/      # RabbitMQ服务层
│   │   └── rabbitmqService.go
│   ├── http/controller/api/
│   │   └── rabbitmqController/       # HTTP控制器
│   │       └── rabbitmq.go
│   └── entity/
│       ├── req/rabbitmqReq.go        # 请求实体
│       └── resp/rabbitmqResp.go      # 响应实体
├── config/
│   ├── local/rabbitmq.yml            # 本地配置
│   ├── bvt/rabbitmq.yml              # 测试配置
│   ├── test/rabbitmq.yml             # 测试配置
│   └── prod/rabbitmq.yml             # 生产配置
├── routes/api/
│   └── rabbitmqRouter/               # 路由注册
│       └── rabbitmq.go
├── submodule/support-go.git/
│   └── bootstrap/rabbitmq.go         # RabbitMQ核心实现
└── README_RABBITMQ.md                # 本文件
```

## 依赖

```go
github.com/rabbitmq/amqp091-go v1.10.0
```

## 常见问题

### Q1: RabbitMQ 连接失败？

**A**: 检查以下几点：
1. RabbitMQ 是否已启动：`docker ps | grep rabbitmq`
2. 端口是否正确：默认 5672（AMQP），15672（管理界面）
3. 用户名密码是否正确
4. 防火墙是否开放端口

### Q2: 为什么 UI 看不到消息？

**A**: RabbitMQ 是消息队列，不是消息日志：
- 消息被消费后就删除了
- 只能看到**未消费**的消息
- 如需查看历史，使用本项目的 Peek API

### Q3: guest 用户无法远程连接？

**A**: RabbitMQ 的 `guest` 用户只能本地连接，解决方法：
1. 创建新用户：
```bash
docker exec rabbitmq rabbitmqctl add_user admin admin123
docker exec rabbitmq rabbitmqctl set_user_tags admin administrator
docker exec rabbitmq rabbitmqctl set_permissions -p / admin ".*" ".*" ".*"
```
2. 或允许 guest 远程连接（不推荐）

### Q4: Peek 和 Consume 有什么区别？

**A**: 
- **Peek**: 查看消息但不删除，适合调试
- **Consume**: 查看并删除消息，适合真正消费

### Q5: 如何选择消息模式？

**A**: 
- **简单队列**: 点对点，一对一
- **Worker**: 负载均衡，多消费者竞争
- **Fanout**: 广播，所有订阅者都收到
- **Direct**: 精确路由，按 routing key 匹配
- **Topic**: 模糊路由，支持通配符

### Q6: 消息持久化如何配置？

**A**: 在配置文件中设置：
```yaml
queue:
  durable: true  # 队列持久化

producer:
  # 发送时使用 DeliveryMode: 2 (持久化)
```

### Q7: 如何处理消息堆积？

**A**: 
1. 查看队列信息：`GET /rabbitmq/queue/info`
2. 增加消费者数量
3. 使用 Worker 模式并发处理
4. 必要时使用 Consume API 手动清理

### Q8: 连接池配置建议？

**A**: 
- **开发环境**: max_open=10, max_idle=5
- **生产环境**: max_open=50, max_idle=20
- 根据实际负载调整

## 性能优化建议

### 1. 连接池优化

```yaml
pool:
  max_open: 50      # 根据并发量调整
  max_idle: 20      # 保持适量空闲连接
  max_lifetime: 7200 # 定期重建连接
```

### 2. 批量发送

```go
// 使用批量发送减少网络开销
service.SendBatch("task-queue", []string{
    "任务1", "任务2", "任务3", // ...
})
```

### 3. 消息持久化

```yaml
# 仅重要消息启用持久化
queue:
  durable: true  # 队列持久化

# 发送时设置 DeliveryMode: 2
```

### 4. 预取优化

```yaml
consumer:
  prefetch_count: 10  # 一次预取多条消息
```

## 最佳实践

### 开发环境

1. ✅ 使用 Peek 接口调试
2. ✅ 为测试队列使用特殊前缀（如 `test-`）
3. ✅ 定期清理测试数据
4. ✅ 使用管理界面监控

### 生产环境

1. ✅ 启用消息持久化
2. ✅ 配置消息确认模式
3. ✅ 设置合理的连接池大小
4. ✅ 配置自动重连
5. ✅ 使用强密码
6. ⚠️ 避免在生产环境使用 Purge/Delete API
7. ✅ 监控队列堆积情况
8. ✅ 配置告警策略

## 监控指标

### 关键指标

1. **队列深度** - 未消费消息数量
2. **消费速率** - 每秒处理消息数
3. **连接数** - 当前活跃连接
4. **消息堆积** - 积压消息数量
5. **错误率** - 发送/接收失败率

### 查询方式

```bash
# 查看队列信息
curl "http://localhost:8080/rabbitmq/queue/info?queue=my-queue"

# 查看健康状态
curl "http://localhost:8080/rabbitmq/health"

# 访问管理界面
open http://localhost:15672
```

## 测试

### 快速测试

```bash
# 1. 发送测试消息
curl -X POST http://localhost:8080/rabbitmq/send \
  -H "Content-Type: application/json" \
  -d '{"queue_name":"test-queue","message":"测试消息"}'

# 2. 查看队列信息
curl "http://localhost:8080/rabbitmq/queue/info?queue=test-queue"

# 3. 查看消息
curl "http://localhost:8080/rabbitmq/queue/peek?queue=test-queue&limit=10"
```

### 完整测试流程

参考项目中的测试脚本和 Postman Collection。

## 故障排查

### 连接问题

```bash
# 1. 检查 RabbitMQ 状态
docker ps | grep rabbitmq

# 2. 查看日志
docker logs rabbitmq

# 3. 测试连接
telnet localhost 5672

# 4. 检查服务健康
curl http://localhost:8080/rabbitmq/health
```

### 消息丢失

1. 检查队列是否持久化
2. 检查消息是否设置持久化
3. 检查消费者是否正确 ACK
4. 查看管理界面的消息统计

### 性能问题

1. 查看队列堆积情况
2. 增加消费者数量
3. 优化消息处理逻辑
4. 调整预取数量
5. 考虑使用批量操作

## RabbitMQ vs Kafka

| 特性 | RabbitMQ | Kafka |
|------|----------|-------|
| 消息模型 | 消息队列 | 消息日志 |
| 消息持久化 | 可选（消费后删除） | 持久化存储 |
| 消息顺序 | 队列级别保证 | 分区级别保证 |
| 消息路由 | ✅ 丰富（5种模式） | ❌ 简单 |
| 历史查询 | ❌ | ✅ |
| 吞吐量 | 中等 | 极高 |
| 延迟 | 低 | 中等 |
| 使用场景 | 任务队列、RPC | 日志、流处理 |

## 版本历史

- **v1.0.0** (2026-01-20): 初始版本
  - ✅ 完成 RabbitMQ 核心功能
  - ✅ 实现 5 种消息模式
  - ✅ 添加查询功能（Peek/Consume）
  - ✅ 实现连接池和自动重连
  - ✅ 提供完整的 HTTP API
  - ✅ 集成管理界面
  - ✅ 完善配置和文档

## 相关资源

- 📖 [RabbitMQ 官方文档](https://www.rabbitmq.com/documentation.html)
- 🎓 [RabbitMQ 教程](https://www.rabbitmq.com/tutorials)
- 🔧 [amqp091-go 文档](https://pkg.go.dev/github.com/rabbitmq/amqp091-go)
- 🖥️ [管理界面](http://localhost:15672)

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可

本项目采用 MIT 许可证。

---

**快速链接**:
- 🚀 [快速开始](#快速开始)
- 📚 [API 参考](#api-参考)
- 🔍 [查询功能](#查询接口)
- ❓ [常见问题](#常见问题)
- 🎯 [最佳实践](#最佳实践)

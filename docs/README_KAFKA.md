# Kafka 集成说明

本项目已成功集成 Apache Kafka 消息队列功能。

## 快速开始

### 1. 配置 Kafka

编辑 `config/{env}/kafka.yml` 文件：

```yaml
brokers:
  - host: "localhost"
    port: 9092
```

### 2. 使用示例

```go
import "develop-template/app/service/kafkaService"

service := kafkaService.NewKafkaService()

// 发送消息
service.SendMessage("test-topic", "Hello Kafka!")

// 发送JSON数据
service.SendJSON("user-events", map[string]interface{}{
    "user_id": 123,
    "action": "login",
})

// 发送日志
service.SendLog("app-logs", "ERROR", "错误消息", map[string]interface{}{
    "request_id": "req-123",
})
```

## 文档

- 📖 [完整使用指南](docs/KAFKA_INTEGRATION_GUIDE.md)
- 📊 [测试报告](tests/kafka/KAFKA_TEST_REPORT.md)
- 🧪 [运行测试](tests/kafka/run_tests.bat)

## 功能特性

- ✅ 同步/异步消息发送
- ✅ 批量消息发送
- ✅ JSON数据自动序列化
- ✅ 结构化日志发送
- ✅ 事件消息发送
- ✅ 指标数据发送
- ✅ 上下文支持 (超时控制)
- ✅ 错误处理和重试
- ✅ 优雅关闭
- ✅ 完整的单元测试 (覆盖率 84.2%)

## 测试

运行测试脚本：

```bash
# Windows
.\tests\kafka\run_tests.bat

# Linux/Mac
bash tests/kafka/run_tests.sh
```

或手动运行：

```bash
# Bootstrap测试
cd submodule/support-go.git
go test -v ./bootstrap/... -run TestKafka

# Service测试
cd ../../
go test -v develop-template/app/service/kafkaService
```

## API 参考

### 基础方法

- `SendMessage(topic, message)` - 同步发送
- `SendMessageAsync(topic, message)` - 异步发送
- `SendMessageWithContext(ctx, topic, message)` - 带上下文发送
- `SendBatch(topic, messages)` - 批量发送

### 高级方法

- `SendJSON(topic, data)` - 发送JSON数据
- `SendLog(topic, level, message, extra)` - 发送日志
- `SendEvent(topic, eventType, eventData)` - 发送事件
- `SendMetric(topic, metricName, value, tags)` - 发送指标

## 配置示例

```yaml
brokers:
  - host: "kafka1.example.com"
    port: 9092
  - host: "kafka2.example.com"
    port: 9092

ssl:
  enable: true  # 生产环境建议启用

producer:
  required_acks: -1      # 等待所有副本确认
  max_retries: 5         # 最大重试5次
  return_successes: true
  return_errors: true

consumer:
  group_id: "gin-develop-template"
  auto_commit: true

version: "3.7.2"
```

## 测试结果

### Bootstrap Kafka 模块
- ✅ 13/14 测试通过 (1个跳过)
- 📊 覆盖率: 13.4% (核心Kafka功能)

### KafkaService 服务层
- ✅ 11/11 测试通过
- 📊 覆盖率: 84.2%

### 性能测试
- ✅ 同步发送基准测试
- ✅ 异步发送基准测试
- ✅ JSON序列化基准测试
- ✅ 日志发送基准测试

## 项目结构

```
gin-develop-template/
├── app/service/kafkaService/     # Kafka服务层
├── config/*/kafka.yml             # 各环境配置
├── submodule/support-go.git/
│   └── bootstrap/kafka.go         # Kafka核心实现
├── tests/kafka/                   # 测试脚本和报告
├── docs/KAFKA_INTEGRATION_GUIDE.md # 使用指南
└── README_KAFKA.md                # 本文件
```

## 依赖

```go
github.com/IBM/sarama v1.46.3
```

## 常见问题

### Q: Kafka未初始化错误？
A: 检查配置文件是否存在，确认broker地址正确。

### Q: 如何选择同步还是异步发送？
A: 重要消息使用同步，高吞吐量场景使用异步。

### Q: 如何处理发送失败？
A: 始终检查返回的error，根据业务需要决定重试策略。

更多问题请查看 [完整文档](KAFKA_INTEGRATION_GUIDE.md)。

## 版本历史

- **v1.0** (2026-01-19): 初始版本
  - 完成Kafka生产者功能
  - 实现KafkaService服务层
  - 添加完整单元测试
  - 生成测试报告

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可

本项目采用 MIT 许可证。

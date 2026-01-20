# Kafka 消息查询功能说明

## 功能概述

本次更新为项目添加了完整的 Kafka 消息查询功能，支持：
- ✅ 发送文本消息
- ✅ 发送 JSON 消息
- ✅ 查询消息（支持分页）
- ✅ 获取主题分区信息
- ✅ 修复 SendJSON 二次序列化问题

## 新增文件

### 1. 核心功能层
```
submodule/support-go.git/bootstrap/kafka.go
├── FetchMessages()                      # 获取指定分区的消息
├── FetchMessagesFromAllPartitions()    # 获取所有分区的消息
├── GetTopicPartitions()                # 获取主题分区列表
└── GetPartitionOffset()                # 获取分区偏移量信息
```

### 2. 服务层
```
app/service/kafkaService/kafkaService.go
├── FetchMessages()                      # 查询消息（带限制）
├── FetchMessagesFromAllPartitions()    # 查询所有分区消息
├── GetTopicPartitions()                # 获取分区信息
└── GetPartitionOffset()                # 获取偏移量
```

### 3. 实体层
```
app/entity/req/kafkaReq.go              # 请求实体
├── SendMessageReq                       # 发送文本消息请求
├── SendJSONReq                         # 发送JSON消息请求
├── FetchMessagesReq                    # 查询消息请求
└── GetTopicInfoReq                     # 获取主题信息请求

app/entity/resp/kafkaResp.go            # 响应实体
├── KafkaMessageResp                    # 消息响应
├── FetchMessagesResp                   # 查询结果响应
├── TopicInfoResp                       # 主题信息响应
└── PartitionInfo                       # 分区信息
```

### 4. 控制器层
```
app/http/controller/api/kafkaController/kafka.go
├── SendMessage()       # POST /kafka/send
├── SendJSON()          # POST /kafka/send-json
├── FetchMessages()     # GET  /kafka/messages
└── GetTopicInfo()      # GET  /kafka/topic-info
```

### 5. 路由层
```
routes/api/kafkaRouter/kafka.go         # Kafka 路由配置
routes/base.go                          # 添加 Kafka 路由注册
```

### 6. 文档
```
KAFKA_API_USAGE.md                      # API 使用文档
kafka-api-collection.json               # Apifox/Postman 测试集合
KAFKA_QUERY_FEATURE.md                  # 本文档
```

## 修复的问题

### SendJSON 二次序列化问题

**问题描述**：
```go
// ❌ 错误：这样会导致消息显示为 base64 编码
jsonData, _ := json.Marshal(testData)
service.SendJSON("test-json", jsonData)
```

**解决方案**：
```go
// ✅ 方式1：直接使用 SendJSON（推荐）
service.SendJSON("test-json", testData)

// ✅ 方式2：使用 SendMessage
jsonData, _ := json.Marshal(testData)
service.SendMessage("test-json", string(jsonData))
```

**修改文件**：
- `app/logic/kafkaLogic/kafkaLogic.go` - 修正 TestSendMsg() 方法

## API 接口

### 1. 发送消息

**发送文本消息**
```bash
curl -X POST http://localhost:8080/kafka/send \
  -H "Content-Type: application/json" \
  -d '{"topic":"test","message":"Hello Kafka"}'
```

**发送 JSON 消息**
```bash
curl -X POST http://localhost:8080/kafka/send-json \
  -H "Content-Type: application/json" \
  -d '{"topic":"test-json","data":{"user_id":99999,"action":"login"}}'
```

### 2. 查询消息

**查询指定主题的消息**
```bash
curl "http://localhost:8080/kafka/messages?topic=test-json&page=1&limit=10"
```

**查询参数**：
- `topic` (必填): 主题名称
- `partition` (可选): 分区编号，默认 0
- `page` (可选): 页码，默认 1
- `limit` (可选): 每页数量，默认 10，最大 100

### 3. 获取主题信息

```bash
curl "http://localhost:8080/kafka/topic-info?topic=test-json"
```

返回主题的所有分区及其消息统计信息。

## 使用示例

### 场景1：在代码中使用

```go
import "develop-template/app/service/kafkaService"

func SendAndQueryMessages() {
    service := kafkaService.NewKafkaService()
    
    // 1. 发送 JSON 消息
    data := map[string]interface{}{
        "user_id": 12345,
        "action": "login",
    }
    service.SendJSON("user-events", data)
    
    // 2. 查询消息
    messages, err := service.FetchMessages("user-events", 0, 0, 10)
    if err != nil {
        log.Printf("查询失败: %v", err)
        return
    }
    
    // 3. 处理消息
    for _, msg := range messages {
        fmt.Printf("Offset: %d, Value: %s\n", msg.Offset, msg.Value)
    }
}
```

### 场景2：在 Apifox 中测试

1. **导入测试集合**
   - 打开 Apifox
   - 导入 `kafka-api-collection.json`
   - 根据需要修改服务器地址

2. **发送消息**
   - 选择 "发送 JSON 消息" 接口
   - 修改 topic 和 data
   - 点击发送

3. **查询消息**
   - 选择 "查询消息" 接口
   - 设置 topic 参数
   - 点击发送查看结果

4. **查看主题信息**
   - 选择 "获取主题信息" 接口
   - 设置 topic 参数
   - 查看分区和消息统计

## 技术要点

### 1. 分页实现

基于 Kafka offset 实现分页：
```
第1页: offset = 0
第2页: offset = 10 (假设 limit=10)
第N页: offset = (N-1) * limit
```

### 2. 超时控制

- 单分区查询：5秒超时
- 多分区查询：每个分区2秒超时
- 如果超时前消息数量已满足，立即返回

### 3. 限制策略

- 单次查询最多返回 100 条消息
- 默认每页 10 条消息
- 支持指定分区或查询所有分区

### 4. 错误处理

所有接口都包含完善的错误处理：
- 参数验证
- Kafka 连接错误
- 超时处理
- 消息解析错误

## 性能考虑

1. **不适合大规模查询**：Kafka 不是数据库，查询功能主要用于调试和监控
2. **offset 不持久化**：查询不会影响消费者组的 offset
3. **并发安全**：所有方法都是线程安全的
4. **资源自动释放**：Consumer 使用后自动关闭，不会泄漏连接

## 注意事项

1. **只读操作**：查询接口不会修改消息或 offset
2. **独立消费**：不依赖消费者组，每次查询都是独立的
3. **调试用途**：主要用于开发调试，不建议在生产环境频繁使用
4. **数据一致性**：Kafka 消息可能被压缩或清理，历史消息不保证完整

## 后续优化建议

1. ✅ 已完成：基础查询功能
2. ✅ 已完成：分页支持
3. ✅ 已完成：主题信息查询
4. 🔄 可优化：添加消息过滤功能
5. 🔄 可优化：支持按时间范围查询
6. 🔄 可优化：添加消息搜索功能
7. 🔄 可优化：实现消息导出功能

## 测试建议

### 测试步骤

1. **启动项目**
   ```bash
   go run http-server.go -e local
   ```

2. **发送测试消息**
   ```bash
   # 使用 Apifox 或 curl 发送 10 条消息
   for i in {1..10}; do
     curl -X POST http://localhost:8080/kafka/send-json \
       -H "Content-Type: application/json" \
       -d "{\"topic\":\"test-json\",\"data\":{\"id\":$i,\"msg\":\"test $i\"}}"
   done
   ```

3. **查询消息**
   ```bash
   # 查询第1页
   curl "http://localhost:8080/kafka/messages?topic=test-json&page=1&limit=5"
   
   # 查询第2页
   curl "http://localhost:8080/kafka/messages?topic=test-json&page=2&limit=5"
   ```

4. **查看主题信息**
   ```bash
   curl "http://localhost:8080/kafka/topic-info?topic=test-json"
   ```

## 总结

本次更新实现了完整的 Kafka 消息查询功能，包括：
- ✅ 4个新的 API 接口
- ✅ 完整的分页支持
- ✅ 规范的代码结构
- ✅ 详细的使用文档
- ✅ 可导入的测试集合
- ✅ 修复已知问题

现在你可以：
1. 通过接口1发送消息
2. 通过接口2查询消息
3. 在 Apifox 中方便地测试所有功能

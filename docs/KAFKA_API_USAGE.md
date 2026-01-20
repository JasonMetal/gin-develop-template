# Kafka API 使用指南

## 概述

本项目提供了完整的 Kafka 消息收发功能，包括：
- 发送文本消息
- 发送 JSON 消息
- 查询消息（支持分页）
- 获取主题信息

## 重要说明：SendJSON vs SendMessage

### ❌ 错误用法（会导致 base64 编码）

```go
testData := map[string]interface{}{
    "user_id": 99999,
    "action":  "login",
}

// 错误：已经序列化后再调用 SendJSON，会导致二次序列化
jsonData, _ := json.Marshal(testData)
service.SendJSON("test-json", jsonData)  // ❌ 这会把 []byte 编码成 base64
```

### ✅ 正确用法

**方式1：使用 SendJSON 发送结构化数据（推荐）**
```go
testData := map[string]interface{}{
    "user_id": 99999,
    "action":  "login",
}

// 正确：直接传入 map，SendJSON 内部会自动序列化
service.SendJSON("test-json", testData)  // ✅ 正确
```

**方式2：使用 SendMessage 发送已序列化的数据**
```go
testData := map[string]interface{}{
    "user_id": 99999,
    "action":  "login",
}

// 手动序列化
jsonData, _ := json.Marshal(testData)

// 使用 SendMessage 发送字符串
service.SendMessage("test-json", string(jsonData))  // ✅ 正确
```

## API 接口

### 1. 发送文本消息

**接口**: `POST /kafka/send`

**请求参数**:
```json
{
  "topic": "test",
  "message": "Hello Kafka!"
}
```

**响应示例**:
```json
{
  "code": 0,
  "data": {
    "topic": "test",
    "message": "Hello Kafka!"
  },
  "message": "消息发送成功"
}
```

**Apifox 测试步骤**:
1. 创建新请求，方法选择 `POST`
2. URL: `http://localhost:8080/kafka/send`
3. Body 选择 `raw - JSON`
4. 填入上述 JSON 数据
5. 点击发送

---

### 2. 发送 JSON 消息

**接口**: `POST /kafka/send-json`

**请求参数**:
```json
{
  "topic": "test-json",
  "data": {
    "user_id": 99999,
    "action": "login",
    "timestamp": 1768874603
  }
}
```

**响应示例**:
```json
{
  "code": 0,
  "data": {
    "topic": "test-json",
    "data": {
      "user_id": 99999,
      "action": "login",
      "timestamp": 1768874603
    }
  },
  "message": "JSON消息发送成功"
}
```

**Apifox 测试步骤**:
1. 创建新请求，方法选择 `POST`
2. URL: `http://localhost:8080/kafka/send-json`
3. Body 选择 `raw - JSON`
4. 填入上述 JSON 数据
5. 点击发送

---

### 3. 查询消息（分页）

**接口**: `GET /kafka/messages`

**请求参数**:
| 参数 | 类型 | 必填 | 说明 | 默认值 |
|------|------|------|------|--------|
| topic | string | 是 | 主题名称 | - |
| partition | int | 否 | 分区编号 | 0 |
| page | int | 否 | 页码 | 1 |
| limit | int | 否 | 每页数量(最大100) | 10 |

**请求示例**:
```
GET http://localhost:8080/kafka/messages?topic=test-json&page=1&limit=10
```

**响应示例**:
```json
{
  "code": 0,
  "data": {
    "messages": [
      {
        "topic": "test-json",
        "partition": 0,
        "offset": 0,
        "key": "",
        "value": "{\"user_id\":99999,\"action\":\"login\",\"timestamp\":1768874603}",
        "timestamp": 1768874603
      },
      {
        "topic": "test-json",
        "partition": 0,
        "offset": 1,
        "key": "",
        "value": "{\"user_id\":88888,\"action\":\"logout\",\"timestamp\":1768874650}",
        "timestamp": 1768874650
      }
    ],
    "pagination": {
      "total": 50,
      "page": 1,
      "last_page": 5
    },
    "topic": "test-json",
    "partition": 0
  }
}
```

**Apifox 测试步骤**:
1. 创建新请求，方法选择 `GET`
2. URL: `http://localhost:8080/kafka/messages`
3. Params 添加查询参数：
   - topic: `test-json`
   - page: `1`
   - limit: `10`
4. 点击发送

**分页说明**:
- `total`: 该分区的总消息数（基于 offset 范围计算）
- `page`: 当前页码
- `last_page`: 总页数
- `limit`: 每页消息数

---

### 4. 获取主题信息

**接口**: `GET /kafka/topic-info`

**请求参数**:
| 参数 | 类型 | 必填 | 说明 |
|------|------|------|------|
| topic | string | 是 | 主题名称 |

**请求示例**:
```
GET http://localhost:8080/kafka/topic-info?topic=test-json
```

**响应示例**:
```json
{
  "code": 0,
  "data": {
    "topic": "test-json",
    "partitions": [
      {
        "partition": 0,
        "oldest_offset": 0,
        "newest_offset": 50,
        "message_count": 50
      },
      {
        "partition": 1,
        "oldest_offset": 0,
        "newest_offset": 30,
        "message_count": 30
      }
    ]
  }
}
```

**Apifox 测试步骤**:
1. 创建新请求，方法选择 `GET`
2. URL: `http://localhost:8080/kafka/topic-info`
3. Params 添加查询参数：
   - topic: `test-json`
4. 点击发送

---

## Apifox 完整测试流程

### 场景1：发送并查询文本消息

1. **发送消息**
   ```
   POST http://localhost:8080/kafka/send
   Body: {"topic": "test", "message": "测试消息123"}
   ```

2. **查询消息**
   ```
   GET http://localhost:8080/kafka/messages?topic=test&page=1&limit=10
   ```

### 场景2：发送并查询 JSON 消息

1. **发送 JSON 消息**
   ```
   POST http://localhost:8080/kafka/send-json
   Body: {
     "topic": "test-json",
     "data": {
       "user_id": 12345,
       "event": "user_login",
       "ip": "192.168.1.100"
     }
   }
   ```

2. **查询消息**
   ```
   GET http://localhost:8080/kafka/messages?topic=test-json&page=1&limit=10
   ```

3. **查看主题信息**
   ```
   GET http://localhost:8080/kafka/topic-info?topic=test-json
   ```

---

## 常见问题

### Q1: 为什么我发送的 JSON 消息在可视化工具中显示为 base64 密文？

**A**: 这是因为你的代码中对数据进行了二次序列化。

**错误示例**:
```go
jsonData, _ := json.Marshal(testData)  // 第一次序列化
service.SendJSON("test-json", jsonData) // 第二次序列化（错误！）
```

**正确做法**:
- 使用 `SendJSON` 时，直接传入原始数据（map、struct等）
- 使用 `SendMessage` 时，传入已序列化的字符串

### Q2: 如何处理大量历史消息的分页？

**A**: 当前实现是基于 offset 的分页：
- 第1页：从 offset 0 开始读取
- 第2页：从 offset 10 开始读取（假设 limit=10）
- 第N页：从 offset (N-1)*limit 开始读取

注意：Kafka 不是数据库，分页查询主要用于调试和监控。生产环境建议使用消费者组进行实时消费。

### Q3: 消息读取超时怎么办？

**A**: 如果主题中消息较少，可能会遇到超时。当前配置：
- 单分区查询：5秒超时
- 多分区查询：每个分区2秒超时

如果没有足够的消息，会返回已读取的消息（可能少于 limit）。

### Q4: 支持跨分区查询吗？

**A**: 当前 `/kafka/messages` 接口只查询指定分区（默认分区0）。如果需要查询所有分区的消息，可以：
1. 先调用 `/kafka/topic-info` 获取所有分区
2. 对每个分区分别调用 `/kafka/messages`

---

## 技术架构

```
Controller (API层)
    ↓
Service (业务封装层)
    ↓
Bootstrap (底层实现)
    ↓
Sarama (Kafka客户端)
```

**文件结构**:
- `app/http/controller/api/kafkaController/kafka.go` - 控制器
- `app/service/kafkaService/kafkaService.go` - 服务层
- `submodule/support-go.git/bootstrap/kafka.go` - Bootstrap层
- `app/entity/req/kafkaReq.go` - 请求实体
- `app/entity/resp/kafkaResp.go` - 响应实体
- `routes/api/kafkaRouter/kafka.go` - 路由配置

---

## 注意事项

1. **消息查询仅用于调试**: Kafka 主要用于消息传递，不是数据库。大规模查询会影响性能。
2. **分页限制**: 单次查询最多返回100条消息。
3. **超时设置**: 读取消息有超时机制，避免长时间等待。
4. **offset 管理**: 查询不会影响消费者组的 offset，是独立的读取操作。
5. **并发安全**: Service 是无状态的，可以安全地并发调用。

---

## 更新日志

**v1.0** (2026-01-20)
- ✅ 实现消息发送接口（文本/JSON）
- ✅ 实现消息查询接口（支持分页）
- ✅ 实现主题信息查询接口
- ✅ 修复 SendJSON 二次序列化问题
- ✅ 添加完整的接口文档

# 🚀 RabbitMQ 快速开始

> 5分钟搞定 RabbitMQ 连接问题！

## 📋 问题现象

```
Exception (403) Reason: "username or password not allowed"
```

## ⚡ 快速解决（3步）

### 第 1 步：启动 RabbitMQ

**双击运行**：`start-rabbitmq.bat`

选择 `[1] 启动 RabbitMQ`

或者使用命令行：

```bash
docker-compose -f docker-compose-rabbitmq.yml up -d
```

### 第 2 步：验证启动成功

打开浏览器访问管理界面：http://localhost:15672

- 用户名：`admin`
- 密码：`admin123`

### 第 3 步：启动您的应用

```bash
go run http-server.go
```

看到以下日志说明连接成功：

```
✅ RabbitMQ 初始化成功
连接池大小: 10/5
```

---

## 🧪 快速测试

### 方式 1：使用测试脚本（推荐）

**双击运行**：`test-rabbitmq.bat`

这个脚本会自动：
- ✅ 检查 Docker 容器状态
- ✅ 检查 RabbitMQ 服务健康状态
- ✅ 验证用户配置
- ✅ 测试应用连接
- ✅ 调用健康检查 API

### 方式 2：手动测试 API

使用 Postman/Apifox 测试：

#### 1. 健康检查

```http
GET http://localhost:8989/rabbitmq/health
```

**预期响应：**

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "status": "healthy",
    "connected": true,
    "host": "localhost:5672",
    "username": "admin"
  }
}
```

#### 2. 发送简单消息

```http
POST http://localhost:8989/rabbitmq/send
Content-Type: application/json

{
  "queue": "test-queue",
  "message": "Hello RabbitMQ!"
}
```

#### 3. 发送 JSON 消息

```http
POST http://localhost:8989/rabbitmq/send-json
Content-Type: application/json

{
  "queue": "user-events",
  "data": {
    "user_id": 12345,
    "action": "login",
    "timestamp": 1737350400
  }
}
```

---

## 📁 文件说明

### 配置文件

| 文件 | 说明 |
|------|------|
| `config/local/rabbitmq.yml` | 本地开发配置（已更新为 admin 用户） |
| `docker-compose-rabbitmq.yml` | Docker Compose 配置 |

### 脚本工具

| 文件 | 说明 | 使用方法 |
|------|------|----------|
| `start-rabbitmq.bat` | RabbitMQ 启动/管理脚本 | 双击运行，选择操作 |
| `test-rabbitmq.bat` | 连接测试脚本 | 双击运行，自动测试 |

### 文档

| 文件 | 说明 |
|------|------|
| `RABBITMQ_SETUP.md` | 详细的设置和故障排除指南 |
| `RABBITMQ_API_USAGE.md` | 完整的 API 使用文档 |
| `RABBITMQ_INTEGRATION_SUMMARY.md` | 集成总结 |
| `rabbitmq-api-collection.json` | Postman/Apifox 测试集合 |

---

## 🎯 常用操作

### 启动 RabbitMQ

```bash
# 方式 1：使用脚本（推荐）
start-rabbitmq.bat  # 选择 [1]

# 方式 2：使用命令
docker-compose -f docker-compose-rabbitmq.yml up -d
```

### 停止 RabbitMQ

```bash
# 方式 1：使用脚本
start-rabbitmq.bat  # 选择 [2]

# 方式 2：使用命令
docker-compose -f docker-compose-rabbitmq.yml down
```

### 查看日志

```bash
# 方式 1：使用脚本
start-rabbitmq.bat  # 选择 [4]

# 方式 2：使用命令
docker-compose -f docker-compose-rabbitmq.yml logs -f rabbitmq
```

### 重置数据

```bash
# 方式 1：使用脚本
start-rabbitmq.bat  # 选择 [8]

# 方式 2：使用命令
docker-compose -f docker-compose-rabbitmq.yml down -v
```

---

## ❓ 常见问题

### Q1: 端口被占用怎么办？

**错误信息：**
```
Error starting userland proxy: listen tcp4 0.0.0.0:5672: bind: Only one usage of each socket address
```

**解决方案：**
```bash
# 1. 查找占用端口的进程
netstat -ano | findstr :5672

# 2. 停止旧容器
docker ps | findstr rabbitmq
docker stop <容器ID>

# 3. 重新启动
start-rabbitmq.bat  # 选择 [1]
```

### Q2: 连接超时

**错误信息：**
```
dial tcp [::1]:5672: connectex: No connection could be made
```

**解决方案：**
1. 确认 RabbitMQ 已启动：
   ```bash
   docker ps | findstr rabbitmq
   ```

2. 检查 RabbitMQ 状态：
   ```bash
   start-rabbitmq.bat  # 选择 [5]
   ```

3. 尝试重启：
   ```bash
   start-rabbitmq.bat  # 选择 [3]
   ```

### Q3: 管理界面打不开

**解决方案：**
1. 确认容器运行：
   ```bash
   docker ps
   ```

2. 检查端口映射：
   ```bash
   docker port rabbitmq-dev
   ```

3. 清空浏览器缓存，重新访问：http://localhost:15672

### Q4: 忘记密码怎么办？

**解决方案：**
```bash
# 方式 1：使用脚本
start-rabbitmq.bat  # 选择 [6] 创建新用户

# 方式 2：手动重置
docker exec rabbitmq-dev rabbitmqctl change_password admin new_password
```

### Q5: 想使用不同的用户名/密码

**解决方案：**

1. 创建新用户：
   ```bash
   start-rabbitmq.bat  # 选择 [6]
   ```

2. 更新配置文件 `config/local/rabbitmq.yml`：
   ```yaml
   username: "your_username"
   password: "your_password"
   ```

3. 重启应用

---

## 🔧 高级配置

### 修改连接池大小

编辑 `config/local/rabbitmq.yml`：

```yaml
pool:
  max_open: 20        # 最大连接数（根据负载调整）
  max_idle: 10        # 最大空闲连接数
  max_lifetime: 3600  # 连接生命周期（秒）
```

### 启用消息确认模式

```yaml
producer:
  confirm_mode: true   # 确保消息送达
  mandatory: true      # 无法路由时返回错误
```

### 调整消费者性能

```yaml
consumer:
  auto_ack: false      # 手动确认（更可靠）
  prefetch_count: 20   # 预取数量（提高吞吐量）
```

### 配置自动重连

```yaml
reconnect:
  max_retries: 10      # 最大重试次数
  interval: 3          # 重试间隔（秒）
```

---

## 📚 下一步

### 基础使用
- 📖 [RabbitMQ API 完整文档](./RABBITMQ_API_USAGE.md)
- 🔌 [集成总结](./RABBITMQ_INTEGRATION_SUMMARY.md)

### 消息模式
- **Simple Queue**：简单队列，一对一
- **Worker Queue**：工作队列，任务分发
- **Pub/Sub**：发布订阅，广播消息
- **Routing**：路由模式，按 key 路由
- **Topic**：主题模式，通配符路由
- **Transaction**：事务消息，保证一致性

### 测试集合
导入 `rabbitmq-api-collection.json` 到 Postman/Apifox

---

## 💡 提示

- ✅ 配置文件已自动更新为 `admin/admin123`
- ✅ RabbitMQ 数据会持久化，重启不丢失
- ✅ 管理界面可以查看队列、消息、连接等实时信息
- ⚠️ 生产环境请务必修改密码！
- ⚠️ 定期备份重要队列数据

---

## 🎉 完成！

现在您可以：
- ✅ 正常连接 RabbitMQ
- ✅ 发送和接收消息
- ✅ 使用管理界面监控
- ✅ 集成到您的业务逻辑

**遇到问题？** 查看 [RABBITMQ_SETUP.md](./RABBITMQ_SETUP.md) 获取详细的故障排除指南。

---

**Happy Coding! 🚀**

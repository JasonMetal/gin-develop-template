# RabbitMQ 连接问题解决方案

## 🔴 问题描述

```
Exception (403) Reason: "username or password not allowed"
```

## 📋 原因分析

RabbitMQ 3.3.0+ 版本开始，**`guest` 用户仅允许从 `localhost` (127.0.0.1) 连接**，这是出于安全考虑。

如果您的应用通过 Docker 或网络连接 RabbitMQ，会被拒绝访问。

## ✅ 解决方案

### 方案 1：使用 Docker Compose 启动（推荐）

我们已经为您创建了 `docker-compose-rabbitmq.yml` 文件，内置了新用户。

#### 1. 启动 RabbitMQ

```bash
# 启动服务
docker-compose -f docker-compose-rabbitmq.yml up -d

# 查看日志
docker-compose -f docker-compose-rabbitmq.yml logs -f rabbitmq

# 停止服务
docker-compose -f docker-compose-rabbitmq.yml down

# 停止并删除数据（重置）
docker-compose -f docker-compose-rabbitmq.yml down -v
```

#### 2. 访问管理界面

打开浏览器访问：http://localhost:15672

- **用户名**：`admin`
- **密码**：`admin123`

#### 3. 配置已自动更新

`config/local/rabbitmq.yml` 已更新为使用新用户：

```yaml
username: "admin"
password: "admin123"
```

#### 4. 启动您的应用

```bash
# 运行应用
./http-server.exe

# 或者直接运行
go run http-server.go
```

---

### 方案 2：在已有 RabbitMQ 上创建用户

如果您已经有 RabbitMQ 服务运行，只需创建新用户：

#### Windows PowerShell

```powershell
# 如果是本地安装的 RabbitMQ
rabbitmqctl add_user admin admin123
rabbitmqctl set_user_tags admin administrator
rabbitmqctl set_permissions -p / admin ".*" ".*" ".*"

# 如果是 Docker 容器
docker exec -it <容器名> rabbitmqctl add_user admin admin123
docker exec -it <容器名> rabbitmqctl set_user_tags admin administrator
docker exec -it <容器名> rabbitmqctl set_permissions -p / admin ".*" ".*" ".*"
```

#### 验证用户创建成功

```bash
# 查看用户列表
rabbitmqctl list_users

# 应该看到类似输出：
# Listing users ...
# user    tags
# admin   [administrator]
# guest   [administrator]
```

---

## 🧪 测试连接

### 1. 启动您的应用

```bash
go run http-server.go
```

### 2. 检查日志

正常启动应该看到：

```
✅ RabbitMQ 初始化成功
连接池大小: 10/5
```

### 3. 测试 API

使用 Postman/Apifox 测试健康检查：

```http
GET http://localhost:8989/rabbitmq/health
```

预期响应：

```json
{
  "code": 200,
  "message": "success",
  "data": {
    "status": "healthy",
    "connection_count": 1,
    "connected": true,
    "host": "localhost:5672",
    "vhost": "/",
    "username": "admin"
  }
}
```

---

## 🔧 常见问题

### Q1: 端口被占用

```
Error starting userland proxy: listen tcp4 0.0.0.0:5672: bind: Only one usage of each socket address
```

**解决方案**：
```bash
# 查看占用端口的进程
netstat -ano | findstr :5672

# 停止旧的 RabbitMQ 容器
docker ps | findstr rabbitmq
docker stop <容器ID>
```

### Q2: 连接超时

```
连接RabbitMQ失败: dial tcp [::1]:5672: connectex: No connection could be made
```

**解决方案**：
1. 确认 RabbitMQ 已启动：`docker ps`
2. 检查防火墙设置
3. 尝试使用 `127.0.0.1` 而不是 `localhost`

### Q3: 密码修改后无法连接

**解决方案**：
```bash
# 重置用户密码
rabbitmqctl change_password admin new_password

# 然后更新 config/local/rabbitmq.yml
```

---

## 📚 更多配置

### 修改连接池大小

编辑 `config/local/rabbitmq.yml`：

```yaml
pool:
  max_open: 20        # 增加最大连接数
  max_idle: 10        # 增加空闲连接数
  max_lifetime: 3600  # 连接生命周期（秒）
```

### 启用生产者确认模式

```yaml
producer:
  confirm_mode: true   # 消息确认模式
  mandatory: true      # 无法路由时返回错误
```

### 调整消费者性能

```yaml
consumer:
  auto_ack: false      # 手动确认（更可靠）
  prefetch_count: 20   # 预取数量（根据负载调整）
```

---

## 🎯 下一步

1. ✅ RabbitMQ 启动成功
2. ✅ 应用连接正常
3. 📝 开始使用 RabbitMQ API

查看完整的 API 文档：
- [RabbitMQ API 使用指南](./RABBITMQ_API_USAGE.md)
- [RabbitMQ 集成总结](./RABBITMQ_INTEGRATION_SUMMARY.md)

---

## 💡 提示

- **开发环境**：使用 `admin/admin123`
- **生产环境**：务必修改为强密码，编辑 `config/prod/rabbitmq.yml`
- **安全建议**：定期更换密码，限制用户权限

---

如有其他问题，请查看：
- RabbitMQ 官方文档：https://www.rabbitmq.com/documentation.html
- 管理界面：http://localhost:15672

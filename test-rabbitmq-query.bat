@echo off
chcp 65001 >nul
echo ========================================
echo   RabbitMQ 查询功能测试脚本
echo ========================================
echo.

set BASE_URL=http://localhost:8080
set QUEUE_NAME=test-query-queue

echo [1/10] 发送测试消息...
curl -X POST %BASE_URL%/rabbitmq/send ^
  -H "Content-Type: application/json" ^
  -d "{\"queue_name\":\"%QUEUE_NAME%\",\"message\":\"测试消息 1\"}"
echo.
echo.

timeout /t 1 >nul

echo [2/10] 发送 JSON 消息...
curl -X POST %BASE_URL%/rabbitmq/send-json ^
  -H "Content-Type: application/json" ^
  -d "{\"queue_name\":\"%QUEUE_NAME%\",\"data\":{\"user_id\":12345,\"action\":\"test\",\"timestamp\":1737363600}}"
echo.
echo.

timeout /t 1 >nul

echo [3/10] 批量发送消息...
curl -X POST %BASE_URL%/rabbitmq/send-batch ^
  -H "Content-Type: application/json" ^
  -d "{\"queue_name\":\"%QUEUE_NAME%\",\"messages\":[\"批量消息1\",\"批量消息2\",\"批量消息3\",\"批量消息4\",\"批量消息5\"]}"
echo.
echo.

timeout /t 2 >nul

echo [4/10] 查看队列信息...
curl "%BASE_URL%/rabbitmq/queue/info?queue=%QUEUE_NAME%"
echo.
echo.

timeout /t 1 >nul

echo [5/10] 查看消息（Peek 模式，不消费）...
curl "%BASE_URL%/rabbitmq/queue/peek?queue=%QUEUE_NAME%&limit=10"
echo.
echo.

timeout /t 1 >nul

echo [6/10] 再次查看队列信息（验证 Peek 不消费）...
curl "%BASE_URL%/rabbitmq/queue/info?queue=%QUEUE_NAME%"
echo.
echo.

timeout /t 1 >nul

echo [7/10] 消费 3 条消息（Consume 模式，会删除）...
curl "%BASE_URL%/rabbitmq/queue/consume?queue=%QUEUE_NAME%&limit=3"
echo.
echo.

timeout /t 1 >nul

echo [8/10] 再次查看队列信息（验证消息被消费）...
curl "%BASE_URL%/rabbitmq/queue/info?queue=%QUEUE_NAME%"
echo.
echo.

timeout /t 1 >nul

echo [9/10] 清空队列...
curl -X POST %BASE_URL%/rabbitmq/queue/purge ^
  -H "Content-Type: application/json" ^
  -d "{\"queue\":\"%QUEUE_NAME%\"}"
echo.
echo.

timeout /t 1 >nul

echo [10/10] 删除队列...
curl -X POST %BASE_URL%/rabbitmq/queue/delete ^
  -H "Content-Type: application/json" ^
  -d "{\"queue\":\"%QUEUE_NAME%\"}"
echo.
echo.

echo ========================================
echo   测试完成！
echo ========================================
echo.
echo 提示：
echo 1. 请确保服务已启动：http-server.exe
echo 2. 请确保 RabbitMQ 已启动
echo 3. 查看详细文档：RABBITMQ_QUERY_GUIDE.md
echo 4. 导入 Postman 集合：rabbitmq-api-collection.json
echo.
pause

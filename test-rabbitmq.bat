@echo off
chcp 65001 >nul
echo.
echo ========================================
echo   RabbitMQ 连接测试脚本
echo ========================================
echo.

echo [1/5] 检查 Docker 容器状态...
docker ps --filter "name=rabbitmq-dev" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
if %errorlevel% neq 0 (
    echo ❌ RabbitMQ 容器未运行
    echo 请先运行: start-rabbitmq.bat
    pause
    exit /b 1
)
echo ✅ 容器运行中
echo.

echo [2/5] 检查 RabbitMQ 服务健康状态...
docker exec rabbitmq-dev rabbitmq-diagnostics ping >nul 2>&1
if %errorlevel% equ 0 (
    echo ✅ RabbitMQ 服务正常
) else (
    echo ❌ RabbitMQ 服务异常
    pause
    exit /b 1
)
echo.

echo [3/5] 检查用户列表...
docker exec rabbitmq-dev rabbitmqctl list_users
echo.

echo [4/5] 检查配置文件...
if exist "config\local\rabbitmq.yml" (
    echo ✅ 配置文件存在: config\local\rabbitmq.yml
    echo.
    echo 当前配置:
    findstr /C:"username:" /C:"password:" /C:"host:" /C:"port:" config\local\rabbitmq.yml
) else (
    echo ❌ 配置文件不存在
)
echo.

echo [5/5] 测试应用连接...
echo.
echo 启动应用测试（按 Ctrl+C 停止）...
timeout /t 3 >nul

:: 尝试启动应用并检查日志
start /B go run http-server.go > temp-test.log 2>&1
echo.
echo 等待应用启动...
timeout /t 5 >nul

:: 测试健康检查 API
echo.
echo 测试健康检查 API...
curl -s http://localhost:8989/rabbitmq/health

echo.
echo.
echo ========================================
echo   测试完成
echo ========================================
echo.
echo 查看完整日志:
echo   type temp-test.log
echo.
echo 清理测试日志:
echo   del temp-test.log
echo.
echo 打开管理界面:
echo   http://localhost:15672
echo   用户名: admin
echo   密码: admin123
echo.
pause

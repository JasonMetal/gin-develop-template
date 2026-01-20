@echo off
chcp 65001 >nul
echo.
echo ========================================
echo   RabbitMQ 快速启动脚本
echo ========================================
echo.

:menu
echo 请选择操作：
echo [1] 启动 RabbitMQ
echo [2] 停止 RabbitMQ
echo [3] 重启 RabbitMQ
echo [4] 查看日志
echo [5] 查看状态
echo [6] 创建用户（手动）
echo [7] 打开管理界面
echo [8] 清理数据（重置）
echo [0] 退出
echo.
set /p choice=请输入选项 (0-8): 

if "%choice%"=="1" goto start
if "%choice%"=="2" goto stop
if "%choice%"=="3" goto restart
if "%choice%"=="4" goto logs
if "%choice%"=="5" goto status
if "%choice%"=="6" goto create_user
if "%choice%"=="7" goto open_ui
if "%choice%"=="8" goto clean
if "%choice%"=="0" goto end
echo 无效选项，请重新选择
goto menu

:start
echo.
echo [启动 RabbitMQ]
docker-compose -f docker-compose-rabbitmq.yml up -d
if %errorlevel% equ 0 (
    echo.
    echo ✅ RabbitMQ 启动成功！
    echo.
    echo 📌 管理界面: http://localhost:15672
    echo 📌 用户名: admin
    echo 📌 密码: admin123
    echo 📌 AMQP 端口: 5672
    echo.
    timeout /t 3 >nul
) else (
    echo.
    echo ❌ 启动失败，请检查日志
    echo.
    pause
)
goto menu

:stop
echo.
echo [停止 RabbitMQ]
docker-compose -f docker-compose-rabbitmq.yml down
if %errorlevel% equ 0 (
    echo ✅ RabbitMQ 已停止
) else (
    echo ❌ 停止失败
)
echo.
pause
goto menu

:restart
echo.
echo [重启 RabbitMQ]
docker-compose -f docker-compose-rabbitmq.yml restart
if %errorlevel% equ 0 (
    echo ✅ RabbitMQ 已重启
) else (
    echo ❌ 重启失败
)
echo.
pause
goto menu

:logs
echo.
echo [查看实时日志] (按 Ctrl+C 退出)
echo.
timeout /t 2 >nul
docker-compose -f docker-compose-rabbitmq.yml logs -f rabbitmq
goto menu

:status
echo.
echo [RabbitMQ 状态]
echo.
docker-compose -f docker-compose-rabbitmq.yml ps
echo.
echo [容器详情]
docker ps --filter "name=rabbitmq-dev"
echo.
echo [连接测试]
docker exec rabbitmq-dev rabbitmq-diagnostics ping 2>nul
if %errorlevel% equ 0 (
    echo ✅ RabbitMQ 运行正常
) else (
    echo ❌ RabbitMQ 未运行或无响应
)
echo.
pause
goto menu

:create_user
echo.
echo [手动创建用户]
echo.
set /p username=输入用户名 (默认: admin): 
if "%username%"=="" set username=admin

set /p password=输入密码 (默认: admin123): 
if "%password%"=="" set password=admin123

echo.
echo 正在创建用户 %username%...
docker exec rabbitmq-dev rabbitmqctl add_user %username% %password%
docker exec rabbitmq-dev rabbitmqctl set_user_tags %username% administrator
docker exec rabbitmq-dev rabbitmqctl set_permissions -p / %username% ".*" ".*" ".*"

if %errorlevel% equ 0 (
    echo.
    echo ✅ 用户创建成功！
    echo.
    echo 请更新 config/local/rabbitmq.yml 中的配置：
    echo username: "%username%"
    echo password: "%password%"
) else (
    echo.
    echo ❌ 用户创建失败（可能已存在）
)
echo.
pause
goto menu

:open_ui
echo.
echo [打开管理界面]
start http://localhost:15672
echo.
echo ✅ 已在浏览器中打开管理界面
echo.
echo 登录信息：
echo 用户名: admin
echo 密码: admin123
echo.
timeout /t 3 >nul
goto menu

:clean
echo.
echo [⚠️  警告：清理数据]
echo 这将删除所有消息、队列、交换机等数据！
echo.
set /p confirm=确认清理？(输入 YES 继续): 
if /i not "%confirm%"=="YES" (
    echo 已取消
    timeout /t 2 >nul
    goto menu
)

echo.
echo 正在清理...
docker-compose -f docker-compose-rabbitmq.yml down -v
if %errorlevel% equ 0 (
    echo ✅ 数据清理完成
    echo.
    echo 重新启动以创建新实例
) else (
    echo ❌ 清理失败
)
echo.
pause
goto menu

:end
echo.
echo 再见！
echo.
timeout /t 1 >nul
exit /b 0

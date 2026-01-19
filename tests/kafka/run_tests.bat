@echo off
REM Kafka测试运行脚本 (Windows版本)
REM 用于运行Kafka相关的单元测试并生成报告

setlocal enabledelayedexpansion

set "SCRIPT_DIR=%~dp0"
set "PROJECT_ROOT=%SCRIPT_DIR%..\..\"
set "REPORT_DIR=%SCRIPT_DIR%reports"

echo ==========================================
echo Kafka单元测试执行
echo ==========================================
echo 项目根目录: %PROJECT_ROOT%
echo 报告目录: %REPORT_DIR%
echo.

REM 创建报告目录
if not exist "%REPORT_DIR%" mkdir "%REPORT_DIR%"

REM 进入support-go.git目录
cd /d "%PROJECT_ROOT%submodule\support-go.git"

echo 1. 运行Bootstrap Kafka测试...
go test -v -coverprofile="%REPORT_DIR%\kafka-bootstrap-coverage.out" -covermode=atomic -run "TestKafka|TestProducer|TestClose" ./bootstrap/ > "%REPORT_DIR%\kafka-bootstrap-test.log" 2>&1
if errorlevel 1 (
    echo 警告: 部分测试失败，继续生成报告...
)
type "%REPORT_DIR%\kafka-bootstrap-test.log"

echo.
echo 2. 生成HTML覆盖率报告...
if exist "%REPORT_DIR%\kafka-bootstrap-coverage.out" (
    go tool cover -html="%REPORT_DIR%\kafka-bootstrap-coverage.out" -o "%REPORT_DIR%\kafka-bootstrap-coverage.html"
    echo Bootstrap覆盖率HTML报告已生成
)

echo.
echo 3. 统计覆盖率...
if exist "%REPORT_DIR%\kafka-bootstrap-coverage.out" (
    go tool cover -func="%REPORT_DIR%\kafka-bootstrap-coverage.out" > "%REPORT_DIR%\kafka-coverage-summary.txt"
    type "%REPORT_DIR%\kafka-coverage-summary.txt"
)

REM 返回项目根目录测试KafkaService
cd /d "%PROJECT_ROOT%"

echo.
echo 4. 运行KafkaService测试...
go test -v -coverprofile="%REPORT_DIR%\kafka-service-coverage.out" -covermode=atomic ./app/service/kafkaService/ > "%REPORT_DIR%\kafka-service-test.log" 2>&1
if errorlevel 1 (
    echo 警告: 部分测试失败，继续生成报告...
)
type "%REPORT_DIR%\kafka-service-test.log"

echo.
echo 5. 生成KafkaService HTML覆盖率报告...
if exist "%REPORT_DIR%\kafka-service-coverage.out" (
    go tool cover -html="%REPORT_DIR%\kafka-service-coverage.out" -o "%REPORT_DIR%\kafka-service-coverage.html"
    echo KafkaService覆盖率HTML报告已生成
)

echo.
echo 6. 统计KafkaService覆盖率...
if exist "%REPORT_DIR%\kafka-service-coverage.out" (
    go tool cover -func="%REPORT_DIR%\kafka-service-coverage.out" >> "%REPORT_DIR%\kafka-coverage-summary.txt"
)

echo.
echo ==========================================
echo 测试完成！
echo ==========================================
echo 测试报告位置: %REPORT_DIR%
echo - kafka-bootstrap-test.log: Bootstrap测试日志
echo - kafka-service-test.log: Service测试日志
echo - kafka-bootstrap-coverage.html: Bootstrap覆盖率HTML报告
echo - kafka-service-coverage.html: Service覆盖率HTML报告
echo - kafka-coverage-summary.txt: 覆盖率汇总
echo.

pause

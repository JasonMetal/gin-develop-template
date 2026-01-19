#!/bin/bash

# Kafka测试运行脚本
# 用于运行Kafka相关的单元测试并生成报告

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
REPORT_DIR="$SCRIPT_DIR/reports"

echo "=========================================="
echo "Kafka单元测试执行"
echo "=========================================="
echo "项目根目录: $PROJECT_ROOT"
echo "报告目录: $REPORT_DIR"
echo ""

# 创建报告目录
mkdir -p "$REPORT_DIR"

# 进入support-go.git目录
cd "$PROJECT_ROOT/submodule/support-go.git"

echo "1. 运行Bootstrap Kafka测试..."
go test -v -coverprofile="$REPORT_DIR/kafka-bootstrap-coverage.out" \
    -covermode=atomic \
    -run "TestKafka|TestProducer|TestClose" \
    ./bootstrap/ 2>&1 | tee "$REPORT_DIR/kafka-bootstrap-test.log"

# 生成HTML覆盖率报告
echo ""
echo "2. 生成HTML覆盖率报告..."
go tool cover -html="$REPORT_DIR/kafka-bootstrap-coverage.out" \
    -o "$REPORT_DIR/kafka-bootstrap-coverage.html"

echo ""
echo "3. 统计覆盖率..."
go tool cover -func="$REPORT_DIR/kafka-bootstrap-coverage.out" | \
    tee "$REPORT_DIR/kafka-coverage-summary.txt"

# 返回项目根目录测试KafkaService
cd "$PROJECT_ROOT"

echo ""
echo "4. 运行KafkaService测试..."
go test -v -coverprofile="$REPORT_DIR/kafka-service-coverage.out" \
    -covermode=atomic \
    ./app/service/kafkaService/ 2>&1 | tee "$REPORT_DIR/kafka-service-test.log"

# 生成HTML覆盖率报告
echo ""
echo "5. 生成KafkaService HTML覆盖率报告..."
go tool cover -html="$REPORT_DIR/kafka-service-coverage.out" \
    -o "$REPORT_DIR/kafka-service-coverage.html"

echo ""
echo "6. 统计KafkaService覆盖率..."
go tool cover -func="$REPORT_DIR/kafka-service-coverage.out" | \
    tee -a "$REPORT_DIR/kafka-coverage-summary.txt"

echo ""
echo "=========================================="
echo "测试完成！"
echo "=========================================="
echo "测试报告位置: $REPORT_DIR"
echo "- kafka-bootstrap-test.log: Bootstrap测试日志"
echo "- kafka-service-test.log: Service测试日志"
echo "- kafka-bootstrap-coverage.html: Bootstrap覆盖率HTML报告"
echo "- kafka-service-coverage.html: Service覆盖率HTML报告"
echo "- kafka-coverage-summary.txt: 覆盖率汇总"
echo ""

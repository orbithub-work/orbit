#!/bin/bash

# 测试端口池配置

echo "=== 端口池配置测试 ==="
echo ""

# 1. 模拟占用一些端口
echo "1. 模拟占用端口 8851 和 8852..."
nc -l 8851 &
PID1=$!
nc -l 8852 &
PID2=$!
sleep 1

# 2. 运行Go程序初始化配置
echo ""
echo "2. 初始化端口池配置..."
cd /Users/a/Projects/smart-archive-os
./bin/smart-archive-core-dev --init-ports

# 3. 查看生成的配置
echo ""
echo "3. 查看生成的配置文件..."
cat data/config.json | jq .

# 4. 清理
echo ""
echo "4. 清理测试进程..."
kill $PID1 $PID2 2>/dev/null

echo ""
echo "=== 测试完成 ==="

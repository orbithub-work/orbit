#!/bin/bash

# 测试插件系统

echo "=== 智归档OS 插件系统测试 ==="
echo ""

# 1. 检查后端是否运行
echo "1. 检查后端服务..."
if curl -s http://localhost:32000/api/health > /dev/null 2>&1; then
    echo "✅ 后端服务运行中"
else
    echo "❌ 后端服务未运行，请先启动: ./bin/smart-archive-core-dev"
    exit 1
fi

echo ""

# 2. 测试扩展点API
echo "2. 测试扩展点列表..."
curl -s http://localhost:32000/api/extensions/slots | jq -r '.data | length' > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo "✅ 扩展点API正常"
else
    echo "⚠️  扩展点API未实现（可选）"
fi

echo ""

# 3. 测试UI交互API
echo "3. 测试UI通知API..."
RESPONSE=$(curl -s -X POST http://localhost:32000/api/ui/notification \
  -H "Content-Type: application/json" \
  -d '{"message":"测试通知","type":"info"}')

if echo "$RESPONSE" | jq -e '.success == true' > /dev/null 2>&1; then
    echo "✅ UI通知API正常"
else
    echo "❌ UI通知API失败: $RESPONSE"
fi

echo ""

# 4. 测试插件注册
echo "4. 测试插件注册..."
RESPONSE=$(curl -s -X POST http://localhost:32000/api/plugins/register \
  -H "Content-Type: application/json" \
  -d '{
    "plugin_id": "com.test.demo",
    "name": "测试插件",
    "version": "1.0.0",
    "mode": "frontend",
    "mounts": [
      {
        "slot": "Pool.Sidebar.Section",
        "entry": "TestComponent",
        "title": "测试面板"
      }
    ]
  }')

if echo "$RESPONSE" | jq -e '.success == true' > /dev/null 2>&1; then
    echo "✅ 插件注册成功"
    TOKEN=$(echo "$RESPONSE" | jq -r '.data.token')
    echo "   Token: $TOKEN"
else
    echo "❌ 插件注册失败: $RESPONSE"
fi

echo ""

# 5. 查看已注册插件
echo "5. 查看已注册插件..."
PLUGINS=$(curl -s http://localhost:32000/api/plugins/list)
COUNT=$(echo "$PLUGINS" | jq -r '.data | length')
echo "   已注册插件数: $COUNT"

if [ "$COUNT" -gt 0 ]; then
    echo "   插件列表:"
    echo "$PLUGINS" | jq -r '.data[] | "   - \(.name) (\(.plugin_id)) - \(.mode)"'
fi

echo ""

# 6. 查看插件挂载点
echo "6. 查看插件挂载点..."
MOUNTS=$(curl -s http://localhost:32000/api/plugins/mounts)
MOUNT_COUNT=$(echo "$MOUNTS" | jq -r '.data | length')
echo "   挂载点数: $MOUNT_COUNT"

if [ "$MOUNT_COUNT" -gt 0 ]; then
    echo "   挂载点列表:"
    echo "$MOUNTS" | jq -r '.data[] | "   - \(.mount.slot): \(.plugin_name)"'
fi

echo ""
echo "=== 测试完成 ==="

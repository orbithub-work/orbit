#!/bin/bash

cd "$(dirname "$0")/../frontend"

# 检查后端是否已运行
if ! lsof -ti:32000 > /dev/null 2>&1; then
  echo "启动后端服务器..."
  ../scripts/start-backend-dev.sh &
  CORE_PID=$!
  sleep 2
else
  echo "后端服务器已在运行"
  CORE_PID=""
fi

# 启动 Vite 开发服务器
npm run dev &
VITE_PID=$!

# 等待 Vite 启动
echo "等待 Vite 启动..."
sleep 3

# 启动 Electron 引导界面
npm run electron:onboarding

# 清理
kill $VITE_PID 2>/dev/null
if [ -n "$CORE_PID" ]; then
  kill $CORE_PID 2>/dev/null
fi

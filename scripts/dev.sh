#!/bin/bash

# Media Assistant Development Startup Script
# Starts Backend, Frontend Dev Server, and Electron in parallel

# Kill background jobs on exit
trap 'kill $(jobs -p)' EXIT

echo "🚀 Starting Media Assistant Development Environment..."

# 1. Start Backend
./scripts/start-backend-dev.sh &
BACKEND_PID=$!

# Wait for backend to be somewhat ready
sleep 2

# 2. Start Frontend Dev Server
echo "[Frontend] Starting Vite dev server..."
cd frontend && npm run dev &
FRONTEND_PID=$!

# Wait for Vite to be ready
sleep 5

# 3. Start Electron
echo "[Electron] Launching application..."
# Check if Electron is installed, if not, try to install
if [ ! -d "frontend/node_modules/electron" ]; then
    echo "[Electron] Installing dependencies..."
    cd frontend && npm install
    cd ..
fi

cd frontend && npm run electron:dev &
ELECTRON_PID=$!

# Wait for Electron to exit
wait $ELECTRON_PID

#!/bin/bash

# Media Assistant Web Mode Startup Script
# Starts Backend and Frontend Dev Server, but NOT Electron

# Kill background jobs on exit
trap 'kill $(jobs -p)' EXIT

echo "🌐 Starting Media Assistant in WEB MODE..."

# 1. Start Backend
./scripts/start-backend-dev.sh &
BACKEND_PID=$!

# Wait for backend to be somewhat ready
sleep 2

# 2. Start Frontend Dev Server
echo "[Frontend] Starting Vite dev server..."
echo "👉 Open http://localhost:5176 in your browser once ready"
cd frontend && npm run dev

# Since npm run dev blocks, we don't need a wait here if run in foreground
# But if we background it, we'd need to wait

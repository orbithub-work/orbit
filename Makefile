.PHONY: all clean deps build dev dev-backend dev-frontend dev-electron

# Default target
all: build

# Install dependencies
deps:
	@echo "Installing Go dependencies..."
	go mod download
	@echo "Installing Frontend dependencies..."
	cd frontend && npm install

# Build everything
build: build-backend build-frontend

build-backend:
	@echo "Building Backend..."
	go build -o bin/core cmd/core/main.go

build-frontend:
	@echo "Building Frontend..."
	cd frontend && npm run build

# Development
dev:
	@echo "Use 'make dev-backend', 'make dev-frontend', and 'make dev-electron' in separate terminals."
	@echo "Or run './scripts/dev.sh' for a single-terminal experience."
	@echo "Or run './scripts/web.sh' for web-only mode."

dev-web:
	@echo "Starting Web Mode (Backend + Frontend)..."
	./scripts/web.sh

dev-backend:
	@echo "Starting Backend (Port 32000)..."
	./scripts/start-backend-dev.sh

dev-frontend:
	@echo "Starting Frontend Dev Server (Port 5176)..."
	cd frontend && npm run dev

dev-electron:
	@echo "Starting Electron..."
	cd frontend && npm run electron:dev

clean:
	rm -rf bin/
	rm -rf frontend/dist/
	rm -rf frontend/dist_electron/

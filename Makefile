.PHONY: dev dev-backend dev-frontend build test clean docker-up docker-down docker-build docker-reset help

# ===== 本地开发 =====

dev: dev-deps
	@echo "Starting development servers..."
	@trap 'kill 0' EXIT; \
		cd backend && go run main.go & \
		cd frontend && npm run dev & \
		wait

dev-deps:
	@if [ ! -f backend/go.sum ]; then cd backend && go mod tidy; fi
	@if [ ! -d frontend/node_modules ]; then cd frontend && npm install; fi

dev-backend:
	cd backend && go run main.go

dev-frontend:
	cd frontend && npm run dev

# ===== 构建 =====

build: build-backend build-frontend

build-backend:
	cd backend && CGO_ENABLED=0 go build -ldflags="-s -w" -o server.exe .

build-frontend:
	cd frontend && npm ci && npm run build

# ===== 测试 =====

test: test-backend test-frontend

test-backend:
	cd backend && go vet ./...

test-frontend:
	cd frontend && npx vue-tsc --noEmit

# ===== Docker 开发（仅启动数据库） =====

docker-db:
	docker compose up -d postgres redis

# ===== Docker 全栈 =====

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

docker-build:
	docker compose build

docker-reset:
	docker compose down -v
	docker compose up -d --build

# ===== Docker 本地开发（数据库用 Docker，后端/前端本地跑） =====

docker-local: docker-db
	@echo "Database & Redis ready. Run 'make dev-backend' and 'make dev-frontend' in separate terminals."

# ===== 清理 =====

clean:
	rm -f backend/server.exe
	rm -rf frontend/dist
	rm -rf frontend/node_modules/.cache

# ===== 帮助 =====

help:
	@echo "AI PR Review 项目 Makefile"
	@echo ""
	@echo "本地开发:"
	@echo "  make dev           一键启动后端 + 前端（自动安装依赖）"
	@echo "  make dev-backend   仅启动后端 :8080"
	@echo "  make dev-frontend  仅启动前端 :5173"
	@echo "  make docker-local  后台启动 DB/Redis，前后端本地跑"
	@echo ""
	@echo "构建:"
	@echo "  make build         构建后端 + 前端"
	@echo "  make build-backend 编译后端"
	@echo "  make build-frontend  构建前端"
	@echo ""
	@echo "测试/检查:"
	@echo "  make test          后端 go vet + 前端 vue-tsc"
	@echo ""
	@echo "Docker 全栈:"
	@echo "  make docker-up     启动全部服务（后端+前端+DB+Redis）"
	@echo "  make docker-down   停止全部服务"
	@echo "  make docker-reset  清除数据卷 + 重新构建启动"
	@echo ""
	@echo "清理:"
	@echo "  make clean         删除构建产物"

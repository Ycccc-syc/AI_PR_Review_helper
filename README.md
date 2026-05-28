# <!-- TODO: 项目名称 -->

<!-- TODO: 一句话描述项目 -->

## Demo 视频

<!-- TODO: 替换为实际的B站/云盘链接 -->

> [点击观看 Demo 视频]()（待上传）

---

## 项目简介

<!-- TODO: 描述项目背景、目标用户、解决的痛点 -->

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go 1.22 + Gin |
| 前端 | Vue 3 + TypeScript |
| 数据库 | PostgreSQL 16 |
| 缓存 | Redis 7 |
| 部署 | Docker + Docker Compose |

## 架构设计

```
frontend (Vue 3)  ──HTTP──>  backend (Go + Gin)
                                  │
                    ┌─────────────┼─────────────┐
                    ▼             ▼             ▼
                PostgreSQL     Redis       <!-- TODO: 外部服务/API -->
```

### 后端分层

```
backend/
├── main.go            # 服务入口
├── config/            # 配置管理
├── logger/            # zap 日志
├── model/             # 数据模型 + 统一响应
├── handler/           # HTTP Handler（只做参数校验+调用service）
├── service/           # 业务逻辑
├── repository/        # 数据访问
├── middleware/         # 中间件（日志/恢复/CORS）
└── router/            # 路由注册
```

## 快速开始

### 前置要求

- Go 1.22+
- Docker & Docker Compose（可选）
- PostgreSQL 16（本地运行需要）
- Redis 7（本地运行需要）

### Docker 一键启动

```bash
cp .env.example .env   # 按需修改配置
docker compose up -d
```

### 本地开发

```bash
cd backend
cp ../.env.example .env
# 按需修改 .env 中的数据库和 Redis 地址
go mod tidy
go run main.go
```

服务启动后访问 `http://localhost:8080/api/v1/health`。

## API 文档

### 统一响应格式

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

| code | 含义 |
|------|------|
| 0 | 成功 |
| -1 | 系统异常 |
| <!-- TODO: 业务错误码 --> |

### 接口列表

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/health` | 健康检查 |
| <!-- TODO: 补充业务接口 --> |

## 第三方库说明

| 库 | 用途 | 使用范围 |
|------|------|------|
| github.com/gin-gonic/gin | HTTP 框架 | 全局 |
| go.uber.org/zap | 结构化日志 | 全局 |
| <!-- TODO: 补充其他第三方库 --> |

<!-- TODO: 如有使用旧代码/开源参考，在此声明来源 -->

## 团队分工

<!-- TODO: 如组队，填写成员及分工；如个人完成，填"独立完成" -->

| 成员 | 分工 |
|------|------|
| <!-- TODO: 姓名 | 负责模块 --> |

---

## 开发日志

<!-- TODO: 记录每日进度和关键决策 -->

| 日期 | 进展 | 备注 |
|------|------|------|
| Day 1 | 项目初始化 + 基础骨架 | |
| Day 2 | | |
| Day 3 | Demo 录制 + README 完善 | |

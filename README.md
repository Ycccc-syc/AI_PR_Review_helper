# AI PR Review 助手

基于 DeepSeek 大模型的 GitHub PR 自动代码评审工具。

## Demo 视频

> [点击观看 Demo 视频](https://www.bilibili.com/)（待上传）

---

## 项目简介

AI PR Review 助手帮助开发者快速理解 Pull Request 变更、发现潜在 Bug 和安全风险。用户只需输入 GitHub PR 链接，系统自动获取代码 diff 并通过 DeepSeek 大模型进行智能分析，输出结构化的 Review 报告（变更总结 + 风险识别 + 改进建议）。

### 目标用户

- 团队 Tech Lead / Senior Developer 需要高效 Review 大量 PR
- 开源项目 Maintainer 需要快速了解外部贡献者的代码质量
- 个人开发者需要 AI 辅助检查代码问题

### 核心功能

- 输入 GitHub PR URL，一键获取代码变更并分析
- 结构化报告：变更总结、风险代码标注（高/中/低严重度）、改进建议
- 支持大 PR 自动分片分析（>8000 行按文件拆分）

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
                PostgreSQL     Redis       DeepSeek API / GitHub API
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

### 接口列表

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/v1/health` | 健康检查 |
| POST | `/api/v1/review` | PR 代码审查（body: `{"pr_url": "..."}`） |

## 第三方库说明

| 库 | 用途 | 使用范围 |
|------|------|------|
| github.com/gin-gonic/gin | HTTP 框架 | 全局 |
| go.uber.org/zap | 结构化日志 | 全局 |
| 无其他第三方依赖 | 纯标准库 + DeepSeek/GitHub HTTP API | |

## 声明

本项目为七牛云 X Engineer 暑期实训营选拔作品，所有代码独立完成，未使用旧代码或开源项目参考。

## 团队分工

独立完成

---

## 开发日志

| 日期 | 进展 | 备注 |
|------|------|------|
| Day 1 | 项目初始化 + 基础骨架 | |
| Day 2 | 数据模型、GitHub API、DeepSeek client、Service/Handler | |
| Day 3 | 前端页面 + Demo 录制 + README 完善 | |

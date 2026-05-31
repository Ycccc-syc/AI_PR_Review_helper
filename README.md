# AI PR Review 助手

基于 DeepSeek 大模型的 GitHub PR 自动代码评审工具。输入 PR 链接，自动拉取代码变更，调用 LLM 智能分析，输出结构化报告。支持 SSE 流式输出、交互式追问、PR 健康分、用户系统与审查历史。

## 在线体验

**线上地址：http://8.163.33.73**

## 🎬 Demo 视频

AI PR Review Assistant - Demo Showcase

https://b23.tv/ClKNnMd

---

## 核心功能

- **一键审查**：输入 GitHub PR URL，自动获取 diff + 元数据，DeepSeek LLM 分析
- **结构化报告**：变更总结 + 风险代码标注（高/中/低，含文件路径和行号）+ 改进建议
- **PR 健康分**：安全/性能/可维护三维度 0-100 评分，SVG 环形图展示
- **SSE 流式输出**：分析文本实时推送到前端，非阻塞等待，即时阅读
- **交互式追问**：对任意风险项点击「追问」，独立对话页面携带上下文进行多轮问答
- **用户系统**：注册即登录（bcrypt + JWT 72h），审查结果自动入库，历史记录可回溯
- **智能过滤**：自动跳过 lock 文件、二进制、字体、图片、生成文件
- **大 PR 分片**：>8000 行按文件拆分为 5 文件/组，逐片分析后合并
- **Redis 缓存**：相同 diff 哈希的结果 1 小时内复用，避免重复调用 LLM

## 快速开始

### 方式一：Docker 全栈部署（推荐）

```bash
git clone <repo-url>
cd AI_PR_Review_helper
cp .env.example .env
# 编辑 .env，填入 GITHUB_TOKEN 和 DEEPSEEK_API_KEY，修改 JWT_SECRET
docker compose up -d --build
```

启动后访问 `http://localhost`，前端 nginx 自动反向代理 `/api` 到后端。

### 方式二：本地开发

前置条件：Go 1.25+、Node 22+、Docker Desktop

```bash
# 启动数据库和缓存
make docker-local

# 另开终端启动前后端
make dev
```

前端 `http://localhost:5173`，后端 `http://localhost:8080`。Vite 自动代理 `/api` 到后端。

### Makefile 常用命令

| 命令 | 说明 |
|------|------|
| `make dev` | 一键启动后端 + 前端（自动安装依赖） |
| `make test` | go vet + vue-tsc |
| `make docker-up` | Docker 全栈部署 |
| `make docker-down` | 停止全部服务 |
| `make docker-reset` | 清除数据卷 + 重建 |
| `make build` | 编译后端 + 构建前端 |
| `make help` | 查看所有命令 |

---

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go 1.25 + Gin 1.12 |
| 前端 | Vue 3 + TypeScript + Vite 8 + Vue Router 5 |
| 数据库 | PostgreSQL 16 |
| 缓存 | Redis 7 |
| 认证 | JWT (golang-jwt/jwt/v5) + bcrypt |
| AI | DeepSeek API（兼容 OpenAI 接口） |
| 部署 | Docker + Docker Compose + Nginx |

## 架构

```
frontend (Vue 3 SPA)  ──HTTP──>  backend (Go + Gin)
         │                            │
  Vue Router (4 routes)    ┌──────────┼──────────┬──────────┐
  / /auth /history /ask    ▼          ▼          ▼          ▼
                       PostgreSQL   Redis   DeepSeek API  GitHub API
                       (用户+历史)  (缓存)   (LLM 分析)  (PR 数据)
```

### 后端分层

```
backend/
├── main.go             # 服务入口 + 数据库/Redis 初始化 + 优雅关闭
├── config/             # 环境变量管理（.env 自动加载 + 不安全默认值警告）
├── logger/             # zap 结构化日志
├── model/              # 数据模型（User / ReviewLog / ReviewResponse / API 响应）
├── handler/            # HTTP Handler
│   ├── health.go       # GET /api/v1/health
│   ├── auth.go         # POST /register + /login（JWT 签发，注册即登录）
│   ├── review.go       # POST /review + /review/stream + /review/ask（SSE）
│   └── history.go      # GET /reviews + /reviews/:id
├── service/            # 业务逻辑（过滤/分片/LLM 编排/追问/缓存）
├── repository/         # 数据访问
│   ├── github.go       # GitHub API（PR 解析/元数据/Diff + 限流重试）
│   ├── db.go           # PostgreSQL 连接 + 自动建表迁移
│   ├── user.go         # 用户 CRUD + bcrypt
│   ├── review_log.go   # 审查记录存取
│   └── redis.go        # Redis 缓存（SHA256 键 + TTL）
├── llm/                # DeepSeek API 客户端（普通 Call + 流式 CallStream）
├── prompt/             # System/User Prompt 模板（审查 + 追问 + 健康分）
├── middleware/          # Recovery / Logger / CORS / Auth / OptionalAuth
└── router/             # 路由注册
```

### 中间件链

```
Recovery → Logger → CORS → [OptionalAuth | Auth] → Handler
```

- **OptionalAuth**：审查接口使用，未登录仍可使用，已登录自动入库
- **Auth**：历史记录接口使用，未登录返回 401

---

## API 文档

### 统一响应格式

```json
{ "code": 0, "message": "success", "data": {} }
```

- `code=0` 成功，`code=401` 未登录，`code=-1` 系统异常
- HTTP 状态码始终为 200，错误信息在 `code`/`message` 中区分

### 接口列表

| 方法 | 路径 | 认证 | 说明 |
|------|------|------|------|
| GET | `/api/v1/health` | 无 | 健康检查，返回 `{"status":"ok"}` |
| POST | `/api/v1/auth/register` | 无 | 用户注册，返回 JWT + 用户信息 |
| POST | `/api/v1/auth/login` | 无 | 用户登录，返回 JWT + 用户信息 |
| POST | `/api/v1/review` | 可选 | PR 代码审查 `{"pr_url":"..."}` → 结构化报告 |
| POST | `/api/v1/review/stream` | 可选 | PR 审查 SSE 流式，`text/event-stream` |
| POST | `/api/v1/review/ask` | 无 | 风险项追问 `{"pr_url","question","context"}` SSE |
| GET | `/api/v1/reviews` | 必须 | 审查历史 `?limit=20&offset=0` |
| GET | `/api/v1/reviews/:id` | 必须 | 单条审查详情，含完整结果 JSON |

> `/review` 和 `/review/stream` 使用 OptionalAuth——未登录可用，已登录自动入库。

### SSE 事件格式

```
data: {"type":"chunk","content":"## 代码分析\n"}

data: {"type":"result","data":{...完整 ReviewResponse...}}

data: {"type":"error","message":"错误描述"}
```

---

## 关键设计决策

### 1. 分片策略：文件级切分

大 PR（>8000 行 diff）按文件拆分为 5 文件/组，分别发给 LLM 后合并。

按文件而非按行数切分——文件是最小逻辑完整单元。按行数硬切可能把同一个函数截断在两个 chunk 中。代价是跨文件分析能力丧失（`auth.go` 变更破坏了 `handler.go` 中的调用方，若落在不同 chunk 则 LLM 无法察觉）。

### 2. 文件过滤：注意力预算管理

跳过 lock 文件、二进制、字体、图片、生成文件。`package-lock.json` 动辄上万行，在 Transformer 的 self-attention 中不存在"跳过不读"——每占据一个上下文位置，就挤占一行真实代码变更的注意力配额。

### 3. SSE 流式 > Polling

LLM 逐 token 生成，SSE 单向推送天然适配：HTTP 原生支持，中间件零改动，连接管理开销远低于 WebSocket。

### 4. Temperature = 0.3

- T=0 确定性最高但丧失概率推理能力，可能漏掉非显而易见的 Bug
- T=1.0 高方差，严重度分级不稳定，偶发幻影问题
- T=0.3 平衡点：JSON schema 遵循度 >95%，同时保留足够方差识别非常规问题

### 5. DeepSeek-V3：成本与可达性

| 维度 | DeepSeek-V3 | GPT-4o |
|------|-------------|--------|
| 上下文窗口 | 64K token | 128K token |
| 国内访问 | 无需代理，直连 | 需代理 |
| 单价 (1M tokens) | ¥1 (输入) / ¥4 (输出) | ¥72 / ¥216 |
| JSON 结构遵循度 | 高（T=0.3 时 >95%） | 高 |

### 6. OptionalAuth：登录与未登录无缝兼容

审查是核心价值，不应设登录门槛。历史记录是增值功能，需要用户标识。两条路径通过一个中间件和平共存——未登录正常使用，已登录自动入库。

### 7. Redis 缓存：SHA256 键 + 1h TTL

缓存键基于 `owner/repo/number` + diff 内容的 SHA256 哈希。同一 PR 的不同 commit 生成不同键，新键写入、旧键随 TTL 淘汰。TTL ≤ 0 时默认 24 小时兜底。

---

## 第三方库

| 库 | 用途 |
|------|------|
| github.com/gin-gonic/gin | HTTP 框架 + 中间件链 |
| go.uber.org/zap | 结构化日志 |
| github.com/golang-jwt/jwt/v5 | JWT 签发与验证 |
| golang.org/x/crypto | bcrypt 密码哈希 |
| github.com/lib/pq | PostgreSQL 驱动 |
| github.com/redis/go-redis/v9 | Redis 客户端 |
| vue-router | Vue SPA 路由 |

---

## 项目结构

```
qiniudemo/
├── backend/
│   ├── main.go
│   ├── config/config.go
│   ├── logger/logger.go
│   ├── model/          (user / review / review_log / response)
│   ├── handler/        (health / auth / review / history)
│   ├── service/        (review + review_test)
│   ├── repository/     (db / github / user / review_log / redis)
│   ├── llm/client.go
│   ├── prompt/review.go
│   ├── middleware/      (middleware + auth)
│   ├── router/router.go
│   ├── Dockerfile
│   ├── go.mod / go.sum
├── frontend/
│   ├── src/
│   │   ├── main.ts
│   │   ├── App.vue
│   │   ├── style.css
│   │   ├── router/index.ts
│   │   ├── pages/       (HomePage / AuthPage / HistoryPage / AskPage)
│   │   ├── components/  (ReviewInput / ReviewResult / ReviewHistory / AuthForm)
│   │   ├── composables/ (useAuth / useReview)
│   │   └── api/         (review / auth / history)
│   ├── index.html
│   ├── Dockerfile
│   ├── nginx.conf
│   ├── vite.config.ts
│   ├── package.json
├── docker-compose.yml
├── Makefile
├── .env.example
├── CLAUDE.md
└── README.md
```

---

## 声明

本项目为七牛云 X Engineer 暑期实训营选拔作品，独立完成。使用的第三方库均通过 Go modules / npm 管理，未修改源码。

---

## 开发日志

| 日期 | 进展 |
|------|------|
| Day 1 | 项目初始化 → Go 后端全链路（GitHub API + DeepSeek Client + Service 编排）→ Vue 前端 → 端到端跑通 |
| Day 2 | SSE 流式审查 + 用户系统（注册/登录/JWT/历史）+ 前端路由重构（4 路由）+ 交互式追问 + PR 健康分环形图 |
| Day 3 | 边界处理 + Redis 缓存 + Docker Compose + Makefile + nginx 部署 + 单元测试 + AI 审查自身代码修复 13 项 |

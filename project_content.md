# project_content.md — 项目文件清单与当前状态

> 自动生成，记录每个文件的用途和状态。

## 根目录

| 文件 | 用途 | 状态 |
|------|------|------|
| `.env.example` | 环境变量模板，含 GITHUB_TOKEN / DEEPSEEK_API_KEY 等 | ✅ |
| `.env` | 实际环境变量（gitignore，存放真实 key） | ✅ |
| `CLAUDE.md` | Claude Code 项目指引（技术栈/架构/规范/进度） | ✅ |
| `TODO.md` | 开发任务清单（Day1-3） | ✅ |
| `README.md` | 项目说明（面向评审） | ✅ |
| `docker-compose.yml` | PostgreSQL + Redis + 后端 + 前端 编排 | ✅ |
| `Dockerfile` | 后端构建镜像 | ✅ |
| `frontend/Dockerfile` | 前端构建镜像 | ✅ |

---

## backend/ — Go 服务

### 入口

| 文件 | 用途 | 状态 |
|------|------|------|
| `main.go` | 服务入口：`config.Load()` → `logger.Init()` → `router.Setup(cfg)` → 启动 HTTP | ✅ |

### config/

| 文件 | 用途 | 状态 |
|------|------|------|
| `config/config.go` | 环境变量加载 + `loadDotEnv()` 自动从当前/父目录加载 .env | ✅ |

### logger/

| 文件 | 用途 | 状态 |
|------|------|------|
| `logger/logger.go` | zap 全局日志初始化 | ✅ |

### middleware/

| 文件 | 用途 | 状态 |
|------|------|------|
| `middleware/middleware.go` | Recovery + Logger + CORS 中间件 | ✅ |

### model/

| 文件 | 用途 | 状态 |
|------|------|------|
| `model/review.go` | ReviewRequest / ReviewResponse / RiskItem / Suggestion / ReviewStats | ✅ |
| `model/response.go` | 统一响应：Success() / Error() | ✅ |

### repository/

| 文件 | 用途 | 状态 |
|------|------|------|
| `repository/github.go` | ParsePRURL / FetchPRMeta / FetchPRDiff / githubGet（含 403 重试） | ✅ |

### llm/

| 文件 | 用途 | 状态 |
|------|------|------|
| `llm/client.go` | DeepSeek API 调用（OpenAI 兼容），temp=0.3，超时 120s，LimitReader 10MB | ✅ |

### prompt/

| 文件 | 用途 | 状态 |
|------|------|------|
| `prompt/review.go` | 中文 system prompt + BuildMessages() 拼接 user prompt | ✅ |

### handler/

| 文件 | 用途 | 状态 |
|------|------|------|
| `handler/health.go` | GET /api/v1/health 健康检查 | ✅ |
| `handler/review.go` | POST /api/v1/review — 参数校验 + 调用 service.AnalyzePR | ✅ |

### service/

| 文件 | 用途 | 状态 |
|------|------|------|
| `service/review.go` | AnalyzePR → filterDiff → buildMessages → LLM Call → extractJSON → ReviewResponse | ✅ |

### router/

| 文件 | 用途 | 状态 |
|------|------|------|
| `router/router.go` | `Setup(cfg *config.Config)` — 注册 /health + /review 路由 | ✅ |

---

## frontend/ — Vue 3 前端

### 入口 + 配置

| 文件 | 用途 | 状态 |
|------|------|------|
| `frontend/src/main.ts` | Vue 应用入口 | ✅ |
| `frontend/src/App.vue` | 根组件 — 整合 ReviewInput + ReviewResult + error 状态 + tab 切换 | ✅ |
| `frontend/vite.config.ts` | Vite 配置，含 `/api` 代理到 `localhost:8080` | ✅ |

### 组件

| 文件 | 用途 | 状态 |
|------|------|------|
| `frontend/src/components/ReviewInput.vue` | PR URL 输入区 + loading spinner + emit submit | ✅ |
| `frontend/src/components/ReviewResult.vue` | 结果展示（统计卡片 + 总结/风险/建议 tab） | ✅ |
| `frontend/src/components/HelloWorld.vue` | Vite 脚手架默认组件 | 🔲 可删除 |

### API 层

| 文件 | 用途 | 状态 |
|------|------|------|
| `frontend/src/api/review.ts` | `analyzePR(prUrl)` — AbortController 120s 超时 + Response 校验 | ✅ |

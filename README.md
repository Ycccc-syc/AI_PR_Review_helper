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

- **一键审查**：输入 GitHub PR URL，自动拉取 diff + 元数据，调用 DeepSeek LLM 分析
- **结构化报告**：变更总结 + 风险代码标注（高/中/低严重度，含文件路径和行号）+ 改进建议
- **SSE 流式输出**：分析过程实时推送到前端，非阻塞等待，用户可即时阅读分析内容
- **交互式追问**：对任意风险项深入追问，携带原始代码上下文 + 对话历史进行多轮分析
- **大 PR 分片**：>8000 行按文件自动拆分，逐片分析后合并结果
- **用户系统**：注册/登录 + JWT 认证，Review 结果自动入库，支持历史记录回溯
- **智能过滤**：自动跳过 lock 文件、二进制、字体、图片、生成文件，聚焦有效变更

## 技术栈

| 层级 | 技术 |
|------|------|
| 后端 | Go 1.22 + Gin |
| 前端 | Vue 3 + TypeScript + Vite |
| 数据库 | PostgreSQL 16 |
| 缓存 | Redis 7 |
| 认证 | JWT (golang-jwt/jwt/v5) + bcrypt |
| AI | DeepSeek API (V3，兼容 OpenAI 接口) |
| 部署 | Docker + Docker Compose |

## 架构设计

```
frontend (Vue 3)  ──HTTP──>  backend (Go + Gin)
         │                         │
         │               ┌─────────┼─────────┬──────────┐
         ▼               ▼         ▼         ▼          ▼
      Vue Router    PostgreSQL   Redis   DeepSeek API  GitHub API
     (SPA 路由)     (用户+历史)  (缓存)   (LLM 分析)  (PR 数据)
```

### 后端分层

```
backend/
├── main.go             # 服务入口 + 数据库初始化 + 优雅关闭
├── config/             # 环境变量管理（支持 .env 自动加载）
├── logger/             # zap 结构化日志
├── model/              # 数据模型（User / ReviewLog / API 统一响应）
├── handler/            # HTTP Handler（参数校验 + 调用 service）
│   ├── health.go       # 健康检查
│   ├── auth.go         # 用户注册/登录 + JWT 签发
│   ├── review.go       # 审查接口（普通/SSE 流式/追问）
│   └── history.go      # 历史记录查询
├── service/            # 业务逻辑（过滤/分片/LLM 编排/追问上下文）
├── repository/         # 数据访问
│   ├── github.go       # GitHub API 封装（PR 解析/元数据/Diff 获取）
│   ├── db.go           # PostgreSQL 连接 + 自动建表迁移
│   ├── user.go         # 用户 CRUD + bcrypt 密码验证
│   └── review_log.go   # 审查记录存取
├── llm/                # DeepSeek API 客户端（普通 + SSE 流式）
├── prompt/             # System/User Prompt 模板
├── middleware/          # Recovery / Logger / CORS / JWT Auth / OptionalAuth
└── router/             # 路由注册 + 中间件链
```

---

## 关键设计决策

### 1. 上下文窗口：不是内存问题，是 Transformer 架构约束

DeepSeek-V3 拥有 64K token 上下文窗口。入门者常将其理解为"内存够不够"，但真正的瓶颈不在显存——

**核心约束来自 Rotary Position Embedding (RoPE)。** Transformer 通过位置嵌入编码 token 顺序，RoPE 对每个 token 施加旋转角度与位置成正比的旋转矩阵。当序列从 8K 扩展到 64K 时，三个不可逾越的物理限制出现：

1. **注意力稀释（Attention Dilution）**：多头自注意力机制将有限的注意力权重分配在所有 token 对上。序列膨胀 8 倍 → 每个 token 获得的平均注意权重下降 8 倍。模型必须在 64K 个位置中自主判断"哪里该关注"，而远端 token 天然在与近邻 token 的竞争中落败。

2. **位置信号退化（Position Signal Decay）**：RoPE 的旋转角度与位置索引成正比。在极端远端位置，相邻 token 之间的**角度差**趋近于零——模型越来越难分辨"第 142 行"和"第 143 行"。对代码审查而言，这意味着**行级精确定位能力随序列长度衰减**。

3. **长距依赖稀释（Long-range Dependency Dilution）**：第 500 行的一个 Bug 若依赖于第 50 行定义的函数，需要模型维持 450 token 的跨距依赖。在 64K 的全注意力分布下，这一信号的数学期望已接近噪声水平。

**我们的 8000 行分片阈值正是基于此推导**：diff 平均每行 ≈30 token，8000 行 ≈240K 输入 token。为 system prompt 和 PR 元数据预留 ~8K token，为 LLM 输出预留 ~8K 余量，我们将单次请求控制在模型**有效注意力区间**内——即行级精确定位能力尚未显著衰减的区间。

### 2. 分片策略：文件级切分的权衡

大 PR（>8000 行）按文件拆分为多个 chunk，每 chunk 5 个文件，分别发给 LLM 分析后合并。

**为什么按文件而不是按行数？**
- 文件是最小的逻辑完整单元。按行数硬切可能把同一个函数截断在两个 chunk 中，导致两边都不完整
- 代价是跨文件分析能力丧失——若 `auth.go` 的变更破坏了 `handler.go` 中的调用方，且它们落在不同 chunk，LLM 无法察觉

**为什么不分发到更多上下文？**
- 引入完整文件内容会增加 5-10 倍 token 消耗，须在准确性提升与成本/延迟之间做取舍
- 实验性方案：仅对 LLM 标记为 high 风险的文件追加完整内容做二次分析（ReAct 一环）

### 3. 文件过滤：注意力预算管理

过滤 lock 文件、二进制、字体、生成文件，表面上为了减少 token 消耗，本质上是**注意力预算的主动分配**。

`package-lock.json` 动辄上万行。在 Transformer 中，self-attention 对**每一对 tokens** 计算注意力分数——不存在"跳过不读"的概念。这份 JSON 每占据一个上下文位置，就挤占了一行真实代码变更的注意力配额。

**过滤逻辑在 service 层而非 repository 层**（当前实现在 service/review.go 中作为预处理），因为它属于 Review 业务逻辑范畴——不同 Agent 可能需要不同的过滤规则。

### 4. 模型选择：DeepSeek-V3 而非 GPT-4/Claude

| 维度 | DeepSeek-V3 | GPT-4o | Claude 4 |
|------|-------------|--------|----------|
| 代码理解 | HumanEval 接近 GPT-4 水平 | 业界标杆 | 业界标杆 |
| 上下文窗口 | 64K | 128K | 200K |
| 国内访问 | 无需代理，直连 | 需代理 | 需代理 |
| 单价 (1M tokens) | ¥1 (输入) / ¥4 (输出) | ¥72 / ¥216 | ¥108 / ¥432 |
| Temperature 稳定性 | JSON 结构遵循度高 | 偶尔自由发挥 | JSON 结构遵循度高 |

选择 DeepSeek 的核心考量：国内直连的稳定性 + 20-50 倍的成本优势 + 64K 上下文对 PR Review 场景够用（我们限制在 ~240K 输入 token 以下）。

### 5. Temperature = 0.3：结构化输出的甜点

Temperature 控制输出概率分布的熵：

- **T=0**：确定性输出。幻觉率最低，但完全丧失概率推理能力——可能漏掉需要"推理跳跃"的非显而易见 Bug
- **T=1.0**：高方差。能发现隐蔽问题，但严重度分级不稳定，偶发幻影问题（phantom issue）
- **T=0.3**：平衡点。经验上保证 JSON schema 遵循度（>95%），同时保留足够方差去识别非常规代码问题

### 6. SSE 流式 > Polling

技术方案选择了 SSE 而非 WebSocket 或轮询：

| 方案 | 复杂度 | 场景适配 |
|------|--------|---------|
| SSE | 单向推送，HTTP 原生支持，中间件零改动 | LLM token 逐字生成，天然适配 |
| WebSocket | 双向通信，需额外握手和连接管理 | 多轮对话需双向的场景 |
| Polling | 最简单但延迟高，无效请求多 | 不适合实时性要求高的场景 |

对于"发 diff → LLM 逐 token 生成 → 推送给前端"的单向流式场景，SSE 是最小化复杂度的方案。

### 7. OptionalAuth：登录/未登录无缝兼容

Review 接口采用 `OptionalAuth` 中间件而非强制 `Auth`：未登录用户可正常使用审查功能，已登录用户审查结果自动入库到个人历史。

设计考量：审查是核心价值，不应设置登录门槛。历史记录是增值功能，需要用户标识。两条路径通过一个中间件和平共存，避免维护两个版本的 Review handler。

---

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

| 方法 | 路径 | 认证 | 说明 |
|------|------|------|------|
| GET | `/api/v1/health` | 无 | 健康检查 |
| POST | `/api/v1/auth/register` | 无 | 用户注册 |
| POST | `/api/v1/auth/login` | 无 | 用户登录，返回 JWT |
| POST | `/api/v1/review` | 可选 | PR 代码审查（`{"pr_url": "..."}`） |
| POST | `/api/v1/review/stream` | 可选 | PR 代码审查（SSE 流式，实时推送分析文本） |
| POST | `/api/v1/review/ask` | 无 | 追问风险项（`{"pr_url","question","context"}`，SSE） |
| GET | `/api/v1/reviews` | 必须 | 当前用户历史记录（`?limit=&offset=`） |

> `/review` 和 `/review/stream` 为可选认证：未登录可使用，已登录则自动入库到历史记录。

## 第三方库说明

| 库 | 用途 |
|------|------|
| github.com/gin-gonic/gin | HTTP 框架 + 中间件链 |
| go.uber.org/zap | 结构化日志 |
| github.com/golang-jwt/jwt/v5 | JWT 签发与验证 |
| golang.org/x/crypto | bcrypt 密码哈希 |
| github.com/lib/pq | PostgreSQL 驱动 |

## 声明

本项目为七牛云 X Engineer 暑期实训营选拔作品，独立完成。使用的第三方库均通过 Go modules / npm 管理，未修改源码。

---

## 开发日志

| 日期 | 进展 |
|------|------|
| Day 1 | 项目初始化 → Go 后端全链路（GitHub API + DeepSeek Client + Service 编排）→ Vue 前端 → 端到端跑通 |
| Day 2 | SSE 流式审查 + 交互式追问 + 用户系统（注册/登录/JWT） + 历史记录 + 前端路由重构 |
| Day 3 | 边界处理完善 + Demo 录制 + README 设计决策文档 |

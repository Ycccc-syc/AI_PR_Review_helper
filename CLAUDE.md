# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目背景

本仓库用于七牛云 X Engineer 暑期实训营选拔。规则：
- 72 小时限时开发，个人或组队
- **必须持续 commit/PR**，最后一天一次性提交视为无效
- 提交物：GitHub 仓库 + README + Demo 视频（B站/云盘）
- 评审维度：完整度与创新(40%) + 开发过程与质量(40%) + 演示与表达(20%)

## 常用命令

```bash
# 后端
cd backend
go mod tidy                    # 安装依赖
go run main.go                 # 启动服务（端口 8080）
go build -o server .           # 编译

# 前端
cd frontend
npm install                    # 安装依赖
npm run dev                    # 启动开发服务器
npm run build                  # 类型检查 + 构建

# Docker
cp .env.example .env
docker compose up -d           # 启动全部服务
docker compose down            # 停止
```

## 技术栈

| 层 | 技术 |
|----|------|
| 后端 | Go 1.25 + Gin + Zap |
| 前端 | Vue 3 + TypeScript + Vite 8 |
| 数据库 | PostgreSQL 16 |
| 缓存 | Redis 7 |
| AI | DeepSeek API（V3 或 R1） |

## 架构约定

### 后端分层（必须遵守）

```
handler → service → repository → model
```

- **handler**：只做参数校验和调用 service，不含业务逻辑
- **service**：业务逻辑层
- **repository**：数据访问层（DB/Redis）
- **model**：数据结构定义 + 统一响应

### API 统一响应

```json
{"code": 0, "message": "success", "data": {}}
```

- code=0 成功，code=-1 系统异常，其他按业务定义
- 即使出错 HTTP status 也为 200，错误信息在 code/message 中区分

### 中间件链

```
Recovery → Logger → CORS → Handler
```

### 配置

全部通过环境变量注入，见 `.env.example`。`config/config.go` 中统一加载。

### Git 提交规范（实训营强制要求）

- 每个 PR 只做一件事，标题清晰，描述含功能说明+实现思路+测试方式
- 大功能拆成多个小 PR
- 每日常态 commit，不能集中到最后一天

## Skills 工作流

本仓库配置了自定义 skills（`skills/` 目录），建议按以下流程使用：

```
/planner              → 拿到议题后先出技术方案和分步计划
/fullstack-builder    → 按计划逐步实现
/review 或 /simplify  → 实现完后 review 代码质量
/debug                → 出 Bug 时分析根因并修复
```

---

## 通用开发原则

1. 优先可维护性
2. 优先小步修改
3. 不要大规模重构
4. 不要过度设计
5. 优先复用已有代码
6. 保持代码风格统一
7. 不要创建重复文件
8. 不要随意修改目录结构
9. 生产环境优先
10. 代码必须易读

## 工作流程

每次开发时：

1. 先理解需求
2. 阅读相关代码
3. 分析问题
4. 给技术方案
5. 再实现代码

不要直接输出大量代码。

## 输出规范

输出时：

1. 先解释原因
2. 再给方案
3. 再逐步实现

不要一次输出完整项目。

## Debug 规范

出现问题时：

1. 先分析根因
2. 给排查思路
3. 再修改代码

不要盲目猜测。

## 数据库规范

1. 优先考虑索引
2. 避免 N+1 查询
3. 合理使用事务
4. 不要重复查询

## Redis 规范

1. key 命名清晰
2. 设置合理 TTL
3. 避免大 key

## 代码修改规则

修改代码时：

1. 尽量最小修改
2. 不要覆盖整个文件
3. 保留原有架构
4. 保持兼容性

## 回复风格

1. 简洁专业
2. 避免废话
3. 优先工程实践
4. 优先真实项目方案
5. 给出具体步骤

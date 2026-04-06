# Research: Go 编程语言学习样例结构

## Decision 1: 项目结构 — 单一 go.mod vs 多 module

**Decision**: 单一 `go.mod`，所有章节作为 `internal/` 子包

**Rationale**: 
- 学习项目的章节依赖基本一致（fmt, time, net/http, database drivers）
- 多 module 增加维护成本（每个目录独立 `go.mod`/`go.sum`）
- 用户只需 `go run ./cmd/hello variables` 即可运行，无需切换目录
- Go 官方项目（Kubernetes, Prometheus）均采用单一 module 模式

**Alternatives considered**:
- 每个章节独立 `go.mod` — 依赖隔离但维护成本过高
- `go.work` workspace — 适合多 module 开发，但学习项目不需要

---

## Decision 2: 示例代码组织 — internal/ vs examples/

**Decision**: `internal/basic/`, `internal/advance/`, `internal/awesome/`

**Rationale**:
- 每个章节为独立子包（如 `internal/basic/variables/`），暴露 `func Run()` 
- `cmd/hello/main.go` 统一入口，通过子命令路由
- 代码不会被外部 import（`internal/` 限制），符合学习样例的定位
- 运行方式统一：`go run ./cmd/hello basic variables`

**Alternatives considered**:
- `examples/` 目录 — Go 惯例放示例，但学习项目需要统一入口调度
- 每个章节独立二进制 — 构建产物过多，用户需要记住多个命令

---

## Decision 3: CLI 入口设计

**Decision**: Cobra 子命令结构

```
hello-go basic variables     # 基础章节
hello-go advance database    # 高级章节
hello-go awesome webservice  # 实战项目
hello-go algo sort           # 算法
hello-go quiz                # 知识检查
```

**Rationale**:
- 项目已有 Cobra 框架（`cmd/foo/`）
- 子命令提供清晰的层级导航
- `--help` 自动生成文档，符合 Constitution Principle III

**Alternatives considered**:
- 纯 `os.Args` 解析 — 简单但缺乏帮助文档
- 交互式菜单 — 用户体验好但不适合自动化测试

---

## Decision 4: 文档结构

**Decision**: mdBook，完全镜像 hello-rust 的目录结构

```
docs/src/
├── basic/          # 基础章节文档
├── advance/        # 高级章节文档
├── awesome/        # 实战项目文档
├── algo/           # 算法
├── leetcode/       # LeetCode 题解
├── quick_reference/# 速查表
├── quiz/           # 题库
├── projects/       # 项目实战
├── glossary.md     # 术语表
├── faq.md          # FAQ
└── SUMMARY.md      # 目录
```

**Rationale**:
- 与 hello-rust 保持一致的学习体验
- mdBook 已配置（book.toml, GitHub Pages 部署）
- 中文编写，技术术语保留英文

---

## Decision 5: 章节代码模板

**Decision**: 每个章节子包统一模板

```go
// internal/basic/variables/main.go
package variables

import "fmt"

// Run 演示变量与表达式的使用
func Run() {
    fmt.Println("=== 变量与表达式 (Variables & Expressions) ===")
    // 示例代码...
}
```

**Rationale**:
- 统一的 `func Run()` 签名，便于 `cmd/hello` 调度
- 包名即章节名，import 路径清晰
- 文档注释使用中文 + 英文术语

---

## Decision 6: 数据库示例处理

**Decision**: SQLite 使用内存模式或临时文件，无需外部服务

**Rationale**:
- go-sqlite3 支持 `:memory:` 模式，零配置
- 项目已有 SQLite 示例（`internal/repository/sqlite/`）
- 无需 Docker 或外部数据库服务

**Alternatives considered**:
- PostgreSQL/MySQL — 需要外部服务，增加复杂度
- Docker compose — 违背"纯 go run 启动"原则

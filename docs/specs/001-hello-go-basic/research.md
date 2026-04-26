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

## Decision 7: Overview 文档结构模板

**Decision**: Basic/Advance 统一模板，Awesome 独立结构

**Rationale**:
- Basic 和 Advance 都是学习章节导览页，内容类型一致（学习目标 + 章节导航 + 学习路径 + 下一步），统一模板保持一致性。
- Awesome 是项目实战导览页，内容类型不同（项目介绍 + 技术栈 + 应用场景），独立结构更合理。

**模板结构 (Basic/Advance)**:
```markdown
# [标题]
  
## 概述 (1-2 段)
## 学习目标清单 (具体可验证的能力描述)
## 章节导航 (每章 1-2 句摘要 + 🔵🟡🔴 难度标记)
## 学习路径建议 (总时长 + 学习策略)
## 下一步 (导航到下一模块)
```

**模板结构 (Awesome)**:
```markdown
# [标题]

## 概述 (1-2 段)
## 实战项目导航 (4 项目 × 名称+技术栈+适合人群/能力点+1-2 句摘要)
## 应用场景建议
## 前置要求
```

**Alternatives considered**:
- 统一三个页面结构 — 不适合 Awesome 的项目特性

---

## Decision 8: Advance 前置知识自检清单实现

**Decision**: 3-5 道自检题目，使用 `<details>/<summary>` HTML 折叠

**Rationale**: 类似 `advance/context.md` 的 "本章适合谁" 模式，可折叠节省页面空间。

---

## No NEEDS CLARIFICATION items remain.

All unknowns were resolved during the `/speckit.clarify` session (9 questions answered).

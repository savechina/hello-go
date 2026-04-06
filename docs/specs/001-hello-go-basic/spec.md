# Feature Specification: Go 编程语言学习样例结构

**Feature Branch**: `001-hello-go-basic`  
**Created**: 2026-04-05  
**Status**: Draft  
**Input**: User description: basic sample and document ,参考../hello-rust 学习样例结构 basic advance awesome 等分不同学习程度构建go 编程语言学习内容

## Clarifications

### Session 2026-04-05

- Q: 内容范围 — 是否包含算法、速查表、题库等附加模块？ → A: 完全镜像 hello-rust：algo + leetcode + quick_reference + quiz + projects + 附录
- Q: 示例代码组织方式 — 如何存放各章节的可运行代码？ → A: 单一 `go.mod`，`basic/`, `advance/`, `awesome/` 各章节为子包，由 `cmd/hello` 统一入口调用
- Q: 项目实战的交付标准？ → A: 每个项目提供 README + 运行说明，纯 `go run` 启动，无 Docker 支持
- Q: 章节文件命名规范？ → A: 遵循 Go 标准库惯例，`concurrency/concurrency.go`（文件名=包名），不用 `main.go`
- Q: 补充缺失的高阶章节？ → A: 优先补充 Context 和高级并发（Mutex/RWMutex/atomic）章节到 advance 级别

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 基础入门学习路径 (Priority: P1)

作为 Go 初学者，我希望能够通过结构化的基础示例逐步掌握 Go 的核心语法和基本概念，每个知识点都有可运行的代码示例和对应的中文文档说明。

**Why this priority**: 这是项目的核心价值——为 Go 初学者提供清晰、可运行的学习材料。没有基础入门内容，项目就失去了存在的意义。

**Independent Test**: 初学者可以按照文档顺序阅读每个章节，并运行对应的代码示例，验证输出与文档描述一致。

**Acceptance Scenarios**:

1. **Given** 用户打开项目文档，**When** 用户按照基础入门章节顺序学习，**Then** 用户能够找到每个知识点对应的可运行代码示例
2. **Given** 用户运行某个基础示例代码，**When** 代码执行完成，**Then** 输出结果与文档描述一致且无编译错误
3. **Given** 用户完成所有基础章节学习，**Then** 用户已掌握变量声明、数据类型、流程控制、函数、结构体、接口、并发等 Go 核心概念

---

### User Story 2 - 高级进阶学习路径 (Priority: P2)

作为已掌握 Go 基础的开发者，我希望能够深入学习 Go 的高级特性（如错误处理模式、泛型、反射、性能优化、数据库操作、Web 开发等），并通过生产级示例理解最佳实践。

**Why this priority**: 基础入门吸引用户，高级进阶留住用户。提供从入门到精通的完整路径是项目区别于简单示例集合的关键。

**Independent Test**: 有 Go 基础的开发者可以直接跳转到高级章节学习，无需依赖基础章节内容即可理解和使用。

**Acceptance Scenarios**:

1. **Given** 用户已具备 Go 基础知识，**When** 用户进入高级章节学习，**Then** 每个高级主题都有独立的、可运行的示例代码
2. **Given** 用户学习数据库或 Web 开发章节，**When** 用户运行示例，**Then** 示例能正确连接数据库或启动 HTTP 服务并返回预期结果

---

### User Story 3 - 精选实战项目 (Priority: P3)

作为完成学习的开发者，我希望能够参考生产级的实战项目（如 CLI 工具、微服务框架、消息处理等），了解如何将 Go 知识点组合成完整的应用。

**Why this priority**: 实战项目展示知识综合运用能力，帮助学习者从"懂语法"过渡到"能干活"。

**Independent Test**: 用户可以独立运行每个实战项目的示例，理解其架构设计和代码组织方式。

**Acceptance Scenarios**:

1. **Given** 用户完成基础和高级学习，**When** 用户查看实战项目，**Then** 项目结构清晰、代码符合 Go 最佳实践
2. **Given** 用户运行实战项目示例，**When** 项目启动，**Then** 能正常提供服务或完成预期功能

---

### Edge Cases

- 当某个示例需要外部服务（如数据库）时，文档是否提供了启动说明或 mock 方案？
- 当 Go 版本不匹配时，是否有明确的版本要求和错误提示？
- 项目实战 (projects) 中的每个项目 MUST 提供独立的 `README.md`，包含运行说明和预期输出，用户通过 `go run` 启动，无需 Docker 或额外部署工具

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 项目 MUST 按照多级结构组织内容：基础入门 (basic)、高级进阶 (advance)、精选实战 (awesome)、算法与练习 (algo)、LeetCode 题解 (leetcode)、代码片段速查 (quick_reference)、知识检查题库 (quiz)、项目实战 (projects)、附录 (glossary/faq/changelog)
- **FR-002**: 每个难度等级 MUST 包含多个独立章节，每个章节对应一个或多个可运行的 Go 代码示例
- **FR-003**: 所有文档 MUST 使用中文编写，技术术语保留英文原文（括号标注）
- **FR-004**: 每个代码示例 MUST 可独立编译运行，无编译错误和未处理的 panic
- **FR-005**: 项目 MUST 使用单一 `go.mod`；`basic/`, `advance/`, `awesome/` 各章节作为 `internal/` 下的子包，由 `cmd/hello` 统一入口通过子命令调用（如 `go run ./cmd/hello basic variables`）
- **FR-006**: 文档 MUST 通过 mdBook 构建为可浏览的静态网站，包含目录导航和章节跳转
- **FR-007**: 基础入门章节 MUST 覆盖 Go 核心语法：变量与数据类型、流程控制、函数、结构体、接口、错误处理、并发（goroutine/channel）、包管理
- **FR-008**: 高级进阶章节 MUST 覆盖：Context 上下文、高级并发（Mutex/RWMutex/atomic/cond）、泛型、反射、性能调优、数据库操作、Web 开发、测试、配置管理
- **FR-009**: 精选实战 MUST 包含至少 3 个完整项目示例（如 CLI 工具、Web 服务、数据处理管道）
- **FR-010**: 每个章节 MUST 包含：概念说明、代码示例（至少 3 个）、知识点总结、练习题或思考题
- **FR-011**: 项目 MUST 包含算法实现章节（algo），提供常见算法的 Go 实现（如链表、排序、搜索）
- **FR-012**: 项目 MUST 包含 LeetCode 题解章节（leetcode），至少覆盖 5 道经典题目
- **FR-013**: 项目 MUST 提供代码片段速查（quick_reference），覆盖常用语法和标准库用法
- **FR-014**: 项目 MUST 包含知识检查题库（quiz），每章配套至少 3 道选择题或判断题
- **FR-015**: 项目 MUST 包含至少 4 个独立项目实战（projects）：命令行工具、HTTP 服务器、并发爬虫、IPC/分布式示例
- **FR-016**: 项目 MUST 包含附录：术语表 (glossary)、常见问题 FAQ、更新日志 (CHANGELOG)

### Key Entities

- **学习章节 (Chapter)**: 每个章节代表一个独立的学习主题，包含文档和对应的代码示例。章节属于某个难度等级。
- **难度等级 (Level)**: basic / advance / awesome 三个等级，代表内容的难度递进关系。
- **代码示例 (Example)**: 可独立运行的 Go 代码，存放在 `internal/basic/`, `internal/advance/`, `internal/awesome/` 下，每个章节为独立子包，由 `cmd/hello` 统一入口调用。项目使用单一 `go.mod`。

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 基础入门章节至少包含 12 个独立主题章节，每章至少 3 个可运行代码示例
- **SC-002**: 高级进阶章节至少包含 8 个独立主题章节
- **SC-003**: 精选实战至少包含 3 个完整项目示例
- **SC-004**: 所有代码示例通过构建命令一次性编译成功率 100%
- **SC-005**: mdBook 文档构建无链接错误（零 404），构建时间 < 5 分钟
- **SC-006**: 初学者按照文档顺序学习，能够在 2 小时内完成基础入门所有章节并运行所有示例

## Assumptions

- 目标用户为中文母语的 Go 学习者，具备基本编程经验（了解至少一门编程语言）
- 示例代码运行环境为 macOS 或 Linux，Go 1.24+
- 部分示例需要 SQLite 数据库（通过 go-sqlite3，CGO 支持）
- 项目结构参考 hello-rust 的 basic/advance/awesome 三级分层模式
- 文档部署通过 GitHub Pages（与 hello-rust 一致）
- 不需要用户认证、登录等 Web 功能——这是纯学习样例项目

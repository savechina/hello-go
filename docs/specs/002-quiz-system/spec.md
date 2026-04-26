# Feature Specification: 知识检查题库系统 (Quiz System)

**Feature Branch**: `002-quiz-system`  
**Created**: 2026-04-26  
**Status**: Draft  
**Input**: FR-014 from spec 001-hello-go-basic — "项目 MUST 包含知识检查题库（quiz），每章配套至少 3 道选择题或判断题"

## Clarifications

### Session 2026-04-26

- Q: 题库格式？ → A: YAML 格式，每章一个独立文件，便于作者维护
- Q: 答题方式？ → A: CLI 交互模式，`go run ./cmd/hello quiz <level> [chapter]`
- Q: 每章题目数量？ → A: ≥3 道，支持混合选择题 + 判断题
- Q: 是否计分/反馈？ → A: 是，答题后立即显示对错和解析，结束时显示总分
- Q: 题库覆盖范围？ → A: Phase 1 基础章节（12章），Phase 2 高级章节（9章），Phase 3 实战项目（4个）

## User Scenarios & Testing *(mandatory)*

### User Story 1 - 章节练习模式 (Priority: P1)

作为一名正在学习 Go 基础章节的初学者，我希望在完成每章学习后，通过 3-5 道选择题和判断题验证自己的理解程度，答错时能看到解析以便查漏补缺。

**Why this priority**: 这是题库系统的核心价值——让学习者能自我评估，而非被动阅读。

**Independent Test**: 用户运行 `go run ./cmd/hello quiz basic variables`，系统依次出 3 道题，用户选择答案后可立即看到对错和解析，最后显示得分。

**Acceptance Scenarios**:

1. **Given** 用户完成了变量章节的学习，**When** 用户运行 `go run ./cmd/hello quiz basic variables`，**Then** 系统逐题展示，用户输入 A/B/C/D 作答
2. **Given** 用户答错题，**When** 提交答案后，**Then** 系统显示"错误"标记 + 解析说明
3. **Given** 用户完成所有题目，**When** 答题结束，**Then** 系统显示总分（如 2/3）和学习建议

---

### User Story 2 - 全章汇总复习模式 (Priority: P2)

作为一名准备进入下一阶段的学习者，我希望能够一次性完成整个基础部分或高级部分的知识检查，了解自己哪些章节掌握薄弱。

**Why this priority**: 阶段复习是学习路径中承上启下的关键环节，需要汇总评估功能。

**Independent Test**: 用户运行 `go run ./cmd/hello quiz basic`，系统从所有基础章节中随机抽取题目或按顺序出题，完成后按章节分类显示薄弱环节。

**Acceptance Scenarios**:

1. **Given** 用户已完成基础部分全部 12 章学习，**When** 用户运行 `go run ./cmd/hello quiz basic`，**Then** 系统按章节顺序出题或随机出 20+ 道题
2. **Given** 用户完成汇总测试，**When** 查看结果，**Then** 系统按章节显示正确率（如 "变量: 3/3, 接口: 1/3 ⚠️"）

---

### User Story 3 - 题库索引页面 (Priority: P3)

作为一名浏览文档的用户，我希望在 mdBook 的题库首页看到所有章节对应的题库链接和题目数量，方便按需练习。

**Why this priority**: 题库页面是用户进入 quiz 系统的入口，需要清晰的目录结构。

**Independent Test**: 用户打开 `docs/src/quiz/index.md`，能看到按基础/高级/实战分类的题库列表，每章标注题目数量。

**Acceptance Scenarios**:

1. **Given** 用户打开题库首页，**When** 浏览章节列表，**Then** 用户能看到每章的题目数量和类型（选择题/判断题）
2. **Given** 用户点击某章链接，**When** 进入该章题库页面，**Then** 用户能看到所有题目、答案和解析（供离线阅读）

---

### Edge Cases

- 当用户输入无效的 quiz 命令参数时，系统显示帮助信息
- 当某章题库文件不存在时，系统提示"该章节题库尚未完成"并跳过
- YAML 文件解析失败时，系统给出明确的文件路径和错误行号

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: 项目 MUST 提供 quiz CLI 命令，支持 `quiz <level>`（全章汇总）和 `quiz <level> <chapter>`（单章练习）两种模式
- **FR-002**: 每章 MUST 包含至少 3 道题，类型为选择题（4 选项）或判断题（正确/错误）
- **FR-003**: 题库数据 MUST 以 YAML 格式存储，每章一个独立文件，存放在 `docs/specs/002-quiz-system/questions/` 下
- **FR-004**: 答题系统 MUST 在用户选择答案后立即反馈对错，并显示解析
- **FR-005**: 答题结束 MUST 显示总分和学习建议（如"建议重新阅读接口章节"）
- **FR-006**: 题目 MUST 使用中文编写，技术术语保留英文原文（括号标注）
- **FR-007**: 汇总模式 MUST 支持按章节分类显示正确率
- **FR-008**: 项目 MUST 提供 quiz 子命令的 `--help` 文档
- **FR-009**: 题库 MUST 覆盖全部已完成的章节：基础 12 章 + 高级 9 章 + 实战 4 章
- **FR-010**: `docs/src/quiz/index.md` MUST 包含完整的题库目录（按基础/高级/实战分类，标注每章题目数量）

### Key Entities

- **题目 (Question)**: 包含题干、选项（选择题）或判断点（判断题）、正确答案、解析说明
- **题库文件 (Quiz File)**: 每章一个 YAML 文件，包含 ≥3 道题目
- **答题记录 (QuizSession)**: 记录用户答题过程、对错结果、总分和用时

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: 基础章节题库完整（12 章 × ≥3 题 = ≥36 道）
- **SC-002**: 高级章节题库完整（9 章 × ≥3 题 = ≥27 道）
- **SC-003**: 实战项目题库完整（4 章 × ≥3 题 = ≥12 道）
- **SC-004**: `quiz` 子命令能通过 `go run ./cmd/hello quiz --help` 查看帮助文档
- **SC-005**: `quiz <level> <chapter>` 模式能正确加载题库、交互式答题、显示得分
- **SC-006**: 题库 YAML 文件格式校验通过，无语法错误
- **SC-007**: `docs/src/quiz/index.md` 页面构建通过 mdBook 零错误

## Assumptions

- YAML 解析使用 `gopkg.in/yaml.v3` 库
- 用户已安装 Go 1.24+，能运行 `go run`
- 每章题库由人工编写（非自动生成），确保题目质量
- 题库与章节内容一一对应，题目基于"你会学到什么"章节目标设计

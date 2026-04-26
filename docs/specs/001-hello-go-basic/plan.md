# Implementation Plan: 完善三个 Overview 文档（basic / advance / awesome）

**Branch**: `001-hello-go-basic` | **Date**: 2026-04-26 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/docs/specs/001-hello-go-basic/spec.md`

## Summary

为 `docs/src/basic/basic-overview.md`、`docs/src/advance/advance-overview.md`、`docs/src/awesome/awesome-overview.md` 三个目前仅剩标题的页面补充完整内容。basic/advance 使用统一模板（学习目标 + 章节导航 + 学习路径建议 + 下一步导航），awesome 使用独立结构（项目导航 + 技术栈标签 + 应用场景）。全部为纯 mdBook 文档编写，不涉及代码改动。

## Technical Context

**Language/Version**: Go 1.24 (toolchain go1.24.3) — 仅文档编写，无代码变更  
**Primary Dependencies**: mdBook 0.4.52（文档构建）、现有 chapter Markdown 文件（内容来源）  
**Storage**: N/A — 纯文档项目  
**Testing**: mdBook build 验证（零 404、零构建错误）  
**Target Platform**: GitHub Pages 静态站点（浏览器访问）  
**Project Type**: documentation-only  
**Performance Goals**: mdBook 完整构建 <5 分钟（SC-005）  
**Constraints**: 中文编写，技术术语保留英文括号标注；字数限制 ~800 字（basic/advance）、~600 字（awesome）  
**Scale/Scope**: 3 个 Markdown 文件，约 2200 字总内容

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Principle I. Code Quality (N/A for docs-only)
- 无代码改动，不涉及编译器错误、lint 等。✅ 不适用

### Principle II. Test-First (adapted for docs)
- 文档"测试" = mdBook build 验证 + 链接检查
- 每个 overview 的内容需通过 `mdbook build` 验证 ✅ 将在实施阶段执行

### Principle III. User Experience Consistency ✅
- 中文主要语言 + 英文技术术语括号标注 — 符合 FR-003
- 章节结构遵循 12-section 模板 — overview 为导览页，独立结构
- 所有链接必须有效 — 将在 build 阶段验证

### Principle IV. Performance Requirements (docs) ✅
- mdbook build 时间 < 5 分钟 — 本次仅 3 个文件，预计 < 30 秒
- 无性能影响

### Principle V. SDD Harness Engineering ✅
- 本命令为 Phase 1（Planning），后续需经 `/speckit.tasks` → `/speckit.implement`
- Constitution Check 通过，无违规

**Gate Result**: PASS — 所有适用原则无违规

## Constitution Check (Post-Design Re-evaluation)

### Re-check after Phase 1 design:

- **Principle I (Code Quality)**: 仍不适用 — 无代码改动 ✅
- **Principle II (Test-First)**: 文档测试 = mdBook build ✓ 计划在实施阶段执行 ✅
- **Principle III (UX Consistency)**: 中文+英文术语、字数限制、结构模板均已明确 ✅
- **Principle IV (Performance)**: mdBook build <5 min，3 个文件预计 <30s ✅
- **Principle V (SDD)**: plan → tasks → implement 流程正确 ✅

**Re-evaluation Result**: PASS — 无新增违规，所有 gates 通过

## Project Structure

### Documentation (this feature)

```text
docs/specs/001-hello-go-basic/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output (entity mapping for docs)
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (N/A for docs-only)
└── tasks.md             # Phase 2 output (NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
docs/src/
├── basic/
│   └── basic-overview.md         # 待完善：~800字
├── advance/
│   └── advance-overview.md       # 待完善：~800字
├── awesome/
│   └── awesome-overview.md       # 待完善：~600字
└── SUMMARY.md                    # 导航结构（无需改动）
```

**Structure Decision**: 纯文档编写任务，仅修改 `docs/src/basic/basic-overview.md`、`docs/src/advance/advance-overview.md`、`docs/src/awesome/awesome-overview.md` 三个文件。不涉及源代码、测试或构建配置变更。

## Complexity Tracking

> No constitution violations. No complexity entries needed.

# Tasks: 完善三个 Overview 文档（basic / advance / awesome）

**Input**: Design documents from `/docs/specs/001-hello-go-basic/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, quickstart.md

**Tests**: N/A — docs-only feature, no code changes. Validation via `mdbook build` for link integrity.

**Organization**: Tasks grouped by user story (overview page) to enable independent editing and validation of each.

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- Documentation files under `docs/src/`
- Module: `hello` (from go.mod)

---

## Phase 1: Setup

**Purpose**: Verify prerequisite — confirm existing chapter content is available for reference

- [x] T001 Read existing chapter files in `docs/src/basic/` to extract chapter titles and key topics (12 chapters: variables, datatype, functions, flowcontrol, structs, interfaces, errorhandling, concurrency, generics, packages, pointers, logging, review-basic)
- [x] T002 Read existing chapter files in `docs/src/advance/` to extract chapter titles and key topics (10 chapters: smartpointers, errorhandling, reflection, database, web, testing, config, concurrency_advanced, context, review-advance)
- [x] T003 Read existing chapter files in `docs/src/awesome/` to extract project names and tech stacks (webservice: net/http + JSON, clidemo: Cobra CLI + CRUD)

---

## Phase 2: Foundational

**Purpose**: Confirm mdBook build passes before making changes (baseline verification)

**⚠️ CRITICAL**: Establish baseline before editing

- [x] T004 Run `cd docs && mdbook build` and confirm zero errors/warnings — record baseline
- [x] T005 [P] Read `docs/src/basic/basic-overview.md` — confirm currently single heading only
- [x] T006 [P] Read `docs/src/advance/advance-overview.md` — confirm currently single heading only
- [x] T007 [P] Read `docs/src/awesome/awesome-overview.md` — confirm currently single heading only

**Checkpoint**: Baseline confirmed. Three overview pages are empty (just `# Title`). mdBook build passes.

---

## Phase 3: User Story 1 - 基础入门 Overview (Priority: P1) 🎯 MVP

**Goal**: `basic-overview.md` populated with ~800 words: verifiable learning objectives, 12 chapter navigations with summaries + difficulty markers, learning path advice, and "next step" links to advance/awesome.

**Independent Test**: User opens `docs/src/basic/basic-overview.md`, can scan all 12 chapters, understand what they will learn, and know "what next" after completion.

### Implementation for User Story 1

- [x] T008 [US1] Write `## 概述` section (1-2 paragraphs) in `docs/src/basic/basic-overview.md` — introduce basic section's purpose and scope
- [x] T009 [US1] Write `## 你会学到什么` section with 12 specific verifiable learning objectives (e.g., "能用 `var` 或 `:=` 声明变量", "能编写 `for` 循环和 `switch` 分支") — one per sub-chapter
- [x] T010 [US1] Write `## 章节导航` section — list all 12 chapters with `./chapter.md` links, 1-2 sentence summary each, and 🔵🟡🔴 difficulty markers: variables(🔵), datatype(🔵), functions(🔵), flowcontrol(🔵), structs(🟡), interfaces(🟡), concurrency(🟡), generics(🟡), packages(🔵), pointers(🟡), logging(🔵), errorhandling(🟡)
- [x] T011 [US1] Write `## 学习路径建议` section — total estimated time (~2 hours), learning strategy tips, 🔵🟡🔴 explanation
- [x] T012 [US1] Write `## 下一步` section — links to [高级进阶](../advance/advance-overview.md) and [精选实战](../awesome/awesome-overview.md) with brief descriptions of what each offers
- [x] T013 [US1] Run `cd docs && mdbook build` and verify zero errors, links valid

**Checkpoint**: `basic-overview.md` ~800 words, all 12 chapters navigable, learning objectives verifiable, next-step links functional. mdBook build passes.

---

## Phase 4: User Story 2 - 高级进阶 Overview (Priority: P2)

**Goal**: `advance-overview.md` populated with ~800 words: prerequisite self-check quiz (3-5 questions), verifiable learning objectives, 8 chapter navigations with summaries + difficulty markers, learning path advice, and "next step" link to awesome.

**Independent Test**: User opens `docs/src/advance/advance-overview.md`, can self-check readiness, understand what they will learn, and navigate to awesome projects.

### Implementation for User Story 2

- [x] T014 [US2] Write `## 概述` section (1-2 paragraphs) in `docs/src/advance/advance-overview.md` — introduce advance section's purpose
- [x] T015 [US2] Write `## 前置知识自检` section — 4 self-check questions using `<details>/<summary>` HTML: Q1 goroutine/channel basics, Q2 interface implicit implementation, Q3 error pattern (`if err != nil`), Q4 struct + method basics
- [x] T016 [US2] Write `## 你会学到什么` section with 8+ specific verifiable learning objectives — one per sub-chapter (context, advanced concurrency, generics, reflection, performance, database, web, testing, config, smartpointers)
- [x] T017 [US2] Write `## 章节导航` section — list all 8+ chapters with `./chapter.md` links, 1-2 sentence summaries, and 🔵🟡🔴 difficulty markers: database(🟡), web(🟡), errorhandling(🟡), context(🔴), concurrency_advanced(🔴), reflection(🟡), testing(🟡), config(🔵), smartpointers(🟡)
- [x] T018 [US2] Write `## 学习路径建议` section — total estimated time (~3 hours), learning strategy, which chapters to prioritize based on goals
- [x] T019 [US2] Write `## 下一步` section — link to [精选实战](../awesome/awesome-overview.md) with description of what awesome projects apply advance concepts
- [x] T020 [US2] Run `cd docs && mdbook build` and verify zero errors, `<details>` tags render correctly

**Checkpoint**: `advance-overview.md` ~800 words, self-check quiz functional, learning objectives verifiable, mdBook build passes.

---

## Phase 5: User Story 3 - 精选实战 Overview (Priority: P3)

**Goal**: `awesome-overview.md` populated with ~600 words: 4 project navigations (name + tech stack tags + target audience/capability + 1-2 sentence summary), application scenario suggestions.

**Independent Test**: User opens `docs/src/awesome/awesome-overview.md`, can browse available projects, understand what tech each uses, and pick based on interest.

### Implementation for User Story 3

- [x] T021 [US3] Write `## 概述` section (1-2 paragraphs) in `docs/src/awesome/awesome-overview.md` — introduce awesome section's purpose, who it's for
- [x] T022 [US3] Write `## 实战项目` section — list 4 projects with:
  1. **Web 服务** — 🔧 net/http + JSON 中间件，适合练习 REST API 设计与并发处理
  2. **CLI 工具** — 🔧 Cobra + CRUD 操作，适合练习命令行架构与子命令路由
  3. **数据处理管道** — 🔧 goroutine pool + channel pipeline，适合练习并发模式与优雅关闭
  4. **工具链实践** — 🔧 Go 工具链深度使用，适合理解编译/测试/基准测试流程
- [x] T023 [US3] Write `## 应用场景建议` section — mapping of project types to real-world use cases (microservices, devops tools, data processing)
- [x] T024 [US3] Write `## 前置要求` section — what knowledge is needed (basic + advance completion recommended)
- [x] T025 [US3] Run `cd docs && mdbook build` and verify zero errors, all cross-links valid

**Checkpoint**: `awesome-overview.md` ~600 words, 4 projects listed with tech stacks and audience guidance, mdBook build passes.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final validation across all three overview pages

- [x] T026 Run `cd docs && mdbook build` full build and verify zero errors, zero warnings, zero 404 links
- [x] T027 [P] Verify all three overview pages have correct cross-links between them (basic→advance, basic→awesome, advance→awesome)
- [x] T028 [P] Verify word counts: basic ~800 (±20%), advance ~800 (±20%), awesome ~600 (±20%)
- [x] T029 [P] Run `mdbook serve docs/` locally and visually verify rendering (spacing, heading hierarchy, table rendering, details elements)
- [x] T030 Verify SUMMARY.md navigation still correct and overview pages appear as expected in sidebar

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **Foundational (Phase 2)**: No blocking dependencies — reads only, no blocking work
- **User Story 1 (Phase 3)**: Depends on Phase 1 content gathering — US1 can start after T001
- **User Story 2 (Phase 4)**: Depends on Phase 1 content gathering — US2 can start after T002
- **User Story 3 (Phase 5)**: Depends on Phase 1 content gathering — US3 can start after T003
- **Polish (Phase 6)**: Depends on all three user story phases being complete

### User Story Dependencies

- **US1 (P1)**: Can start after T001 — No dependencies on US2/US3
- **US2 (P2)**: Can start after T002 — No dependencies on US1/US3
- **US3 (P3)**: Can start after T003 — No dependencies on US1/US2
- All three user stories are **independently implementable and testable**

### Within Each User Story

- Overview section → Learning objectives → Chapter navigation → Learning path → Next steps
- mdBook build verification after each overview page

### Parallel Opportunities

- T001, T002, T003 can run in parallel (reading different directories)
- T005, T006, T007 can run in parallel (reading different files)
- Once Phase 1 & 2 complete, US1, US2, US3 can all be worked on in parallel
- T026's verification sub-tasks (T027-T030) can run in parallel

---

## Parallel Example: All Three Overviews

```bash
# After gathering chapter info (Phase 1), write all three in parallel:
Task US1: "Write basic-overview.md with objectives + 12 chapter nav + learning path + next steps"
Task US2: "Write advance-overview.md with self-check + objectives + 8 chapter nav + learning path + next steps"
Task US3: "Write awesome-overview.md with 4 project nav + tech stacks + scenarios + prerequisites"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Gather chapter info (T001-T003)
2. Complete Phase 2: Baseline build verification (T004-T007)
3. Complete Phase 3: Write `basic-overview.md` (T008-T013)
4. **STOP and VALIDATE**: `mdbook build` passes, links work, word count ~800
5. Preview in browser, verify content quality

### Incremental Delivery

1. Phase 1 + 2 → Chapter info gathered, baseline established
2. Add US1 (basic overview) → Build passes → Preview → Done for MVP
3. Add US2 (advance overview) → Build passes → Preview
4. Add US3 (awesome overview) → Build passes → Preview
5. Phase 6 → Cross-link verification + final polish

### Parallel Team Strategy

With multiple writers:

1. One person gathers chapter info (Phase 1)
2. Once T001-T003 done:
   - Writer A: basic-overview (US1)
   - Writer B: advance-overview (US2)
   - Writer C: awesome-overview (US3)
3. All three validate independently via `mdbook build`

---

## Notes

- [P] tasks = different files, no dependencies
- [US1/US2/US3] label maps task to specific user story for traceability
- Each overview page is independently completable and verifiable
- Word counts are soft targets: basic ~800, advance ~800, awesome ~600 (±20% acceptable)
- No code changes — constitution lint/test gates do not apply
- mdBook build is the sole quality gate: zero errors, zero warnings, zero broken links
- Constitution Principle III: All content in Chinese with English technical terms in parentheses

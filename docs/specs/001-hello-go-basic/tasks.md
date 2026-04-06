# Tasks: Go 编程语言学习样例结构

**Input**: Design documents from `/docs/specs/001-hello-go-basic/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Constitution Principle II requires test-first development. Each chapter package MUST have corresponding test files.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Single project**: `internal/`, `cmd/`, `docs/` at repository root
- Module: `hello` (from go.mod)
- All example code under `internal/basic/`, `internal/advance/`, `internal/awesome/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project structure and documentation scaffolding

- [ ] T001 Create directory structure: `internal/basic/`, `internal/advance/`, `internal/awesome/`, `internal/algo/`, `internal/leetcode/`
- [ ] T002 [P] Create mdBook chapter directories: `docs/src/basic/`, `docs/src/advance/`, `docs/src/awesome/`, `docs/src/algo/`, `docs/src/leetcode/`, `docs/src/quick_reference/`, `docs/src/quiz/`, `docs/src/projects/`
- [ ] T003 [P] Update `docs/src/SUMMARY.md` with full chapter navigation (mirror hello-rust structure)
- [ ] T004 [P] Create `docs/src/glossary.md`, `docs/src/faq.md`, `docs/src/CHANGELOG.md`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core CLI routing and chapter interface that ALL user stories depend on

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T005 Implement CLI subcommand router in `cmd/hello/main.go` using Cobra (commands: basic, advance, awesome, algo, leetcode, quiz)
- [X] T006 [P] Define chapter interface convention: each chapter package exposes `func Run()` — document in `docs/src/about-chapters.md`
- [X] T007 [P] Implement error handling helper in `cmd/hello/errors.go` (consistent error messages with context and suggestions)
- [X] T008 [P] Implement `--help` output for each subcommand level in `cmd/hello/help.go`
- [X] T009 Create placeholder chapter template in `internal/basic/variables/main.go` with `func Run()` and verify `go run ./cmd/hello basic variables` routes correctly
- [X] T010 Add test for CLI routing in `cmd/hello/main_test.go` (table-driven: verify each subcommand dispatches correctly)

**Checkpoint**: Foundation ready — `go run ./cmd/hello --help` shows all subcommands, `go run ./cmd/hello basic variables` runs placeholder successfully

---

## Phase 3: User Story 1 - 基础入门学习路径 (Priority: P1) 🎯 MVP

**Goal**: Deliver 12+ basic chapters covering Go core syntax, each with runnable code examples and Chinese documentation

**Independent Test**: User can run `go run ./cmd/hello basic <chapter>` for any basic chapter and see correct output matching documentation

### Implementation for User Story 1

- [X] T011 [US1] Create `internal/basic/variables/main.go` — 变量与表达式 (Variables & Expressions): var, const, :=, iota
- [X] T012 [US1] Create `docs/src/basic/variables.md` — 变量与表达式文档 (≥500 chars, ≥3 examples, ≥3 quiz questions)
- [X] T013 [US1] Create `internal/basic/variables/main_test.go` — table-driven tests for variables chapter
- [X] T014 [US1] Create `internal/basic/datatype/main.go` — 基础数据类型 (Data Types): int, float, bool, string, slice, map, time
- [X] T015 [US1] Create `docs/src/basic/datatype.md` — 数据类型文档
- [X] T016 [US1] Create `internal/basic/datatype/main_test.go`
- [X] T017 [US1] Create `internal/basic/functions/main.go` — 函数 (Functions): definition, parameters, return values, variadic, closures
- [X] T018 [US1] Create `docs/src/basic/functions.md`
- [X] T019 [US1] Create `internal/basic/functions/main_test.go`
- [X] T020 [US1] Create `internal/basic/flowcontrol/main.go` — 流程控制 (Flow Control): if/else, for, switch, defer
- [X] T021 [US1] Create `docs/src/basic/flowcontrol.md`
- [X] T022 [US1] Create `internal/basic/flowcontrol/main_test.go`
- [X] T023 [US1] Create `internal/basic/structs/main.go` — 结构体 (Structs): definition, fields, methods, embedding
- [X] T024 [US1] Create `docs/src/basic/structs.md`
- [X] T025 [US1] Create `internal/basic/structs/main_test.go`
- [X] T026 [US1] Create `internal/basic/interfaces/main.go` — 接口 (Interfaces): definition, implementation, empty interface
- [X] T027 [US1] Create `docs/src/basic/interfaces.md`
- [X] T028 [US1] Create `internal/basic/interfaces/main_test.go`
- [X] T029 [US1] Create `internal/basic/errorhandling/main.go` — 错误处理 (Error Handling): error interface, wrapping, sentinel errors
- [X] T030 [US1] Create `docs/src/basic/errorhandling.md`
- [X] T031 [US1] Create `internal/basic/errorhandling/main_test.go`
- [X] T032 [US1] Create `internal/basic/concurrency/main.go` — 并发 (Concurrency): goroutines, channels, select, sync.WaitGroup
- [X] T033 [US1] Create `docs/src/basic/concurrency.md`
- [X] T034 [US1] Create `internal/basic/concurrency/main_test.go`
- [X] T035 [US1] Create `internal/basic/generics/main.go` — 泛型 (Generics): type parameters, constraints, comparable
- [X] T036 [US1] Create `docs/src/basic/generics.md`
- [X] T037 [US1] Create `internal/basic/generics/main_test.go`
- [X] T038 [US1] Create `internal/basic/packages/main.go` — 包管理 (Packages): import, visibility, init(), go mod
- [X] T039 [US1] Create `docs/src/basic/packages.md`
- [X] T040 [US1] Create `internal/basic/packages/main_test.go`
- [X] T041 [US1] Create `internal/basic/pointers/main.go` — 指针 (Pointers): &, *, pointer receivers
- [X] T042 [US1] Create `docs/src/basic/pointers.md`
- [X] T043 [US1] Create `internal/basic/pointers/main_test.go`
- [X] T044 [US1] Create `internal/basic/logging/main.go` — 日志记录 (Logging): log, slog, structured logging
- [X] T045 [US1] Create `docs/src/basic/logging.md`
- [X] T046 [US1] Create `internal/basic/logging/main_test.go`
- [X] T047 [US1] Create `internal/basic/review/main.go` — 阶段复习：基础部分 (Review Basic)
- [X] T048 [US1] Create `docs/src/basic/review-basic.md`

**Checkpoint**: All 12+ basic chapters implemented. `go run ./cmd/hello basic <chapter>` works for each. `go test ./internal/basic/...` passes. mdBook builds with all links valid.

---

## Phase 4: User Story 2 - 高级进阶学习路径 (Priority: P2)

**Goal**: Deliver 8+ advance chapters covering Go advanced features: smart pointers patterns, error handling patterns, reflection, database, web development, testing, system programming

**Independent Test**: User can run `go run ./cmd/hello advance <chapter>` for any advance chapter and see correct output matching documentation

### Implementation for User Story 2

- [X] T049 [US2] Create `internal/advance/smartpointers/main.go` — 智能指针模式 (Smart Pointer Patterns): reference counting, pool patterns
- [X] T050 [US2] Create `docs/src/advance/smartpointers.md`
- [X] T051 [US2] Create `internal/advance/smartpointers/main_test.go`
- [X] T052 [US2] Create `internal/advance/errorhandling/main.go` — 高级错误处理 (Advanced Error Handling): errors.Is/As, custom types, stack traces
- [X] T053 [US2] Create `docs/src/advance/errorhandling.md`
- [X] T054 [US2] Create `internal/advance/errorhandling/main_test.go`
- [X] T055 [US2] Create `internal/advance/reflection/main.go` — 反射 (Reflection): reflect package, struct tags, dynamic method calls
- [X] T056 [US2] Create `docs/src/advance/reflection.md`
- [X] T057 [US2] Create `internal/advance/reflection/main_test.go`
- [X] T058 [US2] Create `internal/advance/database/database.go` — 数据库 (Database): GORM CRUD, migrations, relationships, transactions
- [X] T059 [US2] Create `docs/src/advance/database.md`
- [X] T060 [US2] Create `internal/advance/database/database_test.go`
- [X] T061 [US2] Create `internal/advance/web/web.go` — Web 开发 (Web Development): net/http, handlers, middleware, templates
- [X] T062 [US2] Create `docs/src/advance/web.md`
- [X] T063 [US2] Create `internal/advance/web/web_test.go`
- [X] T064 [US2] Create `internal/advance/testing/testing.go` — 测试最佳实践 (Testing Best Practices): table-driven, benchmarks, fuzzing
- [X] T065 [US2] Create `docs/src/advance/testing.md`
- [X] T066 [US2] Create `internal/advance/testing/testing_test.go`
- [X] T067 [US2] Create `internal/advance/config/main.go` — 配置管理 (Configuration): env vars, config files, viper patterns
- [X] T068 [US2] Create `docs/src/advance/config.md`
- [X] T069 [US2] Create `internal/advance/config/main_test.go`
- [X] T070 [US2] Create `internal/advance/review/main.go` — 阶段复习：高级进阶 (Review Advance)
- [X] T071 [US2] Create `docs/src/advance/review-advance.md`

**Checkpoint**: All 8+ advance chapters implemented. `go run ./cmd/hello advance <chapter>` works. `go test ./internal/advance/...` passes.

---

## Phase 5: User Story 3 - 精选实战项目 (Priority: P3)

**Goal**: Deliver 3+ complete project examples showing real-world Go application patterns

**Independent Test**: User can run `go run ./cmd/hello awesome <project>` and see a working application (CLI tool, web service, or data pipeline)

### Implementation for User Story 3

- [ ] T072 [US3] Create `internal/awesome/webservice/` — Web 服务项目: main.go, handler.go, router.go with RESTful API
- [ ] T073 [US3] Create `docs/src/awesome/webservice.md` — 项目文档 (architecture, setup, run instructions)
- [ ] T074 [US3] Create `internal/awesome/webservice/main_test.go`
- [ ] T075 [US3] Create `internal/awesome/clidemo/` — CLI 工具项目: Cobra-based CLI with subcommands
- [ ] T076 [US3] Create `docs/src/awesome/clidemo.md`
- [ ] T077 [US3] Create `internal/awesome/clidemo/main_test.go`
- [ ] T078 [US3] Create `internal/awesome/datapipeline/` — 数据处理管道: worker pool, channels, graceful shutdown
- [ ] T079 [US3] Create `docs/src/awesome/datapipeline.md`
- [ ] T080 [US3] Create `internal/awesome/datapipeline/main_test.go`

**Checkpoint**: All 3 awesome projects implemented and runnable. Each has README with architecture diagram and run instructions.

---

## Phase 6: Additional Modules (algo, leetcode, quiz, projects, appendix)

**Goal**: Complete remaining modules per spec: algorithms, LeetCode solutions, quiz system, projects, glossary, FAQ

- [ ] T081 Create `internal/algo/sort/main.go` — 排序算法 (Sorting): bubble, quick, merge, heap
- [ ] T082 Create `docs/src/algo/algo.md` — 算法实现文档
- [ ] T083 Create `internal/algo/sort/main_test.go`
- [ ] T084 Create `internal/algo/search/main.go` — 搜索算法 (Search): binary, BFS, DFS
- [ ] T085 Create `internal/algo/search/main_test.go`
- [ ] T086 Create `internal/algo/linkedlist/main.go` — 链表 (Linked List)
- [ ] T087 Create `internal/algo/linkedlist/main_test.go`
- [ ] T088 Create `internal/leetcode/twosum/main.go` — LeetCode #1 Two Sum
- [ ] T089 Create `docs/src/leetcode/leetcode.md` — LeetCode 题解文档 (≥5 problems)
- [ ] T090 Create `internal/leetcode/twosum/main_test.go`
- [ ] T091 Create `internal/leetcode/addtwonumbers/main.go` — LeetCode #2 Add Two Numbers
- [ ] T092 Create `internal/leetcode/addtwonumbers/main_test.go`
- [ ] T093 Create `internal/leetcode/longestsubstring/main.go` — LeetCode #3 Longest Substring Without Repeating
- [ ] T094 Create `internal/leetcode/longestsubstring/main_test.go`
- [ ] T095 Create `internal/leetcode/twosum4/main.go` — LeetCode #4 Median of Two Sorted Arrays
- [ ] T096 Create `internal/leetcode/twosum4/main_test.go`
- [ ] T097 Create `internal/leetcode/palindrome/main.go` — LeetCode #5 Longest Palindromic Substring
- [ ] T098 Create `internal/leetcode/palindrome/main_test.go`
- [ ] T099 Create `internal/quiz/engine.go` — 知识检查引擎 (Quiz Engine): question loading, scoring, feedback
- [ ] T100 Create `internal/quiz/engine_test.go`
- [ ] T101 Create `docs/src/quiz/index.md` — 题库首页 (≥3 questions per basic chapter)
- [ ] T102 Create `docs/src/quick_reference/snippets.md` — 代码片段速查
- [ ] T103 Create `docs/src/projects/README.md` — 项目实战索引
- [ ] T104 Create `docs/src/projects/todo-cli/README.md` — 项目: 命令行待办事项
- [ ] T105 Create `docs/src/projects/http-server/README.md` — 项目: 简易 HTTP 服务器
- [ ] T106 Create `docs/src/projects/web-scraper/README.md` — 项目: 并发爬虫
- [ ] T107 Create `docs/src/projects/binaries/README.md` — 项目: IPC 与分布式示例
- [ ] T108 Update `docs/src/glossary.md` — 术语表 (Go terminology Chinese/English)
- [ ] T109 Update `docs/src/faq.md` — 常见问题 FAQ

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Quality gates, documentation validation, and final cleanup

- [ ] T110 [P] Run `gofmt -l .` and fix all formatting issues
- [ ] T111 [P] Run `go vet ./...` and fix all issues
- [ ] T112 [P] Run `golangci-lint run ./...` and fix all issues (or nolint with justification)
- [ ] T113 Run `go test -cover ./...` and verify >80% coverage
- [ ] T114 Run `mdbook build docs/` and verify zero errors/warnings
- [ ] T115 [P] Update `docs/src/SUMMARY.md` with all final chapter links
- [ ] T116 [P] Verify all code examples include GitHub source links in documentation
- [ ] T117 Run `make build` and verify both `hello` and `foo` binaries compile
- [ ] T118 Run `go run ./cmd/hello --help` and verify all subcommands listed
- [ ] T119 [P] Update `README.md` with hello-go learning structure overview
- [ ] T120 [P] Update `AGENTS.md` with new package structure documentation

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion — BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational — P1 (MVP)
- **User Story 2 (Phase 4)**: Depends on Foundational — P2 (can run parallel to US3)
- **User Story 3 (Phase 5)**: Depends on Foundational — P3 (can run parallel to US2)
- **Additional Modules (Phase 6)**: Depends on Foundational — can run parallel to any US
- **Polish (Phase 7)**: Depends on all desired phases being complete

### User Story Dependencies

- **US1 (P1)**: Can start after Foundational — No dependencies on other stories
- **US2 (P2)**: Can start after Foundational — No dependencies on US1/US3
- **US3 (P3)**: Can start after Foundational — No dependencies on US1/US2
- **Additional Modules**: Can start after Foundational — independent of US1/US2/US3

### Within Each User Story

- Models/chapter code before documentation
- Documentation before quiz questions
- Core implementation before tests
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks (T001-T004) can run in parallel
- All Foundational tasks (T005-T010) — T005 must complete first, then T006-T010 in parallel
- Once Foundational phase completes, US1, US2, US3, and Additional Modules can all start in parallel
- Within each story: chapter code + doc + test for different chapters can run in parallel
- Different user stories can be worked on in parallel by different team members

---

## Parallel Example: User Story 1

```bash
# Launch multiple chapters in parallel:
Task: "T011 variables chapter code + T012 variables doc + T013 variables test"
Task: "T014 datatype chapter code + T015 datatype doc + T016 datatype test"
Task: "T017 functions chapter code + T018 functions doc + T019 functions test"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T004)
2. Complete Phase 2: Foundational (T005-T010) — CRITICAL, blocks all stories
3. Complete Phase 3: User Story 1 (T011-T048) — 12 basic chapters
4. **STOP and VALIDATE**: `go run ./cmd/hello basic <chapter>` works for all 12 chapters
5. **STOP and VALIDATE**: `go test ./internal/basic/...` passes
6. **STOP and VALIDATE**: `mdbook build docs/` passes with zero errors
7. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Test independently → Deploy/Demo
4. Add User Story 3 → Test independently → Deploy/Demo
5. Add Additional Modules → Test independently → Deploy/Demo
6. Each phase adds value without breaking previous phases

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (basic chapters)
   - Developer B: User Story 2 (advance chapters)
   - Developer C: User Story 3 (awesome projects)
   - Developer D: Additional Modules (algo, leetcode, quiz)
3. All phases complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [US1/US2/US3] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests pass after each chapter implementation
- Commit after each chapter or logical group
- Stop at any checkpoint to validate story independently
- Constitution Principle II: Tests MUST be written for each chapter package
- Constitution Principle III: All documentation MUST be Chinese with English technical terms
- Constitution Principle I: All code MUST pass `gofmt`, `go vet`, `golangci-lint`

# Tasks: 知识检查题库系统 (Quiz System)

**Input**: Design documents from `/docs/specs/002-quiz-system/`
**Prerequisites**: plan.md (required), spec.md (required for user stories)

**Tests**: Constitution Principle II requires test-first development for quiz engine.

**Organization**: Tasks grouped by development phase — engine → questions (basic → advance → awesome) → CLI integration → testing.

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- Engine code under `internal/quiz/`
- YAML questions under `docs/specs/002-quiz-system/questions/`
- CLI routing in `cmd/hello/`
- Docs under `docs/src/quiz/`

---

## Phase 1: Quiz Engine

**Purpose**: Core Go package for loading YAML questions, running interactive quiz sessions, and scoring.

- [ ] T001 [US1] Create `internal/quiz/types.go` — define `Question`, `QuizSession`, `QuizResult` structs (Question has: ID, chapter, type (multiple-choice/true-false), question text, options, answer index, explanation)
- [ ] T002 [US1] Create `internal/quiz/engine.go` — `LoadQuiz(path string) ([]Question, error)`, `RunQuiz(questions []Question, interactive bool) (*QuizResult, error)` — YAML loading via `gopkg.in/yaml.v3`
- [ ] T003 [US1] Implement interactive prompt in `engine.go` — display question, accept A/B/C/D input, validate input, show immediate feedback (✅/❌ + explanation)
- [ ] T004 [US1] Implement scoring in `engine.go` — track correct/total, calculate percentage, generate learning suggestions based on wrong answers
- [ ] T005 [P] [US1] Add `go.mod` dependency: `gopkg.in/yaml.v3` — run `go get gopkg.in/yaml.v3 && go mod tidy`
- [ ] T006 [US1] Create `internal/quiz/engine_test.go` — table-driven tests for: YAML loading (valid/invalid), scoring calculation, suggestion generation

**Checkpoint**: `go test ./internal/quiz/...` passes. Quiz engine can load a YAML file and run an interactive quiz session.

---

## Phase 2: Basic Chapter Questions (User Story 1) 🎯 MVP

**Goal**: Author ≥36 quiz questions covering all 12 basic chapters.

**Independent Test**: `go run ./cmd/hello quiz basic variables` loads and serves 3 questions.

### YAML Question Files for Basic Chapters

- [ ] T007 [US1] Create `docs/specs/002-quiz-system/questions/basic/variables.yaml` — 3+ questions on var/const/iota/type inference
- [ ] T008 [US1] Create `docs/specs/002-quiz-system/questions/basic/datatype.yaml` — 3+ questions on int/float/bool/string/slice/map/time
- [ ] T009 [US1] Create `docs/specs/002-quiz-system/questions/basic/functions.yaml` — 3+ questions on parameters/returns/closures/variadic
- [ ] T010 [US1] Create `docs/specs/002-quiz-system/questions/basic/flowcontrol.yaml` — 3+ questions on if/for/switch/defer
- [ ] T011 [US1] Create `docs/specs/002-quiz-system/questions/basic/structs.yaml` — 3+ questions on struct definition/methods/embedding
- [ ] T012 [US1] Create `docs/specs/002-quiz-system/questions/basic/interfaces.yaml` — 3+ questions on implicit implementation/empty interface/io.Writer
- [ ] T013 [US1] Create `docs/specs/002-quiz-system/questions/basic/concurrency.yaml` — 3+ questions on goroutines/channels/select
- [ ] T014 [US1] Create `docs/specs/002-quiz-system/questions/basic/generics.yaml` — 3+ questions on type parameters/constraints
- [ ] T015 [US1] Create `docs/specs/002-quiz-system/questions/basic/packages.yaml` — 3+ questions on visibility/init/import
- [ ] T016 [US1] Create `docs/specs/002-quiz-system/questions/basic/pointers.yaml` — 3+ questions on &/*/nil/value-vs-pointer
- [ ] T017 [US1] Create `docs/specs/002-quiz-system/questions/basic/logging.yaml` — 3+ questions on log/slog/levels/handlers
- [ ] T018 [US1] Create `docs/specs/002-quiz-system/questions/basic/errorhandling.yaml` — 3+ questions on errors.Is/As/sentinel errors
- [ ] T019 [US1] Create `docs/specs/002-quiz-system/questions/basic/review.yaml` — 3+ questions integrating basic concepts

**Question Quality Requirements (T007-T019)**:
- Each question MUST have: question text, 4 options (A-D for multiple-choice OR 2 options for true-false), correct answer index, explanation
- Explanations MUST be in Chinese with English technical terms
- Questions MUST be derived from chapter learning objectives (not arbitrary trivia)
- At least 1 true-false question per chapter (the rest can be multiple-choice)

**Checkpoint**: 36+ questions across 12 basic YAML files. Each file validates against Question struct.

---

## Phase 3: Advanced Chapter Questions (User Story 2)

**Goal**: Author ≥27 quiz questions covering all 9 advance chapters.

**Independent Test**: `go run ./cmd/hello quiz advance context` loads and serves 3+ questions.

### YAML Question Files for Advanced Chapters

- [ ] T020 [US2] Create `docs/specs/002-quiz-system/questions/advance/database.yaml` — 3+ questions on GORM CRUD/relationships/transactions
- [ ] T021 [US2] Create `docs/specs/002-quiz-system/questions/advance/web.yaml` — 3+ questions on handlers/middleware/JSON API
- [ ] T022 [US2] Create `docs/specs/002-quiz-system/questions/advance/errorhandling.yaml` — 3+ questions on errors.Is/As/%w/custom types
- [ ] T023 [US2] Create `docs/specs/002-quiz-system/questions/advance/context.yaml` — 3+ questions on WithCancel/WithTimeout/cancellation/goroutine leaks
- [ ] T024 [US2] Create `docs/specs/002-quiz-system/questions/advance/concurrency_advanced.yaml` — 3+ questions on Mutex/RWMutex/atomic/race detection
- [ ] T025 [US2] Create `docs/specs/002-quiz-system/questions/advance/reflection.yaml` — 3+ questions on reflect.Type/Value/struct tags/dynamic method calls
- [ ] T026 [US2] Create `docs/specs/002-quiz-system/questions/advance/testing.yaml` — 3+ questions on table-driven tests/benchmarks/fuzzing
- [ ] T027 [US2] Create `docs/specs/002-quiz-system/questions/advance/config.yaml` — 3+ questions on config loading/env override/reflection binding
- [ ] T028 [US2] Create `docs/specs/002-quiz-system/questions/advance/smartpointers.yaml` — 3+ questions on reference counting/sync.Pool/defer cleanup
- [ ] T029 [US2] Create `docs/specs/002-quiz-system/questions/advance/review.yaml` — 3+ questions integrating advance concepts

**Question Quality Requirements (T020-T029)**: Same as Phase 2.

**Checkpoint**: 27+ questions across 9 advance YAML files.

---

## Phase 4: Awesome Project Questions + Quiz Index Page (User Story 3)

**Goal**: Author ≥12 quiz questions for 4 awesome projects + build quiz index page.

**Independent Test**: Quiz index page renders correctly in mdBook with all chapter links.

- [ ] T030 [US3] Create `docs/specs/002-quiz-system/questions/awesome/webservice.yaml` — 3+ questions on REST/JSON/thread-safe middleware
- [ ] T031 [US3] Create `docs/specs/002-quiz-system/questions/awesome/clidemo.yaml` — 3+ questions on command parsing/subcommand routing/input validation
- [ ] T032 [US3] Create `docs/specs/002-quiz-system/questions/awesome/datapipeline.yaml` — 3+ questions on worker pool/channel pipeline/graceful shutdown
- [ ] T033 [US3] Create `docs/specs/002-quiz-system/questions/awesome/tooling.yaml` — 3+ questions on testing/benchmarking/profiling
- [ ] T034 [US3] Write `docs/src/quiz/index.md` — quiz index page with: overview paragraph, table of all 25 chapters by level, question count per chapter, links to individual question pages

**Checkpoint**: 12+ awesome questions + quiz index page. Total questions: ~75+. mdBook build passes for quiz/index.md.

---

## Phase 5: CLI Integration

**Purpose**: Wire quiz engine into existing `cmd/hello` CLI, add help docs, error handling.

- [ ] T035 Add `quiz` subcommand to `cmd/hello/main.go` — route to `hello quiz <level> [chapter]`
- [ ] T036 [P] Implement `--help` for quiz subcommand in `cmd/hello/help.go` — show usage: `quiz basic [chapter]`, `quiz advance`, `quiz random 10`
- [ ] T037 Implement error handling — invalid level/chapter shows helpful error message, missing YAML file gives clear "题库尚未完成" message
- [ ] T038 Implement summary mode (`quiz basic` without chapter arg) — load all questions for that level, randomize order, show per-chapter breakdown
- [ ] T039 Add `go run ./cmd/hello quiz` validation — verify `quiz basic variables` loads and runs 3 questions interactively

**Checkpoint**: `go run ./cmd/hello quiz --help` shows docs. `quiz basic variables` works. `quiz basic` runs 36-question summary.

---

## Phase 6: Testing, Validation & Polish

**Purpose**: Quality gates, cross-link verification, final cleanup.

- [ ] T040 Run `go build ./internal/quiz/...` — verify zero compilation errors
- [ ] T041 Run `go test -cover ./internal/quiz/...` — verify >80% coverage
- [ ] T042 Run YAML validation script — verify all 25 YAML files parse correctly and have ≥3 questions each
- [ ] T043 Run `mdbook build docs/` — verify quiz/index.md renders with zero errors
- [ ] T044 [P] Verify all quiz index page links in `docs/src/quiz/index.md` point to valid chapters
- [ ] T045 [P] Run `go run ./cmd/hello quiz basic` end-to-end test — manually verify at least 1 chapter flow
- [ ] T046 [P] Verify question quality — spot-check 5 questions per level for correctness and explanation clarity

---

## Dependencies & Execution Order

### Phase Dependencies

- **Engine (Phase 1)**: No dependencies — can start immediately
- **Basic Questions (Phase 2)**: Needs Phase 1 for YAML format validation, but question authoring can proceed independently
- **Advanced Questions (Phase 3)**: Same as Phase 2
- **Awesome Questions (Phase 4)**: Same as Phase 2 + depends on quiz index page design
- **CLI Integration (Phase 5)**: Depends on Phase 1 (engine) — can start after T004
- **Testing (Phase 6)**: Depends on all completion

### Parallel Opportunities

- Question authoring (T007-T019, T020-T029, T030-T033) can all happen in parallel across different chapters
- CLI integration (T035-T039) can start as soon as Phase 1 is done, independent of question writing
- Testing/validation tasks in Phase 6 can run in parallel

---

## Implementation Strategy

### MVP First (Phase 1 + 2)

1. Complete Phase 1: Quiz engine (T001-T006)
2. Complete Phase 2: Basic chapter questions (T007-T019)
3. Complete Phase 5: CLI wiring (T035-T039) — minimal integration
4. **STOP and VALIDATE**: `go run ./cmd/hello quiz basic variables` works
5. **STOP and VALIDATE**: 36 questions load and respond correctly
6. Deploy as MVP — basic quiz system is functional

### Incremental Delivery

1. Engine + Basic → Test independently → MVP
2. Add Advance Questions → Test → Expand coverage
3. Add Awesome Questions + Quiz Index → Test → Complete coverage
4. Polish + Full Test Suite → Final release

---

## Notes

- [P] tasks = different files, no dependencies
- [US1/US2/US3] label maps task to specific user story for traceability
- YAML format must match Question struct exactly — run `go test ./internal/quiz/...` after each YAML file creation
- Questions MUST derive from chapter learning objectives, not generic Go trivia
- All quiz content in Chinese with English technical terms in parentheses
- Constitution Principle II: "Knowledge checkpoint questions MUST validate learning outcomes"

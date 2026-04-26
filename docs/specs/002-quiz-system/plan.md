# Implementation Plan: 知识检查题库系统 (Quiz System)

**Branch**: `002-quiz-system` | **Date**: 2026-04-26 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/docs/specs/002-quiz-system/spec.md`

## Summary

为 hello-go 学习项目实现完整的知识检查题库系统：YAML 题库文件 + Go 问答引擎 + CLI 交互模式 + mdBook 题库首页。共需编写 ≥75 道题目（12 基础 + 9 高级 + 4 实战），支持单章练习和全章汇总两种答题模式。

## Technical Context

**Language/Version**: Go 1.24 (toolchain go1.24.3)  
**Primary Dependencies**: `gopkg.in/yaml.v3` (YAML 解析), Cobra (CLI), `math/rand` (随机出题)  
**Storage**: YAML 题库文件 per chapter, stored in `docs/specs/002-quiz-system/questions/`  
**Testing**: `go test` for quiz engine (table-driven tests for scoring, loading, feedback)  
**Target Platform**: CLI (macOS/Linux) + mdBook (GitHub Pages)  
**Project Type**: CLI tool + documentation  
**Performance Goals**: Quiz loading <100ms for 100+ questions  
**Constraints**: YAML format, Chinese language, ≥3 questions per chapter  
**Scale/Scope**: ~75 questions across 25 chapters, 1 Go package (quiz engine)

## Constitution Check

*GATE: Must pass before Phase 0 research.*

### Principle I. Code Quality (NON-NEGOTIABLE)
- ✅ Quiz engine MUST follow Go 1.24 idioms, pass `gofmt`, `go vet`
- ✅ Exported identifiers MUST have documentation comments
- ✅ No `panic()` in quiz engine — return errors instead

### Principle II. Test-First (NON-NEGOTIABLE)
- ✅ Quiz engine MUST have table-driven unit tests
- ✅ Questions MUST validate learning outcomes (mapped to chapter objectives)
- ✅ Test coverage target: >80% for quiz engine

### Principle III. User Experience Consistency
- ✅ All quiz content in Chinese with English technical terms in parentheses
- ✅ `--help` documentation for quiz subcommand

### Principle IV. Performance Requirements
- ✅ N/A for quiz (non-performance-critical CLI interaction)

### Principle V. SDD Harness Engineering
- ✅ Following spec → plan → tasks → implement workflow
- ✅ Constitution Check passed

**Gate Result**: PASS — All applicable principles satisfied.

## Project Structure

### Feature Artifacts

```text
docs/specs/002-quiz-system/
├── spec.md              # Feature specification
├── plan.md              # This file
└── tasks.md             # Phase 2 output
```

### Source Code

```text
internal/
└── quiz/
    ├── engine.go         # Quiz engine: load YAML, run quiz, score, feedback
    ├── types.go          # Question, QuizSession, QuizResult structs
    ├── engine_test.go    # Table-driven tests for engine
    └── questions/         # YAML question files (generated during implementation)
        ├── basic/
        │   ├── variables.yaml
        │   ├── datatype.yaml
        │   ├── ...
        │   └── review.yaml
        ├── advance/
        │   ├── database.yaml
        │   ├── ...
        │   └── review.yaml
        └── awesome/
            ├── webservice.yaml
            ├── ...
            └── tooling.yaml

cmd/hello/
└── cmd/root.go          # Add quiz subcommand routing

docs/src/
└── quiz/
    └── index.md         # Quiz index page (mdBook)
```

**Structure Decision**: Quiz engine as a standalone `internal/quiz` package. Question files in YAML alongside engine for easy authoring. CLI routing via existing Cobra setup in `cmd/hello/`.

## Complexity Tracking

> No constitution violations. Complexity: Moderate (~75 questions to author, quiz engine is straightforward).

## Phase Breakdown

- **Phase 1**: Quiz Engine — Go package with YAML loading, interactive quiz, scoring, feedback
- **Phase 2**: Basic Chapter Questions — Author 36 questions (12 chapters × 3)
- **Phase 3**: Advanced Chapter Questions — Author 27 questions (9 chapters × 3)
- **Phase 4**: Awesome Project Questions — Author 12 questions (4 projects × 3) + quiz index page
- **Phase 5**: CLI Integration — Wire quiz engine into cmd/hello, help docs, error handling
- **Phase 6**: Testing + Polish — Unit tests, mdBook verification, cross-link check

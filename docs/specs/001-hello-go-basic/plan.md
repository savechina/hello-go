# Implementation Plan: Go 编程语言学习样例结构

**Branch**: `001-hello-go-basic` | **Date**: 2026-04-05 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/docs/specs/001-hello-go-basic/spec.md`

## Summary

参考 hello-rust 的 basic/advance/awesome 三级分层结构，在 hello-go 中构建 Go 编程语言学习样例。使用单一 `go.mod`，所有示例代码作为 `internal/` 下的子包，由 `cmd/hello` 统一入口通过子命令调用。配套 mdBook 中文文档，覆盖基础语法、高级特性、实战项目、算法练习、LeetCode 题解、速查表和知识检查题库。

## Technical Context

**Language/Version**: Go 1.24 (toolchain go1.24.3)  
**Primary Dependencies**: Cobra (CLI), GORM (ORM), go-sqlite3 (SQLite), bbolt (BoltDB), testify (testing)  
**Storage**: SQLite (go-sqlite3, CGO), BoltDB (bbolt) — 仅用于示例演示  
**Testing**: `go test` with table-driven tests, `go test -cover` (>80%), `go test -bench`  
**Target Platform**: macOS / Linux, Go 1.24+  
**Project Type**: CLI tool + learning documentation (mdBook)  
**Performance Goals**: CLI startup <50ms, mdBook build <5min, binary size <20MB  
**Constraints**: 单一 `go.mod`，所有章节共享依赖；`internal/` 包不可被外部 import；CGO required for go-sqlite3  
**Scale/Scope**: 12+ basic chapters, 8+ advance chapters, 3+ awesome projects, 5+ algo, 5+ leetcode, 4+ projects

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Principle I: Code Quality ✅
- Go 1.24 idioms, `gofmt`/`go vet`/`golangci-lint` quality gates
- All exported identifiers documented with `//` comments
- No `panic()` in library code (only in `main` demo entry points)
- All code examples compile with `go build ./...`

### Principle II: Test-First Development ✅
- Table-driven tests for all exported functions
- >80% coverage target via `go test -cover`
- Each chapter example has corresponding test file

### Principle III: User Experience Consistency ✅
- CLI via Cobra with consistent `--help` output
- Chinese documentation with English technical terms in parentheses
- mdBook build passes with zero errors/warnings
- Each chapter: ≥500 Chinese chars, ≥3 code examples, ≥3 quiz questions

### Principle IV: Performance Requirements ✅
- CLI startup <50ms cold start
- No `time.Sleep()` in polling loops
- Memory <100MB for demo applications
- Binary <20MB with `-ldflags="-s -w"`

### Principle V: SDD Harness Engineering ✅
- Spec created via `/speckit.specify` ✅
- Plan created via `/speckit.plan` ✅
- Constitution check passed (all 5 principles) ✅
- Manual commit/push only — no automatic commits

**All gates passed. Proceeding to Phase 0.**

## Project Structure

### Documentation (this feature)

```text
docs/specs/001-hello-go-basic/
├── plan.md              # This file
├── research.md          # Phase 0 output
├── data-model.md        # Phase 1 output
├── quickstart.md        # Phase 1 output
├── contracts/           # Phase 1 output (CLI command schema)
└── tasks.md             # Phase 2 output (NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
hello-go/
├── cmd/
│   ├── hello/              # 主应用，统一入口
│   │   └── main.go         # Cobra root + 子命令路由到各章节
│   └── foo/                # 已有 Cobra CLI
│       └── main.go
├── internal/
│   ├── basic/              # 基础入门章节
│   │   ├── variables/      // package variables → func Run()
│   │   ├── expressions/    // package expressions → func Run()
│   │   ├── datatype/       // package datatype → func Run()
│   │   ├── functions/      // package functions → func Run()
│   │   ├── structs/        // package structs → func Run()
│   │   ├── interfaces/     // package interfaces → func Run()
│   │   ├── generics/       // package generics → func Run()
│   │   ├── concurrency/    // package concurrency → func Run()
│   │   └── ...             # (≥12 chapters total)
│   ├── advance/            # 高级进阶章节
│   │   ├── errorhandling/  // package errorhandling → func Run()
│   │   ├── reflection/     // package reflection → func Run()
│   │   ├── database/       // package database → func Run()
│   │   ├── web/            // package web → func Run()
│   │   ├── testing/        // package testing → func Run()
│   │   └── ...             # (≥8 chapters total)
│   ├── awesome/            # 精选实战
│   │   ├── webservice/     // package webservice → func Run()
│   │   ├── clidemo/        // package clidemo → func Run()
│   │   └── datapipeline/   // package datapipeline → func Run()
│   ├── domain/             # [已有] Domain models
│   ├── repository/         # [已有] SQLite/BoltDB repositories
│   └── business/           # [已有] Business logic interfaces
├── docs/                   # mdBook 文档
│   ├── src/
│   │   ├── basic/          # 基础章节文档
│   │   ├── advance/        # 高级章节文档
│   │   ├── awesome/        # 实战项目文档
│   │   ├── algo/           # 算法实现文档
│   │   ├── leetcode/       # LeetCode 题解文档
│   │   ├── quick_reference/# 代码片段速查
│   │   ├── quiz/           # 知识检查题库
│   │   ├── projects/       # 项目实战文档
│   │   ├── glossary.md     # 术语表
│   │   ├── faq.md          # 常见问题
│   │   └── SUMMARY.md      # 目录
│   └── book.toml
├── configs/                # [已有] Config helpers
├── examples/               # (可选) 独立运行示例
├── data/                   # [已有] Runtime data files
├── go.mod                  # 单一 module: hello
├── go.sum
├── Makefile                # [已有] Build system
└── AGENTS.md               # [已有] Project knowledge base
```

**Structure Decision**: 单一 `go.mod` 方案。`internal/basic/`, `internal/advance/`, `internal/awesome/` 各章节为独立子包，每个包暴露 `func Run()` 入口函数。`cmd/hello/main.go` 作为统一入口，通过子命令路由到对应章节的 `Run()` 函数。文档通过 mdBook 构建，与 hello-rust 保持一致。

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| N/A | All constitution principles pass | — |

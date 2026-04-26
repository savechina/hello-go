<!--
SYNC IMPACT REPORT
==================
Version Change: 1.1.1 → 2.0.0 (MAJOR: Complete rewrite — all principles redefined for Go)
Modified Principles:
  - I. Code Quality: Rust 2024 idioms → Go 1.24 idioms, gofmt/golangci-lint replace cargo clippy/fmt
  - II. Test-First: cargo test/proptest/criterion → go test/table-driven/benchmarks
  - III. UX Consistency: Rust CLI (clap) → Go CLI (Cobra), retained Chinese docs + mdBook
  - IV. Performance: Tokio async/go-sqlite3 → Go stdlib + GORM + BoltDB patterns
  - V. SDD Harness: Retained 8-phase workflow, updated tooling references for Go
Added Sections: None
Removed Sections:
  - All Rust-specific tooling references (Tokio, Axum, Tonic, SQLx, Diesel, proptest, criterion)
  - Rust-specific anti-patterns (unsafe, Arc, Cow, thread::sleep in async)
Templates Requiring Updates:
  - .specify/templates/plan-template.md: ⚠ pending (Language/Version field still shows Rust examples)
  - .specify/templates/tasks-template.md: ⚠ pending (task examples show Python paths)
  - .specify/templates/spec-template.md: ✅ aligned (language-agnostic)
  - .specify/templates/checklist-template.md: ✅ aligned (language-agnostic)
Follow-up TODOs:
  - TODO(PLAN_TEMPLATE): Update Language/Version examples from Rust to Go
  - TODO(TASKS_TEMPLATE): Update task path examples from Python to Go conventions
==================
-->

# Hello-Go Constitution

## Core Principles

### I. Code Quality (NON-NEGOTIABLE)

All code MUST prioritize clarity, simplicity, and idiomatic Go patterns.

**Requirements:**
- Follow Go 1.24 idioms and effective Go best practices
- Zero compiler errors or warnings (`go build ./...`)
- Mandatory `golangci-lint run` with all lints addressed or explicitly nolint'd with justification
- Maximum cyclomatic complexity: 15 per function (exceptions require architectural review)
- Documentation comments (`//`) on all exported identifiers (functions, types, constants, variables)
- No `panic()` in library code — return errors instead
- No blank identifier error suppression (`_ = someFunc()`) without documented rationale

**Documentation Quality:**
- All code examples MUST be from real project code (no fictional examples)
- Code examples MUST include GitHub source links
- All examples MUST compile successfully with `go build ./...`

**Rationale:** Learning resources must demonstrate best practices. Students learn from what they see. Poor code quality compounds as learners replicate patterns.

**Quality Gates:**
- `gofmt -l .` MUST return empty (no unformatted files)
- `go vet ./...` MUST pass with zero issues
- `golangci-lint run` MUST pass (or issues explicitly nolint'd with justification)
- `go test ./...` MUST pass
- All `TODO` and `FIXME` comments MUST have associated issues

### II. Test-First Development (NON-NEGOTIABLE)

Test-driven development is mandatory for all new features and bug fixes.

**Requirements:**
- Tests written and approved BEFORE implementation begins
- Red-Green-Refactor cycle strictly enforced
- Unit tests for all business logic (target: >80% coverage via `go test -cover`)
- Integration tests for all inter-service communication patterns
- Table-driven tests for functions with multiple input cases (idiomatic Go)
- Benchmarks for critical paths using `testing.B`

**Documentation Testing:**
- Documentation code snippets MUST be executable (`go test` examples)
- Knowledge checkpoint questions MUST validate learning outcomes
- All examples MUST have corresponding test coverage

**Rationale:** Tests serve as executable specifications and living documentation. They catch regressions and validate learning outcomes.

**Testing Tiers:**
1. **Unit Tests**: Fast, isolated, comprehensive (`go test -short` for quick runs)
2. **Integration Tests**: Repository boundaries, database interactions, file I/O
3. **End-to-End Tests**: Full CLI workflows using `gstack` browser automation
4. **Benchmarks**: Performance-critical paths via `go test -bench`

**Anti-Patterns:**
- `t.Skip()` without documented reasons and tracking issues
- Tests that depend on execution order
- Mocking internal implementation details instead of interfaces
- Test files importing `internal/` packages they shouldn't access

### III. User Experience Consistency

All user-facing interfaces MUST provide intuitive, consistent, and accessible experiences.

**Requirements:**
- CLI interfaces: Consistent argument parsing via Cobra, helpful error messages, progress indicators
- Repository interfaces: Clear contracts, consistent error types, documented return values
- Documentation: Chinese primary language with English technical terms, searchable, runnable examples
- Error Messages: Actionable, specific, include context and remediation steps
- Response Times: <100ms for CLI startup, <1s for database queries

**Documentation Language Standards:**
- **Primary Language**: Chinese (Simplified) with English technical terms in parentheses
  - Example: 接口 (interface), 结构体 (struct), 并发 (concurrency)
- **Writing Style**: Plain language, avoid academic jargon
- **Chapter Structure**: 12-section template for all tutorial chapters
- **Content Requirements**:
  - Minimum 500 Chinese characters per chapter
  - At least 3 executable code examples
  - At least 3 knowledge checkpoint questions
  - GitHub links to all source code examples

- **Overview Page Exemption**: Overview pages (`*-overview.md`) are navigation/index pages,
  not tutorial chapters. They are exempt from the 12-section template requirement and content
  rules above. Instead, they follow the summary template defined in research.md (Decision 7):
  概述 → 你会学到什么 → 章节导航 → 学习路径建议 → 下一步.

**UX Principles:**
- **Discoverability**: Every feature accessible via `--help` or API documentation
- **Predictability**: Consistent naming, argument order, and output formats
- **Recoverability**: Clear error messages with suggested fixes, no silent failures
- **Accessibility**: WCAG 2.1 AA compliance for web interfaces, screen reader compatible

**gstack Integration:**
- Use `/browse` for manual UX validation before deployment
- Use `/qa` for automated accessibility testing
- Use `/design-review` for visual consistency audits

**Documentation Build Quality:**
- `mdbook build` MUST pass with zero errors and warnings
- All links MUST be valid (no 404 errors)
- All code examples MUST render correctly with syntax highlighting

### IV. Performance Requirements

All code MUST meet defined performance standards and resource constraints.

**Requirements:**
- Memory: No unbounded allocations, explicit limits on data structures
- CPU: No blocking operations in goroutines that should be async (use buffered channels)
- I/O: Streaming for large datasets (no full materialization in memory)
- Database: Prepared statements via GORM, query plan validation, <10ms query latency
- Connection pooling: Properly configured for SQLite and BoltDB

**Performance Standards:**
- CLI startup: <50ms cold start, <10ms warm start
- Database queries: p95 <100ms, p99 <500ms
- Memory footprint: <100MB for demo applications
- Binary size: <20MB for CLI tools (release builds with `-ldflags="-s -w"`)

**Documentation Performance:**
- mdbook build time MUST be <5 minutes for full documentation
- Individual chapter builds MUST be <30 seconds
- Hot reload during development MUST be <2 seconds

**Performance Anti-Patterns:**
- `time.Sleep()` in polling loops (use `time.Ticker` or channels)
- Unnecessary allocations in hot paths (use `sync.Pool` where appropriate)
- Ignoring `defer` overhead in tight loops
- Loading entire database into memory instead of paginating

**Profiling Requirements:**
- Use `go tool pprof` for CPU and memory profiling before optimization
- Use `go test -bench` for benchmarking critical paths
- Document performance characteristics in AGENTS.md

### V. SDD Harness Engineering

Specification Driven Development (SDD) workflows MUST follow the **8-Phase Development Lifecycle** with triple quality gates (Metis + Momus + GStack).

**Development Phases:**

**Phase 0: Product Strategy & Requirements**
- `/office-hours` — Product discovery (YC 6-question forcing framework)
- `/plan-ceo-review` — Scope challenge (4 modes: SCOPE EXPANSION/SELECTIVE/HOLD/REDUCTION)
- `/speckit.specify` — Generate feature specifications
- **Quality Gate**: Metis intent analysis + Momus spec review (≥8/10)

**Phase 1: Technical Architecture & Design**
- `/speckit.plan` — Technical design with constitution check
- `/plan-eng-review` — Engineering review (architecture/data flow/performance)
- `/design-consultation` + `/plan-design-review` — Design system (UI projects)
- **Quality Gate**: Metis deep planning + Momus plan review (≥8/10)

**Phase 2: Task Decomposition**
- `/speckit.tasks` — Granular task breakdown (<4hr per task)
- `/speckit.analyze` — Cross-artifact consistency analysis
- **Quality Gate**: No CRITICAL/HIGH inconsistencies

**Phase 3: Quality Checklists**
- `/speckit.checklist` — Multi-domain checklists (test/security/ux/performance/code-quality/architecture/ai-safety)
- **Quality Gate**: 100% checklist coverage

**Phase 4: Implementation**
- `/speckit.implement` — Test-first execution with task delegation
- **Quality Gate**: `go vet` + `gofmt` + compilation success
- **Manual Review**: Changes MUST be manually reviewed before commit
- **Manual Commit**: ALL commits MUST be manually committed and pushed by user
- **Prohibited**: NO automatic commits or pushes to remote repositories

**Phase 5: Testing & Validation**
- `go test ./...` — Automated testing
- `/review` — Pre-landing PR review
- `/qa` — End-to-end QA testing with browser automation
- **Quality Gate**: 100% tests pass + no CRITICAL issues

**Phase 6: Delivery & Release**
- `/document-release` — Update all documentation
- `/ship` — Merge, version bump, create PR (with user approval)
- **Quality Gate**: All quality gates passed
- **Manual Verification**: User MUST verify all changes before deployment

**Phase 7: Retrospective**
- `/retro` — Engineering retro with trend analysis
- **Output**: Improvement action items for next iteration

**Triple Quality Gates:**

| Gate | Role | Timing | Purpose |
|------|------|--------|---------|
| **Metis** | Pre-planning consultant | Before each phase | Intent analysis, ambiguity detection, AI failure prediction, routing strategy |
| **Momus** | Post-delivery reviewer | After each phase | Clarity/verifiability/completeness/context evaluation, AI failure mode detection |
| **GStack** | Professional specialist | During execution | Domain-specific expertise (CEO review, eng review, design review, QA, PR review) |

**Skill Integration Matrix:**

| Phase | Speckit Commands | GStack Skills | OhMyOpenCode Agents |
|-------|------------------|---------------|---------------------|
| Phase 0 | `specify` | `office-hours`, `plan-ceo-review` | `metis`, `librarian` |
| Phase 1 | `plan` | `plan-eng-review`, `design-consultation`, `plan-design-review` | `metis`, `oracle`, `explore` |
| Phase 2 | `tasks`, `analyze` | - | `metis`, `momus` |
| Phase 3 | `checklist` | - | `momus` |
| Phase 4 | `implement` | - | `task()` delegation |
| Phase 5 | - | `review`, `qa`, `browse` | `momus` |
| Phase 6 | - | `document-release`, `ship` | `momus` |
| Phase 7 | - | `retro` | `momus` |

**Workflow Requirements:**
- Feature specifications via `/speckit.specify` (mandatory for all features)
- Implementation plans via `/speckit.plan` (mandatory before coding)
- Constitution check at Phase 1 (verify all 5 principles)
- Test-first development enforced in Phase 4
- All quality gates MUST pass before proceeding to next phase
- Document all decisions in `docs/specs/{N}-{feature}/`
- **Manual Control**: User MUST manually review, commit, and push all changes

**Workflow Enforcement:**
- No direct commits to `main` branch (use feature branches with PRs)
- All PRs MUST reference a spec document in `docs/specs/`
- All code changes MUST have corresponding test updates
- Breaking changes MUST update version according to semver and migration guide
- **CRITICAL**: NO automatic commits or pushes - user maintains full control

**Automation Standards:**
- CI pipeline: Build → Test → Lint → Benchmark → Deploy
- Deployment: Automated via GitHub Actions, rollback procedures documented
- Monitoring: Structured logging (`slog` or `zap`), metrics (`prometheus`), tracing (`jaeger`)
- Incident Response: Runbooks in `docs/runbooks/`, on-call rotation documented
- **Manual Gates**: User approval required at all deployment stages

**Tool Stack:**
- **Speckit Framework**: 8-phase SDD workflow (`specify`, `plan`, `tasks`, `analyze`, `checklist`, `implement`)
- **GStack Skills**: Quality automation (`office-hours`, `plan-ceo-review`, `plan-eng-review`, `design-consultation`, `plan-design-review`, `review`, `qa`, `browse`, `ship`, `retro`)
- **OhMyOpenCode Agents**: Triple quality gates (`metis` pre-planning, `momus` post-review, `oracle` architecture, `explore` codebase, `librarian` external research)
- **Go Tooling**: `golangci-lint`, `go test -cover`, `go tool pprof`, `govulncheck`
- **Manual Commit Policy**: ALL commits require user review and manual execution

## Technology Stack

**Core:**
- Language: Go 1.24 (toolchain go1.24.3)
- Module: `hello` (local-only, non-VCS-qualified)
- CLI Framework: Cobra 1.8+
- ORM: GORM 1.25+

**Data Layer:**
- SQLite: go-sqlite3 1.14+ (CGO required)
- BoltDB: bbolt 1.4+ (embedded key-value)
- Repository pattern: separate packages per driver (`sqlite/`, `boltdb/`)

**Testing:**
- Unit/Integration: `go test` with table-driven tests
- Coverage: `go test -cover` (target: >80%)
- Benchmarks: `go test -bench`
- Assertions: `testify` 1.10+

**Documentation:**
- mdBook 0.4.52 with plugins (admonish, alerts, pagetoc)
- Primary language: Chinese (Simplified) with English technical terms
- Deployed to: GitHub Pages
- Build verification: `mdbook build` MUST pass with zero errors

**Build/CI:**
- Build: Makefile with dynamic `cmd/*` discovery
- CGO: Enabled by default (required for go-sqlite3 on darwin)
- CI: GitHub Actions (build, test, lint)
- Dev mode: `make start` with `yolo` file watcher

## Development Workflow

### Feature Development Lifecycle

1. **Specification** (`/speckit.specify`)
   - Create feature spec in `docs/specs/<###-feature-name>/spec.md`
   - Define user stories, acceptance criteria, success metrics
   - Quality checklist validation

2. **Planning** (`/speckit.plan`)
   - Technical design document
   - Architecture decisions with rationale
   - Constitution check (verify compliance with all 5 principles)
   - Dependency analysis

3. **Implementation** (`/speckit.tasks` → `/speckit.implement`)
   - Granular task breakdown (<4hr per task)
   - Test-first: Write tests → Tests fail → Implement → Tests pass
   - **Manual Review**: ALL changes MUST be manually reviewed
   - **Manual Commit**: User MUST execute all git commits
   - **Manual Push**: User MUST execute all git pushes
   - **Prohibited**: NO automatic commits or pushes

4. **Quality Assurance** (`/qa`)
   - Automated testing (`go test ./...`)
   - Visual validation (`/browse`, `/design-review`)
   - Performance benchmarks
   - Security scan (`govulncheck`)

5. **Review** (`/review`)
   - Pre-landing code review
   - Constitution compliance check
   - Performance regression check
   - Documentation completeness

6. **Deploy**
   - Merge to `main` via PR
   - Automated CI/CD pipeline
   - Post-deploy monitoring
   - **Manual approval required at all stages**

### Branch Strategy

- `main`: Production-ready code, protected
- `<###-feature-name>`: Feature branches (sequential numbering from speckit)
- All branches MUST have associated spec document
- **Manual Control**: User decides when to create branches and merge

### Commit Conventions

 **DO NOT COMMIT and PUSH**
 **DO NOT COMMIT and PUSH**
 **DO NOT COMMIT and PUSH**

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`

**Manual Commit Examples:**
```bash
# User manually reviews and commits
git add <files>
# git commit -m "docs: amend constitution to v2.0.0"
# git push origin main
```

**Prohibited:**
```bash
# DO NOT automatically commit or push
# git commit -m "auto: ..."  # FORBIDDEN
# git push origin main       # FORBIDDEN without user approval
```

## Governance

**Authority:**
This constitution supersedes all other development practices and guides. In case of conflict with team conventions, constitution principles take precedence.

**Amendment Process:**
1. Propose amendment via GitHub issue with rationale
2. Architectural review for impact assessment
3. Team discussion and approval (consensus required)
4. Update constitution with version bump (MAJOR.MINOR.PATCH)
5. Propagate changes to all dependent templates and documentation
6. Announce changes to all contributors
7. **Manual Execution**: All constitutional amendments require manual user approval and commit

**Versioning Policy:**
- **MAJOR**: Backward incompatible principle removals or redefinitions
- **MINOR**: New principle/section added or materially expanded guidance
- **PATCH**: Clarifications, wording improvements, typo fixes

**Compliance Review:**
- All PRs MUST verify constitution compliance via `/review` command
- Complexity exceptions MUST be justified in PR description with architectural approval
- Violations of NON-NEGOTIABLE principles block merge
- **Manual Review**: User MUST verify all compliance checks before merge

**Runtime Guidance:**
- Use `AGENTS.md` for project-specific technical guidance
- Use `.specify/templates/` for workflow templates
- Use `docs/` for user-facing documentation

**Enforcement:**
- CI checks for linting, formatting, testing, security
- Mandatory code review for all changes
- Quarterly constitution review and update cycle
- **Manual Control**: User has final approval on all changes to main branch

---

**Version**: 2.0.0 | **Ratified**: 2026-04-03 | **Last Amended**: 2026-04-05

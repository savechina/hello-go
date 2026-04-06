# PROJECT KNOWLEDGE BASE

**Generated:** 2026-04-05
**Commit:** f656b3f
**Branch:** main

## OVERVIEW
Go learning/examples project. Two CLI apps (hello, foo) under cmd/, internal packages for domain/repo/business layers. SQLite + BoltDB repositories, GORM ORM, Cobra CLI.

## STRUCTURE
```
hello-go/
├── cmd/
│   ├── hello/     # Sample app — wires internal packages, prints version, demo repo usage
│   └── foo/       # Cobra CLI — root command + "add" subcommand
├── internal/
│   ├── domain/       # Domain models (User, Person)
│   ├── business/     # Business logic interfaces
│   ├── repository/   # Data layer — sqlite/ (3 files), boltdb/ (2 files)
│   ├── first/        # Example: internal package interaction pattern
│   └── version/      # Version constant
├── configs/       # Config helpers + constants
├── docs/          # mdBook-style docs (src/ chapters)
├── data/          # Runtime data files
├── pkg/           # Placeholder (.keep only) — unused
├── examples/      # Example code
├── test/          # Test utilities
├── tools/         # Dev tooling
├── build/         # Build output
├── configs/       # Config files
└── go.mod         # Module: "hello" (non-standard path, local-only)
```

## WHERE TO LOOK
| Task | Location | Notes |
|------|----------|-------|
| Add CLI command | cmd/foo/main.go | Cobra-based, add subcommands here |
| Add domain model | internal/domain/ | Pure models, no external deps |
| Add repository | internal/repository/{sqlite,boltdb}/ | Separate by DB driver |
| Change version | internal/version/version.go | Single constant |
| Build config | Makefile | Dynamic cmd/* compilation |
| Docs | docs/src/ | mdBook format, SUMMARY.md controls nav |

## CONVENTIONS
- Module name: `hello` (local-only, not VCS-qualified)
- Internal packages: domain → business → repository layering
- Repository pattern: separate packages per DB driver (sqlite/, boltdb/)
- Makefile: auto-discovers cmd/* subdirectories for builds
- Comments in go.mod: commented-out deps (viper) and section headers (// indirect)

## ANTI-PATTERNS (THIS PROJECT)
- Do NOT import from pkg/ — it's empty (.keep placeholder)
- Do NOT add business logic to cmd/ — keep wiring only
- Do NOT mix repository implementations in same package — one per DB driver
- Do NOT remove toolchain directive without verifying Go version compatibility

## UNIQUE STYLES
- go.mod has inline comments on require lines and section headers (// indirect)
- Makefile uses dynamic command discovery from cmd/ subdirectories
- docs/ uses mdBook structure (SUMMARY.md, chapter files)

## COMMANDS
```bash
make install    # go mod tidy + download
make build      # compile all cmd/* binaries → bin/
make test       # run all tests
make clean      # clean build artifacts
make start      # watch + auto-restart dev mode (requires yolo)
```

## NOTES
- Go 1.24 + toolchain go1.24.3 — verify toolchain availability before building
- CGO required for go-sqlite3 (cgo enabled by default on darwin)
- pkg/ directory is a placeholder — safe to remove if no library plans

## Active Technologies
- Go 1.24 (toolchain go1.24.3) + Cobra (CLI), GORM (ORM), go-sqlite3 (SQLite), bbolt (BoltDB), testify (testing) (001-hello-go-basic)
- SQLite (go-sqlite3, CGO), BoltDB (bbolt) — 仅用于示例演示 (001-hello-go-basic)

## Recent Changes
- 001-hello-go-basic: Added Go 1.24 (toolchain go1.24.3) + Cobra (CLI), GORM (ORM), go-sqlite3 (SQLite), bbolt (BoltDB), testify (testing)

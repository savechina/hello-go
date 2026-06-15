# Hello Go

English | [简体中文](README_zh.md)

[![Docs](https://img.shields.io/badge/docs-online-blue)](https://savechina.github.io/hello-go/)
[![License](https://img.shields.io/badge/license-Apache%202.0-green)](https://github.com/savechina/hello-go/blob/main/LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/savechina/hello-go?style=social)](https://github.com/savechina/hello-go)

[**Hello Go Tutorial**](https://renyan.org/hello/go) | [**GitHub Pages**](https://savechina.github.io/hello-go/)

A comprehensive sample project and tutorial for learning the Go programming language, from basic syntax to advanced production-grade applications.

## 📖 Online Tutorial

- 🌐 **Official Tutorial**: [renyan.org/hello/go](https://renyan.org/hello/go)
- 📚 **GitHub Pages**: [savechina.github.io/hello-go](https://savechina.github.io/hello-go/)

## 🚀 Quick Start

```bash
# Clone the project
git clone https://github.com/savechina/hello-go.git
cd hello-go

# Build and run the hello binary
make build
bin/hello

# Run the foo CLI app
bin/foo --help

# Run all tests
make test
```

## 📦 Project Modules

### Basic

Core Go syntax and concepts for beginners:

| Module | Content |
|--------|---------|
| [Variables & Expressions](docs/src/basic/variables.md) | Variable binding, zero values, short declaration |
| [Functions](docs/src/basic/functions.md) | Function definitions, multi-return, variadic params |
| [Data Types](docs/src/basic/datatype.md) | Integers, floats, booleans, strings, collections |
| [Control Flow](docs/src/basic/flowcontrol.md) | if/else, for loops, switch, defer, panic/recover |
| [Structs](docs/src/basic/structs.md) | Struct definitions, methods, embedding |
| [Interfaces](docs/src/basic/interfaces.md) | Interface definitions, implicit implementation, type assertion |
| [Concurrency](docs/src/basic/concurrency.md) | Goroutines, channels, select, sync primitives |
| [Generics](docs/src/basic/generics.md) | Type parameters, constraints, generic functions |
| [Packages](docs/src/basic/packages.md) | Module organization, import paths, visibility |
| [Pointers](docs/src/basic/pointers.md) | Pointer types, new, stack vs heap |
| [Logging](docs/src/basic/logging.md) | log package, structured logging, log levels |
| [Error Handling](docs/src/basic/errorhandling.md) | error interface, sentinel errors, wrapping |

### Advance

Deep dive into advanced Go features and ecosystem:

| Module | Content |
|--------|---------|
| [Database](docs/src/advance/database.md) | GORM, go-sqlite3, bbolt, connection pooling |
| [Web Development](docs/src/advance/web.md) | net/http, routing, middleware, RESTful APIs |
| [Error Handling](docs/src/advance/errorhandling.md) | Custom error types, error wrapping, sentinel patterns |
| [Context](docs/src/advance/context.md) | context.Context, cancellation, deadlines, propagation |
| [Advanced Concurrency](docs/src/advance/concurrency_advanced.md) | sync.Map, atomic, worker pools, fan-in/fan-out |
| [Reflection](docs/src/advance/reflection.md) | reflect package, type switching, dynamic invocation |
| [Testing](docs/src/advance/testing.md) | go test, table-driven tests, mocking, testify |
| [Config Management](docs/src/advance/config.md) | Environment variables, YAML, flags, viper patterns |
| [Smart Pointer Patterns](docs/src/advance/smartpointers.md) | sync.Pool, unsafe.Pointer, reference counting |

### Awesome

Production-grade examples and real-world patterns:

| Module | Content |
|--------|---------|
| [Web Service](docs/src/awesome/webservice.md) | RESTful API, middleware, database integration |
| [CLI Tool](docs/src/awesome/clidemo.md) | Cobra command-line framework, subcommands |
| [Data Pipeline](docs/src/awesome/datapipeline.md) | ETL pipelines, streaming, batch processing |
| [Tooling](docs/src/awesome/tooling.md) | Build tools, code generation, dev workflows |

### Algorithm & Practice

| Module | Content |
|--------|---------|
| [Algorithms](docs/src/algo/algo.md) | Sorting, searching, data structures |
| [LeetCode Solutions](docs/src/leetcode/leetcode.md) | Two Sum, Add Two Numbers, common problems |

## 🛠️ Tech Stack

- **Go 1.24** (toolchain go1.24.3)
- **CLI Framework**: Cobra
- **ORM**: GORM
- **Databases**: SQLite (go-sqlite3), BoltDB (bbolt)
- **Testing**: testify
- **Documentation**: mdBook

## 📋 Project Structure

```
hello-go/
├── cmd/
│   ├── hello/        # Sample app — wires internal packages, prints version
│   └── foo/          # Cobra CLI — root command + "add" subcommand
├── internal/
│   ├── domain/       # Domain models (User, Person)
│   ├── business/     # Business logic interfaces
│   ├── repository/   # Data layer — sqlite/, boltdb/
│   ├── basic/        # Basic tutorial example code
│   ├── advance/      # Advanced example code
│   ├── awesome/      # Production-grade examples
│   ├── first/        # Internal package interaction pattern
│   └── version/      # Version constant
├── configs/          # Config helpers + constants
├── docs/             # mdBook tutorial documentation
├── test/             # Test utilities
└── data/             # Runtime data files
```

## 📝 License

Apache License 2.0

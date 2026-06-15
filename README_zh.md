# Hello Go

[English](README.md) | 简体中文

[![Docs](https://img.shields.io/badge/docs-online-blue)](https://savechina.github.io/hello-go/)
[![License](https://img.shields.io/badge/license-Apache%202.0-green)](https://github.com/savechina/hello-go/blob/main/LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/savechina/hello-go?style=social)](https://github.com/savechina/hello-go)

[**Hello Go 教程**](https://renyan.org/hello/go) | [**GitHub Pages**](https://savechina.github.io/hello-go/)

这是学习 Go 编程语言的综合性样例工程，包含从基础语法到高级应用的完整教程和可运行代码示例。

## 📖 在线教程

- 🌐 **官方教程**: [renyan.org/hello/go](https://renyan.org/hello/go)
- 📚 **GitHub Pages**: [savechina.github.io/hello-go](https://savechina.github.io/hello-go/)

## 🚀 快速上手

```bash
# 克隆项目
git clone https://github.com/savechina/hello-go.git
cd hello-go

# 编译运行 hello 应用
make build
bin/hello

# 运行 foo CLI 工具
bin/foo --help

# 运行全部测试
make test
```

## 📦 工程模块

### 基础入门 (Basic)

涵盖 Go 核心语法和概念，适合初学者：

| 模块 | 内容 |
|------|------|
| [变量与表达式](docs/src/basic/variables.md) | 变量绑定、零值、短声明 |
| [函数基础](docs/src/basic/functions.md) | 函数定义、多返回值、可变参数 |
| [基础数据类型](docs/src/basic/datatype.md) | 整数、浮点数、布尔值、字符串、集合类型 |
| [控制流](docs/src/basic/flowcontrol.md) | if/else、for 循环、switch、defer、panic/recover |
| [结构体](docs/src/basic/structs.md) | 结构体定义、方法、类型嵌入 |
| [接口](docs/src/basic/interfaces.md) | 接口定义、隐式实现、类型断言 |
| [并发入门](docs/src/basic/concurrency.md) | goroutine、channel、select、同步原语 |
| [泛型](docs/src/basic/generics.md) | 类型参数、约束、泛型函数 |
| [包管理](docs/src/basic/packages.md) | 模块组织、导入路径、可见性 |
| [指针](docs/src/basic/pointers.md) | 指针类型、new、栈与堆 |
| [日志记录](docs/src/basic/logging.md) | log 包、结构化日志、日志级别 |
| [错误处理](docs/src/basic/errorhandling.md) | error 接口、哨兵错误、错误包装 |

### 高级进阶 (Advance)

深入 Go 高级特性和生态系统：

| 模块 | 内容 |
|------|------|
| [数据库](docs/src/advance/database.md) | GORM、go-sqlite3、bbolt、连接池 |
| [Web 开发](docs/src/advance/web.md) | net/http、路由、中间件、RESTful API |
| [错误处理](docs/src/advance/errorhandling.md) | 自定义错误类型、错误包装、哨兵模式 |
| [Context 上下文](docs/src/advance/context.md) | context.Context、取消、超时、传播 |
| [高级并发](docs/src/advance/concurrency_advanced.md) | sync.Map、atomic、工作池、扇入/扇出 |
| [反射](docs/src/advance/reflection.md) | reflect 包、类型切换、动态调用 |
| [测试](docs/src/advance/testing.md) | go test、表格驱动测试、Mock、testify |
| [配置管理](docs/src/advance/config.md) | 环境变量、YAML、命令行参数、viper 模式 |
| [智能指针模式](docs/src/advance/smartpointers.md) | sync.Pool、unsafe.Pointer、引用计数 |

### 实战精选 (Awesome)

生产级示例和真实世界模式：

| 模块 | 内容 |
|------|------|
| [Web 服务](docs/src/awesome/webservice.md) | RESTful API、中间件、数据库集成 |
| [CLI 工具](docs/src/awesome/clidemo.md) | Cobra 命令行框架、子命令 |
| [数据处理管道](docs/src/awesome/datapipeline.md) | ETL 管道、流处理、批处理 |
| [工具链实践](docs/src/awesome/tooling.md) | 构建工具、代码生成、开发工作流 |

### 算法与练习

| 模块 | 内容 |
|------|------|
| [算法实现](docs/src/algo/algo.md) | 排序、搜索、数据结构 |
| [LeetCode 题解](docs/src/leetcode/leetcode.md) | Two Sum, Add Two Numbers, 常见问题 |

## 🛠️ 技术栈

- **Go 1.24** (toolchain go1.24.3)
- **CLI 框架**: Cobra
- **ORM**: GORM
- **数据库**: SQLite (go-sqlite3), BoltDB (bbolt)
- **测试**: testify
- **文档**: mdBook

## 📋 项目结构

```
hello-go/
├── cmd/
│   ├── hello/        # 示例应用 — 组合内部包，输出版本信息
│   └── foo/          # Cobra CLI — 根命令 + "add" 子命令
├── internal/
│   ├── domain/       # 领域模型 (User, Person)
│   ├── business/     # 业务逻辑接口
│   ├── repository/   # 数据层 — sqlite/, boltdb/
│   ├── basic/        # 基础入门示例代码
│   ├── advance/      # 高级进阶示例代码
│   ├── awesome/      # 生产级实战示例
│   ├── first/        # 内部包交互模式
│   └── version/      # 版本常量
├── configs/          # 配置帮助函数 + 常量
├── docs/             # mdBook 教程文档
├── test/             # 测试工具
└── data/             # 运行时数据文件
```

## 📝 许可证

Apache License 2.0

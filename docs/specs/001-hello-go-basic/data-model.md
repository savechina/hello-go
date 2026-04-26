# Data Model: Go 学习样例结构

## Entities

### Chapter (学习章节)

| Field       | Type   | Description                    | Validation                |
| ----------- | ------ | ------------------------------ | ------------------------- |
| `Name`      | string | 章节名称（英文，如 "variables"） | 小写字母，无空格            |
| `Title`     | string | 章节标题（中文，如 "变量与表达式"） | 非空                      |
| `Level`     | string | 难度等级：basic/advance/awesome | 枚举值                    |
| `Package`   | string | Go 包名（如 `internal/basic/variables`） | 有效 Go import 路径       |
| `RunFunc`   | func   | 章节入口函数 `func Run()`      | 必须存在且可调用          |
| `DocPath`   | string | 对应文档路径（如 `docs/src/basic/variables.md`） | 文件存在且可构建          |
| `Examples`  | []Example | 代码示例列表                | ≥3 个                     |
| `Quiz`      | []Quiz    | 知识检查题列表              | ≥3 道                     |

### Example (代码示例)

| Field       | Type   | Description                    | Validation                |
| ----------- | ------ | ------------------------------ | ------------------------- |
| `Name`      | string | 示例名称                       | 非空                      |
| `Code`      | string | Go 源代码                      | 可编译运行，无 panic        |
| `Output`    | string | 预期输出                       | 与实际运行结果一致        |
| `SourceLink`| string | GitHub 源码链接                | 有效 URL                  |

### Quiz (知识检查题)

| Field       | Type     | Description                    | Validation                |
| ----------- | -------- | ------------------------------ | ------------------------- |
| `Question`  | string   | 题目（中文）                   | 非空                      |
| `Options`   | []string | 选项（A/B/C/D）                | 4 个选项                  |
| `Answer`    | string   | 正确答案                       | 在 Options 中             |
| `Explain`   | string   | 解析说明                       | 非空                      |

### Level (难度等级)

| Field       | Type   | Description                    | Validation                |
| ----------- | ------ | ------------------------------ | ------------------------- |
| `Name`      | string | 等级名称：basic/advance/awesome | 枚举值                    |
| `Title`     | string | 显示名称（如 "基础入门"）        | 非空                      |
| `MinChapters`| int   | 最少章节要求                   | basic≥12, advance≥8, awesome≥3 |

## Relationships

```
Level 1───N Chapter 1───N Example
                   1───N Quiz
```

## State Transitions

不适用 — 学习章节为静态内容，无状态机。

## Validation Rules

1. 每个章节 MUST 有对应的 `func Run()` 入口函数
2. 每个章节 MUST 有对应的 mdBook 文档文件
3. 每个章节的代码示例 MUST 通过 `go build` 编译
4. 每个章节的 Quiz MUST 有 ≥3 道题目
5. 所有文档 MUST 使用中文编写

---

## Overview-Specific Entities (added 2026-04-26)

### OverviewPage

| Field | Type | Description |
|-------|------|-------------|
| `File` | string | 文件路径，如 `docs/src/basic/basic-overview.md` |
| `Level` | string | `basic` / `advance` / `awesome` |
| `WordCountTarget` | int | 800 (basic/advance), 600 (awesome) |
| `Structure` | string | `unified` (basic/advance), `independent` (awesome) |

For validation rules specific to overview pages, see `FR-017` / `FR-018` / `FR-019` in spec.md.

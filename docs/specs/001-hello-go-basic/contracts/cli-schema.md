# CLI Command Contract: hello-go

## Overview

`hello-go` 是学习样例项目的统一入口，通过子命令路由到各章节示例。

## Command Schema

```
hello-go <level> <chapter> [flags]
```

### Levels

| Level     | Description  | Chapters                          |
| --------- | ------------ | --------------------------------- |
| `basic`   | 基础入门     | variables, datatype, functions, structs, interfaces, generics, concurrency, ... |
| `advance` | 高级进阶     | errorhandling, reflection, database, web, testing, ... |
| `awesome` | 精选实战     | webservice, clidemo, datapipeline, ... |
| `algo`    | 算法与练习   | sort, search, linkedlist, ... |
| `leetcode`| LeetCode 题解 | two-sum, add-two-numbers, ... |
| `quiz`    | 知识检查     | (无子命令，随机或按章节出题) |

### Examples

```bash
# 基础章节
hello-go basic variables
hello-go basic concurrency

# 高级章节
hello-go advance database
hello-go advance web

# 实战项目
hello-go awesome webservice

# 算法
hello-go algo sort

# LeetCode
hello-go leetcode two-sum

# 知识检查
hello-go quiz
hello-go quiz --chapter=variables

# 帮助
hello-go --help
hello-go basic --help
```

### Flags

| Flag          | Type   | Default | Description             |
| ------------- | ------ | ------- | ----------------------- |
| `--help, -h`  | bool   | false   | 显示帮助信息            |
| `--version`   | bool   | false   | 显示版本信息            |
| `--chapter`   | string | ""      | quiz 模式：指定章节范围 |

### Output Format

- 章节运行：直接输出到 stdout
- 错误信息：输出到 stderr，包含上下文和建议
- Quiz 模式：交互式问答，显示题目、选项、解析

### Error Handling

| Scenario                | Exit Code | Message                          |
| ----------------------- | --------- | -------------------------------- |
| 未知章节                | 1         | "未知章节: {level} {chapter}"    |
| 章节代码 panic           | 1         | 完整 stack trace                 |
| 缺少参数                | 1         | 显示 `--help` 内容               |
| 编译失败                | 1         | Go 编译器输出                    |

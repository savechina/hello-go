# CLI 工具实战（CLI Tool Demo）

## 开篇故事

想象你是一个餐厅服务员。客人来了，说"我要一份炒饭"，你记下订单交给厨房。客人又说"再加个蛋"，你补充订单。客人最后说"炒饭做好了，结账"，你完成整个服务流程。

CLI 工具就像这个服务员。用户输入命令（客人下单），程序解析命令（记录订单），执行对应操作（交给厨房），返回结果（服务完成）。这个过程看似简单，但要做到优雅和健壮，需要掌握几个核心技能。

命令行工具是开发者的日常伙伴。从 `git commit` 到 `docker run`，从 `npm install` 到 `go build`，每个命令背后都有相似的逻辑：解析参数、路由到对应功能、验证输入、处理错误。理解这些模式，你就能编写出像专业工具一样好用的 CLI 程序。

本章通过一个简单的待办事项工具示例，带你掌握 CLI 开发的四个核心技能：命令解析、子命令路由、输入验证、错误处理。这些都是编写生产级 CLI 工具的基础。

## 本章适合谁

- ✅ 已掌握 Go 基础语法，想编写第一个 CLI 工具的开发者
- ✅ 想理解 `git`、`docker` 等工具背后原理的学习者
- ✅ 需要为项目编写命令行管理脚本的工程师
- ✅ 准备学习 Cobra 等高级 CLI 库，但想先打好基础的技术人员

如果你还没有写过基本的 Go 程序，建议先完成基础章节。

## 你会学到什么

学完本章后，你将能够：

1. **解析命令参数**：理解 `os.Args` 的结构，提取用户输入
2. **实现子命令路由**：用 `switch` 语句实现类似 `git add`、`git commit` 的路由逻辑
3. **编写输入验证**：使用 `strings.TrimSpace` 清理用户输入，拒绝无效数据
4. **处理错误场景**：返回清晰错误信息，让用户知道问题所在

## 前置要求

在开始本章之前，请确保你已经掌握：

- Go 基础语法（函数、结构体、switch 语句）
- 字符串处理基础（strings 包）
- 错误处理基础（fmt.Errorf）
- 数组和切片的基本操作

## 第一个例子

让我们从最简单的 CLI 命令解析开始：

```go
package clidemo

import (
    "fmt"
    "strings"
)

func Run() {
    fmt.Println("=== 实战项目：CLI 工具 (CLI Tool) ===")

    // 示例1: 命令解析 (Command parsing)
    args := []string{"todo", "add", "Learn Go"}
    fmt.Printf("  示例1: 命令解析: %v\n", args)
    
    // 示例2: 子命令路由 (Subcommand routing)
    cmd := "add"
    switch cmd {
    case "add":
        fmt.Println("  示例2: 添加任务 - 'Learn Go'")
    case "list":
        fmt.Println("  示例2: 列出所有任务")
    case "done":
        fmt.Println("  示例2: 标记任务完成")
    default:
        fmt.Println("  示例2: 未知命令")
    }
    
    // 示例3: 输入验证 (Input validation)
    title := "  Learn Go  "
    trimmed := strings.TrimSpace(title)
    if trimmed == "" {
        fmt.Println("  示例3: 输入验证失败 - 标题不能为空")
    } else {
        fmt.Printf("  示例3: 输入验证通过 - '%s'\n", trimmed)
    }
    
    // 示例4: 错误处理 (Error handling)
    if err := validateInput(""); err != nil {
        fmt.Printf("  示例4: 错误处理 - %v\n", err)
    }
}
```

这个例子展示了 CLI 工具的四个核心步骤：
1. **解析参数**：从用户输入中提取命令和数据
2. **路由子命令**：根据命令名称执行不同逻辑
3. **验证输入**：清理和检查用户数据
4. **处理错误**：返回有意义的错误信息

## 原理解析

### 概念 1：命令解析（Command Parsing）

`os.Args` 是 Go 程序接收命令行参数的标准方式：

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    // os.Args[0] 是程序名本身
    // os.Args[1:] 是用户传递的参数
    
    fmt.Println("程序名:", os.Args[0])
    fmt.Println("参数数量:", len(os.Args)-1)
    
    if len(os.Args) > 1 {
        fmt.Println("第一个参数:", os.Args[1])
    }
    
    // 打印所有参数
    for i, arg := range os.Args {
        fmt.Printf("  os.Args[%d] = %s\n", i, arg)
    }
}
```

**运行示例**：
```bash
$ go run main.go todo add "Learn Go"
程序名: main.go
参数数量: 3
第一个参数: todo
  os.Args[0] = main.go
  os.Args[1] = todo
  os.Args[2] = add
  os.Args[3] = Learn Go
```

**关键点**：
- `os.Args[0]` 总是程序名（或 `go run` 时的源文件名）
- 用户参数从 `os.Args[1]` 开始
- 参数数量用 `len(os.Args)-1` 计算
- 参数是原始字符串，不解析类型

### 概念 2：子命令路由（Subcommand Routing）

专业的 CLI 工具通常有多个子命令，像 `git` 有 `add`、`commit`、`push` 等：

```go
func routeCommand(args []string) {
    if len(args) < 2 {
        fmt.Println("用法: todo <命令> [参数]")
        return
    }
    
    // args[1] 是子命令
    cmd := args[1]
    
    switch cmd {
    case "add":
        if len(args) < 3 {
            fmt.Println("用法: todo add <任务标题>")
            return
        }
        title := args[2]
        fmt.Printf("添加任务: %s\n", title)
        
    case "list":
        fmt.Println("列出所有任务:")
        // 这里可以调用 listTasks()
        
    case "done":
        if len(args) < 3 {
            fmt.Println("用法: todo done <任务ID>")
            return
        }
        id := args[2]
        fmt.Printf("完成任务: %s\n", id)
        
    case "help":
        printHelp()
        
    default:
        fmt.Printf("未知命令: %s\n", cmd)
        printHelp()
    }
}

func printHelp() {
    fmt.Println("可用命令:")
    fmt.Println("  todo add <标题>  - 添加新任务")
    fmt.Println("  todo list        - 列出所有任务")
    fmt.Println("  todo done <ID>   - 标记任务完成")
    fmt.Println("  todo help        - 显示帮助")
}
```

**设计原则**：
- 第一个参数通常是主命令（如 `todo`）
- 第二个参数是子命令（如 `add`、`list`）
- 子命令后面的参数是该命令的具体参数
- 每个子命令有自己的参数数量检查

### 概念 3：输入验证（Input Validation）

用户输入可能包含多余的空格、空值甚至恶意内容。验证是 CLI 工具健壮性的基础：

```go
func validateInput(s string) error {
    // 1. 去除首尾空格
    trimmed := strings.TrimSpace(s)
    
    // 2. 检查是否为空
    if trimmed == "" {
        return fmt.Errorf("input cannot be empty")
    }
    
    // 3. 检查长度限制（可选）
    if len(trimmed) > 100 {
        return fmt.Errorf("input too long, max 100 characters")
    }
    
    // 4. 检查非法字符（可选）
    forbiddenChars := []string{"<", ">", "&", "|"}
    for _, char := range forbiddenChars {
        if strings.Contains(trimmed, char) {
            return fmt.Errorf("input contains forbidden character: %s", char)
        }
    }
    
    return nil
}
```

**使用示例**：
```go
func addTask(title string) error {
    // 验证输入
    if err := validateInput(title); err != nil {
        return fmt.Errorf("添加任务失败: %v", err)
    }
    
    // 清理后的数据
    cleanTitle := strings.TrimSpace(title)
    
    // 执行业务逻辑
    fmt.Printf("成功添加任务: %s\n", cleanTitle)
    return nil
}
```

**验证的三个层次**：
1. **格式验证**：去除空格、检查长度、检查格式
2. **内容验证**：检查是否包含非法字符、敏感词
3. **业务验证**：检查是否符合业务规则（如 ID 是否存在）

### 概念 4：错误处理模式（Error Handling Pattern）

CLI 工具的错误处理要让用户能快速定位问题：

```go
func executeCommand(cmd string, args []string) error {
    switch cmd {
    case "add":
        if len(args) == 0 {
            // 清晰的错误信息
            return fmt.Errorf("add 命令需要参数: todo add <任务标题>")
        }
        title := args[0]
        if err := validateInput(title); err != nil {
            // 包装错误，保留原始信息
            return fmt.Errorf("添加任务失败: %w", err)
        }
        return addTask(title)
        
    case "done":
        if len(args) == 0 {
            return fmt.Errorf("done 命令需要参数: todo done <任务ID>")
        }
        id := args[0]
        // 尝试解析 ID
        taskID, err := strconv.Atoi(id)
        if err != nil {
            return fmt.Errorf("任务ID必须是数字: %w", err)
        }
        return markDone(taskID)
        
    default:
        return fmt.Errorf("未知命令: %s，使用 todo help 查看可用命令", cmd)
    }
}

// 错误处理策略
func handleError(err error) {
    if err == nil {
        return
    }
    
    // 1. 打印错误信息（给用户）
    fmt.Fprintf(os.Stderr, "错误: %v\n", err)
    
    // 2. 提供帮助建议（可选）
    if strings.Contains(err.Error(), "未知命令") {
        fmt.Println("运行 'todo help' 查看可用命令")
    }
    
    // 3. 返回非零状态码（给脚本）
    os.Exit(1)
}
```

**错误处理最佳实践**：
- 错误信息要清晰，说明问题和解决方法
- 使用 `%w` 包装错误，保留错误链
- 关键错误输出到 `os.Stderr`（标准错误流）
- 失败时返回非零状态码，方便脚本检测

## 常见错误

### 错误 1：忘记检查参数数量

```go
// ❌ 错误示例
func handleAdd(args []string) {
    title := args[0] // 如果 args 为空，会 panic
    addTask(title)
}

// ✅ 正确示例
func handleAdd(args []string) error {
    if len(args) == 0 {
        return fmt.Errorf("缺少参数，用法: todo add <任务标题>")
    }
    title := args[0]
    return addTask(title)
}
```

### 错误 2：不清理用户输入

```go
// ❌ 错误示例
func addTask(title string) {
    // 用户可能输入 "  hello  " 或 ""
    tasks = append(tasks, title) // 直接使用，不验证
}

// ✅ 正确示例
func addTask(title string) error {
    cleanTitle := strings.TrimSpace(title)
    if cleanTitle == "" {
        return fmt.Errorf("任务标题不能为空")
    }
    tasks = append(tasks, cleanTitle)
    return nil
}
```

### 错误 3：错误信息不够清晰

```go
// ❌ 错误示例
if err != nil {
    fmt.Println("错误") // 用户不知道发生了什么
}

// ✅ 正确示例
if err != nil {
    fmt.Fprintf(os.Stderr, "错误: %v\n", err)
    fmt.Println("提示: 使用 todo help 查看命令用法")
    os.Exit(1)
}
```

### 错误 4：不区分错误类型

```go
// ❌ 错误示例：所有错误一样处理
if err != nil {
    fmt.Println("出错了")
    os.Exit(1)
}

// ✅ 正确示例：区分用户错误和系统错误
if err != nil {
    if isUserError(err) {
        fmt.Fprintf(os.Stderr, "输入错误: %v\n", err)
        fmt.Println("请检查你的命令参数")
    } else {
        fmt.Fprintf(os.Stderr, "系统错误: %v\n", err)
        fmt.Println("请联系管理员或稍后重试")
    }
    os.Exit(1)
}
```

## 动手练习

### 练习 1：实现完整的 add 命令

编写一个 `handleAdd` 函数，完整实现添加任务的逻辑：
- 检查参数数量
- 验证输入（非空、长度限制）
- 清理空格
- 添加到任务列表

<details>
<summary>参考答案</summary>

```go
var tasks []string

func handleAdd(args []string) error {
    // 1. 检查参数数量
    if len(args) == 0 {
        return fmt.Errorf("用法: todo add <任务标题>")
    }
    
    // 2. 获取并清理输入
    title := strings.TrimSpace(args[0])
    
    // 3. 验证非空
    if title == "" {
        return fmt.Errorf("任务标题不能为空")
    }
    
    // 4. 验证长度
    if len(title) > 50 {
        return fmt.Errorf("任务标题过长，最多50个字符")
    }
    
    // 5. 添加任务
    tasks = append(tasks, title)
    fmt.Printf("添加成功: %s\n", title)
    return nil
}
```

</details>

### 练习 2：实现帮助命令

编写一个 `handleHelp` 函数，显示所有可用命令和用法说明。

<details>
<summary>参考答案</summary>

```go
func handleHelp() {
    fmt.Println("Todo CLI 工具 - 简单的任务管理")
    fmt.Println()
    fmt.Println("用法:")
    fmt.Println("  todo <命令> [参数]")
    fmt.Println()
    fmt.Println("可用命令:")
    fmt.Println("  add <标题>   添加新任务")
    fmt.Println("  list         列出所有任务")
    fmt.Println("  done <ID>    标记任务为已完成")
    fmt.Println("  help         显示此帮助信息")
    fmt.Println()
    fmt.Println("示例:")
    fmt.Println("  todo add \"学习 Go 语言\"")
    fmt.Println("  todo list")
    fmt.Println("  todo done 1")
}
```

</details>

### 练习 3：思考题 - 如何支持选项参数？

思考：如果要支持类似 `todo add --priority high "任务标题"` 的选项参数，应该如何设计参数解析逻辑？

提示：考虑以下问题：
- 如何区分选项（`--priority`）和普通参数？
- 如何处理带值的选项（`high`）和不带值的选项（`--verbose`）？
- 如何验证选项的有效性？

<details>
<summary>参考思路</summary>

```go
type Options struct {
    Priority string
    Verbose  bool
}

func parseOptions(args []string) (Options, []string, error) {
    opts := Options{}
    positionalArgs := []string{}
    
    i := 0
    for i < len(args) {
        arg := args[i]
        
        // 检查是否是选项（以 -- 开头）
        if strings.HasPrefix(arg, "--") {
            optionName := strings.TrimPrefix(arg, "--")
            
            switch optionName {
            case "priority":
                if i+1 >= len(args) {
                    return opts, nil, fmt.Errorf("--priority 需要值")
                }
                opts.Priority = args[i+1]
                i += 2 // 跳过选项和值
                
            case "verbose":
                opts.Verbose = true
                i += 1
                
            default:
                return opts, nil, fmt.Errorf("未知选项: --%s", optionName)
            }
        } else {
            // 普通参数
            positionalArgs = append(positionalArgs, arg)
            i += 1
        }
    }
    
    return opts, positionalArgs, nil
}

// 使用示例
func handleAdd(args []string) error {
    opts, positionalArgs, err := parseOptions(args)
    if err != nil {
        return err
    }
    
    if len(positionalArgs) == 0 {
        return fmt.Errorf("缺少任务标题")
    }
    
    title := positionalArgs[0]
    
    // 使用选项
    if opts.Verbose {
        fmt.Printf("添加任务: %s (优先级: %s)\n", title, opts.Priority)
    }
    
    return addTask(title, opts.Priority)
}
```

</details>

## 知识点总结

### 核心技能

| 技能 | 说明 | Go 实现 |
|------|------|---------|
| 命令解析 | 从用户输入提取参数 | `os.Args` 数组 |
| 子命令路由 | 根据命令名执行不同逻辑 | `switch` 语句 |
| 输入验证 | 清理和检查用户数据 | `strings.TrimSpace` + 条件判断 |
| 错误处理 | 返回清晰错误信息 | `fmt.Errorf` + `%w` 包装 |

### CLI 工具设计原则

1. **参数检查优先**：先检查数量，再检查内容，最后执行业务
2. **错误信息友好**：告诉用户问题是什么，如何解决
3. **输入必清理**：用 `TrimSpace` 去除多余空格
4. **帮助信息完善**：提供用法说明和示例
5. **状态码正确**：成功返回 0，失败返回非零

### 进阶方向

当你掌握这些基础后，可以继续学习：

1. **使用 Cobra 库**：专业的 CLI 框架，支持自动帮助生成、参数验证、子命令嵌套
2. **添加配置文件**：支持 `.todo.yaml` 等配置，持久化用户偏好
3. **实现交互模式**：支持用户选择、确认对话框等交互功能
4. **添加颜色输出**：使用颜色区分成功、警告、错误信息
5. **编写单元测试**：覆盖命令解析、输入验证、错误处理

## 工业界应用

### 场景：数据库管理 CLI

一个数据库管理工具可能包含：

```bash
dbtool connect --host localhost --port 3306 --user admin
dbtool query "SELECT * FROM users WHERE active = true"
dbtool backup --output users_backup.sql
dbtool restore --input backup.sql
```

这需要：
- 选项解析（`--host`、`--port`）
- 参数验证（SQL 语句检查）
- 多个子命令的路由
- 清晰的错误提示

### 场景：DevOps 工具

运维脚本可能需要：

```bash
deploytool deploy --env staging --version v1.2.3
deploytool rollback --env production --version v1.2.2
deploytool status --env all
```

这涉及：
- 环境参数验证
- 版本号格式检查
- 操作确认（安全提示）

## 小结

本章介绍了 CLI 工具开发的四个核心技能，这些是编写任何命令行程序的基础。

### 核心概念

- **命令解析**：通过 `os.Args` 获取用户输入
- **子命令路由**：用 `switch` 实现命令分发
- **输入验证**：清理空格、检查有效性
- **错误处理**：返回清晰信息、正确状态码

### 最佳实践

1. 总是检查参数数量，防止 panic
2. 清理用户输入，拒绝空值
3. 错误信息要包含解决建议
4. 提供完善的帮助信息
5. 区分用户错误和系统错误

### 下一步

- 学习 Cobra、urfave/cli 等专业库
- 为你的项目编写管理 CLI
- 参考 `git`、`docker` 等工具的设计

## 术语表

| 术语 | 英文 | 说明 |
|------|------|------|
| 命令行工具 | CLI Tool | 命令行界面程序 |
| 参数 | Arguments | 用户传递给程序的值 |
| 子命令 | Subcommand | 主命令下的具体操作，如 `git add` |
| 路由 | Routing | 根据命令名分发到对应处理逻辑 |
| 输入验证 | Input Validation | 检查用户输入是否有效 |
| 错误处理 | Error Handling | 处理程序异常情况 |
| 状态码 | Exit Code | 程序退出时的数字，0 表示成功 |
| 标准错误流 | stderr | 错误信息的输出通道 |
| 命令解析 | Command Parsing | 从字符串提取参数的过程 |

## 源码

完整示例代码位于：[internal/awesome/clidemo/clidemo.go](https://github.com/savechina/hello-go/blob/main/internal/awesome/clidemo/clidemo.go)
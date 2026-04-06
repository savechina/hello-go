# 关于章节约定

`hello-go` 使用统一的章节接口约定，让 CLI 可以把 `hello <level> <chapter>` 路由到对应的示例代码。

## 章节接口

每个章节包都应暴露同一个入口：

```go
func Run()
```

CLI 路由器只关心这个入口函数，不关心章节内部如何组织示例代码。

## 包命名约定

- 目录按学习层级组织：`internal/basic/`、`internal/advance/`、`internal/awesome/`、`internal/algo/`、`internal/leetcode/`、`internal/quiz/`
- 每个章节使用独立子目录，目录名同时就是 CLI 章节名
- 包名与目录名保持一致，例如 `internal/basic/variables` 使用 `package variables`

## 注册约定

章节通过 `init()` 把自己的 `Run` 函数注册到 CLI 章节表：

```go
package variables

import "hello/internal/chapters"

func init() {
	chapters.Register("basic", "variables", Run)
}

func Run() {
	// chapter demo
}
```

注册完成后，CLI 会按 `level + chapter` 查找并执行对应函数。

## 如何新增章节

1. 在对应层级下创建新目录，例如 `internal/basic/functions/`
2. 实现 `func Run()`
3. 在 `init()` 中调用 `chapters.Register("basic", "functions", Run)`
4. 在 `cmd/hello/main.go` 中添加该章节包的空白导入，确保程序启动时触发注册逻辑
5. 补充对应文档与测试

## 设计原因

- CLI 路由不需要维护大型 `switch` 语句
- 新章节只需要注册自己的入口函数
- 帮助信息可以根据注册表动态展示当前可用章节
- 测试可以直接向注册表注入占位章节，快速验证路由行为

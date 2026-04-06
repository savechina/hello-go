package clidemo

import (
	"fmt"
	"os"
	"strings"

	"hello/internal/chapters"
)

// Run demonstrates a CLI tool with subcommands and flags.
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

func validateInput(s string) error {
	if strings.TrimSpace(s) == "" {
		return fmt.Errorf("input cannot be empty")
	}
	return nil
}

func init() {
	// Use os to avoid unused import error
	_ = os.Args
	chapters.Register("awesome", "clidemo", Run)
}

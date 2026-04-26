package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"hello/internal/chapters"

	"github.com/spf13/cobra"
)

func configureHelp(rootCmd *cobra.Command) {
	helper := func(cmd *cobra.Command, args []string) {
		_, _ = fmt.Fprint(cmd.OutOrStdout(), renderHelp(cmd))
	}

	rootCmd.SetHelpFunc(helper)
	for _, cmd := range rootCmd.Commands() {
		cmd.SetHelpFunc(helper)
	}
}

func renderHelp(cmd *cobra.Command) string {
	if level := cmd.Annotations["level"]; level != "" {
		if level == "quiz" {
			return renderQuizLevelHelp(cmd)
		}
		return renderLevelHelp(cmd, level)
	}

	return renderRootHelp(cmd)
}

func renderRootHelp(cmd *cobra.Command) string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "%s\n\n", cmd.Long)
	fmt.Fprintf(&builder, "Usage:\n  %s\n\n", cmd.UseLine())
	builder.WriteString("Available commands:\n")
	for _, level := range commandLevels {
		subCmd, _, err := cmd.Find([]string{level})
		if err != nil || subCmd == nil {
			continue
		}

		fmt.Fprintf(&builder, "  %-10s %s\n", subCmd.Name(), subCmd.Short)
	}

	if flags := strings.TrimSpace(cmd.Flags().FlagUsages()); flags != "" {
		fmt.Fprintf(&builder, "\nFlags:\n%s", flags)
		if !strings.HasSuffix(flags, "\n") {
			builder.WriteString("\n")
		}
	}

	builder.WriteString("\nExamples:\n")
	builder.WriteString("  hello basic variables\n")
	builder.WriteString("  hello advance <chapter>\n")
	builder.WriteString("  hello quiz <chapter>\n")

	return builder.String()
}

func renderLevelHelp(cmd *cobra.Command, level string) string {
	var builder strings.Builder

	fmt.Fprintf(&builder, "%s\n\n", cmd.Long)
	fmt.Fprintf(&builder, "Usage:\n  hello %s <chapter>\n\n", level)
	builder.WriteString("Available chapters:\n")

	chapterNames := chapters.Names(level)
	if len(chapterNames) == 0 {
		builder.WriteString("  (none registered yet)\n")
	} else {
		for _, chapterName := range chapterNames {
			fmt.Fprintf(&builder, "  - %s\n", chapterName)
		}
	}

	if flags := strings.TrimSpace(cmd.Flags().FlagUsages()); flags != "" {
		fmt.Fprintf(&builder, "\nFlags:\n%s", flags)
		if !strings.HasSuffix(flags, "\n") {
			builder.WriteString("\n")
		}
	}

	fmt.Fprintf(&builder, "\nTip:\n  Run 'hello %s <chapter>' to execute a chapter demo.\n", level)

	return builder.String()
}

func renderQuizLevelHelp(_ *cobra.Command) string {
	var builder strings.Builder

	builder.WriteString("Run a quiz chapter or full quiz entry point from the quiz track.\n\n")
	builder.WriteString("Usage:\n  hello quiz <level>           (all chapters)\n  hello quiz <level> <chapter>   (single chapter)\n\n")
	builder.WriteString("Examples:\n  hello quiz basic\n  hello quiz basic variables\n\n")

	root := findProjectRoot()
	quizDir := filepath.Join(root, "docs", "specs", "002-quiz-system", "questions")

	builder.WriteString("Available levels and chapters:\n")

	entries, err := os.ReadDir(quizDir)
	if err != nil {
		builder.WriteString("  (题目录尚未创建)\n")
	} else {
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}
			level := entry.Name()
			levelDir := filepath.Join(quizDir, level)
			levelEntries, err := os.ReadDir(levelDir)
			if err != nil {
				continue
			}
			fmt.Fprintf(&builder, "  %s:\n", level)
			for _, file := range levelEntries {
				name := file.Name()
				if strings.HasSuffix(name, ".yaml") || strings.HasSuffix(name, ".yml") {
					chapter := strings.TrimSuffix(name, filepath.Ext(name))
					fmt.Fprintf(&builder, "    - %s\n", chapter)
				}
			}
		}
	}

	fmt.Fprintf(&builder, "\nTip:\n  Run 'hello quiz <level>' for a full quiz, or 'hello quiz <level> <chapter>' for a single chapter.\n")

	return builder.String()
}

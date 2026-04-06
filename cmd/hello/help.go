package main

import (
	"fmt"
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

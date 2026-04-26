package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	// Register all basic chapters via init()
	_ "hello/internal/basic/concurrency"
	_ "hello/internal/basic/datatype"
	_ "hello/internal/basic/errorhandling"
	_ "hello/internal/basic/flowcontrol"
	_ "hello/internal/basic/functions"
	_ "hello/internal/basic/generics"
	_ "hello/internal/basic/interfaces"
	_ "hello/internal/basic/logging"
	_ "hello/internal/basic/packages"
	_ "hello/internal/basic/pointers"
	_ "hello/internal/basic/review"
	_ "hello/internal/basic/structs"
	_ "hello/internal/basic/variables"

	// Register all advance chapters via init()
	_ "hello/internal/advance/concurrency_advanced"
	_ "hello/internal/advance/config"
	_ "hello/internal/advance/context"
	_ "hello/internal/advance/database"
	_ "hello/internal/advance/errorhandling"
	_ "hello/internal/advance/reflection"
	_ "hello/internal/advance/review"
	_ "hello/internal/advance/smartpointers"
	_ "hello/internal/advance/testing"
	_ "hello/internal/advance/web"

	// Register all awesome projects via init()
	_ "hello/internal/awesome/clidemo"
	_ "hello/internal/awesome/datapipeline"
	_ "hello/internal/awesome/webservice"

	"hello/internal/chapters"
	"hello/internal/quiz"
	"hello/internal/version"

	"github.com/spf13/cobra"
)

var levelDetails = map[string]struct {
	Short string
	Long  string
}{
	"basic": {
		Short: "Run a basic learning chapter",
		Long:  "Run an entry-level Go learning chapter from the basic track.",
	},
	"advance": {
		Short: "Run an advanced learning chapter",
		Long:  "Run an advanced Go learning chapter from the advance track.",
	},
	"awesome": {
		Short: "Run an awesome project chapter",
		Long:  "Run a hands-on project chapter from the awesome track.",
	},
	"algo": {
		Short: "Run an algorithm chapter",
		Long:  "Run an algorithm exercise chapter from the algo track.",
	},
	"leetcode": {
		Short: "Run a LeetCode chapter",
		Long:  "Run a LeetCode practice chapter from the leetcode track.",
	},
	"quiz": {
		Short: "Run a quiz chapter",
		Long:  "Run a quiz chapter or quiz entry point from the quiz track.",
	},
}

var commandLevels = []string{"basic", "advance", "awesome", "algo", "leetcode", "quiz"}

func main() {
	if err := newRootCmd().Execute(); err != nil {
		handleCommandError(err)
	}
}

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "hello",
		Short:         "hello-go learning project CLI",
		Long:          "hello-go learning project CLI routes commands to runnable chapter demos.",
		Version:       version.VERSION,
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return newCommandError(
					1,
					fmt.Sprintf("unknown command: %s", args[0]),
					"Run 'hello --help' to see available learning tracks.",
				)
			}

			return cmd.Help()
		},
	}

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetVersionTemplate("{{printf \"%s\\n\" .Version}}")

	for _, level := range commandLevels {
		if level == "quiz" {
			rootCmd.AddCommand(newQuizCmd())
		} else {
			rootCmd.AddCommand(newLevelCmd(level))
		}
	}

	configureHelp(rootCmd)

	return rootCmd
}

func newLevelCmd(level string) *cobra.Command {
	details := levelDetails[level]

	return &cobra.Command{
		Use:   fmt.Sprintf("%s <chapter>", level),
		Short: details.Short,
		Long:  details.Long,
		Args:  chapterArgs(level),
		RunE: func(cmd *cobra.Command, args []string) error {
			return dispatchChapter(level, args[0])
		},
		Annotations: map[string]string{
			"level": level,
		},
	}
}

func chapterArgs(level string) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return newCommandError(
				1,
				fmt.Sprintf("expected exactly one chapter argument for %s", level),
				fmt.Sprintf("Run 'hello %s --help' to list available chapters.", level),
			)
		}

		return nil
	}
}

func dispatchChapter(level string, chapter string) error {
	runner, ok := chapters.Lookup(level, chapter)
	if !ok {
		return newCommandError(
			1,
			fmt.Sprintf("unknown chapter: %s %s", level, chapter),
			fmt.Sprintf("Run 'hello %s --help' to list available chapters.", level),
		)
	}

	runner()
	return nil
}

func newQuizCmd() *cobra.Command {
	details := levelDetails["quiz"]

	return &cobra.Command{
		Use:   "quiz <level> [chapter]",
		Short: details.Short,
		Long:  details.Long,
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) < 1 || len(args) > 2 {
				return newCommandError(
					1,
					"quiz requires 1 or 2 arguments: <level> [chapter]",
					"Run 'hello quiz --help' for usage.",
				)
			}
			return nil
		},
		RunE: func(_ *cobra.Command, args []string) error {
			level := args[0]
			root := findProjectRoot()

			if len(args) == 2 {
				chapter := args[1]
				path := filepath.Join(root, "docs", "specs", "002-quiz-system", "questions", level, chapter+".yaml")
				questions, err := quiz.LoadQuiz(path)
				if err != nil {
					return newCommandError(
						1,
						fmt.Sprintf("题库尚未完成: %s", chapter),
						fmt.Sprintf("Run 'hello quiz %s --help' for available chapters.", level),
					)
				}
				if _, err := quiz.RunQuiz(questions, chapter); err != nil {
					return fmt.Errorf("quiz error: %w", err)
				}
				return nil
			}

			qDir := filepath.Join(root, "docs", "specs", "002-quiz-system", "questions", level)
			var questionsByChapter map[string][]quiz.Question = make(map[string][]quiz.Question)

			filepath.WalkDir(qDir, func(path string, d fs.DirEntry, err error) error {
				if err != nil {
					return nil
				}
				if d.IsDir() || !strings.HasSuffix(strings.ToLower(path), ".yaml") {
					return nil
				}
				chapter := strings.TrimSuffix(filepath.Base(path), ".yaml")
				if questions, err := quiz.LoadQuiz(path); err == nil {
					questionsByChapter[chapter] = questions
				}
				return nil
			})

			if len(questionsByChapter) == 0 {
				return newCommandError(1, fmt.Sprintf("没有找到 %s 级别的题目", level), "")
			}

			if _, err := quiz.RunSummaryMode(questionsByChapter); err != nil {
				return fmt.Errorf("quiz error: %w", err)
			}
			return nil
		},
		Annotations: map[string]string{
			"level": "quiz",
		},
	}
}

// findProjectRoot returns the directory containing go.mod, checking cwd first.
func findProjectRoot() string {
	cwd, err := os.Getwd()
	if err == nil {
		if _, err := os.Stat(filepath.Join(cwd, "go.mod")); err == nil {
			return cwd
		}
	}
	dir, _ := os.Getwd()
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return cwd
}

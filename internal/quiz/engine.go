package quiz

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// LoadQuiz loads questions from a YAML file.
func LoadQuiz(path string) ([]Question, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	var questions []Question
	if err := yaml.Unmarshal(data, &questions); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}

	return questions, nil
}

func formatOptions(q Question) string {
	optionLabels := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}
	var sb strings.Builder
	for i, opt := range q.Options {
		if i < len(optionLabels) {
			sb.WriteString(fmt.Sprintf("  %s) %s\n", optionLabels[i], opt))
		} else {
			sb.WriteString(fmt.Sprintf("  %d) %s\n", i+1, opt))
		}
	}
	return strings.TrimSuffix(sb.String(), "\n")
}

func promptQuestion(q Question, reader *bufio.Reader) (int, error) {
	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "%s\n", q.Question)
	fmt.Fprintln(os.Stdout, formatOptions(q))
	fmt.Fprintf(os.Stdout, "请输入选项编号 (0-%d): ", len(q.Options)-1)

	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return -1, fmt.Errorf("unexpected end of input")
			}
			return -1, fmt.Errorf("read input: %w", err)
		}

		input = strings.TrimSpace(input)
		if input == "" {
			fmt.Fprintf(os.Stdout, "无效输入，请输入 0-%d: ", len(q.Options)-1)
			continue
		}
		var choice int
		if _, err := fmt.Sscanf(input, "%d", &choice); err != nil {
			fmt.Fprintf(os.Stdout, "无效输入，请输入 0-%d: ", len(q.Options)-1)
			continue
		}
		if choice < 0 || choice >= len(q.Options) {
			fmt.Fprintf(os.Stdout, "无效输入，请输入 0-%d: ", len(q.Options)-1)
			continue
		}
		return choice, nil
	}
}

// RunQuiz runs an interactive quiz session, printing questions and reading user input from stdin.
func RunQuiz(questions []Question, chapterName string) (*QuizResult, error) {
	return RunQuizWithReader(questions, chapterName, os.Stdin)
}

// RunQuizWithReader runs an interactive quiz session using the provided reader for input.
func RunQuizWithReader(questions []Question, chapterName string, inputReader io.Reader) (*QuizResult, error) {
	return runQuizWithReader(questions, chapterName, bufio.NewReader(inputReader))
}

// runQuizWithReader is the internal implementation that uses a *bufio.Reader directly,
// allowing multiple calls to share the same reader (e.g., summary mode across chapters).
func runQuizWithReader(questions []Question, chapterName string, reader *bufio.Reader) (*QuizResult, error) {
	start := time.Now()

	result := &QuizResult{
		TotalQuestions: len(questions),
	}

	for i, q := range questions {
		fmt.Fprintf(os.Stdout, "--- 问题 %d/%d ---", i+1, len(questions))
		if q.Chapter != "" {
			fmt.Fprintf(os.Stdout, " [章节: %s]", q.Chapter)
		}
		fmt.Fprintln(os.Stdout)

		choice, err := promptQuestion(q, reader)
		if err != nil {
			return nil, err
		}

		if choice == q.Answer {
			result.CorrectAnswers++
			fmt.Fprintln(os.Stdout, "✅ 正确!")
		} else {
			result.IncorrectAnswers++
			chapter := q.Chapter
			if chapter == "" {
				chapter = chapterName
			}
			fmt.Fprintf(os.Stdout, "❌ 错误! 正确答案是 %d\n", q.Answer)
			if !hasSuggestion(result.Suggestions, chapter) {
				result.Suggestions = append(result.Suggestions, fmt.Sprintf("建议重新阅读：%s", chapter))
			}
		}
		if q.Explanation != "" {
			fmt.Fprintf(os.Stdout, "💡 %s\n", q.Explanation)
		}
	}

	elapsed := time.Since(start)
	result.ElapsedSeconds = elapsed.Seconds()

	if result.TotalQuestions > 0 {
		result.Percentage = float64(result.CorrectAnswers) / float64(result.TotalQuestions) * 100
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintf(os.Stdout, "=== 测验结束 ===\n")
	fmt.Fprintf(os.Stdout, "得分: %d/%d (%.1f%%)\n", result.CorrectAnswers, result.TotalQuestions, result.Percentage)
	fmt.Fprintf(os.Stdout, "用时: %.1f 秒\n", result.ElapsedSeconds)

	if len(result.Suggestions) > 0 {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "建议:")
		for _, s := range result.Suggestions {
			fmt.Fprintf(os.Stdout, "  - %s\n", s)
		}
	}

	return result, nil
}

func hasSuggestion(suggestions []string, chapter string) bool {
	target := fmt.Sprintf("建议重新阅读：%s", chapter)
	for _, s := range suggestions {
		if s == target {
			return true
		}
	}
	return false
}

// RunSummaryMode runs quiz across multiple chapters and returns per-chapter breakdown.
func RunSummaryMode(questionsByChapter map[string][]Question) (*QuizResult, error) {
	return RunSummaryModeWithReader(questionsByChapter, os.Stdin)
}

// RunSummaryModeWithReader runs quiz across multiple chapters using the provided reader.
func RunSummaryModeWithReader(questionsByChapter map[string][]Question, inputReader io.Reader) (*QuizResult, error) {
	result := &QuizResult{}

	chapters := make([]string, 0, len(questionsByChapter))
	for ch := range questionsByChapter {
		chapters = append(chapters, ch)
	}
	sort.Strings(chapters)

	// Use a single bufio.Reader for the entire session so it isn't re-buffered per chapter.
	reader := bufio.NewReader(inputReader)

	for _, chapter := range chapters {
		questions := questionsByChapter[chapter]
		chapterResult, err := runQuizWithReader(questions, chapter, reader)
		if err != nil {
			return nil, err
		}

		result.TotalQuestions += chapterResult.TotalQuestions
		result.CorrectAnswers += chapterResult.CorrectAnswers
		result.IncorrectAnswers += chapterResult.IncorrectAnswers
		result.ElapsedSeconds += chapterResult.ElapsedSeconds

		result.ChapterBreakdown = append(result.ChapterBreakdown, ChapterScore{
			Chapter: chapter,
			Correct: chapterResult.CorrectAnswers,
			Total:   chapterResult.TotalQuestions,
		})

		result.Suggestions = append(result.Suggestions, chapterResult.Suggestions...)
	}

	if result.TotalQuestions > 0 {
		result.Percentage = float64(result.CorrectAnswers) / float64(result.TotalQuestions) * 100
	}

	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "=== 总览 ===")
	fmt.Fprintf(os.Stdout, "总分: %d/%d (%.1f%%)\n", result.CorrectAnswers, result.TotalQuestions, result.Percentage)
	fmt.Fprintf(os.Stdout, "总用时: %.1f 秒\n", result.ElapsedSeconds)
	fmt.Fprintln(os.Stdout)
	fmt.Fprintln(os.Stdout, "章节明细:")
	for _, cs := range result.ChapterBreakdown {
		pct := float64(cs.Correct) / float64(cs.Total) * 100
		fmt.Fprintf(os.Stdout, "  %s: %d/%d (%.1f%%)\n", cs.Chapter, cs.Correct, cs.Total, pct)
	}

	if len(result.Suggestions) > 0 {
		fmt.Fprintln(os.Stdout)
		fmt.Fprintln(os.Stdout, "建议:")
		for _, s := range result.Suggestions {
			fmt.Fprintf(os.Stdout, "  - %s\n", s)
		}
	}

	return result, nil
}

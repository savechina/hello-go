package quiz

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadQuiz(t *testing.T) {
	dir := t.TempDir()

	validYAML := `- question: "What is the output of fmt.Println(len(\"hello\"))?"
  options:
    - "4"
    - "5"
    - "6"
    - "0"
  answer: 1
  explanation: "len(\"hello\") returns 5."
  chapter: "strings"
- question: "Is Go a compiled language?"
  options:
    - "正确"
    - "错误"
  answer: 0
  explanation: "Go is a compiled language."
  chapter: "basics"
`
	invalidYAML := `- question: "test"
  options: [a, b
  answer: 0
`
	emptyYAML := ""

	tests := []struct {
		name        string
		content     string
		wantErr     bool
		wantCount   int
	}{
		{
			name:      "valid yaml with 2 questions",
			content:   validYAML,
			wantErr:   false,
			wantCount: 2,
		},
		{
			name:      "invalid yaml",
			content:   invalidYAML,
			wantErr:   true,
			wantCount: 0,
		},
		{
			name:      "empty file",
			content:   emptyYAML,
			wantErr:   false,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(dir, tt.name+".yaml")
			require.NoError(t, os.WriteFile(path, []byte(tt.content), 0o644))

			questions, err := LoadQuiz(path)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, questions, tt.wantCount)
		})
	}
}

func TestRunQuiz_Scoring(t *testing.T) {
	questions := []Question{
		{
			Question:  "Q1: What is 1+1?",
			Options:   []string{"2", "3", "4", "5"},
			Answer:    0,
			Chapter:   "math",
		},
		{
			Question:  "Q2: What is 2+2?",
			Options:   []string{"1", "2", "3", "4"},
			Answer:    1,
			Chapter:   "math",
		},
	}

	tests := []struct {
		name       string
		input      string
		wantCorrect int
		wantTotal   int
	}{
		{
			name:       "all correct",
			input:      "0\n1\n",
			wantCorrect: 2,
			wantTotal:   2,
		},
		{
			name:       "one correct one wrong (answer 0 and 3)",
			input:      "0\n3\n",
			wantCorrect: 1,
			wantTotal:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := io.NopCloser(strings.NewReader(tt.input))
			var buf bytes.Buffer
			origStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			result, err := RunQuizWithReader(questions, "test", input)

			w.Close()
			io.Copy(&buf, r)
			os.Stdout = origStdout

			require.NoError(t, err)
			assert.Equal(t, tt.wantTotal, result.TotalQuestions)
			assert.Equal(t, tt.wantCorrect, result.CorrectAnswers)
			assert.Equal(t, tt.wantTotal-tt.wantCorrect, result.IncorrectAnswers)
		})
	}
}

func TestRunQuiz_Suggestions(t *testing.T) {
	questions := []Question{
		{
			Question:  "Q1: Which keyword defines an interface in Go?",
			Options:   []string{"interface", "type ... interface", "def interface", "kind interface"},
			Answer:    1,
			Chapter:   "interfaces",
		},
		{
			Question:  "Q2: What package provides HTTP functionality?",
			Options:   []string{"net/http", "http", "web/http", "os/http"},
			Answer:    0,
			Chapter:   "net",
		},
	}

	input := io.NopCloser(strings.NewReader("0\n3\n"))
	var buf bytes.Buffer
	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	result, err := RunQuizWithReader(questions, "mixed", input)

	w.Close()
	io.Copy(&buf, r)
	os.Stdout = origStdout

	require.NoError(t, err)
	assert.Equal(t, 0, result.CorrectAnswers)
	assert.Equal(t, 2, result.IncorrectAnswers)
	assert.Equal(t, 2, result.TotalQuestions)

	// Both answers wrong → both chapters should be suggested
	require.Len(t, result.Suggestions, 2)
	assert.Contains(t, result.Suggestions, "建议重新阅读：interfaces")
	assert.Contains(t, result.Suggestions, "建议重新阅读：net")
}

func TestRunSummaryMode_Duplicates(t *testing.T) {
	questionsByChapter := map[string][]Question{
		"basics": {
			{
				Question: "What is var?",
				Options:  []string{"variable", "constant", "function", "type"},
				Answer:   0,
				Chapter:  "basics",
			},
		},
		"interfaces": {
			{
				Question: "Which is an interface?",
				Options:  []string{"type T interface{}", "struct T{}", "func T()", "var T"},
				Answer:   0,
				Chapter:  "interfaces",
			},
		},
	}

	input := strings.NewReader("1\n1\n")
	var buf bytes.Buffer
	origStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	result, err := RunSummaryModeWithReader(questionsByChapter, input)

	w.Close()
	io.Copy(&buf, r)
	os.Stdout = origStdout

	require.NoError(t, err)
	assert.Equal(t, 2, result.TotalQuestions)
	assert.Equal(t, 0, result.CorrectAnswers)
	assert.Equal(t, 2, result.IncorrectAnswers)
	assert.Len(t, result.ChapterBreakdown, 2)
	assert.Len(t, result.Suggestions, 2)
}

// Package quiz defines data types for the quiz engine.
package quiz

// Question represents a single quiz question loaded from YAML.
type Question struct {
	Question    string   `yaml:"question"`              // The question text
	Options     []string `yaml:"options"`               // Answer options (4 for multiple choice, 2 for true/false)
	Answer      int      `yaml:"answer"`                // 0-based index of correct answer
	Explanation string   `yaml:"explanation"`           // Explanation shown after answering
	Type        string   `yaml:"type,omitempty"`        // "multiple-choice" or "true-false" (defaults to multiple-choice if options > 2)
	Chapter     string   `yaml:"chapter,omitempty"`     // Chapter name for grouping and suggestions
}

// ChapterScore holds quiz results for a single chapter.
type ChapterScore struct {
	Chapter string
	Correct int
	Total   int
}

// QuizResult holds the overall quiz session result.
type QuizResult struct {
	TotalQuestions   int            // Total number of questions in the session
	CorrectAnswers   int            // Number of correct answers
	IncorrectAnswers int            // Number of incorrect answers
	Percentage       float64        // Percentage of correct answers
	ElapsedSeconds   float64        // Total time taken in seconds
	ChapterBreakdown []ChapterScore // Per-chapter score breakdown
	Suggestions      []string       // Learning suggestions based on wrong answers
}

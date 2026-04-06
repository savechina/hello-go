package review

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"hello/internal/chapters"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var errInvalidJSON = errors.New("invalid json payload")

type reviewConfig struct {
	ServiceName string `json:"service_name" required:"true"`
	ListenPort  int    `json:"listen_port" min:"1"`
	StorageDSN  string `json:"storage_dsn" required:"true"`
}

type validationError struct {
	Problems []string
}

type courseRecord struct {
	ID         uint `gorm:"primaryKey"`
	Title      string
	Instructor string
}

type courseInput struct {
	Title      string `json:"title" required:"true"`
	Instructor string `json:"instructor" required:"true"`
}

type reviewApp struct {
	config reviewConfig
	db     *gorm.DB
}

func init() {
	chapters.Register("advance", "review", Run)
}

// Run prints runnable examples for the advanced review chapter.
func Run() {
	for _, line := range renderExamples() {
		fmt.Println(line)
	}
}

func renderExamples() []string {
	return []string{
		exampleReflectionValidation(),
		exampleDatabaseWorkflow(),
		exampleWebHandlerErrorHandling(),
	}
}

func exampleReflectionValidation() string {
	err := validateStruct(reviewConfig{ServiceName: "", ListenPort: 0, StorageDSN: ""})
	if err == nil {
		return "[review] example 1 reflection validation => ok"
	}
	return fmt.Sprintf("[review] example 1 reflection validation => %s", err.Error())
}

func exampleDatabaseWorkflow() string {
	app, cleanup, err := newExampleReviewApp()
	if err != nil {
		return fmt.Sprintf("[review] example 2 database workflow error: %v", err)
	}
	defer cleanup()

	if err := app.createCourse(courseInput{Title: "Advanced Review", Instructor: "gopher"}); err != nil {
		return fmt.Sprintf("[review] example 2 database workflow error: %v", err)
	}

	count, err := app.courseCount()
	if err != nil {
		return fmt.Sprintf("[review] example 2 database workflow error: %v", err)
	}

	return fmt.Sprintf("[review] example 2 database workflow => service=%s courses=%d", app.config.ServiceName, count)
}

func exampleWebHandlerErrorHandling() string {
	app, cleanup, err := newExampleReviewApp()
	if err != nil {
		return fmt.Sprintf("[review] example 3 web handler error: %v", err)
	}
	defer cleanup()

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/courses", strings.NewReader(`{"instructor":"gopher"}`))
	app.courseHandler().ServeHTTP(recorder, request)

	return fmt.Sprintf("[review] example 3 web handler => status=%d body=%s", recorder.Code, strings.TrimSpace(recorder.Body.String()))
}

func newExampleReviewApp() (*reviewApp, func(), error) {
	dir, err := os.MkdirTemp("", "hello-review-*")
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		_ = os.RemoveAll(dir)
	}

	app, err := newReviewApp(reviewConfig{
		ServiceName: "review-service",
		ListenPort:  8080,
		StorageDSN:  filepath.Join(dir, "review.db"),
	})
	if err != nil {
		cleanup()
		return nil, nil, err
	}

	return app, cleanup, nil
}

func newReviewApp(cfg reviewConfig) (*reviewApp, error) {
	if err := validateStruct(cfg); err != nil {
		return nil, fmt.Errorf("validate review config: %w", err)
	}

	db, err := gorm.Open(sqlite.Open(cfg.StorageDSN), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("open review database: %w", err)
	}

	if err := db.AutoMigrate(&courseRecord{}); err != nil {
		return nil, fmt.Errorf("migrate review database: %w", err)
	}

	return &reviewApp{config: cfg, db: db}, nil
}

func (a *reviewApp) createCourse(input courseInput) error {
	if err := validateStruct(input); err != nil {
		return fmt.Errorf("validate course input: %w", err)
	}

	record := courseRecord{Title: input.Title, Instructor: input.Instructor}
	if err := a.db.Create(&record).Error; err != nil {
		return fmt.Errorf("create review course: %w", err)
	}

	return nil
}

func (a *reviewApp) courseCount() (int64, error) {
	var count int64
	if err := a.db.Model(&courseRecord{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("count review courses: %w", err)
	}
	return count, nil
}

func (a *reviewApp) courseHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
			return
		}

		var input courseInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": classifyError(fmt.Errorf("decode request: %w", errInvalidJSON))})
			return
		}

		if err := a.createCourse(input); err != nil {
			statusCode := classifyStatus(err)
			writeJSON(w, statusCode, map[string]any{"error": classifyError(err)})
			return
		}

		writeJSON(w, http.StatusCreated, map[string]any{
			"status":  "created",
			"service": a.config.ServiceName,
			"title":   input.Title,
		})
	})
}

func writeJSON(w http.ResponseWriter, statusCode int, payload map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(payload)
}

func classifyStatus(err error) int {
	var validationErr *validationError
	if errors.Is(err, errInvalidJSON) || errors.As(err, &validationErr) {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

func classifyError(err error) string {
	var validationErr *validationError
	switch {
	case errors.Is(err, errInvalidJSON):
		return errInvalidJSON.Error()
	case errors.As(err, &validationErr):
		return validationErr.Error()
	default:
		return "internal server error"
	}
}

func (e *validationError) Error() string {
	return "validation failed: " + strings.Join(e.Problems, "; ")
}

func validateStruct(input any) error {
	value := reflect.ValueOf(input)
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return &validationError{Problems: []string{"nil input"}}
		}
		value = value.Elem()
	}

	if value.Kind() != reflect.Struct {
		return fmt.Errorf("validateStruct expects struct, got %s", value.Kind())
	}

	valueType := value.Type()
	problems := make([]string, 0)
	for index := range value.NumField() {
		fieldValue := value.Field(index)
		fieldType := valueType.Field(index)
		fieldName := fieldType.Tag.Get("json")
		if fieldName == "" {
			fieldName = strings.ToLower(fieldType.Name)
		}

		if fieldType.Tag.Get("required") == "true" && fieldValue.IsZero() {
			problems = append(problems, fmt.Sprintf("%s is required", fieldName))
		}

		if minValue := fieldType.Tag.Get("min"); minValue != "" {
			minimum, err := strconv.Atoi(minValue)
			if err != nil {
				return fmt.Errorf("invalid min tag on %s: %w", fieldType.Name, err)
			}
			if fieldValue.Kind() == reflect.Int && fieldValue.Int() < int64(minimum) {
				problems = append(problems, fmt.Sprintf("%s must be >= %d", fieldName, minimum))
			}
		}
	}

	if len(problems) == 0 {
		return nil
	}

	return &validationError{Problems: problems}
}

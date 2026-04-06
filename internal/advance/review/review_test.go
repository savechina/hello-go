package review

import (
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateStruct(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		wantError string
	}{
		{
			name:  "valid config passes",
			input: reviewConfig{ServiceName: "review-service", ListenPort: 8080, StorageDSN: "review.db"},
		},
		{
			name:      "missing required fields fail",
			input:     reviewConfig{ListenPort: 8080},
			wantError: "service_name is required",
		},
		{
			name:      "minimum port enforced",
			input:     reviewConfig{ServiceName: "svc", ListenPort: 0, StorageDSN: "review.db"},
			wantError: "listen_port must be >= 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateStruct(tt.input)
			if tt.wantError == "" {
				if err != nil {
					t.Fatalf("validateStruct() unexpected error = %v", err)
				}
				return
			}

			if err == nil || !strings.Contains(err.Error(), tt.wantError) {
				t.Fatalf("validateStruct() error = %v, want substring %q", err, tt.wantError)
			}
		})
	}
}

func TestClassifyStatusAndError(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantText   string
	}{
		{name: "invalid json is bad request", err: errInvalidJSON, wantStatus: http.StatusBadRequest, wantText: "invalid json payload"},
		{name: "validation error is bad request", err: &validationError{Problems: []string{"title is required"}}, wantStatus: http.StatusBadRequest, wantText: "title is required"},
		{name: "unknown error is internal", err: assertError("db down"), wantStatus: http.StatusInternalServerError, wantText: "internal server error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := classifyStatus(tt.err); got != tt.wantStatus {
				t.Fatalf("classifyStatus() = %d, want %d", got, tt.wantStatus)
			}
			if got := classifyError(tt.err); !strings.Contains(got, tt.wantText) {
				t.Fatalf("classifyError() = %q, want substring %q", got, tt.wantText)
			}
		})
	}
}

func TestCourseHandler(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		body        string
		wantStatus  int
		wantText    string
		wantRecords int64
	}{
		{name: "invalid json", method: http.MethodPost, body: `{`, wantStatus: http.StatusBadRequest, wantText: "invalid json payload", wantRecords: 0},
		{name: "missing title", method: http.MethodPost, body: `{"instructor":"gopher"}`, wantStatus: http.StatusBadRequest, wantText: "title is required", wantRecords: 0},
		{name: "success", method: http.MethodPost, body: `{"title":"Review","instructor":"gopher"}`, wantStatus: http.StatusCreated, wantText: "created", wantRecords: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app, err := newReviewApp(reviewConfig{
				ServiceName: "review-service",
				ListenPort:  8080,
				StorageDSN:  filepath.Join(t.TempDir(), "review.db"),
			})
			if err != nil {
				t.Fatalf("newReviewApp() error = %v", err)
			}

			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(tt.method, "/courses", strings.NewReader(tt.body))
			app.courseHandler().ServeHTTP(recorder, request)

			if recorder.Code != tt.wantStatus {
				t.Fatalf("courseHandler() status = %d, want %d", recorder.Code, tt.wantStatus)
			}
			if !strings.Contains(recorder.Body.String(), tt.wantText) {
				t.Fatalf("courseHandler() body = %q, want substring %q", recorder.Body.String(), tt.wantText)
			}

			count, err := app.courseCount()
			if err != nil {
				t.Fatalf("courseCount() error = %v", err)
			}
			if count != tt.wantRecords {
				t.Fatalf("courseCount() = %d, want %d", count, tt.wantRecords)
			}
		})
	}
}

func assertError(message string) error {
	return &staticError{message: message}
}

type staticError struct {
	message string
}

func (e *staticError) Error() string {
	return e.message
}

package errorhandling

import (
	"errors"
	"fmt"
	"strings"

	"hello/internal/chapters"
)

var (
	errConfigMissing = errors.New("config missing")
	errPermission    = errors.New("permission denied")
)

func init() {
	chapters.Register("advance", "errorhandling", Run)
}

type validationError struct {
	Op    string
	Field string
	Value any
	Err   error
}

func (e *validationError) Error() string {
	if e == nil {
		return "<nil>"
	}

	return fmt.Sprintf("%s %s=%v: %v", e.Op, e.Field, e.Value, e.Err)
}

func (e *validationError) Unwrap() error {
	if e == nil {
		return nil
	}

	return e.Err
}

type permissionError struct {
	Resource string
	Err      error
}

func (e *permissionError) Error() string {
	return fmt.Sprintf("access %s: %v", e.Resource, e.Err)
}

func (e *permissionError) Unwrap() error {
	return e.Err
}

func validateUserAge(age int) error {
	if age >= 18 {
		return nil
	}

	return &validationError{
		Op:    "validate user",
		Field: "age",
		Value: age,
		Err:   errors.New("must be at least 18"),
	}
}

func loadRuntimeConfig(name string) error {
	return fmt.Errorf("load runtime config %q: %w", name, errConfigMissing)
}

func authorizeAction(resource string, allowed bool) error {
	if allowed {
		return nil
	}

	return fmt.Errorf("service authorize: %w", &permissionError{Resource: resource, Err: errPermission})
}

func buildTraceChain(job string) error {
	base := errors.New("disk full")
	return fmt.Errorf("run %s: %w", job, fmt.Errorf("flush report: %w", fmt.Errorf("write cache: %w", base)))
}

func summarizeValidation(err error) string {
	var target *validationError
	if errors.As(err, &target) {
		return fmt.Sprintf("validation field=%s value=%v", target.Field, target.Value)
	}

	return "no validation error"
}

func summarizeSentinel(err error) string {
	parts := make([]string, 0, 2)
	if errors.Is(err, errConfigMissing) {
		parts = append(parts, "config-missing")
	}
	if errors.Is(err, errPermission) {
		parts = append(parts, "permission-denied")
	}
	if len(parts) == 0 {
		return "no known sentinel"
	}

	return strings.Join(parts, ",")
}

// Run prints the advanced error handling examples.
func Run() {
	validationErr := validateUserAge(16)
	configErr := loadRuntimeConfig("prod")
	authErr := authorizeAction("reports", false)
	traceErr := buildTraceChain("nightly-job")

	fmt.Println("[errorhandling] example 1:", summarizeValidation(validationErr), "|", validationErr)
	fmt.Println("[errorhandling] example 2:", summarizeSentinel(configErr), "|", summarizeSentinel(authErr))
	fmt.Println("[errorhandling] example 3:", traceErr)
}

package errorhandling

import (
	"errors"
	"fmt"
	"strconv"

	"hello/internal/chapters"
)

var (
	errAmountMustBePositive = errors.New("amount must be positive")
	errSettingNotFound      = errors.New("setting not found")
)

func init() {
	chapters.Register("basic", "errorhandling", Run)
}

func main() {
	Run()
}

// Run executes the error handling chapter examples.
func Run() {
	basicErr := validateAmount(-1)
	_, wrappedErr := lookupSetting(map[string]string{"mode": "dev"}, "timeout")
	_, parseErr := parseRetryCount("abc")

	examples := []string{
		"1) errors.New => " + basicErr.Error(),
		"2) fmt.Errorf + %w => " + wrappedErr.Error(),
		"3) errors.Is / errors.As => " + summarizeError(wrappedErr) + " | " + summarizeError(parseErr),
	}

	for _, example := range examples {
		fmt.Println(example)
	}
}

type fieldError struct {
	field string
	value string
	err   error
}

func (e *fieldError) Error() string {
	return fmt.Sprintf("%s %q: %v", e.field, e.value, e.err)
}

func (e *fieldError) Unwrap() error {
	return e.err
}

func validateAmount(amount int) error {
	if amount <= 0 {
		return errAmountMustBePositive
	}

	return nil
}

func lookupSetting(settings map[string]string, key string) (string, error) {
	value, ok := settings[key]
	if !ok {
		return "", fmt.Errorf("lookup %q: %w", key, errSettingNotFound)
	}

	return value, nil
}

func parseRetryCount(raw string) (int, error) {
	value, err := strconv.Atoi(raw)
	if err != nil {
		return 0, &fieldError{
			field: "retry count",
			value: raw,
			err:   err,
		}
	}

	return value, nil
}

func summarizeError(err error) string {
	if err == nil {
		return "no error"
	}

	if errors.Is(err, errSettingNotFound) {
		return "missing setting detected"
	}

	var fieldErr *fieldError
	if errors.As(err, &fieldErr) {
		return fmt.Sprintf("field error on %s with value %q", fieldErr.field, fieldErr.value)
	}

	return err.Error()
}

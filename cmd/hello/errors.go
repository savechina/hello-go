package main

import (
	"errors"
	"fmt"
	"os"
)

type commandError struct {
	code       int
	message    string
	suggestion string
}

func (e *commandError) Error() string {
	return e.message
}

func newCommandError(code int, msg string, suggestion string) error {
	return &commandError{
		code:       code,
		message:    msg,
		suggestion: suggestion,
	}
}

func handleCommandError(err error) {
	var cmdErr *commandError
	if errors.As(err, &cmdErr) {
		ExitWithError(cmdErr.code, cmdErr.message, cmdErr.suggestion)
	}

	ExitWithError(1, err.Error(), "Run 'hello --help' for usage information.")
}

func ExitWithError(code int, msg string, suggestion string) {
	_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
	if suggestion != "" {
		_, _ = fmt.Fprintf(os.Stderr, "Suggestion: %s\n", suggestion)
	}

	os.Exit(code)
}

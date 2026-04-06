package main

import (
	"bytes"
	"fmt"
	"testing"

	"hello/internal/chapters"
)

func TestLevelCommandsDispatchRegisteredChapters(t *testing.T) {
	tests := []struct {
		level   string
		chapter string
	}{
		{level: "basic", chapter: "test-basic"},
		{level: "advance", chapter: "test-advance"},
		{level: "awesome", chapter: "test-awesome"},
		{level: "algo", chapter: "test-algo"},
		{level: "leetcode", chapter: "test-leetcode"},
		{level: "quiz", chapter: "test-quiz"},
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			var called bool

			chapters.Register(tt.level, tt.chapter, func() {
				called = true
			})

			cmd := newRootCmd()
			cmd.SetArgs([]string{tt.level, tt.chapter})

			var stdout bytes.Buffer
			var stderr bytes.Buffer
			cmd.SetOut(&stdout)
			cmd.SetErr(&stderr)

			if err := cmd.Execute(); err != nil {
				t.Fatalf("execute %s %s: %v", tt.level, tt.chapter, err)
			}

			if !called {
				t.Fatalf("expected runner to be called for %s %s", tt.level, tt.chapter)
			}
		})
	}
}

func TestRootCommandContainsExpectedLevelCommands(t *testing.T) {
	cmd := newRootCmd()

	tests := []string{"basic", "advance", "awesome", "algo", "leetcode", "quiz"}
	for _, level := range tests {
		subCmd, _, err := cmd.Find([]string{level})
		if err != nil {
			t.Fatalf("find %s: %v", level, err)
		}

		if subCmd == nil {
			t.Fatalf("expected subcommand %s to exist", level)
		}

		if subCmd.RunE == nil {
			t.Fatalf("expected subcommand %s to have a RunE handler", level)
		}
	}
}

func TestUnknownChapterReturnsHelpfulError(t *testing.T) {
	cmd := newRootCmd()
	cmd.SetArgs([]string{"basic", "missing-chapter"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for unknown chapter")
	}

	got := err.Error()
	want := "unknown chapter: basic missing-chapter"
	if got != want {
		t.Fatalf("unexpected error message: got %q want %q", got, want)
	}
}

func TestLevelHelpListsRegisteredChapters(t *testing.T) {
	chapterName := fmt.Sprintf("help-basic-%s", t.Name())
	chapters.Register("basic", chapterName, func() {})

	cmd := newRootCmd()
	cmd.SetArgs([]string{"basic", "--help"})

	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetErr(&stdout)

	if err := cmd.Execute(); err != nil {
		t.Fatalf("execute help: %v", err)
	}

	if !bytes.Contains(stdout.Bytes(), []byte(chapterName)) {
		t.Fatalf("expected help output to mention chapter %q, got %q", chapterName, stdout.String())
	}
}

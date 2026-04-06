package errorhandling

import (
	"errors"
	"strings"
	"testing"
)

func TestValidateUserAge(t *testing.T) {
	tests := []struct {
		name    string
		age     int
		wantErr bool
		want    string
	}{
		{name: "adult passes", age: 20, wantErr: false},
		{name: "minor returns validation error", age: 16, wantErr: true, want: "validation field=age value=16"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateUserAge(tt.age)
			if (err != nil) != tt.wantErr {
				t.Fatalf("validateUserAge() err=%v, wantErr=%t", err, tt.wantErr)
			}
			if tt.wantErr {
				if got := summarizeValidation(err); got != tt.want {
					t.Fatalf("summarizeValidation() = %q, want %q", got, tt.want)
				}
			}
		})
	}
}

func TestSentinelMatching(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{name: "config missing chain", err: loadRuntimeConfig("prod"), want: "config-missing"},
		{name: "permission chain", err: authorizeAction("reports", false), want: "permission-denied"},
		{name: "clean path", err: nil, want: "no known sentinel"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := summarizeSentinel(tt.err); got != tt.want {
				t.Fatalf("summarizeSentinel() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTraceAndAs(t *testing.T) {
	tests := []struct {
		name         string
		err          error
		wantContains string
		wantIs       error
	}{
		{name: "trace preserves message", err: buildTraceChain("job"), wantContains: "write cache", wantIs: nil},
		{name: "permission error extracts sentinel", err: authorizeAction("dashboard", false), wantContains: "dashboard", wantIs: errPermission},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantContains != "" && !strings.Contains(tt.err.Error(), tt.wantContains) {
				t.Fatalf("error %q missing %q", tt.err.Error(), tt.wantContains)
			}
			if tt.wantIs != nil && !errors.Is(tt.err, tt.wantIs) {
				t.Fatalf("errors.Is(%v, %v) = false, want true", tt.err, tt.wantIs)
			}
		})
	}
}

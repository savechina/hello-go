package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadEnvConfig(t *testing.T) {
	tests := []struct {
		name      string
		env       map[string]string
		wantApp   string
		wantPort  int
		wantDSN   string
		wantError string
	}{
		{
			name: "override defaults from environment",
			env: map[string]string{
				"HELLO_APP_NAME":     "chapter-env",
				"HELLO_SERVER_PORT":  "9099",
				"HELLO_DATABASE_DSN": "file:env-test.db",
			},
			wantApp:  "chapter-env",
			wantPort: 9099,
			wantDSN:  "file:env-test.db",
		},
		{
			name: "invalid integer returns error",
			env: map[string]string{
				"HELLO_SERVER_PORT": "oops",
			},
			wantError: "HELLO_SERVER_PORT",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := loadEnvConfig(defaultConfig(), "HELLO", mapLookup(tt.env))
			if tt.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantError) {
					t.Fatalf("loadEnvConfig() error = %v, want substring %q", err, tt.wantError)
				}
				return
			}

			if err != nil {
				t.Fatalf("loadEnvConfig() unexpected error = %v", err)
			}

			if cfg.AppName != tt.wantApp || cfg.Server.Port != tt.wantPort || cfg.Database.DSN != tt.wantDSN {
				t.Fatalf("loadEnvConfig() = %+v, want app=%q port=%d dsn=%q", cfg, tt.wantApp, tt.wantPort, tt.wantDSN)
			}
		})
	}
}

func TestLoadConfigFile(t *testing.T) {
	tests := []struct {
		name      string
		fileName  string
		content   string
		wantHost  string
		wantPort  int
		wantPool  int
		wantError string
	}{
		{
			name:     "load json config",
			fileName: "app.json",
			content:  `{"server":{"host":"0.0.0.0","port":8088},"database":{"max_open_conns":6}}`,
			wantHost: "0.0.0.0",
			wantPort: 8088,
			wantPool: 6,
		},
		{
			name:     "load yaml config",
			fileName: "app.yaml",
			content:  "server:\n  host: yaml.local\n  port: 8082\ndatabase:\n  max_open_conns: 9\n",
			wantHost: "yaml.local",
			wantPort: 8082,
			wantPool: 9,
		},
		{
			name:      "reject unsupported extension",
			fileName:  "app.toml",
			content:   "port = 80",
			wantError: "unsupported config extension",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), tt.fileName)
			if err := os.WriteFile(path, []byte(tt.content), 0o600); err != nil {
				t.Fatalf("WriteFile() error = %v", err)
			}

			cfg, err := loadConfigFile(defaultConfig(), path)
			if tt.wantError != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantError) {
					t.Fatalf("loadConfigFile() error = %v, want substring %q", err, tt.wantError)
				}
				return
			}

			if err != nil {
				t.Fatalf("loadConfigFile() unexpected error = %v", err)
			}

			if cfg.Server.Host != tt.wantHost || cfg.Server.Port != tt.wantPort || cfg.Database.MaxOpenConns != tt.wantPool {
				t.Fatalf("loadConfigFile() = %+v, want host=%q port=%d pool=%d", cfg, tt.wantHost, tt.wantPort, tt.wantPool)
			}
		})
	}
}

func TestResolveConfig(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		env      map[string]string
		wantHost string
		wantLog  string
		wantPool int
	}{
		{
			name:    "environment wins over file",
			content: "log_level: info\nserver:\n  host: yaml.local\ndatabase:\n  max_open_conns: 4\n",
			env: map[string]string{
				"HELLO_LOG_LEVEL":               "debug",
				"HELLO_SERVER_HOST":             "env.local",
				"HELLO_DATABASE_MAX_OPEN_CONNS": "12",
			},
			wantHost: "env.local",
			wantLog:  "debug",
			wantPool: 12,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "config.yaml")
			if err := os.WriteFile(path, []byte(tt.content), 0o600); err != nil {
				t.Fatalf("WriteFile() error = %v", err)
			}

			cfg, err := resolveConfig([]string{path}, "HELLO", mapLookup(tt.env))
			if err != nil {
				t.Fatalf("resolveConfig() unexpected error = %v", err)
			}

			if cfg.Server.Host != tt.wantHost || cfg.LogLevel != tt.wantLog || cfg.Database.MaxOpenConns != tt.wantPool {
				t.Fatalf("resolveConfig() = %+v, want host=%q log=%q pool=%d", cfg, tt.wantHost, tt.wantLog, tt.wantPool)
			}
		})
	}
}

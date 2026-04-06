package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"hello/internal/chapters"
)

type appConfig struct {
	AppName  string         `json:"app_name" config:"app_name" env:"APP_NAME"`
	LogLevel string         `json:"log_level" config:"log_level" env:"LOG_LEVEL"`
	Server   serverConfig   `json:"server" config:"server"`
	Database databaseConfig `json:"database" config:"database"`
}

type serverConfig struct {
	Host string `json:"host" config:"host" env:"SERVER_HOST"`
	Port int    `json:"port" config:"port" env:"SERVER_PORT"`
}

type databaseConfig struct {
	Driver       string `json:"driver" config:"driver" env:"DATABASE_DRIVER"`
	DSN          string `json:"dsn" config:"dsn" env:"DATABASE_DSN"`
	MaxOpenConns int    `json:"max_open_conns" config:"max_open_conns" env:"DATABASE_MAX_OPEN_CONNS"`
}

func init() {
	chapters.Register("advance", "config", Run)
}

// Run prints runnable examples for configuration management patterns.
func Run() {
	for _, line := range renderExamples() {
		fmt.Println(line)
	}
}

func renderExamples() []string {
	return []string{
		exampleEnvironmentOverride(),
		exampleJSONConfigFile(),
		exampleLayeredConfigSources(),
	}
}

func exampleEnvironmentOverride() string {
	lookup := mapLookup(map[string]string{
		"HELLO_APP_NAME":     "hello-go-env",
		"HELLO_SERVER_PORT":  "9090",
		"HELLO_DATABASE_DSN": "file:env.db",
	})

	cfg, err := loadEnvConfig(defaultConfig(), "HELLO", lookup)
	if err != nil {
		return fmt.Sprintf("[config] example 1 error: %v", err)
	}

	return fmt.Sprintf("[config] example 1 env override => %s", summarizeConfig(cfg))
}

func exampleJSONConfigFile() string {
	file, err := os.CreateTemp("", "hello-config-*.json")
	if err != nil {
		return fmt.Sprintf("[config] example 2 error: %v", err)
	}
	defer os.Remove(file.Name())

	content := `{
		"app_name": "hello-go-json",
		"log_level": "debug",
		"server": {"host": "0.0.0.0", "port": 8088},
		"database": {"driver": "sqlite", "dsn": "file:data/json.db", "max_open_conns": 8}
	}`

	if _, err := file.WriteString(content); err != nil {
		return fmt.Sprintf("[config] example 2 error: %v", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Sprintf("[config] example 2 error: %v", err)
	}

	cfg, err := loadConfigFile(defaultConfig(), file.Name())
	if err != nil {
		return fmt.Sprintf("[config] example 2 error: %v", err)
	}

	return fmt.Sprintf("[config] example 2 json file => %s", summarizeConfig(cfg))
}

func exampleLayeredConfigSources() string {
	file, err := os.CreateTemp("", "hello-config-*.yaml")
	if err != nil {
		return fmt.Sprintf("[config] example 3 error: %v", err)
	}
	defer os.Remove(file.Name())

	content := `app_name: hello-go-yaml
log_level: info
server:
  host: 127.0.0.1
  port: 8081
database:
  driver: sqlite
  dsn: file:data/yaml.db
  max_open_conns: 4
`

	if _, err := file.WriteString(content); err != nil {
		return fmt.Sprintf("[config] example 3 error: %v", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Sprintf("[config] example 3 error: %v", err)
	}

	lookup := mapLookup(map[string]string{
		"HELLO_LOG_LEVEL":               "warn",
		"HELLO_DATABASE_MAX_OPEN_CONNS": "16",
		"HELLO_SERVER_HOST":             "config.service.local",
	})

	cfg, err := resolveConfig([]string{file.Name()}, "HELLO", lookup)
	if err != nil {
		return fmt.Sprintf("[config] example 3 error: %v", err)
	}

	return fmt.Sprintf("[config] example 3 layered sources => %s", summarizeConfig(cfg))
}

func defaultConfig() appConfig {
	return appConfig{
		AppName:  "hello-go",
		LogLevel: "info",
		Server: serverConfig{
			Host: "127.0.0.1",
			Port: 8080,
		},
		Database: databaseConfig{
			Driver:       "sqlite",
			DSN:          "file:data/default.db",
			MaxOpenConns: 2,
		},
	}
}

func resolveConfig(paths []string, prefix string, lookup func(string) (string, bool)) (appConfig, error) {
	cfg := defaultConfig()

	for _, path := range paths {
		next, err := loadConfigFile(cfg, path)
		if err != nil {
			return appConfig{}, err
		}
		cfg = next
	}

	return loadEnvConfig(cfg, prefix, lookup)
}

func loadEnvConfig(base appConfig, prefix string, lookup func(string) (string, bool)) (appConfig, error) {
	cfg := base
	if err := bindEnv(&cfg, prefix, lookup); err != nil {
		return appConfig{}, err
	}
	return cfg, nil
}

func loadConfigFile(base appConfig, path string) (appConfig, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return appConfig{}, fmt.Errorf("read config file %q: %w", path, err)
	}

	values := map[string]any{}
	switch strings.ToLower(filepath.Ext(path)) {
	case ".json":
		if err := json.Unmarshal(content, &values); err != nil {
			return appConfig{}, fmt.Errorf("decode json config %q: %w", path, err)
		}
	case ".yaml", ".yml":
		values, err = parseSimpleYAML(string(content))
		if err != nil {
			return appConfig{}, fmt.Errorf("decode yaml config %q: %w", path, err)
		}
	default:
		return appConfig{}, fmt.Errorf("unsupported config extension: %s", filepath.Ext(path))
	}

	cfg := base
	if err := bindMap(&cfg, values); err != nil {
		return appConfig{}, err
	}

	return cfg, nil
}

func bindMap(target any, values map[string]any) error {
	root := reflect.ValueOf(target)
	if root.Kind() != reflect.Pointer || root.IsNil() {
		return errors.New("config target must be a non-nil pointer")
	}

	return bindMapValue(root.Elem(), values)
}

func bindMapValue(target reflect.Value, values map[string]any) error {
	targetType := target.Type()
	for index := range target.NumField() {
		fieldValue := target.Field(index)
		fieldType := targetType.Field(index)
		key := fieldType.Tag.Get("config")
		if key == "" {
			continue
		}

		rawValue, ok := values[key]
		if !ok {
			continue
		}

		if fieldValue.Kind() == reflect.Struct {
			nestedValues, ok := rawValue.(map[string]any)
			if !ok {
				return fmt.Errorf("config key %q must be an object", key)
			}
			if err := bindMapValue(fieldValue, nestedValues); err != nil {
				return err
			}
			continue
		}

		if err := setValueFromAny(fieldValue, rawValue); err != nil {
			return fmt.Errorf("config key %q: %w", key, err)
		}
	}

	return nil
}

func bindEnv(target any, prefix string, lookup func(string) (string, bool)) error {
	root := reflect.ValueOf(target)
	if root.Kind() != reflect.Pointer || root.IsNil() {
		return errors.New("config target must be a non-nil pointer")
	}

	return bindEnvValue(root.Elem(), prefix, lookup)
}

func bindEnvValue(target reflect.Value, prefix string, lookup func(string) (string, bool)) error {
	targetType := target.Type()
	for index := range target.NumField() {
		fieldValue := target.Field(index)
		fieldType := targetType.Field(index)

		if fieldValue.Kind() == reflect.Struct {
			if err := bindEnvValue(fieldValue, prefix, lookup); err != nil {
				return err
			}
			continue
		}

		suffix := fieldType.Tag.Get("env")
		if suffix == "" {
			continue
		}

		key := suffix
		if prefix != "" {
			key = prefix + "_" + suffix
		}

		rawValue, ok := lookup(key)
		if !ok {
			continue
		}

		if err := setValueFromString(fieldValue, rawValue); err != nil {
			return fmt.Errorf("environment key %q: %w", key, err)
		}
	}

	return nil
}

func setValueFromAny(target reflect.Value, rawValue any) error {
	switch value := rawValue.(type) {
	case string:
		return setValueFromString(target, value)
	case bool:
		if target.Kind() != reflect.Bool {
			return fmt.Errorf("expected %s, got bool", target.Kind())
		}
		target.SetBool(value)
		return nil
	case float64:
		if target.Kind() != reflect.Int {
			return fmt.Errorf("expected %s, got number", target.Kind())
		}
		target.SetInt(int64(value))
		return nil
	case int:
		if target.Kind() != reflect.Int {
			return fmt.Errorf("expected %s, got int", target.Kind())
		}
		target.SetInt(int64(value))
		return nil
	default:
		return fmt.Errorf("unsupported value type %T", rawValue)
	}
}

func setValueFromString(target reflect.Value, rawValue string) error {
	switch target.Kind() {
	case reflect.String:
		target.SetString(rawValue)
		return nil
	case reflect.Int:
		parsed, err := strconv.Atoi(rawValue)
		if err != nil {
			return err
		}
		target.SetInt(int64(parsed))
		return nil
	case reflect.Bool:
		parsed, err := strconv.ParseBool(rawValue)
		if err != nil {
			return err
		}
		target.SetBool(parsed)
		return nil
	default:
		return fmt.Errorf("unsupported target kind %s", target.Kind())
	}
}

func parseSimpleYAML(content string) (map[string]any, error) {
	result := map[string]any{}
	currentSection := ""

	for lineNumber, rawLine := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(rawLine)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		parts := strings.SplitN(trimmed, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: invalid yaml entry", lineNumber+1)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		indent := len(rawLine) - len(strings.TrimLeft(rawLine, " "))

		if indent == 0 {
			if value == "" {
				result[key] = map[string]any{}
				currentSection = key
				continue
			}

			result[key] = parseScalar(value)
			currentSection = ""
			continue
		}

		if indent < 2 || currentSection == "" {
			return nil, fmt.Errorf("line %d: unsupported yaml indentation", lineNumber+1)
		}

		sectionValues, ok := result[currentSection].(map[string]any)
		if !ok {
			return nil, fmt.Errorf("line %d: invalid yaml section %q", lineNumber+1, currentSection)
		}

		sectionValues[key] = parseScalar(value)
	}

	return result, nil
}

func parseScalar(value string) any {
	unquoted := strings.Trim(value, `"'`)
	if parsed, err := strconv.Atoi(unquoted); err == nil {
		return parsed
	}
	if parsed, err := strconv.ParseBool(unquoted); err == nil {
		return parsed
	}
	return unquoted
}

func mapLookup(values map[string]string) func(string) (string, bool) {
	return func(key string) (string, bool) {
		value, ok := values[key]
		return value, ok
	}
}

func summarizeConfig(cfg appConfig) string {
	return fmt.Sprintf(
		"app=%s log=%s server=%s:%d db=%s(%s) pool=%d",
		cfg.AppName,
		cfg.LogLevel,
		cfg.Server.Host,
		cfg.Server.Port,
		cfg.Database.Driver,
		cfg.Database.DSN,
		cfg.Database.MaxOpenConns,
	)
}

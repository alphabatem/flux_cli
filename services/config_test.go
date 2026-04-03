package services

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/alphabatem/flux_cli/dto"
)

func newTestConfigService(t *testing.T) *ConfigService {
	t.Helper()
	dir := t.TempDir()
	svc := &ConfigService{
		configDir:  dir,
		configPath: filepath.Join(dir, "config.json"),
		cfg:        dto.DefaultConfig(),
	}
	return svc
}

func TestConfigService_Id(t *testing.T) {
	svc := &ConfigService{}
	if svc.Id() != CONFIG_SVC {
		t.Errorf("expected %s, got %s", CONFIG_SVC, svc.Id())
	}
}

func TestConfigService_Set_ValidKeys(t *testing.T) {
	tests := []struct {
		key   string
		value string
		check func(cfg *dto.Config) string
	}{
		{"datastream.api_key", "ds-key-123", func(c *dto.Config) string { return c.DataStream.APIKey }},
		{"datastream.base_url", "https://custom.url", func(c *dto.Config) string { return c.DataStream.BaseURL }},
		{"fluxrpc.api_key", "rpc-key-456", func(c *dto.Config) string { return c.FluxRPC.APIKey }},
		{"fluxrpc.base_url", "https://rpc.custom", func(c *dto.Config) string { return c.FluxRPC.BaseURL }},
		{"fluxrpc.region", "eu", func(c *dto.Config) string { return c.FluxRPC.Region }},
		{"rugcheck.api_key", "rc-key-789", func(c *dto.Config) string { return c.RugCheck.APIKey }},
		{"rugcheck.base_url", "https://rc.custom", func(c *dto.Config) string { return c.RugCheck.BaseURL }},
		{"output.format", "table", func(c *dto.Config) string { return c.Output.Format }},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			svc := newTestConfigService(t)
			err := svc.Set(tt.key, tt.value)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got := tt.check(&svc.cfg); got != tt.value {
				t.Errorf("expected %q, got %q", tt.value, got)
			}
		})
	}
}

func TestConfigService_Set_InvalidKeys(t *testing.T) {
	tests := []struct {
		name string
		key  string
	}{
		{"no dot", "nodot"},
		{"unknown section", "unknown.api_key"},
		{"unknown field", "datastream.unknown"},
		{"unknown fluxrpc field", "fluxrpc.unknown"},
		{"unknown rugcheck field", "rugcheck.unknown"},
		{"unknown output field", "output.unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newTestConfigService(t)
			err := svc.Set(tt.key, "value")
			if err == nil {
				t.Error("expected error for invalid key")
			}
		})
	}
}

func TestConfigService_Set_InvalidRegion(t *testing.T) {
	svc := newTestConfigService(t)
	err := svc.Set("fluxrpc.region", "asia")
	if err == nil {
		t.Error("expected error for invalid region")
	}
}

func TestConfigService_Get_ValidKeys(t *testing.T) {
	svc := newTestConfigService(t)
	svc.cfg.DataStream.APIKey = "test-key"
	svc.cfg.FluxRPC.BaseURL = "https://rpc.test"

	val, err := svc.Get("datastream.api_key")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "test-key" {
		t.Errorf("expected test-key, got %s", val)
	}

	val, err = svc.Get("fluxrpc.base_url")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "https://rpc.test" {
		t.Errorf("expected https://rpc.test, got %s", val)
	}
}

func TestConfigService_Get_InvalidKey(t *testing.T) {
	svc := newTestConfigService(t)
	_, err := svc.Get("invalid")
	if err == nil {
		t.Error("expected error for invalid key format")
	}

	_, err = svc.Get("unknown.field")
	if err == nil {
		t.Error("expected error for unknown key")
	}
}

func TestConfigService_List(t *testing.T) {
	svc := newTestConfigService(t)
	svc.cfg.DataStream.APIKey = "long-api-key-here"
	svc.cfg.FluxRPC.APIKey = ""

	list := svc.List()

	// Check redacted key
	dsKey, ok := list["datastream.api_key"].(string)
	if !ok {
		t.Fatal("expected string for datastream.api_key")
	}
	if dsKey == "long-api-key-here" {
		t.Error("API key should be redacted")
	}
	if dsKey != "long****here" {
		t.Errorf("expected long****here, got %s", dsKey)
	}

	// Check unset key
	rpcKey, ok := list["fluxrpc.api_key"].(string)
	if !ok {
		t.Fatal("expected string for fluxrpc.api_key")
	}
	if rpcKey != "(not set)" {
		t.Errorf("expected (not set), got %s", rpcKey)
	}

	// Check non-sensitive field not redacted
	baseURL, ok := list["datastream.base_url"].(string)
	if !ok {
		t.Fatal("expected string for datastream.base_url")
	}
	if baseURL != "https://data.fluxbeam.xyz" {
		t.Errorf("expected default base URL, got %s", baseURL)
	}
}

func TestConfigService_Set_PersistsToFile(t *testing.T) {
	svc := newTestConfigService(t)
	err := svc.Set("datastream.api_key", "persist-test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Read the file back
	data, err := os.ReadFile(svc.configPath)
	if err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}

	var cfg dto.Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("failed to unmarshal config: %v", err)
	}
	if cfg.DataStream.APIKey != "persist-test" {
		t.Errorf("expected persist-test in file, got %s", cfg.DataStream.APIKey)
	}
}

func TestConfigService_ConfigPath(t *testing.T) {
	svc := newTestConfigService(t)
	if svc.ConfigPath() == "" {
		t.Error("expected non-empty config path")
	}
}

func TestConfigService_Exists(t *testing.T) {
	svc := newTestConfigService(t)
	if svc.Exists() {
		t.Fatal("expected config to not exist yet")
	}

	if err := svc.Set("fluxrpc.region", "eu"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !svc.Exists() {
		t.Fatal("expected config to exist after save")
	}
}

func TestConfigService_SaveConfig(t *testing.T) {
	svc := newTestConfigService(t)
	cfg := dto.DefaultConfig()
	cfg.FluxRPC.Region = "eu"
	cfg.Output.Format = "table"
	cfg.RugCheck.APIKey = "rc-key"

	if err := svc.SaveConfig(cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	data, err := os.ReadFile(svc.configPath)
	if err != nil {
		t.Fatalf("failed to read config file: %v", err)
	}

	var saved dto.Config
	if err := json.Unmarshal(data, &saved); err != nil {
		t.Fatalf("failed to unmarshal config: %v", err)
	}
	if saved.FluxRPC.Region != "eu" {
		t.Fatalf("expected region eu, got %s", saved.FluxRPC.Region)
	}
	if saved.Output.Format != "table" {
		t.Fatalf("expected output format table, got %s", saved.Output.Format)
	}
}

func TestRedact(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", "(not set)"},
		{"abc", "****"},
		{"12345678", "****"},
		{"123456789", "1234****6789"},
		{"abcdefghijklmnop", "abcd****mnop"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := redact(tt.input)
			if got != tt.expected {
				t.Errorf("redact(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestLoadDotEnv(t *testing.T) {
	// Create a temp .env file
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	content := `# comment
FLUX_TEST_VAR=hello
FLUX_TEST_QUOTED="world"
FLUX_TEST_SINGLE='single'
  FLUX_TEST_SPACES = spaced

`
	os.WriteFile(envPath, []byte(content), 0600)

	// Clear any existing values
	os.Unsetenv("FLUX_TEST_VAR")
	os.Unsetenv("FLUX_TEST_QUOTED")
	os.Unsetenv("FLUX_TEST_SINGLE")
	os.Unsetenv("FLUX_TEST_SPACES")
	defer func() {
		os.Unsetenv("FLUX_TEST_VAR")
		os.Unsetenv("FLUX_TEST_QUOTED")
		os.Unsetenv("FLUX_TEST_SINGLE")
		os.Unsetenv("FLUX_TEST_SPACES")
	}()

	loadDotEnv(envPath)

	if v := os.Getenv("FLUX_TEST_VAR"); v != "hello" {
		t.Errorf("expected hello, got %q", v)
	}
	if v := os.Getenv("FLUX_TEST_QUOTED"); v != "world" {
		t.Errorf("expected world (unquoted), got %q", v)
	}
	if v := os.Getenv("FLUX_TEST_SINGLE"); v != "single" {
		t.Errorf("expected single (unquoted), got %q", v)
	}
	if v := os.Getenv("FLUX_TEST_SPACES"); v != "spaced" {
		t.Errorf("expected spaced (trimmed), got %q", v)
	}
}

func TestLoadDotEnv_NoOverride(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")
	os.WriteFile(envPath, []byte("FLUX_TEST_EXISTING=from-file\n"), 0600)

	os.Setenv("FLUX_TEST_EXISTING", "from-env")
	defer os.Unsetenv("FLUX_TEST_EXISTING")

	loadDotEnv(envPath)

	if v := os.Getenv("FLUX_TEST_EXISTING"); v != "from-env" {
		t.Errorf("expected from-env (should not override), got %q", v)
	}
}

func TestLoadDotEnv_MissingFile(t *testing.T) {
	// Should not panic or error
	loadDotEnv("/nonexistent/.env")
}

func TestConfigService_EnvOverride(t *testing.T) {
	svc := newTestConfigService(t)

	// Write a config file with a known key
	svc.cfg.DataStream.APIKey = "from-file"

	// Set env var (higher priority)
	os.Setenv("FLUX_DATASTREAM_API_KEY", "from-env")
	defer os.Unsetenv("FLUX_DATASTREAM_API_KEY")

	// Simulate the env overlay from Configure
	if v := os.Getenv("FLUX_DATASTREAM_API_KEY"); v != "" {
		svc.cfg.DataStream.APIKey = v
	}

	val, _ := svc.Get("datastream.api_key")
	if val != "from-env" {
		t.Errorf("expected from-env, got %s", val)
	}
}

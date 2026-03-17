package services

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alphabatem/common/context"
	"github.com/alphabatem/flux_cli/dto"
)

const CONFIG_SVC = "config_svc"

type ConfigService struct {
	context.DefaultService

	configDir  string
	configPath string
	cfg        dto.Config
}

func (s *ConfigService) Id() string {
	return CONFIG_SVC
}

func (s *ConfigService) Configure(ctx *context.Context) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("unable to determine home directory: %w", err)
	}

	s.configDir = filepath.Join(home, ".flux-cli")
	s.configPath = filepath.Join(s.configDir, "config.json")

	// Load existing config or start with defaults
	s.cfg = dto.DefaultConfig()
	if data, err := os.ReadFile(s.configPath); err == nil {
		if err := json.Unmarshal(data, &s.cfg); err != nil {
			return fmt.Errorf("parsing config file: %w", err)
		}
	}

	// Load .env file if present (does not override existing env vars)
	loadDotEnv(".env")

	// Overlay environment variables (higher priority than file)
	if v := os.Getenv("FLUX_DATASTREAM_API_KEY"); v != "" {
		s.cfg.DataStream.APIKey = v
	}
	if v := os.Getenv("FLUX_RPC_API_KEY"); v != "" {
		s.cfg.FluxRPC.APIKey = v
	}
	if v := os.Getenv("FLUX_RPC_REGION"); v != "" {
		s.cfg.FluxRPC.Region = v
	}
	if v := os.Getenv("FLUX_RUGCHECK_API_KEY"); v != "" {
		s.cfg.RugCheck.APIKey = v
	}

	return s.DefaultService.Configure(ctx)
}

func (s *ConfigService) Config() *dto.Config {
	return &s.cfg
}

func (s *ConfigService) ConfigPath() string {
	return s.configPath
}

// Set updates a config value by dot-notation key (e.g. "datastream.api_key").
func (s *ConfigService) Set(key, value string) error {
	parts := strings.SplitN(key, ".", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid key format: use <section>.<field> (e.g. datastream.api_key)")
	}

	section, field := parts[0], parts[1]

	switch section {
	case "datastream":
		switch field {
		case "api_key":
			s.cfg.DataStream.APIKey = value
		case "base_url":
			s.cfg.DataStream.BaseURL = value
		default:
			return fmt.Errorf("unknown field: %s.%s", section, field)
		}
	case "fluxrpc":
		switch field {
		case "api_key":
			s.cfg.FluxRPC.APIKey = value
		case "base_url":
			s.cfg.FluxRPC.BaseURL = value
		case "region":
			if value != "eu" && value != "us" {
				return fmt.Errorf("invalid region: %s (must be 'eu' or 'us')", value)
			}
			s.cfg.FluxRPC.Region = value
		default:
			return fmt.Errorf("unknown field: %s.%s", section, field)
		}
	case "rugcheck":
		switch field {
		case "api_key":
			s.cfg.RugCheck.APIKey = value
		case "base_url":
			s.cfg.RugCheck.BaseURL = value
		default:
			return fmt.Errorf("unknown field: %s.%s", section, field)
		}
	case "output":
		switch field {
		case "format":
			s.cfg.Output.Format = value
		default:
			return fmt.Errorf("unknown field: %s.%s", section, field)
		}
	default:
		return fmt.Errorf("unknown section: %s", section)
	}

	return s.save()
}

// Get retrieves a config value by dot-notation key.
func (s *ConfigService) Get(key string) (string, error) {
	parts := strings.SplitN(key, ".", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid key format: use <section>.<field> (e.g. datastream.api_key)")
	}

	section, field := parts[0], parts[1]

	switch section {
	case "datastream":
		switch field {
		case "api_key":
			return s.cfg.DataStream.APIKey, nil
		case "base_url":
			return s.cfg.DataStream.BaseURL, nil
		}
	case "fluxrpc":
		switch field {
		case "api_key":
			return s.cfg.FluxRPC.APIKey, nil
		case "base_url":
			return s.cfg.FluxRPC.BaseURL, nil
		case "region":
			return s.cfg.FluxRPC.Region, nil
		}
	case "rugcheck":
		switch field {
		case "api_key":
			return s.cfg.RugCheck.APIKey, nil
		case "base_url":
			return s.cfg.RugCheck.BaseURL, nil
		}
	case "output":
		switch field {
		case "format":
			return s.cfg.Output.Format, nil
		}
	}

	return "", fmt.Errorf("unknown key: %s", key)
}

// List returns a redacted copy of the config suitable for display.
func (s *ConfigService) List() map[string]interface{} {
	return map[string]interface{}{
		"datastream.api_key":  redact(s.cfg.DataStream.APIKey),
		"datastream.base_url": s.cfg.DataStream.BaseURL,
		"fluxrpc.api_key":     redact(s.cfg.FluxRPC.APIKey),
		"fluxrpc.base_url":    ResolveRPCBaseURL(&s.cfg.FluxRPC),
		"fluxrpc.region":      s.cfg.FluxRPC.Region,
		"rugcheck.api_key":    redact(s.cfg.RugCheck.APIKey),
		"rugcheck.base_url":   s.cfg.RugCheck.BaseURL,
		"output.format":       s.cfg.Output.Format,
	}
}

func (s *ConfigService) save() error {
	if err := os.MkdirAll(s.configDir, 0700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}

	data, err := json.MarshalIndent(s.cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding config: %w", err)
	}

	return os.WriteFile(s.configPath, data, 0600)
}

// loadDotEnv reads a .env file and sets env vars that are not already set.
func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)
		// Strip surrounding quotes
		if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'')) {
			value = value[1 : len(value)-1]
		}
		// Don't override existing env vars
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
}

func redact(s string) string {
	if s == "" {
		return "(not set)"
	}
	if len(s) <= 8 {
		return "****"
	}
	return s[:4] + "****" + s[len(s)-4:]
}

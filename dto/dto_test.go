package dto

import (
	"encoding/json"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.DataStream.BaseURL != "https://data.fluxbeam.xyz" {
		t.Errorf("expected default datastream URL, got %s", cfg.DataStream.BaseURL)
	}
	if cfg.FluxRPC.Region != "us" {
		t.Errorf("expected default fluxrpc region us, got %s", cfg.FluxRPC.Region)
	}
	if cfg.RugCheck.BaseURL != "https://api.rugcheck.xyz" {
		t.Errorf("expected default rugcheck URL, got %s", cfg.RugCheck.BaseURL)
	}
	if cfg.Output.Format != "json" {
		t.Errorf("expected default format json, got %s", cfg.Output.Format)
	}
	if cfg.DataStream.APIKey != "" {
		t.Error("expected empty API key by default")
	}
}

func TestCLIResponse_MarshalJSON(t *testing.T) {
	resp := CLIResponse{
		Success: true,
		Data:    map[string]string{"key": "value"},
		Meta:    &CLIMeta{Service: "test", Endpoint: "/test", DurationMs: 100},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var parsed map[string]interface{}
	json.Unmarshal(data, &parsed)

	if parsed["success"] != true {
		t.Error("expected success=true")
	}
	if parsed["error"] != nil {
		t.Error("expected error=null")
	}

	meta := parsed["meta"].(map[string]interface{})
	if meta["service"] != "test" {
		t.Errorf("expected service=test, got %v", meta["service"])
	}
}

func TestCLIResponse_ErrorMarshal(t *testing.T) {
	resp := CLIResponse{
		Success: false,
		Error:   &CLIError{Code: "ERR_CODE", Message: "something broke"},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var parsed map[string]interface{}
	json.Unmarshal(data, &parsed)

	if parsed["success"] != false {
		t.Error("expected success=false")
	}
	if parsed["data"] != nil {
		t.Error("expected data=null")
	}

	errObj := parsed["error"].(map[string]interface{})
	if errObj["code"] != "ERR_CODE" {
		t.Errorf("expected ERR_CODE, got %v", errObj["code"])
	}
}

func TestCLIMeta_OmitEmpty(t *testing.T) {
	resp := CLIResponse{
		Success: true,
		Data:    "ok",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Meta should be omitted when nil
	var parsed map[string]interface{}
	json.Unmarshal(data, &parsed)
	if _, exists := parsed["meta"]; exists {
		t.Error("expected meta to be omitted when nil")
	}
}

func TestConfig_MarshalRoundtrip(t *testing.T) {
	cfg := DefaultConfig()
	cfg.DataStream.APIKey = "ds-key"
	cfg.FluxRPC.APIKey = "rpc-key"
	cfg.RugCheck.APIKey = "rc-key"

	data, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var parsed Config
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if parsed.DataStream.APIKey != "ds-key" {
		t.Errorf("expected ds-key, got %s", parsed.DataStream.APIKey)
	}
	if parsed.FluxRPC.APIKey != "rpc-key" {
		t.Errorf("expected rpc-key, got %s", parsed.FluxRPC.APIKey)
	}
	if parsed.RugCheck.APIKey != "rc-key" {
		t.Errorf("expected rc-key, got %s", parsed.RugCheck.APIKey)
	}
}

func TestExitCodes(t *testing.T) {
	// Verify exit code constants are distinct
	codes := map[int]string{
		ExitSuccess:      "Success",
		ExitGeneralError: "GeneralError",
		ExitUsageError:   "UsageError",
		ExitAPIError:     "APIError",
		ExitAuthError:    "AuthError",
		ExitServiceDown:  "ServiceDown",
		ExitConfigError:  "ConfigError",
	}

	if len(codes) != 7 {
		t.Errorf("expected 7 distinct exit codes, got %d", len(codes))
	}

	if ExitSuccess != 0 {
		t.Error("ExitSuccess should be 0")
	}
}

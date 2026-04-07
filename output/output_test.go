package output

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/alphabatem/flux_cli/dto"
	"github.com/spf13/cobra"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func captureStderr(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	fn()

	w.Close()
	os.Stderr = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func newTestCmd(format string) *cobra.Command {
	cmd := &cobra.Command{Use: "test"}
	cmd.Flags().String("format", format, "")
	return cmd
}

func TestPrintJSON_Success(t *testing.T) {
	cmd := newTestCmd("json")
	out := captureStdout(t, func() {
		PrintSuccess(cmd, map[string]string{"key": "value"}, &dto.CLIMeta{Service: "test", Endpoint: "/test"})
	})

	var resp dto.CLIResponse
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("failed to parse JSON output: %v\nOutput: %s", err, out)
	}
	if !resp.Success {
		t.Error("expected success=true")
	}
	if resp.Error != nil {
		t.Error("expected error=nil")
	}
	data, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("expected map data, got %T", resp.Data)
	}
	if data["key"] != "value" {
		t.Errorf("expected key=value, got %v", data["key"])
	}
}

func TestPrintJSON_Error(t *testing.T) {
	cmd := newTestCmd("json")
	out := captureStdout(t, func() {
		PrintError(cmd, "API_KEY_MISSING", "key not set", nil)
	})

	var resp dto.CLIResponse
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	if resp.Success {
		t.Error("expected success=false")
	}
	if resp.Error == nil {
		t.Fatal("expected error to be set")
	}
	if resp.Error.Code != "API_KEY_MISSING" {
		t.Errorf("expected API_KEY_MISSING, got %s", resp.Error.Code)
	}
	if resp.Error.Message != "key not set" {
		t.Errorf("expected 'key not set', got %s", resp.Error.Message)
	}
}

func TestPrintJSON_NilData(t *testing.T) {
	cmd := newTestCmd("json")
	out := captureStdout(t, func() {
		PrintSuccess(cmd, nil, nil)
	})

	var resp dto.CLIResponse
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	if resp.Data != nil {
		t.Errorf("expected nil data, got %v", resp.Data)
	}
}

func TestPrintJSON_SliceData(t *testing.T) {
	cmd := newTestCmd("json")
	out := captureStdout(t, func() {
		PrintSuccess(cmd, []string{"a", "b", "c"}, nil)
	})

	var resp dto.CLIResponse
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	items, ok := resp.Data.([]interface{})
	if !ok {
		t.Fatalf("expected slice data, got %T", resp.Data)
	}
	if len(items) != 3 {
		t.Errorf("expected 3 items, got %d", len(items))
	}
}

func TestPrintJSON_Indented(t *testing.T) {
	cmd := newTestCmd("json")
	out := captureStdout(t, func() {
		PrintSuccess(cmd, map[string]string{"a": "1"}, nil)
	})

	// Check indentation
	if !strings.Contains(out, "  ") {
		t.Error("expected indented JSON output")
	}
}

func TestPrintJSON_OmitsMeta(t *testing.T) {
	cmd := newTestCmd("json")
	out := captureStdout(t, func() {
		PrintSuccess(cmd, 1.23, &dto.CLIMeta{Service: "datastream", Endpoint: "/prices"})
	})

	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(out), &parsed); err != nil {
		t.Fatalf("failed to parse JSON: %v", err)
	}
	if _, exists := parsed["meta"]; exists {
		t.Fatalf("expected meta to be omitted, got: %s", out)
	}
}

func TestPrintJSON_DoesNotEscapeAngleBrackets(t *testing.T) {
	cmd := newTestCmd("json")
	out := captureStdout(t, func() {
		PrintError(cmd, "USAGE_ERROR", "Usage: flux data tokens candles <mint>", nil)
	})

	if strings.Contains(out, `\u003c`) || strings.Contains(out, `\u003e`) {
		t.Fatalf("expected angle brackets to remain unescaped, got: %s", out)
	}
	if !strings.Contains(out, "<mint>") {
		t.Fatalf("expected literal placeholder in output, got: %s", out)
	}
}

func TestPrintTable_Success(t *testing.T) {
	cmd := newTestCmd("table")
	out := captureStdout(t, func() {
		Print(cmd, &dto.CLIResponse{
			Success: true,
			Data: map[string]interface{}{
				"name":  "SOL",
				"price": 150.5,
			},
		})
	})

	if !strings.Contains(out, "KEY") || !strings.Contains(out, "VALUE") {
		t.Errorf("expected table headers, got: %s", out)
	}
	if !strings.Contains(out, "name") || !strings.Contains(out, "SOL") {
		t.Errorf("expected name/SOL in output, got: %s", out)
	}
}

func TestPrintTable_Error(t *testing.T) {
	cmd := newTestCmd("table")
	out := captureStderr(t, func() {
		Print(cmd, &dto.CLIResponse{
			Success: false,
			Error:   &dto.CLIError{Code: "ERR", Message: "something failed"},
		})
	})

	if !strings.Contains(out, "ERR") || !strings.Contains(out, "something failed") {
		t.Errorf("expected error output, got: %s", out)
	}
}

func TestPrintTable_NilData(t *testing.T) {
	cmd := newTestCmd("table")
	out := captureStdout(t, func() {
		Print(cmd, &dto.CLIResponse{
			Success: true,
			Data:    nil,
		})
	})

	if !strings.Contains(out, "No data") {
		t.Errorf("expected 'No data', got: %s", out)
	}
}

func TestPrintTable_SliceData(t *testing.T) {
	cmd := newTestCmd("table")
	out := captureStdout(t, func() {
		Print(cmd, &dto.CLIResponse{
			Success: true,
			Data: []interface{}{
				map[string]interface{}{"mint": "abc", "price": 1.5},
				map[string]interface{}{"mint": "def", "price": 2.5},
			},
		})
	})

	if !strings.Contains(out, "abc") || !strings.Contains(out, "def") {
		t.Errorf("expected mint values in table, got: %s", out)
	}
}

func TestPrint_DefaultsToJSON(t *testing.T) {
	cmd := newTestCmd("")
	out := captureStdout(t, func() {
		PrintSuccess(cmd, "test", nil)
	})

	var resp dto.CLIResponse
	if err := json.Unmarshal([]byte(out), &resp); err != nil {
		t.Fatalf("expected JSON output by default, got: %s", out)
	}
}

func TestFormatValue(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{"nil", nil, ""},
		{"string", "hello", "hello"},
		{"int", 42, "42"},
		{"float", 1.5, "1.5"},
		{"bool", true, "true"},
		{"slice", []int{1, 2}, "[1,2]"},
		{"map", map[string]int{"a": 1}, `{"a":1}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatValue(tt.input)
			if got != tt.expected {
				t.Errorf("formatValue(%v) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
